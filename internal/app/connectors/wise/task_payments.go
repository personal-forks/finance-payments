package wise

import (
	"context"
	"math/big"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/formancehq/payments/internal/app/connectors/currency"
	"github.com/formancehq/payments/internal/app/connectors/wise/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/storage"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

var (
	initiateTransferAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "initiate_transfer"))...)
	initiatePayoutAttrs   = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "initiate_payout"))...)
)

func taskInitiatePayment(logger logging.Logger, wiseClient *client.Client, transferID string) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		storageReader storage.Reader,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info("initiate payment for transfer-initiation %s", transferID)

		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)

		attrs := metric.WithAttributes(connectorAttrs...)
		var err error
		defer func() {
			if err != nil {
				ctx, cancel := contextutil.Detached(ctx)
				defer cancel()
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, attrs)
				if err := ingester.UpdateTransferInitiationStatus(ctx, transferInitiationID, models.TransferInitiationStatusFailed, err.Error(), time.Now()); err != nil {
					logger.Error("failed to update transfer initiation status: %v", err)
				}
			}
		}()

		err = ingester.UpdateTransferInitiationStatus(ctx, transferInitiationID, models.TransferInitiationStatusProcessing, "", time.Now())
		if err != nil {
			return err
		}

		var transfer *models.TransferInitiation
		transfer, err = getTransfer(ctx, storageReader, transferInitiationID, true)
		if err != nil {
			return err
		}

		attrs = initiateTransferAttrs
		if transfer.Type == models.TransferInitiationTypePayout {
			attrs = initiatePayoutAttrs
		}

		logger.Info("initiate payment between", transfer.SourceAccountID, " and %s", transfer.DestinationAccountID)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), attrs)
		}()

		if transfer.SourceAccount == nil {
			err = errors.New("missing source account")
			return err
		}

		if transfer.SourceAccount.Type == models.AccountTypeExternal {
			err = errors.New("payin not implemented: source account must be an internal account")
			return err
		}

		profileID, ok := transfer.SourceAccount.Metadata["profile_id"]
		if !ok || profileID == "" {
			err = errors.New("missing user_id in source account metadata")
			return err
		}

		var curr string
		curr, _, err = currency.GetCurrencyAndPrecisionFromAsset(transfer.Asset)
		if err != nil {
			return err
		}

		amount := big.NewFloat(0).SetInt(transfer.Amount)
		amount = amount.Quo(amount, big.NewFloat(100))

		quote, err := wiseClient.CreateQuote(profileID, curr, amount)
		if err != nil {
			return err
		}

		var connectorPaymentID uint64
		var paymentType models.PaymentType
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		switch transfer.DestinationAccount.Type {
		case models.AccountTypeInternal:
			// Transfer between internal accounts
			destinationAccount, err := strconv.ParseUint(transfer.DestinationAccount.Metadata["profile_id"], 10, 64)
			if err != nil {
				return err
			}

			var resp *client.Transfer
			resp, err = wiseClient.CreateTransfer(quote, destinationAccount, transfer.ID.Reference)
			if err != nil {
				return err
			}

			connectorPaymentID = resp.ID
			paymentType = models.PaymentTypeTransfer
		case models.AccountTypeExternal:
			// Payout to an external account

			destinationAccount, err := strconv.ParseUint(transfer.DestinationAccount.Reference, 10, 64)
			if err != nil {
				return err
			}

			var resp *client.Payout
			resp, err = wiseClient.CreatePayout(quote, destinationAccount, transfer.ID.Reference)
			if err != nil {
				return err
			}

			connectorPaymentID = resp.ID
			paymentType = models.PaymentTypePayOut
		}
		metricsRegistry.ConnectorObjects().Add(ctx, 1, attrs)

		err = ingester.UpdateTransferInitiationPaymentID(ctx, transferInitiationID, models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: strconv.FormatUint(connectorPaymentID, 10),
				Type:      paymentType,
			},
			Provider: models.ConnectorProviderWise,
		}, time.Now())
		if err != nil {
			return err
		}

		taskDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name:       "Update transfer initiation status",
			Key:        taskNameUpdatePaymentStatus,
			TransferID: transfer.ID.String(),
			Attempt:    1,
		})
		if err != nil {
			return err
		}

		ctx, _ = contextutil.DetachedWithTimeout(ctx, 10*time.Second)
		err = scheduler.Schedule(ctx, taskDescriptor, models.TaskSchedulerOptions{
			// We want to polling every c.cfg.PollingPeriod.Duration seconds the users
			// and their transactions.
			ScheduleOption: models.OPTIONS_RUN_NOW,
			// No need to restart this task, since the connector is not existing or
			// was uninstalled previously, the task does not exists in the database
			Restart: true,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			return err
		}

		return nil
	}
}

var (
	updateTransferAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "update_transfer"))...)
	updatePayoutAttrs   = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "update_payout"))...)
)

func taskUpdatePaymentStatus(
	logger logging.Logger,
	wiseClient *client.Client,
	transferID string,
	attempt int,
) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		storageReader storage.Reader,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		transferInitiationID := models.MustTransferInitiationIDFromString(transferID)
		transfer, err := getTransfer(ctx, storageReader, transferInitiationID, false)
		if err != nil {
			return err
		}
		logger.Info("attempt: ", attempt, " fetching status of ", transfer.PaymentID)

		attrs := updateTransferAttrs
		if transfer.Type == models.TransferInitiationTypePayout {
			attrs = updatePayoutAttrs
		}

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), attrs)
		}()

		defer func() {
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, attrs)
			}
		}()

		var status string
		switch transfer.Type {
		case models.TransferInitiationTypeTransfer:
			var resp *client.Transfer
			resp, err = wiseClient.GetTransfer(ctx, transfer.PaymentID.Reference)
			if err != nil {
				return err
			}

			status = resp.Status
		case models.TransferInitiationTypePayout:
			var resp *client.Payout
			resp, err = wiseClient.GetPayout(ctx, transfer.PaymentID.Reference)
			if err != nil {
				return err
			}

			status = resp.Status
		}

		switch status {
		case "incoming_payment_waiting",
			"incoming_payment_initiated",
			"processing",
			"funds_converted",
			"bounced_back",
			"unknown":
			taskDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:       "Update transfer initiation status",
				Key:        taskNameUpdatePaymentStatus,
				TransferID: transfer.ID.String(),
				Attempt:    attempt + 1,
			})
			if err != nil {
				return err
			}

			err = scheduler.Schedule(ctx, taskDescriptor, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_IN_DURATION,
				Duration:       2 * time.Minute,
				Restart:        true,
			})
			if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				return err
			}
		case "outgoing_payment_sent", "funds_refunded":
			err = ingester.UpdateTransferInitiationStatus(ctx, transferInitiationID, models.TransferInitiationStatusProcessed, "", time.Now())
			if err != nil {
				return err
			}

			return nil
		case "charged_back", "cancelled":
			err = ingester.UpdateTransferInitiationStatus(ctx, transferInitiationID, models.TransferInitiationStatusFailed, "", time.Now())
			if err != nil {
				return err
			}

			return nil
		}

		return nil
	}
}

func getTransfer(
	ctx context.Context,
	reader storage.Reader,
	transferID models.TransferInitiationID,
	expand bool,
) (*models.TransferInitiation, error) {
	transfer, err := reader.ReadTransferInitiation(ctx, transferID)
	if err != nil {
		return nil, err
	}

	if expand {
		if transfer.SourceAccountID.Reference != "" {
			sourceAccount, err := reader.GetAccount(ctx, transfer.SourceAccountID.String())
			if err != nil {
				return nil, err
			}
			transfer.SourceAccount = sourceAccount
		}

		destinationAccount, err := reader.GetAccount(ctx, transfer.DestinationAccountID.String())
		if err != nil {
			return nil, err
		}
		transfer.DestinationAccount = destinationAccount
	}

	return transfer, nil
}

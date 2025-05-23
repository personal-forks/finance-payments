package services

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/go-libs/v3/bun/bunpaginate"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
	"github.com/pkg/errors"
)

func (s *Service) PaymentInitiationsApprove(ctx context.Context, id models.PaymentInitiationID, waitResult bool) (models.Task, error) {
	cursor, err := s.storage.PaymentInitiationAdjustmentsList(
		ctx,
		id,
		storage.NewListPaymentInitiationAdjustmentsQuery(
			bunpaginate.NewPaginatedQueryOptions(storage.PaymentInitiationAdjustmentsQuery{}).
				WithPageSize(1),
		),
	)
	if err != nil {
		return models.Task{}, newStorageError(err, "cannot list payment initiation adjustments")
	}

	if len(cursor.Data) == 0 {
		return models.Task{}, errors.New("payment initiation adjustments not found")
	}

	lastAdjustment := cursor.Data[0]

	if lastAdjustment.Status != models.PAYMENT_INITIATION_ADJUSTMENT_STATUS_WAITING_FOR_VALIDATION {
		return models.Task{}, fmt.Errorf("cannot approve an already approved payment initiation: %w", ErrValidation)
	}

	pi, err := s.storage.PaymentInitiationsGet(ctx, id)
	if err != nil {
		return models.Task{}, newStorageError(err, "cannot get payment initiation")
	}

	if !pi.ScheduledAt.IsZero() && pi.ScheduledAt.After(time.Now()) {
		// In any case, if the payment initiation is scheduled for the future,
		// we do not want to wait for the results
		waitResult = false
	}

	switch pi.Type {
	case models.PAYMENT_INITIATION_TYPE_TRANSFER:
		task, err := s.engine.CreateTransfer(ctx, pi.ID, 1, waitResult)
		if err != nil {
			return models.Task{}, handleEngineErrors(err)
		}
		return task, nil
	case models.PAYMENT_INITIATION_TYPE_PAYOUT:
		task, err := s.engine.CreatePayout(ctx, pi.ID, 1, waitResult)
		if err != nil {
			return models.Task{}, handleEngineErrors(err)
		}
		return task, nil
	}

	return models.Task{}, nil
}

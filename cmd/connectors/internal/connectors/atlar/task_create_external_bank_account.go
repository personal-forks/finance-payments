package atlar

import (
	"context"
	"errors"
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func CreateExternalBankAccountTask(config Config, client *client.Client, newExternalBankAccount *models.BankAccount) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
	) error {
		err := validateExternalBankAccount(newExternalBankAccount)
		if err != nil {
			return err
		}

		externalAccountID, err := createExternalBankAccount(ctx, client, newExternalBankAccount)
		if err != nil {
			return err
		}
		if externalAccountID == nil {
			return errors.New("no external account id returned")
		}

		err = ingestExternalAccountFromAtlar(
			ctx,
			logger,
			connectorID,
			ingester,
			client,
			*externalAccountID,
		)
		if err != nil {
			return err
		}

		return nil
	}
}

// TODO: validation (also metadata) needs to return a 400
func validateExternalBankAccount(newExternalBankAccount *models.BankAccount) error {
	_, err := ExtractNamespacedMetadata(newExternalBankAccount.Metadata, "owner/name")
	if err != nil {
		return fmt.Errorf("required metadata field %sowner/name is missing", atlarMetadataSpecNamespace)
	}
	ownerType, err := ExtractNamespacedMetadata(newExternalBankAccount.Metadata, "owner/type")
	if err != nil {
		return fmt.Errorf("required metadata field %sowner/type is missing", atlarMetadataSpecNamespace)
	}
	if *ownerType != "INDIVIDUAL" && *ownerType != "COMPANY" {
		return fmt.Errorf("metadata field %sowner/type needs to be one of [ INDIVIDUAL COMPANY ]", atlarMetadataSpecNamespace)
	}

	return nil
}

func createExternalBankAccount(ctx context.Context, client *client.Client, newExternalBankAccount *models.BankAccount) (*string, error) {
	return client.CreateCounterParty(ctx, newExternalBankAccount)
}

func ingestExternalAccountFromAtlar(
	ctx context.Context,
	logger logging.Logger,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	client *client.Client,
	externalAccountID string,
) error {
	accountsBatch := ingestion.AccountBatch{}

	externalAccountResponse, err := client.GetV1ExternalAccountsID(ctx, externalAccountID)
	if err != nil {
		return err
	}

	counterpartyResponse, err := client.GetV1CounterpartiesID(ctx, externalAccountResponse.Payload.CounterpartyID)
	if err != nil {
		return err
	}

	newAccount, err := ExternalAccountFromAtlarData(connectorID, externalAccountResponse.Payload, counterpartyResponse.Payload)
	if err != nil {
		return err
	}
	logger.WithContext(ctx).Info("Got external Account from Atlar", newAccount)

	accountsBatch = append(accountsBatch, newAccount)

	err = ingester.IngestAccounts(ctx, accountsBatch)
	if err != nil {
		return err
	}

	return nil
}

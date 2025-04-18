package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/formancehq/payments/internal/connectors/metrics"
	errorsutils "github.com/formancehq/payments/internal/utils/errors"
)

type PayoutRequest struct {
	BeneficiaryID   string      `json:"beneficiary_id"`
	Currency        string      `json:"currency"`
	Amount          json.Number `json:"amount"`
	Reference       string      `json:"reference"`
	Reason          string      `json:"reason"`
	UniqueRequestID string      `json:"unique_request_id"`
}

func (pr *PayoutRequest) ToFormData() url.Values {
	form := url.Values{}
	form.Set("beneficiary_id", pr.BeneficiaryID)
	form.Set("currency", pr.Currency)
	form.Set("amount", pr.Amount.String())
	form.Set("reason", pr.Reason)
	form.Set("payment_type", "regular")
	form.Set("reference", pr.Reference)
	if pr.UniqueRequestID != "" {
		form.Set("unique_request_id", pr.UniqueRequestID)
	}

	return form
}

type PayoutResponse struct {
	ID               string      `json:"id"`
	Amount           json.Number `json:"amount"`
	BeneficiaryID    string      `json:"beneficiary_id"`
	Currency         string      `json:"currency"`
	Reference        string      `json:"reference"`
	Status           string      `json:"status"`
	Reason           string      `json:"reason"`
	CreatorContactID string      `json:"creator_contact_id"`
	PaymentType      string      `json:"payment_type"`
	TransferredAt    string      `json:"transferred_at"`
	PaymentDate      string      `json:"payment_date"`
	FailureReason    string      `json:"failure_reason"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	UniqueRequestID  string      `json:"unique_request_id"`
}

func (c *client) InitiatePayout(ctx context.Context, payoutRequest *PayoutRequest) (*PayoutResponse, error) {
	ctx = context.WithValue(ctx, metrics.MetricOperationContextKey, "initiate_payout")

	if err := c.ensureLogin(ctx); err != nil {
		return nil, err
	}

	form := payoutRequest.ToFormData()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.buildEndpoint("v2/payments/create"), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var payoutResponse PayoutResponse
	var errRes currencyCloudError
	_, err = c.httpClient.Do(ctx, req, &payoutResponse, &errRes)
	if err != nil {
		return nil, errorsutils.NewWrappedError(
			fmt.Errorf("failed to create payout: %v", errRes.Error()),
			err,
		)
	}

	return &payoutResponse, nil
}

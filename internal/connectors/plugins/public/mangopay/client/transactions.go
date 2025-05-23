package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/formancehq/payments/internal/connectors/metrics"
	errorsutils "github.com/formancehq/payments/internal/utils/errors"
)

type Payment struct {
	Id               string `json:"Id"`
	Tag              string `json:"Tag"`
	CreationDate     int64  `json:"CreationDate"`
	AuthorId         string `json:"AuthorId"`
	CreditedUserId   string `json:"CreditedUserId"`
	DebitedFunds     Funds  `json:"DebitedFunds"`
	CreditedFunds    Funds  `json:"CreditedFunds"`
	Fees             Funds  `json:"Fees"`
	Status           string `json:"Status"`
	ResultCode       string `json:"ResultCode"`
	ResultMessage    string `json:"ResultMessage"`
	ExecutionDate    int64  `json:"ExecutionDate"`
	Type             string `json:"Type"`
	Nature           string `json:"Nature"`
	CreditedWalletID string `json:"CreditedWalletId"`
	DebitedWalletID  string `json:"DebitedWalletId"`
}

func (c *client) GetTransactions(ctx context.Context, walletsID string, page, pageSize int, afterCreatedAt time.Time) ([]Payment, error) {
	ctx = context.WithValue(ctx, metrics.MetricOperationContextKey, "list_transactions")

	endpoint := fmt.Sprintf("%s/v2.01/%s/wallets/%s/transactions", c.endpoint, c.clientID, walletsID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

	q := req.URL.Query()
	q.Add("per_page", strconv.Itoa(pageSize))
	q.Add("page", fmt.Sprint(page))
	q.Add("Sort", "CreationDate:ASC")
	if !afterCreatedAt.IsZero() {
		q.Add("AfterDate", strconv.FormatInt(afterCreatedAt.UTC().Unix(), 10))
	}
	req.URL.RawQuery = q.Encode()

	var payments []Payment
	statusCode, err := c.httpClient.Do(ctx, req, &payments, nil)
	if err != nil {
		return nil, errorsutils.NewWrappedError(
			fmt.Errorf("failed to get transactions: status code %d", statusCode),
			err,
		)
	}
	return payments, nil
}

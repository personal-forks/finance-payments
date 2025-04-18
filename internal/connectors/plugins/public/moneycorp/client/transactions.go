package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/formancehq/payments/internal/connectors/metrics"
	errorsutils "github.com/formancehq/payments/internal/utils/errors"
)

type transactionsResponse struct {
	Transactions []*Transaction `json:"data"`
}

type fetchTransactionRequest struct {
	Data struct {
		Attributes struct {
			TransactionDateTimeFrom string `json:"transactionDateTimeFrom"`
		} `json:"attributes"`
	} `json:"data"`
}

type Transaction struct {
	ID            string                `json:"id"`
	Type          string                `json:"type"`
	Attributes    TransactionAttributes `json:"attributes"`
	Relationships RelationShips         `json:"relationships"`
}

type Data struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type RelationShips struct {
	Data Data `json:"data"`
}

type TransactionAttributes struct {
	AccountID            int32       `json:"accountId"`
	CreatedAt            string      `json:"createdAt"`
	Currency             string      `json:"transactionCurrency"`
	Amount               json.Number `json:"transactionAmount"`
	Direction            string      `json:"transactionDirection"`
	Type                 string      `json:"transactionType"`
	ClientReference      string      `json:"clientReference"`
	TransactionReference string      `json:"transactionReference"`
}

func (c *client) GetTransactions(ctx context.Context, accountID string, page, pageSize int, lastCreatedAt time.Time) ([]*Transaction, error) {
	ctx = context.WithValue(ctx, metrics.MetricOperationContextKey, "list_transactions")

	var body io.Reader
	if !lastCreatedAt.IsZero() {
		reqBody := fetchTransactionRequest{
			Data: struct {
				Attributes struct {
					TransactionDateTimeFrom string "json:\"transactionDateTimeFrom\""
				} "json:\"attributes\""
			}{
				Attributes: struct {
					TransactionDateTimeFrom string "json:\"transactionDateTimeFrom\""
				}{
					TransactionDateTimeFrom: lastCreatedAt.Format("2006-01-02T15:04:05.999999999"),
				},
			},
		}

		raw, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal transfer request: %w", err)
		}

		body = bytes.NewBuffer(raw)
	} else {
		body = http.NoBody
	}

	endpoint := fmt.Sprintf("%s/accounts/%s/transactions/find", c.endpoint, accountID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactions request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("page[size]", strconv.Itoa(pageSize))
	q.Add("page[number]", fmt.Sprint(page))
	q.Add("sortBy", "createdAt.asc")
	req.URL.RawQuery = q.Encode()

	transactions := transactionsResponse{Transactions: make([]*Transaction, 0)}
	var errRes moneycorpErrors
	_, err = c.httpClient.Do(ctx, req, &transactions, &errRes)
	if err != nil {
		return nil, errorsutils.NewWrappedError(
			fmt.Errorf("failed to get transactions: %v", errRes.Error()),
			err,
		)
	}

	return transactions.Transactions, nil
}

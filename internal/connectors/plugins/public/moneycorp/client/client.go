package client

import (
	"net/http"
	"strings"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	httpClient httpwrapper.Client
	endpoint   string
}

func New(clientID, apiKey, endpoint string) (*Client, error) {
	config := &httpwrapper.Config{
		Transport: &apiTransport{
			clientID:   clientID,
			apiKey:     apiKey,
			endpoint:   endpoint,
			underlying: otelhttp.NewTransport(http.DefaultTransport),
		},
		HttpErrorCheckerFn: func(statusCode int) error {
			if statusCode == http.StatusNotFound {
				return nil
			} else if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
				return httpwrapper.ErrStatusCodeClientError
			} else if statusCode >= http.StatusInternalServerError {
				return httpwrapper.ErrStatusCodeServerError
			}
			return nil

		},
	}
	endpoint = strings.TrimSuffix(endpoint, "/")

	httpClient, err := httpwrapper.NewClient(config)
	c := &Client{
		httpClient: httpClient,
		endpoint:   endpoint,
	}
	return c, err
}

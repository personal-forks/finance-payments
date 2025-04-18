package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
	"github.com/formancehq/payments/internal/connectors/metrics"
	"github.com/formancehq/payments/internal/models"
	errorsutils "github.com/formancehq/payments/internal/utils/errors"
)

//go:generate mockgen -source client.go -destination client_generated.go -package client . Client
type Client interface {
	GetAccounts(ctx context.Context, page int, pageSize int, fromOpeningDate time.Time) ([]Account, error)
	GetAccount(ctx context.Context, accountID string) (*Account, error)
	GetPayments(ctx context.Context, page int, pageSize int) ([]Payment, error)
	GetPayment(ctx context.Context, paymentID string) (*Payment, error)
	GetPaymentStatus(ctx context.Context, paymentID string) (*StatusResponse, error)
	InitiateTransferOrPayouts(ctx context.Context, transferRequest *PaymentRequest) (*PaymentResponse, error)
}

type client struct {
	httpClient httpwrapper.Client

	username string
	password string

	endpoint              string
	authorizationEndpoint string

	accessToken          string
	accessTokenExpiresAt time.Time
}

func New(
	connectorName string,
	username, password,
	endpoint, authorizationEndpoint,
	uCertificate, uCertificateKey string,
) (Client, error) {
	cert, err := tls.X509KeyPair([]byte(uCertificate), []byte(uCertificateKey))
	if err != nil {
		return nil, errorsutils.NewWrappedError(
			fmt.Errorf("failed to load user certificate: %w", err),
			models.ErrInvalidConfig,
		)
	}

	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	config := &httpwrapper.Config{
		Transport: metrics.NewTransport(connectorName, metrics.TransportOpts{Transport: tr}),
	}

	c := &client{
		httpClient: httpwrapper.NewClient(config),

		username:              username,
		password:              password,
		endpoint:              endpoint,
		authorizationEndpoint: authorizationEndpoint,
	}

	return c, nil
}

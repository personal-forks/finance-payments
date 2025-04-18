package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/formancehq/payments/internal/connectors/engine"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestPaymentInitiationsCreate(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storage.NewMockStorage(ctrl)
	eng := engine.NewMockEngine(ctrl)

	s := New(store, eng, false)

	piWithoutScheduledAt := models.PaymentInitiation{
		Type: models.PAYMENT_INITIATION_TYPE_TRANSFER,
	}
	piWithScheduledAt := models.PaymentInitiation{
		Type:        models.PAYMENT_INITIATION_TYPE_PAYOUT,
		ScheduledAt: time.Now(),
	}

	tests := []struct {
		name                          string
		sendToPSP                     bool
		pi                            models.PaymentInitiation
		engineCreatePaymentInitiation error
		engineErr                     error
		expectedPIError               error
		expectedEngineError           error
		typedError                    bool
	}{
		{
			name:      "success without scheduled at and transfer",
			pi:        piWithoutScheduledAt,
			sendToPSP: true,
		},
		{
			name:      "success with scheduled at and payout",
			pi:        piWithScheduledAt,
			sendToPSP: true,
		},
		{
			name:      "success without sending to PSP",
			pi:        piWithoutScheduledAt,
			sendToPSP: false,
		},
		{
			name:                "not found error",
			sendToPSP:           true,
			pi:                  piWithoutScheduledAt,
			engineErr:           engine.ErrNotFound,
			expectedEngineError: ErrNotFound,
			typedError:          true,
		},
		{
			name:                "other error",
			sendToPSP:           true,
			pi:                  piWithoutScheduledAt,
			engineErr:           fmt.Errorf("error"),
			expectedEngineError: fmt.Errorf("error"),
		},
		{
			name:                          "not found error",
			pi:                            piWithoutScheduledAt,
			engineCreatePaymentInitiation: engine.ErrNotFound,
			expectedEngineError:           ErrNotFound,
			typedError:                    true,
		},
		{
			name:                          "other error",
			pi:                            piWithoutScheduledAt,
			engineCreatePaymentInitiation: fmt.Errorf("error"),
			expectedEngineError:           fmt.Errorf("error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			eng.EXPECT().CreateFormancePaymentInitiation(gomock.Any(), test.pi, gomock.Any()).Return(test.engineCreatePaymentInitiation)
			if test.sendToPSP {
				switch test.pi.Type {
				case models.PAYMENT_INITIATION_TYPE_TRANSFER:
					eng.EXPECT().CreateTransfer(gomock.Any(), models.PaymentInitiationID{}, 1, false).Return(models.Task{}, test.engineErr)
				case models.PAYMENT_INITIATION_TYPE_PAYOUT:
					eng.EXPECT().CreatePayout(gomock.Any(), models.PaymentInitiationID{}, 1, false).Return(models.Task{}, test.engineErr)
				}
			}

			_, err := s.PaymentInitiationsCreate(context.Background(), test.pi, test.sendToPSP, false)
			switch {
			case test.expectedPIError == nil && test.expectedEngineError == nil:
				require.NoError(t, err)
			case test.expectedPIError != nil:
				if test.typedError {
					require.ErrorIs(t, err, test.expectedPIError)
				} else {
					require.Equal(t, test.expectedPIError.Error(), err.Error())
				}
			case test.expectedEngineError != nil:
				if test.typedError {
					require.ErrorIs(t, err, test.expectedEngineError)
				} else {
					require.Equal(t, test.expectedEngineError, err)
				}
			}
		})
	}
}

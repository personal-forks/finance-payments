package v2

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/go-libs/v3/api"
	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/api/common"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

func paymentsGet(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v2_paymentsGet")
		defer span.End()

		span.SetAttributes(attribute.String("paymentID", paymentID(r)))
		id, err := models.PaymentIDFromString(paymentID(r))
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		payment, err := backend.PaymentsGet(ctx, id)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		data := PaymentResponse{
			ID:            payment.ID.String(),
			Reference:     payment.Reference,
			Type:          payment.Type.String(),
			Provider:      toV2Provider(payment.ConnectorID.Provider),
			ConnectorID:   payment.ConnectorID.String(),
			Status:        payment.Status.String(),
			Amount:        payment.Amount,
			InitialAmount: payment.InitialAmount,
			Scheme:        toV2PaymentScheme(payment.Scheme),
			Asset:         payment.Asset,
			CreatedAt:     payment.CreatedAt,
			Metadata:      payment.Metadata,
		}

		if payment.SourceAccountID != nil {
			data.SourceAccountID = payment.SourceAccountID.String()
		}

		if payment.DestinationAccountID != nil {
			data.DestinationAccountID = payment.DestinationAccountID.String()
		}

		data.Adjustments = make([]paymentAdjustment, len(payment.Adjustments))
		for i := range payment.Adjustments {
			data.Adjustments[i] = paymentAdjustment{
				Reference: payment.Adjustments[i].ID.Reference,
				CreatedAt: payment.Adjustments[i].CreatedAt,
				Status:    payment.Adjustments[i].Status.String(),
				Amount:    payment.Adjustments[i].Amount,
				Raw:       payment.Adjustments[i].Raw,
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[PaymentResponse]{
			Data: &data,
		})
		if err != nil {
			otel.RecordError(span, err)
			common.InternalServerError(w, r, err)
			return
		}
	}
}

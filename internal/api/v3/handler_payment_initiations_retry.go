package v3

import (
	"net/http"

	"github.com/formancehq/go-libs/v3/api"
	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

type PaymentInitiationsRetryResponse struct {
	TaskID string `json:"taskID"`
}

func paymentInitiationsRetry(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v3_paymentInitiationsRetry")
		defer span.End()

		span.SetAttributes(attribute.String("paymentInitiationID", paymentInitiationID(r)))
		id, err := models.PaymentInitiationIDFromString(paymentInitiationID(r))
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		task, err := backend.PaymentInitiationsRetry(ctx, id, false)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		api.Accepted(w, PaymentInitiationsRetryResponse{
			TaskID: task.ID.String(),
		})
	}
}

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

func accountsGet(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v2_accountsGet")
		defer span.End()

		span.SetAttributes(attribute.String("accountID", accountID(r)))
		id, err := models.AccountIDFromString(accountID(r))
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		account, err := backend.AccountsGet(ctx, id)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		data := &accountResponse{
			ID:          account.ID.String(),
			Reference:   account.Reference,
			CreatedAt:   account.CreatedAt,
			ConnectorID: account.ConnectorID.String(),
			Provider:    toV2Provider(account.ConnectorID.Provider),
			Type:        string(account.Type),
			Metadata:    account.Metadata,
			Raw:         account.Raw,
		}

		if account.DefaultAsset != nil {
			data.DefaultCurrency = *account.DefaultAsset
			data.DefaultAsset = *account.DefaultAsset
		}

		if account.Name != nil {
			data.AccountName = *account.Name
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[accountResponse]{
			Data: data,
		})
		if err != nil {
			otel.RecordError(span, err)
			common.InternalServerError(w, r, err)
			return
		}
	}
}

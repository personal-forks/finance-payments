package v2

import (
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/go-libs/v3/api"
	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/api/common"
	"github.com/formancehq/payments/internal/otel"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
)

// NOTE: in order to maintain previous version compatibility, we need to keep the
// same response structure as the previous version of the API
type poolBalancesResponse struct {
	Balances []*poolBalanceResponse `json:"balances"`
}

type poolBalanceResponse struct {
	Amount *big.Int `json:"amount"`
	Asset  string   `json:"asset"`
}

func poolsBalancesAt(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v2_poolsBalancesAt")
		defer span.End()

		span.SetAttributes(attribute.String("poolID", poolID(r)))
		id, err := uuid.Parse(poolID(r))
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		span.SetAttributes(attribute.String("at", r.URL.Query().Get("at")))
		atTime := r.URL.Query().Get("at")
		if atTime == "" {
			otel.RecordError(span, errors.New("missing atTime"))
			api.BadRequest(w, ErrValidation, errors.New("missing atTime"))
			return
		}

		at, err := time.Parse(time.RFC3339, atTime)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, errors.Wrap(err, "invalid atTime"))
			return
		}

		balances, err := backend.PoolsBalancesAt(ctx, id, at)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		data := &poolBalancesResponse{
			Balances: make([]*poolBalanceResponse, len(balances)),
		}

		for i := range balances {
			data.Balances[i] = &poolBalanceResponse{
				Amount: balances[i].Amount,
				Asset:  balances[i].Asset,
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[poolBalancesResponse]{
			Data: data,
		})
		if err != nil {
			otel.RecordError(span, err)
			common.InternalServerError(w, r, err)
			return
		}
	}
}

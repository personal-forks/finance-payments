package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/formancehq/go-libs/v2/pointer"
)

type PaymentInitiationAdjustment struct {
	// Unique ID
	ID PaymentInitiationAdjustmentID `json:"id"`

	// Related Payment Initiation ID
	PaymentInitiationID PaymentInitiationID `json:"paymentInitiationID"`
	// Creation date of the adjustment
	CreatedAt time.Time `json:"createdAt"`
	// Last status of the adjustment
	Status PaymentInitiationAdjustmentStatus `json:"status"`
	// Error description if we had one
	Error error `json:"error"`
	// Additional metadata
	Metadata map[string]string `json:"metadata"`
}

func (pia PaymentInitiationAdjustment) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID                  string                            `json:"id"`
		PaymentInitiationID string                            `json:"paymentInitiationID"`
		CreatedAt           time.Time                         `json:"createdAt"`
		Status              PaymentInitiationAdjustmentStatus `json:"status"`
		Error               *string                           `json:"error,omitempty"`
		Metadata            map[string]string                 `json:"metadata"`
	}{
		ID:                  pia.ID.String(),
		PaymentInitiationID: pia.PaymentInitiationID.String(),
		CreatedAt:           pia.CreatedAt,
		Status:              pia.Status,
		Error: func() *string {
			if pia.Error == nil {
				return nil
			}

			return pointer.For(pia.Error.Error())
		}(),
		Metadata: pia.Metadata,
	})
}

func (pia *PaymentInitiationAdjustment) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID                  string                            `json:"id"`
		PaymentInitiationID string                            `json:"paymentInitiationID"`
		CreatedAt           time.Time                         `json:"createdAt"`
		Status              PaymentInitiationAdjustmentStatus `json:"status"`
		Error               *string                           `json:"error"`
		Metadata            map[string]string                 `json:"metadata"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	id, err := PaymentInitiationAdjustmentIDFromString(aux.ID)
	if err != nil {
		return err
	}

	piID, err := PaymentInitiationIDFromString(aux.PaymentInitiationID)
	if err != nil {
		return err
	}

	pia.ID = id
	pia.PaymentInitiationID = piID
	pia.CreatedAt = aux.CreatedAt
	pia.Status = aux.Status
	if aux.Error != nil {
		pia.Error = errors.New(*aux.Error)
	}
	pia.Metadata = aux.Metadata

	return nil
}

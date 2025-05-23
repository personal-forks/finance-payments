// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package components

import (
	"encoding/json"
	"fmt"
	"github.com/formancehq/payments/pkg/client/internal/utils"
	"math/big"
	"time"
)

type TransferInitiationType string

const (
	TransferInitiationTypeTransfer TransferInitiationType = "TRANSFER"
	TransferInitiationTypePayout   TransferInitiationType = "PAYOUT"
)

func (e TransferInitiationType) ToPointer() *TransferInitiationType {
	return &e
}
func (e *TransferInitiationType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "TRANSFER":
		fallthrough
	case "PAYOUT":
		*e = TransferInitiationType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for TransferInitiationType: %v", v)
	}
}

type TransferInitiation struct {
	ID                   string                         `json:"id"`
	Reference            string                         `json:"reference"`
	CreatedAt            time.Time                      `json:"createdAt"`
	ScheduledAt          time.Time                      `json:"scheduledAt"`
	Description          string                         `json:"description"`
	SourceAccountID      string                         `json:"sourceAccountID"`
	DestinationAccountID string                         `json:"destinationAccountID"`
	ConnectorID          string                         `json:"connectorID"`
	Provider             *string                        `json:"provider"`
	Type                 TransferInitiationType         `json:"type"`
	Amount               *big.Int                       `json:"amount"`
	InitialAmount        *big.Int                       `json:"initialAmount"`
	Asset                string                         `json:"asset"`
	Status               TransferInitiationStatus       `json:"status"`
	Error                *string                        `json:"error,omitempty"`
	Metadata             map[string]string              `json:"metadata,omitempty"`
	RelatedPayments      []TransferInitiationPayments   `json:"relatedPayments,omitempty"`
	RelatedAdjustments   []TransferInitiationAdjusments `json:"relatedAdjustments,omitempty"`
}

func (t TransferInitiation) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(t, "", false)
}

func (t *TransferInitiation) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &t, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *TransferInitiation) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *TransferInitiation) GetReference() string {
	if o == nil {
		return ""
	}
	return o.Reference
}

func (o *TransferInitiation) GetCreatedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.CreatedAt
}

func (o *TransferInitiation) GetScheduledAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.ScheduledAt
}

func (o *TransferInitiation) GetDescription() string {
	if o == nil {
		return ""
	}
	return o.Description
}

func (o *TransferInitiation) GetSourceAccountID() string {
	if o == nil {
		return ""
	}
	return o.SourceAccountID
}

func (o *TransferInitiation) GetDestinationAccountID() string {
	if o == nil {
		return ""
	}
	return o.DestinationAccountID
}

func (o *TransferInitiation) GetConnectorID() string {
	if o == nil {
		return ""
	}
	return o.ConnectorID
}

func (o *TransferInitiation) GetProvider() *string {
	if o == nil {
		return nil
	}
	return o.Provider
}

func (o *TransferInitiation) GetType() TransferInitiationType {
	if o == nil {
		return TransferInitiationType("")
	}
	return o.Type
}

func (o *TransferInitiation) GetAmount() *big.Int {
	if o == nil {
		return big.NewInt(0)
	}
	return o.Amount
}

func (o *TransferInitiation) GetInitialAmount() *big.Int {
	if o == nil {
		return big.NewInt(0)
	}
	return o.InitialAmount
}

func (o *TransferInitiation) GetAsset() string {
	if o == nil {
		return ""
	}
	return o.Asset
}

func (o *TransferInitiation) GetStatus() TransferInitiationStatus {
	if o == nil {
		return TransferInitiationStatus("")
	}
	return o.Status
}

func (o *TransferInitiation) GetError() *string {
	if o == nil {
		return nil
	}
	return o.Error
}

func (o *TransferInitiation) GetMetadata() map[string]string {
	if o == nil {
		return nil
	}
	return o.Metadata
}

func (o *TransferInitiation) GetRelatedPayments() []TransferInitiationPayments {
	if o == nil {
		return nil
	}
	return o.RelatedPayments
}

func (o *TransferInitiation) GetRelatedAdjustments() []TransferInitiationAdjusments {
	if o == nil {
		return nil
	}
	return o.RelatedAdjustments
}

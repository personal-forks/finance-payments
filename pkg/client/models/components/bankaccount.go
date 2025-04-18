// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package components

import (
	"github.com/formancehq/payments/pkg/client/internal/utils"
	"time"
)

type BankAccount struct {
	ID              string                       `json:"id"`
	Name            string                       `json:"name"`
	CreatedAt       time.Time                    `json:"createdAt"`
	Country         string                       `json:"country"`
	ConnectorID     *string                      `json:"connectorID,omitempty"`
	AccountID       *string                      `json:"accountID,omitempty"`
	Provider        *string                      `json:"provider,omitempty"`
	Iban            *string                      `json:"iban,omitempty"`
	AccountNumber   *string                      `json:"accountNumber,omitempty"`
	SwiftBicCode    *string                      `json:"swiftBicCode,omitempty"`
	RelatedAccounts []BankAccountRelatedAccounts `json:"relatedAccounts,omitempty"`
	Metadata        map[string]string            `json:"metadata,omitempty"`
}

func (b BankAccount) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(b, "", false)
}

func (b *BankAccount) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &b, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *BankAccount) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *BankAccount) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *BankAccount) GetCreatedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.CreatedAt
}

func (o *BankAccount) GetCountry() string {
	if o == nil {
		return ""
	}
	return o.Country
}

func (o *BankAccount) GetConnectorID() *string {
	if o == nil {
		return nil
	}
	return o.ConnectorID
}

func (o *BankAccount) GetAccountID() *string {
	if o == nil {
		return nil
	}
	return o.AccountID
}

func (o *BankAccount) GetProvider() *string {
	if o == nil {
		return nil
	}
	return o.Provider
}

func (o *BankAccount) GetIban() *string {
	if o == nil {
		return nil
	}
	return o.Iban
}

func (o *BankAccount) GetAccountNumber() *string {
	if o == nil {
		return nil
	}
	return o.AccountNumber
}

func (o *BankAccount) GetSwiftBicCode() *string {
	if o == nil {
		return nil
	}
	return o.SwiftBicCode
}

func (o *BankAccount) GetRelatedAccounts() []BankAccountRelatedAccounts {
	if o == nil {
		return nil
	}
	return o.RelatedAccounts
}

func (o *BankAccount) GetMetadata() map[string]string {
	if o == nil {
		return nil
	}
	return o.Metadata
}

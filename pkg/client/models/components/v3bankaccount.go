// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package components

import (
	"github.com/formancehq/payments/pkg/client/internal/utils"
	"time"
)

type V3BankAccount struct {
	ID              string                        `json:"id"`
	CreatedAt       time.Time                     `json:"createdAt"`
	Name            string                        `json:"name"`
	AccountNumber   *string                       `json:"accountNumber,omitempty"`
	Iban            *string                       `json:"iban,omitempty"`
	SwiftBicCode    *string                       `json:"swiftBicCode,omitempty"`
	Country         *string                       `json:"country,omitempty"`
	Metadata        map[string]string             `json:"metadata,omitempty"`
	RelatedAccounts []V3BankAccountRelatedAccount `json:"relatedAccounts,omitempty"`
}

func (v V3BankAccount) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(v, "", false)
}

func (v *V3BankAccount) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &v, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *V3BankAccount) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *V3BankAccount) GetCreatedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.CreatedAt
}

func (o *V3BankAccount) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *V3BankAccount) GetAccountNumber() *string {
	if o == nil {
		return nil
	}
	return o.AccountNumber
}

func (o *V3BankAccount) GetIban() *string {
	if o == nil {
		return nil
	}
	return o.Iban
}

func (o *V3BankAccount) GetSwiftBicCode() *string {
	if o == nil {
		return nil
	}
	return o.SwiftBicCode
}

func (o *V3BankAccount) GetCountry() *string {
	if o == nil {
		return nil
	}
	return o.Country
}

func (o *V3BankAccount) GetMetadata() map[string]string {
	if o == nil {
		return nil
	}
	return o.Metadata
}

func (o *V3BankAccount) GetRelatedAccounts() []V3BankAccountRelatedAccount {
	if o == nil {
		return nil
	}
	return o.RelatedAccounts
}

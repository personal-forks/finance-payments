// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package components

import (
	"encoding/json"
	"fmt"
)

type PaymentScheme string

const (
	PaymentSchemeUnknown    PaymentScheme = "unknown"
	PaymentSchemeOther      PaymentScheme = "other"
	PaymentSchemeVisa       PaymentScheme = "visa"
	PaymentSchemeMastercard PaymentScheme = "mastercard"
	PaymentSchemeAmex       PaymentScheme = "amex"
	PaymentSchemeDiners     PaymentScheme = "diners"
	PaymentSchemeDiscover   PaymentScheme = "discover"
	PaymentSchemeJcb        PaymentScheme = "jcb"
	PaymentSchemeUnionpay   PaymentScheme = "unionpay"
	PaymentSchemeAlipay     PaymentScheme = "alipay"
	PaymentSchemeCup        PaymentScheme = "cup"
	PaymentSchemeSepaDebit  PaymentScheme = "sepa debit"
	PaymentSchemeSepaCredit PaymentScheme = "sepa credit"
	PaymentSchemeSepa       PaymentScheme = "sepa"
	PaymentSchemeApplePay   PaymentScheme = "apple pay"
	PaymentSchemeGooglePay  PaymentScheme = "google pay"
	PaymentSchemeDoku       PaymentScheme = "doku"
	PaymentSchemeDragonpay  PaymentScheme = "dragonpay"
	PaymentSchemeMaestro    PaymentScheme = "maestro"
	PaymentSchemeMolpay     PaymentScheme = "molpay"
	PaymentSchemeA2a        PaymentScheme = "a2a"
	PaymentSchemeAchDebit   PaymentScheme = "ach debit"
	PaymentSchemeAch        PaymentScheme = "ach"
	PaymentSchemeRtp        PaymentScheme = "rtp"
)

func (e PaymentScheme) ToPointer() *PaymentScheme {
	return &e
}
func (e *PaymentScheme) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "unknown":
		fallthrough
	case "other":
		fallthrough
	case "visa":
		fallthrough
	case "mastercard":
		fallthrough
	case "amex":
		fallthrough
	case "diners":
		fallthrough
	case "discover":
		fallthrough
	case "jcb":
		fallthrough
	case "unionpay":
		fallthrough
	case "alipay":
		fallthrough
	case "cup":
		fallthrough
	case "sepa debit":
		fallthrough
	case "sepa credit":
		fallthrough
	case "sepa":
		fallthrough
	case "apple pay":
		fallthrough
	case "google pay":
		fallthrough
	case "doku":
		fallthrough
	case "dragonpay":
		fallthrough
	case "maestro":
		fallthrough
	case "molpay":
		fallthrough
	case "a2a":
		fallthrough
	case "ach debit":
		fallthrough
	case "ach":
		fallthrough
	case "rtp":
		*e = PaymentScheme(v)
		return nil
	default:
		return fmt.Errorf("invalid value for PaymentScheme: %v", v)
	}
}

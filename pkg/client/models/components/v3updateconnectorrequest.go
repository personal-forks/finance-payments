// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package components

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/formancehq/payments/pkg/client/internal/utils"
)

type V3UpdateConnectorRequestType string

const (
	V3UpdateConnectorRequestTypeAdyen         V3UpdateConnectorRequestType = "Adyen"
	V3UpdateConnectorRequestTypeAtlar         V3UpdateConnectorRequestType = "Atlar"
	V3UpdateConnectorRequestTypeBankingcircle V3UpdateConnectorRequestType = "Bankingcircle"
	V3UpdateConnectorRequestTypeColumn        V3UpdateConnectorRequestType = "Column"
	V3UpdateConnectorRequestTypeCurrencycloud V3UpdateConnectorRequestType = "Currencycloud"
	V3UpdateConnectorRequestTypeDummypay      V3UpdateConnectorRequestType = "Dummypay"
	V3UpdateConnectorRequestTypeGeneric       V3UpdateConnectorRequestType = "Generic"
	V3UpdateConnectorRequestTypeMangopay      V3UpdateConnectorRequestType = "Mangopay"
	V3UpdateConnectorRequestTypeModulr        V3UpdateConnectorRequestType = "Modulr"
	V3UpdateConnectorRequestTypeMoneycorp     V3UpdateConnectorRequestType = "Moneycorp"
	V3UpdateConnectorRequestTypeStripe        V3UpdateConnectorRequestType = "Stripe"
	V3UpdateConnectorRequestTypeWise          V3UpdateConnectorRequestType = "Wise"
)

type V3UpdateConnectorRequest struct {
	V3AdyenConfig         *V3AdyenConfig         `queryParam:"inline"`
	V3AtlarConfig         *V3AtlarConfig         `queryParam:"inline"`
	V3BankingcircleConfig *V3BankingcircleConfig `queryParam:"inline"`
	V3ColumnConfig        *V3ColumnConfig        `queryParam:"inline"`
	V3CurrencycloudConfig *V3CurrencycloudConfig `queryParam:"inline"`
	V3DummypayConfig      *V3DummypayConfig      `queryParam:"inline"`
	V3GenericConfig       *V3GenericConfig       `queryParam:"inline"`
	V3MangopayConfig      *V3MangopayConfig      `queryParam:"inline"`
	V3ModulrConfig        *V3ModulrConfig        `queryParam:"inline"`
	V3MoneycorpConfig     *V3MoneycorpConfig     `queryParam:"inline"`
	V3StripeConfig        *V3StripeConfig        `queryParam:"inline"`
	V3WiseConfig          *V3WiseConfig          `queryParam:"inline"`

	Type V3UpdateConnectorRequestType
}

func CreateV3UpdateConnectorRequestAdyen(adyen V3AdyenConfig) V3UpdateConnectorRequest {
	typ := V3UpdateConnectorRequestTypeAdyen

	typStr := string(typ)
	adyen.Provider = &typStr

	return V3UpdateConnectorRequest{
		V3AdyenConfig: &adyen,
		Type:          typ,
	}
}

func CreateV3UpdateConnectorRequestAtlar(atlar V3AtlarConfig) V3UpdateConnectorRequest {
	typ := V3UpdateConnectorRequestTypeAtlar

	typStr := string(typ)
	atlar.Provider = &typStr

	return V3UpdateConnectorRequest{
		V3AtlarConfig: &atlar,
		Type:          typ,
	}
}

func CreateV3UpdateConnectorRequestBankingcircle(bankingcircle V3BankingcircleConfig) V3UpdateConnectorRequest {
	typ := V3UpdateConnectorRequestTypeBankingcircle

	typStr := string(typ)
	bankingcircle.Provider = &typStr

	return V3UpdateConnectorRequest{
		V3BankingcircleConfig: &bankingcircle,
		Type:                  typ,
	}
}

func CreateV3UpdateConnectorRequestColumn(column V3ColumnConfig) V3UpdateConnectorRequest {
	typ := V3UpdateConnectorRequestTypeColumn

	typStr := string(typ)
	column.Provider = &typStr

	return V3UpdateConnectorRequest{
		V3ColumnConfig: &column,
		Type:           typ,
	}
}

func CreateV3UpdateConnectorRequestCurrencycloud(currencycloud V3CurrencycloudConfig) V3UpdateConnectorRequest {
	typ := V3UpdateConnectorRequestTypeCurrencycloud

	typStr := string(typ)
	currencycloud.Provider = &typStr

	return V3UpdateConnectorRequest{
		V3CurrencycloudConfig: &currencycloud,
		Type:                  typ,
	}
}

func CreateV3UpdateConnectorRequestDummypay(dummypay V3DummypayConfig) V3UpdateConnectorRequest {
	typ := V3UpdateConnectorRequestTypeDummypay

	typStr := string(typ)
	dummypay.Provider = &typStr

	return V3UpdateConnectorRequest{
		V3DummypayConfig: &dummypay,
		Type:             typ,
	}
}

func CreateV3UpdateConnectorRequestGeneric(generic V3GenericConfig) V3UpdateConnectorRequest {
	typ := V3UpdateConnectorRequestTypeGeneric

	typStr := string(typ)
	generic.Provider = &typStr

	return V3UpdateConnectorRequest{
		V3GenericConfig: &generic,
		Type:            typ,
	}
}

func CreateV3UpdateConnectorRequestMangopay(mangopay V3MangopayConfig) V3UpdateConnectorRequest {
	typ := V3UpdateConnectorRequestTypeMangopay

	typStr := string(typ)
	mangopay.Provider = &typStr

	return V3UpdateConnectorRequest{
		V3MangopayConfig: &mangopay,
		Type:             typ,
	}
}

func CreateV3UpdateConnectorRequestModulr(modulr V3ModulrConfig) V3UpdateConnectorRequest {
	typ := V3UpdateConnectorRequestTypeModulr

	typStr := string(typ)
	modulr.Provider = &typStr

	return V3UpdateConnectorRequest{
		V3ModulrConfig: &modulr,
		Type:           typ,
	}
}

func CreateV3UpdateConnectorRequestMoneycorp(moneycorp V3MoneycorpConfig) V3UpdateConnectorRequest {
	typ := V3UpdateConnectorRequestTypeMoneycorp

	typStr := string(typ)
	moneycorp.Provider = &typStr

	return V3UpdateConnectorRequest{
		V3MoneycorpConfig: &moneycorp,
		Type:              typ,
	}
}

func CreateV3UpdateConnectorRequestStripe(stripe V3StripeConfig) V3UpdateConnectorRequest {
	typ := V3UpdateConnectorRequestTypeStripe

	typStr := string(typ)
	stripe.Provider = &typStr

	return V3UpdateConnectorRequest{
		V3StripeConfig: &stripe,
		Type:           typ,
	}
}

func CreateV3UpdateConnectorRequestWise(wise V3WiseConfig) V3UpdateConnectorRequest {
	typ := V3UpdateConnectorRequestTypeWise

	typStr := string(typ)
	wise.Provider = &typStr

	return V3UpdateConnectorRequest{
		V3WiseConfig: &wise,
		Type:         typ,
	}
}

func (u *V3UpdateConnectorRequest) UnmarshalJSON(data []byte) error {

	type discriminator struct {
		Provider string `json:"provider"`
	}

	dis := new(discriminator)
	if err := json.Unmarshal(data, &dis); err != nil {
		return fmt.Errorf("could not unmarshal discriminator: %w", err)
	}

	switch dis.Provider {
	case "Adyen":
		v3AdyenConfig := new(V3AdyenConfig)
		if err := utils.UnmarshalJSON(data, &v3AdyenConfig, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Provider == Adyen) type V3AdyenConfig within V3UpdateConnectorRequest: %w", string(data), err)
		}

		u.V3AdyenConfig = v3AdyenConfig
		u.Type = V3UpdateConnectorRequestTypeAdyen
		return nil
	case "Atlar":
		v3AtlarConfig := new(V3AtlarConfig)
		if err := utils.UnmarshalJSON(data, &v3AtlarConfig, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Provider == Atlar) type V3AtlarConfig within V3UpdateConnectorRequest: %w", string(data), err)
		}

		u.V3AtlarConfig = v3AtlarConfig
		u.Type = V3UpdateConnectorRequestTypeAtlar
		return nil
	case "Bankingcircle":
		v3BankingcircleConfig := new(V3BankingcircleConfig)
		if err := utils.UnmarshalJSON(data, &v3BankingcircleConfig, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Provider == Bankingcircle) type V3BankingcircleConfig within V3UpdateConnectorRequest: %w", string(data), err)
		}

		u.V3BankingcircleConfig = v3BankingcircleConfig
		u.Type = V3UpdateConnectorRequestTypeBankingcircle
		return nil
	case "Column":
		v3ColumnConfig := new(V3ColumnConfig)
		if err := utils.UnmarshalJSON(data, &v3ColumnConfig, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Provider == Column) type V3ColumnConfig within V3UpdateConnectorRequest: %w", string(data), err)
		}

		u.V3ColumnConfig = v3ColumnConfig
		u.Type = V3UpdateConnectorRequestTypeColumn
		return nil
	case "Currencycloud":
		v3CurrencycloudConfig := new(V3CurrencycloudConfig)
		if err := utils.UnmarshalJSON(data, &v3CurrencycloudConfig, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Provider == Currencycloud) type V3CurrencycloudConfig within V3UpdateConnectorRequest: %w", string(data), err)
		}

		u.V3CurrencycloudConfig = v3CurrencycloudConfig
		u.Type = V3UpdateConnectorRequestTypeCurrencycloud
		return nil
	case "Dummypay":
		v3DummypayConfig := new(V3DummypayConfig)
		if err := utils.UnmarshalJSON(data, &v3DummypayConfig, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Provider == Dummypay) type V3DummypayConfig within V3UpdateConnectorRequest: %w", string(data), err)
		}

		u.V3DummypayConfig = v3DummypayConfig
		u.Type = V3UpdateConnectorRequestTypeDummypay
		return nil
	case "Generic":
		v3GenericConfig := new(V3GenericConfig)
		if err := utils.UnmarshalJSON(data, &v3GenericConfig, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Provider == Generic) type V3GenericConfig within V3UpdateConnectorRequest: %w", string(data), err)
		}

		u.V3GenericConfig = v3GenericConfig
		u.Type = V3UpdateConnectorRequestTypeGeneric
		return nil
	case "Mangopay":
		v3MangopayConfig := new(V3MangopayConfig)
		if err := utils.UnmarshalJSON(data, &v3MangopayConfig, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Provider == Mangopay) type V3MangopayConfig within V3UpdateConnectorRequest: %w", string(data), err)
		}

		u.V3MangopayConfig = v3MangopayConfig
		u.Type = V3UpdateConnectorRequestTypeMangopay
		return nil
	case "Modulr":
		v3ModulrConfig := new(V3ModulrConfig)
		if err := utils.UnmarshalJSON(data, &v3ModulrConfig, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Provider == Modulr) type V3ModulrConfig within V3UpdateConnectorRequest: %w", string(data), err)
		}

		u.V3ModulrConfig = v3ModulrConfig
		u.Type = V3UpdateConnectorRequestTypeModulr
		return nil
	case "Moneycorp":
		v3MoneycorpConfig := new(V3MoneycorpConfig)
		if err := utils.UnmarshalJSON(data, &v3MoneycorpConfig, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Provider == Moneycorp) type V3MoneycorpConfig within V3UpdateConnectorRequest: %w", string(data), err)
		}

		u.V3MoneycorpConfig = v3MoneycorpConfig
		u.Type = V3UpdateConnectorRequestTypeMoneycorp
		return nil
	case "Stripe":
		v3StripeConfig := new(V3StripeConfig)
		if err := utils.UnmarshalJSON(data, &v3StripeConfig, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Provider == Stripe) type V3StripeConfig within V3UpdateConnectorRequest: %w", string(data), err)
		}

		u.V3StripeConfig = v3StripeConfig
		u.Type = V3UpdateConnectorRequestTypeStripe
		return nil
	case "Wise":
		v3WiseConfig := new(V3WiseConfig)
		if err := utils.UnmarshalJSON(data, &v3WiseConfig, "", true, false); err != nil {
			return fmt.Errorf("could not unmarshal `%s` into expected (Provider == Wise) type V3WiseConfig within V3UpdateConnectorRequest: %w", string(data), err)
		}

		u.V3WiseConfig = v3WiseConfig
		u.Type = V3UpdateConnectorRequestTypeWise
		return nil
	}

	return fmt.Errorf("could not unmarshal `%s` into any supported union types for V3UpdateConnectorRequest", string(data))
}

func (u V3UpdateConnectorRequest) MarshalJSON() ([]byte, error) {
	if u.V3AdyenConfig != nil {
		return utils.MarshalJSON(u.V3AdyenConfig, "", true)
	}

	if u.V3AtlarConfig != nil {
		return utils.MarshalJSON(u.V3AtlarConfig, "", true)
	}

	if u.V3BankingcircleConfig != nil {
		return utils.MarshalJSON(u.V3BankingcircleConfig, "", true)
	}

	if u.V3ColumnConfig != nil {
		return utils.MarshalJSON(u.V3ColumnConfig, "", true)
	}

	if u.V3CurrencycloudConfig != nil {
		return utils.MarshalJSON(u.V3CurrencycloudConfig, "", true)
	}

	if u.V3DummypayConfig != nil {
		return utils.MarshalJSON(u.V3DummypayConfig, "", true)
	}

	if u.V3GenericConfig != nil {
		return utils.MarshalJSON(u.V3GenericConfig, "", true)
	}

	if u.V3MangopayConfig != nil {
		return utils.MarshalJSON(u.V3MangopayConfig, "", true)
	}

	if u.V3ModulrConfig != nil {
		return utils.MarshalJSON(u.V3ModulrConfig, "", true)
	}

	if u.V3MoneycorpConfig != nil {
		return utils.MarshalJSON(u.V3MoneycorpConfig, "", true)
	}

	if u.V3StripeConfig != nil {
		return utils.MarshalJSON(u.V3StripeConfig, "", true)
	}

	if u.V3WiseConfig != nil {
		return utils.MarshalJSON(u.V3WiseConfig, "", true)
	}

	return nil, errors.New("could not marshal union type V3UpdateConnectorRequest: all fields are null")
}

// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package components

import (
	"github.com/formancehq/payments/pkg/client/internal/utils"
)

type V3WiseConfig struct {
	APIKey           string  `json:"apiKey"`
	Name             string  `json:"name"`
	PageSize         *int64  `default:"25" json:"pageSize"`
	PollingPeriod    *string `default:"2m" json:"pollingPeriod"`
	Provider         *string `default:"Wise" json:"provider"`
	WebhookPublicKey string  `json:"webhookPublicKey"`
}

func (v V3WiseConfig) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(v, "", false)
}

func (v *V3WiseConfig) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &v, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *V3WiseConfig) GetAPIKey() string {
	if o == nil {
		return ""
	}
	return o.APIKey
}

func (o *V3WiseConfig) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *V3WiseConfig) GetPageSize() *int64 {
	if o == nil {
		return nil
	}
	return o.PageSize
}

func (o *V3WiseConfig) GetPollingPeriod() *string {
	if o == nil {
		return nil
	}
	return o.PollingPeriod
}

func (o *V3WiseConfig) GetProvider() *string {
	if o == nil {
		return nil
	}
	return o.Provider
}

func (o *V3WiseConfig) GetWebhookPublicKey() string {
	if o == nil {
		return ""
	}
	return o.WebhookPublicKey
}

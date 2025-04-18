// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/payments/pkg/client/models/components"
)

type V3AddAccountToPoolRequest struct {
	// The pool ID
	PoolID string `pathParam:"style=simple,explode=false,name=poolID"`
	// The account ID
	AccountID string `pathParam:"style=simple,explode=false,name=accountID"`
}

func (o *V3AddAccountToPoolRequest) GetPoolID() string {
	if o == nil {
		return ""
	}
	return o.PoolID
}

func (o *V3AddAccountToPoolRequest) GetAccountID() string {
	if o == nil {
		return ""
	}
	return o.AccountID
}

type V3AddAccountToPoolResponse struct {
	HTTPMeta components.HTTPMetadata `json:"-"`
	// Error
	V3ErrorResponse *components.V3ErrorResponse
}

func (o *V3AddAccountToPoolResponse) GetHTTPMeta() components.HTTPMetadata {
	if o == nil {
		return components.HTTPMetadata{}
	}
	return o.HTTPMeta
}

func (o *V3AddAccountToPoolResponse) GetV3ErrorResponse() *components.V3ErrorResponse {
	if o == nil {
		return nil
	}
	return o.V3ErrorResponse
}

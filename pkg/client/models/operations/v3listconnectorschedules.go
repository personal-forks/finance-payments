// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/formancehq/payments/pkg/client/models/components"
)

type V3ListConnectorSchedulesRequest struct {
	// The connector ID
	ConnectorID string `pathParam:"style=simple,explode=false,name=connectorID"`
	// The number of items to return
	PageSize *int64 `queryParam:"style=form,explode=true,name=pageSize"`
	// Parameter used in pagination requests. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set.
	//
	Cursor      *string        `queryParam:"style=form,explode=true,name=cursor"`
	RequestBody map[string]any `request:"mediaType=application/json"`
}

func (o *V3ListConnectorSchedulesRequest) GetConnectorID() string {
	if o == nil {
		return ""
	}
	return o.ConnectorID
}

func (o *V3ListConnectorSchedulesRequest) GetPageSize() *int64 {
	if o == nil {
		return nil
	}
	return o.PageSize
}

func (o *V3ListConnectorSchedulesRequest) GetCursor() *string {
	if o == nil {
		return nil
	}
	return o.Cursor
}

func (o *V3ListConnectorSchedulesRequest) GetRequestBody() map[string]any {
	if o == nil {
		return nil
	}
	return o.RequestBody
}

type V3ListConnectorSchedulesResponse struct {
	HTTPMeta components.HTTPMetadata `json:"-"`
	// OK
	V3ConnectorSchedulesCursorResponse *components.V3ConnectorSchedulesCursorResponse
	// Error
	V3ErrorResponse *components.V3ErrorResponse
}

func (o *V3ListConnectorSchedulesResponse) GetHTTPMeta() components.HTTPMetadata {
	if o == nil {
		return components.HTTPMetadata{}
	}
	return o.HTTPMeta
}

func (o *V3ListConnectorSchedulesResponse) GetV3ConnectorSchedulesCursorResponse() *components.V3ConnectorSchedulesCursorResponse {
	if o == nil {
		return nil
	}
	return o.V3ConnectorSchedulesCursorResponse
}

func (o *V3ListConnectorSchedulesResponse) GetV3ErrorResponse() *components.V3ErrorResponse {
	if o == nil {
		return nil
	}
	return o.V3ErrorResponse
}

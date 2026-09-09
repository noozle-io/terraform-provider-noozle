package client

import (
	"context"
	"fmt"
)

// OutletTransformState is the selected and effective transform state returned
// by a typed outlet transform endpoint.
type OutletTransformState struct {
	SelectedTransform  *string `json:"selected_transform"`
	EffectiveTransform *string `json:"effective_transform"`
	State              *string `json:"state"`
	Applied            bool    `json:"applied"`
	Generation         *int64  `json:"generation"`
	LastError          *string `json:"last_error"`
}

type outletTransformSetRequest struct {
	Transform string `json:"transform"`
}

// GetTypedOutletTransform reads transform state using the provider-specific
// Terraform route selected by outletType.
func (c *Client) GetTypedOutletTransform(ctx context.Context, outletType string, outletID int64) (OutletTransformState, error) {
	var response OutletTransformState
	if err := c.doJSON(ctx, httpMethodGet, typedOutletTransformPath(outletType, outletID), nil, &response, statusOK); err != nil {
		return OutletTransformState{}, err
	}
	return response, nil
}

func (c *Client) SetTypedOutletTransform(ctx context.Context, outletType string, outletID int64, transform string) (OutletTransformState, error) {
	var response OutletTransformState
	if err := c.doJSON(ctx, httpMethodPut, typedOutletTransformPath(outletType, outletID), outletTransformSetRequest{Transform: transform}, &response, statusOK); err != nil {
		return OutletTransformState{}, err
	}
	return response, nil
}

func (c *Client) ClearTypedOutletTransform(ctx context.Context, outletType string, outletID int64) (OutletTransformState, error) {
	var response OutletTransformState
	if err := c.doJSON(ctx, httpMethodDelete, typedOutletTransformPath(outletType, outletID), nil, &response, statusOK); err != nil {
		return OutletTransformState{}, err
	}
	return response, nil
}

func typedOutletTransformPath(outletType string, outletID int64) string {
	return fmt.Sprintf("/v1/terraform/%s/%d/transform", outletType, outletID)
}

package client

import (
	"context"
	"fmt"
	"net/url"
)

// OutletRuntimeError describes a terminal runtime realization error.
type OutletRuntimeError struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// OutletRuntimeRealization is the runtime observation for a requested Outlet generation.
type OutletRuntimeRealization struct {
	ActiveGeneration   int64               `json:"active_generation"`
	RealizationAttempt int64               `json:"realization_attempt"`
	Phase              string              `json:"phase"`
	Ready              bool                `json:"ready"`
	Degraded           bool                `json:"degraded"`
	LastError          *OutletRuntimeError `json:"last_error"`
}

// OutletOperation identifies a durable Outlet mutation that can be observed.
type OutletOperation struct {
	Kind              string                   `json:"kind"`
	DesiredGeneration int64                    `json:"desired_generation"`
	StatusLink        string                   `json:"status_link"`
	Realization       OutletRuntimeRealization `json:"realization"`
}

// OutletStatus is the typed realization status for one requested Outlet generation.
type OutletStatus struct {
	OutletID            int64                    `json:"outlet_id"`
	TenantID            string                   `json:"tenant_id"`
	DesiredStatus       string                   `json:"desired_status"`
	RequestedGeneration int64                    `json:"requested_generation"`
	DesiredGeneration   int64                    `json:"desired_generation"`
	Superseded          bool                     `json:"superseded"`
	RuntimeStatus       string                   `json:"runtime_status"`
	Realization         OutletRuntimeRealization `json:"realization"`
}

// GetTypedOutletStatus observes one requested generation using the typed Outlet route.
func (c *Client) GetTypedOutletStatus(ctx context.Context, outletType string, outletID, generation int64) (OutletStatus, error) {
	var response OutletStatus
	path := fmt.Sprintf("/v1/terraform/%s/%d/status?generation=%d", url.PathEscape(outletType), outletID, generation)
	if err := c.doJSON(ctx, httpMethodGet, path, nil, &response, statusOK); err != nil {
		return OutletStatus{}, err
	}

	return response, nil
}

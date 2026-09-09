package client

import (
	"context"
	"fmt"
)

type OutletRef struct {
	ID       int64
	QueryIDs []int64
}

type outletLinkRequest struct {
	QueryIDs []int64 `json:"query_ids"`
}

type outletResponse struct {
	ID       int64   `json:"id"`
	QueryIDs []int64 `json:"queries"`
}

func (c *Client) LinkOutletQuery(ctx context.Context, outletID, queryID int64) error {
	if err := c.requireTenantContext(); err != nil {
		return err
	}

	body := outletLinkRequest{QueryIDs: []int64{queryID}}
	return c.doTenantJSON(ctx, httpMethodPost, fmt.Sprintf("/v1/tenants/%s/outlets/%d/queries", c.tenantContext.TenantID, outletID), body, nil, statusOK)
}

func (c *Client) UnlinkOutletQuery(ctx context.Context, outletID, queryID int64) error {
	if err := c.requireTenantContext(); err != nil {
		return err
	}

	body := outletLinkRequest{QueryIDs: []int64{queryID}}
	return c.doTenantJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/tenants/%s/outlets/%d/queries", c.tenantContext.TenantID, outletID), body, nil, statusOK)
}

func (c *Client) GetOutlet(ctx context.Context, id int64) (OutletRef, error) {
	if err := c.requireTenantContext(); err != nil {
		return OutletRef{}, err
	}

	var response outletResponse
	if err := c.doTenantJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/tenants/%s/outlets/%d", c.tenantContext.TenantID, id), nil, &response, statusOK); err != nil {
		return OutletRef{}, err
	}

	return OutletRef(response), nil
}

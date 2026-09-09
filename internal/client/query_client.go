package client

import (
	"context"
	"fmt"
	"time"
)

type Query struct {
	ID             int64
	Name           string
	Description    string
	FullExpression string
	Tags           []string
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type QueryRequest struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	FullExpression string   `json:"full_expression"`
	Tags           []string `json:"tags"`
	IsActive       bool     `json:"is_active"`
}

func (c *Client) CreateQuery(ctx context.Context, req QueryRequest) (Query, error) {
	if err := c.requireTenantContext(); err != nil {
		return Query{}, err
	}

	var response queryResponse
	if err := c.doTenantJSON(ctx, httpMethodPost, fmt.Sprintf("/v1/tenants/%s/queries", c.tenantContext.TenantID), req, &response, statusCreated); err != nil {
		return Query{}, err
	}

	return response.toQuery(), nil
}

func (c *Client) GetQuery(ctx context.Context, id int64) (Query, error) {
	if err := c.requireTenantContext(); err != nil {
		return Query{}, err
	}

	var response queryDetailResponse
	if err := c.doTenantJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/tenants/%s/queries/%d", c.tenantContext.TenantID, id), nil, &response, statusOK); err != nil {
		return Query{}, err
	}

	return response.toQuery(), nil
}

func (c *Client) UpdateQuery(ctx context.Context, id int64, req QueryRequest) (Query, error) {
	if err := c.requireTenantContext(); err != nil {
		return Query{}, err
	}

	var response queryResponse
	if err := c.doTenantJSON(ctx, httpMethodPut, fmt.Sprintf("/v1/tenants/%s/queries/%d", c.tenantContext.TenantID, id), req, &response, statusOK); err != nil {
		return Query{}, err
	}

	return response.toQuery(), nil
}

func (c *Client) DeleteQuery(ctx context.Context, id int64) error {
	if err := c.requireTenantContext(); err != nil {
		return err
	}

	return c.doTenantJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/tenants/%s/queries/%d", c.tenantContext.TenantID, id), nil, nil, statusNoContent)
}

type queryResponse struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	FullExpression string    `json:"full_expression"`
	Tags           []string  `json:"tags"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (q queryResponse) toQuery() Query {
	return Query(q)
}

type queryDetailResponse struct {
	queryResponse
	ActiveTraceSessionIDs []string `json:"active_trace_session_ids"`
}

const (
	httpMethodGet    = "GET"
	httpMethodPost   = "POST"
	httpMethodPut    = "PUT"
	httpMethodDelete = "DELETE"
	httpMethodPatch  = "PATCH"

	statusOK        = 200
	statusCreated   = 201
	statusAccepted  = 202
	statusNoContent = 204
)

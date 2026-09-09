package client

import (
	"context"
	"fmt"
)

type SseOutlet struct {
	ID                int64
	Name              string
	Description       string
	Enabled           bool
	EventType         string
	Path              string
	RetryMS           int64
	NotifyPolicy      *string
	StreamURL         string
	ApplyStatus       string
	LastApplyError    string
	ConfigGeneration  int64
	AppliedGeneration int64
}

type SseOutletCreateRequest struct {
	Name         string
	Description  string
	Enabled      bool
	EventType    string
	Path         string
	RetryMS      int64
	NotifyPolicy *string
}

type SseOutletUpdateRequest struct {
	Name         string
	Description  string
	Enabled      bool
	EventType    string
	Path         string
	RetryMS      int64
	NotifyPolicy *string
}

func (c *Client) CreateSseOutlet(ctx context.Context, req SseOutletCreateRequest) (SseOutlet, error) {
	fields := map[string]any{
		"sse.event_type": req.EventType,
		"sse.path":       req.Path,
		"sse.retry_ms":   req.RetryMS,
	}
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)

	var response typedOutletResponse
	if err := c.doMultipartPayload(ctx, httpMethodPost, "/v1/terraform/sse", typedOutletCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, &response, statusCreated); err != nil {
		return SseOutlet{}, err
	}

	return decodeSseOutlet(response)
}

func (c *Client) GetSseOutlet(ctx context.Context, id int64) (SseOutlet, error) {
	var response typedOutletResponse
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/sse/%d", id), nil, &response, statusOK); err != nil {
		return SseOutlet{}, err
	}

	return decodeSseOutlet(response)
}

func (c *Client) UpdateSseOutlet(ctx context.Context, id int64, req SseOutletUpdateRequest) error {
	fields := map[string]any{
		"sse.event_type": req.EventType,
		"sse.path":       req.Path,
		"sse.retry_ms":   req.RetryMS,
	}
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/sse/%d", id), typedOutletUpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, nil, statusOK)
}

func (c *Client) DeleteSseOutlet(ctx context.Context, id int64) error {
	_, err := c.deleteTypedOutlet(ctx, fmt.Sprintf("/v1/terraform/sse/%d", id))
	return err
}

func (c *Client) SetSseEventType(ctx context.Context, id int64, value string) error {
	return c.patchOutletField(ctx, fmt.Sprintf("/v1/terraform/sse/%d/fields/sse.event_type", id), value)
}

func (c *Client) SetSsePath(ctx context.Context, id int64, value string) error {
	return c.patchOutletField(ctx, fmt.Sprintf("/v1/terraform/sse/%d/fields/sse.path", id), value)
}

func (c *Client) SetSseRetryMS(ctx context.Context, id int64, value int64) error {
	return c.patchOutletField(ctx, fmt.Sprintf("/v1/terraform/sse/%d/fields/sse.retry_ms", id), value)
}

func (c *Client) ClearSseNotifyPolicy(ctx context.Context, id int64) error {
	return c.deleteOutletField(ctx, fmt.Sprintf("/v1/terraform/sse/%d/fields/notify_policy", id))
}

func decodeSseOutlet(response typedOutletResponse) (SseOutlet, error) {
	eventType, err := requiredStringField(response.Config, "sse.event_type")
	if err != nil {
		return SseOutlet{}, err
	}
	path, err := requiredStringField(response.Config, "sse.path")
	if err != nil {
		return SseOutlet{}, err
	}
	retryMS, err := requiredInt64Field(response.Config, "sse.retry_ms")
	if err != nil {
		return SseOutlet{}, err
	}
	notifyPolicy, err := stringFieldPointer(response.Config, "notify_policy")
	if err != nil {
		return SseOutlet{}, err
	}

	return SseOutlet{
		ID:                response.ID,
		Name:              response.Name,
		Description:       response.Description,
		Enabled:           response.Enabled,
		EventType:         eventType,
		Path:              path,
		RetryMS:           retryMS,
		NotifyPolicy:      notifyPolicy,
		StreamURL:         response.StreamURL,
		ApplyStatus:       response.ApplyStatus,
		LastApplyError:    response.LastApplyError,
		ConfigGeneration:  response.ConfigGeneration,
		AppliedGeneration: response.AppliedGeneration,
	}, nil
}

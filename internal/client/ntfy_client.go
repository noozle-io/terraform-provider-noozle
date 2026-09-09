package client

import (
	"context"
	"fmt"
)

type NtfyOutlet struct {
	ID                int64
	Name              string
	Description       string
	Enabled           bool
	Topic             string
	NotifyPolicy      *string
	ServerURL         *string
	StreamURL         string
	ApplyStatus       string
	LastApplyError    string
	ConfigGeneration  int64
	AppliedGeneration int64
}

type NtfyOutletCreateRequest struct {
	Name         string
	Description  string
	Enabled      bool
	Topic        string
	NotifyPolicy *string
	ServerURL    *string
}

type NtfyOutletUpdateRequest struct {
	Name         string
	Description  string
	Enabled      bool
	Topic        string
	NotifyPolicy *string
	ServerURL    *string
}

func (c *Client) CreateNtfyOutlet(ctx context.Context, req NtfyOutletCreateRequest) (NtfyOutlet, error) {
	fields := map[string]any{
		"ntfy.topic": req.Topic,
	}
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalStringField(fields, "ntfy.server_url", req.ServerURL)

	var response typedOutletResponse
	if err := c.doMultipartPayload(ctx, httpMethodPost, "/v1/terraform/ntfy", typedOutletCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, &response, statusCreated); err != nil {
		return NtfyOutlet{}, err
	}

	return decodeNtfyOutlet(response)
}

func (c *Client) GetNtfyOutlet(ctx context.Context, id int64) (NtfyOutlet, error) {
	var response typedOutletResponse
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/ntfy/%d", id), nil, &response, statusOK); err != nil {
		return NtfyOutlet{}, err
	}

	return decodeNtfyOutlet(response)
}

func (c *Client) UpdateNtfyOutlet(ctx context.Context, id int64, req NtfyOutletUpdateRequest) error {
	fields := map[string]any{
		"ntfy.topic": req.Topic,
	}
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalStringField(fields, "ntfy.server_url", req.ServerURL)

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/ntfy/%d", id), typedOutletUpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, nil, statusOK)
}

func (c *Client) DeleteNtfyOutlet(ctx context.Context, id int64) error {
	_, err := c.DeleteNtfyOutletOperation(ctx, id)
	return err
}

func (c *Client) DeleteNtfyOutletOperation(ctx context.Context, id int64) (OutletOperation, error) {
	return c.deleteTypedOutlet(ctx, fmt.Sprintf("/v1/terraform/ntfy/%d", id))
}

func (c *Client) SetNtfyTopic(ctx context.Context, id int64, value string) error {
	return c.patchOutletField(ctx, fmt.Sprintf("/v1/terraform/ntfy/%d/fields/ntfy.topic", id), value)
}

func (c *Client) ClearNtfyNotifyPolicy(ctx context.Context, id int64) error {
	return c.deleteOutletField(ctx, fmt.Sprintf("/v1/terraform/ntfy/%d/fields/notify_policy", id))
}

func (c *Client) ClearNtfyServerURL(ctx context.Context, id int64) error {
	return c.deleteOutletField(ctx, fmt.Sprintf("/v1/terraform/ntfy/%d/fields/ntfy.server_url", id))
}

func decodeNtfyOutlet(response typedOutletResponse) (NtfyOutlet, error) {
	topic, err := requiredStringField(response.Config, "ntfy.topic")
	if err != nil {
		return NtfyOutlet{}, err
	}
	notifyPolicy, err := stringFieldPointer(response.Config, "notify_policy")
	if err != nil {
		return NtfyOutlet{}, err
	}
	serverURL, err := stringFieldPointer(response.Config, "ntfy.server_url")
	if err != nil {
		return NtfyOutlet{}, err
	}

	return NtfyOutlet{
		ID:                response.ID,
		Name:              response.Name,
		Description:       response.Description,
		Enabled:           response.Enabled,
		Topic:             topic,
		NotifyPolicy:      notifyPolicy,
		ServerURL:         serverURL,
		StreamURL:         response.StreamURL,
		ApplyStatus:       response.ApplyStatus,
		LastApplyError:    response.LastApplyError,
		ConfigGeneration:  response.ConfigGeneration,
		AppliedGeneration: response.AppliedGeneration,
	}, nil
}

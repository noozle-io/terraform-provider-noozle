package client

import (
	"context"
	"fmt"
)

const (
	httpServerTLSFeatureID    = "http_server.tls"
	httpServerTLSCertRole     = "server_cert"
	httpServerTLSKeyRole      = "server_key"
	httpServerTLSCertPartName = "server_cert_part"
	httpServerTLSKeyPartName  = "server_key_part"
)

type HttpServerOutlet struct {
	ID                 int64
	Name               string
	Description        string
	Enabled            bool
	Address            *string
	AllowedVerbs       []string
	CorsAllowedHeaders []string
	CorsAllowedMethods []string
	CorsAllowedOrigins []string
	CorsEnabled        *bool
	Heartbeat          *string
	NotifyPolicy       *string
	Path               string
	PingPeriod         *string
	PongWait           *string
	StreamFormat       *string
	StreamPath         *string
	Timeout            *string
	WriteWait          *string
	WsMessageType      *string
	WsPath             *string
	ApplyStatus        string
	LastApplyError     string
	ConfigGeneration   int64
	AppliedGeneration  int64
}

type HttpServerOutletCreateRequest struct {
	Name               string
	Description        string
	Enabled            bool
	Address            *string
	AllowedVerbs       []string
	CorsAllowedHeaders []string
	CorsAllowedMethods []string
	CorsAllowedOrigins []string
	CorsEnabled        *bool
	Heartbeat          *string
	NotifyPolicy       *string
	Path               string
	PingPeriod         *string
	PongWait           *string
	StreamFormat       *string
	StreamPath         *string
	Timeout            *string
	WriteWait          *string
	WsMessageType      *string
	WsPath             *string
}

type HttpServerOutletUpdateRequest = HttpServerOutletCreateRequest

type OutletFeatureInput struct {
	Key        string `json:"key"`
	Kind       string `json:"kind"`
	Required   bool   `json:"required"`
	Set        bool   `json:"set"`
	Filename   string `json:"filename"`
	Sensitive  bool   `json:"sensitive"`
	FileBacked bool   `json:"file_backed"`
}

type OutletFeatureDetail struct {
	Feature           string               `json:"feature"`
	Description       string               `json:"description"`
	Enabled           bool                 `json:"enabled"`
	State             string               `json:"state"`
	DesiredGeneration *int64               `json:"desired_generation"`
	AppliedGeneration *int64               `json:"applied_generation"`
	LastError         string               `json:"last_error"`
	Conflicts         []string             `json:"conflicts"`
	Inputs            []OutletFeatureInput `json:"inputs"`
}

type HttpServerTLSFeatureRequest struct {
	ServerCertFilename      string
	ServerCertContentBase64 string
	ServerKeyFilename       string
	ServerKeyContentBase64  string
}

func (c *Client) CreateHttpServerOutlet(ctx context.Context, req HttpServerOutletCreateRequest) (HttpServerOutlet, error) {
	fields := map[string]any{
		"http_server.path": req.Path,
	}
	setOptionalStringField(fields, "http_server.address", req.Address)
	setOptionalStringListField(fields, "http_server.allowed_verbs", req.AllowedVerbs)
	setOptionalStringListField(fields, "http_server.cors.allowed_headers", req.CorsAllowedHeaders)
	setOptionalStringListField(fields, "http_server.cors.allowed_methods", req.CorsAllowedMethods)
	setOptionalStringListField(fields, "http_server.cors.allowed_origins", req.CorsAllowedOrigins)
	setOptionalBoolField(fields, "http_server.cors.enabled", req.CorsEnabled)
	setOptionalStringField(fields, "http_server.heartbeat", req.Heartbeat)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalStringField(fields, "http_server.ping_period", req.PingPeriod)
	setOptionalStringField(fields, "http_server.pong_wait", req.PongWait)
	setOptionalStringField(fields, "http_server.stream_format", req.StreamFormat)
	setOptionalStringField(fields, "http_server.stream_path", req.StreamPath)
	setOptionalStringField(fields, "http_server.timeout", req.Timeout)
	setOptionalStringField(fields, "http_server.write_wait", req.WriteWait)
	setOptionalStringField(fields, "http_server.ws_message_type", req.WsMessageType)
	setOptionalStringField(fields, "http_server.ws_path", req.WsPath)

	var response typedOutletResponse
	if err := c.doMultipartPayload(ctx, httpMethodPost, "/v1/terraform/http_server", typedOutletCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, &response, statusCreated); err != nil {
		return HttpServerOutlet{}, err
	}

	return decodeHttpServerOutlet(response)
}

func (c *Client) GetHttpServerOutlet(ctx context.Context, id int64) (HttpServerOutlet, error) {
	var response typedOutletResponse
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/http_server/%d", id), nil, &response, statusOK); err != nil {
		return HttpServerOutlet{}, err
	}

	return decodeHttpServerOutlet(response)
}

func (c *Client) UpdateHttpServerOutlet(ctx context.Context, id int64, req HttpServerOutletUpdateRequest) error {
	fields := map[string]any{
		"http_server.path": req.Path,
	}
	setOptionalStringField(fields, "http_server.address", req.Address)
	setOptionalStringListField(fields, "http_server.allowed_verbs", req.AllowedVerbs)
	setOptionalStringListField(fields, "http_server.cors.allowed_headers", req.CorsAllowedHeaders)
	setOptionalStringListField(fields, "http_server.cors.allowed_methods", req.CorsAllowedMethods)
	setOptionalStringListField(fields, "http_server.cors.allowed_origins", req.CorsAllowedOrigins)
	setOptionalBoolField(fields, "http_server.cors.enabled", req.CorsEnabled)
	setOptionalStringField(fields, "http_server.heartbeat", req.Heartbeat)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalStringField(fields, "http_server.ping_period", req.PingPeriod)
	setOptionalStringField(fields, "http_server.pong_wait", req.PongWait)
	setOptionalStringField(fields, "http_server.stream_format", req.StreamFormat)
	setOptionalStringField(fields, "http_server.stream_path", req.StreamPath)
	setOptionalStringField(fields, "http_server.timeout", req.Timeout)
	setOptionalStringField(fields, "http_server.write_wait", req.WriteWait)
	setOptionalStringField(fields, "http_server.ws_message_type", req.WsMessageType)
	setOptionalStringField(fields, "http_server.ws_path", req.WsPath)

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/http_server/%d", id), typedOutletUpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, nil, statusOK)
}

func (c *Client) DeleteHttpServerOutlet(ctx context.Context, id int64) error {
	_, err := c.DeleteHttpServerOutletOperation(ctx, id)
	return err
}

func (c *Client) DeleteHttpServerOutletOperation(ctx context.Context, id int64) (OutletOperation, error) {
	return c.deleteTypedOutlet(ctx, fmt.Sprintf("/v1/terraform/http_server/%d", id))
}

func (c *Client) SetHttpServerField(ctx context.Context, id int64, field string, value any) error {
	return c.patchOutletField(ctx, fmt.Sprintf("/v1/terraform/http_server/%d/fields/%s", id, field), value)
}

func (c *Client) DeleteHttpServerField(ctx context.Context, id int64, field string) error {
	return c.deleteOutletField(ctx, fmt.Sprintf("/v1/terraform/http_server/%d/fields/%s", id, field))
}

func (c *Client) ListHttpServerCredentials(ctx context.Context, outletID int64) ([]OutletMaterialSummary, error) {
	return c.listOutletMaterials(ctx, fmt.Sprintf("/v1/terraform/http_server/%d/credentials", outletID))
}

func (c *Client) DeleteHttpServerCredential(ctx context.Context, outletID int64, role string) (OutletMaterialDeleteResponse, error) {
	return c.deleteOutletMaterial(ctx, fmt.Sprintf("/v1/terraform/http_server/%d/credentials/%s", outletID, role))
}

func (c *Client) GetHttpServerTLSFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	var response OutletFeatureDetail
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/http_server/%d/features/%s", outletID, httpServerTLSFeatureID), nil, &response, statusOK); err != nil {
		return OutletFeatureDetail{}, err
	}

	return response, nil
}

func (c *Client) ConfigureHttpServerTLSFeature(ctx context.Context, outletID int64, req HttpServerTLSFeatureRequest) error {
	serverCertContent, err := decodeBase64File(req.ServerCertContentBase64)
	if err != nil {
		return err
	}
	serverKeyContent, err := decodeBase64File(req.ServerKeyContentBase64)
	if err != nil {
		return err
	}

	return c.doMultipartPayload(ctx, httpMethodPost, fmt.Sprintf("/v1/terraform/http_server/%d/features/%s", outletID, httpServerTLSFeatureID), map[string]any{
		"inputs": map[string]any{
			httpServerTLSCertRole: payloadFileRef{File: httpServerTLSCertPartName},
			httpServerTLSKeyRole:  payloadFileRef{File: httpServerTLSKeyPartName},
		},
	}, []multipartFilePart{
		{FieldName: httpServerTLSCertPartName, FileName: req.ServerCertFilename, Content: serverCertContent},
		{FieldName: httpServerTLSKeyPartName, FileName: req.ServerKeyFilename, Content: serverKeyContent},
	}, nil, statusOK)
}

func (c *Client) UpdateHttpServerTLSFeature(ctx context.Context, outletID int64, req HttpServerTLSFeatureRequest) error {
	serverCertContent, err := decodeBase64File(req.ServerCertContentBase64)
	if err != nil {
		return err
	}
	serverKeyContent, err := decodeBase64File(req.ServerKeyContentBase64)
	if err != nil {
		return err
	}

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/http_server/%d/features/%s", outletID, httpServerTLSFeatureID), map[string]any{
		"inputs": map[string]any{
			httpServerTLSCertRole: payloadFileRef{File: httpServerTLSCertPartName},
			httpServerTLSKeyRole:  payloadFileRef{File: httpServerTLSKeyPartName},
		},
	}, []multipartFilePart{
		{FieldName: httpServerTLSCertPartName, FileName: req.ServerCertFilename, Content: serverCertContent},
		{FieldName: httpServerTLSKeyPartName, FileName: req.ServerKeyFilename, Content: serverKeyContent},
	}, nil, statusOK)
}

func (c *Client) DeleteHttpServerTLSFeature(ctx context.Context, outletID int64) error {
	return c.doJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/terraform/http_server/%d/features/%s", outletID, httpServerTLSFeatureID), nil, nil, statusOK)
}

func decodeHttpServerOutlet(response typedOutletResponse) (HttpServerOutlet, error) {
	path, err := requiredStringField(response.Config, "http_server.path")
	if err != nil {
		return HttpServerOutlet{}, err
	}

	address, err := stringFieldPointer(response.Config, "http_server.address")
	if err != nil {
		return HttpServerOutlet{}, err
	}
	allowedVerbs, err := stringListField(response.Config, "http_server.allowed_verbs")
	if err != nil {
		return HttpServerOutlet{}, err
	}
	corsAllowedHeaders, err := stringListField(response.Config, "http_server.cors.allowed_headers")
	if err != nil {
		return HttpServerOutlet{}, err
	}
	corsAllowedMethods, err := stringListField(response.Config, "http_server.cors.allowed_methods")
	if err != nil {
		return HttpServerOutlet{}, err
	}
	corsAllowedOrigins, err := stringListField(response.Config, "http_server.cors.allowed_origins")
	if err != nil {
		return HttpServerOutlet{}, err
	}
	corsEnabled, err := boolFieldPointer(response.Config, "http_server.cors.enabled")
	if err != nil {
		return HttpServerOutlet{}, err
	}
	heartbeat, err := stringFieldPointer(response.Config, "http_server.heartbeat")
	if err != nil {
		return HttpServerOutlet{}, err
	}
	notifyPolicy, err := notifyPolicyPointerPreservingEmpty(response.Config)
	if err != nil {
		return HttpServerOutlet{}, err
	}
	pingPeriod, err := stringFieldPointer(response.Config, "http_server.ping_period")
	if err != nil {
		return HttpServerOutlet{}, err
	}
	pongWait, err := stringFieldPointer(response.Config, "http_server.pong_wait")
	if err != nil {
		return HttpServerOutlet{}, err
	}
	streamFormat, err := stringFieldPointer(response.Config, "http_server.stream_format")
	if err != nil {
		return HttpServerOutlet{}, err
	}
	streamPath, err := stringFieldPointer(response.Config, "http_server.stream_path")
	if err != nil {
		return HttpServerOutlet{}, err
	}
	timeout, err := stringFieldPointer(response.Config, "http_server.timeout")
	if err != nil {
		return HttpServerOutlet{}, err
	}
	writeWait, err := stringFieldPointer(response.Config, "http_server.write_wait")
	if err != nil {
		return HttpServerOutlet{}, err
	}
	wsMessageType, err := stringFieldPointer(response.Config, "http_server.ws_message_type")
	if err != nil {
		return HttpServerOutlet{}, err
	}
	wsPath, err := stringFieldPointer(response.Config, "http_server.ws_path")
	if err != nil {
		return HttpServerOutlet{}, err
	}

	return HttpServerOutlet{
		ID:                 response.ID,
		Name:               response.Name,
		Description:        response.Description,
		Enabled:            response.Enabled,
		Address:            address,
		AllowedVerbs:       allowedVerbs,
		CorsAllowedHeaders: corsAllowedHeaders,
		CorsAllowedMethods: corsAllowedMethods,
		CorsAllowedOrigins: corsAllowedOrigins,
		CorsEnabled:        corsEnabled,
		Heartbeat:          heartbeat,
		NotifyPolicy:       notifyPolicy,
		Path:               path,
		PingPeriod:         pingPeriod,
		PongWait:           pongWait,
		StreamFormat:       streamFormat,
		StreamPath:         streamPath,
		Timeout:            timeout,
		WriteWait:          writeWait,
		WsMessageType:      wsMessageType,
		WsPath:             wsPath,
		ApplyStatus:        response.ApplyStatus,
		LastApplyError:     response.LastApplyError,
		ConfigGeneration:   response.ConfigGeneration,
		AppliedGeneration:  response.AppliedGeneration,
	}, nil
}

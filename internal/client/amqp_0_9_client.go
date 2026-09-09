package client

import (
	"context"
	"fmt"
)

const (
	amqp09URLsRole         = "amqp_0_9_urls"
	amqp09TLSFeatureID     = "amqp_0_9.tls"
	amqp09MTLSFeatureID    = "amqp_0_9.mtls"
	amqpTLSCARole          = "ca_cert"
	amqpMTLSClientCertRole = "client_cert"
	amqpMTLSClientKeyRole  = "client_key"
	amqpTLSCAPartName      = "ca_cert_part"
	amqpMTLSClientCertPart = "client_cert_part"
	amqpMTLSClientKeyPart  = "client_key_part"
)

type Amqp09Outlet struct {
	ID                      int64
	Name                    string
	Description             string
	Enabled                 bool
	AppID                   *string
	ContentEncoding         *string
	ContentType             *string
	CorrelationID           *string
	Exchange                string
	ExchangeDeclareDurable  *bool
	ExchangeDeclareEnabled  *bool
	ExchangeDeclareType     *string
	Expiration              *string
	Immediate               *bool
	Key                     string
	Mandatory               *bool
	MaxInFlight             *int64
	MessageID               *string
	MetadataExcludePrefixes []string
	NotifyPolicy            *string
	Persistent              *bool
	Priority                *string
	ReplyTo                 *string
	Timeout                 *string
	TLSEnableRenegotiation  *bool
	TLSSkipCertVerify       *bool
	Type                    *string
	UserID                  *string
	ApplyStatus             string
	LastApplyError          string
	ConfigGeneration        int64
	AppliedGeneration       int64
}

type Amqp09OutletCreateRequest struct {
	Name                    string
	Description             string
	Enabled                 bool
	AppID                   *string
	ContentEncoding         *string
	ContentType             *string
	CorrelationID           *string
	Exchange                string
	ExchangeDeclareDurable  *bool
	ExchangeDeclareEnabled  *bool
	ExchangeDeclareType     *string
	Expiration              *string
	Immediate               *bool
	Key                     string
	Mandatory               *bool
	MaxInFlight             *int64
	MessageID               *string
	MetadataExcludePrefixes []string
	NotifyPolicy            *string
	Persistent              *bool
	Priority                *string
	ReplyTo                 *string
	Timeout                 *string
	TLSEnableRenegotiation  *bool
	TLSSkipCertVerify       *bool
	Type                    *string
	URLs                    []string
	UserID                  *string
}

type Amqp09OutletUpdateRequest = Amqp09OutletCreateRequest

type AmqpTLSFeatureRequest struct {
	CACertFilename      string
	CACertContentBase64 string
}

type AmqpMTLSFeatureRequest struct {
	CACertFilename          string
	CACertContentBase64     string
	ClientCertFilename      string
	ClientCertContentBase64 string
	ClientKeyFilename       string
	ClientKeyContentBase64  string
}

func (c *Client) CreateAmqp09Outlet(ctx context.Context, req Amqp09OutletCreateRequest) (Amqp09Outlet, error) {
	fields := map[string]any{
		"amqp_0_9.exchange": req.Exchange,
		"amqp_0_9.key":      req.Key,
		"amqp_0_9.urls":     req.URLs,
	}
	setOptionalStringField(fields, "amqp_0_9.app_id", req.AppID)
	setOptionalStringField(fields, "amqp_0_9.content_encoding", req.ContentEncoding)
	setOptionalStringField(fields, "amqp_0_9.content_type", req.ContentType)
	setOptionalStringField(fields, "amqp_0_9.correlation_id", req.CorrelationID)
	setOptionalBoolField(fields, "amqp_0_9.exchange_declare.durable", req.ExchangeDeclareDurable)
	setOptionalBoolField(fields, "amqp_0_9.exchange_declare.enabled", req.ExchangeDeclareEnabled)
	setOptionalStringField(fields, "amqp_0_9.exchange_declare.type", req.ExchangeDeclareType)
	setOptionalStringField(fields, "amqp_0_9.expiration", req.Expiration)
	setOptionalBoolField(fields, "amqp_0_9.immediate", req.Immediate)
	setOptionalBoolField(fields, "amqp_0_9.mandatory", req.Mandatory)
	setOptionalInt64Field(fields, "amqp_0_9.max_in_flight", req.MaxInFlight)
	setOptionalStringField(fields, "amqp_0_9.message_id", req.MessageID)
	setOptionalStringListField(fields, "amqp_0_9.metadata.exclude_prefixes", req.MetadataExcludePrefixes)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalBoolField(fields, "amqp_0_9.persistent", req.Persistent)
	setOptionalStringField(fields, "amqp_0_9.priority", req.Priority)
	setOptionalStringField(fields, "amqp_0_9.reply_to", req.ReplyTo)
	setOptionalStringField(fields, "amqp_0_9.timeout", req.Timeout)
	setOptionalBoolField(fields, "amqp_0_9.tls.enable_renegotiation", req.TLSEnableRenegotiation)
	setOptionalBoolField(fields, "amqp_0_9.tls.skip_cert_verify", req.TLSSkipCertVerify)
	setOptionalStringField(fields, "amqp_0_9.type", req.Type)
	setOptionalStringField(fields, "amqp_0_9.user_id", req.UserID)

	var response typedOutletResponse
	if err := c.doMultipartPayload(ctx, httpMethodPost, "/v1/terraform/amqp_0_9", typedOutletCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, &response, statusCreated); err != nil {
		return Amqp09Outlet{}, err
	}

	return decodeAmqp09Outlet(response)
}

func (c *Client) GetAmqp09Outlet(ctx context.Context, id int64) (Amqp09Outlet, error) {
	var response typedOutletResponse
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/amqp_0_9/%d", id), nil, &response, statusOK); err != nil {
		return Amqp09Outlet{}, err
	}

	return decodeAmqp09Outlet(response)
}

func (c *Client) UpdateAmqp09Outlet(ctx context.Context, id int64, req Amqp09OutletUpdateRequest) error {
	fields := map[string]any{
		"amqp_0_9.exchange": req.Exchange,
		"amqp_0_9.key":      req.Key,
		"amqp_0_9.urls":     req.URLs,
	}
	setOptionalStringField(fields, "amqp_0_9.app_id", req.AppID)
	setOptionalStringField(fields, "amqp_0_9.content_encoding", req.ContentEncoding)
	setOptionalStringField(fields, "amqp_0_9.content_type", req.ContentType)
	setOptionalStringField(fields, "amqp_0_9.correlation_id", req.CorrelationID)
	setOptionalBoolField(fields, "amqp_0_9.exchange_declare.durable", req.ExchangeDeclareDurable)
	setOptionalBoolField(fields, "amqp_0_9.exchange_declare.enabled", req.ExchangeDeclareEnabled)
	setOptionalStringField(fields, "amqp_0_9.exchange_declare.type", req.ExchangeDeclareType)
	setOptionalStringField(fields, "amqp_0_9.expiration", req.Expiration)
	setOptionalBoolField(fields, "amqp_0_9.immediate", req.Immediate)
	setOptionalBoolField(fields, "amqp_0_9.mandatory", req.Mandatory)
	setOptionalInt64Field(fields, "amqp_0_9.max_in_flight", req.MaxInFlight)
	setOptionalStringField(fields, "amqp_0_9.message_id", req.MessageID)
	setOptionalStringListField(fields, "amqp_0_9.metadata.exclude_prefixes", req.MetadataExcludePrefixes)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalBoolField(fields, "amqp_0_9.persistent", req.Persistent)
	setOptionalStringField(fields, "amqp_0_9.priority", req.Priority)
	setOptionalStringField(fields, "amqp_0_9.reply_to", req.ReplyTo)
	setOptionalStringField(fields, "amqp_0_9.timeout", req.Timeout)
	setOptionalBoolField(fields, "amqp_0_9.tls.enable_renegotiation", req.TLSEnableRenegotiation)
	setOptionalBoolField(fields, "amqp_0_9.tls.skip_cert_verify", req.TLSSkipCertVerify)
	setOptionalStringField(fields, "amqp_0_9.type", req.Type)
	setOptionalStringField(fields, "amqp_0_9.user_id", req.UserID)

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/amqp_0_9/%d", id), typedOutletUpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, nil, statusOK)
}

func (c *Client) DeleteAmqp09Outlet(ctx context.Context, id int64) error {
	_, err := c.DeleteAmqp09OutletOperation(ctx, id)
	return err
}

func (c *Client) DeleteAmqp09OutletOperation(ctx context.Context, id int64) (OutletOperation, error) {
	return c.deleteTypedOutlet(ctx, fmt.Sprintf("/v1/terraform/amqp_0_9/%d", id))
}

func (c *Client) SetAmqp09Field(ctx context.Context, id int64, field string, value any) error {
	return c.patchOutletField(ctx, fmt.Sprintf("/v1/terraform/amqp_0_9/%d/fields/%s", id, field), value)
}

func (c *Client) DeleteAmqp09Field(ctx context.Context, id int64, field string) error {
	return c.deleteOutletField(ctx, fmt.Sprintf("/v1/terraform/amqp_0_9/%d/fields/%s", id, field))
}

func (c *Client) ListAmqp09Credentials(ctx context.Context, outletID int64) ([]OutletMaterialSummary, error) {
	return c.listOutletMaterials(ctx, fmt.Sprintf("/v1/terraform/amqp_0_9/%d/credentials", outletID))
}

func (c *Client) DeleteAmqp09Credential(ctx context.Context, outletID int64, role string) (OutletMaterialDeleteResponse, error) {
	return c.deleteOutletMaterial(ctx, fmt.Sprintf("/v1/terraform/amqp_0_9/%d/credentials/%s", outletID, role))
}

func (c *Client) GetAmqp09TLSFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	var response OutletFeatureDetail
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/amqp_0_9/%d/features/%s", outletID, amqp09TLSFeatureID), nil, &response, statusOK); err != nil {
		return OutletFeatureDetail{}, err
	}
	return response, nil
}

func (c *Client) ConfigureAmqp09TLSFeature(ctx context.Context, outletID int64, req AmqpTLSFeatureRequest) error {
	content, err := decodeBase64File(req.CACertContentBase64)
	if err != nil {
		return err
	}
	return c.doMultipartPayload(ctx, httpMethodPost, fmt.Sprintf("/v1/terraform/amqp_0_9/%d/features/%s", outletID, amqp09TLSFeatureID), map[string]any{
		"inputs": map[string]any{
			amqpTLSCARole: payloadFileRef{File: amqpTLSCAPartName},
		},
	}, []multipartFilePart{{FieldName: amqpTLSCAPartName, FileName: req.CACertFilename, Content: content}}, nil, statusOK)
}

func (c *Client) UpdateAmqp09TLSFeature(ctx context.Context, outletID int64, req AmqpTLSFeatureRequest) error {
	content, err := decodeBase64File(req.CACertContentBase64)
	if err != nil {
		return err
	}
	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/amqp_0_9/%d/features/%s", outletID, amqp09TLSFeatureID), map[string]any{
		"inputs": map[string]any{
			amqpTLSCARole: payloadFileRef{File: amqpTLSCAPartName},
		},
	}, []multipartFilePart{{FieldName: amqpTLSCAPartName, FileName: req.CACertFilename, Content: content}}, nil, statusOK)
}

func (c *Client) DeleteAmqp09TLSFeature(ctx context.Context, outletID int64) error {
	return c.doJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/terraform/amqp_0_9/%d/features/%s", outletID, amqp09TLSFeatureID), nil, nil, statusOK)
}

func (c *Client) GetAmqp09MTLSFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	var response OutletFeatureDetail
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/amqp_0_9/%d/features/%s", outletID, amqp09MTLSFeatureID), nil, &response, statusOK); err != nil {
		return OutletFeatureDetail{}, err
	}
	return response, nil
}

func (c *Client) ConfigureAmqp09MTLSFeature(ctx context.Context, outletID int64, req AmqpMTLSFeatureRequest) error {
	caCert, err := decodeBase64File(req.CACertContentBase64)
	if err != nil {
		return err
	}
	clientCert, err := decodeBase64File(req.ClientCertContentBase64)
	if err != nil {
		return err
	}
	clientKey, err := decodeBase64File(req.ClientKeyContentBase64)
	if err != nil {
		return err
	}
	return c.doMultipartPayload(ctx, httpMethodPost, fmt.Sprintf("/v1/terraform/amqp_0_9/%d/features/%s", outletID, amqp09MTLSFeatureID), map[string]any{
		"inputs": map[string]any{
			amqpTLSCARole:          payloadFileRef{File: amqpTLSCAPartName},
			amqpMTLSClientCertRole: payloadFileRef{File: amqpMTLSClientCertPart},
			amqpMTLSClientKeyRole:  payloadFileRef{File: amqpMTLSClientKeyPart},
		},
	}, []multipartFilePart{
		{FieldName: amqpTLSCAPartName, FileName: req.CACertFilename, Content: caCert},
		{FieldName: amqpMTLSClientCertPart, FileName: req.ClientCertFilename, Content: clientCert},
		{FieldName: amqpMTLSClientKeyPart, FileName: req.ClientKeyFilename, Content: clientKey},
	}, nil, statusOK)
}

func (c *Client) UpdateAmqp09MTLSFeature(ctx context.Context, outletID int64, req AmqpMTLSFeatureRequest) error {
	caCert, err := decodeBase64File(req.CACertContentBase64)
	if err != nil {
		return err
	}
	clientCert, err := decodeBase64File(req.ClientCertContentBase64)
	if err != nil {
		return err
	}
	clientKey, err := decodeBase64File(req.ClientKeyContentBase64)
	if err != nil {
		return err
	}
	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/amqp_0_9/%d/features/%s", outletID, amqp09MTLSFeatureID), map[string]any{
		"inputs": map[string]any{
			amqpTLSCARole:          payloadFileRef{File: amqpTLSCAPartName},
			amqpMTLSClientCertRole: payloadFileRef{File: amqpMTLSClientCertPart},
			amqpMTLSClientKeyRole:  payloadFileRef{File: amqpMTLSClientKeyPart},
		},
	}, []multipartFilePart{
		{FieldName: amqpTLSCAPartName, FileName: req.CACertFilename, Content: caCert},
		{FieldName: amqpMTLSClientCertPart, FileName: req.ClientCertFilename, Content: clientCert},
		{FieldName: amqpMTLSClientKeyPart, FileName: req.ClientKeyFilename, Content: clientKey},
	}, nil, statusOK)
}

func (c *Client) DeleteAmqp09MTLSFeature(ctx context.Context, outletID int64) error {
	return c.doJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/terraform/amqp_0_9/%d/features/%s", outletID, amqp09MTLSFeatureID), nil, nil, statusOK)
}

func decodeAmqp09Outlet(response typedOutletResponse) (Amqp09Outlet, error) {
	exchange, err := requiredStringField(response.Config, "amqp_0_9.exchange")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	key, err := requiredStringField(response.Config, "amqp_0_9.key")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	appID, err := stringFieldPointer(response.Config, "amqp_0_9.app_id")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	contentEncoding, err := stringFieldPointer(response.Config, "amqp_0_9.content_encoding")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	contentType, err := stringFieldPointer(response.Config, "amqp_0_9.content_type")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	correlationID, err := stringFieldPointer(response.Config, "amqp_0_9.correlation_id")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	exchangeDeclareDurable, err := boolFieldPointer(response.Config, "amqp_0_9.exchange_declare.durable")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	exchangeDeclareEnabled, err := boolFieldPointer(response.Config, "amqp_0_9.exchange_declare.enabled")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	exchangeDeclareType, err := stringFieldPointer(response.Config, "amqp_0_9.exchange_declare.type")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	expiration, err := stringFieldPointer(response.Config, "amqp_0_9.expiration")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	immediate, err := boolFieldPointer(response.Config, "amqp_0_9.immediate")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	mandatory, err := boolFieldPointer(response.Config, "amqp_0_9.mandatory")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	maxInFlight, err := int64FieldPointer(response.Config, "amqp_0_9.max_in_flight")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	messageID, err := stringFieldPointer(response.Config, "amqp_0_9.message_id")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	metadataExcludePrefixes, err := stringListField(response.Config, "amqp_0_9.metadata.exclude_prefixes")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	notifyPolicy, err := notifyPolicyPointerPreservingEmpty(response.Config)
	if err != nil {
		return Amqp09Outlet{}, err
	}
	persistent, err := boolFieldPointer(response.Config, "amqp_0_9.persistent")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	priority, err := stringFieldPointer(response.Config, "amqp_0_9.priority")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	replyTo, err := stringFieldPointer(response.Config, "amqp_0_9.reply_to")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	timeout, err := stringFieldPointer(response.Config, "amqp_0_9.timeout")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	tlsEnableRenegotiation, err := boolFieldPointer(response.Config, "amqp_0_9.tls.enable_renegotiation")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	tlsSkipCertVerify, err := boolFieldPointer(response.Config, "amqp_0_9.tls.skip_cert_verify")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	amqpType, err := stringFieldPointer(response.Config, "amqp_0_9.type")
	if err != nil {
		return Amqp09Outlet{}, err
	}
	userID, err := stringFieldPointer(response.Config, "amqp_0_9.user_id")
	if err != nil {
		return Amqp09Outlet{}, err
	}

	return Amqp09Outlet{
		ID:                      response.ID,
		Name:                    response.Name,
		Description:             response.Description,
		Enabled:                 response.Enabled,
		AppID:                   appID,
		ContentEncoding:         contentEncoding,
		ContentType:             contentType,
		CorrelationID:           correlationID,
		Exchange:                exchange,
		ExchangeDeclareDurable:  exchangeDeclareDurable,
		ExchangeDeclareEnabled:  exchangeDeclareEnabled,
		ExchangeDeclareType:     exchangeDeclareType,
		Expiration:              expiration,
		Immediate:               immediate,
		Key:                     key,
		Mandatory:               mandatory,
		MaxInFlight:             maxInFlight,
		MessageID:               messageID,
		MetadataExcludePrefixes: metadataExcludePrefixes,
		NotifyPolicy:            notifyPolicy,
		Persistent:              persistent,
		Priority:                priority,
		ReplyTo:                 replyTo,
		Timeout:                 timeout,
		TLSEnableRenegotiation:  tlsEnableRenegotiation,
		TLSSkipCertVerify:       tlsSkipCertVerify,
		Type:                    amqpType,
		UserID:                  userID,
		ApplyStatus:             response.ApplyStatus,
		LastApplyError:          response.LastApplyError,
		ConfigGeneration:        response.ConfigGeneration,
		AppliedGeneration:       response.AppliedGeneration,
	}, nil
}

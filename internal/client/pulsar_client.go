package client

import (
	"context"
	"fmt"
)

const (
	pulsarTLSFeatureID             = "pulsar.tls"
	pulsarOAuth2FeatureID          = "pulsar.oauth2"
	pulsarTokenFeatureID           = "pulsar.token"
	pulsarTLSCARole                = "ca_cert"
	pulsarOAuth2PrivateKeyRole     = "oauth2_private_key"
	pulsarTokenRole                = "pulsar_access_token"
	pulsarTLSCAPartName            = "ca_cert_part"
	pulsarOAuth2PrivateKeyPartName = "oauth2_private_key_part"
)

type PulsarOutlet struct {
	ID                int64
	Name              string
	Description       string
	Enabled           bool
	Key               *string
	MaxInFlight       *int64
	NotifyPolicy      *string
	OrderingKey       *string
	Topic             string
	URL               string
	ApplyStatus       string
	LastApplyError    string
	ConfigGeneration  int64
	AppliedGeneration int64
}

type PulsarOutletCreateRequest struct {
	Name         string
	Description  string
	Enabled      bool
	Key          *string
	MaxInFlight  *int64
	NotifyPolicy *string
	OrderingKey  *string
	Topic        string
	URL          string
}

type PulsarOutletUpdateRequest = PulsarOutletCreateRequest

type PulsarTLSFeatureRequest struct {
	CACertFilename      string
	CACertContentBase64 string
}

type PulsarOAuth2FeatureRequest struct {
	IssuerURL               string
	Audience                string
	PrivateKeyFilename      string
	PrivateKeyContentBase64 string
}

type PulsarTokenFeatureRequest struct {
	AccessToken string
}

func (c *Client) CreatePulsarOutlet(ctx context.Context, req PulsarOutletCreateRequest) (PulsarOutlet, error) {
	fields := map[string]any{
		"pulsar.topic": req.Topic,
		"pulsar.url":   req.URL,
	}
	setOptionalStringField(fields, "pulsar.key", req.Key)
	setOptionalInt64Field(fields, "pulsar.max_in_flight", req.MaxInFlight)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalStringField(fields, "pulsar.ordering_key", req.OrderingKey)

	var response typedOutletResponse
	if err := c.doMultipartPayload(ctx, httpMethodPost, "/v1/terraform/pulsar", typedOutletCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, &response, statusCreated); err != nil {
		return PulsarOutlet{}, err
	}

	return decodePulsarOutlet(response)
}

func (c *Client) GetPulsarOutlet(ctx context.Context, id int64) (PulsarOutlet, error) {
	var response typedOutletResponse
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/pulsar/%d", id), nil, &response, statusOK); err != nil {
		return PulsarOutlet{}, err
	}

	return decodePulsarOutlet(response)
}

func (c *Client) UpdatePulsarOutlet(ctx context.Context, id int64, req PulsarOutletUpdateRequest) error {
	fields := map[string]any{
		"pulsar.topic": req.Topic,
		"pulsar.url":   req.URL,
	}
	setOptionalStringField(fields, "pulsar.key", req.Key)
	setOptionalInt64Field(fields, "pulsar.max_in_flight", req.MaxInFlight)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalStringField(fields, "pulsar.ordering_key", req.OrderingKey)

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/pulsar/%d", id), typedOutletUpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, nil, statusOK)
}

func (c *Client) DeletePulsarOutlet(ctx context.Context, id int64) error {
	_, err := c.DeletePulsarOutletOperation(ctx, id)
	return err
}

func (c *Client) DeletePulsarOutletOperation(ctx context.Context, id int64) (OutletOperation, error) {
	return c.deleteTypedOutlet(ctx, fmt.Sprintf("/v1/terraform/pulsar/%d", id))
}

func (c *Client) SetPulsarField(ctx context.Context, id int64, field string, value any) error {
	return c.patchOutletField(ctx, fmt.Sprintf("/v1/terraform/pulsar/%d/fields/%s", id, field), value)
}

func (c *Client) DeletePulsarField(ctx context.Context, id int64, field string) error {
	return c.deleteOutletField(ctx, fmt.Sprintf("/v1/terraform/pulsar/%d/fields/%s", id, field))
}

func (c *Client) ListPulsarCredentials(ctx context.Context, outletID int64) ([]OutletMaterialSummary, error) {
	return c.listOutletMaterials(ctx, fmt.Sprintf("/v1/terraform/pulsar/%d/credentials", outletID))
}

func (c *Client) DeletePulsarCredential(ctx context.Context, outletID int64, role string) (OutletMaterialDeleteResponse, error) {
	return c.deleteOutletMaterial(ctx, fmt.Sprintf("/v1/terraform/pulsar/%d/credentials/%s", outletID, role))
}

func (c *Client) GetPulsarTLSFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	var response OutletFeatureDetail
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/pulsar/%d/features/%s", outletID, pulsarTLSFeatureID), nil, &response, statusOK); err != nil {
		return OutletFeatureDetail{}, err
	}

	return response, nil
}

func (c *Client) ConfigurePulsarTLSFeature(ctx context.Context, outletID int64, req PulsarTLSFeatureRequest) error {
	content, err := decodeBase64File(req.CACertContentBase64)
	if err != nil {
		return err
	}

	return c.doMultipartPayload(ctx, httpMethodPost, fmt.Sprintf("/v1/terraform/pulsar/%d/features/%s", outletID, pulsarTLSFeatureID), map[string]any{
		"inputs": map[string]any{
			pulsarTLSCARole: payloadFileRef{File: pulsarTLSCAPartName},
		},
	}, []multipartFilePart{{
		FieldName: pulsarTLSCAPartName,
		FileName:  req.CACertFilename,
		Content:   content,
	}}, nil, statusOK)
}

func (c *Client) UpdatePulsarTLSFeature(ctx context.Context, outletID int64, req PulsarTLSFeatureRequest) error {
	content, err := decodeBase64File(req.CACertContentBase64)
	if err != nil {
		return err
	}

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/pulsar/%d/features/%s", outletID, pulsarTLSFeatureID), map[string]any{
		"inputs": map[string]any{
			pulsarTLSCARole: payloadFileRef{File: pulsarTLSCAPartName},
		},
	}, []multipartFilePart{{
		FieldName: pulsarTLSCAPartName,
		FileName:  req.CACertFilename,
		Content:   content,
	}}, nil, statusOK)
}

func (c *Client) DeletePulsarTLSFeature(ctx context.Context, outletID int64) error {
	return c.doJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/terraform/pulsar/%d/features/%s", outletID, pulsarTLSFeatureID), nil, nil, statusOK)
}

func (c *Client) GetPulsarOAuth2Feature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	var response OutletFeatureDetail
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/pulsar/%d/features/%s", outletID, pulsarOAuth2FeatureID), nil, &response, statusOK); err != nil {
		return OutletFeatureDetail{}, err
	}

	return response, nil
}

func (c *Client) ConfigurePulsarOAuth2Feature(ctx context.Context, outletID int64, req PulsarOAuth2FeatureRequest) error {
	content, err := decodeBase64File(req.PrivateKeyContentBase64)
	if err != nil {
		return err
	}

	return c.doMultipartPayload(ctx, httpMethodPost, fmt.Sprintf("/v1/terraform/pulsar/%d/features/%s", outletID, pulsarOAuth2FeatureID), map[string]any{
		"inputs": map[string]any{
			"issuer_url":               req.IssuerURL,
			"audience":                 req.Audience,
			pulsarOAuth2PrivateKeyRole: payloadFileRef{File: pulsarOAuth2PrivateKeyPartName},
		},
	}, []multipartFilePart{{
		FieldName: pulsarOAuth2PrivateKeyPartName,
		FileName:  req.PrivateKeyFilename,
		Content:   content,
	}}, nil, statusOK)
}

func (c *Client) UpdatePulsarOAuth2Feature(ctx context.Context, outletID int64, req PulsarOAuth2FeatureRequest) error {
	content, err := decodeBase64File(req.PrivateKeyContentBase64)
	if err != nil {
		return err
	}

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/pulsar/%d/features/%s", outletID, pulsarOAuth2FeatureID), map[string]any{
		"inputs": map[string]any{
			"issuer_url":               req.IssuerURL,
			"audience":                 req.Audience,
			pulsarOAuth2PrivateKeyRole: payloadFileRef{File: pulsarOAuth2PrivateKeyPartName},
		},
	}, []multipartFilePart{{
		FieldName: pulsarOAuth2PrivateKeyPartName,
		FileName:  req.PrivateKeyFilename,
		Content:   content,
	}}, nil, statusOK)
}

func (c *Client) DeletePulsarOAuth2Feature(ctx context.Context, outletID int64) error {
	return c.doJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/terraform/pulsar/%d/features/%s", outletID, pulsarOAuth2FeatureID), nil, nil, statusOK)
}

func (c *Client) GetPulsarTokenFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	var response OutletFeatureDetail
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/pulsar/%d/features/%s", outletID, pulsarTokenFeatureID), nil, &response, statusOK); err != nil {
		return OutletFeatureDetail{}, err
	}

	return response, nil
}

func (c *Client) ConfigurePulsarTokenFeature(ctx context.Context, outletID int64, req PulsarTokenFeatureRequest) error {
	return c.doMultipartPayload(ctx, httpMethodPost, fmt.Sprintf("/v1/terraform/pulsar/%d/features/%s", outletID, pulsarTokenFeatureID), map[string]any{
		"inputs": map[string]any{
			pulsarTokenRole: req.AccessToken,
		},
	}, nil, nil, statusOK)
}

func (c *Client) UpdatePulsarTokenFeature(ctx context.Context, outletID int64, req PulsarTokenFeatureRequest) error {
	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/pulsar/%d/features/%s", outletID, pulsarTokenFeatureID), map[string]any{
		"inputs": map[string]any{
			pulsarTokenRole: req.AccessToken,
		},
	}, nil, nil, statusOK)
}

func (c *Client) DeletePulsarTokenFeature(ctx context.Context, outletID int64) error {
	return c.doJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/terraform/pulsar/%d/features/%s", outletID, pulsarTokenFeatureID), nil, nil, statusOK)
}

func decodePulsarOutlet(response typedOutletResponse) (PulsarOutlet, error) {
	topic, err := requiredStringField(response.Config, "pulsar.topic")
	if err != nil {
		return PulsarOutlet{}, err
	}
	url, err := requiredStringField(response.Config, "pulsar.url")
	if err != nil {
		return PulsarOutlet{}, err
	}
	key, err := stringFieldPointer(response.Config, "pulsar.key")
	if err != nil {
		return PulsarOutlet{}, err
	}
	maxInFlight, err := int64FieldPointer(response.Config, "pulsar.max_in_flight")
	if err != nil {
		return PulsarOutlet{}, err
	}
	notifyPolicy, err := notifyPolicyPointerPreservingEmpty(response.Config)
	if err != nil {
		return PulsarOutlet{}, err
	}
	orderingKey, err := stringFieldPointer(response.Config, "pulsar.ordering_key")
	if err != nil {
		return PulsarOutlet{}, err
	}

	return PulsarOutlet{
		ID:                response.ID,
		Name:              response.Name,
		Description:       response.Description,
		Enabled:           response.Enabled,
		Key:               key,
		MaxInFlight:       maxInFlight,
		NotifyPolicy:      notifyPolicy,
		OrderingKey:       orderingKey,
		Topic:             topic,
		URL:               url,
		ApplyStatus:       response.ApplyStatus,
		LastApplyError:    response.LastApplyError,
		ConfigGeneration:  response.ConfigGeneration,
		AppliedGeneration: response.AppliedGeneration,
	}, nil
}

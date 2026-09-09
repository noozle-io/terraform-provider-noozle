package client

import (
	"context"
	"fmt"
)

const (
	amqp1URLsRole             = "amqp_1_urls"
	amqp1TLSFeatureID         = "amqp_1.tls"
	amqp1MTLSFeatureID        = "amqp_1.mtls"
	amqp1SASLAnonymousFeature = "amqp_1.sasl.anonymous"
	amqp1SASLPlainFeature     = "amqp_1.sasl.plain"
	amqp1SASLPasswordRole     = "sasl_password"
)

type Amqp1Outlet struct {
	ID                       int64
	Name                     string
	Description              string
	Enabled                  bool
	ApplicationPropertiesMap *string
	MaxInFlight              *int64
	MetadataExcludePrefixes  []string
	NotifyPolicy             *string
	TargetAddress            string
	TLSEnableRenegotiation   *bool
	TLSSkipCertVerify        *bool
	ApplyStatus              string
	LastApplyError           string
	ConfigGeneration         int64
	AppliedGeneration        int64
}

type Amqp1OutletCreateRequest struct {
	Name                     string
	Description              string
	Enabled                  bool
	ApplicationPropertiesMap *string
	MaxInFlight              *int64
	MetadataExcludePrefixes  []string
	NotifyPolicy             *string
	TargetAddress            string
	TLSEnableRenegotiation   *bool
	TLSSkipCertVerify        *bool
	URLs                     []string
}

type Amqp1OutletUpdateRequest = Amqp1OutletCreateRequest

type Amqp1SASLPlainFeatureRequest struct {
	User     string
	Password string
}

func (c *Client) CreateAmqp1Outlet(ctx context.Context, req Amqp1OutletCreateRequest) (Amqp1Outlet, error) {
	fields := map[string]any{
		"amqp_1.target_address": req.TargetAddress,
		"amqp_1.urls":           req.URLs,
	}
	setOptionalStringField(fields, "amqp_1.application_properties_map", req.ApplicationPropertiesMap)
	setOptionalInt64Field(fields, "amqp_1.max_in_flight", req.MaxInFlight)
	setOptionalStringListField(fields, "amqp_1.metadata.exclude_prefixes", req.MetadataExcludePrefixes)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalBoolField(fields, "amqp_1.tls.enable_renegotiation", req.TLSEnableRenegotiation)
	setOptionalBoolField(fields, "amqp_1.tls.skip_cert_verify", req.TLSSkipCertVerify)

	var response typedOutletResponse
	if err := c.doMultipartPayload(ctx, httpMethodPost, "/v1/terraform/amqp_1", typedOutletCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, &response, statusCreated); err != nil {
		return Amqp1Outlet{}, err
	}
	return decodeAmqp1Outlet(response)
}

func (c *Client) GetAmqp1Outlet(ctx context.Context, id int64) (Amqp1Outlet, error) {
	var response typedOutletResponse
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/amqp_1/%d", id), nil, &response, statusOK); err != nil {
		return Amqp1Outlet{}, err
	}
	return decodeAmqp1Outlet(response)
}

func (c *Client) UpdateAmqp1Outlet(ctx context.Context, id int64, req Amqp1OutletUpdateRequest) error {
	fields := map[string]any{
		"amqp_1.target_address": req.TargetAddress,
		"amqp_1.urls":           req.URLs,
	}
	setOptionalStringField(fields, "amqp_1.application_properties_map", req.ApplicationPropertiesMap)
	setOptionalInt64Field(fields, "amqp_1.max_in_flight", req.MaxInFlight)
	setOptionalStringListField(fields, "amqp_1.metadata.exclude_prefixes", req.MetadataExcludePrefixes)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalBoolField(fields, "amqp_1.tls.enable_renegotiation", req.TLSEnableRenegotiation)
	setOptionalBoolField(fields, "amqp_1.tls.skip_cert_verify", req.TLSSkipCertVerify)

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/amqp_1/%d", id), typedOutletUpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, nil, statusOK)
}

func (c *Client) DeleteAmqp1Outlet(ctx context.Context, id int64) error {
	_, err := c.DeleteAmqp1OutletOperation(ctx, id)
	return err
}

func (c *Client) DeleteAmqp1OutletOperation(ctx context.Context, id int64) (OutletOperation, error) {
	return c.deleteTypedOutlet(ctx, fmt.Sprintf("/v1/terraform/amqp_1/%d", id))
}

func (c *Client) SetAmqp1Field(ctx context.Context, id int64, field string, value any) error {
	return c.patchOutletField(ctx, fmt.Sprintf("/v1/terraform/amqp_1/%d/fields/%s", id, field), value)
}

func (c *Client) DeleteAmqp1Field(ctx context.Context, id int64, field string) error {
	return c.deleteOutletField(ctx, fmt.Sprintf("/v1/terraform/amqp_1/%d/fields/%s", id, field))
}

func (c *Client) ListAmqp1Credentials(ctx context.Context, outletID int64) ([]OutletMaterialSummary, error) {
	return c.listOutletMaterials(ctx, fmt.Sprintf("/v1/terraform/amqp_1/%d/credentials", outletID))
}

func (c *Client) DeleteAmqp1Credential(ctx context.Context, outletID int64, role string) (OutletMaterialDeleteResponse, error) {
	return c.deleteOutletMaterial(ctx, fmt.Sprintf("/v1/terraform/amqp_1/%d/credentials/%s", outletID, role))
}

func (c *Client) GetAmqp1TLSFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	var response OutletFeatureDetail
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1TLSFeatureID), nil, &response, statusOK); err != nil {
		return OutletFeatureDetail{}, err
	}
	return response, nil
}

func (c *Client) ConfigureAmqp1TLSFeature(ctx context.Context, outletID int64, req AmqpTLSFeatureRequest) error {
	content, err := decodeBase64File(req.CACertContentBase64)
	if err != nil {
		return err
	}
	return c.doMultipartPayload(ctx, httpMethodPost, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1TLSFeatureID), map[string]any{
		"inputs": map[string]any{
			amqpTLSCARole: payloadFileRef{File: amqpTLSCAPartName},
		},
	}, []multipartFilePart{{FieldName: amqpTLSCAPartName, FileName: req.CACertFilename, Content: content}}, nil, statusOK)
}

func (c *Client) UpdateAmqp1TLSFeature(ctx context.Context, outletID int64, req AmqpTLSFeatureRequest) error {
	content, err := decodeBase64File(req.CACertContentBase64)
	if err != nil {
		return err
	}
	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1TLSFeatureID), map[string]any{
		"inputs": map[string]any{
			amqpTLSCARole: payloadFileRef{File: amqpTLSCAPartName},
		},
	}, []multipartFilePart{{FieldName: amqpTLSCAPartName, FileName: req.CACertFilename, Content: content}}, nil, statusOK)
}

func (c *Client) DeleteAmqp1TLSFeature(ctx context.Context, outletID int64) error {
	return c.doJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1TLSFeatureID), nil, nil, statusOK)
}

func (c *Client) GetAmqp1MTLSFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	var response OutletFeatureDetail
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1MTLSFeatureID), nil, &response, statusOK); err != nil {
		return OutletFeatureDetail{}, err
	}
	return response, nil
}

func (c *Client) ConfigureAmqp1MTLSFeature(ctx context.Context, outletID int64, req AmqpMTLSFeatureRequest) error {
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
	return c.doMultipartPayload(ctx, httpMethodPost, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1MTLSFeatureID), map[string]any{
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

func (c *Client) UpdateAmqp1MTLSFeature(ctx context.Context, outletID int64, req AmqpMTLSFeatureRequest) error {
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
	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1MTLSFeatureID), map[string]any{
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

func (c *Client) DeleteAmqp1MTLSFeature(ctx context.Context, outletID int64) error {
	return c.doJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1MTLSFeatureID), nil, nil, statusOK)
}

func (c *Client) GetAmqp1SASLAnonymousFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	var response OutletFeatureDetail
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1SASLAnonymousFeature), nil, &response, statusOK); err != nil {
		return OutletFeatureDetail{}, err
	}
	return response, nil
}

func (c *Client) ConfigureAmqp1SASLAnonymousFeature(ctx context.Context, outletID int64) error {
	return c.doMultipartPayload(ctx, httpMethodPost, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1SASLAnonymousFeature), map[string]any{
		"inputs": map[string]any{},
	}, nil, nil, statusOK)
}

func (c *Client) UpdateAmqp1SASLAnonymousFeature(ctx context.Context, outletID int64) error {
	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1SASLAnonymousFeature), map[string]any{
		"inputs": map[string]any{},
	}, nil, nil, statusOK)
}

func (c *Client) DeleteAmqp1SASLAnonymousFeature(ctx context.Context, outletID int64) error {
	return c.doJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1SASLAnonymousFeature), nil, nil, statusOK)
}

func (c *Client) GetAmqp1SASLPlainFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	var response OutletFeatureDetail
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1SASLPlainFeature), nil, &response, statusOK); err != nil {
		return OutletFeatureDetail{}, err
	}
	return response, nil
}

func (c *Client) ConfigureAmqp1SASLPlainFeature(ctx context.Context, outletID int64, req Amqp1SASLPlainFeatureRequest) error {
	return c.doMultipartPayload(ctx, httpMethodPost, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1SASLPlainFeature), map[string]any{
		"inputs": map[string]any{
			"user":                req.User,
			amqp1SASLPasswordRole: req.Password,
		},
	}, nil, nil, statusOK)
}

func (c *Client) UpdateAmqp1SASLPlainFeature(ctx context.Context, outletID int64, req Amqp1SASLPlainFeatureRequest) error {
	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1SASLPlainFeature), map[string]any{
		"inputs": map[string]any{
			"user":                req.User,
			amqp1SASLPasswordRole: req.Password,
		},
	}, nil, nil, statusOK)
}

func (c *Client) DeleteAmqp1SASLPlainFeature(ctx context.Context, outletID int64) error {
	return c.doJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/terraform/amqp_1/%d/features/%s", outletID, amqp1SASLPlainFeature), nil, nil, statusOK)
}

func decodeAmqp1Outlet(response typedOutletResponse) (Amqp1Outlet, error) {
	targetAddress, err := requiredStringField(response.Config, "amqp_1.target_address")
	if err != nil {
		return Amqp1Outlet{}, err
	}
	applicationPropertiesMap, err := stringFieldPointer(response.Config, "amqp_1.application_properties_map")
	if err != nil {
		return Amqp1Outlet{}, err
	}
	maxInFlight, err := int64FieldPointer(response.Config, "amqp_1.max_in_flight")
	if err != nil {
		return Amqp1Outlet{}, err
	}
	metadataExcludePrefixes, err := stringListField(response.Config, "amqp_1.metadata.exclude_prefixes")
	if err != nil {
		return Amqp1Outlet{}, err
	}
	notifyPolicy, err := notifyPolicyPointerPreservingEmpty(response.Config)
	if err != nil {
		return Amqp1Outlet{}, err
	}
	tlsEnableRenegotiation, err := boolFieldPointer(response.Config, "amqp_1.tls.enable_renegotiation")
	if err != nil {
		return Amqp1Outlet{}, err
	}
	tlsSkipCertVerify, err := boolFieldPointer(response.Config, "amqp_1.tls.skip_cert_verify")
	if err != nil {
		return Amqp1Outlet{}, err
	}

	return Amqp1Outlet{
		ID:                       response.ID,
		Name:                     response.Name,
		Description:              response.Description,
		Enabled:                  response.Enabled,
		ApplicationPropertiesMap: applicationPropertiesMap,
		MaxInFlight:              maxInFlight,
		MetadataExcludePrefixes:  metadataExcludePrefixes,
		NotifyPolicy:             notifyPolicy,
		TargetAddress:            targetAddress,
		TLSEnableRenegotiation:   tlsEnableRenegotiation,
		TLSSkipCertVerify:        tlsSkipCertVerify,
		ApplyStatus:              response.ApplyStatus,
		LastApplyError:           response.LastApplyError,
		ConfigGeneration:         response.ConfigGeneration,
		AppliedGeneration:        response.AppliedGeneration,
	}, nil
}

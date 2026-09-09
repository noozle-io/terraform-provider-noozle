package client

import (
	"context"
	"fmt"
)

const (
	natsJetstreamURLsRole            = "nats_jetstream_urls"
	natsJetstreamJWTUserRole         = "nats_user_jwt"
	natsJetstreamNkeyFileRole        = "nkey_file"
	natsJetstreamNkeySeedRole        = "nkey_seed"
	natsJetstreamUserCredentialsRole = "user_credentials"

	natsJetstreamJWTSeedFeatureID         = "nats.auth.jwt_seed"
	natsJetstreamNkeyFeatureID            = "nats.auth.nkey"
	natsJetstreamUserCredentialsFeatureID = "nats.auth.user_credentials"

	natsJetstreamNkeyPartName            = "nkey_file_part"
	natsJetstreamUserCredentialsPartName = "user_credentials_part"
)

type NatsJetstreamOutlet struct {
	ID                      int64
	Name                    string
	Description             string
	Enabled                 bool
	BackoffInitialInterval  *string
	BackoffJitter           *float64
	BackoffMaxElapsedTime   *string
	BackoffMaxInterval      *string
	Headers                 map[string]string
	InjectTracingMap        *string
	MaxInFlight             *int64
	MaxRetries              *int64
	MetadataIncludePatterns []string
	MetadataIncludePrefixes []string
	NotifyPolicy            *string
	Subject                 string
	TLSEnableRenegotiation  *bool
	TLSSkipCertVerify       *bool
	ApplyStatus             string
	LastApplyError          string
	ConfigGeneration        int64
	AppliedGeneration       int64
}

type NatsJetstreamOutletCreateRequest struct {
	Name                    string
	Description             string
	Enabled                 bool
	BackoffInitialInterval  *string
	BackoffJitter           *float64
	BackoffMaxElapsedTime   *string
	BackoffMaxInterval      *string
	Headers                 map[string]string
	InjectTracingMap        *string
	MaxInFlight             *int64
	MaxRetries              *int64
	MetadataIncludePatterns []string
	MetadataIncludePrefixes []string
	NotifyPolicy            *string
	Subject                 string
	TLSEnableRenegotiation  *bool
	TLSSkipCertVerify       *bool
	URLs                    []string
}

type NatsJetstreamOutletUpdateRequest = NatsJetstreamOutletCreateRequest

type NatsJetstreamJWTSeedFeatureRequest struct {
	UserJWT  string
	NkeySeed string
}

type NatsJetstreamNkeyFeatureRequest struct {
	Filename      string
	ContentBase64 string
}

type NatsJetstreamUserCredentialsFeatureRequest struct {
	Filename      string
	ContentBase64 string
}

func (c *Client) CreateNatsJetstreamOutlet(ctx context.Context, req NatsJetstreamOutletCreateRequest) (NatsJetstreamOutlet, error) {
	fields := map[string]any{
		"nats_jetstream.subject": req.Subject,
		"nats_jetstream.urls":    req.URLs,
	}
	setOptionalStringField(fields, "nats_jetstream.backoff.initial_interval", req.BackoffInitialInterval)
	setOptionalFloat64Field(fields, "nats_jetstream.backoff.jitter", req.BackoffJitter)
	setOptionalStringField(fields, "nats_jetstream.backoff.max_elapsed_time", req.BackoffMaxElapsedTime)
	setOptionalStringField(fields, "nats_jetstream.backoff.max_interval", req.BackoffMaxInterval)
	setOptionalStringMapField(fields, "nats_jetstream.headers", req.Headers)
	setOptionalStringField(fields, "nats_jetstream.inject_tracing_map", req.InjectTracingMap)
	setOptionalInt64Field(fields, "nats_jetstream.max_in_flight", req.MaxInFlight)
	setOptionalInt64Field(fields, "nats_jetstream.max_retries", req.MaxRetries)
	setOptionalStringListField(fields, "nats_jetstream.metadata.include_patterns", req.MetadataIncludePatterns)
	setOptionalStringListField(fields, "nats_jetstream.metadata.include_prefixes", req.MetadataIncludePrefixes)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalBoolField(fields, "nats_jetstream.tls.enable_renegotiation", req.TLSEnableRenegotiation)
	setOptionalBoolField(fields, "nats_jetstream.tls.skip_cert_verify", req.TLSSkipCertVerify)

	var response typedOutletResponse
	if err := c.doMultipartPayload(ctx, httpMethodPost, "/v1/terraform/nats_jetstream", typedOutletCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, &response, statusCreated); err != nil {
		return NatsJetstreamOutlet{}, err
	}

	return decodeNatsJetstreamOutlet(response)
}

func (c *Client) GetNatsJetstreamOutlet(ctx context.Context, id int64) (NatsJetstreamOutlet, error) {
	var response typedOutletResponse
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/nats_jetstream/%d", id), nil, &response, statusOK); err != nil {
		return NatsJetstreamOutlet{}, err
	}

	return decodeNatsJetstreamOutlet(response)
}

func (c *Client) UpdateNatsJetstreamOutlet(ctx context.Context, id int64, req NatsJetstreamOutletUpdateRequest) error {
	fields := map[string]any{
		"nats_jetstream.subject": req.Subject,
		"nats_jetstream.urls":    req.URLs,
	}
	setOptionalStringField(fields, "nats_jetstream.backoff.initial_interval", req.BackoffInitialInterval)
	setOptionalFloat64Field(fields, "nats_jetstream.backoff.jitter", req.BackoffJitter)
	setOptionalStringField(fields, "nats_jetstream.backoff.max_elapsed_time", req.BackoffMaxElapsedTime)
	setOptionalStringField(fields, "nats_jetstream.backoff.max_interval", req.BackoffMaxInterval)
	setOptionalStringMapField(fields, "nats_jetstream.headers", req.Headers)
	setOptionalStringField(fields, "nats_jetstream.inject_tracing_map", req.InjectTracingMap)
	setOptionalInt64Field(fields, "nats_jetstream.max_in_flight", req.MaxInFlight)
	setOptionalInt64Field(fields, "nats_jetstream.max_retries", req.MaxRetries)
	setOptionalStringListField(fields, "nats_jetstream.metadata.include_patterns", req.MetadataIncludePatterns)
	setOptionalStringListField(fields, "nats_jetstream.metadata.include_prefixes", req.MetadataIncludePrefixes)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalBoolField(fields, "nats_jetstream.tls.enable_renegotiation", req.TLSEnableRenegotiation)
	setOptionalBoolField(fields, "nats_jetstream.tls.skip_cert_verify", req.TLSSkipCertVerify)

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/nats_jetstream/%d", id), typedOutletUpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, nil, statusOK)
}

func (c *Client) DeleteNatsJetstreamOutlet(ctx context.Context, id int64) error {
	_, err := c.DeleteNatsJetstreamOutletOperation(ctx, id)
	return err
}

func (c *Client) DeleteNatsJetstreamOutletOperation(ctx context.Context, id int64) (OutletOperation, error) {
	return c.deleteTypedOutlet(ctx, fmt.Sprintf("/v1/terraform/nats_jetstream/%d", id))
}

func (c *Client) SetNatsJetstreamField(ctx context.Context, id int64, field string, value any) error {
	return c.patchOutletField(ctx, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/fields/%s", id, field), value)
}

func (c *Client) DeleteNatsJetstreamField(ctx context.Context, id int64, field string) error {
	return c.deleteOutletField(ctx, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/fields/%s", id, field))
}

func (c *Client) ListNatsJetstreamCredentials(ctx context.Context, outletID int64) ([]OutletMaterialSummary, error) {
	return c.listOutletMaterials(ctx, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/credentials", outletID))
}

func (c *Client) DeleteNatsJetstreamCredential(ctx context.Context, outletID int64, role string) (OutletMaterialDeleteResponse, error) {
	return c.deleteOutletMaterial(ctx, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/credentials/%s", outletID, role))
}

func (c *Client) GetNatsJetstreamJWTSeedFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	var response OutletFeatureDetail
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/features/%s", outletID, natsJetstreamJWTSeedFeatureID), nil, &response, statusOK); err != nil {
		return OutletFeatureDetail{}, err
	}
	return response, nil
}

func (c *Client) ConfigureNatsJetstreamJWTSeedFeature(ctx context.Context, outletID int64, req NatsJetstreamJWTSeedFeatureRequest) error {
	return c.doMultipartPayload(ctx, httpMethodPost, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/features/%s", outletID, natsJetstreamJWTSeedFeatureID), map[string]any{
		"inputs": map[string]any{
			natsJetstreamJWTUserRole:  req.UserJWT,
			natsJetstreamNkeySeedRole: req.NkeySeed,
		},
	}, nil, nil, statusOK)
}

func (c *Client) UpdateNatsJetstreamJWTSeedFeature(ctx context.Context, outletID int64, req NatsJetstreamJWTSeedFeatureRequest) error {
	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/features/%s", outletID, natsJetstreamJWTSeedFeatureID), map[string]any{
		"inputs": map[string]any{
			natsJetstreamJWTUserRole:  req.UserJWT,
			natsJetstreamNkeySeedRole: req.NkeySeed,
		},
	}, nil, nil, statusOK)
}

func (c *Client) DeleteNatsJetstreamJWTSeedFeature(ctx context.Context, outletID int64) error {
	return c.doJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/features/%s", outletID, natsJetstreamJWTSeedFeatureID), nil, nil, statusOK)
}

func (c *Client) GetNatsJetstreamNkeyFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	var response OutletFeatureDetail
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/features/%s", outletID, natsJetstreamNkeyFeatureID), nil, &response, statusOK); err != nil {
		return OutletFeatureDetail{}, err
	}
	return response, nil
}

func (c *Client) ConfigureNatsJetstreamNkeyFeature(ctx context.Context, outletID int64, req NatsJetstreamNkeyFeatureRequest) error {
	content, err := decodeBase64File(req.ContentBase64)
	if err != nil {
		return err
	}

	return c.doMultipartPayload(ctx, httpMethodPost, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/features/%s", outletID, natsJetstreamNkeyFeatureID), map[string]any{
		"inputs": map[string]any{
			natsJetstreamNkeyFileRole: payloadFileRef{File: natsJetstreamNkeyPartName},
		},
	}, []multipartFilePart{{
		FieldName: natsJetstreamNkeyPartName,
		FileName:  req.Filename,
		Content:   content,
	}}, nil, statusOK)
}

func (c *Client) UpdateNatsJetstreamNkeyFeature(ctx context.Context, outletID int64, req NatsJetstreamNkeyFeatureRequest) error {
	content, err := decodeBase64File(req.ContentBase64)
	if err != nil {
		return err
	}

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/features/%s", outletID, natsJetstreamNkeyFeatureID), map[string]any{
		"inputs": map[string]any{
			natsJetstreamNkeyFileRole: payloadFileRef{File: natsJetstreamNkeyPartName},
		},
	}, []multipartFilePart{{
		FieldName: natsJetstreamNkeyPartName,
		FileName:  req.Filename,
		Content:   content,
	}}, nil, statusOK)
}

func (c *Client) DeleteNatsJetstreamNkeyFeature(ctx context.Context, outletID int64) error {
	return c.doJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/features/%s", outletID, natsJetstreamNkeyFeatureID), nil, nil, statusOK)
}

func (c *Client) GetNatsJetstreamUserCredentialsFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	var response OutletFeatureDetail
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/features/%s", outletID, natsJetstreamUserCredentialsFeatureID), nil, &response, statusOK); err != nil {
		return OutletFeatureDetail{}, err
	}
	return response, nil
}

func (c *Client) ConfigureNatsJetstreamUserCredentialsFeature(ctx context.Context, outletID int64, req NatsJetstreamUserCredentialsFeatureRequest) error {
	content, err := decodeBase64File(req.ContentBase64)
	if err != nil {
		return err
	}

	return c.doMultipartPayload(ctx, httpMethodPost, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/features/%s", outletID, natsJetstreamUserCredentialsFeatureID), map[string]any{
		"inputs": map[string]any{
			natsJetstreamUserCredentialsRole: payloadFileRef{File: natsJetstreamUserCredentialsPartName},
		},
	}, []multipartFilePart{{
		FieldName: natsJetstreamUserCredentialsPartName,
		FileName:  req.Filename,
		Content:   content,
	}}, nil, statusOK)
}

func (c *Client) UpdateNatsJetstreamUserCredentialsFeature(ctx context.Context, outletID int64, req NatsJetstreamUserCredentialsFeatureRequest) error {
	content, err := decodeBase64File(req.ContentBase64)
	if err != nil {
		return err
	}

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/features/%s", outletID, natsJetstreamUserCredentialsFeatureID), map[string]any{
		"inputs": map[string]any{
			natsJetstreamUserCredentialsRole: payloadFileRef{File: natsJetstreamUserCredentialsPartName},
		},
	}, []multipartFilePart{{
		FieldName: natsJetstreamUserCredentialsPartName,
		FileName:  req.Filename,
		Content:   content,
	}}, nil, statusOK)
}

func (c *Client) DeleteNatsJetstreamUserCredentialsFeature(ctx context.Context, outletID int64) error {
	return c.doJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/terraform/nats_jetstream/%d/features/%s", outletID, natsJetstreamUserCredentialsFeatureID), nil, nil, statusOK)
}

func decodeNatsJetstreamOutlet(response typedOutletResponse) (NatsJetstreamOutlet, error) {
	backoffInitialInterval, err := stringFieldPointer(response.Config, "nats_jetstream.backoff.initial_interval")
	if err != nil {
		return NatsJetstreamOutlet{}, err
	}
	backoffJitter, err := float64FieldPointer(response.Config, "nats_jetstream.backoff.jitter")
	if err != nil {
		return NatsJetstreamOutlet{}, err
	}
	backoffMaxElapsedTime, err := stringFieldPointer(response.Config, "nats_jetstream.backoff.max_elapsed_time")
	if err != nil {
		return NatsJetstreamOutlet{}, err
	}
	backoffMaxInterval, err := stringFieldPointer(response.Config, "nats_jetstream.backoff.max_interval")
	if err != nil {
		return NatsJetstreamOutlet{}, err
	}
	headers, err := stringMapField(response.Config, "nats_jetstream.headers")
	if err != nil {
		return NatsJetstreamOutlet{}, err
	}
	injectTracingMap, err := stringFieldPointer(response.Config, "nats_jetstream.inject_tracing_map")
	if err != nil {
		return NatsJetstreamOutlet{}, err
	}
	maxInFlight, err := int64FieldPointer(response.Config, "nats_jetstream.max_in_flight")
	if err != nil {
		return NatsJetstreamOutlet{}, err
	}
	maxRetries, err := int64FieldPointer(response.Config, "nats_jetstream.max_retries")
	if err != nil {
		return NatsJetstreamOutlet{}, err
	}
	metadataIncludePatterns, err := stringListField(response.Config, "nats_jetstream.metadata.include_patterns")
	if err != nil {
		return NatsJetstreamOutlet{}, err
	}
	metadataIncludePrefixes, err := stringListField(response.Config, "nats_jetstream.metadata.include_prefixes")
	if err != nil {
		return NatsJetstreamOutlet{}, err
	}
	notifyPolicy, err := notifyPolicyPointerPreservingEmpty(response.Config)
	if err != nil {
		return NatsJetstreamOutlet{}, err
	}
	subject, err := requiredStringField(response.Config, "nats_jetstream.subject")
	if err != nil {
		return NatsJetstreamOutlet{}, err
	}
	tlsEnableRenegotiation, err := boolFieldPointer(response.Config, "nats_jetstream.tls.enable_renegotiation")
	if err != nil {
		return NatsJetstreamOutlet{}, err
	}
	tlsSkipCertVerify, err := boolFieldPointer(response.Config, "nats_jetstream.tls.skip_cert_verify")
	if err != nil {
		return NatsJetstreamOutlet{}, err
	}

	return NatsJetstreamOutlet{
		ID:                      response.ID,
		Name:                    response.Name,
		Description:             response.Description,
		Enabled:                 response.Enabled,
		BackoffInitialInterval:  backoffInitialInterval,
		BackoffJitter:           backoffJitter,
		BackoffMaxElapsedTime:   backoffMaxElapsedTime,
		BackoffMaxInterval:      backoffMaxInterval,
		Headers:                 headers,
		InjectTracingMap:        injectTracingMap,
		MaxInFlight:             maxInFlight,
		MaxRetries:              maxRetries,
		MetadataIncludePatterns: metadataIncludePatterns,
		MetadataIncludePrefixes: metadataIncludePrefixes,
		NotifyPolicy:            notifyPolicy,
		Subject:                 subject,
		TLSEnableRenegotiation:  tlsEnableRenegotiation,
		TLSSkipCertVerify:       tlsSkipCertVerify,
		ApplyStatus:             response.ApplyStatus,
		LastApplyError:          response.LastApplyError,
		ConfigGeneration:        response.ConfigGeneration,
		AppliedGeneration:       response.AppliedGeneration,
	}, nil
}

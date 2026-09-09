package client

import (
	"context"
	"fmt"
)

type GcpPubsubOutlet struct {
	ID                                int64
	Name                              string
	Description                       string
	Enabled                           bool
	BatchingByteSize                  *int64
	BatchingCheck                     *string
	BatchingCount                     *int64
	BatchingJitter                    *float64
	BatchingPeriod                    *string
	ByteThreshold                     *int64
	CountThreshold                    *int64
	DelayThreshold                    *string
	Endpoint                          *string
	FlowControlLimitExceededBehavior  *string
	FlowControlMaxOutstandingBytes    *int64
	FlowControlMaxOutstandingMessages *int64
	MaxInFlight                       *int64
	MetadataExcludePrefixes           []string
	NotifyPolicy                      *string
	OrderingKey                       *string
	Project                           string
	PublishTimeout                    *string
	Topic                             string
	ApplyStatus                       string
	LastApplyError                    string
	ConfigGeneration                  int64
	AppliedGeneration                 int64
}

type GcpPubsubOutletCreateRequest struct {
	Name                              string
	Description                       string
	Enabled                           bool
	BatchingByteSize                  *int64
	BatchingCheck                     *string
	BatchingCount                     *int64
	BatchingJitter                    *float64
	BatchingPeriod                    *string
	ByteThreshold                     *int64
	CountThreshold                    *int64
	DelayThreshold                    *string
	Endpoint                          *string
	FlowControlLimitExceededBehavior  *string
	FlowControlMaxOutstandingBytes    *int64
	FlowControlMaxOutstandingMessages *int64
	MaxInFlight                       *int64
	MetadataExcludePrefixes           []string
	NotifyPolicy                      *string
	OrderingKey                       *string
	Project                           string
	PublishTimeout                    *string
	Topic                             string
	CredentialsFilename               *string
	CredentialsContentBase64          *string
}

type GcpPubsubOutletUpdateRequest = GcpPubsubOutletCreateRequest

func (c *Client) CreateGcpPubsubOutlet(ctx context.Context, req GcpPubsubOutletCreateRequest) (GcpPubsubOutlet, error) {
	fields := map[string]any{
		"gcp_pubsub.project": req.Project,
		"gcp_pubsub.topic":   req.Topic,
	}
	setOptionalInt64Field(fields, "gcp_pubsub.batching.byte_size", req.BatchingByteSize)
	setOptionalStringField(fields, "gcp_pubsub.batching.check", req.BatchingCheck)
	setOptionalInt64Field(fields, "gcp_pubsub.batching.count", req.BatchingCount)
	setOptionalFloat64Field(fields, "gcp_pubsub.batching.jitter", req.BatchingJitter)
	setOptionalStringField(fields, "gcp_pubsub.batching.period", req.BatchingPeriod)
	setOptionalInt64Field(fields, "gcp_pubsub.byte_threshold", req.ByteThreshold)
	setOptionalInt64Field(fields, "gcp_pubsub.count_threshold", req.CountThreshold)
	setOptionalStringField(fields, "gcp_pubsub.delay_threshold", req.DelayThreshold)
	setOptionalStringField(fields, "gcp_pubsub.endpoint", req.Endpoint)
	setOptionalStringField(fields, "gcp_pubsub.flow_control.limit_exceeded_behavior", req.FlowControlLimitExceededBehavior)
	setOptionalInt64Field(fields, "gcp_pubsub.flow_control.max_outstanding_bytes", req.FlowControlMaxOutstandingBytes)
	setOptionalInt64Field(fields, "gcp_pubsub.flow_control.max_outstanding_messages", req.FlowControlMaxOutstandingMessages)
	setOptionalInt64Field(fields, "gcp_pubsub.max_in_flight", req.MaxInFlight)
	setOptionalStringListField(fields, "gcp_pubsub.metadata.exclude_prefixes", req.MetadataExcludePrefixes)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalStringField(fields, "gcp_pubsub.ordering_key", req.OrderingKey)
	setOptionalStringField(fields, "gcp_pubsub.publish_timeout", req.PublishTimeout)

	var files []multipartFilePart
	if req.CredentialsFilename != nil && req.CredentialsContentBase64 != nil && *req.CredentialsFilename != "" && *req.CredentialsContentBase64 != "" {
		content, err := decodeBase64File(*req.CredentialsContentBase64)
		if err != nil {
			return GcpPubsubOutlet{}, err
		}
		fields["gcp_credentials_file"] = payloadFileRef{File: "gcp_credentials_file_part"}
		files = append(files, multipartFilePart{
			FieldName: "gcp_credentials_file_part",
			FileName:  *req.CredentialsFilename,
			Content:   content,
		})
	}

	var response typedOutletResponse
	if err := c.doMultipartPayload(ctx, httpMethodPost, "/v1/terraform/gcp_pubsub", typedOutletCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, files, &response, statusCreated); err != nil {
		return GcpPubsubOutlet{}, err
	}

	return decodeGcpPubsubOutlet(response)
}

func (c *Client) GetGcpPubsubOutlet(ctx context.Context, id int64) (GcpPubsubOutlet, error) {
	var response typedOutletResponse
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/gcp_pubsub/%d", id), nil, &response, statusOK); err != nil {
		return GcpPubsubOutlet{}, err
	}

	return decodeGcpPubsubOutlet(response)
}

func (c *Client) UpdateGcpPubsubOutlet(ctx context.Context, id int64, req GcpPubsubOutletUpdateRequest) error {
	fields := map[string]any{
		"gcp_pubsub.project": req.Project,
		"gcp_pubsub.topic":   req.Topic,
	}
	setOptionalInt64Field(fields, "gcp_pubsub.batching.byte_size", req.BatchingByteSize)
	setOptionalStringField(fields, "gcp_pubsub.batching.check", req.BatchingCheck)
	setOptionalInt64Field(fields, "gcp_pubsub.batching.count", req.BatchingCount)
	setOptionalFloat64Field(fields, "gcp_pubsub.batching.jitter", req.BatchingJitter)
	setOptionalStringField(fields, "gcp_pubsub.batching.period", req.BatchingPeriod)
	setOptionalInt64Field(fields, "gcp_pubsub.byte_threshold", req.ByteThreshold)
	setOptionalInt64Field(fields, "gcp_pubsub.count_threshold", req.CountThreshold)
	setOptionalStringField(fields, "gcp_pubsub.delay_threshold", req.DelayThreshold)
	setOptionalStringField(fields, "gcp_pubsub.endpoint", req.Endpoint)
	setOptionalStringField(fields, "gcp_pubsub.flow_control.limit_exceeded_behavior", req.FlowControlLimitExceededBehavior)
	setOptionalInt64Field(fields, "gcp_pubsub.flow_control.max_outstanding_bytes", req.FlowControlMaxOutstandingBytes)
	setOptionalInt64Field(fields, "gcp_pubsub.flow_control.max_outstanding_messages", req.FlowControlMaxOutstandingMessages)
	setOptionalInt64Field(fields, "gcp_pubsub.max_in_flight", req.MaxInFlight)
	setOptionalStringListField(fields, "gcp_pubsub.metadata.exclude_prefixes", req.MetadataExcludePrefixes)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalStringField(fields, "gcp_pubsub.ordering_key", req.OrderingKey)
	setOptionalStringField(fields, "gcp_pubsub.publish_timeout", req.PublishTimeout)

	var files []multipartFilePart
	if req.CredentialsFilename != nil && req.CredentialsContentBase64 != nil && *req.CredentialsFilename != "" && *req.CredentialsContentBase64 != "" {
		content, err := decodeBase64File(*req.CredentialsContentBase64)
		if err != nil {
			return err
		}
		fields["gcp_credentials_file"] = payloadFileRef{File: "gcp_credentials_file_part"}
		files = append(files, multipartFilePart{
			FieldName: "gcp_credentials_file_part",
			FileName:  *req.CredentialsFilename,
			Content:   content,
		})
	}

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/gcp_pubsub/%d", id), typedOutletUpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, files, nil, statusOK)
}

func (c *Client) DeleteGcpPubsubOutlet(ctx context.Context, id int64) error {
	_, err := c.DeleteGcpPubsubOutletOperation(ctx, id)
	return err
}

func (c *Client) DeleteGcpPubsubOutletOperation(ctx context.Context, id int64) (OutletOperation, error) {
	return c.deleteTypedOutlet(ctx, fmt.Sprintf("/v1/terraform/gcp_pubsub/%d", id))
}

func (c *Client) SetGcpPubsubField(ctx context.Context, id int64, field string, value any) error {
	return c.patchOutletField(ctx, fmt.Sprintf("/v1/terraform/gcp_pubsub/%d/fields/%s", id, field), value)
}

func (c *Client) DeleteGcpPubsubField(ctx context.Context, id int64, field string) error {
	return c.deleteOutletField(ctx, fmt.Sprintf("/v1/terraform/gcp_pubsub/%d/fields/%s", id, field))
}

func decodeGcpPubsubOutlet(response typedOutletResponse) (GcpPubsubOutlet, error) {
	project, err := requiredStringField(response.Config, "gcp_pubsub.project")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	topic, err := requiredStringField(response.Config, "gcp_pubsub.topic")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}

	batchingByteSize, err := int64FieldPointer(response.Config, "gcp_pubsub.batching.byte_size")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	batchingCheck, err := stringFieldPointer(response.Config, "gcp_pubsub.batching.check")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	batchingCount, err := int64FieldPointer(response.Config, "gcp_pubsub.batching.count")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	batchingJitter, err := float64FieldPointer(response.Config, "gcp_pubsub.batching.jitter")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	batchingPeriod, err := stringFieldPointer(response.Config, "gcp_pubsub.batching.period")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	byteThreshold, err := int64FieldPointer(response.Config, "gcp_pubsub.byte_threshold")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	countThreshold, err := int64FieldPointer(response.Config, "gcp_pubsub.count_threshold")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	delayThreshold, err := stringFieldPointer(response.Config, "gcp_pubsub.delay_threshold")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	endpoint, err := stringFieldPointer(response.Config, "gcp_pubsub.endpoint")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	flowControlLimitExceededBehavior, err := stringFieldPointer(response.Config, "gcp_pubsub.flow_control.limit_exceeded_behavior")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	flowControlMaxOutstandingBytes, err := int64FieldPointer(response.Config, "gcp_pubsub.flow_control.max_outstanding_bytes")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	flowControlMaxOutstandingMessages, err := int64FieldPointer(response.Config, "gcp_pubsub.flow_control.max_outstanding_messages")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	maxInFlight, err := int64FieldPointer(response.Config, "gcp_pubsub.max_in_flight")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	metadataExcludePrefixes, err := stringListField(response.Config, "gcp_pubsub.metadata.exclude_prefixes")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	notifyPolicy, err := notifyPolicyPointerPreservingEmpty(response.Config)
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	orderingKey, err := stringFieldPointer(response.Config, "gcp_pubsub.ordering_key")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}
	publishTimeout, err := stringFieldPointer(response.Config, "gcp_pubsub.publish_timeout")
	if err != nil {
		return GcpPubsubOutlet{}, err
	}

	return GcpPubsubOutlet{
		ID:                                response.ID,
		Name:                              response.Name,
		Description:                       response.Description,
		Enabled:                           response.Enabled,
		BatchingByteSize:                  batchingByteSize,
		BatchingCheck:                     batchingCheck,
		BatchingCount:                     batchingCount,
		BatchingJitter:                    batchingJitter,
		BatchingPeriod:                    batchingPeriod,
		ByteThreshold:                     byteThreshold,
		CountThreshold:                    countThreshold,
		DelayThreshold:                    delayThreshold,
		Endpoint:                          endpoint,
		FlowControlLimitExceededBehavior:  flowControlLimitExceededBehavior,
		FlowControlMaxOutstandingBytes:    flowControlMaxOutstandingBytes,
		FlowControlMaxOutstandingMessages: flowControlMaxOutstandingMessages,
		MaxInFlight:                       maxInFlight,
		MetadataExcludePrefixes:           metadataExcludePrefixes,
		NotifyPolicy:                      notifyPolicy,
		OrderingKey:                       orderingKey,
		Project:                           project,
		PublishTimeout:                    publishTimeout,
		Topic:                             topic,
		ApplyStatus:                       response.ApplyStatus,
		LastApplyError:                    response.LastApplyError,
		ConfigGeneration:                  response.ConfigGeneration,
		AppliedGeneration:                 response.AppliedGeneration,
	}, nil
}

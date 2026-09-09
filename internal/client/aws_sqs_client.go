package client

import (
	"context"
	"fmt"
)

type AwsSqsOutlet struct {
	ID                      int64
	Name                    string
	Description             string
	Enabled                 bool
	BackoffInitialInterval  *string
	BackoffMaxElapsedTime   *string
	BackoffMaxInterval      *string
	BatchingByteSize        *int64
	BatchingCheck           *string
	BatchingCount           *int64
	BatchingJitter          *float64
	BatchingPeriod          *string
	DelaySeconds            *int64
	EndpointURL             *string
	ExpiryWindow            *string
	ExternalID              *string
	MaxInFlight             *int64
	MaxRetries              *int64
	MessageDeduplicationID  *string
	MessageGroupID          *string
	MetadataExcludePrefixes []string
	NotifyPolicy            *string
	Profile                 *string
	QueueURL                string
	Region                  *string
	RoleARN                 *string
	UseEC2InstanceRole      *bool
	ApplyStatus             string
	LastApplyError          string
	ConfigGeneration        int64
	AppliedGeneration       int64
}

type AwsSqsOutletCreateRequest struct {
	Name                    string
	Description             string
	Enabled                 bool
	AccessKeyID             *string
	BackoffInitialInterval  *string
	BackoffMaxElapsedTime   *string
	BackoffMaxInterval      *string
	BatchingByteSize        *int64
	BatchingCheck           *string
	BatchingCount           *int64
	BatchingJitter          *float64
	BatchingPeriod          *string
	DelaySeconds            *int64
	EndpointURL             *string
	ExpiryWindow            *string
	ExternalID              *string
	MaxInFlight             *int64
	MaxRetries              *int64
	MessageDeduplicationID  *string
	MessageGroupID          *string
	MetadataExcludePrefixes []string
	NotifyPolicy            *string
	Profile                 *string
	QueueURL                string
	Region                  *string
	RoleARN                 *string
	SecretAccessKey         *string
	SessionToken            *string
	UseEC2InstanceRole      *bool
}

type AwsSqsOutletUpdateRequest = AwsSqsOutletCreateRequest

func (c *Client) CreateAwsSqsOutlet(ctx context.Context, req AwsSqsOutletCreateRequest) (AwsSqsOutlet, error) {
	fields := map[string]any{
		"aws_sqs.queue_url": req.QueueURL,
	}
	setOptionalStringField(fields, "access_key_id", req.AccessKeyID)
	setOptionalStringField(fields, "aws_sqs.backoff.initial_interval", req.BackoffInitialInterval)
	setOptionalStringField(fields, "aws_sqs.backoff.max_elapsed_time", req.BackoffMaxElapsedTime)
	setOptionalStringField(fields, "aws_sqs.backoff.max_interval", req.BackoffMaxInterval)
	setOptionalInt64Field(fields, "aws_sqs.batching.byte_size", req.BatchingByteSize)
	setOptionalStringField(fields, "aws_sqs.batching.check", req.BatchingCheck)
	setOptionalInt64Field(fields, "aws_sqs.batching.count", req.BatchingCount)
	setOptionalFloat64Field(fields, "aws_sqs.batching.jitter", req.BatchingJitter)
	setOptionalStringField(fields, "aws_sqs.batching.period", req.BatchingPeriod)
	setOptionalInt64Field(fields, "aws_sqs.delay_seconds", req.DelaySeconds)
	setOptionalStringField(fields, "aws_sqs.endpoint_url", req.EndpointURL)
	setOptionalStringField(fields, "aws_sqs.expiry_window", req.ExpiryWindow)
	setOptionalStringField(fields, "aws_sqs.external_id", req.ExternalID)
	setOptionalInt64Field(fields, "aws_sqs.max_in_flight", req.MaxInFlight)
	setOptionalInt64Field(fields, "aws_sqs.max_retries", req.MaxRetries)
	setOptionalStringField(fields, "aws_sqs.message_deduplication_id", req.MessageDeduplicationID)
	setOptionalStringField(fields, "aws_sqs.message_group_id", req.MessageGroupID)
	setOptionalStringListField(fields, "aws_sqs.metadata.exclude_prefixes", req.MetadataExcludePrefixes)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalStringField(fields, "aws_sqs.profile", req.Profile)
	setOptionalStringField(fields, "aws_sqs.region", req.Region)
	setOptionalStringField(fields, "aws_sqs.role_arn", req.RoleARN)
	setOptionalStringField(fields, "secret_access_key", req.SecretAccessKey)
	setOptionalStringField(fields, "session_token", req.SessionToken)
	setOptionalBoolField(fields, "aws_sqs.use_ec2_instance_role", req.UseEC2InstanceRole)

	var response typedOutletResponse
	if err := c.doMultipartPayload(ctx, httpMethodPost, "/v1/terraform/aws_sqs", typedOutletCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, &response, statusCreated); err != nil {
		return AwsSqsOutlet{}, err
	}

	return decodeAwsSqsOutlet(response)
}

func (c *Client) GetAwsSqsOutlet(ctx context.Context, id int64) (AwsSqsOutlet, error) {
	var response typedOutletResponse
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/aws_sqs/%d", id), nil, &response, statusOK); err != nil {
		return AwsSqsOutlet{}, err
	}

	return decodeAwsSqsOutlet(response)
}

func (c *Client) UpdateAwsSqsOutlet(ctx context.Context, id int64, req AwsSqsOutletUpdateRequest) error {
	fields := map[string]any{
		"aws_sqs.queue_url": req.QueueURL,
	}
	setOptionalStringField(fields, "access_key_id", req.AccessKeyID)
	setOptionalStringField(fields, "aws_sqs.backoff.initial_interval", req.BackoffInitialInterval)
	setOptionalStringField(fields, "aws_sqs.backoff.max_elapsed_time", req.BackoffMaxElapsedTime)
	setOptionalStringField(fields, "aws_sqs.backoff.max_interval", req.BackoffMaxInterval)
	setOptionalInt64Field(fields, "aws_sqs.batching.byte_size", req.BatchingByteSize)
	setOptionalStringField(fields, "aws_sqs.batching.check", req.BatchingCheck)
	setOptionalInt64Field(fields, "aws_sqs.batching.count", req.BatchingCount)
	setOptionalFloat64Field(fields, "aws_sqs.batching.jitter", req.BatchingJitter)
	setOptionalStringField(fields, "aws_sqs.batching.period", req.BatchingPeriod)
	setOptionalInt64Field(fields, "aws_sqs.delay_seconds", req.DelaySeconds)
	setOptionalStringField(fields, "aws_sqs.endpoint_url", req.EndpointURL)
	setOptionalStringField(fields, "aws_sqs.expiry_window", req.ExpiryWindow)
	setOptionalStringField(fields, "aws_sqs.external_id", req.ExternalID)
	setOptionalInt64Field(fields, "aws_sqs.max_in_flight", req.MaxInFlight)
	setOptionalInt64Field(fields, "aws_sqs.max_retries", req.MaxRetries)
	setOptionalStringField(fields, "aws_sqs.message_deduplication_id", req.MessageDeduplicationID)
	setOptionalStringField(fields, "aws_sqs.message_group_id", req.MessageGroupID)
	setOptionalStringListField(fields, "aws_sqs.metadata.exclude_prefixes", req.MetadataExcludePrefixes)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalStringField(fields, "aws_sqs.profile", req.Profile)
	setOptionalStringField(fields, "aws_sqs.region", req.Region)
	setOptionalStringField(fields, "aws_sqs.role_arn", req.RoleARN)
	setOptionalStringField(fields, "secret_access_key", req.SecretAccessKey)
	setOptionalStringField(fields, "session_token", req.SessionToken)
	setOptionalBoolField(fields, "aws_sqs.use_ec2_instance_role", req.UseEC2InstanceRole)

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/aws_sqs/%d", id), typedOutletUpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, nil, statusOK)
}

func (c *Client) DeleteAwsSqsOutlet(ctx context.Context, id int64) error {
	_, err := c.DeleteAwsSqsOutletOperation(ctx, id)
	return err
}

func (c *Client) DeleteAwsSqsOutletOperation(ctx context.Context, id int64) (OutletOperation, error) {
	return c.deleteTypedOutlet(ctx, fmt.Sprintf("/v1/terraform/aws_sqs/%d", id))
}

func (c *Client) SetAwsSqsField(ctx context.Context, id int64, field string, value any) error {
	return c.patchOutletField(ctx, fmt.Sprintf("/v1/terraform/aws_sqs/%d/fields/%s", id, field), value)
}

func (c *Client) DeleteAwsSqsField(ctx context.Context, id int64, field string) error {
	return c.deleteOutletField(ctx, fmt.Sprintf("/v1/terraform/aws_sqs/%d/fields/%s", id, field))
}

func decodeAwsSqsOutlet(response typedOutletResponse) (AwsSqsOutlet, error) {
	queueURL, err := requiredStringField(response.Config, "aws_sqs.queue_url")
	if err != nil {
		return AwsSqsOutlet{}, err
	}

	backoffInitialInterval, err := stringFieldPointer(response.Config, "aws_sqs.backoff.initial_interval")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	backoffMaxElapsedTime, err := stringFieldPointer(response.Config, "aws_sqs.backoff.max_elapsed_time")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	backoffMaxInterval, err := stringFieldPointer(response.Config, "aws_sqs.backoff.max_interval")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	batchingByteSize, err := int64FieldPointer(response.Config, "aws_sqs.batching.byte_size")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	batchingCheck, err := stringFieldPointer(response.Config, "aws_sqs.batching.check")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	batchingCount, err := int64FieldPointer(response.Config, "aws_sqs.batching.count")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	batchingJitter, err := float64FieldPointer(response.Config, "aws_sqs.batching.jitter")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	batchingPeriod, err := stringFieldPointer(response.Config, "aws_sqs.batching.period")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	delaySeconds, err := int64FieldPointer(response.Config, "aws_sqs.delay_seconds")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	endpointURL, err := stringFieldPointer(response.Config, "aws_sqs.endpoint_url")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	expiryWindow, err := stringFieldPointer(response.Config, "aws_sqs.expiry_window")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	externalID, err := stringFieldPointer(response.Config, "aws_sqs.external_id")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	maxInFlight, err := int64FieldPointer(response.Config, "aws_sqs.max_in_flight")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	maxRetries, err := int64FieldPointer(response.Config, "aws_sqs.max_retries")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	messageDeduplicationID, err := stringFieldPointer(response.Config, "aws_sqs.message_deduplication_id")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	messageGroupID, err := stringFieldPointer(response.Config, "aws_sqs.message_group_id")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	metadataExcludePrefixes, err := stringListField(response.Config, "aws_sqs.metadata.exclude_prefixes")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	notifyPolicy, err := notifyPolicyPointerPreservingEmpty(response.Config)
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	profile, err := stringFieldPointer(response.Config, "aws_sqs.profile")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	region, err := stringFieldPointer(response.Config, "aws_sqs.region")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	roleARN, err := stringFieldPointer(response.Config, "aws_sqs.role_arn")
	if err != nil {
		return AwsSqsOutlet{}, err
	}
	useEC2InstanceRole, err := boolFieldPointer(response.Config, "aws_sqs.use_ec2_instance_role")
	if err != nil {
		return AwsSqsOutlet{}, err
	}

	return AwsSqsOutlet{
		ID:                      response.ID,
		Name:                    response.Name,
		Description:             response.Description,
		Enabled:                 response.Enabled,
		BackoffInitialInterval:  backoffInitialInterval,
		BackoffMaxElapsedTime:   backoffMaxElapsedTime,
		BackoffMaxInterval:      backoffMaxInterval,
		BatchingByteSize:        batchingByteSize,
		BatchingCheck:           batchingCheck,
		BatchingCount:           batchingCount,
		BatchingJitter:          batchingJitter,
		BatchingPeriod:          batchingPeriod,
		DelaySeconds:            delaySeconds,
		EndpointURL:             endpointURL,
		ExpiryWindow:            expiryWindow,
		ExternalID:              externalID,
		MaxInFlight:             maxInFlight,
		MaxRetries:              maxRetries,
		MessageDeduplicationID:  messageDeduplicationID,
		MessageGroupID:          messageGroupID,
		MetadataExcludePrefixes: metadataExcludePrefixes,
		NotifyPolicy:            notifyPolicy,
		Profile:                 profile,
		QueueURL:                queueURL,
		Region:                  region,
		RoleARN:                 roleARN,
		UseEC2InstanceRole:      useEC2InstanceRole,
		ApplyStatus:             response.ApplyStatus,
		LastApplyError:          response.LastApplyError,
		ConfigGeneration:        response.ConfigGeneration,
		AppliedGeneration:       response.AppliedGeneration,
	}, nil
}

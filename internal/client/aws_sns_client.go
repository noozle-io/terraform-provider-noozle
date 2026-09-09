package client

import (
	"context"
	"fmt"
)

type AwsSnsOutlet struct {
	ID                      int64
	Name                    string
	Description             string
	Enabled                 bool
	EndpointURL             *string
	ExpiryWindow            *string
	ExternalID              *string
	MaxInFlight             *int64
	MessageDeduplicationID  *string
	MessageGroupID          *string
	MetadataExcludePrefixes []string
	NotifyPolicy            *string
	Profile                 *string
	Region                  *string
	RoleARN                 *string
	Timeout                 *string
	TopicARN                string
	UseEC2InstanceRole      *bool
	ApplyStatus             string
	LastApplyError          string
	ConfigGeneration        int64
	AppliedGeneration       int64
}

type AwsSnsOutletCreateRequest struct {
	Name                    string
	Description             string
	Enabled                 bool
	AccessKeyID             *string
	EndpointURL             *string
	ExpiryWindow            *string
	ExternalID              *string
	MaxInFlight             *int64
	MessageDeduplicationID  *string
	MessageGroupID          *string
	MetadataExcludePrefixes []string
	NotifyPolicy            *string
	Profile                 *string
	Region                  *string
	RoleARN                 *string
	SecretAccessKey         *string
	SessionToken            *string
	Timeout                 *string
	TopicARN                string
	UseEC2InstanceRole      *bool
}

type AwsSnsOutletUpdateRequest = AwsSnsOutletCreateRequest

func (c *Client) CreateAwsSnsOutlet(ctx context.Context, req AwsSnsOutletCreateRequest) (AwsSnsOutlet, error) {
	fields := map[string]any{
		"aws_sns.topic_arn": req.TopicARN,
	}
	setOptionalStringField(fields, "access_key_id", req.AccessKeyID)
	setOptionalStringField(fields, "aws_sns.endpoint_url", req.EndpointURL)
	setOptionalStringField(fields, "aws_sns.expiry_window", req.ExpiryWindow)
	setOptionalStringField(fields, "aws_sns.external_id", req.ExternalID)
	setOptionalInt64Field(fields, "aws_sns.max_in_flight", req.MaxInFlight)
	setOptionalStringField(fields, "aws_sns.message_deduplication_id", req.MessageDeduplicationID)
	setOptionalStringField(fields, "aws_sns.message_group_id", req.MessageGroupID)
	setOptionalStringListField(fields, "aws_sns.metadata.exclude_prefixes", req.MetadataExcludePrefixes)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalStringField(fields, "aws_sns.profile", req.Profile)
	setOptionalStringField(fields, "aws_sns.region", req.Region)
	setOptionalStringField(fields, "aws_sns.role_arn", req.RoleARN)
	setOptionalStringField(fields, "secret_access_key", req.SecretAccessKey)
	setOptionalStringField(fields, "session_token", req.SessionToken)
	setOptionalStringField(fields, "aws_sns.timeout", req.Timeout)
	setOptionalBoolField(fields, "aws_sns.use_ec2_instance_role", req.UseEC2InstanceRole)

	var response typedOutletResponse
	if err := c.doMultipartPayload(ctx, httpMethodPost, "/v1/terraform/aws_sns", typedOutletCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, &response, statusCreated); err != nil {
		return AwsSnsOutlet{}, err
	}

	return decodeAwsSnsOutlet(response)
}

func (c *Client) GetAwsSnsOutlet(ctx context.Context, id int64) (AwsSnsOutlet, error) {
	var response typedOutletResponse
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/aws_sns/%d", id), nil, &response, statusOK); err != nil {
		return AwsSnsOutlet{}, err
	}

	return decodeAwsSnsOutlet(response)
}

func (c *Client) UpdateAwsSnsOutlet(ctx context.Context, id int64, req AwsSnsOutletUpdateRequest) error {
	fields := map[string]any{
		"aws_sns.topic_arn": req.TopicARN,
	}
	setOptionalStringField(fields, "access_key_id", req.AccessKeyID)
	setOptionalStringField(fields, "aws_sns.endpoint_url", req.EndpointURL)
	setOptionalStringField(fields, "aws_sns.expiry_window", req.ExpiryWindow)
	setOptionalStringField(fields, "aws_sns.external_id", req.ExternalID)
	setOptionalInt64Field(fields, "aws_sns.max_in_flight", req.MaxInFlight)
	setOptionalStringField(fields, "aws_sns.message_deduplication_id", req.MessageDeduplicationID)
	setOptionalStringField(fields, "aws_sns.message_group_id", req.MessageGroupID)
	setOptionalStringListField(fields, "aws_sns.metadata.exclude_prefixes", req.MetadataExcludePrefixes)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalStringField(fields, "aws_sns.profile", req.Profile)
	setOptionalStringField(fields, "aws_sns.region", req.Region)
	setOptionalStringField(fields, "aws_sns.role_arn", req.RoleARN)
	setOptionalStringField(fields, "secret_access_key", req.SecretAccessKey)
	setOptionalStringField(fields, "session_token", req.SessionToken)
	setOptionalStringField(fields, "aws_sns.timeout", req.Timeout)
	setOptionalBoolField(fields, "aws_sns.use_ec2_instance_role", req.UseEC2InstanceRole)

	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/aws_sns/%d", id), typedOutletUpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, nil, statusOK)
}

func (c *Client) DeleteAwsSnsOutlet(ctx context.Context, id int64) error {
	_, err := c.DeleteAwsSnsOutletOperation(ctx, id)
	return err
}

func (c *Client) DeleteAwsSnsOutletOperation(ctx context.Context, id int64) (OutletOperation, error) {
	return c.deleteTypedOutlet(ctx, fmt.Sprintf("/v1/terraform/aws_sns/%d", id))
}

func (c *Client) SetAwsSnsField(ctx context.Context, id int64, field string, value any) error {
	return c.patchOutletField(ctx, fmt.Sprintf("/v1/terraform/aws_sns/%d/fields/%s", id, field), value)
}

func (c *Client) DeleteAwsSnsField(ctx context.Context, id int64, field string) error {
	return c.deleteOutletField(ctx, fmt.Sprintf("/v1/terraform/aws_sns/%d/fields/%s", id, field))
}

func decodeAwsSnsOutlet(response typedOutletResponse) (AwsSnsOutlet, error) {
	topicARN, err := requiredStringField(response.Config, "aws_sns.topic_arn")
	if err != nil {
		return AwsSnsOutlet{}, err
	}

	endpointURL, err := stringFieldPointer(response.Config, "aws_sns.endpoint_url")
	if err != nil {
		return AwsSnsOutlet{}, err
	}
	expiryWindow, err := stringFieldPointer(response.Config, "aws_sns.expiry_window")
	if err != nil {
		return AwsSnsOutlet{}, err
	}
	externalID, err := stringFieldPointer(response.Config, "aws_sns.external_id")
	if err != nil {
		return AwsSnsOutlet{}, err
	}
	maxInFlight, err := int64FieldPointer(response.Config, "aws_sns.max_in_flight")
	if err != nil {
		return AwsSnsOutlet{}, err
	}
	messageDeduplicationID, err := stringFieldPointer(response.Config, "aws_sns.message_deduplication_id")
	if err != nil {
		return AwsSnsOutlet{}, err
	}
	messageGroupID, err := stringFieldPointer(response.Config, "aws_sns.message_group_id")
	if err != nil {
		return AwsSnsOutlet{}, err
	}
	metadataExcludePrefixes, err := stringListField(response.Config, "aws_sns.metadata.exclude_prefixes")
	if err != nil {
		return AwsSnsOutlet{}, err
	}
	notifyPolicy, err := notifyPolicyPointerPreservingEmpty(response.Config)
	if err != nil {
		return AwsSnsOutlet{}, err
	}
	profile, err := stringFieldPointer(response.Config, "aws_sns.profile")
	if err != nil {
		return AwsSnsOutlet{}, err
	}
	region, err := stringFieldPointer(response.Config, "aws_sns.region")
	if err != nil {
		return AwsSnsOutlet{}, err
	}
	roleARN, err := stringFieldPointer(response.Config, "aws_sns.role_arn")
	if err != nil {
		return AwsSnsOutlet{}, err
	}
	timeout, err := stringFieldPointer(response.Config, "aws_sns.timeout")
	if err != nil {
		return AwsSnsOutlet{}, err
	}
	useEC2InstanceRole, err := boolFieldPointer(response.Config, "aws_sns.use_ec2_instance_role")
	if err != nil {
		return AwsSnsOutlet{}, err
	}

	return AwsSnsOutlet{
		ID:                      response.ID,
		Name:                    response.Name,
		Description:             response.Description,
		Enabled:                 response.Enabled,
		EndpointURL:             endpointURL,
		ExpiryWindow:            expiryWindow,
		ExternalID:              externalID,
		MaxInFlight:             maxInFlight,
		MessageDeduplicationID:  messageDeduplicationID,
		MessageGroupID:          messageGroupID,
		MetadataExcludePrefixes: metadataExcludePrefixes,
		NotifyPolicy:            notifyPolicy,
		Profile:                 profile,
		Region:                  region,
		RoleARN:                 roleARN,
		Timeout:                 timeout,
		TopicARN:                topicARN,
		UseEC2InstanceRole:      useEC2InstanceRole,
		ApplyStatus:             response.ApplyStatus,
		LastApplyError:          response.LastApplyError,
		ConfigGeneration:        response.ConfigGeneration,
		AppliedGeneration:       response.AppliedGeneration,
	}, nil
}

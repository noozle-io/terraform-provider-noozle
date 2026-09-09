package client

import (
	"context"
	"fmt"
)

const (
	kafkaTLSFeatureID                = "kafka.tls"
	kafkaSASLPlainFeatureID          = "kafka.sasl.plain"
	kafkaSASLScram256FeatureID       = "kafka.sasl.scram256"
	kafkaSASLScram512FeatureID       = "kafka.sasl.scram512"
	kafkaSASLOAuthbearerFeatureID    = "kafka.sasl.oauthbearer"
	kafkaSASLOAuthbearerCacheFeature = "kafka.sasl.oauthbearer.cache"
	kafkaSASLAwsMSKIAMFeatureID      = "kafka.sasl.aws_msk_iam"

	kafkaCACertRole            = "ca_cert"
	kafkaClientCertRole        = "client_cert"
	kafkaClientKeyRole         = "client_key"
	kafkaSASLPasswordRole      = "sasl_password"
	kafkaSASLAccessTokenRole   = "kafka_sasl_access_token"
	kafkaAwsSecretRole         = "aws_secret"
	kafkaAwsTokenRole          = "aws_token"
	kafkaTLSCACertPartName     = "kafka_ca_cert_part"
	kafkaTLSClientCertPartName = "kafka_client_cert_part"
	kafkaTLSClientKeyPartName  = "kafka_client_key_part"
)

type KafkaOutlet struct {
	ID                                   int64
	Name                                 string
	Description                          string
	Enabled                              bool
	AckReplicas                          *bool
	Addresses                            []string
	BackoffInitialInterval               *string
	BackoffMaxElapsedTime                *string
	BackoffMaxInterval                   *string
	BatchingByteSize                     *int64
	BatchingCheck                        *string
	BatchingCount                        *int64
	BatchingJitter                       *float64
	BatchingPeriod                       *string
	ClientID                             *string
	Compression                          *string
	CustomTopicCreationEnabled           *bool
	CustomTopicCreationPartitions        *int64
	CustomTopicCreationReplicationFactor *int64
	IdempotentWrite                      *bool
	InjectTracingMap                     *string
	Key                                  *string
	MaxInFlight                          *int64
	MaxMsgBytes                          *int64
	MaxRetries                           *int64
	MetadataExcludePrefixes              []string
	NotifyPolicy                         *string
	Partition                            *string
	Partitioner                          *string
	RackID                               *string
	RetryAsBatch                         *bool
	StaticHeaders                        map[string]string
	TargetVersion                        *string
	Timeout                              *string
	Topic                                string
	ApplyStatus                          string
	LastApplyError                       string
	ConfigGeneration                     int64
	AppliedGeneration                    int64
}

type KafkaOutletCreateRequest struct {
	Name                                 string
	Description                          string
	Enabled                              bool
	AckReplicas                          *bool
	Addresses                            []string
	BackoffInitialInterval               *string
	BackoffMaxElapsedTime                *string
	BackoffMaxInterval                   *string
	BatchingByteSize                     *int64
	BatchingCheck                        *string
	BatchingCount                        *int64
	BatchingJitter                       *float64
	BatchingPeriod                       *string
	ClientID                             *string
	Compression                          *string
	CustomTopicCreationEnabled           *bool
	CustomTopicCreationPartitions        *int64
	CustomTopicCreationReplicationFactor *int64
	IdempotentWrite                      *bool
	InjectTracingMap                     *string
	Key                                  *string
	MaxInFlight                          *int64
	MaxMsgBytes                          *int64
	MaxRetries                           *int64
	MetadataExcludePrefixes              []string
	NotifyPolicy                         *string
	Partition                            *string
	Partitioner                          *string
	RackID                               *string
	RetryAsBatch                         *bool
	StaticHeaders                        map[string]string
	TargetVersion                        *string
	Timeout                              *string
	Topic                                string
}

type KafkaOutletUpdateRequest = KafkaOutletCreateRequest

type KafkaTLSFeatureRequest struct {
	CACertFilename          string
	CACertContentBase64     string
	ClientCertFilename      string
	ClientCertContentBase64 string
	ClientKeyFilename       string
	ClientKeyContentBase64  string
}

type KafkaSASLUserPassFeatureRequest struct {
	User     string
	Password string
}

type KafkaSASLOAuthbearerFeatureRequest struct {
	AccessToken string
}

type KafkaSASLOAuthbearerCacheFeatureRequest struct {
	TokenCache string
	TokenKey   string
}

type KafkaSASLAwsMSKIAMFeatureRequest struct {
	Region         string
	Endpoint       *string
	Profile        *string
	ID             *string
	FromEC2Role    *string
	Role           *string
	RoleExternalID *string
	ExpiryWindow   *string
	AwsSecret      *string
	AwsToken       *string
}

func (c *Client) CreateKafkaOutlet(ctx context.Context, req KafkaOutletCreateRequest) (KafkaOutlet, error) {
	fields := kafkaOutletFields(req)
	var response typedOutletResponse
	if err := c.doMultipartPayload(ctx, httpMethodPost, "/v1/terraform/kafka", typedOutletCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      fields,
	}, nil, &response, statusCreated); err != nil {
		return KafkaOutlet{}, err
	}
	return decodeKafkaOutlet(response)
}

func (c *Client) GetKafkaOutlet(ctx context.Context, id int64) (KafkaOutlet, error) {
	var response typedOutletResponse
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/kafka/%d", id), nil, &response, statusOK); err != nil {
		return KafkaOutlet{}, err
	}
	return decodeKafkaOutlet(response)
}

func (c *Client) UpdateKafkaOutlet(ctx context.Context, id int64, req KafkaOutletUpdateRequest) error {
	return c.doMultipartPayload(ctx, httpMethodPut, fmt.Sprintf("/v1/terraform/kafka/%d", id), typedOutletUpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Fields:      kafkaOutletFields(req),
	}, nil, nil, statusOK)
}

func (c *Client) DeleteKafkaOutlet(ctx context.Context, id int64) error {
	_, err := c.DeleteKafkaOutletOperation(ctx, id)
	return err
}

func (c *Client) DeleteKafkaOutletOperation(ctx context.Context, id int64) (OutletOperation, error) {
	return c.deleteTypedOutlet(ctx, fmt.Sprintf("/v1/terraform/kafka/%d", id))
}

func (c *Client) SetKafkaField(ctx context.Context, id int64, field string, value any) error {
	return c.patchOutletField(ctx, fmt.Sprintf("/v1/terraform/kafka/%d/fields/%s", id, field), value)
}

func (c *Client) DeleteKafkaField(ctx context.Context, id int64, field string) error {
	return c.deleteOutletField(ctx, fmt.Sprintf("/v1/terraform/kafka/%d/fields/%s", id, field))
}

func (c *Client) ListKafkaCredentials(ctx context.Context, outletID int64) ([]OutletMaterialSummary, error) {
	return c.listOutletMaterials(ctx, fmt.Sprintf("/v1/terraform/kafka/%d/credentials", outletID))
}

func (c *Client) DeleteKafkaCredential(ctx context.Context, outletID int64, role string) (OutletMaterialDeleteResponse, error) {
	return c.deleteOutletMaterial(ctx, fmt.Sprintf("/v1/terraform/kafka/%d/credentials/%s", outletID, role))
}

func (c *Client) GetKafkaTLSFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	return c.getKafkaFeature(ctx, outletID, kafkaTLSFeatureID)
}

func (c *Client) ConfigureKafkaTLSFeature(ctx context.Context, outletID int64, req KafkaTLSFeatureRequest) error {
	return c.putKafkaTLSFeature(ctx, httpMethodPost, outletID, req)
}

func (c *Client) UpdateKafkaTLSFeature(ctx context.Context, outletID int64, req KafkaTLSFeatureRequest) error {
	return c.putKafkaTLSFeature(ctx, httpMethodPut, outletID, req)
}

func (c *Client) DeleteKafkaTLSFeature(ctx context.Context, outletID int64) error {
	return c.deleteKafkaFeature(ctx, outletID, kafkaTLSFeatureID)
}

func (c *Client) GetKafkaSASLPlainFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	return c.getKafkaFeature(ctx, outletID, kafkaSASLPlainFeatureID)
}

func (c *Client) ConfigureKafkaSASLPlainFeature(ctx context.Context, outletID int64, req KafkaSASLUserPassFeatureRequest) error {
	return c.putKafkaSASLUserPassFeature(ctx, httpMethodPost, outletID, kafkaSASLPlainFeatureID, req)
}

func (c *Client) UpdateKafkaSASLPlainFeature(ctx context.Context, outletID int64, req KafkaSASLUserPassFeatureRequest) error {
	return c.putKafkaSASLUserPassFeature(ctx, httpMethodPut, outletID, kafkaSASLPlainFeatureID, req)
}

func (c *Client) DeleteKafkaSASLPlainFeature(ctx context.Context, outletID int64) error {
	return c.deleteKafkaFeature(ctx, outletID, kafkaSASLPlainFeatureID)
}

func (c *Client) GetKafkaSASLScram256Feature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	return c.getKafkaFeature(ctx, outletID, kafkaSASLScram256FeatureID)
}

func (c *Client) ConfigureKafkaSASLScram256Feature(ctx context.Context, outletID int64, req KafkaSASLUserPassFeatureRequest) error {
	return c.putKafkaSASLUserPassFeature(ctx, httpMethodPost, outletID, kafkaSASLScram256FeatureID, req)
}

func (c *Client) UpdateKafkaSASLScram256Feature(ctx context.Context, outletID int64, req KafkaSASLUserPassFeatureRequest) error {
	return c.putKafkaSASLUserPassFeature(ctx, httpMethodPut, outletID, kafkaSASLScram256FeatureID, req)
}

func (c *Client) DeleteKafkaSASLScram256Feature(ctx context.Context, outletID int64) error {
	return c.deleteKafkaFeature(ctx, outletID, kafkaSASLScram256FeatureID)
}

func (c *Client) GetKafkaSASLScram512Feature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	return c.getKafkaFeature(ctx, outletID, kafkaSASLScram512FeatureID)
}

func (c *Client) ConfigureKafkaSASLScram512Feature(ctx context.Context, outletID int64, req KafkaSASLUserPassFeatureRequest) error {
	return c.putKafkaSASLUserPassFeature(ctx, httpMethodPost, outletID, kafkaSASLScram512FeatureID, req)
}

func (c *Client) UpdateKafkaSASLScram512Feature(ctx context.Context, outletID int64, req KafkaSASLUserPassFeatureRequest) error {
	return c.putKafkaSASLUserPassFeature(ctx, httpMethodPut, outletID, kafkaSASLScram512FeatureID, req)
}

func (c *Client) DeleteKafkaSASLScram512Feature(ctx context.Context, outletID int64) error {
	return c.deleteKafkaFeature(ctx, outletID, kafkaSASLScram512FeatureID)
}

func (c *Client) GetKafkaSASLOAuthbearerFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	return c.getKafkaFeature(ctx, outletID, kafkaSASLOAuthbearerFeatureID)
}

func (c *Client) ConfigureKafkaSASLOAuthbearerFeature(ctx context.Context, outletID int64, req KafkaSASLOAuthbearerFeatureRequest) error {
	return c.putKafkaOAuthFeature(ctx, httpMethodPost, outletID, req)
}

func (c *Client) UpdateKafkaSASLOAuthbearerFeature(ctx context.Context, outletID int64, req KafkaSASLOAuthbearerFeatureRequest) error {
	return c.putKafkaOAuthFeature(ctx, httpMethodPut, outletID, req)
}

func (c *Client) DeleteKafkaSASLOAuthbearerFeature(ctx context.Context, outletID int64) error {
	return c.deleteKafkaFeature(ctx, outletID, kafkaSASLOAuthbearerFeatureID)
}

func (c *Client) GetKafkaSASLOAuthbearerCacheFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	return c.getKafkaFeature(ctx, outletID, kafkaSASLOAuthbearerCacheFeature)
}

func (c *Client) ConfigureKafkaSASLOAuthbearerCacheFeature(ctx context.Context, outletID int64, req KafkaSASLOAuthbearerCacheFeatureRequest) error {
	return c.putKafkaOAuthCacheFeature(ctx, httpMethodPost, outletID, req)
}

func (c *Client) UpdateKafkaSASLOAuthbearerCacheFeature(ctx context.Context, outletID int64, req KafkaSASLOAuthbearerCacheFeatureRequest) error {
	return c.putKafkaOAuthCacheFeature(ctx, httpMethodPut, outletID, req)
}

func (c *Client) DeleteKafkaSASLOAuthbearerCacheFeature(ctx context.Context, outletID int64) error {
	return c.deleteKafkaFeature(ctx, outletID, kafkaSASLOAuthbearerCacheFeature)
}

func (c *Client) GetKafkaSASLAwsMSKIAMFeature(ctx context.Context, outletID int64) (OutletFeatureDetail, error) {
	return c.getKafkaFeature(ctx, outletID, kafkaSASLAwsMSKIAMFeatureID)
}

func (c *Client) ConfigureKafkaSASLAwsMSKIAMFeature(ctx context.Context, outletID int64, req KafkaSASLAwsMSKIAMFeatureRequest) error {
	return c.putKafkaAwsMSKIAMFeature(ctx, httpMethodPost, outletID, req)
}

func (c *Client) UpdateKafkaSASLAwsMSKIAMFeature(ctx context.Context, outletID int64, req KafkaSASLAwsMSKIAMFeatureRequest) error {
	return c.putKafkaAwsMSKIAMFeature(ctx, httpMethodPut, outletID, req)
}

func (c *Client) DeleteKafkaSASLAwsMSKIAMFeature(ctx context.Context, outletID int64) error {
	return c.deleteKafkaFeature(ctx, outletID, kafkaSASLAwsMSKIAMFeatureID)
}

func kafkaOutletFields(req KafkaOutletCreateRequest) map[string]any {
	fields := map[string]any{
		"kafka.addresses": req.Addresses,
		"kafka.topic":     req.Topic,
	}
	setOptionalBoolField(fields, "kafka.ack_replicas", req.AckReplicas)
	setOptionalStringField(fields, "kafka.backoff.initial_interval", req.BackoffInitialInterval)
	setOptionalStringField(fields, "kafka.backoff.max_elapsed_time", req.BackoffMaxElapsedTime)
	setOptionalStringField(fields, "kafka.backoff.max_interval", req.BackoffMaxInterval)
	setOptionalInt64Field(fields, "kafka.batching.byte_size", req.BatchingByteSize)
	setOptionalStringField(fields, "kafka.batching.check", req.BatchingCheck)
	setOptionalInt64Field(fields, "kafka.batching.count", req.BatchingCount)
	setOptionalFloat64Field(fields, "kafka.batching.jitter", req.BatchingJitter)
	setOptionalStringField(fields, "kafka.batching.period", req.BatchingPeriod)
	setOptionalStringField(fields, "kafka.client_id", req.ClientID)
	setOptionalStringField(fields, "kafka.compression", req.Compression)
	setOptionalBoolField(fields, "kafka.custom_topic_creation.enabled", req.CustomTopicCreationEnabled)
	setOptionalInt64Field(fields, "kafka.custom_topic_creation.partitions", req.CustomTopicCreationPartitions)
	setOptionalInt64Field(fields, "kafka.custom_topic_creation.replication_factor", req.CustomTopicCreationReplicationFactor)
	setOptionalBoolField(fields, "kafka.idempotent_write", req.IdempotentWrite)
	setOptionalStringField(fields, "kafka.inject_tracing_map", req.InjectTracingMap)
	setOptionalStringField(fields, "kafka.key", req.Key)
	setOptionalInt64Field(fields, "kafka.max_in_flight", req.MaxInFlight)
	setOptionalInt64Field(fields, "kafka.max_msg_bytes", req.MaxMsgBytes)
	setOptionalInt64Field(fields, "kafka.max_retries", req.MaxRetries)
	setOptionalStringListField(fields, "kafka.metadata.exclude_prefixes", req.MetadataExcludePrefixes)
	setOptionalStringField(fields, "notify_policy", req.NotifyPolicy)
	setOptionalStringField(fields, "kafka.partition", req.Partition)
	setOptionalStringField(fields, "kafka.partitioner", req.Partitioner)
	setOptionalStringField(fields, "kafka.rack_id", req.RackID)
	setOptionalBoolField(fields, "kafka.retry_as_batch", req.RetryAsBatch)
	setOptionalStringMapField(fields, "kafka.static_headers", req.StaticHeaders)
	setOptionalStringField(fields, "kafka.target_version", req.TargetVersion)
	setOptionalStringField(fields, "kafka.timeout", req.Timeout)
	return fields
}

func (c *Client) getKafkaFeature(ctx context.Context, outletID int64, featureID string) (OutletFeatureDetail, error) {
	var response OutletFeatureDetail
	if err := c.doJSON(ctx, httpMethodGet, fmt.Sprintf("/v1/terraform/kafka/%d/features/%s", outletID, featureID), nil, &response, statusOK); err != nil {
		return OutletFeatureDetail{}, err
	}
	return response, nil
}

func (c *Client) deleteKafkaFeature(ctx context.Context, outletID int64, featureID string) error {
	return c.doJSON(ctx, httpMethodDelete, fmt.Sprintf("/v1/terraform/kafka/%d/features/%s", outletID, featureID), nil, nil, statusOK)
}

func (c *Client) putKafkaTLSFeature(ctx context.Context, method string, outletID int64, req KafkaTLSFeatureRequest) error {
	caCert, err := decodeBase64File(req.CACertContentBase64)
	if err != nil {
		return err
	}
	inputs := map[string]any{
		kafkaCACertRole: payloadFileRef{File: kafkaTLSCACertPartName},
	}
	parts := []multipartFilePart{{
		FieldName: kafkaTLSCACertPartName,
		FileName:  req.CACertFilename,
		Content:   caCert,
	}}
	if req.ClientCertContentBase64 != "" {
		content, err := decodeBase64File(req.ClientCertContentBase64)
		if err != nil {
			return err
		}
		inputs[kafkaClientCertRole] = payloadFileRef{File: kafkaTLSClientCertPartName}
		parts = append(parts, multipartFilePart{FieldName: kafkaTLSClientCertPartName, FileName: req.ClientCertFilename, Content: content})
	}
	if req.ClientKeyContentBase64 != "" {
		content, err := decodeBase64File(req.ClientKeyContentBase64)
		if err != nil {
			return err
		}
		inputs[kafkaClientKeyRole] = payloadFileRef{File: kafkaTLSClientKeyPartName}
		parts = append(parts, multipartFilePart{FieldName: kafkaTLSClientKeyPartName, FileName: req.ClientKeyFilename, Content: content})
	}
	return c.doMultipartPayload(ctx, method, fmt.Sprintf("/v1/terraform/kafka/%d/features/%s", outletID, kafkaTLSFeatureID), map[string]any{"inputs": inputs}, parts, nil, statusOK)
}

func (c *Client) putKafkaSASLUserPassFeature(ctx context.Context, method string, outletID int64, featureID string, req KafkaSASLUserPassFeatureRequest) error {
	return c.doMultipartPayload(ctx, method, fmt.Sprintf("/v1/terraform/kafka/%d/features/%s", outletID, featureID), map[string]any{
		"inputs": map[string]any{
			"user":                req.User,
			kafkaSASLPasswordRole: req.Password,
		},
	}, nil, nil, statusOK)
}

func (c *Client) putKafkaOAuthFeature(ctx context.Context, method string, outletID int64, req KafkaSASLOAuthbearerFeatureRequest) error {
	return c.doMultipartPayload(ctx, method, fmt.Sprintf("/v1/terraform/kafka/%d/features/%s", outletID, kafkaSASLOAuthbearerFeatureID), map[string]any{
		"inputs": map[string]any{
			kafkaSASLAccessTokenRole: req.AccessToken,
		},
	}, nil, nil, statusOK)
}

func (c *Client) putKafkaOAuthCacheFeature(ctx context.Context, method string, outletID int64, req KafkaSASLOAuthbearerCacheFeatureRequest) error {
	return c.doMultipartPayload(ctx, method, fmt.Sprintf("/v1/terraform/kafka/%d/features/%s", outletID, kafkaSASLOAuthbearerCacheFeature), map[string]any{
		"inputs": map[string]any{
			"token_cache": req.TokenCache,
			"token_key":   req.TokenKey,
		},
	}, nil, nil, statusOK)
}

func (c *Client) putKafkaAwsMSKIAMFeature(ctx context.Context, method string, outletID int64, req KafkaSASLAwsMSKIAMFeatureRequest) error {
	inputs := map[string]any{
		"region": req.Region,
	}
	setOptionalStringInput(inputs, "endpoint", req.Endpoint)
	setOptionalStringInput(inputs, "profile", req.Profile)
	setOptionalStringInput(inputs, "id", req.ID)
	setOptionalStringInput(inputs, "from_ec2_role", req.FromEC2Role)
	setOptionalStringInput(inputs, "role", req.Role)
	setOptionalStringInput(inputs, "role_external_id", req.RoleExternalID)
	setOptionalStringInput(inputs, "expiry_window", req.ExpiryWindow)
	setOptionalStringInput(inputs, kafkaAwsSecretRole, req.AwsSecret)
	setOptionalStringInput(inputs, kafkaAwsTokenRole, req.AwsToken)
	return c.doMultipartPayload(ctx, method, fmt.Sprintf("/v1/terraform/kafka/%d/features/%s", outletID, kafkaSASLAwsMSKIAMFeatureID), map[string]any{
		"inputs": inputs,
	}, nil, nil, statusOK)
}

func setOptionalStringInput(inputs map[string]any, key string, value *string) {
	if value != nil {
		inputs[key] = *value
	}
}

func decodeKafkaOutlet(response typedOutletResponse) (KafkaOutlet, error) {
	fields := response.Config
	addresses, err := stringListField(fields, "kafka.addresses")
	if err != nil {
		return KafkaOutlet{}, err
	}
	topic, err := requiredStringField(fields, "kafka.topic")
	if err != nil {
		return KafkaOutlet{}, err
	}
	ackReplicas, err := boolFieldPointer(fields, "kafka.ack_replicas")
	if err != nil {
		return KafkaOutlet{}, err
	}
	backoffInitialInterval, err := stringFieldPointer(fields, "kafka.backoff.initial_interval")
	if err != nil {
		return KafkaOutlet{}, err
	}
	backoffMaxElapsedTime, err := stringFieldPointer(fields, "kafka.backoff.max_elapsed_time")
	if err != nil {
		return KafkaOutlet{}, err
	}
	backoffMaxInterval, err := stringFieldPointer(fields, "kafka.backoff.max_interval")
	if err != nil {
		return KafkaOutlet{}, err
	}
	batchingByteSize, err := int64FieldPointer(fields, "kafka.batching.byte_size")
	if err != nil {
		return KafkaOutlet{}, err
	}
	batchingCheck, err := stringFieldPointer(fields, "kafka.batching.check")
	if err != nil {
		return KafkaOutlet{}, err
	}
	batchingCount, err := int64FieldPointer(fields, "kafka.batching.count")
	if err != nil {
		return KafkaOutlet{}, err
	}
	batchingJitter, err := float64FieldPointer(fields, "kafka.batching.jitter")
	if err != nil {
		return KafkaOutlet{}, err
	}
	batchingPeriod, err := stringFieldPointer(fields, "kafka.batching.period")
	if err != nil {
		return KafkaOutlet{}, err
	}
	clientID, err := stringFieldPointer(fields, "kafka.client_id")
	if err != nil {
		return KafkaOutlet{}, err
	}
	compression, err := stringFieldPointer(fields, "kafka.compression")
	if err != nil {
		return KafkaOutlet{}, err
	}
	customTopicCreationEnabled, err := boolFieldPointer(fields, "kafka.custom_topic_creation.enabled")
	if err != nil {
		return KafkaOutlet{}, err
	}
	customTopicCreationPartitions, err := int64FieldPointer(fields, "kafka.custom_topic_creation.partitions")
	if err != nil {
		return KafkaOutlet{}, err
	}
	customTopicCreationReplicationFactor, err := int64FieldPointer(fields, "kafka.custom_topic_creation.replication_factor")
	if err != nil {
		return KafkaOutlet{}, err
	}
	idempotentWrite, err := boolFieldPointer(fields, "kafka.idempotent_write")
	if err != nil {
		return KafkaOutlet{}, err
	}
	injectTracingMap, err := stringFieldPointer(fields, "kafka.inject_tracing_map")
	if err != nil {
		return KafkaOutlet{}, err
	}
	key, err := stringFieldPointer(fields, "kafka.key")
	if err != nil {
		return KafkaOutlet{}, err
	}
	maxInFlight, err := int64FieldPointer(fields, "kafka.max_in_flight")
	if err != nil {
		return KafkaOutlet{}, err
	}
	maxMsgBytes, err := int64FieldPointer(fields, "kafka.max_msg_bytes")
	if err != nil {
		return KafkaOutlet{}, err
	}
	maxRetries, err := int64FieldPointer(fields, "kafka.max_retries")
	if err != nil {
		return KafkaOutlet{}, err
	}
	metadataExcludePrefixes, err := stringListField(fields, "kafka.metadata.exclude_prefixes")
	if err != nil {
		return KafkaOutlet{}, err
	}
	notifyPolicy, err := notifyPolicyPointerPreservingEmpty(fields)
	if err != nil {
		return KafkaOutlet{}, err
	}
	partition, err := stringFieldPointer(fields, "kafka.partition")
	if err != nil {
		return KafkaOutlet{}, err
	}
	partitioner, err := stringFieldPointer(fields, "kafka.partitioner")
	if err != nil {
		return KafkaOutlet{}, err
	}
	rackID, err := stringFieldPointer(fields, "kafka.rack_id")
	if err != nil {
		return KafkaOutlet{}, err
	}
	retryAsBatch, err := boolFieldPointer(fields, "kafka.retry_as_batch")
	if err != nil {
		return KafkaOutlet{}, err
	}
	staticHeaders, err := stringMapField(fields, "kafka.static_headers")
	if err != nil {
		return KafkaOutlet{}, err
	}
	targetVersion, err := stringFieldPointer(fields, "kafka.target_version")
	if err != nil {
		return KafkaOutlet{}, err
	}
	timeout, err := stringFieldPointer(fields, "kafka.timeout")
	if err != nil {
		return KafkaOutlet{}, err
	}
	return KafkaOutlet{
		ID:                                   response.ID,
		Name:                                 response.Name,
		Description:                          response.Description,
		Enabled:                              response.Enabled,
		AckReplicas:                          ackReplicas,
		Addresses:                            addresses,
		BackoffInitialInterval:               backoffInitialInterval,
		BackoffMaxElapsedTime:                backoffMaxElapsedTime,
		BackoffMaxInterval:                   backoffMaxInterval,
		BatchingByteSize:                     batchingByteSize,
		BatchingCheck:                        batchingCheck,
		BatchingCount:                        batchingCount,
		BatchingJitter:                       batchingJitter,
		BatchingPeriod:                       batchingPeriod,
		ClientID:                             clientID,
		Compression:                          compression,
		CustomTopicCreationEnabled:           customTopicCreationEnabled,
		CustomTopicCreationPartitions:        customTopicCreationPartitions,
		CustomTopicCreationReplicationFactor: customTopicCreationReplicationFactor,
		IdempotentWrite:                      idempotentWrite,
		InjectTracingMap:                     injectTracingMap,
		Key:                                  key,
		MaxInFlight:                          maxInFlight,
		MaxMsgBytes:                          maxMsgBytes,
		MaxRetries:                           maxRetries,
		MetadataExcludePrefixes:              metadataExcludePrefixes,
		NotifyPolicy:                         notifyPolicy,
		Partition:                            partition,
		Partitioner:                          partitioner,
		RackID:                               rackID,
		RetryAsBatch:                         retryAsBatch,
		StaticHeaders:                        staticHeaders,
		TargetVersion:                        targetVersion,
		Timeout:                              timeout,
		Topic:                                topic,
		ApplyStatus:                          response.ApplyStatus,
		LastApplyError:                       response.LastApplyError,
		ConfigGeneration:                     response.ConfigGeneration,
		AppliedGeneration:                    response.AppliedGeneration,
	}, nil
}

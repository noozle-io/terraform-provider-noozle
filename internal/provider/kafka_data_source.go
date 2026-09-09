package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

var (
	_ datasource.DataSource              = &kafkaOutletDataSource{}
	_ datasource.DataSourceWithConfigure = &kafkaOutletDataSource{}
)

type kafkaOutletDataSource struct{ client kafkaOutletClient }

type kafkaOutletDataSourceModel struct {
	ID                                   types.Int64   `tfsdk:"id"`
	Name                                 types.String  `tfsdk:"name"`
	Description                          types.String  `tfsdk:"description"`
	Enabled                              types.Bool    `tfsdk:"enabled"`
	AckReplicas                          types.Bool    `tfsdk:"ack_replicas"`
	Addresses                            types.List    `tfsdk:"addresses"`
	BackoffInitialInterval               types.String  `tfsdk:"backoff_initial_interval"`
	BackoffMaxElapsedTime                types.String  `tfsdk:"backoff_max_elapsed_time"`
	BackoffMaxInterval                   types.String  `tfsdk:"backoff_max_interval"`
	BatchingByteSize                     types.Int64   `tfsdk:"batching_byte_size"`
	BatchingCheck                        types.String  `tfsdk:"batching_check"`
	BatchingCount                        types.Int64   `tfsdk:"batching_count"`
	BatchingJitter                       types.Float64 `tfsdk:"batching_jitter"`
	BatchingPeriod                       types.String  `tfsdk:"batching_period"`
	ClientID                             types.String  `tfsdk:"client_id"`
	Compression                          types.String  `tfsdk:"compression"`
	CustomTopicCreationEnabled           types.Bool    `tfsdk:"custom_topic_creation_enabled"`
	CustomTopicCreationPartitions        types.Int64   `tfsdk:"custom_topic_creation_partitions"`
	CustomTopicCreationReplicationFactor types.Int64   `tfsdk:"custom_topic_creation_replication_factor"`
	IdempotentWrite                      types.Bool    `tfsdk:"idempotent_write"`
	InjectTracingMap                     types.String  `tfsdk:"inject_tracing_map"`
	Key                                  types.String  `tfsdk:"key"`
	MaxInFlight                          types.Int64   `tfsdk:"max_in_flight"`
	MaxMsgBytes                          types.Int64   `tfsdk:"max_msg_bytes"`
	MaxRetries                           types.Int64   `tfsdk:"max_retries"`
	MetadataExcludePrefixes              types.List    `tfsdk:"metadata_exclude_prefixes"`
	NotifyPolicy                         types.String  `tfsdk:"notify_policy"`
	Partition                            types.String  `tfsdk:"partition"`
	Partitioner                          types.String  `tfsdk:"partitioner"`
	RackID                               types.String  `tfsdk:"rack_id"`
	RetryAsBatch                         types.Bool    `tfsdk:"retry_as_batch"`
	StaticHeaders                        types.Map     `tfsdk:"static_headers"`
	TargetVersion                        types.String  `tfsdk:"target_version"`
	Timeout                              types.String  `tfsdk:"timeout"`
	Topic                                types.String  `tfsdk:"topic"`
	TLS                                  types.Object  `tfsdk:"tls"`
	SASLPlain                            types.Object  `tfsdk:"sasl_plain"`
	SASLScram256                         types.Object  `tfsdk:"sasl_scram256"`
	SASLScram512                         types.Object  `tfsdk:"sasl_scram512"`
	SASLOAuthbearer                      types.Object  `tfsdk:"sasl_oauthbearer"`
	SASLOAuthbearerCache                 types.Object  `tfsdk:"sasl_oauthbearer_cache"`
	SASLAwsMSKIAM                        types.Object  `tfsdk:"sasl_aws_msk_iam"`
	ApplyStatus                          types.String  `tfsdk:"apply_status"`
	LastApplyError                       types.String  `tfsdk:"last_apply_error"`
	ConfigGeneration                     types.Int64   `tfsdk:"config_generation"`
	AppliedGeneration                    types.Int64   `tfsdk:"applied_generation"`
	Transform                            types.Object  `tfsdk:"transform"`
}

func NewKafkaOutletDataSource() datasource.DataSource { return &kafkaOutletDataSource{} }

func (d *kafkaOutletDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kafka_outlet"
}

func (d *kafkaOutletDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Reads a Noozle Kafka outlet by ID.",
		Attributes: map[string]datasourceschema.Attribute{
			"id": datasourceschema.Int64Attribute{Required: true}, "name": datasourceschema.StringAttribute{Computed: true}, "description": datasourceschema.StringAttribute{Computed: true}, "enabled": datasourceschema.BoolAttribute{Computed: true},
			"ack_replicas": datasourceschema.BoolAttribute{Computed: true}, "addresses": datasourceschema.ListAttribute{Computed: true, ElementType: types.StringType},
			"backoff_initial_interval": datasourceschema.StringAttribute{Computed: true}, "backoff_max_elapsed_time": datasourceschema.StringAttribute{Computed: true}, "backoff_max_interval": datasourceschema.StringAttribute{Computed: true},
			"batching_byte_size": datasourceschema.Int64Attribute{Computed: true}, "batching_check": datasourceschema.StringAttribute{Computed: true}, "batching_count": datasourceschema.Int64Attribute{Computed: true},
			"batching_jitter": datasourceschema.Float64Attribute{Computed: true}, "batching_period": datasourceschema.StringAttribute{Computed: true}, "client_id": datasourceschema.StringAttribute{Computed: true}, "compression": datasourceschema.StringAttribute{Computed: true},
			"custom_topic_creation_enabled": datasourceschema.BoolAttribute{Computed: true}, "custom_topic_creation_partitions": datasourceschema.Int64Attribute{Computed: true}, "custom_topic_creation_replication_factor": datasourceschema.Int64Attribute{Computed: true},
			"idempotent_write": datasourceschema.BoolAttribute{Computed: true}, "inject_tracing_map": datasourceschema.StringAttribute{Computed: true}, "key": datasourceschema.StringAttribute{Computed: true},
			"max_in_flight": datasourceschema.Int64Attribute{Computed: true}, "max_msg_bytes": datasourceschema.Int64Attribute{Computed: true}, "max_retries": datasourceschema.Int64Attribute{Computed: true},
			"metadata_exclude_prefixes": datasourceschema.ListAttribute{Computed: true, ElementType: types.StringType}, "notify_policy": datasourceschema.StringAttribute{Computed: true}, "partition": datasourceschema.StringAttribute{Computed: true}, "partitioner": datasourceschema.StringAttribute{Computed: true},
			"rack_id": datasourceschema.StringAttribute{Computed: true}, "retry_as_batch": datasourceschema.BoolAttribute{Computed: true}, "static_headers": datasourceschema.MapAttribute{Computed: true, ElementType: types.StringType},
			"target_version": datasourceschema.StringAttribute{Computed: true}, "timeout": datasourceschema.StringAttribute{Computed: true}, "topic": datasourceschema.StringAttribute{Computed: true},
			"apply_status": datasourceschema.StringAttribute{Computed: true}, "last_apply_error": datasourceschema.StringAttribute{Computed: true}, "config_generation": datasourceschema.Int64Attribute{Computed: true}, "applied_generation": datasourceschema.Int64Attribute{Computed: true},
			"transform": datasourceschema.SingleNestedAttribute{Computed: true, Attributes: map[string]datasourceschema.Attribute{
				"selected": datasourceschema.StringAttribute{Computed: true}, "effective": datasourceschema.StringAttribute{Computed: true}, "state": datasourceschema.StringAttribute{Computed: true}, "applied": datasourceschema.BoolAttribute{Computed: true}, "generation": datasourceschema.Int64Attribute{Computed: true}, "last_error": datasourceschema.StringAttribute{Computed: true},
			}},
		},
		Blocks: map[string]datasourceschema.Block{
			"tls": datasourceschema.SingleNestedBlock{Attributes: map[string]datasourceschema.Attribute{
				"enabled": datasourceschema.BoolAttribute{Computed: true}, "state": datasourceschema.StringAttribute{Computed: true}, "desired_generation": datasourceschema.Int64Attribute{Computed: true}, "applied_generation": datasourceschema.Int64Attribute{Computed: true}, "last_error": datasourceschema.StringAttribute{Computed: true},
				"ca_cert_filename": datasourceschema.StringAttribute{Computed: true}, "ca_cert_checksum_sha256": datasourceschema.StringAttribute{Computed: true},
				"client_cert_filename": datasourceschema.StringAttribute{Computed: true}, "client_cert_checksum_sha256": datasourceschema.StringAttribute{Computed: true},
				"client_key_filename": datasourceschema.StringAttribute{Computed: true}, "client_key_checksum_sha256": datasourceschema.StringAttribute{Computed: true},
			}},
			"sasl_plain":             datasourceschema.SingleNestedBlock{Attributes: kafkaDataSourceFeatureAttributes()},
			"sasl_scram256":          datasourceschema.SingleNestedBlock{Attributes: kafkaDataSourceFeatureAttributes()},
			"sasl_scram512":          datasourceschema.SingleNestedBlock{Attributes: kafkaDataSourceFeatureAttributes()},
			"sasl_oauthbearer":       datasourceschema.SingleNestedBlock{Attributes: kafkaDataSourceFeatureAttributes()},
			"sasl_oauthbearer_cache": datasourceschema.SingleNestedBlock{Attributes: kafkaDataSourceFeatureAttributes()},
			"sasl_aws_msk_iam":       datasourceschema.SingleNestedBlock{Attributes: kafkaDataSourceFeatureAttributes()},
		},
	}
}

func kafkaDataSourceFeatureAttributes() map[string]datasourceschema.Attribute {
	return map[string]datasourceschema.Attribute{
		"enabled":            datasourceschema.BoolAttribute{Computed: true},
		"state":              datasourceschema.StringAttribute{Computed: true},
		"desired_generation": datasourceschema.Int64Attribute{Computed: true},
		"applied_generation": datasourceschema.Int64Attribute{Computed: true},
		"last_error":         datasourceschema.StringAttribute{Computed: true},
	}
}

func (d *kafkaOutletDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	outletClient, ok := req.ProviderData.(kafkaOutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected Kafka outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}
	d.client = outletClient
}

func (d *kafkaOutletDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var config kafkaOutletDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	outlet, err := d.client.GetKafkaOutlet(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Kafka Outlet", err.Error())
		return
	}
	credentials, err := d.client.ListKafkaCredentials(ctx, config.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Read Kafka Credentials", err.Error())
		return
	}
	tlsFeature, plainFeature, scram256Feature, scram512Feature, oauthFeature, oauthCacheFeature, awsFeature, diags := (&kafkaOutletResource{client: d.client}).readFeatures(ctx, config.ID.ValueInt64())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tlsObject, diags := buildKafkaTLSDataSourceObject(tlsFeature, credentials)
	resp.Diagnostics.Append(diags...)
	plainObject, diags := buildKafkaFeatureDataSourceObject(plainFeature)
	resp.Diagnostics.Append(diags...)
	scram256Object, diags := buildKafkaFeatureDataSourceObject(scram256Feature)
	resp.Diagnostics.Append(diags...)
	scram512Object, diags := buildKafkaFeatureDataSourceObject(scram512Feature)
	resp.Diagnostics.Append(diags...)
	oauthObject, diags := buildKafkaFeatureDataSourceObject(oauthFeature)
	resp.Diagnostics.Append(diags...)
	oauthCacheObject, diags := buildKafkaFeatureDataSourceObject(oauthCacheFeature)
	resp.Diagnostics.Append(diags...)
	awsObject, diags := buildKafkaFeatureDataSourceObject(awsFeature)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	transform, err := d.client.GetTypedOutletTransform(ctx, "kafka", config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Kafka Transform", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &kafkaOutletDataSourceModel{
		ID: types.Int64Value(outlet.ID), Name: types.StringValue(outlet.Name), Description: nullableString(outlet.Description), Enabled: types.BoolValue(outlet.Enabled),
		AckReplicas: nullableBoolPointer(outlet.AckReplicas), Addresses: nullableStringList(outlet.Addresses), BackoffInitialInterval: nullableStringPointer(outlet.BackoffInitialInterval),
		BackoffMaxElapsedTime: nullableStringPointer(outlet.BackoffMaxElapsedTime), BackoffMaxInterval: nullableStringPointer(outlet.BackoffMaxInterval), BatchingByteSize: nullableInt64Pointer(outlet.BatchingByteSize),
		BatchingCheck: nullableStringPointer(outlet.BatchingCheck), BatchingCount: nullableInt64Pointer(outlet.BatchingCount), BatchingJitter: nullableFloat64Pointer(outlet.BatchingJitter),
		BatchingPeriod: nullableStringPointer(outlet.BatchingPeriod), ClientID: nullableStringPointer(outlet.ClientID), Compression: nullableStringPointer(outlet.Compression),
		CustomTopicCreationEnabled: nullableBoolPointer(outlet.CustomTopicCreationEnabled), CustomTopicCreationPartitions: nullableInt64Pointer(outlet.CustomTopicCreationPartitions),
		CustomTopicCreationReplicationFactor: nullableInt64Pointer(outlet.CustomTopicCreationReplicationFactor), IdempotentWrite: nullableBoolPointer(outlet.IdempotentWrite), InjectTracingMap: nullableStringPointer(outlet.InjectTracingMap),
		Key: nullableStringPointer(outlet.Key), MaxInFlight: nullableInt64Pointer(outlet.MaxInFlight), MaxMsgBytes: nullableInt64Pointer(outlet.MaxMsgBytes), MaxRetries: nullableInt64Pointer(outlet.MaxRetries),
		MetadataExcludePrefixes: nullableStringList(outlet.MetadataExcludePrefixes), NotifyPolicy: nullableStringPointerPreservingEmpty(outlet.NotifyPolicy), Partition: nullableStringPointer(outlet.Partition), Partitioner: nullableStringPointer(outlet.Partitioner), RackID: nullableStringPointer(outlet.RackID),
		RetryAsBatch: nullableBoolPointer(outlet.RetryAsBatch), StaticHeaders: nullableStringMap(outlet.StaticHeaders), TargetVersion: nullableStringPointer(outlet.TargetVersion), Timeout: nullableStringPointer(outlet.Timeout), Topic: types.StringValue(outlet.Topic),
		TLS: tlsObject, SASLPlain: plainObject, SASLScram256: scram256Object, SASLScram512: scram512Object, SASLOAuthbearer: oauthObject, SASLOAuthbearerCache: oauthCacheObject, SASLAwsMSKIAM: awsObject,
		ApplyStatus: nullableString(outlet.ApplyStatus), LastApplyError: nullableString(outlet.LastApplyError), ConfigGeneration: types.Int64Value(outlet.ConfigGeneration), AppliedGeneration: types.Int64Value(outlet.AppliedGeneration),
		Transform: transformObject(transform),
	})...)
}

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

var (
	_ resource.Resource                   = &kafkaOutletResource{}
	_ resource.ResourceWithConfigure      = &kafkaOutletResource{}
	_ resource.ResourceWithImportState    = &kafkaOutletResource{}
	_ resource.ResourceWithValidateConfig = &kafkaOutletResource{}
)

type kafkaOutletClient interface {
	CreateKafkaOutlet(context.Context, client.KafkaOutletCreateRequest) (client.KafkaOutlet, error)
	GetKafkaOutlet(context.Context, int64) (client.KafkaOutlet, error)
	UpdateKafkaOutlet(context.Context, int64, client.KafkaOutletUpdateRequest) error
	DeleteKafkaOutlet(context.Context, int64) error
	SetKafkaField(context.Context, int64, string, any) error
	DeleteKafkaField(context.Context, int64, string) error
	DeleteKafkaOutletOperation(context.Context, int64) (client.OutletOperation, error)
	GetTypedOutletStatus(context.Context, string, int64, int64) (client.OutletStatus, error)
	ListKafkaCredentials(context.Context, int64) ([]client.OutletMaterialSummary, error)
	DeleteKafkaCredential(context.Context, int64, string) (client.OutletMaterialDeleteResponse, error)
	GetKafkaTLSFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureKafkaTLSFeature(context.Context, int64, client.KafkaTLSFeatureRequest) error
	UpdateKafkaTLSFeature(context.Context, int64, client.KafkaTLSFeatureRequest) error
	DeleteKafkaTLSFeature(context.Context, int64) error
	GetKafkaSASLPlainFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureKafkaSASLPlainFeature(context.Context, int64, client.KafkaSASLUserPassFeatureRequest) error
	UpdateKafkaSASLPlainFeature(context.Context, int64, client.KafkaSASLUserPassFeatureRequest) error
	DeleteKafkaSASLPlainFeature(context.Context, int64) error
	GetKafkaSASLScram256Feature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureKafkaSASLScram256Feature(context.Context, int64, client.KafkaSASLUserPassFeatureRequest) error
	UpdateKafkaSASLScram256Feature(context.Context, int64, client.KafkaSASLUserPassFeatureRequest) error
	DeleteKafkaSASLScram256Feature(context.Context, int64) error
	GetKafkaSASLScram512Feature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureKafkaSASLScram512Feature(context.Context, int64, client.KafkaSASLUserPassFeatureRequest) error
	UpdateKafkaSASLScram512Feature(context.Context, int64, client.KafkaSASLUserPassFeatureRequest) error
	DeleteKafkaSASLScram512Feature(context.Context, int64) error
	GetKafkaSASLOAuthbearerFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureKafkaSASLOAuthbearerFeature(context.Context, int64, client.KafkaSASLOAuthbearerFeatureRequest) error
	UpdateKafkaSASLOAuthbearerFeature(context.Context, int64, client.KafkaSASLOAuthbearerFeatureRequest) error
	DeleteKafkaSASLOAuthbearerFeature(context.Context, int64) error
	GetKafkaSASLOAuthbearerCacheFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureKafkaSASLOAuthbearerCacheFeature(context.Context, int64, client.KafkaSASLOAuthbearerCacheFeatureRequest) error
	UpdateKafkaSASLOAuthbearerCacheFeature(context.Context, int64, client.KafkaSASLOAuthbearerCacheFeatureRequest) error
	DeleteKafkaSASLOAuthbearerCacheFeature(context.Context, int64) error
	GetKafkaSASLAwsMSKIAMFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureKafkaSASLAwsMSKIAMFeature(context.Context, int64, client.KafkaSASLAwsMSKIAMFeatureRequest) error
	UpdateKafkaSASLAwsMSKIAMFeature(context.Context, int64, client.KafkaSASLAwsMSKIAMFeatureRequest) error
	DeleteKafkaSASLAwsMSKIAMFeature(context.Context, int64) error
	typedOutletTransformClient
}

type kafkaOutletFieldUpdater struct{ client kafkaOutletClient }

func (u kafkaOutletFieldUpdater) SetField(ctx context.Context, outletID int64, field string, value any) error {
	return u.client.SetKafkaField(ctx, outletID, field, value)
}

func (u kafkaOutletFieldUpdater) DeleteField(ctx context.Context, outletID int64, field string) error {
	return u.client.DeleteKafkaField(ctx, outletID, field)
}

type kafkaOutletResource struct {
	outletLifecycle[kafkaOutletResourceModel]
	client kafkaOutletClient
}

type kafkaOutletAdapter struct{ resource *kafkaOutletResource }

type kafkaOutletResourceModel struct {
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

func NewKafkaOutletResource() resource.Resource {
	r := &kafkaOutletResource{}
	r.outletLifecycle = newOutletLifecycle("Kafka Outlet", func(providerData any) (outletLifecycleAdapter[kafkaOutletResourceModel], diag.Diagnostics) {
		outletClient, ok := providerData.(kafkaOutletClient)
		if !ok {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Unexpected Provider Data Type", fmt.Sprintf("Expected Kafka outlet client, got: %T. Please report this provider bug.", providerData))}
		}
		r.client = outletClient
		return kafkaOutletAdapter{resource: r}, nil
	})
	return r
}

func (r *kafkaOutletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kafka_outlet"
}

func (r *kafkaOutletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		MarkdownDescription: "Manages a Noozle Kafka outlet.",
		Attributes: map[string]resourceschema.Attribute{
			"id":                               resourceschema.Int64Attribute{Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}},
			"name":                             resourceschema.StringAttribute{Required: true},
			"description":                      resourceschema.StringAttribute{Optional: true},
			"enabled":                          resourceschema.BoolAttribute{Required: true},
			"ack_replicas":                     resourceschema.BoolAttribute{Optional: true},
			"addresses":                        resourceschema.ListAttribute{Required: true, ElementType: types.StringType},
			"backoff_initial_interval":         resourceschema.StringAttribute{Optional: true},
			"backoff_max_elapsed_time":         resourceschema.StringAttribute{Optional: true},
			"backoff_max_interval":             resourceschema.StringAttribute{Optional: true},
			"batching_byte_size":               resourceschema.Int64Attribute{Optional: true},
			"batching_check":                   resourceschema.StringAttribute{Optional: true},
			"batching_count":                   resourceschema.Int64Attribute{Optional: true},
			"batching_jitter":                  resourceschema.Float64Attribute{Optional: true},
			"batching_period":                  resourceschema.StringAttribute{Optional: true},
			"client_id":                        resourceschema.StringAttribute{Optional: true},
			"compression":                      resourceschema.StringAttribute{Optional: true},
			"custom_topic_creation_enabled":    resourceschema.BoolAttribute{Optional: true},
			"custom_topic_creation_partitions": resourceschema.Int64Attribute{Optional: true},
			"custom_topic_creation_replication_factor": resourceschema.Int64Attribute{Optional: true},
			"idempotent_write":                         resourceschema.BoolAttribute{Optional: true},
			"inject_tracing_map":                       resourceschema.StringAttribute{Optional: true},
			"key":                                      resourceschema.StringAttribute{Optional: true},
			"max_in_flight":                            resourceschema.Int64Attribute{Optional: true},
			"max_msg_bytes":                            resourceschema.Int64Attribute{Optional: true},
			"max_retries":                              resourceschema.Int64Attribute{Optional: true},
			"metadata_exclude_prefixes":                resourceschema.ListAttribute{Optional: true, ElementType: types.StringType},
			"notify_policy":                            resourceschema.StringAttribute{Optional: true},
			"partition":                                resourceschema.StringAttribute{Optional: true},
			"partitioner":                              resourceschema.StringAttribute{Optional: true},
			"rack_id":                                  resourceschema.StringAttribute{Optional: true},
			"retry_as_batch":                           resourceschema.BoolAttribute{Optional: true},
			"static_headers":                           resourceschema.MapAttribute{Optional: true, ElementType: types.StringType},
			"target_version":                           resourceschema.StringAttribute{Optional: true},
			"timeout":                                  resourceschema.StringAttribute{Optional: true},
			"topic":                                    resourceschema.StringAttribute{Required: true},
			"apply_status":                             resourceschema.StringAttribute{Computed: true},
			"last_apply_error":                         resourceschema.StringAttribute{Computed: true},
			"config_generation":                        resourceschema.Int64Attribute{Computed: true},
			"applied_generation":                       resourceschema.Int64Attribute{Computed: true},
			"transform": resourceschema.SingleNestedAttribute{MarkdownDescription: "Optional outlet transform selection and observed runtime state.", Optional: true, Computed: true, Attributes: map[string]resourceschema.Attribute{
				"selected": resourceschema.StringAttribute{MarkdownDescription: "Explicitly selected transform preset.", Optional: true}, "effective": resourceschema.StringAttribute{MarkdownDescription: "Effective transform after defaults are applied.", Computed: true}, "state": resourceschema.StringAttribute{MarkdownDescription: "Observed transform runtime state.", Computed: true}, "applied": resourceschema.BoolAttribute{MarkdownDescription: "Whether the transform has been applied.", Computed: true}, "generation": resourceschema.Int64Attribute{MarkdownDescription: "Observed transform generation.", Computed: true}, "last_error": resourceschema.StringAttribute{MarkdownDescription: "Last observed transform apply error.", Computed: true},
			}},
		},
		Blocks: map[string]resourceschema.Block{
			"tls": resourceschema.SingleNestedBlock{Attributes: map[string]resourceschema.Attribute{
				"enabled": resourceschema.BoolAttribute{Computed: true}, "state": resourceschema.StringAttribute{Computed: true}, "desired_generation": resourceschema.Int64Attribute{Computed: true}, "applied_generation": resourceschema.Int64Attribute{Computed: true}, "last_error": resourceschema.StringAttribute{Computed: true},
				"ca_cert_filename": resourceschema.StringAttribute{Optional: true}, "ca_cert_content_base64": resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true}, "ca_cert_version": resourceschema.Int64Attribute{Optional: true}, "ca_cert_checksum_sha256": resourceschema.StringAttribute{Computed: true},
				"client_cert_filename": resourceschema.StringAttribute{Optional: true}, "client_cert_content_base64": resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true}, "client_cert_version": resourceschema.Int64Attribute{Optional: true}, "client_cert_checksum_sha256": resourceschema.StringAttribute{Computed: true},
				"client_key_filename": resourceschema.StringAttribute{Optional: true}, "client_key_content_base64": resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true}, "client_key_version": resourceschema.Int64Attribute{Optional: true}, "client_key_checksum_sha256": resourceschema.StringAttribute{Computed: true},
			}},
			"sasl_plain": resourceschema.SingleNestedBlock{Attributes: map[string]resourceschema.Attribute{
				"enabled": resourceschema.BoolAttribute{Computed: true}, "state": resourceschema.StringAttribute{Computed: true}, "desired_generation": resourceschema.Int64Attribute{Computed: true}, "applied_generation": resourceschema.Int64Attribute{Computed: true}, "last_error": resourceschema.StringAttribute{Computed: true},
				"user": resourceschema.StringAttribute{Optional: true}, "password": resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true}, "password_version": resourceschema.Int64Attribute{Optional: true},
			}},
			"sasl_scram256": resourceschema.SingleNestedBlock{Attributes: map[string]resourceschema.Attribute{
				"enabled": resourceschema.BoolAttribute{Computed: true}, "state": resourceschema.StringAttribute{Computed: true}, "desired_generation": resourceschema.Int64Attribute{Computed: true}, "applied_generation": resourceschema.Int64Attribute{Computed: true}, "last_error": resourceschema.StringAttribute{Computed: true},
				"user": resourceschema.StringAttribute{Optional: true}, "password": resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true}, "password_version": resourceschema.Int64Attribute{Optional: true},
			}},
			"sasl_scram512": resourceschema.SingleNestedBlock{Attributes: map[string]resourceschema.Attribute{
				"enabled": resourceschema.BoolAttribute{Computed: true}, "state": resourceschema.StringAttribute{Computed: true}, "desired_generation": resourceschema.Int64Attribute{Computed: true}, "applied_generation": resourceschema.Int64Attribute{Computed: true}, "last_error": resourceschema.StringAttribute{Computed: true},
				"user": resourceschema.StringAttribute{Optional: true}, "password": resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true}, "password_version": resourceschema.Int64Attribute{Optional: true},
			}},
			"sasl_oauthbearer": resourceschema.SingleNestedBlock{Attributes: map[string]resourceschema.Attribute{
				"enabled": resourceschema.BoolAttribute{Computed: true}, "state": resourceschema.StringAttribute{Computed: true}, "desired_generation": resourceschema.Int64Attribute{Computed: true}, "applied_generation": resourceschema.Int64Attribute{Computed: true}, "last_error": resourceschema.StringAttribute{Computed: true},
				"access_token": resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true}, "access_token_version": resourceschema.Int64Attribute{Optional: true},
			}},
			"sasl_oauthbearer_cache": resourceschema.SingleNestedBlock{Attributes: map[string]resourceschema.Attribute{
				"enabled": resourceschema.BoolAttribute{Computed: true}, "state": resourceschema.StringAttribute{Computed: true}, "desired_generation": resourceschema.Int64Attribute{Computed: true}, "applied_generation": resourceschema.Int64Attribute{Computed: true}, "last_error": resourceschema.StringAttribute{Computed: true},
				"token_cache": resourceschema.StringAttribute{Optional: true}, "token_key": resourceschema.StringAttribute{Optional: true},
			}},
			"sasl_aws_msk_iam": resourceschema.SingleNestedBlock{Attributes: map[string]resourceschema.Attribute{
				"enabled": resourceschema.BoolAttribute{Computed: true}, "state": resourceschema.StringAttribute{Computed: true}, "desired_generation": resourceschema.Int64Attribute{Computed: true}, "applied_generation": resourceschema.Int64Attribute{Computed: true}, "last_error": resourceschema.StringAttribute{Computed: true},
				"region": resourceschema.StringAttribute{Optional: true}, "endpoint": resourceschema.StringAttribute{Optional: true}, "profile": resourceschema.StringAttribute{Optional: true}, "id": resourceschema.StringAttribute{Optional: true},
				"from_ec2_role": resourceschema.StringAttribute{Optional: true}, "role": resourceschema.StringAttribute{Optional: true}, "role_external_id": resourceschema.StringAttribute{Optional: true}, "expiry_window": resourceschema.StringAttribute{Optional: true},
				"aws_secret": resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true}, "aws_secret_version": resourceschema.Int64Attribute{Optional: true},
				"aws_token": resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true}, "aws_token_version": resourceschema.Int64Attribute{Optional: true},
			}},
		},
	}
}

func (r *kafkaOutletResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config kafkaOutletResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	authBlocks := 0
	if !config.TLS.IsNull() && !config.TLS.IsUnknown() {
		tls, diags := expandKafkaTLS(ctx, config.TLS)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if tls.CACertFilename.IsNull() || tls.CACertContentBase64.IsNull() || tls.CACertVersion.IsNull() ||
			tls.CACertFilename.IsUnknown() || tls.CACertContentBase64.IsUnknown() || tls.CACertVersion.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete Kafka TLS Configuration", "When `tls` is configured, `ca_cert_filename`, `ca_cert_content_base64`, and `ca_cert_version` must all be set.")
		}
		if anyKafkaTLSClientMaterialConfigured(tls) && !fullKafkaTLSClientMaterialConfigured(tls) {
			resp.Diagnostics.AddError("Incomplete Kafka TLS Client Material", "When using optional Kafka TLS client credentials, set `client_cert_*` and `client_key_*` filename/content/version values together.")
		}
	}
	if !config.SASLPlain.IsNull() && !config.SASLPlain.IsUnknown() {
		authBlocks++
		validateKafkaSASLUserPass(ctx, config.SASLPlain, "Kafka SASL/PLAIN", &resp.Diagnostics)
	}
	if !config.SASLScram256.IsNull() && !config.SASLScram256.IsUnknown() {
		authBlocks++
		validateKafkaSASLUserPass(ctx, config.SASLScram256, "Kafka SASL SCRAM-256", &resp.Diagnostics)
	}
	if !config.SASLScram512.IsNull() && !config.SASLScram512.IsUnknown() {
		authBlocks++
		validateKafkaSASLUserPass(ctx, config.SASLScram512, "Kafka SASL SCRAM-512", &resp.Diagnostics)
	}
	if !config.SASLOAuthbearer.IsNull() && !config.SASLOAuthbearer.IsUnknown() {
		authBlocks++
		oauth, diags := expandKafkaSASLOAuthbearer(ctx, config.SASLOAuthbearer)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if oauth.AccessToken.IsNull() || oauth.AccessToken.IsUnknown() || oauth.AccessTokenVersion.IsNull() || oauth.AccessTokenVersion.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete Kafka SASL OAUTHBEARER Configuration", "When `sasl_oauthbearer` is configured, `access_token` and `access_token_version` must both be set.")
		}
	}
	if !config.SASLOAuthbearerCache.IsNull() && !config.SASLOAuthbearerCache.IsUnknown() {
		authBlocks++
		cache, diags := expandKafkaSASLOAuthbearerCache(ctx, config.SASLOAuthbearerCache)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if cache.TokenCache.IsNull() || cache.TokenCache.IsUnknown() || cache.TokenKey.IsNull() || cache.TokenKey.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete Kafka SASL OAUTHBEARER Cache Configuration", "When `sasl_oauthbearer_cache` is configured, `token_cache` and `token_key` must both be set.")
		}
	}
	if !config.SASLAwsMSKIAM.IsNull() && !config.SASLAwsMSKIAM.IsUnknown() {
		authBlocks++
		aws, diags := expandKafkaSASLAwsMSKIAM(ctx, config.SASLAwsMSKIAM)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if aws.Region.IsNull() || aws.Region.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete Kafka SASL AWS_MSK_IAM Configuration", "When `sasl_aws_msk_iam` is configured, `region` must be set.")
		}
		if (aws.AwsSecret.IsNull() != aws.AwsSecretVersion.IsNull()) || (aws.AwsSecret.IsUnknown() != aws.AwsSecretVersion.IsUnknown()) {
			resp.Diagnostics.AddError("Incomplete Kafka AWS Secret Configuration", "When `aws_secret` is configured, `aws_secret_version` must also be set.")
		}
		if (aws.AwsToken.IsNull() != aws.AwsTokenVersion.IsNull()) || (aws.AwsToken.IsUnknown() != aws.AwsTokenVersion.IsUnknown()) {
			resp.Diagnostics.AddError("Incomplete Kafka AWS Token Configuration", "When `aws_token` is configured, `aws_token_version` must also be set.")
		}
	}
	if authBlocks > 1 {
		resp.Diagnostics.AddError("Conflicting Kafka SASL Features", "Only one of `sasl_plain`, `sasl_scram256`, `sasl_scram512`, `sasl_oauthbearer`, `sasl_oauthbearer_cache`, or `sasl_aws_msk_iam` can be configured at the same time.")
	}
}

func (a kafkaOutletAdapter) Create(ctx context.Context, config tfsdk.Config, plan kafkaOutletResourceModel) (outletLifecycleResult[kafkaOutletResourceModel], diag.Diagnostics) {
	r := a.resource
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	var diags diag.Diagnostics
	cfg, ok := r.readConfig(ctx, config, plan, &diags)
	if !ok {
		return outletLifecycleResult[kafkaOutletResourceModel]{}, diags
	}
	configuredTransform, transformDiags := selectedTransformFromConfig(ctx, config)
	diags.Append(transformDiags...)
	if diags.HasError() {
		return outletLifecycleResult[kafkaOutletResourceModel]{}, diags
	}
	addresses, headers, metadataExcludePrefixes, ok := r.expandRootValues(ctx, plan, &diags)
	if !ok {
		return outletLifecycleResult[kafkaOutletResourceModel]{}, diags
	}
	outlet, err := r.client.CreateKafkaOutlet(ctx, client.KafkaOutletCreateRequest{
		Name: plan.Name.ValueString(), Description: stringValueOrEmpty(plan.Description), Enabled: plan.Enabled.ValueBool(),
		AckReplicas: boolPointerValue(plan.AckReplicas), Addresses: addresses, BackoffInitialInterval: stringPointerValue(plan.BackoffInitialInterval),
		BackoffMaxElapsedTime: stringPointerValue(plan.BackoffMaxElapsedTime), BackoffMaxInterval: stringPointerValue(plan.BackoffMaxInterval),
		BatchingByteSize: int64PointerValue(plan.BatchingByteSize), BatchingCheck: stringPointerValue(plan.BatchingCheck), BatchingCount: int64PointerValue(plan.BatchingCount),
		BatchingJitter: float64PointerValue(plan.BatchingJitter), BatchingPeriod: stringPointerValue(plan.BatchingPeriod), ClientID: stringPointerValue(plan.ClientID),
		Compression: stringPointerValue(plan.Compression), CustomTopicCreationEnabled: boolPointerValue(plan.CustomTopicCreationEnabled),
		CustomTopicCreationPartitions: int64PointerValue(plan.CustomTopicCreationPartitions), CustomTopicCreationReplicationFactor: int64PointerValue(plan.CustomTopicCreationReplicationFactor),
		IdempotentWrite: boolPointerValue(plan.IdempotentWrite), InjectTracingMap: stringPointerValue(plan.InjectTracingMap), Key: stringPointerValue(plan.Key),
		MaxInFlight: int64PointerValue(plan.MaxInFlight), MaxMsgBytes: int64PointerValue(plan.MaxMsgBytes), MaxRetries: int64PointerValue(plan.MaxRetries),
		MetadataExcludePrefixes: metadataExcludePrefixes, NotifyPolicy: stringPointerValue(plan.NotifyPolicy), Partition: stringPointerValue(plan.Partition), Partitioner: stringPointerValue(plan.Partitioner),
		RackID: stringPointerValue(plan.RackID), RetryAsBatch: boolPointerValue(plan.RetryAsBatch), StaticHeaders: headers, TargetVersion: stringPointerValue(plan.TargetVersion),
		Timeout: stringPointerValue(plan.Timeout), Topic: plan.Topic.ValueString(),
	})
	if err != nil {
		return outletLifecycleResult[kafkaOutletResourceModel]{}, diag.Diagnostics{diag.NewErrorDiagnostic("Unable to Create Kafka Outlet", err.Error())}
	}
	if err := r.applyFeatureConfig(ctx, outlet.ID, kafkaFeatureState{}, cfg); err != nil {
		if deleteErr := r.client.DeleteKafkaOutlet(ctx, outlet.ID); deleteErr != nil && !client.IsNotFound(deleteErr) {
			return outletLifecycleResult[kafkaOutletResourceModel]{}, diag.Diagnostics{diag.NewErrorDiagnostic("Unable to Configure Kafka Features", fmt.Sprintf("%s; cleanup failed: %s", err, deleteErr))}
		}
		return outletLifecycleResult[kafkaOutletResourceModel]{}, diag.Diagnostics{diag.NewErrorDiagnostic("Unable to Configure Kafka Features", err.Error())}
	}
	if !configuredTransform.IsNull() && !configuredTransform.IsUnknown() {
		if _, err := r.client.SetTypedOutletTransform(ctx, "kafka", outlet.ID, configuredTransform.ValueString()); err != nil {
			if deleteErr := r.client.DeleteKafkaOutlet(ctx, outlet.ID); deleteErr != nil && !client.IsNotFound(deleteErr) {
				return outletLifecycleResult[kafkaOutletResourceModel]{}, diag.Diagnostics{diag.NewErrorDiagnostic("Unable to Configure Kafka Transform", fmt.Sprintf("%s; cleanup failed: %s", err, deleteErr))}
			}
			return outletLifecycleResult[kafkaOutletResourceModel]{}, diag.Diagnostics{diag.NewErrorDiagnostic("Unable to Configure Kafka Transform", err.Error())}
		}
	}
	newState, removed, diags := r.buildState(ctx, outlet.ID, plan.TLS, plan.SASLPlain, plan.SASLScram256, plan.SASLScram512, plan.SASLOAuthbearer, plan.SASLOAuthbearerCache, plan.SASLAwsMSKIAM)
	if diags.HasError() {
		return outletLifecycleResult[kafkaOutletResourceModel]{}, diags
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "kafka", outlet.ID, plan.Enabled.ValueBool()); err != nil {
		return outletLifecycleResult[kafkaOutletResourceModel]{}, diag.Diagnostics{diag.NewErrorDiagnostic("Unable to Observe Kafka Outlet Lifecycle", err.Error())}
	}
	return outletLifecycleResult[kafkaOutletResourceModel]{State: newState, Removed: removed}, nil
}

func (a kafkaOutletAdapter) Read(ctx context.Context, state kafkaOutletResourceModel) (outletLifecycleResult[kafkaOutletResourceModel], diag.Diagnostics) {
	r := a.resource
	newState, removed, diags := r.buildState(ctx, state.ID.ValueInt64(), state.TLS, state.SASLPlain, state.SASLScram256, state.SASLScram512, state.SASLOAuthbearer, state.SASLOAuthbearerCache, state.SASLAwsMSKIAM)
	if diags.HasError() {
		return outletLifecycleResult[kafkaOutletResourceModel]{}, diags
	}
	return outletLifecycleResult[kafkaOutletResourceModel]{State: newState, Removed: removed}, nil
}

func (a kafkaOutletAdapter) Update(ctx context.Context, config tfsdk.Config, plan, state kafkaOutletResourceModel) (result outletLifecycleResult[kafkaOutletResourceModel], diags diag.Diagnostics) {
	r := a.resource
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	var resp struct{ Diagnostics diag.Diagnostics }
	defer func() { diags = resp.Diagnostics }()
	cfg, ok := r.readConfig(ctx, config, plan, &resp.Diagnostics)
	if !ok {
		return
	}
	configuredTransform, transformDiags := selectedTransformFromConfig(ctx, config)
	resp.Diagnostics.Append(transformDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	addresses, headers, metadataExcludePrefixes, ok := r.expandRootValues(ctx, plan, &resp.Diagnostics)
	if !ok {
		return
	}
	if err := r.client.UpdateKafkaOutlet(ctx, state.ID.ValueInt64(), client.KafkaOutletUpdateRequest{
		Name: plan.Name.ValueString(), Description: stringValueOrEmpty(plan.Description), Enabled: plan.Enabled.ValueBool(),
		AckReplicas: boolPointerValue(plan.AckReplicas), Addresses: addresses, BackoffInitialInterval: stringPointerValue(plan.BackoffInitialInterval),
		BackoffMaxElapsedTime: stringPointerValue(plan.BackoffMaxElapsedTime), BackoffMaxInterval: stringPointerValue(plan.BackoffMaxInterval),
		BatchingByteSize: int64PointerValue(plan.BatchingByteSize), BatchingCheck: stringPointerValue(plan.BatchingCheck), BatchingCount: int64PointerValue(plan.BatchingCount),
		BatchingJitter: float64PointerValue(plan.BatchingJitter), BatchingPeriod: stringPointerValue(plan.BatchingPeriod), ClientID: stringPointerValue(plan.ClientID),
		Compression: stringPointerValue(plan.Compression), CustomTopicCreationEnabled: boolPointerValue(plan.CustomTopicCreationEnabled),
		CustomTopicCreationPartitions: int64PointerValue(plan.CustomTopicCreationPartitions), CustomTopicCreationReplicationFactor: int64PointerValue(plan.CustomTopicCreationReplicationFactor),
		IdempotentWrite: boolPointerValue(plan.IdempotentWrite), InjectTracingMap: stringPointerValue(plan.InjectTracingMap), Key: stringPointerValue(plan.Key),
		MaxInFlight: int64PointerValue(plan.MaxInFlight), MaxMsgBytes: int64PointerValue(plan.MaxMsgBytes), MaxRetries: int64PointerValue(plan.MaxRetries),
		MetadataExcludePrefixes: metadataExcludePrefixes, NotifyPolicy: stringPointerValue(plan.NotifyPolicy), Partition: stringPointerValue(plan.Partition), Partitioner: stringPointerValue(plan.Partitioner),
		RackID: stringPointerValue(plan.RackID), RetryAsBatch: boolPointerValue(plan.RetryAsBatch), StaticHeaders: headers, TargetVersion: stringPointerValue(plan.TargetVersion),
		Timeout: stringPointerValue(plan.Timeout), Topic: plan.Topic.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Unable to Update Kafka Outlet", err.Error())
		return
	}
	fieldUpdater := kafkaOutletFieldUpdater{client: r.client}
	mutations := []struct {
		name string
		err  error
	}{
		{"ack_replicas", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.ack_replicas", plan.AckReplicas, state.AckReplicas)},
		{"backoff_initial_interval", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.backoff.initial_interval", plan.BackoffInitialInterval, state.BackoffInitialInterval)},
		{"backoff_max_elapsed_time", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.backoff.max_elapsed_time", plan.BackoffMaxElapsedTime, state.BackoffMaxElapsedTime)},
		{"backoff_max_interval", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.backoff.max_interval", plan.BackoffMaxInterval, state.BackoffMaxInterval)},
		{"batching_byte_size", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.batching.byte_size", plan.BatchingByteSize, state.BatchingByteSize)},
		{"batching_check", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.batching.check", plan.BatchingCheck, state.BatchingCheck)},
		{"batching_count", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.batching.count", plan.BatchingCount, state.BatchingCount)},
		{"batching_jitter", syncFloat64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.batching.jitter", plan.BatchingJitter, state.BatchingJitter)},
		{"batching_period", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.batching.period", plan.BatchingPeriod, state.BatchingPeriod)},
		{"client_id", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.client_id", plan.ClientID, state.ClientID)},
		{"compression", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.compression", plan.Compression, state.Compression)},
		{"custom_topic_creation_enabled", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.custom_topic_creation.enabled", plan.CustomTopicCreationEnabled, state.CustomTopicCreationEnabled)},
		{"custom_topic_creation_partitions", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.custom_topic_creation.partitions", plan.CustomTopicCreationPartitions, state.CustomTopicCreationPartitions)},
		{"custom_topic_creation_replication_factor", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.custom_topic_creation.replication_factor", plan.CustomTopicCreationReplicationFactor, state.CustomTopicCreationReplicationFactor)},
		{"idempotent_write", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.idempotent_write", plan.IdempotentWrite, state.IdempotentWrite)},
		{"inject_tracing_map", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.inject_tracing_map", plan.InjectTracingMap, state.InjectTracingMap)},
		{"key", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.key", plan.Key, state.Key)},
		{"max_in_flight", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.max_in_flight", plan.MaxInFlight, state.MaxInFlight)},
		{"max_msg_bytes", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.max_msg_bytes", plan.MaxMsgBytes, state.MaxMsgBytes)},
		{"max_retries", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.max_retries", plan.MaxRetries, state.MaxRetries)},
		{"metadata_exclude_prefixes", syncStringListField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.metadata.exclude_prefixes", plan.MetadataExcludePrefixes, state.MetadataExcludePrefixes)},
		{"notify_policy", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "notify_policy", plan.NotifyPolicy, state.NotifyPolicy)},
		{"partition", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.partition", plan.Partition, state.Partition)},
		{"partitioner", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.partitioner", plan.Partitioner, state.Partitioner)},
		{"rack_id", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.rack_id", plan.RackID, state.RackID)},
		{"retry_as_batch", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.retry_as_batch", plan.RetryAsBatch, state.RetryAsBatch)},
		{"static_headers", syncStringMapField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.static_headers", plan.StaticHeaders, state.StaticHeaders)},
		{"target_version", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.target_version", plan.TargetVersion, state.TargetVersion)},
		{"timeout", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "kafka.timeout", plan.Timeout, state.Timeout)},
	}
	for _, mutation := range mutations {
		if mutation.err != nil {
			resp.Diagnostics.AddError("Unable to Update Kafka Field", fmt.Sprintf("%s: %s", mutation.name, mutation.err))
			return
		}
	}
	if err := r.applyFeatureConfig(ctx, state.ID.ValueInt64(), kafkaFeatureState{
		tls: state.TLS, plain: state.SASLPlain, scram256: state.SASLScram256, scram512: state.SASLScram512,
		oauth: state.SASLOAuthbearer, oauthCache: state.SASLOAuthbearerCache, awsMSKIAM: state.SASLAwsMSKIAM,
	}, cfg); err != nil {
		resp.Diagnostics.AddError("Unable to Update Kafka Features", err.Error())
		return
	}
	if err := syncTypedOutletTransform(ctx, r.client, "kafka", state.ID.ValueInt64(), configuredTransform, selectedTransform(state.Transform)); err != nil {
		resp.Diagnostics.AddError("Unable to Update Kafka Transform", err.Error())
		return
	}
	newState, removed, diags := r.buildState(ctx, state.ID.ValueInt64(), plan.TLS, plan.SASLPlain, plan.SASLScram256, plan.SASLScram512, plan.SASLOAuthbearer, plan.SASLOAuthbearerCache, plan.SASLAwsMSKIAM)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if removed {
		result.Removed = true
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "kafka", state.ID.ValueInt64(), plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe Kafka Outlet Lifecycle", err.Error())
		return
	}
	result.State = newState
	return
}

func (a kafkaOutletAdapter) Delete(ctx context.Context, state kafkaOutletResourceModel) diag.Diagnostics {
	operation, err := a.resource.client.DeleteKafkaOutletOperation(ctx, state.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Unable to Delete Kafka Outlet", err.Error())}
	}
	if err == nil && operation.DesiredGeneration != 0 {
		if err := observeTypedOutletLifecycle(ctx, a.resource.client, "kafka", state.ID.ValueInt64(), operation.DesiredGeneration, false, true); err != nil {
			return diag.Diagnostics{diag.NewErrorDiagnostic("Unable to Observe Kafka Outlet Lifecycle", err.Error())}
		}
	}
	return nil
}

type kafkaConfigBlocks struct {
	tls        kafkaTLSModel
	plain      kafkaSASLUserPassModel
	scram256   kafkaSASLUserPassModel
	scram512   kafkaSASLUserPassModel
	oauth      kafkaSASLOAuthbearerModel
	oauthCache kafkaSASLOAuthbearerCacheModel
	awsMSKIAM  kafkaSASLAwsMSKIAMModel
}

type kafkaFeatureState struct {
	tls        types.Object
	plain      types.Object
	scram256   types.Object
	scram512   types.Object
	oauth      types.Object
	oauthCache types.Object
	awsMSKIAM  types.Object
}

func (r *kafkaOutletResource) readConfig(ctx context.Context, cfg tfsdk.Config, plan kafkaOutletResourceModel, diags *diag.Diagnostics) (kafkaConfigBlocks, bool) {
	blocks := kafkaConfigBlocks{}
	var d diag.Diagnostics
	blocks.tls, d = expandKafkaTLS(ctx, plan.TLS)
	diags.Append(d...)
	blocks.plain, d = expandKafkaSASLUserPass(ctx, plan.SASLPlain)
	diags.Append(d...)
	blocks.scram256, d = expandKafkaSASLUserPass(ctx, plan.SASLScram256)
	diags.Append(d...)
	blocks.scram512, d = expandKafkaSASLUserPass(ctx, plan.SASLScram512)
	diags.Append(d...)
	blocks.oauth, d = expandKafkaSASLOAuthbearer(ctx, plan.SASLOAuthbearer)
	diags.Append(d...)
	blocks.oauthCache, d = expandKafkaSASLOAuthbearerCache(ctx, plan.SASLOAuthbearerCache)
	diags.Append(d...)
	blocks.awsMSKIAM, d = expandKafkaSASLAwsMSKIAM(ctx, plan.SASLAwsMSKIAM)
	diags.Append(d...)
	if diags.HasError() {
		return blocks, false
	}
	blocks.tls, d = enrichKafkaTLSFromConfig(ctx, cfg, "tls", blocks.tls)
	diags.Append(d...)
	blocks.plain, d = enrichKafkaSASLUserPassFromConfig(ctx, cfg, "sasl_plain", blocks.plain)
	diags.Append(d...)
	blocks.scram256, d = enrichKafkaSASLUserPassFromConfig(ctx, cfg, "sasl_scram256", blocks.scram256)
	diags.Append(d...)
	blocks.scram512, d = enrichKafkaSASLUserPassFromConfig(ctx, cfg, "sasl_scram512", blocks.scram512)
	diags.Append(d...)
	blocks.oauth, d = enrichKafkaSASLOAuthbearerFromConfig(ctx, cfg, blocks.oauth)
	diags.Append(d...)
	blocks.awsMSKIAM, d = enrichKafkaSASLAwsMSKIAMFromConfig(ctx, cfg, blocks.awsMSKIAM)
	diags.Append(d...)
	if diags.HasError() {
		return blocks, false
	}
	return blocks, true
}

func (r *kafkaOutletResource) expandRootValues(ctx context.Context, plan kafkaOutletResourceModel, diags *diag.Diagnostics) ([]string, map[string]string, []string, bool) {
	addresses, err := stringListValueOrNil(ctx, plan.Addresses)
	if err != nil {
		diags.AddError("Invalid Kafka Addresses", err.Error())
		return nil, nil, nil, false
	}
	headers, err := stringMapValueOrNil(ctx, plan.StaticHeaders)
	if err != nil {
		diags.AddError("Invalid Kafka Static Headers", err.Error())
		return nil, nil, nil, false
	}
	metadataExcludePrefixes, err := stringListValueOrNil(ctx, plan.MetadataExcludePrefixes)
	if err != nil {
		diags.AddError("Invalid Kafka Metadata Exclude Prefixes", err.Error())
		return nil, nil, nil, false
	}
	return addresses, headers, metadataExcludePrefixes, true
}

func (r *kafkaOutletResource) buildState(ctx context.Context, outletID int64, previousTLS, previousPlain, previousScram256, previousScram512, previousOAuth, previousOAuthCache, previousAwsMSKIAM types.Object) (kafkaOutletResourceModel, bool, diag.Diagnostics) {
	outlet, err := r.client.GetKafkaOutlet(ctx, outletID)
	if err != nil {
		if client.IsNotFound(err) {
			return kafkaOutletResourceModel{}, true, nil
		}
		var diags diag.Diagnostics
		diags.AddError("Unable to Read Kafka Outlet", err.Error())
		return kafkaOutletResourceModel{}, false, diags
	}
	credentials, err := r.client.ListKafkaCredentials(ctx, outletID)
	if err != nil && !client.IsNotFound(err) {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read Kafka Credentials", err.Error())
		return kafkaOutletResourceModel{}, false, diags
	}
	tlsFeature, plainFeature, scram256Feature, scram512Feature, oauthFeature, oauthCacheFeature, awsFeature, diags := r.readFeatures(ctx, outletID)
	if diags.HasError() {
		return kafkaOutletResourceModel{}, false, diags
	}
	tlsObject, d := buildKafkaTLSObject(ctx, previousTLS, tlsFeature, credentials)
	diags.Append(d...)
	plainObject, d := buildKafkaSASLUserPassObject(ctx, previousPlain, plainFeature)
	diags.Append(d...)
	scram256Object, d := buildKafkaSASLUserPassObject(ctx, previousScram256, scram256Feature)
	diags.Append(d...)
	scram512Object, d := buildKafkaSASLUserPassObject(ctx, previousScram512, scram512Feature)
	diags.Append(d...)
	oauthObject, d := buildKafkaSASLOAuthbearerObject(ctx, previousOAuth, oauthFeature)
	diags.Append(d...)
	oauthCacheObject, d := buildKafkaSASLOAuthbearerCacheObject(ctx, previousOAuthCache, oauthCacheFeature)
	diags.Append(d...)
	awsObject, d := buildKafkaSASLAwsMSKIAMObject(ctx, previousAwsMSKIAM, awsFeature)
	diags.Append(d...)
	if diags.HasError() {
		return kafkaOutletResourceModel{}, false, diags
	}
	transform, err := r.client.GetTypedOutletTransform(ctx, "kafka", outletID)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read Kafka Transform", err.Error())
		return kafkaOutletResourceModel{}, false, diags
	}
	return kafkaOutletModelFromAPI(outlet, tlsObject, plainObject, scram256Object, scram512Object, oauthObject, oauthCacheObject, awsObject, transform), false, nil
}

func (r *kafkaOutletResource) readFeatures(ctx context.Context, outletID int64) (*client.OutletFeatureDetail, *client.OutletFeatureDetail, *client.OutletFeatureDetail, *client.OutletFeatureDetail, *client.OutletFeatureDetail, *client.OutletFeatureDetail, *client.OutletFeatureDetail, diag.Diagnostics) {
	var diags diag.Diagnostics
	tls := readKafkaFeature(ctx, r.client.GetKafkaTLSFeature, outletID, "Kafka TLS Feature", &diags)
	plain := readKafkaFeature(ctx, r.client.GetKafkaSASLPlainFeature, outletID, "Kafka SASL/PLAIN Feature", &diags)
	scram256 := readKafkaFeature(ctx, r.client.GetKafkaSASLScram256Feature, outletID, "Kafka SASL SCRAM-256 Feature", &diags)
	scram512 := readKafkaFeature(ctx, r.client.GetKafkaSASLScram512Feature, outletID, "Kafka SASL SCRAM-512 Feature", &diags)
	oauth := readKafkaFeature(ctx, r.client.GetKafkaSASLOAuthbearerFeature, outletID, "Kafka SASL OAUTHBEARER Feature", &diags)
	oauthCache := readKafkaFeature(ctx, r.client.GetKafkaSASLOAuthbearerCacheFeature, outletID, "Kafka SASL OAUTHBEARER Cache Feature", &diags)
	awsMSKIAM := readKafkaFeature(ctx, r.client.GetKafkaSASLAwsMSKIAMFeature, outletID, "Kafka SASL AWS_MSK_IAM Feature", &diags)
	return tls, plain, scram256, scram512, oauth, oauthCache, awsMSKIAM, diags
}

func (r *kafkaOutletResource) applyFeatureConfig(ctx context.Context, outletID int64, state kafkaFeatureState, plan kafkaConfigBlocks) error {
	if !state.tls.IsNull() && !amqpFeatureEnabled(state.tls) {
		state.tls = types.ObjectNull(kafkaTLSAttributeTypes)
	}
	if !state.plain.IsNull() && !amqpFeatureEnabled(state.plain) {
		state.plain = types.ObjectNull(kafkaSASLUserPassAttributeTypes)
	}
	if !state.scram256.IsNull() && !amqpFeatureEnabled(state.scram256) {
		state.scram256 = types.ObjectNull(kafkaSASLUserPassAttributeTypes)
	}
	if !state.scram512.IsNull() && !amqpFeatureEnabled(state.scram512) {
		state.scram512 = types.ObjectNull(kafkaSASLUserPassAttributeTypes)
	}
	if !state.oauth.IsNull() && !amqpFeatureEnabled(state.oauth) {
		state.oauth = types.ObjectNull(kafkaSASLOAuthbearerAttributeTypes)
	}
	if !state.oauthCache.IsNull() && !amqpFeatureEnabled(state.oauthCache) {
		state.oauthCache = types.ObjectNull(kafkaSASLOAuthbearerCacheAttributeTypes)
	}
	if !state.awsMSKIAM.IsNull() && !amqpFeatureEnabled(state.awsMSKIAM) {
		state.awsMSKIAM = types.ObjectNull(kafkaSASLAwsMSKIAMAttributeTypes)
	}

	activePlan := r.activeSASLPlan(plan)
	if activePlan != "plain" {
		if err := r.deleteKafkaPlain(ctx, outletID, state.plain); err != nil {
			return err
		}
	}
	if activePlan != "scram256" {
		if err := r.deleteKafkaScram256(ctx, outletID, state.scram256); err != nil {
			return err
		}
	}
	if activePlan != "scram512" {
		if err := r.deleteKafkaScram512(ctx, outletID, state.scram512); err != nil {
			return err
		}
	}
	if activePlan != "oauth" {
		if err := r.deleteKafkaOAuth(ctx, outletID, state.oauth); err != nil {
			return err
		}
	}
	if activePlan != "oauth_cache" {
		if err := r.deleteKafkaOAuthCache(ctx, outletID, state.oauthCache); err != nil {
			return err
		}
	}
	if activePlan != "aws_msk_iam" {
		if err := r.deleteKafkaAwsMSKIAM(ctx, outletID, state.awsMSKIAM); err != nil {
			return err
		}
	}

	if err := r.applyKafkaTLS(ctx, outletID, state.tls, plan.tls); err != nil {
		return err
	}
	switch activePlan {
	case "plain":
		return r.applyKafkaUserPass(ctx, outletID, state.plain, plan.plain, r.client.ConfigureKafkaSASLPlainFeature, r.client.UpdateKafkaSASLPlainFeature)
	case "scram256":
		return r.applyKafkaUserPass(ctx, outletID, state.scram256, plan.scram256, r.client.ConfigureKafkaSASLScram256Feature, r.client.UpdateKafkaSASLScram256Feature)
	case "scram512":
		return r.applyKafkaUserPass(ctx, outletID, state.scram512, plan.scram512, r.client.ConfigureKafkaSASLScram512Feature, r.client.UpdateKafkaSASLScram512Feature)
	case "oauth":
		return r.applyKafkaOAuth(ctx, outletID, state.oauth, plan.oauth)
	case "oauth_cache":
		return r.applyKafkaOAuthCache(ctx, outletID, state.oauthCache, plan.oauthCache)
	case "aws_msk_iam":
		return r.applyKafkaAwsMSKIAM(ctx, outletID, state.awsMSKIAM, plan.awsMSKIAM)
	default:
		return nil
	}
}

func (r *kafkaOutletResource) activeSASLPlan(plan kafkaConfigBlocks) string {
	switch {
	case !plan.plain.PasswordVersion.IsNull():
		return "plain"
	case !plan.scram256.PasswordVersion.IsNull():
		return "scram256"
	case !plan.scram512.PasswordVersion.IsNull():
		return "scram512"
	case !plan.oauth.AccessTokenVersion.IsNull():
		return "oauth"
	case !plan.oauthCache.TokenCache.IsNull():
		return "oauth_cache"
	case !plan.awsMSKIAM.Region.IsNull():
		return "aws_msk_iam"
	default:
		return ""
	}
}

func (r *kafkaOutletResource) applyKafkaTLS(ctx context.Context, outletID int64, state types.Object, plan kafkaTLSModel) error {
	if plan.CACertVersion.IsNull() {
		if !state.IsNull() && !state.IsUnknown() {
			if err := r.client.DeleteKafkaTLSFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
				return err
			}
			for _, role := range []string{kafkaCACertRole, kafkaClientCertRole, kafkaClientKeyRole} {
				if _, err := r.client.DeleteKafkaCredential(ctx, outletID, role); err != nil && !client.IsNotFound(err) {
					return err
				}
			}
		}
		return nil
	}
	req := client.KafkaTLSFeatureRequest{
		CACertFilename:      plan.CACertFilename.ValueString(),
		CACertContentBase64: plan.CACertContentBase64.ValueString(),
	}
	if !plan.ClientCertVersion.IsNull() {
		req.ClientCertFilename = stringValueOrEmpty(plan.ClientCertFilename)
		req.ClientCertContentBase64 = stringValueOrEmpty(plan.ClientCertContentBase64)
	}
	if !plan.ClientKeyVersion.IsNull() {
		req.ClientKeyFilename = stringValueOrEmpty(plan.ClientKeyFilename)
		req.ClientKeyContentBase64 = stringValueOrEmpty(plan.ClientKeyContentBase64)
	}
	if state.IsNull() || state.IsUnknown() {
		return r.client.ConfigureKafkaTLSFeature(ctx, outletID, req)
	}
	stateModel, diags := expandKafkaTLS(ctx, state)
	if diags.HasError() {
		return errors.New(diags.Errors()[0].Detail())
	}
	changed := !plan.CACertVersion.Equal(stateModel.CACertVersion) || !plan.CACertFilename.Equal(stateModel.CACertFilename) ||
		!plan.ClientCertVersion.Equal(stateModel.ClientCertVersion) || !plan.ClientCertFilename.Equal(stateModel.ClientCertFilename) ||
		!plan.ClientKeyVersion.Equal(stateModel.ClientKeyVersion) || !plan.ClientKeyFilename.Equal(stateModel.ClientKeyFilename)
	if changed {
		if err := r.client.UpdateKafkaTLSFeature(ctx, outletID, req); err != nil {
			return err
		}
	}
	if plan.ClientCertVersion.IsNull() {
		if _, err := r.client.DeleteKafkaCredential(ctx, outletID, kafkaClientCertRole); err != nil && !client.IsNotFound(err) {
			return err
		}
	}
	if plan.ClientKeyVersion.IsNull() {
		if _, err := r.client.DeleteKafkaCredential(ctx, outletID, kafkaClientKeyRole); err != nil && !client.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (r *kafkaOutletResource) applyKafkaUserPass(ctx context.Context, outletID int64, state types.Object, plan kafkaSASLUserPassModel, createFn func(context.Context, int64, client.KafkaSASLUserPassFeatureRequest) error, updateFn func(context.Context, int64, client.KafkaSASLUserPassFeatureRequest) error) error {
	req := client.KafkaSASLUserPassFeatureRequest{User: plan.User.ValueString(), Password: plan.Password.ValueString()}
	if state.IsNull() || state.IsUnknown() {
		return createFn(ctx, outletID, req)
	}
	stateModel, diags := expandKafkaSASLUserPass(ctx, state)
	if diags.HasError() {
		return errors.New(diags.Errors()[0].Detail())
	}
	if !plan.User.Equal(stateModel.User) || !plan.PasswordVersion.Equal(stateModel.PasswordVersion) {
		return updateFn(ctx, outletID, req)
	}
	return nil
}

func (r *kafkaOutletResource) applyKafkaOAuth(ctx context.Context, outletID int64, state types.Object, plan kafkaSASLOAuthbearerModel) error {
	req := client.KafkaSASLOAuthbearerFeatureRequest{AccessToken: plan.AccessToken.ValueString()}
	if state.IsNull() || state.IsUnknown() {
		return r.client.ConfigureKafkaSASLOAuthbearerFeature(ctx, outletID, req)
	}
	stateModel, diags := expandKafkaSASLOAuthbearer(ctx, state)
	if diags.HasError() {
		return errors.New(diags.Errors()[0].Detail())
	}
	if !plan.AccessTokenVersion.Equal(stateModel.AccessTokenVersion) {
		return r.client.UpdateKafkaSASLOAuthbearerFeature(ctx, outletID, req)
	}
	return nil
}

func (r *kafkaOutletResource) applyKafkaOAuthCache(ctx context.Context, outletID int64, state types.Object, plan kafkaSASLOAuthbearerCacheModel) error {
	req := client.KafkaSASLOAuthbearerCacheFeatureRequest{TokenCache: plan.TokenCache.ValueString(), TokenKey: plan.TokenKey.ValueString()}
	if state.IsNull() || state.IsUnknown() {
		return r.client.ConfigureKafkaSASLOAuthbearerCacheFeature(ctx, outletID, req)
	}
	stateModel, diags := expandKafkaSASLOAuthbearerCache(ctx, state)
	if diags.HasError() {
		return errors.New(diags.Errors()[0].Detail())
	}
	if !plan.TokenCache.Equal(stateModel.TokenCache) || !plan.TokenKey.Equal(stateModel.TokenKey) {
		return r.client.UpdateKafkaSASLOAuthbearerCacheFeature(ctx, outletID, req)
	}
	return nil
}

func (r *kafkaOutletResource) applyKafkaAwsMSKIAM(ctx context.Context, outletID int64, state types.Object, plan kafkaSASLAwsMSKIAMModel) error {
	req := client.KafkaSASLAwsMSKIAMFeatureRequest{
		Region: plan.Region.ValueString(), Endpoint: stringPointerValue(plan.Endpoint), Profile: stringPointerValue(plan.Profile), ID: stringPointerValue(plan.ID),
		FromEC2Role: stringPointerValue(plan.FromEC2Role), Role: stringPointerValue(plan.Role), RoleExternalID: stringPointerValue(plan.RoleExternalID), ExpiryWindow: stringPointerValue(plan.ExpiryWindow),
		AwsSecret: stringPointerValue(plan.AwsSecret), AwsToken: stringPointerValue(plan.AwsToken),
	}
	if state.IsNull() || state.IsUnknown() {
		return r.client.ConfigureKafkaSASLAwsMSKIAMFeature(ctx, outletID, req)
	}
	stateModel, diags := expandKafkaSASLAwsMSKIAM(ctx, state)
	if diags.HasError() {
		return errors.New(diags.Errors()[0].Detail())
	}
	changed := !plan.Region.Equal(stateModel.Region) || !plan.Endpoint.Equal(stateModel.Endpoint) || !plan.Profile.Equal(stateModel.Profile) ||
		!plan.ID.Equal(stateModel.ID) || !plan.FromEC2Role.Equal(stateModel.FromEC2Role) || !plan.Role.Equal(stateModel.Role) ||
		!plan.RoleExternalID.Equal(stateModel.RoleExternalID) || !plan.ExpiryWindow.Equal(stateModel.ExpiryWindow) ||
		!plan.AwsSecretVersion.Equal(stateModel.AwsSecretVersion) || !plan.AwsTokenVersion.Equal(stateModel.AwsTokenVersion)
	if changed {
		if err := r.client.UpdateKafkaSASLAwsMSKIAMFeature(ctx, outletID, req); err != nil {
			return err
		}
	}
	if plan.AwsSecretVersion.IsNull() {
		if _, err := r.client.DeleteKafkaCredential(ctx, outletID, kafkaAwsSecretRole); err != nil && !client.IsNotFound(err) {
			return err
		}
	}
	if plan.AwsTokenVersion.IsNull() {
		if _, err := r.client.DeleteKafkaCredential(ctx, outletID, kafkaAwsTokenRole); err != nil && !client.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (r *kafkaOutletResource) deleteKafkaPlain(ctx context.Context, outletID int64, state types.Object) error {
	return r.deleteKafkaFeature(ctx, outletID, state, r.client.DeleteKafkaSASLPlainFeature, []string{kafkaSASLPasswordRole})
}

func (r *kafkaOutletResource) deleteKafkaScram256(ctx context.Context, outletID int64, state types.Object) error {
	return r.deleteKafkaFeature(ctx, outletID, state, r.client.DeleteKafkaSASLScram256Feature, []string{kafkaSASLPasswordRole})
}

func (r *kafkaOutletResource) deleteKafkaScram512(ctx context.Context, outletID int64, state types.Object) error {
	return r.deleteKafkaFeature(ctx, outletID, state, r.client.DeleteKafkaSASLScram512Feature, []string{kafkaSASLPasswordRole})
}

func (r *kafkaOutletResource) deleteKafkaOAuth(ctx context.Context, outletID int64, state types.Object) error {
	return r.deleteKafkaFeature(ctx, outletID, state, r.client.DeleteKafkaSASLOAuthbearerFeature, []string{kafkaSASLAccessTokenRole})
}

func (r *kafkaOutletResource) deleteKafkaOAuthCache(ctx context.Context, outletID int64, state types.Object) error {
	return r.deleteKafkaFeature(ctx, outletID, state, r.client.DeleteKafkaSASLOAuthbearerCacheFeature, nil)
}

func (r *kafkaOutletResource) deleteKafkaAwsMSKIAM(ctx context.Context, outletID int64, state types.Object) error {
	return r.deleteKafkaFeature(ctx, outletID, state, r.client.DeleteKafkaSASLAwsMSKIAMFeature, []string{kafkaAwsSecretRole, kafkaAwsTokenRole})
}

func (r *kafkaOutletResource) deleteKafkaFeature(ctx context.Context, outletID int64, state types.Object, deleteFn func(context.Context, int64) error, credentials []string) error {
	if state.IsNull() || state.IsUnknown() {
		return nil
	}
	if err := deleteFn(ctx, outletID); err != nil && !client.IsNotFound(err) {
		return err
	}
	for _, role := range credentials {
		if _, err := r.client.DeleteKafkaCredential(ctx, outletID, role); err != nil && !client.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func validateKafkaSASLUserPass(ctx context.Context, object types.Object, label string, diags *diag.Diagnostics) {
	model, d := expandKafkaSASLUserPass(ctx, object)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	if model.User.IsNull() || model.User.IsUnknown() || model.Password.IsNull() || model.Password.IsUnknown() || model.PasswordVersion.IsNull() || model.PasswordVersion.IsUnknown() {
		diags.AddError("Incomplete "+label+" Configuration", fmt.Sprintf("When `%s` is configured, `user`, `password`, and `password_version` must all be set.", label))
	}
}

func anyKafkaTLSClientMaterialConfigured(model kafkaTLSModel) bool {
	return !model.ClientCertFilename.IsNull() || !model.ClientCertContentBase64.IsNull() || !model.ClientCertVersion.IsNull() || !model.ClientKeyFilename.IsNull() || !model.ClientKeyContentBase64.IsNull() || !model.ClientKeyVersion.IsNull()
}

func fullKafkaTLSClientMaterialConfigured(model kafkaTLSModel) bool {
	return !model.ClientCertFilename.IsNull() && !model.ClientCertContentBase64.IsNull() && !model.ClientCertVersion.IsNull() &&
		!model.ClientKeyFilename.IsNull() && !model.ClientKeyContentBase64.IsNull() && !model.ClientKeyVersion.IsNull() &&
		!model.ClientCertFilename.IsUnknown() && !model.ClientCertContentBase64.IsUnknown() && !model.ClientCertVersion.IsUnknown() &&
		!model.ClientKeyFilename.IsUnknown() && !model.ClientKeyContentBase64.IsUnknown() && !model.ClientKeyVersion.IsUnknown()
}

func readKafkaFeature(ctx context.Context, getter func(context.Context, int64) (client.OutletFeatureDetail, error), outletID int64, label string, diags *diag.Diagnostics) *client.OutletFeatureDetail {
	feature, err := getter(ctx, outletID)
	if err != nil {
		if !client.IsNotFound(err) {
			diags.AddError("Unable to Read "+label, err.Error())
		}
		return nil
	}
	if feature.Feature == "" {
		return nil
	}
	return &feature
}

func kafkaOutletModelFromAPI(outlet client.KafkaOutlet, tls, plain, scram256, scram512, oauth, oauthCache, aws types.Object, transform client.OutletTransformState) kafkaOutletResourceModel {
	return kafkaOutletResourceModel{
		ID: types.Int64Value(outlet.ID), Name: types.StringValue(outlet.Name), Description: nullableString(outlet.Description), Enabled: types.BoolValue(outlet.Enabled),
		AckReplicas: nullableBoolPointer(outlet.AckReplicas), Addresses: nullableStringList(outlet.Addresses), BackoffInitialInterval: nullableStringPointer(outlet.BackoffInitialInterval),
		BackoffMaxElapsedTime: nullableStringPointer(outlet.BackoffMaxElapsedTime), BackoffMaxInterval: nullableStringPointer(outlet.BackoffMaxInterval),
		BatchingByteSize: nullableInt64Pointer(outlet.BatchingByteSize), BatchingCheck: nullableStringPointer(outlet.BatchingCheck), BatchingCount: nullableInt64Pointer(outlet.BatchingCount),
		BatchingJitter: nullableFloat64Pointer(outlet.BatchingJitter), BatchingPeriod: nullableStringPointer(outlet.BatchingPeriod), ClientID: nullableStringPointer(outlet.ClientID),
		Compression: nullableStringPointer(outlet.Compression), CustomTopicCreationEnabled: nullableBoolPointer(outlet.CustomTopicCreationEnabled),
		CustomTopicCreationPartitions: nullableInt64Pointer(outlet.CustomTopicCreationPartitions), CustomTopicCreationReplicationFactor: nullableInt64Pointer(outlet.CustomTopicCreationReplicationFactor),
		IdempotentWrite: nullableBoolPointer(outlet.IdempotentWrite), InjectTracingMap: nullableStringPointer(outlet.InjectTracingMap), Key: nullableStringPointer(outlet.Key),
		MaxInFlight: nullableInt64Pointer(outlet.MaxInFlight), MaxMsgBytes: nullableInt64Pointer(outlet.MaxMsgBytes), MaxRetries: nullableInt64Pointer(outlet.MaxRetries),
		MetadataExcludePrefixes: nullableStringList(outlet.MetadataExcludePrefixes), NotifyPolicy: nullableStringPointerPreservingEmpty(outlet.NotifyPolicy), Partition: nullableStringPointer(outlet.Partition), Partitioner: nullableStringPointer(outlet.Partitioner),
		RackID: nullableStringPointer(outlet.RackID), RetryAsBatch: nullableBoolPointer(outlet.RetryAsBatch), StaticHeaders: nullableStringMap(outlet.StaticHeaders), TargetVersion: nullableStringPointer(outlet.TargetVersion),
		Timeout: nullableStringPointer(outlet.Timeout), Topic: types.StringValue(outlet.Topic), TLS: tls, SASLPlain: plain, SASLScram256: scram256, SASLScram512: scram512,
		SASLOAuthbearer: oauth, SASLOAuthbearerCache: oauthCache, SASLAwsMSKIAM: aws, ApplyStatus: nullableString(outlet.ApplyStatus), LastApplyError: nullableString(outlet.LastApplyError),
		ConfigGeneration: types.Int64Value(outlet.ConfigGeneration), AppliedGeneration: types.Int64Value(outlet.AppliedGeneration),
		Transform: transformObject(transform),
	}
}

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

var (
	_ resource.Resource                = &awsSqsOutletResource{}
	_ resource.ResourceWithConfigure   = &awsSqsOutletResource{}
	_ resource.ResourceWithImportState = &awsSqsOutletResource{}
)

type awsSqsOutletClient interface {
	CreateAwsSqsOutlet(context.Context, client.AwsSqsOutletCreateRequest) (client.AwsSqsOutlet, error)
	GetAwsSqsOutlet(context.Context, int64) (client.AwsSqsOutlet, error)
	UpdateAwsSqsOutlet(context.Context, int64, client.AwsSqsOutletUpdateRequest) error
	DeleteAwsSqsOutlet(context.Context, int64) error
	DeleteAwsSqsOutletOperation(context.Context, int64) (client.OutletOperation, error)
	GetTypedOutletStatus(context.Context, string, int64, int64) (client.OutletStatus, error)
	ListAwsSqsCredentials(context.Context, int64) ([]client.OutletMaterialSummary, error)
	DeleteAwsSqsMaterial(context.Context, int64, string) (client.OutletMaterialDeleteResponse, error)
	SetAwsSqsField(context.Context, int64, string, any) error
	DeleteAwsSqsField(context.Context, int64, string) error
}

type awsSqsOutletFieldUpdater struct {
	client awsSqsOutletClient
}

func (u awsSqsOutletFieldUpdater) SetField(ctx context.Context, outletID int64, field string, value any) error {
	return u.client.SetAwsSqsField(ctx, outletID, field, value)
}

func (u awsSqsOutletFieldUpdater) DeleteField(ctx context.Context, outletID int64, field string) error {
	return u.client.DeleteAwsSqsField(ctx, outletID, field)
}

type awsSqsOutletResource struct {
	client awsSqsOutletClient
}

type awsSqsOutletResourceModel struct {
	ID                      types.Int64   `tfsdk:"id"`
	Name                    types.String  `tfsdk:"name"`
	Description             types.String  `tfsdk:"description"`
	Enabled                 types.Bool    `tfsdk:"enabled"`
	BackoffInitialInterval  types.String  `tfsdk:"backoff_initial_interval"`
	BackoffMaxElapsedTime   types.String  `tfsdk:"backoff_max_elapsed_time"`
	BackoffMaxInterval      types.String  `tfsdk:"backoff_max_interval"`
	BatchingByteSize        types.Int64   `tfsdk:"batching_byte_size"`
	BatchingCheck           types.String  `tfsdk:"batching_check"`
	BatchingCount           types.Int64   `tfsdk:"batching_count"`
	BatchingJitter          types.Float64 `tfsdk:"batching_jitter"`
	BatchingPeriod          types.String  `tfsdk:"batching_period"`
	DelaySeconds            types.Int64   `tfsdk:"delay_seconds"`
	EndpointURL             types.String  `tfsdk:"endpoint_url"`
	ExpiryWindow            types.String  `tfsdk:"expiry_window"`
	ExternalID              types.String  `tfsdk:"external_id"`
	MaxInFlight             types.Int64   `tfsdk:"max_in_flight"`
	MaxRetries              types.Int64   `tfsdk:"max_retries"`
	MessageDeduplicationID  types.String  `tfsdk:"message_deduplication_id"`
	MessageGroupID          types.String  `tfsdk:"message_group_id"`
	MetadataExcludePrefixes types.List    `tfsdk:"metadata_exclude_prefixes"`
	NotifyPolicy            types.String  `tfsdk:"notify_policy"`
	Profile                 types.String  `tfsdk:"profile"`
	QueueURL                types.String  `tfsdk:"queue_url"`
	Region                  types.String  `tfsdk:"region"`
	RoleARN                 types.String  `tfsdk:"role_arn"`
	UseEC2InstanceRole      types.Bool    `tfsdk:"use_ec2_instance_role"`
	Credentials             types.Object  `tfsdk:"credentials"`
	ApplyStatus             types.String  `tfsdk:"apply_status"`
	LastApplyError          types.String  `tfsdk:"last_apply_error"`
	ConfigGeneration        types.Int64   `tfsdk:"config_generation"`
	AppliedGeneration       types.Int64   `tfsdk:"applied_generation"`
}

func NewAwsSqsOutletResource() resource.Resource {
	return &awsSqsOutletResource{}
}

func (r *awsSqsOutletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_sqs_outlet"
}

func (r *awsSqsOutletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		MarkdownDescription: "Manages a Noozle AWS SQS outlet.",
		Attributes: map[string]resourceschema.Attribute{
			"id": resourceschema.Int64Attribute{
				MarkdownDescription: "Unique identifier of the outlet.",
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name":                      resourceschema.StringAttribute{MarkdownDescription: "Name presented for the outlet.", Required: true},
			"description":               resourceschema.StringAttribute{MarkdownDescription: "Optional note explaining the outlet purpose.", Optional: true},
			"enabled":                   resourceschema.BoolAttribute{MarkdownDescription: "Whether the outlet is active.", Required: true},
			"backoff_initial_interval":  resourceschema.StringAttribute{MarkdownDescription: "Initial SQS retry backoff interval.", Optional: true},
			"backoff_max_elapsed_time":  resourceschema.StringAttribute{MarkdownDescription: "Maximum total SQS retry backoff time.", Optional: true},
			"backoff_max_interval":      resourceschema.StringAttribute{MarkdownDescription: "Maximum SQS retry backoff interval.", Optional: true},
			"batching_byte_size":        resourceschema.Int64Attribute{MarkdownDescription: "Maximum bytes per SQS batch.", Optional: true},
			"batching_check":            resourceschema.StringAttribute{MarkdownDescription: "Bloblang condition for flushing an SQS batch.", Optional: true},
			"batching_count":            resourceschema.Int64Attribute{MarkdownDescription: "Maximum messages per SQS batch.", Optional: true},
			"batching_jitter":           resourceschema.Float64Attribute{MarkdownDescription: "Randomized jitter to apply to the SQS batching period.", Optional: true},
			"batching_period":           resourceschema.StringAttribute{MarkdownDescription: "Maximum wait period before flushing an SQS batch.", Optional: true},
			"delay_seconds":             resourceschema.Int64Attribute{MarkdownDescription: "Delay in seconds before an SQS message becomes visible.", Optional: true},
			"endpoint_url":              resourceschema.StringAttribute{MarkdownDescription: "AWS service endpoint URL for SQS.", Optional: true},
			"expiry_window":             resourceschema.StringAttribute{MarkdownDescription: "Expiry window for refreshed SQS credentials.", Optional: true},
			"external_id":               resourceschema.StringAttribute{MarkdownDescription: "External ID for assuming an AWS IAM role.", Optional: true},
			"max_in_flight":             resourceschema.Int64Attribute{MarkdownDescription: "Maximum number of in-flight SQS messages.", Optional: true},
			"max_retries":               resourceschema.Int64Attribute{MarkdownDescription: "Maximum number of SQS send retries.", Optional: true},
			"message_deduplication_id":  resourceschema.StringAttribute{MarkdownDescription: "AWS FIFO message deduplication ID.", Optional: true},
			"message_group_id":          resourceschema.StringAttribute{MarkdownDescription: "AWS FIFO message group ID.", Optional: true},
			"metadata_exclude_prefixes": resourceschema.ListAttribute{MarkdownDescription: "Metadata header prefixes to exclude from SQS messages.", Optional: true, ElementType: types.StringType},
			"notify_policy":             resourceschema.StringAttribute{MarkdownDescription: "Delivery policy applied when notifying this outlet.", Optional: true},
			"profile":                   resourceschema.StringAttribute{MarkdownDescription: "AWS shared credential profile name for SQS.", Optional: true},
			"queue_url":                 resourceschema.StringAttribute{MarkdownDescription: "AWS SQS queue URL.", Required: true},
			"region":                    resourceschema.StringAttribute{MarkdownDescription: "AWS region for SQS.", Optional: true},
			"role_arn":                  resourceschema.StringAttribute{MarkdownDescription: "AWS IAM role ARN to assume for SQS.", Optional: true},
			"use_ec2_instance_role":     resourceschema.BoolAttribute{MarkdownDescription: "Use the EC2 instance role for SQS credentials.", Optional: true},
			"apply_status":              resourceschema.StringAttribute{MarkdownDescription: "Observed apply/runtime status for the outlet.", Computed: true},
			"last_apply_error":          resourceschema.StringAttribute{MarkdownDescription: "Last observed apply/runtime error message.", Computed: true},
			"config_generation":         resourceschema.Int64Attribute{MarkdownDescription: "Desired outlet config generation.", Computed: true},
			"applied_generation":        resourceschema.Int64Attribute{MarkdownDescription: "Last observed applied config generation.", Computed: true},
		},
		Blocks: map[string]resourceschema.Block{
			"credentials": resourceschema.SingleNestedBlock{
				MarkdownDescription: "Optional inline AWS credentials for this outlet. Values are write-only; rotate each credential by updating its value and incrementing the matching version.",
				Attributes: map[string]resourceschema.Attribute{
					"access_key_id":             resourceschema.StringAttribute{MarkdownDescription: "AWS access key ID.", Optional: true, Sensitive: true, WriteOnly: true},
					"access_key_id_version":     resourceschema.Int64Attribute{MarkdownDescription: "Monotonic version used to rotate `access_key_id`.", Optional: true},
					"secret_access_key":         resourceschema.StringAttribute{MarkdownDescription: "AWS secret access key.", Optional: true, Sensitive: true, WriteOnly: true},
					"secret_access_key_version": resourceschema.Int64Attribute{MarkdownDescription: "Monotonic version used to rotate `secret_access_key`.", Optional: true},
					"session_token":             resourceschema.StringAttribute{MarkdownDescription: "AWS session token.", Optional: true, Sensitive: true, WriteOnly: true},
					"session_token_version":     resourceschema.Int64Attribute{MarkdownDescription: "Monotonic version used to rotate `session_token`.", Optional: true},
				},
			},
		},
	}
}

func (r *awsSqsOutletResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	outletClient, ok := req.ProviderData.(awsSqsOutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected AWS SQS outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}
	r.client = outletClient
}

func (r *awsSqsOutletResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var plan awsSqsOutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	metadataExcludePrefixes, err := stringListValueOrNil(ctx, plan.MetadataExcludePrefixes)
	if err != nil {
		resp.Diagnostics.AddError("Invalid AWS SQS Metadata Exclude Prefixes", err.Error())
		return
	}
	credentials, diags := expandAwsCredentials(ctx, plan.Credentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	credentials, diags = enrichAwsCredentialsFromConfig(ctx, req.Config, credentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	outlet, err := r.client.CreateAwsSqsOutlet(ctx, client.AwsSqsOutletCreateRequest{
		Name:                    plan.Name.ValueString(),
		Description:             stringValueOrEmpty(plan.Description),
		Enabled:                 plan.Enabled.ValueBool(),
		AccessKeyID:             stringPointerValue(credentials.AccessKeyID),
		BackoffInitialInterval:  stringPointerValue(plan.BackoffInitialInterval),
		BackoffMaxElapsedTime:   stringPointerValue(plan.BackoffMaxElapsedTime),
		BackoffMaxInterval:      stringPointerValue(plan.BackoffMaxInterval),
		BatchingByteSize:        int64PointerValue(plan.BatchingByteSize),
		BatchingCheck:           stringPointerValue(plan.BatchingCheck),
		BatchingCount:           int64PointerValue(plan.BatchingCount),
		BatchingJitter:          float64PointerValue(plan.BatchingJitter),
		BatchingPeriod:          stringPointerValue(plan.BatchingPeriod),
		DelaySeconds:            int64PointerValue(plan.DelaySeconds),
		EndpointURL:             stringPointerValue(plan.EndpointURL),
		ExpiryWindow:            stringPointerValue(plan.ExpiryWindow),
		ExternalID:              stringPointerValue(plan.ExternalID),
		MaxInFlight:             int64PointerValue(plan.MaxInFlight),
		MaxRetries:              int64PointerValue(plan.MaxRetries),
		MessageDeduplicationID:  stringPointerValue(plan.MessageDeduplicationID),
		MessageGroupID:          stringPointerValue(plan.MessageGroupID),
		MetadataExcludePrefixes: metadataExcludePrefixes,
		NotifyPolicy:            stringPointerValue(plan.NotifyPolicy),
		Profile:                 stringPointerValue(plan.Profile),
		QueueURL:                plan.QueueURL.ValueString(),
		Region:                  stringPointerValue(plan.Region),
		RoleARN:                 stringPointerValue(plan.RoleARN),
		SecretAccessKey:         stringPointerValue(credentials.SecretAccessKey),
		SessionToken:            stringPointerValue(credentials.SessionToken),
		UseEC2InstanceRole:      boolPointerValue(plan.UseEC2InstanceRole),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create AWS SQS Outlet", err.Error())
		return
	}
	remoteCredentials, err := r.client.ListAwsSqsCredentials(ctx, outlet.ID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read AWS SQS Credentials", err.Error())
		return
	}
	credentialsObject, diags := buildAwsCredentialsObject(ctx, plan.Credentials, remoteCredentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "aws_sqs", outlet.ID, plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe AWS SQS Outlet Lifecycle", err.Error())
		return
	}
	state := awsSqsOutletModelFromAPI(outlet, credentialsObject)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *awsSqsOutletResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var state awsSqsOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	outlet, err := r.client.GetAwsSqsOutlet(ctx, state.ID.ValueInt64())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to Read AWS SQS Outlet", err.Error())
		return
	}
	remoteCredentials, err := r.client.ListAwsSqsCredentials(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read AWS SQS Credentials", err.Error())
		return
	}
	credentialsObject, diags := buildAwsCredentialsObject(ctx, state.Credentials, remoteCredentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	newState := awsSqsOutletModelFromAPI(outlet, credentialsObject)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *awsSqsOutletResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var plan awsSqsOutletResourceModel
	var state awsSqsOutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	credentials, diags := expandAwsCredentials(ctx, plan.Credentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	credentials, diags = enrichAwsCredentialsFromConfig(ctx, req.Config, credentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	stateCredentials, diags := expandAwsCredentials(ctx, state.Credentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.UpdateAwsSqsOutlet(ctx, state.ID.ValueInt64(), client.AwsSqsOutletUpdateRequest{
		Name:                    plan.Name.ValueString(),
		Description:             stringValueOrEmpty(plan.Description),
		Enabled:                 plan.Enabled.ValueBool(),
		AccessKeyID:             stringPointerValue(credentials.AccessKeyID),
		BackoffInitialInterval:  stringPointerValue(plan.BackoffInitialInterval),
		BackoffMaxElapsedTime:   stringPointerValue(plan.BackoffMaxElapsedTime),
		BackoffMaxInterval:      stringPointerValue(plan.BackoffMaxInterval),
		BatchingByteSize:        int64PointerValue(plan.BatchingByteSize),
		BatchingCheck:           stringPointerValue(plan.BatchingCheck),
		BatchingCount:           int64PointerValue(plan.BatchingCount),
		BatchingJitter:          float64PointerValue(plan.BatchingJitter),
		BatchingPeriod:          stringPointerValue(plan.BatchingPeriod),
		DelaySeconds:            int64PointerValue(plan.DelaySeconds),
		EndpointURL:             stringPointerValue(plan.EndpointURL),
		ExpiryWindow:            stringPointerValue(plan.ExpiryWindow),
		ExternalID:              stringPointerValue(plan.ExternalID),
		MaxInFlight:             int64PointerValue(plan.MaxInFlight),
		MaxRetries:              int64PointerValue(plan.MaxRetries),
		MessageDeduplicationID:  stringPointerValue(plan.MessageDeduplicationID),
		MessageGroupID:          stringPointerValue(plan.MessageGroupID),
		MetadataExcludePrefixes: mustStringListOrEmpty(ctx, plan.MetadataExcludePrefixes),
		NotifyPolicy:            stringPointerValue(plan.NotifyPolicy),
		Profile:                 stringPointerValue(plan.Profile),
		QueueURL:                plan.QueueURL.ValueString(),
		Region:                  stringPointerValue(plan.Region),
		RoleARN:                 stringPointerValue(plan.RoleARN),
		SecretAccessKey:         stringPointerValue(credentials.SecretAccessKey),
		SessionToken:            stringPointerValue(credentials.SessionToken),
		UseEC2InstanceRole:      boolPointerValue(plan.UseEC2InstanceRole),
	}); err != nil {
		resp.Diagnostics.AddError("Unable to Update AWS SQS Outlet", err.Error())
		return
	}

	fieldUpdater := awsSqsOutletFieldUpdater{client: r.client}
	mutations := []struct {
		name string
		err  error
	}{
		{"backoff_initial_interval", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.backoff.initial_interval", plan.BackoffInitialInterval, state.BackoffInitialInterval)},
		{"backoff_max_elapsed_time", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.backoff.max_elapsed_time", plan.BackoffMaxElapsedTime, state.BackoffMaxElapsedTime)},
		{"backoff_max_interval", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.backoff.max_interval", plan.BackoffMaxInterval, state.BackoffMaxInterval)},
		{"batching_byte_size", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.batching.byte_size", plan.BatchingByteSize, state.BatchingByteSize)},
		{"batching_check", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.batching.check", plan.BatchingCheck, state.BatchingCheck)},
		{"batching_count", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.batching.count", plan.BatchingCount, state.BatchingCount)},
		{"batching_jitter", syncFloat64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.batching.jitter", plan.BatchingJitter, state.BatchingJitter)},
		{"batching_period", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.batching.period", plan.BatchingPeriod, state.BatchingPeriod)},
		{"delay_seconds", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.delay_seconds", plan.DelaySeconds, state.DelaySeconds)},
		{"endpoint_url", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.endpoint_url", plan.EndpointURL, state.EndpointURL)},
		{"expiry_window", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.expiry_window", plan.ExpiryWindow, state.ExpiryWindow)},
		{"external_id", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.external_id", plan.ExternalID, state.ExternalID)},
		{"max_in_flight", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.max_in_flight", plan.MaxInFlight, state.MaxInFlight)},
		{"max_retries", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.max_retries", plan.MaxRetries, state.MaxRetries)},
		{"message_deduplication_id", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.message_deduplication_id", plan.MessageDeduplicationID, state.MessageDeduplicationID)},
		{"message_group_id", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.message_group_id", plan.MessageGroupID, state.MessageGroupID)},
		{"metadata_exclude_prefixes", syncStringListField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.metadata.exclude_prefixes", plan.MetadataExcludePrefixes, state.MetadataExcludePrefixes)},
		{"notify_policy", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "notify_policy", plan.NotifyPolicy, state.NotifyPolicy)},
		{"profile", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.profile", plan.Profile, state.Profile)},
		{"region", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.region", plan.Region, state.Region)},
		{"role_arn", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.role_arn", plan.RoleARN, state.RoleARN)},
		{"use_ec2_instance_role", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sqs.use_ec2_instance_role", plan.UseEC2InstanceRole, state.UseEC2InstanceRole)},
	}
	for _, mutation := range mutations {
		if mutation.err != nil {
			resp.Diagnostics.AddError("Unable to Update AWS SQS Field", fmt.Sprintf("%s: %s", mutation.name, mutation.err))
			return
		}
	}

	for _, role := range awsCredentialRolesToDelete(credentials, stateCredentials) {
		if _, err := r.client.DeleteAwsSqsMaterial(ctx, state.ID.ValueInt64(), role); err != nil && !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Delete AWS SQS Credential", fmt.Sprintf("%s: %s", role, err))
			return
		}
	}
	outlet, err := r.client.GetAwsSqsOutlet(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Refresh AWS SQS Outlet", err.Error())
		return
	}
	remoteCredentials, err := r.client.ListAwsSqsCredentials(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read AWS SQS Credentials", err.Error())
		return
	}
	credentialsObject, diags := buildAwsCredentialsObject(ctx, plan.Credentials, remoteCredentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "aws_sqs", outlet.ID, plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe AWS SQS Outlet Lifecycle", err.Error())
		return
	}
	newState := awsSqsOutletModelFromAPI(outlet, credentialsObject)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *awsSqsOutletResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var state awsSqsOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	operation, err := r.client.DeleteAwsSqsOutletOperation(ctx, state.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Delete AWS SQS Outlet", err.Error())
		return
	}
	if err == nil && operation.DesiredGeneration != 0 {
		if err := observeTypedOutletLifecycle(ctx, r.client, "aws_sqs", state.ID.ValueInt64(), operation.DesiredGeneration, false, true); err != nil {
			resp.Diagnostics.AddError("Unable to Observe AWS SQS Outlet Lifecycle", err.Error())
		}
	}
}

func (r *awsSqsOutletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid AWS SQS Outlet Import Identifier", fmt.Sprintf("Expected numeric outlet ID, got %q: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func awsSqsOutletModelFromAPI(outlet client.AwsSqsOutlet, credentials types.Object) awsSqsOutletResourceModel {
	return awsSqsOutletResourceModel{
		ID:                      types.Int64Value(outlet.ID),
		Name:                    types.StringValue(outlet.Name),
		Description:             nullableString(outlet.Description),
		Enabled:                 types.BoolValue(outlet.Enabled),
		BackoffInitialInterval:  nullableStringPointer(outlet.BackoffInitialInterval),
		BackoffMaxElapsedTime:   nullableStringPointer(outlet.BackoffMaxElapsedTime),
		BackoffMaxInterval:      nullableStringPointer(outlet.BackoffMaxInterval),
		BatchingByteSize:        nullableInt64Pointer(outlet.BatchingByteSize),
		BatchingCheck:           nullableStringPointer(outlet.BatchingCheck),
		BatchingCount:           nullableInt64Pointer(outlet.BatchingCount),
		BatchingJitter:          nullableFloat64Pointer(outlet.BatchingJitter),
		BatchingPeriod:          nullableStringPointer(outlet.BatchingPeriod),
		DelaySeconds:            nullableInt64Pointer(outlet.DelaySeconds),
		EndpointURL:             nullableStringPointer(outlet.EndpointURL),
		ExpiryWindow:            nullableStringPointer(outlet.ExpiryWindow),
		ExternalID:              nullableStringPointer(outlet.ExternalID),
		MaxInFlight:             nullableInt64Pointer(outlet.MaxInFlight),
		MaxRetries:              nullableInt64Pointer(outlet.MaxRetries),
		MessageDeduplicationID:  nullableStringPointer(outlet.MessageDeduplicationID),
		MessageGroupID:          nullableStringPointer(outlet.MessageGroupID),
		MetadataExcludePrefixes: nullableStringList(outlet.MetadataExcludePrefixes),
		NotifyPolicy:            nullableStringPointerPreservingEmpty(outlet.NotifyPolicy),
		Profile:                 nullableStringPointer(outlet.Profile),
		QueueURL:                types.StringValue(outlet.QueueURL),
		Region:                  nullableStringPointer(outlet.Region),
		RoleARN:                 nullableStringPointer(outlet.RoleARN),
		UseEC2InstanceRole:      nullableBoolPointer(outlet.UseEC2InstanceRole),
		Credentials:             credentials,
		ApplyStatus:             nullableString(outlet.ApplyStatus),
		LastApplyError:          nullableString(outlet.LastApplyError),
		ConfigGeneration:        types.Int64Value(outlet.ConfigGeneration),
		AppliedGeneration:       types.Int64Value(outlet.AppliedGeneration),
	}
}

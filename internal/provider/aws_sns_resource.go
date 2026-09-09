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
	_ resource.Resource                = &awsSnsOutletResource{}
	_ resource.ResourceWithConfigure   = &awsSnsOutletResource{}
	_ resource.ResourceWithImportState = &awsSnsOutletResource{}
)

type awsSnsOutletClient interface {
	CreateAwsSnsOutlet(context.Context, client.AwsSnsOutletCreateRequest) (client.AwsSnsOutlet, error)
	GetAwsSnsOutlet(context.Context, int64) (client.AwsSnsOutlet, error)
	UpdateAwsSnsOutlet(context.Context, int64, client.AwsSnsOutletUpdateRequest) error
	DeleteAwsSnsOutlet(context.Context, int64) error
	DeleteAwsSnsOutletOperation(context.Context, int64) (client.OutletOperation, error)
	GetTypedOutletStatus(context.Context, string, int64, int64) (client.OutletStatus, error)
	ListAwsSnsCredentials(context.Context, int64) ([]client.OutletMaterialSummary, error)
	DeleteAwsSnsMaterial(context.Context, int64, string) (client.OutletMaterialDeleteResponse, error)
	SetAwsSnsField(context.Context, int64, string, any) error
	DeleteAwsSnsField(context.Context, int64, string) error
}

type awsSnsOutletFieldUpdater struct {
	client awsSnsOutletClient
}

func (u awsSnsOutletFieldUpdater) SetField(ctx context.Context, outletID int64, field string, value any) error {
	return u.client.SetAwsSnsField(ctx, outletID, field, value)
}

func (u awsSnsOutletFieldUpdater) DeleteField(ctx context.Context, outletID int64, field string) error {
	return u.client.DeleteAwsSnsField(ctx, outletID, field)
}

type awsSnsOutletResource struct {
	client awsSnsOutletClient
}

type awsSnsOutletResourceModel struct {
	ID                      types.Int64  `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
	EndpointURL             types.String `tfsdk:"endpoint_url"`
	ExpiryWindow            types.String `tfsdk:"expiry_window"`
	ExternalID              types.String `tfsdk:"external_id"`
	MaxInFlight             types.Int64  `tfsdk:"max_in_flight"`
	MessageDeduplicationID  types.String `tfsdk:"message_deduplication_id"`
	MessageGroupID          types.String `tfsdk:"message_group_id"`
	MetadataExcludePrefixes types.List   `tfsdk:"metadata_exclude_prefixes"`
	NotifyPolicy            types.String `tfsdk:"notify_policy"`
	Profile                 types.String `tfsdk:"profile"`
	Region                  types.String `tfsdk:"region"`
	RoleARN                 types.String `tfsdk:"role_arn"`
	Timeout                 types.String `tfsdk:"timeout"`
	TopicARN                types.String `tfsdk:"topic_arn"`
	UseEC2InstanceRole      types.Bool   `tfsdk:"use_ec2_instance_role"`
	Credentials             types.Object `tfsdk:"credentials"`
	ApplyStatus             types.String `tfsdk:"apply_status"`
	LastApplyError          types.String `tfsdk:"last_apply_error"`
	ConfigGeneration        types.Int64  `tfsdk:"config_generation"`
	AppliedGeneration       types.Int64  `tfsdk:"applied_generation"`
}

func NewAwsSnsOutletResource() resource.Resource {
	return &awsSnsOutletResource{}
}

func (r *awsSnsOutletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_sns_outlet"
}

func (r *awsSnsOutletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		MarkdownDescription: "Manages a Noozle AWS SNS outlet.",
		Attributes: map[string]resourceschema.Attribute{
			"id": resourceschema.Int64Attribute{
				MarkdownDescription: "Unique identifier of the outlet.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name":                      resourceschema.StringAttribute{MarkdownDescription: "Name presented for the outlet.", Required: true},
			"description":               resourceschema.StringAttribute{MarkdownDescription: "Optional note explaining the outlet purpose.", Optional: true},
			"enabled":                   resourceschema.BoolAttribute{MarkdownDescription: "Whether the outlet is active.", Required: true},
			"endpoint_url":              resourceschema.StringAttribute{MarkdownDescription: "AWS service endpoint URL for SNS.", Optional: true},
			"expiry_window":             resourceschema.StringAttribute{MarkdownDescription: "Expiry window for refreshed SNS credentials.", Optional: true},
			"external_id":               resourceschema.StringAttribute{MarkdownDescription: "External ID for assuming an AWS IAM role.", Optional: true},
			"max_in_flight":             resourceschema.Int64Attribute{MarkdownDescription: "Maximum number of in-flight SNS messages.", Optional: true},
			"message_deduplication_id":  resourceschema.StringAttribute{MarkdownDescription: "AWS FIFO message deduplication ID.", Optional: true},
			"message_group_id":          resourceschema.StringAttribute{MarkdownDescription: "AWS FIFO message group ID.", Optional: true},
			"metadata_exclude_prefixes": resourceschema.ListAttribute{MarkdownDescription: "Metadata header prefixes to exclude from SNS messages.", Optional: true, ElementType: types.StringType},
			"notify_policy":             resourceschema.StringAttribute{MarkdownDescription: "Delivery policy applied when notifying this outlet.", Optional: true},
			"profile":                   resourceschema.StringAttribute{MarkdownDescription: "AWS shared credential profile name for SNS.", Optional: true},
			"region":                    resourceschema.StringAttribute{MarkdownDescription: "AWS region for SNS.", Optional: true},
			"role_arn":                  resourceschema.StringAttribute{MarkdownDescription: "AWS IAM role ARN to assume for SNS.", Optional: true},
			"timeout":                   resourceschema.StringAttribute{MarkdownDescription: "Timeout for SNS publish requests.", Optional: true},
			"topic_arn":                 resourceschema.StringAttribute{MarkdownDescription: "AWS SNS topic ARN.", Required: true},
			"use_ec2_instance_role":     resourceschema.BoolAttribute{MarkdownDescription: "Use the EC2 instance role for SNS credentials.", Optional: true},
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

func (r *awsSnsOutletResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	outletClient, ok := req.ProviderData.(awsSnsOutletClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected AWS SNS outlet client, got: %T. Please report this provider bug.", req.ProviderData),
		)
		return
	}

	r.client = outletClient
}

func (r *awsSnsOutletResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var plan awsSnsOutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	metadataExcludePrefixes, err := stringListValueOrNil(ctx, plan.MetadataExcludePrefixes)
	if err != nil {
		resp.Diagnostics.AddError("Invalid AWS SNS Metadata Exclude Prefixes", err.Error())
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

	outlet, err := r.client.CreateAwsSnsOutlet(ctx, client.AwsSnsOutletCreateRequest{
		Name:                    plan.Name.ValueString(),
		Description:             stringValueOrEmpty(plan.Description),
		Enabled:                 plan.Enabled.ValueBool(),
		AccessKeyID:             stringPointerValue(credentials.AccessKeyID),
		EndpointURL:             stringPointerValue(plan.EndpointURL),
		ExpiryWindow:            stringPointerValue(plan.ExpiryWindow),
		ExternalID:              stringPointerValue(plan.ExternalID),
		MaxInFlight:             int64PointerValue(plan.MaxInFlight),
		MessageDeduplicationID:  stringPointerValue(plan.MessageDeduplicationID),
		MessageGroupID:          stringPointerValue(plan.MessageGroupID),
		MetadataExcludePrefixes: metadataExcludePrefixes,
		NotifyPolicy:            stringPointerValue(plan.NotifyPolicy),
		Profile:                 stringPointerValue(plan.Profile),
		Region:                  stringPointerValue(plan.Region),
		RoleARN:                 stringPointerValue(plan.RoleARN),
		SecretAccessKey:         stringPointerValue(credentials.SecretAccessKey),
		SessionToken:            stringPointerValue(credentials.SessionToken),
		Timeout:                 stringPointerValue(plan.Timeout),
		TopicARN:                plan.TopicARN.ValueString(),
		UseEC2InstanceRole:      boolPointerValue(plan.UseEC2InstanceRole),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create AWS SNS Outlet", err.Error())
		return
	}

	remoteCredentials, err := r.client.ListAwsSnsCredentials(ctx, outlet.ID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read AWS SNS Credentials", err.Error())
		return
	}
	credentialsObject, diags := buildAwsCredentialsObject(ctx, plan.Credentials, remoteCredentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "aws_sns", outlet.ID, plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe AWS SNS Outlet Lifecycle", err.Error())
		return
	}
	state := awsSnsOutletModelFromAPI(outlet, credentialsObject)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *awsSnsOutletResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var state awsSnsOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	outlet, err := r.client.GetAwsSnsOutlet(ctx, state.ID.ValueInt64())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Unable to Read AWS SNS Outlet", err.Error())
		return
	}

	remoteCredentials, err := r.client.ListAwsSnsCredentials(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read AWS SNS Credentials", err.Error())
		return
	}
	credentialsObject, diags := buildAwsCredentialsObject(ctx, state.Credentials, remoteCredentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState := awsSnsOutletModelFromAPI(outlet, credentialsObject)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *awsSnsOutletResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var plan awsSnsOutletResourceModel
	var state awsSnsOutletResourceModel
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

	if err := r.client.UpdateAwsSnsOutlet(ctx, state.ID.ValueInt64(), client.AwsSnsOutletUpdateRequest{
		Name:                    plan.Name.ValueString(),
		Description:             stringValueOrEmpty(plan.Description),
		Enabled:                 plan.Enabled.ValueBool(),
		AccessKeyID:             stringPointerValue(credentials.AccessKeyID),
		EndpointURL:             stringPointerValue(plan.EndpointURL),
		ExpiryWindow:            stringPointerValue(plan.ExpiryWindow),
		ExternalID:              stringPointerValue(plan.ExternalID),
		MaxInFlight:             int64PointerValue(plan.MaxInFlight),
		MessageDeduplicationID:  stringPointerValue(plan.MessageDeduplicationID),
		MessageGroupID:          stringPointerValue(plan.MessageGroupID),
		MetadataExcludePrefixes: mustStringListOrEmpty(ctx, plan.MetadataExcludePrefixes),
		NotifyPolicy:            stringPointerValue(plan.NotifyPolicy),
		Profile:                 stringPointerValue(plan.Profile),
		Region:                  stringPointerValue(plan.Region),
		RoleARN:                 stringPointerValue(plan.RoleARN),
		SecretAccessKey:         stringPointerValue(credentials.SecretAccessKey),
		SessionToken:            stringPointerValue(credentials.SessionToken),
		Timeout:                 stringPointerValue(plan.Timeout),
		TopicARN:                plan.TopicARN.ValueString(),
		UseEC2InstanceRole:      boolPointerValue(plan.UseEC2InstanceRole),
	}); err != nil {
		resp.Diagnostics.AddError("Unable to Update AWS SNS Outlet", err.Error())
		return
	}

	fieldUpdater := awsSnsOutletFieldUpdater{client: r.client}
	mutations := []struct {
		name string
		err  error
	}{
		{"endpoint_url", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sns.endpoint_url", plan.EndpointURL, state.EndpointURL)},
		{"expiry_window", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sns.expiry_window", plan.ExpiryWindow, state.ExpiryWindow)},
		{"external_id", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sns.external_id", plan.ExternalID, state.ExternalID)},
		{"max_in_flight", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sns.max_in_flight", plan.MaxInFlight, state.MaxInFlight)},
		{"message_deduplication_id", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sns.message_deduplication_id", plan.MessageDeduplicationID, state.MessageDeduplicationID)},
		{"message_group_id", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sns.message_group_id", plan.MessageGroupID, state.MessageGroupID)},
		{"metadata_exclude_prefixes", syncStringListField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sns.metadata.exclude_prefixes", plan.MetadataExcludePrefixes, state.MetadataExcludePrefixes)},
		{"notify_policy", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "notify_policy", plan.NotifyPolicy, state.NotifyPolicy)},
		{"profile", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sns.profile", plan.Profile, state.Profile)},
		{"region", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sns.region", plan.Region, state.Region)},
		{"role_arn", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sns.role_arn", plan.RoleARN, state.RoleARN)},
		{"timeout", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sns.timeout", plan.Timeout, state.Timeout)},
		{"use_ec2_instance_role", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "aws_sns.use_ec2_instance_role", plan.UseEC2InstanceRole, state.UseEC2InstanceRole)},
	}
	for _, mutation := range mutations {
		if mutation.err != nil {
			resp.Diagnostics.AddError("Unable to Update AWS SNS Field", fmt.Sprintf("%s: %s", mutation.name, mutation.err))
			return
		}
	}

	for _, role := range awsCredentialRolesToDelete(credentials, stateCredentials) {
		if _, err := r.client.DeleteAwsSnsMaterial(ctx, state.ID.ValueInt64(), role); err != nil && !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Delete AWS SNS Credential", fmt.Sprintf("%s: %s", role, err))
			return
		}
	}

	outlet, err := r.client.GetAwsSnsOutlet(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Refresh AWS SNS Outlet", err.Error())
		return
	}
	remoteCredentials, err := r.client.ListAwsSnsCredentials(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read AWS SNS Credentials", err.Error())
		return
	}
	credentialsObject, diags := buildAwsCredentialsObject(ctx, plan.Credentials, remoteCredentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "aws_sns", outlet.ID, plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe AWS SNS Outlet Lifecycle", err.Error())
		return
	}
	newState := awsSnsOutletModelFromAPI(outlet, credentialsObject)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *awsSnsOutletResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var state awsSnsOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operation, err := r.client.DeleteAwsSnsOutletOperation(ctx, state.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Delete AWS SNS Outlet", err.Error())
		return
	}
	if err == nil && operation.DesiredGeneration != 0 {
		if err := observeTypedOutletLifecycle(ctx, r.client, "aws_sns", state.ID.ValueInt64(), operation.DesiredGeneration, false, true); err != nil {
			resp.Diagnostics.AddError("Unable to Observe AWS SNS Outlet Lifecycle", err.Error())
		}
	}
}

func (r *awsSnsOutletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid AWS SNS Outlet Import Identifier",
			fmt.Sprintf("Expected numeric outlet ID, got %q: %s", req.ID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func awsSnsOutletModelFromAPI(outlet client.AwsSnsOutlet, credentials types.Object) awsSnsOutletResourceModel {
	return awsSnsOutletResourceModel{
		ID:                      types.Int64Value(outlet.ID),
		Name:                    types.StringValue(outlet.Name),
		Description:             nullableString(outlet.Description),
		Enabled:                 types.BoolValue(outlet.Enabled),
		EndpointURL:             nullableStringPointer(outlet.EndpointURL),
		ExpiryWindow:            nullableStringPointer(outlet.ExpiryWindow),
		ExternalID:              nullableStringPointer(outlet.ExternalID),
		MaxInFlight:             nullableInt64Pointer(outlet.MaxInFlight),
		MessageDeduplicationID:  nullableStringPointer(outlet.MessageDeduplicationID),
		MessageGroupID:          nullableStringPointer(outlet.MessageGroupID),
		MetadataExcludePrefixes: nullableStringList(outlet.MetadataExcludePrefixes),
		NotifyPolicy:            nullableStringPointerPreservingEmpty(outlet.NotifyPolicy),
		Profile:                 nullableStringPointer(outlet.Profile),
		Region:                  nullableStringPointer(outlet.Region),
		RoleARN:                 nullableStringPointer(outlet.RoleARN),
		Timeout:                 nullableStringPointer(outlet.Timeout),
		TopicARN:                types.StringValue(outlet.TopicARN),
		UseEC2InstanceRole:      nullableBoolPointer(outlet.UseEC2InstanceRole),
		Credentials:             credentials,
		ApplyStatus:             nullableString(outlet.ApplyStatus),
		LastApplyError:          nullableString(outlet.LastApplyError),
		ConfigGeneration:        types.Int64Value(outlet.ConfigGeneration),
		AppliedGeneration:       types.Int64Value(outlet.AppliedGeneration),
	}
}

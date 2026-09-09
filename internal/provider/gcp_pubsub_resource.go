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
	_ resource.Resource                = &gcpPubsubOutletResource{}
	_ resource.ResourceWithConfigure   = &gcpPubsubOutletResource{}
	_ resource.ResourceWithImportState = &gcpPubsubOutletResource{}
)

type gcpPubsubOutletClient interface {
	CreateGcpPubsubOutlet(context.Context, client.GcpPubsubOutletCreateRequest) (client.GcpPubsubOutlet, error)
	GetGcpPubsubOutlet(context.Context, int64) (client.GcpPubsubOutlet, error)
	UpdateGcpPubsubOutlet(context.Context, int64, client.GcpPubsubOutletUpdateRequest) error
	DeleteGcpPubsubOutlet(context.Context, int64) error
	DeleteGcpPubsubOutletOperation(context.Context, int64) (client.OutletOperation, error)
	GetTypedOutletStatus(context.Context, string, int64, int64) (client.OutletStatus, error)
	ListGcpPubsubCredentials(context.Context, int64) ([]client.OutletMaterialSummary, error)
	DeleteGcpPubsubMaterial(context.Context, int64, string) (client.OutletMaterialDeleteResponse, error)
	SetGcpPubsubField(context.Context, int64, string, any) error
	DeleteGcpPubsubField(context.Context, int64, string) error
	typedOutletTransformClient
}

type gcpPubsubOutletFieldUpdater struct {
	client gcpPubsubOutletClient
}

func (u gcpPubsubOutletFieldUpdater) SetField(ctx context.Context, outletID int64, field string, value any) error {
	return u.client.SetGcpPubsubField(ctx, outletID, field, value)
}

func (u gcpPubsubOutletFieldUpdater) DeleteField(ctx context.Context, outletID int64, field string) error {
	return u.client.DeleteGcpPubsubField(ctx, outletID, field)
}

type gcpPubsubOutletResource struct {
	client gcpPubsubOutletClient
}

type gcpPubsubOutletResourceModel struct {
	ID                                types.Int64   `tfsdk:"id"`
	Name                              types.String  `tfsdk:"name"`
	Description                       types.String  `tfsdk:"description"`
	Enabled                           types.Bool    `tfsdk:"enabled"`
	BatchingByteSize                  types.Int64   `tfsdk:"batching_byte_size"`
	BatchingCheck                     types.String  `tfsdk:"batching_check"`
	BatchingCount                     types.Int64   `tfsdk:"batching_count"`
	BatchingJitter                    types.Float64 `tfsdk:"batching_jitter"`
	BatchingPeriod                    types.String  `tfsdk:"batching_period"`
	ByteThreshold                     types.Int64   `tfsdk:"byte_threshold"`
	CountThreshold                    types.Int64   `tfsdk:"count_threshold"`
	DelayThreshold                    types.String  `tfsdk:"delay_threshold"`
	Endpoint                          types.String  `tfsdk:"endpoint"`
	FlowControlLimitExceededBehavior  types.String  `tfsdk:"flow_control_limit_exceeded_behavior"`
	FlowControlMaxOutstandingBytes    types.Int64   `tfsdk:"flow_control_max_outstanding_bytes"`
	FlowControlMaxOutstandingMessages types.Int64   `tfsdk:"flow_control_max_outstanding_messages"`
	MaxInFlight                       types.Int64   `tfsdk:"max_in_flight"`
	MetadataExcludePrefixes           types.List    `tfsdk:"metadata_exclude_prefixes"`
	NotifyPolicy                      types.String  `tfsdk:"notify_policy"`
	OrderingKey                       types.String  `tfsdk:"ordering_key"`
	Project                           types.String  `tfsdk:"project"`
	PublishTimeout                    types.String  `tfsdk:"publish_timeout"`
	Topic                             types.String  `tfsdk:"topic"`
	Credentials                       types.Object  `tfsdk:"credentials"`
	ApplyStatus                       types.String  `tfsdk:"apply_status"`
	LastApplyError                    types.String  `tfsdk:"last_apply_error"`
	ConfigGeneration                  types.Int64   `tfsdk:"config_generation"`
	AppliedGeneration                 types.Int64   `tfsdk:"applied_generation"`
	Transform                         types.Object  `tfsdk:"transform"`
}

func NewGcpPubsubOutletResource() resource.Resource { return &gcpPubsubOutletResource{} }

func (r *gcpPubsubOutletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gcp_pubsub_outlet"
}

func (r *gcpPubsubOutletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		MarkdownDescription: "Manages a Noozle GCP Pub/Sub outlet.",
		Attributes: map[string]resourceschema.Attribute{
			"id":                                    resourceschema.Int64Attribute{MarkdownDescription: "Unique identifier of the outlet.", Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}},
			"name":                                  resourceschema.StringAttribute{MarkdownDescription: "Name presented for the outlet.", Required: true},
			"description":                           resourceschema.StringAttribute{MarkdownDescription: "Optional note explaining the outlet purpose.", Optional: true},
			"enabled":                               resourceschema.BoolAttribute{MarkdownDescription: "Whether the outlet is active.", Required: true},
			"batching_byte_size":                    resourceschema.Int64Attribute{MarkdownDescription: "Maximum bytes per Pub/Sub batch.", Optional: true},
			"batching_check":                        resourceschema.StringAttribute{MarkdownDescription: "Bloblang condition for flushing a Pub/Sub batch.", Optional: true},
			"batching_count":                        resourceschema.Int64Attribute{MarkdownDescription: "Maximum messages per Pub/Sub batch.", Optional: true},
			"batching_jitter":                       resourceschema.Float64Attribute{MarkdownDescription: "Randomized jitter to apply to the Pub/Sub batching period.", Optional: true},
			"batching_period":                       resourceschema.StringAttribute{MarkdownDescription: "Maximum wait period before flushing a Pub/Sub batch.", Optional: true},
			"byte_threshold":                        resourceschema.Int64Attribute{MarkdownDescription: "Publish when the Pub/Sub buffer reaches this many bytes.", Optional: true},
			"count_threshold":                       resourceschema.Int64Attribute{MarkdownDescription: "Publish when the Pub/Sub buffer reaches this many messages.", Optional: true},
			"delay_threshold":                       resourceschema.StringAttribute{MarkdownDescription: "Publish a non-empty Pub/Sub buffer after this delay.", Optional: true},
			"endpoint":                              resourceschema.StringAttribute{MarkdownDescription: "Optional Google Cloud Pub/Sub API endpoint override.", Optional: true},
			"flow_control_limit_exceeded_behavior":  resourceschema.StringAttribute{MarkdownDescription: "Behavior when the Pub/Sub flow-control buffer is full.", Optional: true},
			"flow_control_max_outstanding_bytes":    resourceschema.Int64Attribute{MarkdownDescription: "Maximum buffered Pub/Sub bytes waiting to be published.", Optional: true},
			"flow_control_max_outstanding_messages": resourceschema.Int64Attribute{MarkdownDescription: "Maximum buffered Pub/Sub messages waiting to be published.", Optional: true},
			"max_in_flight":                         resourceschema.Int64Attribute{MarkdownDescription: "Maximum number of in-flight Pub/Sub messages.", Optional: true},
			"metadata_exclude_prefixes":             resourceschema.ListAttribute{MarkdownDescription: "Metadata prefixes to exclude from Pub/Sub attributes.", Optional: true, ElementType: types.StringType},
			"notify_policy":                         resourceschema.StringAttribute{MarkdownDescription: "Delivery policy applied when notifying this outlet.", Optional: true},
			"ordering_key":                          resourceschema.StringAttribute{MarkdownDescription: "Pub/Sub ordering key interpolation.", Optional: true},
			"project":                               resourceschema.StringAttribute{MarkdownDescription: "Google Cloud project ID.", Required: true},
			"publish_timeout":                       resourceschema.StringAttribute{MarkdownDescription: "Maximum time to wait before abandoning a Pub/Sub publish attempt.", Optional: true},
			"topic":                                 resourceschema.StringAttribute{MarkdownDescription: "Google Cloud Pub/Sub topic.", Required: true},
			"apply_status":                          resourceschema.StringAttribute{MarkdownDescription: "Observed apply/runtime status for the outlet.", Computed: true},
			"last_apply_error":                      resourceschema.StringAttribute{MarkdownDescription: "Last observed apply/runtime error message.", Computed: true},
			"config_generation":                     resourceschema.Int64Attribute{MarkdownDescription: "Desired outlet config generation.", Computed: true},
			"applied_generation":                    resourceschema.Int64Attribute{MarkdownDescription: "Last observed applied config generation.", Computed: true},
			"transform": resourceschema.SingleNestedAttribute{MarkdownDescription: "Optional outlet transform selection and observed runtime state.", Optional: true, Computed: true, Attributes: map[string]resourceschema.Attribute{
				"selected": resourceschema.StringAttribute{MarkdownDescription: "Explicitly selected transform preset.", Optional: true}, "effective": resourceschema.StringAttribute{MarkdownDescription: "Effective transform after defaults are applied.", Computed: true}, "state": resourceschema.StringAttribute{MarkdownDescription: "Observed transform runtime state.", Computed: true}, "applied": resourceschema.BoolAttribute{MarkdownDescription: "Whether the transform has been applied.", Computed: true}, "generation": resourceschema.Int64Attribute{MarkdownDescription: "Observed transform generation.", Computed: true}, "last_error": resourceschema.StringAttribute{MarkdownDescription: "Last observed transform apply error.", Computed: true},
			}},
		},
		Blocks: map[string]resourceschema.Block{
			"credentials": resourceschema.SingleNestedBlock{
				MarkdownDescription: "Optional inline GCP credentials file for this outlet. The file content is write-only; rotate it by updating the content and incrementing `content_version`.",
				Attributes: map[string]resourceschema.Attribute{
					"filename":        resourceschema.StringAttribute{MarkdownDescription: "Filename to report for the uploaded service account credentials file.", Required: true},
					"content_base64":  resourceschema.StringAttribute{MarkdownDescription: "Base64-encoded service account credentials file.", Required: true, Sensitive: true, WriteOnly: true},
					"content_version": resourceschema.Int64Attribute{MarkdownDescription: "Monotonic version used to rotate the credentials file.", Required: true},
				},
			},
		},
	}
}

func (r *gcpPubsubOutletResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	outletClient, ok := req.ProviderData.(gcpPubsubOutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected GCP Pub/Sub outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}
	r.client = outletClient
}

func (r *gcpPubsubOutletResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var plan gcpPubsubOutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	metadataExcludePrefixes, err := stringListValueOrNil(ctx, plan.MetadataExcludePrefixes)
	if err != nil {
		resp.Diagnostics.AddError("Invalid GCP Pub/Sub Metadata Exclude Prefixes", err.Error())
		return
	}
	credentials, diags := expandGcpCredentials(ctx, plan.Credentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	credentials, diags = enrichGcpCredentialsFromConfig(ctx, req.Config, credentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	configuredTransform, transformDiags := selectedTransformFromConfig(ctx, req.Config)
	resp.Diagnostics.Append(transformDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	outlet, err := r.client.CreateGcpPubsubOutlet(ctx, client.GcpPubsubOutletCreateRequest{
		Name:                              plan.Name.ValueString(),
		Description:                       stringValueOrEmpty(plan.Description),
		Enabled:                           plan.Enabled.ValueBool(),
		BatchingByteSize:                  int64PointerValue(plan.BatchingByteSize),
		BatchingCheck:                     stringPointerValue(plan.BatchingCheck),
		BatchingCount:                     int64PointerValue(plan.BatchingCount),
		BatchingJitter:                    float64PointerValue(plan.BatchingJitter),
		BatchingPeriod:                    stringPointerValue(plan.BatchingPeriod),
		ByteThreshold:                     int64PointerValue(plan.ByteThreshold),
		CountThreshold:                    int64PointerValue(plan.CountThreshold),
		DelayThreshold:                    stringPointerValue(plan.DelayThreshold),
		Endpoint:                          stringPointerValue(plan.Endpoint),
		FlowControlLimitExceededBehavior:  stringPointerValue(plan.FlowControlLimitExceededBehavior),
		FlowControlMaxOutstandingBytes:    int64PointerValue(plan.FlowControlMaxOutstandingBytes),
		FlowControlMaxOutstandingMessages: int64PointerValue(plan.FlowControlMaxOutstandingMessages),
		MaxInFlight:                       int64PointerValue(plan.MaxInFlight),
		MetadataExcludePrefixes:           metadataExcludePrefixes,
		NotifyPolicy:                      stringPointerValue(plan.NotifyPolicy),
		OrderingKey:                       stringPointerValue(plan.OrderingKey),
		Project:                           plan.Project.ValueString(),
		PublishTimeout:                    stringPointerValue(plan.PublishTimeout),
		Topic:                             plan.Topic.ValueString(),
		CredentialsFilename:               stringPointerValue(credentials.Filename),
		CredentialsContentBase64:          stringPointerValue(credentials.ContentBase64),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create GCP Pub/Sub Outlet", err.Error())
		return
	}
	if !configuredTransform.IsNull() && !configuredTransform.IsUnknown() {
		if _, err := r.client.SetTypedOutletTransform(ctx, "gcp_pubsub", outlet.ID, configuredTransform.ValueString()); err != nil {
			if deleteErr := r.client.DeleteGcpPubsubOutlet(ctx, outlet.ID); deleteErr != nil && !client.IsNotFound(deleteErr) {
				resp.Diagnostics.AddError("Unable to Configure GCP Pub/Sub Transform", fmt.Sprintf("%s; cleanup failed: %s", err, deleteErr))
				return
			}
			resp.Diagnostics.AddError("Unable to Configure GCP Pub/Sub Transform", err.Error())
			return
		}
	}
	remoteCredentials, err := r.client.ListGcpPubsubCredentials(ctx, outlet.ID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read GCP Pub/Sub Credentials", err.Error())
		return
	}
	credentialsObject, diags := buildGcpCredentialsObject(ctx, plan.Credentials, remoteCredentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	transform, err := r.client.GetTypedOutletTransform(ctx, "gcp_pubsub", outlet.ID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read GCP Pub/Sub Transform", err.Error())
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "gcp_pubsub", outlet.ID, plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe GCP Pub/Sub Outlet Lifecycle", err.Error())
		return
	}
	state := gcpPubsubOutletModelFromAPI(outlet, credentialsObject, transform)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *gcpPubsubOutletResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var state gcpPubsubOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	outlet, err := r.client.GetGcpPubsubOutlet(ctx, state.ID.ValueInt64())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to Read GCP Pub/Sub Outlet", err.Error())
		return
	}
	remoteCredentials, err := r.client.ListGcpPubsubCredentials(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read GCP Pub/Sub Credentials", err.Error())
		return
	}
	credentialsObject, diags := buildGcpCredentialsObject(ctx, state.Credentials, remoteCredentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	transform, err := r.client.GetTypedOutletTransform(ctx, "gcp_pubsub", state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read GCP Pub/Sub Transform", err.Error())
		return
	}
	newState := gcpPubsubOutletModelFromAPI(outlet, credentialsObject, transform)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *gcpPubsubOutletResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var plan gcpPubsubOutletResourceModel
	var state gcpPubsubOutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	credentials, diags := expandGcpCredentials(ctx, plan.Credentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	credentials, diags = enrichGcpCredentialsFromConfig(ctx, req.Config, credentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	configuredTransform, transformDiags := selectedTransformFromConfig(ctx, req.Config)
	resp.Diagnostics.Append(transformDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	stateCredentials, diags := expandGcpCredentials(ctx, state.Credentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.UpdateGcpPubsubOutlet(ctx, state.ID.ValueInt64(), client.GcpPubsubOutletUpdateRequest{
		Name:                              plan.Name.ValueString(),
		Description:                       stringValueOrEmpty(plan.Description),
		Enabled:                           plan.Enabled.ValueBool(),
		BatchingByteSize:                  int64PointerValue(plan.BatchingByteSize),
		BatchingCheck:                     stringPointerValue(plan.BatchingCheck),
		BatchingCount:                     int64PointerValue(plan.BatchingCount),
		BatchingJitter:                    float64PointerValue(plan.BatchingJitter),
		BatchingPeriod:                    stringPointerValue(plan.BatchingPeriod),
		ByteThreshold:                     int64PointerValue(plan.ByteThreshold),
		CountThreshold:                    int64PointerValue(plan.CountThreshold),
		DelayThreshold:                    stringPointerValue(plan.DelayThreshold),
		Endpoint:                          stringPointerValue(plan.Endpoint),
		FlowControlLimitExceededBehavior:  stringPointerValue(plan.FlowControlLimitExceededBehavior),
		FlowControlMaxOutstandingBytes:    int64PointerValue(plan.FlowControlMaxOutstandingBytes),
		FlowControlMaxOutstandingMessages: int64PointerValue(plan.FlowControlMaxOutstandingMessages),
		MaxInFlight:                       int64PointerValue(plan.MaxInFlight),
		MetadataExcludePrefixes:           mustStringListOrEmpty(ctx, plan.MetadataExcludePrefixes),
		NotifyPolicy:                      stringPointerValue(plan.NotifyPolicy),
		OrderingKey:                       stringPointerValue(plan.OrderingKey),
		Project:                           plan.Project.ValueString(),
		PublishTimeout:                    stringPointerValue(plan.PublishTimeout),
		Topic:                             plan.Topic.ValueString(),
		CredentialsFilename:               stringPointerValue(credentials.Filename),
		CredentialsContentBase64:          stringPointerValue(credentials.ContentBase64),
	}); err != nil {
		resp.Diagnostics.AddError("Unable to Update GCP Pub/Sub Outlet", err.Error())
		return
	}

	fieldUpdater := gcpPubsubOutletFieldUpdater{client: r.client}
	mutations := []struct {
		name string
		err  error
	}{
		{"batching_byte_size", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.batching.byte_size", plan.BatchingByteSize, state.BatchingByteSize)},
		{"batching_check", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.batching.check", plan.BatchingCheck, state.BatchingCheck)},
		{"batching_count", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.batching.count", plan.BatchingCount, state.BatchingCount)},
		{"batching_jitter", syncFloat64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.batching.jitter", plan.BatchingJitter, state.BatchingJitter)},
		{"batching_period", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.batching.period", plan.BatchingPeriod, state.BatchingPeriod)},
		{"byte_threshold", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.byte_threshold", plan.ByteThreshold, state.ByteThreshold)},
		{"count_threshold", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.count_threshold", plan.CountThreshold, state.CountThreshold)},
		{"delay_threshold", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.delay_threshold", plan.DelayThreshold, state.DelayThreshold)},
		{"endpoint", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.endpoint", plan.Endpoint, state.Endpoint)},
		{"flow_control_limit_exceeded_behavior", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.flow_control.limit_exceeded_behavior", plan.FlowControlLimitExceededBehavior, state.FlowControlLimitExceededBehavior)},
		{"flow_control_max_outstanding_bytes", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.flow_control.max_outstanding_bytes", plan.FlowControlMaxOutstandingBytes, state.FlowControlMaxOutstandingBytes)},
		{"flow_control_max_outstanding_messages", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.flow_control.max_outstanding_messages", plan.FlowControlMaxOutstandingMessages, state.FlowControlMaxOutstandingMessages)},
		{"max_in_flight", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.max_in_flight", plan.MaxInFlight, state.MaxInFlight)},
		{"metadata_exclude_prefixes", syncStringListField(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.metadata.exclude_prefixes", plan.MetadataExcludePrefixes, state.MetadataExcludePrefixes)},
		{"notify_policy", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "notify_policy", plan.NotifyPolicy, state.NotifyPolicy)},
		{"ordering_key", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.ordering_key", plan.OrderingKey, state.OrderingKey)},
		{"publish_timeout", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "gcp_pubsub.publish_timeout", plan.PublishTimeout, state.PublishTimeout)},
	}
	for _, mutation := range mutations {
		if mutation.err != nil {
			resp.Diagnostics.AddError("Unable to Update GCP Pub/Sub Field", fmt.Sprintf("%s: %s", mutation.name, mutation.err))
			return
		}
	}

	if gcpCredentialRemoved(credentials, stateCredentials) {
		if _, err := r.client.DeleteGcpPubsubMaterial(ctx, state.ID.ValueInt64(), gcpPubsubCredentialsRole); err != nil && !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Delete GCP Pub/Sub Credential", err.Error())
			return
		}
	}
	if err := syncTypedOutletTransform(ctx, r.client, "gcp_pubsub", state.ID.ValueInt64(), configuredTransform, selectedTransform(state.Transform)); err != nil {
		resp.Diagnostics.AddError("Unable to Update GCP Pub/Sub Transform", err.Error())
		return
	}
	outlet, err := r.client.GetGcpPubsubOutlet(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Refresh GCP Pub/Sub Outlet", err.Error())
		return
	}
	remoteCredentials, err := r.client.ListGcpPubsubCredentials(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read GCP Pub/Sub Credentials", err.Error())
		return
	}
	credentialsObject, diags := buildGcpCredentialsObject(ctx, plan.Credentials, remoteCredentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	transform, err := r.client.GetTypedOutletTransform(ctx, "gcp_pubsub", state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read GCP Pub/Sub Transform", err.Error())
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "gcp_pubsub", outlet.ID, plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe GCP Pub/Sub Outlet Lifecycle", err.Error())
		return
	}
	newState := gcpPubsubOutletModelFromAPI(outlet, credentialsObject, transform)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *gcpPubsubOutletResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var state gcpPubsubOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	operation, err := r.client.DeleteGcpPubsubOutletOperation(ctx, state.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Delete GCP Pub/Sub Outlet", err.Error())
		return
	}
	if err == nil && operation.DesiredGeneration != 0 {
		if err := observeTypedOutletLifecycle(ctx, r.client, "gcp_pubsub", state.ID.ValueInt64(), operation.DesiredGeneration, false, true); err != nil {
			resp.Diagnostics.AddError("Unable to Observe GCP Pub/Sub Outlet Lifecycle", err.Error())
		}
	}
}

func (r *gcpPubsubOutletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid GCP Pub/Sub Outlet Import Identifier", fmt.Sprintf("Expected numeric outlet ID, got %q: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func gcpPubsubOutletModelFromAPI(outlet client.GcpPubsubOutlet, credentials types.Object, transform client.OutletTransformState) gcpPubsubOutletResourceModel {
	return gcpPubsubOutletResourceModel{
		ID:                                types.Int64Value(outlet.ID),
		Name:                              types.StringValue(outlet.Name),
		Description:                       nullableString(outlet.Description),
		Enabled:                           types.BoolValue(outlet.Enabled),
		BatchingByteSize:                  nullableInt64Pointer(outlet.BatchingByteSize),
		BatchingCheck:                     nullableStringPointer(outlet.BatchingCheck),
		BatchingCount:                     nullableInt64Pointer(outlet.BatchingCount),
		BatchingJitter:                    nullableFloat64Pointer(outlet.BatchingJitter),
		BatchingPeriod:                    nullableStringPointer(outlet.BatchingPeriod),
		ByteThreshold:                     nullableInt64Pointer(outlet.ByteThreshold),
		CountThreshold:                    nullableInt64Pointer(outlet.CountThreshold),
		DelayThreshold:                    nullableStringPointer(outlet.DelayThreshold),
		Endpoint:                          nullableStringPointer(outlet.Endpoint),
		FlowControlLimitExceededBehavior:  nullableStringPointer(outlet.FlowControlLimitExceededBehavior),
		FlowControlMaxOutstandingBytes:    nullableInt64Pointer(outlet.FlowControlMaxOutstandingBytes),
		FlowControlMaxOutstandingMessages: nullableInt64Pointer(outlet.FlowControlMaxOutstandingMessages),
		MaxInFlight:                       nullableInt64Pointer(outlet.MaxInFlight),
		MetadataExcludePrefixes:           nullableStringList(outlet.MetadataExcludePrefixes),
		NotifyPolicy:                      nullableStringPointerPreservingEmpty(outlet.NotifyPolicy),
		OrderingKey:                       nullableStringPointer(outlet.OrderingKey),
		Project:                           types.StringValue(outlet.Project),
		PublishTimeout:                    nullableStringPointer(outlet.PublishTimeout),
		Topic:                             types.StringValue(outlet.Topic),
		Credentials:                       credentials,
		ApplyStatus:                       nullableString(outlet.ApplyStatus),
		LastApplyError:                    nullableString(outlet.LastApplyError),
		ConfigGeneration:                  types.Int64Value(outlet.ConfigGeneration),
		AppliedGeneration:                 types.Int64Value(outlet.AppliedGeneration),
		Transform:                         transformObject(transform),
	}
}

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
	_ resource.Resource                = &sseOutletResource{}
	_ resource.ResourceWithConfigure   = &sseOutletResource{}
	_ resource.ResourceWithImportState = &sseOutletResource{}
)

type sseOutletClient interface {
	CreateSseOutlet(context.Context, client.SseOutletCreateRequest) (client.SseOutlet, error)
	GetSseOutlet(context.Context, int64) (client.SseOutlet, error)
	UpdateSseOutlet(context.Context, int64, client.SseOutletUpdateRequest) error
	DeleteSseOutlet(context.Context, int64) error
	SetSseEventType(context.Context, int64, string) error
	SetSsePath(context.Context, int64, string) error
	SetSseRetryMS(context.Context, int64, int64) error
	ClearSseNotifyPolicy(context.Context, int64) error
}

type sseOutletResource struct {
	client sseOutletClient
}

type sseOutletResourceModel struct {
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	EventType         types.String `tfsdk:"event_type"`
	Path              types.String `tfsdk:"path"`
	RetryMS           types.Int64  `tfsdk:"retry_ms"`
	NotifyPolicy      types.String `tfsdk:"notify_policy"`
	StreamURL         types.String `tfsdk:"stream_url"`
	ApplyStatus       types.String `tfsdk:"apply_status"`
	LastApplyError    types.String `tfsdk:"last_apply_error"`
	ConfigGeneration  types.Int64  `tfsdk:"config_generation"`
	AppliedGeneration types.Int64  `tfsdk:"applied_generation"`
}

func NewSseOutletResource() resource.Resource {
	return &sseOutletResource{}
}

func (r *sseOutletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sse_outlet"
}

func (r *sseOutletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		MarkdownDescription: "Manages a Noozle SSE outlet.",
		Attributes: map[string]resourceschema.Attribute{
			"id": resourceschema.Int64Attribute{
				MarkdownDescription: "Unique identifier of the outlet.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name":               resourceschema.StringAttribute{MarkdownDescription: "Name presented for the outlet.", Required: true},
			"description":        resourceschema.StringAttribute{MarkdownDescription: "Optional note explaining the outlet purpose.", Optional: true},
			"enabled":            resourceschema.BoolAttribute{MarkdownDescription: "Whether the outlet is active.", Required: true},
			"event_type":         resourceschema.StringAttribute{MarkdownDescription: "SSE event type.", Required: true},
			"path":               resourceschema.StringAttribute{MarkdownDescription: "SSE endpoint path.", Required: true},
			"retry_ms":           resourceschema.Int64Attribute{MarkdownDescription: "SSE retry interval in milliseconds.", Required: true},
			"notify_policy":      resourceschema.StringAttribute{MarkdownDescription: "Notification policy for routed matches. Use FIRST_ONLY to suppress duplicates or ALL to emit every match.", Optional: true},
			"stream_url":         resourceschema.StringAttribute{MarkdownDescription: "Stream URL for the SSE outlet.", Computed: true},
			"apply_status":       resourceschema.StringAttribute{MarkdownDescription: "Observed apply/runtime status for the outlet.", Computed: true},
			"last_apply_error":   resourceschema.StringAttribute{MarkdownDescription: "Last observed apply/runtime error message.", Computed: true},
			"config_generation":  resourceschema.Int64Attribute{MarkdownDescription: "Desired outlet config generation.", Computed: true},
			"applied_generation": resourceschema.Int64Attribute{MarkdownDescription: "Last observed applied config generation.", Computed: true},
		},
	}
}

func (r *sseOutletResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	outletClient, ok := req.ProviderData.(sseOutletClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected SSE outlet client, got: %T. Please report this provider bug.", req.ProviderData),
		)
		return
	}

	r.client = outletClient
}

func (r *sseOutletResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var plan sseOutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	outlet, err := r.client.CreateSseOutlet(ctx, client.SseOutletCreateRequest{
		Name:         plan.Name.ValueString(),
		Description:  stringValueOrEmpty(plan.Description),
		Enabled:      plan.Enabled.ValueBool(),
		EventType:    plan.EventType.ValueString(),
		Path:         plan.Path.ValueString(),
		RetryMS:      plan.RetryMS.ValueInt64(),
		NotifyPolicy: stringPointerValue(plan.NotifyPolicy),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create SSE Outlet", err.Error())
		return
	}

	state := sseOutletModelFromAPI(outlet)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sseOutletResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var state sseOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	outlet, err := r.client.GetSseOutlet(ctx, state.ID.ValueInt64())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Unable to Read SSE Outlet", err.Error())
		return
	}

	newState := sseOutletModelFromAPI(outlet)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *sseOutletResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var plan sseOutletResourceModel
	var state sseOutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateSseOutlet(ctx, state.ID.ValueInt64(), client.SseOutletUpdateRequest{
		Name:         plan.Name.ValueString(),
		Description:  stringValueOrEmpty(plan.Description),
		Enabled:      plan.Enabled.ValueBool(),
		EventType:    plan.EventType.ValueString(),
		Path:         plan.Path.ValueString(),
		RetryMS:      plan.RetryMS.ValueInt64(),
		NotifyPolicy: stringPointerValue(plan.NotifyPolicy),
	}); err != nil {
		resp.Diagnostics.AddError("Unable to Update SSE Outlet", err.Error())
		return
	}
	if plan.NotifyPolicy.IsNull() && !state.NotifyPolicy.IsNull() {
		if err := r.client.ClearSseNotifyPolicy(ctx, state.ID.ValueInt64()); err != nil {
			resp.Diagnostics.AddError("Unable to Clear SSE Notify Policy", err.Error())
			return
		}
	}

	outlet, err := r.client.GetSseOutlet(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Refresh SSE Outlet", err.Error())
		return
	}

	newState := sseOutletModelFromAPI(outlet)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *sseOutletResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var state sseOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSseOutlet(ctx, state.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Delete SSE Outlet", err.Error())
	}
}

func (r *sseOutletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid SSE Outlet Import Identifier",
			fmt.Sprintf("Expected numeric outlet ID, got %q: %s", req.ID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func sseOutletModelFromAPI(outlet client.SseOutlet) sseOutletResourceModel {
	return sseOutletResourceModel{
		ID:                types.Int64Value(outlet.ID),
		Name:              types.StringValue(outlet.Name),
		Description:       nullableString(outlet.Description),
		Enabled:           types.BoolValue(outlet.Enabled),
		EventType:         types.StringValue(outlet.EventType),
		Path:              types.StringValue(outlet.Path),
		RetryMS:           types.Int64Value(outlet.RetryMS),
		NotifyPolicy:      nullableStringPointer(outlet.NotifyPolicy),
		StreamURL:         nullableString(outlet.StreamURL),
		ApplyStatus:       nullableString(outlet.ApplyStatus),
		LastApplyError:    nullableString(outlet.LastApplyError),
		ConfigGeneration:  types.Int64Value(outlet.ConfigGeneration),
		AppliedGeneration: types.Int64Value(outlet.AppliedGeneration),
	}
}

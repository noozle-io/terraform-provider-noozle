package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

var (
	_ resource.Resource                = &queryOutletLinkResource{}
	_ resource.ResourceWithConfigure   = &queryOutletLinkResource{}
	_ resource.ResourceWithImportState = &queryOutletLinkResource{}
)

type queryOutletLinkClient interface {
	LinkOutletQuery(context.Context, int64, int64) error
	UnlinkOutletQuery(context.Context, int64, int64) error
	GetOutlet(context.Context, int64) (client.OutletRef, error)
}

type queryOutletLinkResource struct {
	client queryOutletLinkClient
}

type queryOutletLinkResourceModel struct {
	ID       types.String `tfsdk:"id"`
	OutletID types.Int64  `tfsdk:"outlet_id"`
	QueryID  types.Int64  `tfsdk:"query_id"`
}

func NewQueryOutletLinkResource() resource.Resource {
	return &queryOutletLinkResource{}
}

func (r *queryOutletLinkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_query_outlet_link"
}

func (r *queryOutletLinkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		MarkdownDescription: "Manages one query-to-outlet association for a pre-existing Noozle outlet.",
		Attributes: map[string]resourceschema.Attribute{
			"id": resourceschema.StringAttribute{
				MarkdownDescription: "Composite identifier in `<outlet_id>:<query_id>` format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"outlet_id": resourceschema.Int64Attribute{
				MarkdownDescription: "Identifier of the pre-existing outlet.",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"query_id": resourceschema.Int64Attribute{
				MarkdownDescription: "Identifier of the query to associate with the outlet.",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *queryOutletLinkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	linkClient, ok := req.ProviderData.(queryOutletLinkClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected query outlet link client, got: %T. Please report this provider bug.", req.ProviderData),
		)
		return
	}

	r.client = linkClient
}

func (r *queryOutletLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var plan queryOutletLinkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.LinkOutletQuery(ctx, plan.OutletID.ValueInt64(), plan.QueryID.ValueInt64()); err != nil {
		if addMissingTenantContextDiagnosticForError(err, &resp.Diagnostics) {
			return
		}
		resp.Diagnostics.AddError("Unable to Create Outlet Query Link", err.Error())
		return
	}

	plan.ID = types.StringValue(outletQueryLinkID(plan.OutletID.ValueInt64(), plan.QueryID.ValueInt64()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *queryOutletLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var state queryOutletLinkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	outlet, err := r.client.GetOutlet(ctx, state.OutletID.ValueInt64())
	if addMissingTenantContextDiagnosticForError(err, &resp.Diagnostics) {
		return
	}
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Unable to Read Outlet Query Link", err.Error())
		return
	}

	if !outletHasQuery(outlet, state.QueryID.ValueInt64()) {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(outletQueryLinkID(state.OutletID.ValueInt64(), state.QueryID.ValueInt64()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *queryOutletLinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan queryOutletLinkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(outletQueryLinkID(plan.OutletID.ValueInt64(), plan.QueryID.ValueInt64()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *queryOutletLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var state queryOutletLinkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UnlinkOutletQuery(ctx, state.OutletID.ValueInt64(), state.QueryID.ValueInt64())
	if addMissingTenantContextDiagnosticForError(err, &resp.Diagnostics) {
		return
	}
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Delete Outlet Query Link", err.Error())
	}
}

func (r *queryOutletLinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	outletID, queryID, err := parseOutletQueryLinkID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Outlet Query Link Import Identifier", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), outletQueryLinkID(outletID, queryID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("outlet_id"), outletID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("query_id"), queryID)...)
}

func outletQueryLinkID(outletID, queryID int64) string {
	return fmt.Sprintf("%d:%d", outletID, queryID)
}

func parseOutletQueryLinkID(value string) (int64, int64, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected import identifier in <outlet_id>:<query_id> format, got %q", value)
	}

	outletID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid outlet_id %q: %w", parts[0], err)
	}

	queryID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid query_id %q: %w", parts[1], err)
	}

	return outletID, queryID, nil
}

func outletHasQuery(outlet client.OutletRef, queryID int64) bool {
	for _, existing := range outlet.QueryIDs {
		if existing == queryID {
			return true
		}
	}

	return false
}

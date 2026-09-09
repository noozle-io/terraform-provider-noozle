package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.Resource                = &queryResource{}
	_ resource.ResourceWithConfigure   = &queryResource{}
	_ resource.ResourceWithImportState = &queryResource{}
)

type queryClient interface {
	CreateQuery(context.Context, client.QueryRequest) (client.Query, error)
	GetQuery(context.Context, int64) (client.Query, error)
	UpdateQuery(context.Context, int64, client.QueryRequest) (client.Query, error)
	DeleteQuery(context.Context, int64) error
}

type queryResource struct {
	client queryClient
}

type queryResourceModel struct {
	ID             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	FullExpression types.String `tfsdk:"full_expression"`
	Tags           types.List   `tfsdk:"tags"`
	IsActive       types.Bool   `tfsdk:"is_active"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

func NewQueryResource() resource.Resource {
	return &queryResource{}
}

func (r *queryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_query"
}

func (r *queryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		MarkdownDescription: "Manages a saved Noozle query.",
		Attributes: map[string]resourceschema.Attribute{
			"id": resourceschema.Int64Attribute{
				MarkdownDescription: "Unique identifier of the query.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": resourceschema.StringAttribute{
				MarkdownDescription: "Name presented for the saved query.",
				Required:            true,
			},
			"description": resourceschema.StringAttribute{
				MarkdownDescription: "Optional note explaining the query purpose.",
				Optional:            true,
			},
			"full_expression": resourceschema.StringAttribute{
				MarkdownDescription: "Serialized query expression DSL string.",
				Required:            true,
			},
			"tags": resourceschema.ListAttribute{
				MarkdownDescription: "Labels assigned to the query.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"is_active": resourceschema.BoolAttribute{
				MarkdownDescription: "Whether the saved query is active.",
				Required:            true,
			},
			"created_at": resourceschema.StringAttribute{
				MarkdownDescription: "Timestamp when the query was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": resourceschema.StringAttribute{
				MarkdownDescription: "Timestamp of the most recent update.",
				Computed:            true,
			},
		},
	}
}

func (r *queryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	queryClient, ok := req.ProviderData.(queryClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected query client, got: %T. Please report this provider bug.", req.ProviderData),
		)
		return
	}

	r.client = queryClient
}

func (r *queryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var plan queryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query, err := r.client.CreateQuery(ctx, plan.toQueryRequest(ctx, &resp.Diagnostics))
	if resp.Diagnostics.HasError() {
		return
	}
	if addMissingTenantContextDiagnosticForError(err, &resp.Diagnostics) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Query", err.Error())
		return
	}

	state, diags := queryModelFromAPI(ctx, query, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *queryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var state queryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query, err := r.client.GetQuery(ctx, state.ID.ValueInt64())
	if addMissingTenantContextDiagnosticForError(err, &resp.Diagnostics) {
		return
	}
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Unable to Read Query", err.Error())
		return
	}

	newState, diags := queryModelFromAPI(ctx, query, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *queryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var plan queryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query, err := r.client.UpdateQuery(ctx, plan.ID.ValueInt64(), plan.toQueryRequest(ctx, &resp.Diagnostics))
	if resp.Diagnostics.HasError() {
		return
	}
	if addMissingTenantContextDiagnosticForError(err, &resp.Diagnostics) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Query", err.Error())
		return
	}

	state, diags := queryModelFromAPI(ctx, query, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *queryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var state queryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteQuery(ctx, state.ID.ValueInt64())
	if addMissingTenantContextDiagnosticForError(err, &resp.Diagnostics) {
		return
	}
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Delete Query", err.Error())
	}
}

func (r *queryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Query Import Identifier",
			fmt.Sprintf("Expected numeric query ID, got %q: %s", req.ID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (m queryResourceModel) toQueryRequest(ctx context.Context, diags *diag.Diagnostics) client.QueryRequest {
	tags := make([]string, 0)
	if !m.Tags.IsNull() && !m.Tags.IsUnknown() {
		diags.Append(m.Tags.ElementsAs(ctx, &tags, false)...)
	}

	description := ""
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		description = m.Description.ValueString()
	}

	return client.QueryRequest{
		Name:           m.Name.ValueString(),
		Description:    description,
		FullExpression: m.FullExpression.ValueString(),
		Tags:           tags,
		IsActive:       m.IsActive.ValueBool(),
	}
}

func queryModelFromAPI(ctx context.Context, query client.Query, prior queryResourceModel) (queryResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	var tags types.List
	if len(query.Tags) == 0 && prior.Tags.IsNull() {
		tags = types.ListNull(types.StringType)
	} else {
		var tagDiags diag.Diagnostics
		tags, tagDiags = types.ListValueFrom(ctx, types.StringType, query.Tags)
		diags.Append(tagDiags...)
	}

	description := types.StringNull()
	if query.Description != "" {
		description = types.StringValue(query.Description)
	}

	return queryResourceModel{
		ID:             types.Int64Value(query.ID),
		Name:           types.StringValue(query.Name),
		Description:    description,
		FullExpression: types.StringValue(query.FullExpression),
		Tags:           tags,
		IsActive:       types.BoolValue(query.IsActive),
		CreatedAt:      types.StringValue(query.CreatedAt.Format(time.RFC3339)),
		UpdatedAt:      types.StringValue(query.UpdatedAt.Format(time.RFC3339)),
	}, diags
}

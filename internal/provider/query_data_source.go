package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

var (
	_ datasource.DataSource              = &queryDataSource{}
	_ datasource.DataSourceWithConfigure = &queryDataSource{}
)

type queryDataSource struct {
	client queryClient
}

type queryDataSourceModel struct {
	ID             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	FullExpression types.String `tfsdk:"full_expression"`
	Tags           types.List   `tfsdk:"tags"`
	IsActive       types.Bool   `tfsdk:"is_active"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

func NewQueryDataSource() datasource.DataSource {
	return &queryDataSource{}
}

func (d *queryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_query"
}

func (d *queryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Reads a saved Noozle query by ID.",
		Attributes: map[string]datasourceschema.Attribute{
			"id": datasourceschema.Int64Attribute{
				MarkdownDescription: "Unique identifier of the query.",
				Required:            true,
			},
			"name": datasourceschema.StringAttribute{
				MarkdownDescription: "Name presented for the saved query.",
				Computed:            true,
			},
			"description": datasourceschema.StringAttribute{
				MarkdownDescription: "Optional note explaining the query purpose.",
				Computed:            true,
			},
			"full_expression": datasourceschema.StringAttribute{
				MarkdownDescription: "Serialized query expression DSL string.",
				Computed:            true,
			},
			"tags": datasourceschema.ListAttribute{
				MarkdownDescription: "Labels assigned to the query.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"is_active": datasourceschema.BoolAttribute{
				MarkdownDescription: "Whether the saved query is active.",
				Computed:            true,
			},
			"created_at": datasourceschema.StringAttribute{
				MarkdownDescription: "Timestamp when the query was created.",
				Computed:            true,
			},
			"updated_at": datasourceschema.StringAttribute{
				MarkdownDescription: "Timestamp of the most recent update.",
				Computed:            true,
			},
		},
	}
}

func (d *queryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = queryClient
}

func (d *queryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var config queryDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query, err := d.client.GetQuery(ctx, config.ID.ValueInt64())
	if addMissingTenantContextDiagnosticForError(err, &resp.Diagnostics) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Query", err.Error())
		return
	}

	state, diags := queryDataSourceModelFromAPI(ctx, query)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func queryDataSourceModelFromAPI(ctx context.Context, query client.Query) (queryDataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	tags, tagDiags := types.ListValueFrom(ctx, types.StringType, query.Tags)
	diags.Append(tagDiags...)

	description := types.StringNull()
	if query.Description != "" {
		description = types.StringValue(query.Description)
	}

	return queryDataSourceModel{
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

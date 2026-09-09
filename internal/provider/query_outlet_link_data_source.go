package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &queryOutletLinkDataSource{}
	_ datasource.DataSourceWithConfigure = &queryOutletLinkDataSource{}
)

type queryOutletLinkDataSource struct {
	client queryOutletLinkClient
}

type queryOutletLinkDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	OutletID types.Int64  `tfsdk:"outlet_id"`
	QueryID  types.Int64  `tfsdk:"query_id"`
}

func NewQueryOutletLinkDataSource() datasource.DataSource {
	return &queryOutletLinkDataSource{}
}

func (d *queryOutletLinkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_query_outlet_link"
}

func (d *queryOutletLinkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Reads one query-to-outlet association for a pre-existing Noozle outlet.",
		Attributes: map[string]datasourceschema.Attribute{
			"id": datasourceschema.StringAttribute{
				MarkdownDescription: "Composite identifier in `<outlet_id>:<query_id>` format.",
				Computed:            true,
			},
			"outlet_id": datasourceschema.Int64Attribute{
				MarkdownDescription: "Identifier of the outlet.",
				Required:            true,
			},
			"query_id": datasourceschema.Int64Attribute{
				MarkdownDescription: "Identifier of the linked query.",
				Required:            true,
			},
		},
	}
}

func (d *queryOutletLinkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = linkClient
}

func (d *queryOutletLinkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var config queryOutletLinkDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	outlet, err := d.client.GetOutlet(ctx, config.OutletID.ValueInt64())
	if addMissingTenantContextDiagnosticForError(err, &resp.Diagnostics) {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Outlet Query Link", err.Error())
		return
	}

	if !outletHasQuery(outlet, config.QueryID.ValueInt64()) {
		resp.Diagnostics.AddError(
			"Outlet Query Link Not Found",
			fmt.Sprintf("Outlet %d is not linked to query %d.", config.OutletID.ValueInt64(), config.QueryID.ValueInt64()),
		)
		return
	}

	config.ID = types.StringValue(outletQueryLinkID(config.OutletID.ValueInt64(), config.QueryID.ValueInt64()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

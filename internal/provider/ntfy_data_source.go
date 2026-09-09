package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ntfyOutletDataSource{}
	_ datasource.DataSourceWithConfigure = &ntfyOutletDataSource{}
)

type ntfyOutletDataSource struct {
	client ntfyOutletClient
}

type ntfyOutletDataSourceModel struct {
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	Topic             types.String `tfsdk:"topic"`
	NotifyPolicy      types.String `tfsdk:"notify_policy"`
	ServerURL         types.String `tfsdk:"server_url"`
	ApplyStatus       types.String `tfsdk:"apply_status"`
	LastApplyError    types.String `tfsdk:"last_apply_error"`
	ConfigGeneration  types.Int64  `tfsdk:"config_generation"`
	AppliedGeneration types.Int64  `tfsdk:"applied_generation"`
}

func NewNtfyOutletDataSource() datasource.DataSource {
	return &ntfyOutletDataSource{}
}

func (d *ntfyOutletDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ntfy_outlet"
}

func (d *ntfyOutletDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Reads a Noozle ntfy outlet by ID.",
		Attributes: map[string]datasourceschema.Attribute{
			"id":                 datasourceschema.Int64Attribute{MarkdownDescription: "Unique identifier of the outlet.", Required: true},
			"name":               datasourceschema.StringAttribute{MarkdownDescription: "Name presented for the outlet.", Computed: true},
			"description":        datasourceschema.StringAttribute{MarkdownDescription: "Optional note explaining the outlet purpose.", Computed: true},
			"enabled":            datasourceschema.BoolAttribute{MarkdownDescription: "Whether the outlet is active.", Computed: true},
			"topic":              datasourceschema.StringAttribute{MarkdownDescription: "ntfy topic.", Computed: true},
			"notify_policy":      datasourceschema.StringAttribute{MarkdownDescription: "Notification policy for routed matches.", Computed: true},
			"server_url":         datasourceschema.StringAttribute{MarkdownDescription: "Optional ntfy base URL.", Computed: true},
			"apply_status":       datasourceschema.StringAttribute{MarkdownDescription: "Observed apply/runtime status for the outlet.", Computed: true},
			"last_apply_error":   datasourceschema.StringAttribute{MarkdownDescription: "Last observed apply/runtime error message.", Computed: true},
			"config_generation":  datasourceschema.Int64Attribute{MarkdownDescription: "Desired outlet config generation.", Computed: true},
			"applied_generation": datasourceschema.Int64Attribute{MarkdownDescription: "Last observed applied config generation.", Computed: true},
		},
	}
}

func (d *ntfyOutletDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	outletClient, ok := req.ProviderData.(ntfyOutletClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected ntfy outlet client, got: %T. Please report this provider bug.", req.ProviderData),
		)
		return
	}

	d.client = outletClient
}

func (d *ntfyOutletDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var config ntfyOutletDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	outlet, err := d.client.GetNtfyOutlet(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read ntfy Outlet", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &ntfyOutletDataSourceModel{
		ID:                types.Int64Value(outlet.ID),
		Name:              types.StringValue(outlet.Name),
		Description:       nullableString(outlet.Description),
		Enabled:           types.BoolValue(outlet.Enabled),
		Topic:             types.StringValue(outlet.Topic),
		NotifyPolicy:      nullableStringPointer(outlet.NotifyPolicy),
		ServerURL:         nullableStringPointer(outlet.ServerURL),
		ApplyStatus:       nullableString(outlet.ApplyStatus),
		LastApplyError:    nullableString(outlet.LastApplyError),
		ConfigGeneration:  types.Int64Value(outlet.ConfigGeneration),
		AppliedGeneration: types.Int64Value(outlet.AppliedGeneration),
	})...)
}

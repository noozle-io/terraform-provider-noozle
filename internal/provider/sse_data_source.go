package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &sseOutletDataSource{}
	_ datasource.DataSourceWithConfigure = &sseOutletDataSource{}
)

type sseOutletDataSource struct {
	client sseOutletClient
}

type sseOutletDataSourceModel struct {
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

func NewSseOutletDataSource() datasource.DataSource {
	return &sseOutletDataSource{}
}

func (d *sseOutletDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sse_outlet"
}

func (d *sseOutletDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Reads a Noozle SSE outlet by ID.",
		Attributes: map[string]datasourceschema.Attribute{
			"id":                 datasourceschema.Int64Attribute{MarkdownDescription: "Unique identifier of the outlet.", Required: true},
			"name":               datasourceschema.StringAttribute{MarkdownDescription: "Name presented for the outlet.", Computed: true},
			"description":        datasourceschema.StringAttribute{MarkdownDescription: "Optional note explaining the outlet purpose.", Computed: true},
			"enabled":            datasourceschema.BoolAttribute{MarkdownDescription: "Whether the outlet is active.", Computed: true},
			"event_type":         datasourceschema.StringAttribute{MarkdownDescription: "SSE event type.", Computed: true},
			"path":               datasourceschema.StringAttribute{MarkdownDescription: "SSE endpoint path.", Computed: true},
			"retry_ms":           datasourceschema.Int64Attribute{MarkdownDescription: "SSE retry interval in milliseconds.", Computed: true},
			"notify_policy":      datasourceschema.StringAttribute{MarkdownDescription: "Notification policy for routed matches.", Computed: true},
			"stream_url":         datasourceschema.StringAttribute{MarkdownDescription: "Stream URL for the SSE outlet.", Computed: true},
			"apply_status":       datasourceschema.StringAttribute{MarkdownDescription: "Observed apply/runtime status for the outlet.", Computed: true},
			"last_apply_error":   datasourceschema.StringAttribute{MarkdownDescription: "Last observed apply/runtime error message.", Computed: true},
			"config_generation":  datasourceschema.Int64Attribute{MarkdownDescription: "Desired outlet config generation.", Computed: true},
			"applied_generation": datasourceschema.Int64Attribute{MarkdownDescription: "Last observed applied config generation.", Computed: true},
		},
	}
}

func (d *sseOutletDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = outletClient
}

func (d *sseOutletDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var config sseOutletDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	outlet, err := d.client.GetSseOutlet(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read SSE Outlet", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &sseOutletDataSourceModel{
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
	})...)
}

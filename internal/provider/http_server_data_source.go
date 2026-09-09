package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

var (
	_ datasource.DataSource              = &httpServerOutletDataSource{}
	_ datasource.DataSourceWithConfigure = &httpServerOutletDataSource{}
)

type httpServerOutletDataSource struct {
	client httpServerOutletClient
}

type httpServerOutletDataSourceModel struct {
	ID                 types.Int64  `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	Address            types.String `tfsdk:"address"`
	AllowedVerbs       types.List   `tfsdk:"allowed_verbs"`
	CorsAllowedHeaders types.List   `tfsdk:"cors_allowed_headers"`
	CorsAllowedMethods types.List   `tfsdk:"cors_allowed_methods"`
	CorsAllowedOrigins types.List   `tfsdk:"cors_allowed_origins"`
	CorsEnabled        types.Bool   `tfsdk:"cors_enabled"`
	Heartbeat          types.String `tfsdk:"heartbeat"`
	NotifyPolicy       types.String `tfsdk:"notify_policy"`
	Path               types.String `tfsdk:"path"`
	PingPeriod         types.String `tfsdk:"ping_period"`
	PongWait           types.String `tfsdk:"pong_wait"`
	StreamFormat       types.String `tfsdk:"stream_format"`
	StreamPath         types.String `tfsdk:"stream_path"`
	Timeout            types.String `tfsdk:"timeout"`
	WriteWait          types.String `tfsdk:"write_wait"`
	WsMessageType      types.String `tfsdk:"ws_message_type"`
	WsPath             types.String `tfsdk:"ws_path"`
	TLS                types.Object `tfsdk:"tls"`
	ApplyStatus        types.String `tfsdk:"apply_status"`
	LastApplyError     types.String `tfsdk:"last_apply_error"`
	ConfigGeneration   types.Int64  `tfsdk:"config_generation"`
	AppliedGeneration  types.Int64  `tfsdk:"applied_generation"`
	Transform          types.Object `tfsdk:"transform"`
}

func NewHttpServerOutletDataSource() datasource.DataSource {
	return &httpServerOutletDataSource{}
}

func (d *httpServerOutletDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_server_outlet"
}

func (d *httpServerOutletDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Reads a Noozle HTTP Server outlet by ID.",
		Attributes: map[string]datasourceschema.Attribute{
			"id":                   datasourceschema.Int64Attribute{MarkdownDescription: "Unique identifier of the outlet.", Required: true},
			"name":                 datasourceschema.StringAttribute{MarkdownDescription: "Name presented for the outlet.", Computed: true},
			"description":          datasourceschema.StringAttribute{MarkdownDescription: "Optional note explaining the outlet purpose.", Computed: true},
			"enabled":              datasourceschema.BoolAttribute{MarkdownDescription: "Whether the outlet is active.", Computed: true},
			"address":              datasourceschema.StringAttribute{MarkdownDescription: "Custom HTTP server bind address.", Computed: true},
			"allowed_verbs":        datasourceschema.ListAttribute{MarkdownDescription: "HTTP verbs allowed for path and stream path endpoints.", Computed: true, ElementType: types.StringType},
			"cors_allowed_headers": datasourceschema.ListAttribute{MarkdownDescription: "Additional CORS headers to allow.", Computed: true, ElementType: types.StringType},
			"cors_allowed_methods": datasourceschema.ListAttribute{MarkdownDescription: "CORS methods to allow.", Computed: true, ElementType: types.StringType},
			"cors_allowed_origins": datasourceschema.ListAttribute{MarkdownDescription: "CORS origins to allow.", Computed: true, ElementType: types.StringType},
			"cors_enabled":         datasourceschema.BoolAttribute{MarkdownDescription: "Enable CORS handling for a custom HTTP server address.", Computed: true},
			"heartbeat":            datasourceschema.StringAttribute{MarkdownDescription: "HTTP stream heartbeat interval.", Computed: true},
			"notify_policy":        datasourceschema.StringAttribute{MarkdownDescription: "Delivery policy applied when notifying this outlet.", Computed: true},
			"path":                 datasourceschema.StringAttribute{MarkdownDescription: "HTTP server path.", Computed: true},
			"ping_period":          datasourceschema.StringAttribute{MarkdownDescription: "Websocket ping period.", Computed: true},
			"pong_wait":            datasourceschema.StringAttribute{MarkdownDescription: "Websocket pong wait timeout.", Computed: true},
			"stream_format":        datasourceschema.StringAttribute{MarkdownDescription: "HTTP stream endpoint format.", Computed: true},
			"stream_path":          datasourceschema.StringAttribute{MarkdownDescription: "HTTP server stream path.", Computed: true},
			"timeout":              datasourceschema.StringAttribute{MarkdownDescription: "HTTP path endpoint timeout.", Computed: true},
			"write_wait":           datasourceschema.StringAttribute{MarkdownDescription: "Websocket write timeout.", Computed: true},
			"ws_message_type":      datasourceschema.StringAttribute{MarkdownDescription: "Websocket message type.", Computed: true},
			"ws_path":              datasourceschema.StringAttribute{MarkdownDescription: "HTTP server websocket path.", Computed: true},
			"apply_status":         datasourceschema.StringAttribute{MarkdownDescription: "Observed apply/runtime status for the outlet.", Computed: true},
			"last_apply_error":     datasourceschema.StringAttribute{MarkdownDescription: "Last observed apply/runtime error message.", Computed: true},
			"config_generation":    datasourceschema.Int64Attribute{MarkdownDescription: "Desired outlet config generation.", Computed: true},
			"applied_generation":   datasourceschema.Int64Attribute{MarkdownDescription: "Last observed applied config generation.", Computed: true},
			"transform": datasourceschema.SingleNestedAttribute{Computed: true, Attributes: map[string]datasourceschema.Attribute{
				"selected": datasourceschema.StringAttribute{Computed: true}, "effective": datasourceschema.StringAttribute{Computed: true}, "state": datasourceschema.StringAttribute{Computed: true}, "applied": datasourceschema.BoolAttribute{Computed: true}, "generation": datasourceschema.Int64Attribute{Computed: true}, "last_error": datasourceschema.StringAttribute{Computed: true},
			}},
		},
		Blocks: map[string]datasourceschema.Block{
			"tls": datasourceschema.SingleNestedBlock{
				MarkdownDescription: "Observed TLS feature state for the HTTP server outlet.",
				Attributes: map[string]datasourceschema.Attribute{
					"enabled":                     datasourceschema.BoolAttribute{MarkdownDescription: "Whether the TLS feature is enabled.", Computed: true},
					"state":                       datasourceschema.StringAttribute{MarkdownDescription: "Observed TLS feature state.", Computed: true},
					"desired_generation":          datasourceschema.Int64Attribute{MarkdownDescription: "Desired TLS feature generation.", Computed: true},
					"applied_generation":          datasourceschema.Int64Attribute{MarkdownDescription: "Applied TLS feature generation.", Computed: true},
					"last_error":                  datasourceschema.StringAttribute{MarkdownDescription: "Last observed TLS feature error.", Computed: true},
					"server_cert_filename":        datasourceschema.StringAttribute{MarkdownDescription: "Observed TLS certificate filename.", Computed: true},
					"server_cert_checksum_sha256": datasourceschema.StringAttribute{MarkdownDescription: "Observed checksum of the TLS certificate file.", Computed: true},
					"server_key_filename":         datasourceschema.StringAttribute{MarkdownDescription: "Observed TLS private key filename.", Computed: true},
					"server_key_checksum_sha256":  datasourceschema.StringAttribute{MarkdownDescription: "Observed checksum of the TLS private key file.", Computed: true},
				},
			},
		},
	}
}

func (d *httpServerOutletDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	outletClient, ok := req.ProviderData.(httpServerOutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected HTTP server outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}

	d.client = outletClient
}

func (d *httpServerOutletDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var config httpServerOutletDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	outlet, err := d.client.GetHttpServerOutlet(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read HTTP Server Outlet", err.Error())
		return
	}

	tlsFeature, err := d.client.GetHttpServerTLSFeature(ctx, config.ID.ValueInt64())
	var tlsFeaturePtr *client.OutletFeatureDetail
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Read HTTP Server TLS Feature", err.Error())
			return
		}
	} else {
		tlsFeaturePtr = &tlsFeature
	}

	credentials, err := d.client.ListHttpServerCredentials(ctx, config.ID.ValueInt64())
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Read HTTP Server Credentials", err.Error())
			return
		}
		credentials = nil
	}

	tlsObject, diags := buildHttpServerTLSDataSourceObject(tlsFeaturePtr, credentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	transform, err := d.client.GetTypedOutletTransform(ctx, "http_server", config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read HTTP Server Transform", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &httpServerOutletDataSourceModel{
		ID:                 types.Int64Value(outlet.ID),
		Name:               types.StringValue(outlet.Name),
		Description:        nullableString(outlet.Description),
		Enabled:            types.BoolValue(outlet.Enabled),
		Address:            nullableStringPointer(outlet.Address),
		AllowedVerbs:       nullableStringList(outlet.AllowedVerbs),
		CorsAllowedHeaders: nullableStringList(outlet.CorsAllowedHeaders),
		CorsAllowedMethods: nullableStringList(outlet.CorsAllowedMethods),
		CorsAllowedOrigins: nullableStringList(outlet.CorsAllowedOrigins),
		CorsEnabled:        nullableBoolPointer(outlet.CorsEnabled),
		Heartbeat:          nullableStringPointer(outlet.Heartbeat),
		NotifyPolicy:       nullableStringPointerPreservingEmpty(outlet.NotifyPolicy),
		Path:               types.StringValue(outlet.Path),
		PingPeriod:         nullableStringPointer(outlet.PingPeriod),
		PongWait:           nullableStringPointer(outlet.PongWait),
		StreamFormat:       nullableStringPointer(outlet.StreamFormat),
		StreamPath:         nullableStringPointer(outlet.StreamPath),
		Timeout:            nullableStringPointer(outlet.Timeout),
		WriteWait:          nullableStringPointer(outlet.WriteWait),
		WsMessageType:      nullableStringPointer(outlet.WsMessageType),
		WsPath:             nullableStringPointer(outlet.WsPath),
		TLS:                tlsObject,
		ApplyStatus:        nullableString(outlet.ApplyStatus),
		LastApplyError:     nullableString(outlet.LastApplyError),
		ConfigGeneration:   types.Int64Value(outlet.ConfigGeneration),
		AppliedGeneration:  types.Int64Value(outlet.AppliedGeneration),
		Transform:          transformObject(transform),
	})...)
}

func buildHttpServerTLSDataSourceObject(feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	serverCert := findMaterialOrNil(credentials, httpServerTLSCertRole)
	serverKey := findMaterialOrNil(credentials, httpServerTLSKeyRole)
	if feature == nil && serverCert == nil && serverKey == nil {
		return types.ObjectNull(map[string]attr.Type{
			"enabled":                     types.BoolType,
			"state":                       types.StringType,
			"desired_generation":          types.Int64Type,
			"applied_generation":          types.Int64Type,
			"last_error":                  types.StringType,
			"server_cert_filename":        types.StringType,
			"server_cert_checksum_sha256": types.StringType,
			"server_key_filename":         types.StringType,
			"server_key_checksum_sha256":  types.StringType,
		}), nil
	}

	enabled := types.BoolNull()
	state := types.StringNull()
	desiredGeneration := types.Int64Null()
	appliedGeneration := types.Int64Null()
	lastError := types.StringNull()
	if feature != nil {
		enabled = types.BoolValue(feature.Enabled)
		state = nullableString(feature.State)
		desiredGeneration = nullableInt64Pointer(feature.DesiredGeneration)
		appliedGeneration = nullableInt64Pointer(feature.AppliedGeneration)
		lastError = nullableString(feature.LastError)
	}

	return types.ObjectValue(map[string]attr.Type{
		"enabled":                     types.BoolType,
		"state":                       types.StringType,
		"desired_generation":          types.Int64Type,
		"applied_generation":          types.Int64Type,
		"last_error":                  types.StringType,
		"server_cert_filename":        types.StringType,
		"server_cert_checksum_sha256": types.StringType,
		"server_key_filename":         types.StringType,
		"server_key_checksum_sha256":  types.StringType,
	}, map[string]attr.Value{
		"enabled":                     enabled,
		"state":                       state,
		"desired_generation":          desiredGeneration,
		"applied_generation":          appliedGeneration,
		"last_error":                  lastError,
		"server_cert_filename":        nullableString(materialField(serverCert, func(m *client.OutletMaterialSummary) string { return m.Filename })),
		"server_cert_checksum_sha256": nullableString(materialField(serverCert, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 })),
		"server_key_filename":         nullableString(materialField(serverKey, func(m *client.OutletMaterialSummary) string { return m.Filename })),
		"server_key_checksum_sha256":  nullableString(materialField(serverKey, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 })),
	})
}

func materialField(material *client.OutletMaterialSummary, selector func(*client.OutletMaterialSummary) string) string {
	if material == nil {
		return ""
	}
	return selector(material)
}

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

var (
	_ resource.Resource                   = &httpServerOutletResource{}
	_ resource.ResourceWithConfigure      = &httpServerOutletResource{}
	_ resource.ResourceWithImportState    = &httpServerOutletResource{}
	_ resource.ResourceWithValidateConfig = &httpServerOutletResource{}
)

type httpServerOutletClient interface {
	CreateHttpServerOutlet(context.Context, client.HttpServerOutletCreateRequest) (client.HttpServerOutlet, error)
	GetHttpServerOutlet(context.Context, int64) (client.HttpServerOutlet, error)
	UpdateHttpServerOutlet(context.Context, int64, client.HttpServerOutletUpdateRequest) error
	DeleteHttpServerOutlet(context.Context, int64) error
	DeleteHttpServerOutletOperation(context.Context, int64) (client.OutletOperation, error)
	GetTypedOutletStatus(context.Context, string, int64, int64) (client.OutletStatus, error)
	SetHttpServerField(context.Context, int64, string, any) error
	DeleteHttpServerField(context.Context, int64, string) error
	GetHttpServerTLSFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureHttpServerTLSFeature(context.Context, int64, client.HttpServerTLSFeatureRequest) error
	UpdateHttpServerTLSFeature(context.Context, int64, client.HttpServerTLSFeatureRequest) error
	DeleteHttpServerTLSFeature(context.Context, int64) error
	ListHttpServerCredentials(context.Context, int64) ([]client.OutletMaterialSummary, error)
	DeleteHttpServerCredential(context.Context, int64, string) (client.OutletMaterialDeleteResponse, error)
	typedOutletTransformClient
}

type httpServerOutletFieldUpdater struct {
	client httpServerOutletClient
}

func (u httpServerOutletFieldUpdater) SetField(ctx context.Context, outletID int64, field string, value any) error {
	return u.client.SetHttpServerField(ctx, outletID, field, value)
}

func (u httpServerOutletFieldUpdater) DeleteField(ctx context.Context, outletID int64, field string) error {
	return u.client.DeleteHttpServerField(ctx, outletID, field)
}

type httpServerOutletResource struct {
	client httpServerOutletClient
}

type httpServerOutletResourceModel struct {
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

func NewHttpServerOutletResource() resource.Resource {
	return &httpServerOutletResource{}
}

func (r *httpServerOutletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_server_outlet"
}

func (r *httpServerOutletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		MarkdownDescription: "Manages a Noozle HTTP Server outlet.",
		Attributes: map[string]resourceschema.Attribute{
			"id":                   resourceschema.Int64Attribute{MarkdownDescription: "Unique identifier of the outlet.", Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}},
			"name":                 resourceschema.StringAttribute{MarkdownDescription: "Name presented for the outlet.", Required: true},
			"description":          resourceschema.StringAttribute{MarkdownDescription: "Optional note explaining the outlet purpose.", Optional: true},
			"enabled":              resourceschema.BoolAttribute{MarkdownDescription: "Whether the outlet is active.", Required: true},
			"address":              resourceschema.StringAttribute{MarkdownDescription: "Custom HTTP server bind address.", Optional: true},
			"allowed_verbs":        resourceschema.ListAttribute{MarkdownDescription: "HTTP verbs allowed for path and stream path endpoints.", Optional: true, ElementType: types.StringType},
			"cors_allowed_headers": resourceschema.ListAttribute{MarkdownDescription: "Additional CORS headers to allow.", Optional: true, ElementType: types.StringType},
			"cors_allowed_methods": resourceschema.ListAttribute{MarkdownDescription: "CORS methods to allow.", Optional: true, ElementType: types.StringType},
			"cors_allowed_origins": resourceschema.ListAttribute{MarkdownDescription: "CORS origins to allow.", Optional: true, ElementType: types.StringType},
			"cors_enabled":         resourceschema.BoolAttribute{MarkdownDescription: "Enable CORS handling for a custom HTTP server address.", Optional: true},
			"heartbeat":            resourceschema.StringAttribute{MarkdownDescription: "HTTP stream heartbeat interval.", Optional: true},
			"notify_policy":        resourceschema.StringAttribute{MarkdownDescription: "Delivery policy applied when notifying this outlet.", Optional: true},
			"path":                 resourceschema.StringAttribute{MarkdownDescription: "HTTP server path.", Required: true},
			"ping_period":          resourceschema.StringAttribute{MarkdownDescription: "Websocket ping period.", Optional: true},
			"pong_wait":            resourceschema.StringAttribute{MarkdownDescription: "Websocket pong wait timeout.", Optional: true},
			"stream_format":        resourceschema.StringAttribute{MarkdownDescription: "HTTP stream endpoint format.", Optional: true},
			"stream_path":          resourceschema.StringAttribute{MarkdownDescription: "HTTP server stream path.", Optional: true},
			"timeout":              resourceschema.StringAttribute{MarkdownDescription: "HTTP path endpoint timeout.", Optional: true},
			"write_wait":           resourceschema.StringAttribute{MarkdownDescription: "Websocket write timeout.", Optional: true},
			"ws_message_type":      resourceschema.StringAttribute{MarkdownDescription: "Websocket message type.", Optional: true},
			"ws_path":              resourceschema.StringAttribute{MarkdownDescription: "HTTP server websocket path.", Optional: true},
			"apply_status":         resourceschema.StringAttribute{MarkdownDescription: "Observed apply/runtime status for the outlet.", Computed: true},
			"last_apply_error":     resourceschema.StringAttribute{MarkdownDescription: "Last observed apply/runtime error message.", Computed: true},
			"config_generation":    resourceschema.Int64Attribute{MarkdownDescription: "Desired outlet config generation.", Computed: true},
			"applied_generation":   resourceschema.Int64Attribute{MarkdownDescription: "Last observed applied config generation.", Computed: true},
			"transform": resourceschema.SingleNestedAttribute{MarkdownDescription: "Optional outlet transform selection and observed runtime state.", Optional: true, Computed: true, Attributes: map[string]resourceschema.Attribute{
				"selected": resourceschema.StringAttribute{MarkdownDescription: "Explicitly selected transform preset.", Optional: true}, "effective": resourceschema.StringAttribute{MarkdownDescription: "Effective transform after defaults are applied.", Computed: true}, "state": resourceschema.StringAttribute{MarkdownDescription: "Observed transform runtime state.", Computed: true}, "applied": resourceschema.BoolAttribute{MarkdownDescription: "Whether the transform has been applied.", Computed: true}, "generation": resourceschema.Int64Attribute{MarkdownDescription: "Observed transform generation.", Computed: true}, "last_error": resourceschema.StringAttribute{MarkdownDescription: "Last observed transform apply error.", Computed: true},
			}},
		},
		Blocks: map[string]resourceschema.Block{
			"tls": resourceschema.SingleNestedBlock{
				MarkdownDescription: "Optional inline TLS feature for the HTTP server outlet.",
				Attributes: map[string]resourceschema.Attribute{
					"enabled":                     resourceschema.BoolAttribute{MarkdownDescription: "Whether the TLS feature is enabled.", Computed: true},
					"state":                       resourceschema.StringAttribute{MarkdownDescription: "Observed TLS feature state.", Computed: true},
					"desired_generation":          resourceschema.Int64Attribute{MarkdownDescription: "Desired TLS feature generation.", Computed: true},
					"applied_generation":          resourceschema.Int64Attribute{MarkdownDescription: "Applied TLS feature generation.", Computed: true},
					"last_error":                  resourceschema.StringAttribute{MarkdownDescription: "Last observed TLS feature error.", Computed: true},
					"server_cert_filename":        resourceschema.StringAttribute{MarkdownDescription: "Filename to report for the TLS certificate file.", Optional: true},
					"server_cert_content_base64":  resourceschema.StringAttribute{MarkdownDescription: "Base64-encoded TLS certificate file.", Optional: true, Sensitive: true, WriteOnly: true},
					"server_cert_version":         resourceschema.Int64Attribute{MarkdownDescription: "Monotonic version used to rotate the TLS certificate file.", Optional: true},
					"server_cert_checksum_sha256": resourceschema.StringAttribute{MarkdownDescription: "Observed checksum of the uploaded TLS certificate file.", Computed: true},
					"server_key_filename":         resourceschema.StringAttribute{MarkdownDescription: "Filename to report for the TLS private key file.", Optional: true},
					"server_key_content_base64":   resourceschema.StringAttribute{MarkdownDescription: "Base64-encoded TLS private key file.", Optional: true, Sensitive: true, WriteOnly: true},
					"server_key_version":          resourceschema.Int64Attribute{MarkdownDescription: "Monotonic version used to rotate the TLS private key file.", Optional: true},
					"server_key_checksum_sha256":  resourceschema.StringAttribute{MarkdownDescription: "Observed checksum of the uploaded TLS private key file.", Computed: true},
				},
			},
		},
	}
}

func (r *httpServerOutletResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	outletClient, ok := req.ProviderData.(httpServerOutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected HTTP server outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}

	r.client = outletClient
}

func (r *httpServerOutletResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config httpServerOutletResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.TLS.IsNull() || config.TLS.IsUnknown() {
		return
	}

	tls, diags := expandHttpServerTLS(ctx, config.TLS)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if tls.ServerCertFilename.IsNull() || tls.ServerCertFilename.IsUnknown() ||
		tls.ServerCertContentBase64.IsNull() || tls.ServerCertContentBase64.IsUnknown() ||
		tls.ServerCertVersion.IsNull() || tls.ServerCertVersion.IsUnknown() ||
		tls.ServerKeyFilename.IsNull() || tls.ServerKeyFilename.IsUnknown() ||
		tls.ServerKeyContentBase64.IsNull() || tls.ServerKeyContentBase64.IsUnknown() ||
		tls.ServerKeyVersion.IsNull() || tls.ServerKeyVersion.IsUnknown() {
		resp.Diagnostics.AddError("Incomplete HTTP Server TLS Configuration", "When `tls` is configured, `server_cert_filename`, `server_cert_content_base64`, `server_cert_version`, `server_key_filename`, `server_key_content_base64`, and `server_key_version` must all be set.")
	}
}

func (r *httpServerOutletResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var plan httpServerOutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	allowedVerbs, err := stringListValueOrNil(ctx, plan.AllowedVerbs)
	if err != nil {
		resp.Diagnostics.AddError("Invalid HTTP Server Allowed Verbs", err.Error())
		return
	}
	corsAllowedHeaders, err := stringListValueOrNil(ctx, plan.CorsAllowedHeaders)
	if err != nil {
		resp.Diagnostics.AddError("Invalid HTTP Server CORS Allowed Headers", err.Error())
		return
	}
	corsAllowedMethods, err := stringListValueOrNil(ctx, plan.CorsAllowedMethods)
	if err != nil {
		resp.Diagnostics.AddError("Invalid HTTP Server CORS Allowed Methods", err.Error())
		return
	}
	corsAllowedOrigins, err := stringListValueOrNil(ctx, plan.CorsAllowedOrigins)
	if err != nil {
		resp.Diagnostics.AddError("Invalid HTTP Server CORS Allowed Origins", err.Error())
		return
	}
	tls, diags := expandHttpServerTLS(ctx, plan.TLS)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tls, diags = enrichHttpServerTLSFromConfig(ctx, req.Config, tls)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	configuredTransform, transformDiags := selectedTransformFromConfig(ctx, req.Config)
	resp.Diagnostics.Append(transformDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	outlet, err := r.client.CreateHttpServerOutlet(ctx, client.HttpServerOutletCreateRequest{
		Name:               plan.Name.ValueString(),
		Description:        stringValueOrEmpty(plan.Description),
		Enabled:            plan.Enabled.ValueBool(),
		Address:            stringPointerValue(plan.Address),
		AllowedVerbs:       allowedVerbs,
		CorsAllowedHeaders: corsAllowedHeaders,
		CorsAllowedMethods: corsAllowedMethods,
		CorsAllowedOrigins: corsAllowedOrigins,
		CorsEnabled:        boolPointerValue(plan.CorsEnabled),
		Heartbeat:          stringPointerValue(plan.Heartbeat),
		NotifyPolicy:       stringPointerValue(plan.NotifyPolicy),
		Path:               plan.Path.ValueString(),
		PingPeriod:         stringPointerValue(plan.PingPeriod),
		PongWait:           stringPointerValue(plan.PongWait),
		StreamFormat:       stringPointerValue(plan.StreamFormat),
		StreamPath:         stringPointerValue(plan.StreamPath),
		Timeout:            stringPointerValue(plan.Timeout),
		WriteWait:          stringPointerValue(plan.WriteWait),
		WsMessageType:      stringPointerValue(plan.WsMessageType),
		WsPath:             stringPointerValue(plan.WsPath),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create HTTP Server Outlet", err.Error())
		return
	}

	if !plan.TLS.IsNull() {
		err = r.client.ConfigureHttpServerTLSFeature(ctx, outlet.ID, client.HttpServerTLSFeatureRequest{
			ServerCertFilename:      tls.ServerCertFilename.ValueString(),
			ServerCertContentBase64: tls.ServerCertContentBase64.ValueString(),
			ServerKeyFilename:       tls.ServerKeyFilename.ValueString(),
			ServerKeyContentBase64:  tls.ServerKeyContentBase64.ValueString(),
		})
		if err != nil {
			if deleteErr := r.client.DeleteHttpServerOutlet(ctx, outlet.ID); deleteErr != nil && !client.IsNotFound(deleteErr) {
				resp.Diagnostics.AddError("Unable to Configure HTTP Server TLS Feature", fmt.Sprintf("%s; cleanup failed: %s", err, deleteErr))
				return
			}
			resp.Diagnostics.AddError("Unable to Configure HTTP Server TLS Feature", err.Error())
			return
		}
	}
	if !configuredTransform.IsNull() && !configuredTransform.IsUnknown() {
		if _, err := r.client.SetTypedOutletTransform(ctx, "http_server", outlet.ID, configuredTransform.ValueString()); err != nil {
			if deleteErr := r.client.DeleteHttpServerOutlet(ctx, outlet.ID); deleteErr != nil && !client.IsNotFound(deleteErr) {
				resp.Diagnostics.AddError("Unable to Configure HTTP Server Transform", fmt.Sprintf("%s; cleanup failed: %s", err, deleteErr))
				return
			}
			resp.Diagnostics.AddError("Unable to Configure HTTP Server Transform", err.Error())
			return
		}
	}

	newState, removed, diags := r.buildState(ctx, outlet.ID, plan.TLS)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if removed {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "http_server", outlet.ID, plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe HTTP Server Outlet Lifecycle", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *httpServerOutletResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var state httpServerOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, removed, diags := r.buildState(ctx, state.ID.ValueInt64(), state.TLS)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if removed {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *httpServerOutletResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var plan httpServerOutletResourceModel
	var state httpServerOutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	allowedVerbs := mustStringListOrEmpty(ctx, plan.AllowedVerbs)
	corsAllowedHeaders := mustStringListOrEmpty(ctx, plan.CorsAllowedHeaders)
	corsAllowedMethods := mustStringListOrEmpty(ctx, plan.CorsAllowedMethods)
	corsAllowedOrigins := mustStringListOrEmpty(ctx, plan.CorsAllowedOrigins)
	planTLS, diags := expandHttpServerTLS(ctx, plan.TLS)
	resp.Diagnostics.Append(diags...)
	stateTLS, diags := expandHttpServerTLS(ctx, state.TLS)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	planTLS, diags = enrichHttpServerTLSFromConfig(ctx, req.Config, planTLS)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	configuredTransform, transformDiags := selectedTransformFromConfig(ctx, req.Config)
	resp.Diagnostics.Append(transformDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateHttpServerOutlet(ctx, state.ID.ValueInt64(), client.HttpServerOutletUpdateRequest{
		Name:               plan.Name.ValueString(),
		Description:        stringValueOrEmpty(plan.Description),
		Enabled:            plan.Enabled.ValueBool(),
		Address:            stringPointerValue(plan.Address),
		AllowedVerbs:       allowedVerbs,
		CorsAllowedHeaders: corsAllowedHeaders,
		CorsAllowedMethods: corsAllowedMethods,
		CorsAllowedOrigins: corsAllowedOrigins,
		CorsEnabled:        boolPointerValue(plan.CorsEnabled),
		Heartbeat:          stringPointerValue(plan.Heartbeat),
		NotifyPolicy:       stringPointerValue(plan.NotifyPolicy),
		Path:               plan.Path.ValueString(),
		PingPeriod:         stringPointerValue(plan.PingPeriod),
		PongWait:           stringPointerValue(plan.PongWait),
		StreamFormat:       stringPointerValue(plan.StreamFormat),
		StreamPath:         stringPointerValue(plan.StreamPath),
		Timeout:            stringPointerValue(plan.Timeout),
		WriteWait:          stringPointerValue(plan.WriteWait),
		WsMessageType:      stringPointerValue(plan.WsMessageType),
		WsPath:             stringPointerValue(plan.WsPath),
	}); err != nil {
		resp.Diagnostics.AddError("Unable to Update HTTP Server Outlet", err.Error())
		return
	}

	fieldUpdater := httpServerOutletFieldUpdater{client: r.client}
	mutations := []struct {
		name string
		err  error
	}{
		{"address", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.address", plan.Address, state.Address)},
		{"allowed_verbs", syncStringListField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.allowed_verbs", plan.AllowedVerbs, state.AllowedVerbs)},
		{"cors_allowed_headers", syncStringListField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.cors.allowed_headers", plan.CorsAllowedHeaders, state.CorsAllowedHeaders)},
		{"cors_allowed_methods", syncStringListField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.cors.allowed_methods", plan.CorsAllowedMethods, state.CorsAllowedMethods)},
		{"cors_allowed_origins", syncStringListField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.cors.allowed_origins", plan.CorsAllowedOrigins, state.CorsAllowedOrigins)},
		{"cors_enabled", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.cors.enabled", plan.CorsEnabled, state.CorsEnabled)},
		{"heartbeat", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.heartbeat", plan.Heartbeat, state.Heartbeat)},
		{"notify_policy", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "notify_policy", plan.NotifyPolicy, state.NotifyPolicy)},
		{"ping_period", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.ping_period", plan.PingPeriod, state.PingPeriod)},
		{"pong_wait", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.pong_wait", plan.PongWait, state.PongWait)},
		{"stream_format", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.stream_format", plan.StreamFormat, state.StreamFormat)},
		{"stream_path", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.stream_path", plan.StreamPath, state.StreamPath)},
		{"timeout", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.timeout", plan.Timeout, state.Timeout)},
		{"write_wait", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.write_wait", plan.WriteWait, state.WriteWait)},
		{"ws_message_type", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.ws_message_type", plan.WsMessageType, state.WsMessageType)},
		{"ws_path", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "http_server.ws_path", plan.WsPath, state.WsPath)},
	}
	for _, mutation := range mutations {
		if mutation.err != nil {
			resp.Diagnostics.AddError("Unable to Update HTTP Server Field", fmt.Sprintf("%s: %s", mutation.name, mutation.err))
			return
		}
	}

	switch {
	case plan.TLS.IsNull() && !state.TLS.IsNull():
		if err := r.client.DeleteHttpServerTLSFeature(ctx, state.ID.ValueInt64()); err != nil && !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Disable HTTP Server TLS Feature", err.Error())
			return
		}
		for _, role := range []string{httpServerTLSCertRole, httpServerTLSKeyRole} {
			if _, err := r.client.DeleteHttpServerCredential(ctx, state.ID.ValueInt64(), role); err != nil && !client.IsNotFound(err) {
				resp.Diagnostics.AddError("Unable to Delete HTTP Server TLS Credential", fmt.Sprintf("%s: %s", role, err))
				return
			}
		}
	case !plan.TLS.IsNull():
		tlsRequest := client.HttpServerTLSFeatureRequest{
			ServerCertFilename:      planTLS.ServerCertFilename.ValueString(),
			ServerCertContentBase64: planTLS.ServerCertContentBase64.ValueString(),
			ServerKeyFilename:       planTLS.ServerKeyFilename.ValueString(),
			ServerKeyContentBase64:  planTLS.ServerKeyContentBase64.ValueString(),
		}
		if state.TLS.IsNull() || !httpServerTLSEnabled(stateTLS) {
			if err := r.client.ConfigureHttpServerTLSFeature(ctx, state.ID.ValueInt64(), tlsRequest); err != nil {
				resp.Diagnostics.AddError("Unable to Configure HTTP Server TLS Feature", err.Error())
				return
			}
		} else if httpServerTLSChanged(planTLS, stateTLS) {
			if err := r.client.UpdateHttpServerTLSFeature(ctx, state.ID.ValueInt64(), tlsRequest); err != nil {
				resp.Diagnostics.AddError("Unable to Update HTTP Server TLS Feature", err.Error())
				return
			}
		}
	}
	if err := syncTypedOutletTransform(ctx, r.client, "http_server", state.ID.ValueInt64(), configuredTransform, selectedTransform(state.Transform)); err != nil {
		resp.Diagnostics.AddError("Unable to Update HTTP Server Transform", err.Error())
		return
	}

	newState, removed, diags := r.buildState(ctx, state.ID.ValueInt64(), plan.TLS)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if removed {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "http_server", state.ID.ValueInt64(), plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe HTTP Server Outlet Lifecycle", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *httpServerOutletResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var state httpServerOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operation, err := r.client.DeleteHttpServerOutletOperation(ctx, state.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Delete HTTP Server Outlet", err.Error())
		return
	}
	if err == nil && operation.DesiredGeneration != 0 {
		if err := observeTypedOutletLifecycle(ctx, r.client, "http_server", state.ID.ValueInt64(), operation.DesiredGeneration, false, true); err != nil {
			resp.Diagnostics.AddError("Unable to Observe HTTP Server Outlet Lifecycle", err.Error())
		}
	}
}

func (r *httpServerOutletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid HTTP Server Outlet Import Identifier", fmt.Sprintf("Expected numeric outlet ID, got %q: %s", req.ID, err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (r *httpServerOutletResource) buildState(ctx context.Context, outletID int64, previousTLS types.Object) (httpServerOutletResourceModel, bool, diag.Diagnostics) {
	outlet, err := r.client.GetHttpServerOutlet(ctx, outletID)
	if err != nil {
		if client.IsNotFound(err) {
			return httpServerOutletResourceModel{}, true, nil
		}
		var diags diag.Diagnostics
		diags.AddError("Unable to Read HTTP Server Outlet", err.Error())
		return httpServerOutletResourceModel{}, false, diags
	}

	tlsFeature, err := r.client.GetHttpServerTLSFeature(ctx, outletID)
	var tlsFeaturePtr *client.OutletFeatureDetail
	if err != nil {
		if !client.IsNotFound(err) {
			var diags diag.Diagnostics
			diags.AddError("Unable to Read HTTP Server TLS Feature", err.Error())
			return httpServerOutletResourceModel{}, false, diags
		}
	} else {
		tlsFeaturePtr = &tlsFeature
	}

	credentials, err := r.client.ListHttpServerCredentials(ctx, outletID)
	if err != nil {
		if !client.IsNotFound(err) {
			var diags diag.Diagnostics
			diags.AddError("Unable to Read HTTP Server Credentials", err.Error())
			return httpServerOutletResourceModel{}, false, diags
		}
		credentials = nil
	}

	tlsObject, tlsDiags := buildHttpServerTLSObject(ctx, previousTLS, tlsFeaturePtr, credentials)
	if tlsDiags.HasError() {
		return httpServerOutletResourceModel{}, false, tlsDiags
	}

	transform, err := r.client.GetTypedOutletTransform(ctx, "http_server", outletID)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read HTTP Server Transform", err.Error())
		return httpServerOutletResourceModel{}, false, diags
	}

	return httpServerOutletModelFromAPI(outlet, tlsObject, transform), false, nil
}

func httpServerOutletModelFromAPI(outlet client.HttpServerOutlet, tls types.Object, transform client.OutletTransformState) httpServerOutletResourceModel {
	return httpServerOutletResourceModel{
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
		TLS:                tls,
		ApplyStatus:        nullableString(outlet.ApplyStatus),
		LastApplyError:     nullableString(outlet.LastApplyError),
		ConfigGeneration:   types.Int64Value(outlet.ConfigGeneration),
		AppliedGeneration:  types.Int64Value(outlet.AppliedGeneration),
		Transform:          transformObject(transform),
	}
}

func enrichHttpServerTLSFromConfig(ctx context.Context, config tfsdk.Config, tls httpServerTLSModel) (httpServerTLSModel, diag.Diagnostics) {
	if tls.ServerCertContentBase64.IsNull() || tls.ServerCertContentBase64.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root("tls").AtName("server_cert_content_base64"), &value)
		if diags.HasError() {
			return tls, diags
		}
		tls.ServerCertContentBase64 = value
	}

	if tls.ServerKeyContentBase64.IsNull() || tls.ServerKeyContentBase64.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root("tls").AtName("server_key_content_base64"), &value)
		if diags.HasError() {
			return tls, diags
		}
		tls.ServerKeyContentBase64 = value
	}

	return tls, nil
}

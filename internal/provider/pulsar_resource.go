package provider

import (
	"context"
	"errors"
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
	_ resource.Resource                   = &pulsarOutletResource{}
	_ resource.ResourceWithConfigure      = &pulsarOutletResource{}
	_ resource.ResourceWithImportState    = &pulsarOutletResource{}
	_ resource.ResourceWithValidateConfig = &pulsarOutletResource{}
)

type pulsarOutletClient interface {
	CreatePulsarOutlet(context.Context, client.PulsarOutletCreateRequest) (client.PulsarOutlet, error)
	GetPulsarOutlet(context.Context, int64) (client.PulsarOutlet, error)
	UpdatePulsarOutlet(context.Context, int64, client.PulsarOutletUpdateRequest) error
	DeletePulsarOutlet(context.Context, int64) error
	DeletePulsarOutletOperation(context.Context, int64) (client.OutletOperation, error)
	GetTypedOutletStatus(context.Context, string, int64, int64) (client.OutletStatus, error)
	SetPulsarField(context.Context, int64, string, any) error
	DeletePulsarField(context.Context, int64, string) error
	ListPulsarCredentials(context.Context, int64) ([]client.OutletMaterialSummary, error)
	DeletePulsarCredential(context.Context, int64, string) (client.OutletMaterialDeleteResponse, error)
	GetPulsarTLSFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigurePulsarTLSFeature(context.Context, int64, client.PulsarTLSFeatureRequest) error
	UpdatePulsarTLSFeature(context.Context, int64, client.PulsarTLSFeatureRequest) error
	DeletePulsarTLSFeature(context.Context, int64) error
	GetPulsarOAuth2Feature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigurePulsarOAuth2Feature(context.Context, int64, client.PulsarOAuth2FeatureRequest) error
	UpdatePulsarOAuth2Feature(context.Context, int64, client.PulsarOAuth2FeatureRequest) error
	DeletePulsarOAuth2Feature(context.Context, int64) error
	GetPulsarTokenFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigurePulsarTokenFeature(context.Context, int64, client.PulsarTokenFeatureRequest) error
	UpdatePulsarTokenFeature(context.Context, int64, client.PulsarTokenFeatureRequest) error
	DeletePulsarTokenFeature(context.Context, int64) error
	typedOutletTransformClient
}

type pulsarOutletFieldUpdater struct {
	client pulsarOutletClient
}

func (u pulsarOutletFieldUpdater) SetField(ctx context.Context, outletID int64, field string, value any) error {
	return u.client.SetPulsarField(ctx, outletID, field, value)
}

func (u pulsarOutletFieldUpdater) DeleteField(ctx context.Context, outletID int64, field string) error {
	return u.client.DeletePulsarField(ctx, outletID, field)
}

type pulsarOutletResource struct {
	client pulsarOutletClient
}

type pulsarOutletResourceModel struct {
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	Key               types.String `tfsdk:"key"`
	MaxInFlight       types.Int64  `tfsdk:"max_in_flight"`
	NotifyPolicy      types.String `tfsdk:"notify_policy"`
	OrderingKey       types.String `tfsdk:"ordering_key"`
	Topic             types.String `tfsdk:"topic"`
	URL               types.String `tfsdk:"url"`
	TLS               types.Object `tfsdk:"tls"`
	OAuth2            types.Object `tfsdk:"oauth2"`
	Token             types.Object `tfsdk:"token"`
	ApplyStatus       types.String `tfsdk:"apply_status"`
	LastApplyError    types.String `tfsdk:"last_apply_error"`
	ConfigGeneration  types.Int64  `tfsdk:"config_generation"`
	AppliedGeneration types.Int64  `tfsdk:"applied_generation"`
	Transform         types.Object `tfsdk:"transform"`
}

func NewPulsarOutletResource() resource.Resource { return &pulsarOutletResource{} }

func (r *pulsarOutletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pulsar_outlet"
}

func (r *pulsarOutletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		MarkdownDescription: "Manages a Noozle Pulsar outlet.",
		Attributes: map[string]resourceschema.Attribute{
			"id":                 resourceschema.Int64Attribute{MarkdownDescription: "Unique identifier of the outlet.", Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}},
			"name":               resourceschema.StringAttribute{MarkdownDescription: "Name presented for the outlet.", Required: true},
			"description":        resourceschema.StringAttribute{MarkdownDescription: "Optional note explaining the outlet purpose.", Optional: true},
			"enabled":            resourceschema.BoolAttribute{MarkdownDescription: "Whether the outlet is active.", Required: true},
			"key":                resourceschema.StringAttribute{MarkdownDescription: "Pulsar message key interpolation.", Optional: true},
			"max_in_flight":      resourceschema.Int64Attribute{MarkdownDescription: "Maximum number of in-flight Pulsar messages.", Optional: true},
			"notify_policy":      resourceschema.StringAttribute{MarkdownDescription: "Delivery policy applied when notifying this outlet.", Optional: true},
			"ordering_key":       resourceschema.StringAttribute{MarkdownDescription: "Pulsar ordering key interpolation.", Optional: true},
			"topic":              resourceschema.StringAttribute{MarkdownDescription: "Pulsar topic.", Required: true},
			"url":                resourceschema.StringAttribute{MarkdownDescription: "Pulsar broker URL.", Required: true},
			"apply_status":       resourceschema.StringAttribute{MarkdownDescription: "Observed apply/runtime status for the outlet.", Computed: true},
			"last_apply_error":   resourceschema.StringAttribute{MarkdownDescription: "Last observed apply/runtime error message.", Computed: true},
			"config_generation":  resourceschema.Int64Attribute{MarkdownDescription: "Desired outlet config generation.", Computed: true},
			"applied_generation": resourceschema.Int64Attribute{MarkdownDescription: "Last observed applied config generation.", Computed: true},
			"transform": resourceschema.SingleNestedAttribute{MarkdownDescription: "Optional outlet transform selection and observed runtime state.", Optional: true, Computed: true, Attributes: map[string]resourceschema.Attribute{
				"selected": resourceschema.StringAttribute{MarkdownDescription: "Explicitly selected transform preset.", Optional: true}, "effective": resourceschema.StringAttribute{MarkdownDescription: "Effective transform after defaults are applied.", Computed: true}, "state": resourceschema.StringAttribute{MarkdownDescription: "Observed transform runtime state.", Computed: true}, "applied": resourceschema.BoolAttribute{MarkdownDescription: "Whether the transform has been applied.", Computed: true}, "generation": resourceschema.Int64Attribute{MarkdownDescription: "Observed transform generation.", Computed: true}, "last_error": resourceschema.StringAttribute{MarkdownDescription: "Last observed transform apply error.", Computed: true},
			}},
		},
		Blocks: map[string]resourceschema.Block{
			"tls": resourceschema.SingleNestedBlock{
				MarkdownDescription: "Optional inline TLS feature for the Pulsar outlet.",
				Attributes: map[string]resourceschema.Attribute{
					"enabled":                 resourceschema.BoolAttribute{MarkdownDescription: "Whether the TLS feature is enabled.", Computed: true},
					"state":                   resourceschema.StringAttribute{MarkdownDescription: "Observed TLS feature state.", Computed: true},
					"desired_generation":      resourceschema.Int64Attribute{MarkdownDescription: "Desired TLS feature generation.", Computed: true},
					"applied_generation":      resourceschema.Int64Attribute{MarkdownDescription: "Applied TLS feature generation.", Computed: true},
					"last_error":              resourceschema.StringAttribute{MarkdownDescription: "Last observed TLS feature error.", Computed: true},
					"ca_cert_filename":        resourceschema.StringAttribute{MarkdownDescription: "Filename to report for the TLS CA certificate file.", Optional: true},
					"ca_cert_content_base64":  resourceschema.StringAttribute{MarkdownDescription: "Base64-encoded TLS CA certificate file.", Optional: true, Sensitive: true, WriteOnly: true},
					"ca_cert_version":         resourceschema.Int64Attribute{MarkdownDescription: "Monotonic version used to rotate the TLS CA certificate file.", Optional: true},
					"ca_cert_checksum_sha256": resourceschema.StringAttribute{MarkdownDescription: "Observed checksum of the uploaded TLS CA certificate file.", Computed: true},
				},
			},
			"oauth2": resourceschema.SingleNestedBlock{
				MarkdownDescription: "Optional inline OAuth2 feature for the Pulsar outlet. Conflicts with `token`.",
				Attributes: map[string]resourceschema.Attribute{
					"enabled":                     resourceschema.BoolAttribute{MarkdownDescription: "Whether the OAuth2 feature is enabled.", Computed: true},
					"state":                       resourceschema.StringAttribute{MarkdownDescription: "Observed OAuth2 feature state.", Computed: true},
					"desired_generation":          resourceschema.Int64Attribute{MarkdownDescription: "Desired OAuth2 feature generation.", Computed: true},
					"applied_generation":          resourceschema.Int64Attribute{MarkdownDescription: "Applied OAuth2 feature generation.", Computed: true},
					"last_error":                  resourceschema.StringAttribute{MarkdownDescription: "Last observed OAuth2 feature error.", Computed: true},
					"issuer_url":                  resourceschema.StringAttribute{MarkdownDescription: "OAuth2 issuer URL.", Optional: true},
					"audience":                    resourceschema.StringAttribute{MarkdownDescription: "OAuth2 audience.", Optional: true},
					"private_key_filename":        resourceschema.StringAttribute{MarkdownDescription: "Filename to report for the OAuth2 private key file.", Optional: true},
					"private_key_content_base64":  resourceschema.StringAttribute{MarkdownDescription: "Base64-encoded OAuth2 private key file.", Optional: true, Sensitive: true, WriteOnly: true},
					"private_key_version":         resourceschema.Int64Attribute{MarkdownDescription: "Monotonic version used to rotate the OAuth2 private key file.", Optional: true},
					"private_key_checksum_sha256": resourceschema.StringAttribute{MarkdownDescription: "Observed checksum of the uploaded OAuth2 private key file.", Computed: true},
				},
			},
			"token": resourceschema.SingleNestedBlock{
				MarkdownDescription: "Optional inline token feature for the Pulsar outlet. Conflicts with `oauth2`.",
				Attributes: map[string]resourceschema.Attribute{
					"enabled":              resourceschema.BoolAttribute{MarkdownDescription: "Whether the token feature is enabled.", Computed: true},
					"state":                resourceschema.StringAttribute{MarkdownDescription: "Observed token feature state.", Computed: true},
					"desired_generation":   resourceschema.Int64Attribute{MarkdownDescription: "Desired token feature generation.", Computed: true},
					"applied_generation":   resourceschema.Int64Attribute{MarkdownDescription: "Applied token feature generation.", Computed: true},
					"last_error":           resourceschema.StringAttribute{MarkdownDescription: "Last observed token feature error.", Computed: true},
					"access_token":         resourceschema.StringAttribute{MarkdownDescription: "Pulsar access token.", Optional: true, Sensitive: true, WriteOnly: true},
					"access_token_version": resourceschema.Int64Attribute{MarkdownDescription: "Monotonic version used to rotate the Pulsar access token.", Optional: true},
				},
			},
		},
	}
}

func (r *pulsarOutletResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	outletClient, ok := req.ProviderData.(pulsarOutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected Pulsar outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}
	r.client = outletClient
}

func (r *pulsarOutletResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config pulsarOutletResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.OAuth2.IsNull() && !config.OAuth2.IsUnknown() && !config.Token.IsNull() && !config.Token.IsUnknown() {
		resp.Diagnostics.AddError("Conflicting Pulsar Features", "Only one of `oauth2` or `token` can be configured at the same time. Both blocks are optional, but they are mutually exclusive.")
	}

	if !config.TLS.IsNull() && !config.TLS.IsUnknown() {
		tls, diags := expandPulsarTLS(ctx, config.TLS)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if tls.CACertFilename.IsNull() || tls.CACertFilename.IsUnknown() || tls.CACertContentBase64.IsNull() || tls.CACertContentBase64.IsUnknown() || tls.CACertVersion.IsNull() || tls.CACertVersion.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete Pulsar TLS Configuration", "When `tls` is configured, `ca_cert_filename`, `ca_cert_content_base64`, and `ca_cert_version` must all be set.")
		}
	}

	if !config.OAuth2.IsNull() && !config.OAuth2.IsUnknown() {
		oauth2, diags := expandPulsarOAuth2(ctx, config.OAuth2)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if oauth2.IssuerURL.IsNull() || oauth2.IssuerURL.IsUnknown() ||
			oauth2.Audience.IsNull() || oauth2.Audience.IsUnknown() ||
			oauth2.PrivateKeyFilename.IsNull() || oauth2.PrivateKeyFilename.IsUnknown() ||
			oauth2.PrivateKeyContentBase64.IsNull() || oauth2.PrivateKeyContentBase64.IsUnknown() ||
			oauth2.PrivateKeyVersion.IsNull() || oauth2.PrivateKeyVersion.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete Pulsar OAuth2 Configuration", "When `oauth2` is configured, `issuer_url`, `audience`, `private_key_filename`, `private_key_content_base64`, and `private_key_version` must all be set.")
		}
	}

	if !config.Token.IsNull() && !config.Token.IsUnknown() {
		token, diags := expandPulsarToken(ctx, config.Token)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if token.AccessToken.IsNull() || token.AccessToken.IsUnknown() || token.AccessTokenVersion.IsNull() || token.AccessTokenVersion.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete Pulsar Token Configuration", "When `token` is configured, `access_token` and `access_token_version` must both be set.")
		}
	}
}

func (r *pulsarOutletResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var plan pulsarOutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tls, oauth2, token, ok := r.readFeaturesFromConfig(ctx, req.Config, plan, &resp.Diagnostics)
	if !ok {
		return
	}
	configuredTransform, transformDiags := selectedTransformFromConfig(ctx, req.Config)
	resp.Diagnostics.Append(transformDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	outlet, err := r.client.CreatePulsarOutlet(ctx, client.PulsarOutletCreateRequest{
		Name:         plan.Name.ValueString(),
		Description:  stringValueOrEmpty(plan.Description),
		Enabled:      plan.Enabled.ValueBool(),
		Key:          stringPointerValue(plan.Key),
		MaxInFlight:  int64PointerValue(plan.MaxInFlight),
		NotifyPolicy: stringPointerValue(plan.NotifyPolicy),
		OrderingKey:  stringPointerValue(plan.OrderingKey),
		Topic:        plan.Topic.ValueString(),
		URL:          plan.URL.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Pulsar Outlet", err.Error())
		return
	}
	if err := r.applyFeatureConfig(ctx, outlet.ID, types.ObjectNull(pulsarTLSAttributeTypes), types.ObjectNull(pulsarOAuth2AttributeTypes), types.ObjectNull(pulsarTokenAttributeTypes), tls, oauth2, token); err != nil {
		if deleteErr := r.client.DeletePulsarOutlet(ctx, outlet.ID); deleteErr != nil && !client.IsNotFound(deleteErr) {
			resp.Diagnostics.AddError("Unable to Configure Pulsar Features", fmt.Sprintf("%s; cleanup failed: %s", err, deleteErr))
			return
		}
		resp.Diagnostics.AddError("Unable to Configure Pulsar Features", err.Error())
		return
	}
	if !configuredTransform.IsNull() && !configuredTransform.IsUnknown() {
		if _, err := r.client.SetTypedOutletTransform(ctx, "pulsar", outlet.ID, configuredTransform.ValueString()); err != nil {
			if deleteErr := r.client.DeletePulsarOutlet(ctx, outlet.ID); deleteErr != nil && !client.IsNotFound(deleteErr) {
				resp.Diagnostics.AddError("Unable to Configure Pulsar Transform", fmt.Sprintf("%s; cleanup failed: %s", err, deleteErr))
				return
			}
			resp.Diagnostics.AddError("Unable to Configure Pulsar Transform", err.Error())
			return
		}
	}
	newState, removed, diags := r.buildState(ctx, outlet.ID, plan.TLS, plan.OAuth2, plan.Token)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if removed {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "pulsar", outlet.ID, plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe Pulsar Outlet Lifecycle", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *pulsarOutletResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var state pulsarOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	newState, removed, diags := r.buildState(ctx, state.ID.ValueInt64(), state.TLS, state.OAuth2, state.Token)
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

func (r *pulsarOutletResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var plan pulsarOutletResourceModel
	var state pulsarOutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tls, oauth2, token, ok := r.readFeaturesFromConfig(ctx, req.Config, plan, &resp.Diagnostics)
	if !ok {
		return
	}
	configuredTransform, transformDiags := selectedTransformFromConfig(ctx, req.Config)
	resp.Diagnostics.Append(transformDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.UpdatePulsarOutlet(ctx, state.ID.ValueInt64(), client.PulsarOutletUpdateRequest{
		Name:         plan.Name.ValueString(),
		Description:  stringValueOrEmpty(plan.Description),
		Enabled:      plan.Enabled.ValueBool(),
		Key:          stringPointerValue(plan.Key),
		MaxInFlight:  int64PointerValue(plan.MaxInFlight),
		NotifyPolicy: stringPointerValue(plan.NotifyPolicy),
		OrderingKey:  stringPointerValue(plan.OrderingKey),
		Topic:        plan.Topic.ValueString(),
		URL:          plan.URL.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Unable to Update Pulsar Outlet", err.Error())
		return
	}
	fieldUpdater := pulsarOutletFieldUpdater{client: r.client}
	mutations := []struct {
		name string
		err  error
	}{
		{"key", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "pulsar.key", plan.Key, state.Key)},
		{"max_in_flight", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "pulsar.max_in_flight", plan.MaxInFlight, state.MaxInFlight)},
		{"notify_policy", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "notify_policy", plan.NotifyPolicy, state.NotifyPolicy)},
		{"ordering_key", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "pulsar.ordering_key", plan.OrderingKey, state.OrderingKey)},
	}
	for _, mutation := range mutations {
		if mutation.err != nil {
			resp.Diagnostics.AddError("Unable to Update Pulsar Field", fmt.Sprintf("%s: %s", mutation.name, mutation.err))
			return
		}
	}
	if err := r.applyFeatureConfig(ctx, state.ID.ValueInt64(), state.TLS, state.OAuth2, state.Token, tls, oauth2, token); err != nil {
		resp.Diagnostics.AddError("Unable to Update Pulsar Features", err.Error())
		return
	}
	if err := syncTypedOutletTransform(ctx, r.client, "pulsar", state.ID.ValueInt64(), configuredTransform, selectedTransform(state.Transform)); err != nil {
		resp.Diagnostics.AddError("Unable to Update Pulsar Transform", err.Error())
		return
	}
	newState, removed, diags := r.buildState(ctx, state.ID.ValueInt64(), plan.TLS, plan.OAuth2, plan.Token)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if removed {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "pulsar", state.ID.ValueInt64(), plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe Pulsar Outlet Lifecycle", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *pulsarOutletResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var state pulsarOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	operation, err := r.client.DeletePulsarOutletOperation(ctx, state.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Delete Pulsar Outlet", err.Error())
		return
	}
	if err == nil && operation.DesiredGeneration != 0 {
		if err := observeTypedOutletLifecycle(ctx, r.client, "pulsar", state.ID.ValueInt64(), operation.DesiredGeneration, false, true); err != nil {
			resp.Diagnostics.AddError("Unable to Observe Pulsar Outlet Lifecycle", err.Error())
		}
	}
}

func (r *pulsarOutletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Pulsar Outlet Import Identifier", fmt.Sprintf("Expected numeric outlet ID, got %q: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (r *pulsarOutletResource) readFeaturesFromConfig(ctx context.Context, cfg tfsdk.Config, plan pulsarOutletResourceModel, diags *diag.Diagnostics) (pulsarTLSModel, pulsarOAuth2Model, pulsarTokenModel, bool) {
	tls, tlsDiags := expandPulsarTLS(ctx, plan.TLS)
	diags.Append(tlsDiags...)
	oauth2, oauth2Diags := expandPulsarOAuth2(ctx, plan.OAuth2)
	diags.Append(oauth2Diags...)
	token, tokenDiags := expandPulsarToken(ctx, plan.Token)
	diags.Append(tokenDiags...)
	if diags.HasError() {
		return tls, oauth2, token, false
	}
	tls, tlsDiags = enrichPulsarTLSFromConfig(ctx, cfg, tls)
	diags.Append(tlsDiags...)
	oauth2, oauth2Diags = enrichPulsarOAuth2FromConfig(ctx, cfg, oauth2)
	diags.Append(oauth2Diags...)
	token, tokenDiags = enrichPulsarTokenFromConfig(ctx, cfg, token)
	diags.Append(tokenDiags...)
	if diags.HasError() {
		return tls, oauth2, token, false
	}
	if !plan.OAuth2.IsNull() && !plan.Token.IsNull() {
		diags.AddError("Conflicting Pulsar Features", "Only one of `oauth2` or `token` can be configured at the same time.")
		return tls, oauth2, token, false
	}
	return tls, oauth2, token, true
}

func (r *pulsarOutletResource) applyFeatureConfig(ctx context.Context, outletID int64, stateTLS, stateOAuth2, stateToken types.Object, planTLS pulsarTLSModel, planOAuth2 pulsarOAuth2Model, planToken pulsarTokenModel) error {
	if !stateOAuth2.IsNull() && !pulsarFeatureEnabled(stateOAuth2) {
		stateOAuth2 = types.ObjectNull(pulsarOAuth2AttributeTypes)
	}
	if !stateToken.IsNull() && !pulsarFeatureEnabled(stateToken) {
		stateToken = types.ObjectNull(pulsarTokenAttributeTypes)
	}
	if !stateTLS.IsNull() && !pulsarFeatureEnabled(stateTLS) {
		stateTLS = types.ObjectNull(pulsarTLSAttributeTypes)
	}

	if !stateOAuth2.IsNull() && !stateOAuth2.IsUnknown() && !planToken.AccessTokenVersion.IsNull() {
		if err := r.client.DeletePulsarOAuth2Feature(ctx, outletID); err != nil && !client.IsNotFound(err) {
			return err
		}
		if _, err := r.client.DeletePulsarCredential(ctx, outletID, pulsarOAuth2PrivateKeyRole); err != nil && !client.IsNotFound(err) {
			return err
		}
		stateOAuth2 = types.ObjectNull(pulsarOAuth2AttributeTypes)
	}
	if !stateToken.IsNull() && !stateToken.IsUnknown() && !planOAuth2.PrivateKeyVersion.IsNull() {
		if err := r.client.DeletePulsarTokenFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
			return err
		}
		if _, err := r.client.DeletePulsarCredential(ctx, outletID, pulsarTokenRole); err != nil && !client.IsNotFound(err) {
			return err
		}
		stateToken = types.ObjectNull(pulsarTokenAttributeTypes)
	}

	if planTLS.CACertVersion.IsNull() {
		if !stateTLS.IsNull() && !stateTLS.IsUnknown() {
			if err := r.client.DeletePulsarTLSFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
				return err
			}
			if _, err := r.client.DeletePulsarCredential(ctx, outletID, pulsarTLSCARole); err != nil && !client.IsNotFound(err) {
				return err
			}
		}
	} else if stateTLS.IsNull() || stateTLS.IsUnknown() {
		if err := r.client.ConfigurePulsarTLSFeature(ctx, outletID, client.PulsarTLSFeatureRequest{
			CACertFilename:      planTLS.CACertFilename.ValueString(),
			CACertContentBase64: planTLS.CACertContentBase64.ValueString(),
		}); err != nil {
			return err
		}
	} else {
		stateTLSModel, diags := expandPulsarTLS(ctx, stateTLS)
		if diags.HasError() {
			return errors.New(diags.Errors()[0].Detail())
		}
		if !planTLS.CACertVersion.Equal(stateTLSModel.CACertVersion) || !planTLS.CACertFilename.Equal(stateTLSModel.CACertFilename) {
			if err := r.client.UpdatePulsarTLSFeature(ctx, outletID, client.PulsarTLSFeatureRequest{
				CACertFilename:      planTLS.CACertFilename.ValueString(),
				CACertContentBase64: planTLS.CACertContentBase64.ValueString(),
			}); err != nil {
				return err
			}
		}
	}

	if planOAuth2.PrivateKeyVersion.IsNull() {
		if !stateOAuth2.IsNull() && !stateOAuth2.IsUnknown() {
			if err := r.client.DeletePulsarOAuth2Feature(ctx, outletID); err != nil && !client.IsNotFound(err) {
				return err
			}
			if _, err := r.client.DeletePulsarCredential(ctx, outletID, pulsarOAuth2PrivateKeyRole); err != nil && !client.IsNotFound(err) {
				return err
			}
		}
	} else if stateOAuth2.IsNull() || stateOAuth2.IsUnknown() {
		if err := r.client.ConfigurePulsarOAuth2Feature(ctx, outletID, client.PulsarOAuth2FeatureRequest{
			IssuerURL:               planOAuth2.IssuerURL.ValueString(),
			Audience:                planOAuth2.Audience.ValueString(),
			PrivateKeyFilename:      planOAuth2.PrivateKeyFilename.ValueString(),
			PrivateKeyContentBase64: planOAuth2.PrivateKeyContentBase64.ValueString(),
		}); err != nil {
			return err
		}
	} else {
		stateOAuth2Model, diags := expandPulsarOAuth2(ctx, stateOAuth2)
		if diags.HasError() {
			return errors.New(diags.Errors()[0].Detail())
		}
		if !planOAuth2.PrivateKeyVersion.Equal(stateOAuth2Model.PrivateKeyVersion) || !planOAuth2.PrivateKeyFilename.Equal(stateOAuth2Model.PrivateKeyFilename) || !planOAuth2.IssuerURL.Equal(stateOAuth2Model.IssuerURL) || !planOAuth2.Audience.Equal(stateOAuth2Model.Audience) {
			if err := r.client.UpdatePulsarOAuth2Feature(ctx, outletID, client.PulsarOAuth2FeatureRequest{
				IssuerURL:               planOAuth2.IssuerURL.ValueString(),
				Audience:                planOAuth2.Audience.ValueString(),
				PrivateKeyFilename:      planOAuth2.PrivateKeyFilename.ValueString(),
				PrivateKeyContentBase64: planOAuth2.PrivateKeyContentBase64.ValueString(),
			}); err != nil {
				return err
			}
		}
	}

	if planToken.AccessTokenVersion.IsNull() {
		if !stateToken.IsNull() && !stateToken.IsUnknown() {
			if err := r.client.DeletePulsarTokenFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
				return err
			}
			if _, err := r.client.DeletePulsarCredential(ctx, outletID, pulsarTokenRole); err != nil && !client.IsNotFound(err) {
				return err
			}
		}
	} else if stateToken.IsNull() || stateToken.IsUnknown() {
		if err := r.client.ConfigurePulsarTokenFeature(ctx, outletID, client.PulsarTokenFeatureRequest{AccessToken: planToken.AccessToken.ValueString()}); err != nil {
			return err
		}
	} else {
		stateTokenModel, diags := expandPulsarToken(ctx, stateToken)
		if diags.HasError() {
			return errors.New(diags.Errors()[0].Detail())
		}
		if !planToken.AccessTokenVersion.Equal(stateTokenModel.AccessTokenVersion) {
			if err := r.client.UpdatePulsarTokenFeature(ctx, outletID, client.PulsarTokenFeatureRequest{AccessToken: planToken.AccessToken.ValueString()}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *pulsarOutletResource) buildState(ctx context.Context, outletID int64, previousTLS, previousOAuth2, previousToken types.Object) (pulsarOutletResourceModel, bool, diag.Diagnostics) {
	outlet, err := r.client.GetPulsarOutlet(ctx, outletID)
	if err != nil {
		if client.IsNotFound(err) {
			return pulsarOutletResourceModel{}, true, nil
		}
		var diags diag.Diagnostics
		diags.AddError("Unable to Read Pulsar Outlet", err.Error())
		return pulsarOutletResourceModel{}, false, diags
	}
	credentials, err := r.client.ListPulsarCredentials(ctx, outletID)
	if err != nil {
		if !client.IsNotFound(err) {
			var diags diag.Diagnostics
			diags.AddError("Unable to Read Pulsar Credentials", err.Error())
			return pulsarOutletResourceModel{}, false, diags
		}
		credentials = nil
	}
	tlsFeature, err := r.client.GetPulsarTLSFeature(ctx, outletID)
	var tlsFeaturePtr *client.OutletFeatureDetail
	if err != nil && !client.IsNotFound(err) {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read Pulsar TLS Feature", err.Error())
		return pulsarOutletResourceModel{}, false, diags
	}
	if err == nil && tlsFeature.Feature != "" {
		tlsFeaturePtr = &tlsFeature
	}
	oauth2Feature, err := r.client.GetPulsarOAuth2Feature(ctx, outletID)
	var oauth2FeaturePtr *client.OutletFeatureDetail
	if err != nil && !client.IsNotFound(err) {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read Pulsar OAuth2 Feature", err.Error())
		return pulsarOutletResourceModel{}, false, diags
	}
	if err == nil && oauth2Feature.Feature != "" {
		oauth2FeaturePtr = &oauth2Feature
	}
	tokenFeature, err := r.client.GetPulsarTokenFeature(ctx, outletID)
	var tokenFeaturePtr *client.OutletFeatureDetail
	if err != nil && !client.IsNotFound(err) {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read Pulsar Token Feature", err.Error())
		return pulsarOutletResourceModel{}, false, diags
	}
	if err == nil && tokenFeature.Feature != "" {
		tokenFeaturePtr = &tokenFeature
	}
	tlsObject, diags := buildPulsarTLSObject(ctx, previousTLS, tlsFeaturePtr, credentials)
	if diags.HasError() {
		return pulsarOutletResourceModel{}, false, diags
	}
	oauth2Object, diags := buildPulsarOAuth2Object(ctx, previousOAuth2, oauth2FeaturePtr, credentials)
	if diags.HasError() {
		return pulsarOutletResourceModel{}, false, diags
	}
	tokenObject, diags := buildPulsarTokenObject(ctx, previousToken, tokenFeaturePtr, credentials)
	if diags.HasError() {
		return pulsarOutletResourceModel{}, false, diags
	}
	transform, err := r.client.GetTypedOutletTransform(ctx, "pulsar", outletID)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read Pulsar Transform", err.Error())
		return pulsarOutletResourceModel{}, false, diags
	}
	return pulsarOutletModelFromAPI(outlet, tlsObject, oauth2Object, tokenObject, transform), false, nil
}

func pulsarOutletModelFromAPI(outlet client.PulsarOutlet, tls, oauth2, token types.Object, transform client.OutletTransformState) pulsarOutletResourceModel {
	return pulsarOutletResourceModel{
		ID:                types.Int64Value(outlet.ID),
		Name:              types.StringValue(outlet.Name),
		Description:       nullableString(outlet.Description),
		Enabled:           types.BoolValue(outlet.Enabled),
		Key:               nullableStringPointer(outlet.Key),
		MaxInFlight:       nullableInt64Pointer(outlet.MaxInFlight),
		NotifyPolicy:      nullableStringPointerPreservingEmpty(outlet.NotifyPolicy),
		OrderingKey:       nullableStringPointer(outlet.OrderingKey),
		Topic:             types.StringValue(outlet.Topic),
		URL:               types.StringValue(outlet.URL),
		TLS:               tls,
		OAuth2:            oauth2,
		Token:             token,
		ApplyStatus:       nullableString(outlet.ApplyStatus),
		LastApplyError:    nullableString(outlet.LastApplyError),
		ConfigGeneration:  types.Int64Value(outlet.ConfigGeneration),
		AppliedGeneration: types.Int64Value(outlet.AppliedGeneration),
		Transform:         transformObject(transform),
	}
}

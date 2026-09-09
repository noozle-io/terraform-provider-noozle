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
	_ resource.Resource                   = &amqp1OutletResource{}
	_ resource.ResourceWithConfigure      = &amqp1OutletResource{}
	_ resource.ResourceWithImportState    = &amqp1OutletResource{}
	_ resource.ResourceWithValidateConfig = &amqp1OutletResource{}
)

type amqp1OutletClient interface {
	CreateAmqp1Outlet(context.Context, client.Amqp1OutletCreateRequest) (client.Amqp1Outlet, error)
	GetAmqp1Outlet(context.Context, int64) (client.Amqp1Outlet, error)
	UpdateAmqp1Outlet(context.Context, int64, client.Amqp1OutletUpdateRequest) error
	DeleteAmqp1Outlet(context.Context, int64) error
	DeleteAmqp1OutletOperation(context.Context, int64) (client.OutletOperation, error)
	GetTypedOutletStatus(context.Context, string, int64, int64) (client.OutletStatus, error)
	SetAmqp1Field(context.Context, int64, string, any) error
	DeleteAmqp1Field(context.Context, int64, string) error
	ListAmqp1Credentials(context.Context, int64) ([]client.OutletMaterialSummary, error)
	DeleteAmqp1Credential(context.Context, int64, string) (client.OutletMaterialDeleteResponse, error)
	GetAmqp1TLSFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureAmqp1TLSFeature(context.Context, int64, client.AmqpTLSFeatureRequest) error
	UpdateAmqp1TLSFeature(context.Context, int64, client.AmqpTLSFeatureRequest) error
	DeleteAmqp1TLSFeature(context.Context, int64) error
	GetAmqp1MTLSFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureAmqp1MTLSFeature(context.Context, int64, client.AmqpMTLSFeatureRequest) error
	UpdateAmqp1MTLSFeature(context.Context, int64, client.AmqpMTLSFeatureRequest) error
	DeleteAmqp1MTLSFeature(context.Context, int64) error
	GetAmqp1SASLAnonymousFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureAmqp1SASLAnonymousFeature(context.Context, int64) error
	UpdateAmqp1SASLAnonymousFeature(context.Context, int64) error
	DeleteAmqp1SASLAnonymousFeature(context.Context, int64) error
	GetAmqp1SASLPlainFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureAmqp1SASLPlainFeature(context.Context, int64, client.Amqp1SASLPlainFeatureRequest) error
	UpdateAmqp1SASLPlainFeature(context.Context, int64, client.Amqp1SASLPlainFeatureRequest) error
	DeleteAmqp1SASLPlainFeature(context.Context, int64) error
	typedOutletTransformClient
}

type amqp1OutletFieldUpdater struct{ client amqp1OutletClient }

func (u amqp1OutletFieldUpdater) SetField(ctx context.Context, outletID int64, field string, value any) error {
	return u.client.SetAmqp1Field(ctx, outletID, field, value)
}

func (u amqp1OutletFieldUpdater) DeleteField(ctx context.Context, outletID int64, field string) error {
	return u.client.DeleteAmqp1Field(ctx, outletID, field)
}

type amqp1OutletResource struct{ client amqp1OutletClient }

type amqp1OutletResourceModel struct {
	ID                       types.Int64  `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Description              types.String `tfsdk:"description"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	ApplicationPropertiesMap types.String `tfsdk:"application_properties_map"`
	MaxInFlight              types.Int64  `tfsdk:"max_in_flight"`
	MetadataExcludePrefixes  types.List   `tfsdk:"metadata_exclude_prefixes"`
	NotifyPolicy             types.String `tfsdk:"notify_policy"`
	TargetAddress            types.String `tfsdk:"target_address"`
	TLSEnableRenegotiation   types.Bool   `tfsdk:"tls_enable_renegotiation"`
	TLSSkipCertVerify        types.Bool   `tfsdk:"tls_skip_cert_verify"`
	Credentials              types.Object `tfsdk:"credentials"`
	TLS                      types.Object `tfsdk:"tls"`
	MTLS                     types.Object `tfsdk:"mtls"`
	SASLAnonymous            types.Object `tfsdk:"sasl_anonymous"`
	SASLPlain                types.Object `tfsdk:"sasl_plain"`
	ApplyStatus              types.String `tfsdk:"apply_status"`
	LastApplyError           types.String `tfsdk:"last_apply_error"`
	ConfigGeneration         types.Int64  `tfsdk:"config_generation"`
	AppliedGeneration        types.Int64  `tfsdk:"applied_generation"`
	Transform                types.Object `tfsdk:"transform"`
}

func NewAmqp1OutletResource() resource.Resource { return &amqp1OutletResource{} }

func (r *amqp1OutletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_amqp_1_outlet"
}

func (r *amqp1OutletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		MarkdownDescription: "Manages a Noozle AMQP 1.0 outlet.",
		Attributes: map[string]resourceschema.Attribute{
			"id":                         resourceschema.Int64Attribute{Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}},
			"name":                       resourceschema.StringAttribute{Required: true},
			"description":                resourceschema.StringAttribute{Optional: true},
			"enabled":                    resourceschema.BoolAttribute{Required: true},
			"application_properties_map": resourceschema.StringAttribute{Optional: true},
			"max_in_flight":              resourceschema.Int64Attribute{Optional: true},
			"metadata_exclude_prefixes":  resourceschema.ListAttribute{Optional: true, ElementType: types.StringType},
			"notify_policy":              resourceschema.StringAttribute{Optional: true},
			"target_address":             resourceschema.StringAttribute{Required: true},
			"tls_enable_renegotiation":   resourceschema.BoolAttribute{Optional: true},
			"tls_skip_cert_verify":       resourceschema.BoolAttribute{Optional: true},
			"apply_status":               resourceschema.StringAttribute{Computed: true},
			"last_apply_error":           resourceschema.StringAttribute{Computed: true},
			"config_generation":          resourceschema.Int64Attribute{Computed: true},
			"applied_generation":         resourceschema.Int64Attribute{Computed: true},
			"transform": resourceschema.SingleNestedAttribute{Optional: true, Computed: true, Attributes: map[string]resourceschema.Attribute{
				"selected": resourceschema.StringAttribute{Optional: true}, "effective": resourceschema.StringAttribute{Computed: true}, "state": resourceschema.StringAttribute{Computed: true}, "applied": resourceschema.BoolAttribute{Computed: true}, "generation": resourceschema.Int64Attribute{Computed: true}, "last_error": resourceschema.StringAttribute{Computed: true},
			}},
		},
		Blocks: map[string]resourceschema.Block{
			"credentials": resourceschema.SingleNestedBlock{
				Attributes: map[string]resourceschema.Attribute{
					"urls":         resourceschema.ListAttribute{Optional: true, Sensitive: true, WriteOnly: true, ElementType: types.StringType},
					"urls_version": resourceschema.Int64Attribute{Optional: true},
				},
			},
			"tls": resourceschema.SingleNestedBlock{Attributes: map[string]resourceschema.Attribute{
				"enabled": resourceschema.BoolAttribute{Computed: true}, "state": resourceschema.StringAttribute{Computed: true}, "desired_generation": resourceschema.Int64Attribute{Computed: true}, "applied_generation": resourceschema.Int64Attribute{Computed: true}, "last_error": resourceschema.StringAttribute{Computed: true}, "ca_cert_filename": resourceschema.StringAttribute{Optional: true}, "ca_cert_content_base64": resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true}, "ca_cert_version": resourceschema.Int64Attribute{Optional: true}, "ca_cert_checksum_sha256": resourceschema.StringAttribute{Computed: true},
			}},
			"mtls": resourceschema.SingleNestedBlock{Attributes: map[string]resourceschema.Attribute{
				"enabled": resourceschema.BoolAttribute{Computed: true}, "state": resourceschema.StringAttribute{Computed: true}, "desired_generation": resourceschema.Int64Attribute{Computed: true}, "applied_generation": resourceschema.Int64Attribute{Computed: true}, "last_error": resourceschema.StringAttribute{Computed: true},
				"ca_cert_filename": resourceschema.StringAttribute{Optional: true}, "ca_cert_content_base64": resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true}, "ca_cert_version": resourceschema.Int64Attribute{Optional: true}, "ca_cert_checksum_sha256": resourceschema.StringAttribute{Computed: true},
				"client_cert_filename": resourceschema.StringAttribute{Optional: true}, "client_cert_content_base64": resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true}, "client_cert_version": resourceschema.Int64Attribute{Optional: true}, "client_cert_checksum_sha256": resourceschema.StringAttribute{Computed: true},
				"client_key_filename": resourceschema.StringAttribute{Optional: true}, "client_key_content_base64": resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true}, "client_key_version": resourceschema.Int64Attribute{Optional: true}, "client_key_checksum_sha256": resourceschema.StringAttribute{Computed: true},
			}},
			"sasl_anonymous": resourceschema.SingleNestedBlock{Attributes: map[string]resourceschema.Attribute{
				"enabled": resourceschema.BoolAttribute{Computed: true}, "state": resourceschema.StringAttribute{Computed: true}, "desired_generation": resourceschema.Int64Attribute{Computed: true}, "applied_generation": resourceschema.Int64Attribute{Computed: true}, "last_error": resourceschema.StringAttribute{Computed: true},
			}},
			"sasl_plain": resourceschema.SingleNestedBlock{Attributes: map[string]resourceschema.Attribute{
				"enabled": resourceschema.BoolAttribute{Computed: true}, "state": resourceschema.StringAttribute{Computed: true}, "desired_generation": resourceschema.Int64Attribute{Computed: true}, "applied_generation": resourceschema.Int64Attribute{Computed: true}, "last_error": resourceschema.StringAttribute{Computed: true}, "user": resourceschema.StringAttribute{Optional: true}, "password": resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true}, "password_version": resourceschema.Int64Attribute{Optional: true},
			}},
		},
	}
}

func (r *amqp1OutletResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	outletClient, ok := req.ProviderData.(amqp1OutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected AMQP 1.0 outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}
	r.client = outletClient
}

func (r *amqp1OutletResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config amqp1OutletResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	credentials, diags := expandAmqpURLsCredentials(ctx, config.Credentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.Credentials.IsNull() || config.Credentials.IsUnknown() || credentials.URLs.IsNull() || credentials.URLs.IsUnknown() || credentials.URLsVersion.IsNull() || credentials.URLsVersion.IsUnknown() {
		resp.Diagnostics.AddError("Missing AMQP 1.0 Credentials", "The `credentials` block must set both `urls` and `urls_version`.")
	}
	if !config.TLS.IsNull() && !config.TLS.IsUnknown() && !config.MTLS.IsNull() && !config.MTLS.IsUnknown() {
		resp.Diagnostics.AddError("Conflicting AMQP 1.0 TLS Features", "Only one of `tls` or `mtls` can be configured at the same time.")
	}
	if !config.SASLAnonymous.IsNull() && !config.SASLAnonymous.IsUnknown() && !config.SASLPlain.IsNull() && !config.SASLPlain.IsUnknown() {
		resp.Diagnostics.AddError("Conflicting AMQP 1.0 SASL Features", "Only one of `sasl_anonymous` or `sasl_plain` can be configured at the same time.")
	}
	if !config.TLS.IsNull() && !config.TLS.IsUnknown() {
		tls, diags := expandAMQPTLS(ctx, config.TLS)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if tls.CACertFilename.IsNull() || tls.CACertFilename.IsUnknown() || tls.CACertContentBase64.IsNull() || tls.CACertContentBase64.IsUnknown() || tls.CACertVersion.IsNull() || tls.CACertVersion.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete AMQP 1.0 TLS Configuration", "When `tls` is configured, `ca_cert_filename`, `ca_cert_content_base64`, and `ca_cert_version` must all be set.")
		}
	}
	if !config.MTLS.IsNull() && !config.MTLS.IsUnknown() {
		mtls, diags := expandAMQPMTLS(ctx, config.MTLS)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if mtls.CACertFilename.IsNull() || mtls.CACertFilename.IsUnknown() ||
			mtls.CACertContentBase64.IsNull() || mtls.CACertContentBase64.IsUnknown() || mtls.CACertVersion.IsNull() || mtls.CACertVersion.IsUnknown() ||
			mtls.ClientCertFilename.IsNull() || mtls.ClientCertFilename.IsUnknown() || mtls.ClientCertContentBase64.IsNull() || mtls.ClientCertContentBase64.IsUnknown() || mtls.ClientCertVersion.IsNull() || mtls.ClientCertVersion.IsUnknown() ||
			mtls.ClientKeyFilename.IsNull() || mtls.ClientKeyFilename.IsUnknown() || mtls.ClientKeyContentBase64.IsNull() || mtls.ClientKeyContentBase64.IsUnknown() || mtls.ClientKeyVersion.IsNull() || mtls.ClientKeyVersion.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete AMQP 1.0 mTLS Configuration", "When `mtls` is configured, CA cert, client cert, and client key filename/content/version values must all be set.")
		}
	}
	if !config.SASLPlain.IsNull() && !config.SASLPlain.IsUnknown() {
		plain, diags := expandAMQPSASLPlain(ctx, config.SASLPlain)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if plain.User.IsNull() || plain.User.IsUnknown() || plain.Password.IsNull() || plain.Password.IsUnknown() || plain.PasswordVersion.IsNull() || plain.PasswordVersion.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete AMQP 1.0 SASL/PLAIN Configuration", "When `sasl_plain` is configured, `user`, `password`, and `password_version` must all be set.")
		}
	}
}

func (r *amqp1OutletResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var plan amqp1OutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	credentials, tls, mtls, saslAnonymous, saslPlain, ok := r.readConfig(ctx, req.Config, plan, &resp.Diagnostics)
	if !ok {
		return
	}
	configuredTransform, transformDiags := selectedTransformFromConfig(ctx, req.Config)
	resp.Diagnostics.Append(transformDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	urls, err := stringListValueOrNil(ctx, credentials.URLs)
	if err != nil {
		resp.Diagnostics.AddError("Invalid AMQP 1.0 URLs", err.Error())
		return
	}
	metadataExcludePrefixes, err := stringListValueOrNil(ctx, plan.MetadataExcludePrefixes)
	if err != nil {
		resp.Diagnostics.AddError("Invalid AMQP 1.0 Metadata Exclude Prefixes", err.Error())
		return
	}
	outlet, err := r.client.CreateAmqp1Outlet(ctx, client.Amqp1OutletCreateRequest{
		Name: plan.Name.ValueString(), Description: stringValueOrEmpty(plan.Description), Enabled: plan.Enabled.ValueBool(),
		ApplicationPropertiesMap: stringPointerValue(plan.ApplicationPropertiesMap), MaxInFlight: int64PointerValue(plan.MaxInFlight), MetadataExcludePrefixes: metadataExcludePrefixes, NotifyPolicy: stringPointerValue(plan.NotifyPolicy),
		TargetAddress: plan.TargetAddress.ValueString(), TLSEnableRenegotiation: boolPointerValue(plan.TLSEnableRenegotiation), TLSSkipCertVerify: boolPointerValue(plan.TLSSkipCertVerify), URLs: urls,
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create AMQP 1.0 Outlet", err.Error())
		return
	}
	if err := r.applyFeatureConfig(ctx, outlet.ID, types.ObjectNull(amqpTLSAttributeTypes), types.ObjectNull(amqpMTLSAttributeTypes), types.ObjectNull(amqpSASLAnonymousAttributeTypes), types.ObjectNull(amqpSASLPlainAttributeTypes), tls, mtls, saslAnonymous, saslPlain); err != nil {
		if deleteErr := r.client.DeleteAmqp1Outlet(ctx, outlet.ID); deleteErr != nil && !client.IsNotFound(deleteErr) {
			resp.Diagnostics.AddError("Unable to Configure AMQP 1.0 Features", fmt.Sprintf("%s; cleanup failed: %s", err, deleteErr))
			return
		}
		resp.Diagnostics.AddError("Unable to Configure AMQP 1.0 Features", err.Error())
		return
	}
	if !configuredTransform.IsNull() && !configuredTransform.IsUnknown() {
		if _, err := r.client.SetTypedOutletTransform(ctx, "amqp_1", outlet.ID, configuredTransform.ValueString()); err != nil {
			if deleteErr := r.client.DeleteAmqp1Outlet(ctx, outlet.ID); deleteErr != nil && !client.IsNotFound(deleteErr) {
				resp.Diagnostics.AddError("Unable to Configure AMQP 1.0 Transform", fmt.Sprintf("%s; cleanup failed: %s", err, deleteErr))
				return
			}
			resp.Diagnostics.AddError("Unable to Configure AMQP 1.0 Transform", err.Error())
			return
		}
	}
	newState, removed, diags := r.buildState(ctx, outlet.ID, plan.Credentials, plan.TLS, plan.MTLS, plan.SASLPlain)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if removed {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "amqp_1", outlet.ID, plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe AMQP 1.0 Outlet Lifecycle", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *amqp1OutletResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var state amqp1OutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	newState, removed, diags := r.buildState(ctx, state.ID.ValueInt64(), state.Credentials, state.TLS, state.MTLS, state.SASLPlain)
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

func (r *amqp1OutletResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var plan, state amqp1OutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	credentials, tls, mtls, saslAnonymous, saslPlain, ok := r.readConfig(ctx, req.Config, plan, &resp.Diagnostics)
	if !ok {
		return
	}
	configuredTransform, transformDiags := selectedTransformFromConfig(ctx, req.Config)
	resp.Diagnostics.Append(transformDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	urls, err := stringListValueOrNil(ctx, credentials.URLs)
	if err != nil {
		resp.Diagnostics.AddError("Invalid AMQP 1.0 URLs", err.Error())
		return
	}
	metadataExcludePrefixes, err := stringListValueOrNil(ctx, plan.MetadataExcludePrefixes)
	if err != nil {
		resp.Diagnostics.AddError("Invalid AMQP 1.0 Metadata Exclude Prefixes", err.Error())
		return
	}
	if err := r.client.UpdateAmqp1Outlet(ctx, state.ID.ValueInt64(), client.Amqp1OutletUpdateRequest{
		Name: plan.Name.ValueString(), Description: stringValueOrEmpty(plan.Description), Enabled: plan.Enabled.ValueBool(),
		ApplicationPropertiesMap: stringPointerValue(plan.ApplicationPropertiesMap), MaxInFlight: int64PointerValue(plan.MaxInFlight), MetadataExcludePrefixes: metadataExcludePrefixes, NotifyPolicy: stringPointerValue(plan.NotifyPolicy),
		TargetAddress: plan.TargetAddress.ValueString(), TLSEnableRenegotiation: boolPointerValue(plan.TLSEnableRenegotiation), TLSSkipCertVerify: boolPointerValue(plan.TLSSkipCertVerify), URLs: urls,
	}); err != nil {
		resp.Diagnostics.AddError("Unable to Update AMQP 1.0 Outlet", err.Error())
		return
	}
	fieldUpdater := amqp1OutletFieldUpdater{client: r.client}
	mutations := []struct {
		name string
		err  error
	}{
		{"application_properties_map", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_1.application_properties_map", plan.ApplicationPropertiesMap, state.ApplicationPropertiesMap)},
		{"max_in_flight", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_1.max_in_flight", plan.MaxInFlight, state.MaxInFlight)},
		{"metadata_exclude_prefixes", syncStringListField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_1.metadata.exclude_prefixes", plan.MetadataExcludePrefixes, state.MetadataExcludePrefixes)},
		{"notify_policy", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "notify_policy", plan.NotifyPolicy, state.NotifyPolicy)},
		{"tls_enable_renegotiation", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_1.tls.enable_renegotiation", plan.TLSEnableRenegotiation, state.TLSEnableRenegotiation)},
		{"tls_skip_cert_verify", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_1.tls.skip_cert_verify", plan.TLSSkipCertVerify, state.TLSSkipCertVerify)},
	}
	for _, mutation := range mutations {
		if mutation.err != nil {
			resp.Diagnostics.AddError("Unable to Update AMQP 1.0 Field", fmt.Sprintf("%s: %s", mutation.name, mutation.err))
			return
		}
	}
	if err := syncTypedOutletTransform(ctx, r.client, "amqp_1", state.ID.ValueInt64(), configuredTransform, selectedTransform(state.Transform)); err != nil {
		resp.Diagnostics.AddError("Unable to Update AMQP 1.0 Transform", err.Error())
		return
	}
	if err := r.applyFeatureConfig(ctx, state.ID.ValueInt64(), state.TLS, state.MTLS, state.SASLAnonymous, state.SASLPlain, tls, mtls, saslAnonymous, saslPlain); err != nil {
		resp.Diagnostics.AddError("Unable to Update AMQP 1.0 Features", err.Error())
		return
	}
	newState, removed, diags := r.buildState(ctx, state.ID.ValueInt64(), plan.Credentials, plan.TLS, plan.MTLS, plan.SASLPlain)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if removed {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "amqp_1", state.ID.ValueInt64(), plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe AMQP 1.0 Outlet Lifecycle", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *amqp1OutletResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var state amqp1OutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	operation, err := r.client.DeleteAmqp1OutletOperation(ctx, state.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Delete AMQP 1.0 Outlet", err.Error())
		return
	}
	if err == nil && operation.DesiredGeneration != 0 {
		if err := observeTypedOutletLifecycle(ctx, r.client, "amqp_1", state.ID.ValueInt64(), operation.DesiredGeneration, false, true); err != nil {
			resp.Diagnostics.AddError("Unable to Observe AMQP 1.0 Outlet Lifecycle", err.Error())
		}
	}
}

func (r *amqp1OutletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid AMQP 1.0 Outlet Import Identifier", fmt.Sprintf("Expected numeric outlet ID, got %q: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (r *amqp1OutletResource) readConfig(ctx context.Context, cfg tfsdk.Config, plan amqp1OutletResourceModel, diags *diag.Diagnostics) (amqpURLsCredentialsModel, amqpTLSModel, amqpMTLSModel, amqpSASLAnonymousModel, amqpSASLPlainModel, bool) {
	credentials, d1 := expandAmqpURLsCredentials(ctx, plan.Credentials)
	diags.Append(d1...)
	tls, d2 := expandAMQPTLS(ctx, plan.TLS)
	diags.Append(d2...)
	mtls, d3 := expandAMQPMTLS(ctx, plan.MTLS)
	diags.Append(d3...)
	saslAnonymous, d4 := expandAMQPSASLAnonymous(ctx, plan.SASLAnonymous)
	diags.Append(d4...)
	saslPlain, d5 := expandAMQPSASLPlain(ctx, plan.SASLPlain)
	diags.Append(d5...)
	if diags.HasError() {
		return credentials, tls, mtls, saslAnonymous, saslPlain, false
	}
	credentials, d1 = enrichAmqpURLsCredentialsFromConfig(ctx, cfg, credentials)
	diags.Append(d1...)
	tls, d2 = enrichAMQPTLSFromConfig(ctx, cfg, "tls", tls)
	diags.Append(d2...)
	mtls, d3 = enrichAMQPMTLSFromConfig(ctx, cfg, "mtls", mtls)
	diags.Append(d3...)
	saslPlain, d5 = enrichAMQPSASLPlainFromConfig(ctx, cfg, "sasl_plain", saslPlain)
	diags.Append(d5...)
	if diags.HasError() {
		return credentials, tls, mtls, saslAnonymous, saslPlain, false
	}
	return credentials, tls, mtls, saslAnonymous, saslPlain, true
}

func (r *amqp1OutletResource) applyFeatureConfig(ctx context.Context, outletID int64, stateTLS, stateMTLS, stateSASLAnonymous, stateSASLPlain types.Object, planTLS amqpTLSModel, planMTLS amqpMTLSModel, planSASLAnonymous amqpSASLAnonymousModel, planSASLPlain amqpSASLPlainModel) error {
	if !stateTLS.IsNull() && !amqpFeatureEnabled(stateTLS) {
		stateTLS = types.ObjectNull(amqpTLSAttributeTypes)
	}
	if !stateMTLS.IsNull() && !amqpFeatureEnabled(stateMTLS) {
		stateMTLS = types.ObjectNull(amqpMTLSAttributeTypes)
	}
	if !stateSASLAnonymous.IsNull() && !amqpFeatureEnabled(stateSASLAnonymous) {
		stateSASLAnonymous = types.ObjectNull(amqpSASLAnonymousAttributeTypes)
	}
	if !stateSASLPlain.IsNull() && !amqpFeatureEnabled(stateSASLPlain) {
		stateSASLPlain = types.ObjectNull(amqpSASLPlainAttributeTypes)
	}

	if !stateTLS.IsNull() && !stateTLS.IsUnknown() && !planMTLS.ClientKeyVersion.IsNull() {
		if err := r.client.DeleteAmqp1TLSFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
			return err
		}
		if _, err := r.client.DeleteAmqp1Credential(ctx, outletID, amqpCACertRole); err != nil && !client.IsNotFound(err) {
			return err
		}
		stateTLS = types.ObjectNull(amqpTLSAttributeTypes)
	}
	if !stateMTLS.IsNull() && !stateMTLS.IsUnknown() && !planTLS.CACertVersion.IsNull() {
		if err := r.client.DeleteAmqp1MTLSFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
			return err
		}
		for _, role := range []string{amqpCACertRole, amqpClientCertRole, amqpClientKeyRole} {
			if _, err := r.client.DeleteAmqp1Credential(ctx, outletID, role); err != nil && !client.IsNotFound(err) {
				return err
			}
		}
		stateMTLS = types.ObjectNull(amqpMTLSAttributeTypes)
	}
	if !stateSASLAnonymous.IsNull() && !stateSASLAnonymous.IsUnknown() && !planSASLPlain.PasswordVersion.IsNull() {
		if err := r.client.DeleteAmqp1SASLAnonymousFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
			return err
		}
		stateSASLAnonymous = types.ObjectNull(amqpSASLAnonymousAttributeTypes)
	}
	if !stateSASLPlain.IsNull() && !stateSASLPlain.IsUnknown() && !planSASLAnonymous.Enabled.IsNull() {
		if err := r.client.DeleteAmqp1SASLPlainFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
			return err
		}
		if _, err := r.client.DeleteAmqp1Credential(ctx, outletID, amqpSASLPasswordRole); err != nil && !client.IsNotFound(err) {
			return err
		}
		stateSASLPlain = types.ObjectNull(amqpSASLPlainAttributeTypes)
	}

	if planTLS.CACertVersion.IsNull() {
		if !stateTLS.IsNull() && !stateTLS.IsUnknown() {
			if err := r.client.DeleteAmqp1TLSFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
				return err
			}
			if _, err := r.client.DeleteAmqp1Credential(ctx, outletID, amqpCACertRole); err != nil && !client.IsNotFound(err) {
				return err
			}
		}
	} else if stateTLS.IsNull() || stateTLS.IsUnknown() {
		if err := r.client.ConfigureAmqp1TLSFeature(ctx, outletID, client.AmqpTLSFeatureRequest{CACertFilename: planTLS.CACertFilename.ValueString(), CACertContentBase64: planTLS.CACertContentBase64.ValueString()}); err != nil {
			return err
		}
	} else {
		stateModel, diags := expandAMQPTLS(ctx, stateTLS)
		if diags.HasError() {
			return errors.New(diags.Errors()[0].Detail())
		}
		if !planTLS.CACertVersion.Equal(stateModel.CACertVersion) || !planTLS.CACertFilename.Equal(stateModel.CACertFilename) {
			if err := r.client.UpdateAmqp1TLSFeature(ctx, outletID, client.AmqpTLSFeatureRequest{CACertFilename: planTLS.CACertFilename.ValueString(), CACertContentBase64: planTLS.CACertContentBase64.ValueString()}); err != nil {
				return err
			}
		}
	}
	if planMTLS.ClientKeyVersion.IsNull() {
		if !stateMTLS.IsNull() && !stateMTLS.IsUnknown() {
			if err := r.client.DeleteAmqp1MTLSFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
				return err
			}
			for _, role := range []string{amqpCACertRole, amqpClientCertRole, amqpClientKeyRole} {
				if _, err := r.client.DeleteAmqp1Credential(ctx, outletID, role); err != nil && !client.IsNotFound(err) {
					return err
				}
			}
		}
	} else if stateMTLS.IsNull() || stateMTLS.IsUnknown() {
		if err := r.client.ConfigureAmqp1MTLSFeature(ctx, outletID, client.AmqpMTLSFeatureRequest{
			CACertFilename: planMTLS.CACertFilename.ValueString(), CACertContentBase64: planMTLS.CACertContentBase64.ValueString(), ClientCertFilename: planMTLS.ClientCertFilename.ValueString(), ClientCertContentBase64: planMTLS.ClientCertContentBase64.ValueString(), ClientKeyFilename: planMTLS.ClientKeyFilename.ValueString(), ClientKeyContentBase64: planMTLS.ClientKeyContentBase64.ValueString(),
		}); err != nil {
			return err
		}
	} else {
		stateModel, diags := expandAMQPMTLS(ctx, stateMTLS)
		if diags.HasError() {
			return errors.New(diags.Errors()[0].Detail())
		}
		if !planMTLS.CACertVersion.Equal(stateModel.CACertVersion) || !planMTLS.CACertFilename.Equal(stateModel.CACertFilename) || !planMTLS.ClientCertVersion.Equal(stateModel.ClientCertVersion) || !planMTLS.ClientCertFilename.Equal(stateModel.ClientCertFilename) || !planMTLS.ClientKeyVersion.Equal(stateModel.ClientKeyVersion) || !planMTLS.ClientKeyFilename.Equal(stateModel.ClientKeyFilename) {
			if err := r.client.UpdateAmqp1MTLSFeature(ctx, outletID, client.AmqpMTLSFeatureRequest{
				CACertFilename: planMTLS.CACertFilename.ValueString(), CACertContentBase64: planMTLS.CACertContentBase64.ValueString(), ClientCertFilename: planMTLS.ClientCertFilename.ValueString(), ClientCertContentBase64: planMTLS.ClientCertContentBase64.ValueString(), ClientKeyFilename: planMTLS.ClientKeyFilename.ValueString(), ClientKeyContentBase64: planMTLS.ClientKeyContentBase64.ValueString(),
			}); err != nil {
				return err
			}
		}
	}
	if planSASLPlain.PasswordVersion.IsNull() {
		if !stateSASLPlain.IsNull() && !stateSASLPlain.IsUnknown() {
			if err := r.client.DeleteAmqp1SASLPlainFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
				return err
			}
			if _, err := r.client.DeleteAmqp1Credential(ctx, outletID, amqpSASLPasswordRole); err != nil && !client.IsNotFound(err) {
				return err
			}
		}
	} else if stateSASLPlain.IsNull() || stateSASLPlain.IsUnknown() {
		if err := r.client.ConfigureAmqp1SASLPlainFeature(ctx, outletID, client.Amqp1SASLPlainFeatureRequest{User: planSASLPlain.User.ValueString(), Password: planSASLPlain.Password.ValueString()}); err != nil {
			return err
		}
	} else {
		stateModel, diags := expandAMQPSASLPlain(ctx, stateSASLPlain)
		if diags.HasError() {
			return errors.New(diags.Errors()[0].Detail())
		}
		if !planSASLPlain.User.Equal(stateModel.User) || !planSASLPlain.PasswordVersion.Equal(stateModel.PasswordVersion) {
			if err := r.client.UpdateAmqp1SASLPlainFeature(ctx, outletID, client.Amqp1SASLPlainFeatureRequest{User: planSASLPlain.User.ValueString(), Password: planSASLPlain.Password.ValueString()}); err != nil {
				return err
			}
		}
	}
	if planSASLAnonymous.Enabled.IsNull() {
		if !stateSASLAnonymous.IsNull() && !stateSASLAnonymous.IsUnknown() {
			if err := r.client.DeleteAmqp1SASLAnonymousFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
				return err
			}
		}
	} else if stateSASLAnonymous.IsNull() || stateSASLAnonymous.IsUnknown() {
		if err := r.client.ConfigureAmqp1SASLAnonymousFeature(ctx, outletID); err != nil {
			return err
		}
	} else {
		if err := r.client.UpdateAmqp1SASLAnonymousFeature(ctx, outletID); err != nil {
			return err
		}
	}
	return nil
}

func (r *amqp1OutletResource) buildState(ctx context.Context, outletID int64, previousCredentials, previousTLS, previousMTLS, previousSASLPlain types.Object) (amqp1OutletResourceModel, bool, diag.Diagnostics) {
	outlet, err := r.client.GetAmqp1Outlet(ctx, outletID)
	if err != nil {
		if client.IsNotFound(err) {
			return amqp1OutletResourceModel{}, true, nil
		}
		var diags diag.Diagnostics
		diags.AddError("Unable to Read AMQP 1.0 Outlet", err.Error())
		return amqp1OutletResourceModel{}, false, diags
	}
	credentials, err := r.client.ListAmqp1Credentials(ctx, outletID)
	if err != nil {
		if !client.IsNotFound(err) {
			var diags diag.Diagnostics
			diags.AddError("Unable to Read AMQP 1.0 Credentials", err.Error())
			return amqp1OutletResourceModel{}, false, diags
		}
		credentials = nil
	}
	tlsFeature, err := r.client.GetAmqp1TLSFeature(ctx, outletID)
	var tlsFeaturePtr *client.OutletFeatureDetail
	if err != nil && !client.IsNotFound(err) {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read AMQP 1.0 TLS Feature", err.Error())
		return amqp1OutletResourceModel{}, false, diags
	}
	if err == nil && tlsFeature.Feature != "" {
		tlsFeaturePtr = &tlsFeature
	}
	mtlsFeature, err := r.client.GetAmqp1MTLSFeature(ctx, outletID)
	var mtlsFeaturePtr *client.OutletFeatureDetail
	if err != nil && !client.IsNotFound(err) {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read AMQP 1.0 mTLS Feature", err.Error())
		return amqp1OutletResourceModel{}, false, diags
	}
	if err == nil && mtlsFeature.Feature != "" {
		mtlsFeaturePtr = &mtlsFeature
	}
	saslAnonymousFeature, err := r.client.GetAmqp1SASLAnonymousFeature(ctx, outletID)
	var saslAnonymousFeaturePtr *client.OutletFeatureDetail
	if err != nil && !client.IsNotFound(err) {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read AMQP 1.0 SASL Anonymous Feature", err.Error())
		return amqp1OutletResourceModel{}, false, diags
	}
	if err == nil && saslAnonymousFeature.Feature != "" {
		saslAnonymousFeaturePtr = &saslAnonymousFeature
	}
	saslPlainFeature, err := r.client.GetAmqp1SASLPlainFeature(ctx, outletID)
	var saslPlainFeaturePtr *client.OutletFeatureDetail
	if err != nil && !client.IsNotFound(err) {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read AMQP 1.0 SASL Plain Feature", err.Error())
		return amqp1OutletResourceModel{}, false, diags
	}
	if err == nil && saslPlainFeature.Feature != "" {
		saslPlainFeaturePtr = &saslPlainFeature
	}
	credentialsObject, diags := buildAmqpURLsCredentialsObject(ctx, previousCredentials, credentials, amqp1URLsCredentialRole)
	if diags.HasError() {
		return amqp1OutletResourceModel{}, false, diags
	}
	tlsObject, diags := buildAMQPTLSObject(ctx, previousTLS, tlsFeaturePtr, credentials)
	if diags.HasError() {
		return amqp1OutletResourceModel{}, false, diags
	}
	mtlsObject, diags := buildAMQPMTLSObject(ctx, previousMTLS, mtlsFeaturePtr, credentials)
	if diags.HasError() {
		return amqp1OutletResourceModel{}, false, diags
	}
	saslAnonymousObject, diags := buildAMQPSASLAnonymousObject(saslAnonymousFeaturePtr)
	if diags.HasError() {
		return amqp1OutletResourceModel{}, false, diags
	}
	saslPlainObject, diags := buildAMQPSASLPlainObject(ctx, previousSASLPlain, saslPlainFeaturePtr, credentials)
	if diags.HasError() {
		return amqp1OutletResourceModel{}, false, diags
	}
	transform, err := r.client.GetTypedOutletTransform(ctx, "amqp_1", outletID)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read AMQP 1.0 Transform", err.Error())
		return amqp1OutletResourceModel{}, false, diags
	}
	return amqp1OutletModelFromAPI(outlet, credentialsObject, tlsObject, mtlsObject, saslAnonymousObject, saslPlainObject, transform), false, nil
}

func amqp1OutletModelFromAPI(outlet client.Amqp1Outlet, credentials, tls, mtls, saslAnonymous, saslPlain types.Object, transform client.OutletTransformState) amqp1OutletResourceModel {
	return amqp1OutletResourceModel{
		ID: types.Int64Value(outlet.ID), Name: types.StringValue(outlet.Name), Description: nullableString(outlet.Description), Enabled: types.BoolValue(outlet.Enabled),
		ApplicationPropertiesMap: nullableStringPointer(outlet.ApplicationPropertiesMap), MaxInFlight: nullableInt64Pointer(outlet.MaxInFlight), MetadataExcludePrefixes: nullableStringList(outlet.MetadataExcludePrefixes), NotifyPolicy: nullableStringPointerPreservingEmpty(outlet.NotifyPolicy),
		TargetAddress: types.StringValue(outlet.TargetAddress), TLSEnableRenegotiation: nullableBoolPointer(outlet.TLSEnableRenegotiation), TLSSkipCertVerify: nullableBoolPointer(outlet.TLSSkipCertVerify),
		Credentials: credentials, TLS: tls, MTLS: mtls, SASLAnonymous: saslAnonymous, SASLPlain: saslPlain,
		ApplyStatus: nullableString(outlet.ApplyStatus), LastApplyError: nullableString(outlet.LastApplyError), ConfigGeneration: types.Int64Value(outlet.ConfigGeneration), AppliedGeneration: types.Int64Value(outlet.AppliedGeneration),
		Transform: transformObject(transform),
	}
}

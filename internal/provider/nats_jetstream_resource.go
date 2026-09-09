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
	_ resource.Resource                   = &natsJetstreamOutletResource{}
	_ resource.ResourceWithConfigure      = &natsJetstreamOutletResource{}
	_ resource.ResourceWithImportState    = &natsJetstreamOutletResource{}
	_ resource.ResourceWithValidateConfig = &natsJetstreamOutletResource{}
)

type natsJetstreamOutletClient interface {
	CreateNatsJetstreamOutlet(context.Context, client.NatsJetstreamOutletCreateRequest) (client.NatsJetstreamOutlet, error)
	GetNatsJetstreamOutlet(context.Context, int64) (client.NatsJetstreamOutlet, error)
	UpdateNatsJetstreamOutlet(context.Context, int64, client.NatsJetstreamOutletUpdateRequest) error
	DeleteNatsJetstreamOutlet(context.Context, int64) error
	DeleteNatsJetstreamOutletOperation(context.Context, int64) (client.OutletOperation, error)
	GetTypedOutletStatus(context.Context, string, int64, int64) (client.OutletStatus, error)
	SetNatsJetstreamField(context.Context, int64, string, any) error
	DeleteNatsJetstreamField(context.Context, int64, string) error
	ListNatsJetstreamCredentials(context.Context, int64) ([]client.OutletMaterialSummary, error)
	DeleteNatsJetstreamCredential(context.Context, int64, string) (client.OutletMaterialDeleteResponse, error)
	GetNatsJetstreamJWTSeedFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureNatsJetstreamJWTSeedFeature(context.Context, int64, client.NatsJetstreamJWTSeedFeatureRequest) error
	UpdateNatsJetstreamJWTSeedFeature(context.Context, int64, client.NatsJetstreamJWTSeedFeatureRequest) error
	DeleteNatsJetstreamJWTSeedFeature(context.Context, int64) error
	GetNatsJetstreamNkeyFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureNatsJetstreamNkeyFeature(context.Context, int64, client.NatsJetstreamNkeyFeatureRequest) error
	UpdateNatsJetstreamNkeyFeature(context.Context, int64, client.NatsJetstreamNkeyFeatureRequest) error
	DeleteNatsJetstreamNkeyFeature(context.Context, int64) error
	GetNatsJetstreamUserCredentialsFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureNatsJetstreamUserCredentialsFeature(context.Context, int64, client.NatsJetstreamUserCredentialsFeatureRequest) error
	UpdateNatsJetstreamUserCredentialsFeature(context.Context, int64, client.NatsJetstreamUserCredentialsFeatureRequest) error
	DeleteNatsJetstreamUserCredentialsFeature(context.Context, int64) error
	typedOutletTransformClient
}

type natsJetstreamOutletFieldUpdater struct{ client natsJetstreamOutletClient }

func (u natsJetstreamOutletFieldUpdater) SetField(ctx context.Context, outletID int64, field string, value any) error {
	return u.client.SetNatsJetstreamField(ctx, outletID, field, value)
}

func (u natsJetstreamOutletFieldUpdater) DeleteField(ctx context.Context, outletID int64, field string) error {
	return u.client.DeleteNatsJetstreamField(ctx, outletID, field)
}

type natsJetstreamOutletResource struct{ client natsJetstreamOutletClient }

type natsJetstreamOutletResourceModel struct {
	ID                      types.Int64   `tfsdk:"id"`
	Name                    types.String  `tfsdk:"name"`
	Description             types.String  `tfsdk:"description"`
	Enabled                 types.Bool    `tfsdk:"enabled"`
	BackoffInitialInterval  types.String  `tfsdk:"backoff_initial_interval"`
	BackoffJitter           types.Float64 `tfsdk:"backoff_jitter"`
	BackoffMaxElapsedTime   types.String  `tfsdk:"backoff_max_elapsed_time"`
	BackoffMaxInterval      types.String  `tfsdk:"backoff_max_interval"`
	Headers                 types.Map     `tfsdk:"headers"`
	InjectTracingMap        types.String  `tfsdk:"inject_tracing_map"`
	MaxInFlight             types.Int64   `tfsdk:"max_in_flight"`
	MaxRetries              types.Int64   `tfsdk:"max_retries"`
	MetadataIncludePatterns types.List    `tfsdk:"metadata_include_patterns"`
	MetadataIncludePrefixes types.List    `tfsdk:"metadata_include_prefixes"`
	NotifyPolicy            types.String  `tfsdk:"notify_policy"`
	Subject                 types.String  `tfsdk:"subject"`
	TLSEnableRenegotiation  types.Bool    `tfsdk:"tls_enable_renegotiation"`
	TLSSkipCertVerify       types.Bool    `tfsdk:"tls_skip_cert_verify"`
	Credentials             types.Object  `tfsdk:"credentials"`
	JWTSeed                 types.Object  `tfsdk:"jwt_seed"`
	Nkey                    types.Object  `tfsdk:"nkey"`
	UserCredentials         types.Object  `tfsdk:"user_credentials"`
	ApplyStatus             types.String  `tfsdk:"apply_status"`
	LastApplyError          types.String  `tfsdk:"last_apply_error"`
	ConfigGeneration        types.Int64   `tfsdk:"config_generation"`
	AppliedGeneration       types.Int64   `tfsdk:"applied_generation"`
	Transform               types.Object  `tfsdk:"transform"`
}

func NewNatsJetstreamOutletResource() resource.Resource { return &natsJetstreamOutletResource{} }

func (r *natsJetstreamOutletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nats_jetstream_outlet"
}

func (r *natsJetstreamOutletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		MarkdownDescription: "Manages a Noozle NATS JetStream outlet.",
		Attributes: map[string]resourceschema.Attribute{
			"id":                        resourceschema.Int64Attribute{Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}},
			"name":                      resourceschema.StringAttribute{Required: true},
			"description":               resourceschema.StringAttribute{Optional: true},
			"enabled":                   resourceschema.BoolAttribute{Required: true},
			"backoff_initial_interval":  resourceschema.StringAttribute{Optional: true},
			"backoff_jitter":            resourceschema.Float64Attribute{Optional: true},
			"backoff_max_elapsed_time":  resourceschema.StringAttribute{Optional: true},
			"backoff_max_interval":      resourceschema.StringAttribute{Optional: true},
			"headers":                   resourceschema.MapAttribute{Optional: true, ElementType: types.StringType},
			"inject_tracing_map":        resourceschema.StringAttribute{Optional: true},
			"max_in_flight":             resourceschema.Int64Attribute{Optional: true},
			"max_retries":               resourceschema.Int64Attribute{Optional: true},
			"metadata_include_patterns": resourceschema.ListAttribute{Optional: true, ElementType: types.StringType},
			"metadata_include_prefixes": resourceschema.ListAttribute{Optional: true, ElementType: types.StringType},
			"notify_policy":             resourceschema.StringAttribute{Optional: true},
			"subject":                   resourceschema.StringAttribute{Required: true},
			"tls_enable_renegotiation":  resourceschema.BoolAttribute{Optional: true},
			"tls_skip_cert_verify":      resourceschema.BoolAttribute{Optional: true},
			"apply_status":              resourceschema.StringAttribute{Computed: true},
			"last_apply_error":          resourceschema.StringAttribute{Computed: true},
			"config_generation":         resourceschema.Int64Attribute{Computed: true},
			"applied_generation":        resourceschema.Int64Attribute{Computed: true},
			"transform": resourceschema.SingleNestedAttribute{MarkdownDescription: "Optional outlet transform selection and observed runtime state.", Optional: true, Computed: true, Attributes: map[string]resourceschema.Attribute{
				"selected": resourceschema.StringAttribute{MarkdownDescription: "Explicitly selected transform preset.", Optional: true}, "effective": resourceschema.StringAttribute{MarkdownDescription: "Effective transform after defaults are applied.", Computed: true}, "state": resourceschema.StringAttribute{MarkdownDescription: "Observed transform runtime state.", Computed: true}, "applied": resourceschema.BoolAttribute{MarkdownDescription: "Whether the transform has been applied.", Computed: true}, "generation": resourceschema.Int64Attribute{MarkdownDescription: "Observed transform generation.", Computed: true}, "last_error": resourceschema.StringAttribute{MarkdownDescription: "Last observed transform apply error.", Computed: true},
			}},
		},
		Blocks: map[string]resourceschema.Block{
			"credentials": resourceschema.SingleNestedBlock{
				MarkdownDescription: "Connection URLs for this outlet. The URLs are write-only; rotate them by updating the list and incrementing `urls_version`.",
				Attributes: map[string]resourceschema.Attribute{
					"urls":         resourceschema.ListAttribute{Optional: true, Sensitive: true, WriteOnly: true, ElementType: types.StringType},
					"urls_version": resourceschema.Int64Attribute{Optional: true},
				},
			},
			"jwt_seed": resourceschema.SingleNestedBlock{
				MarkdownDescription: "Optional JWT seed authentication feature. Conflicts with `nkey` and `user_credentials`.",
				Attributes: map[string]resourceschema.Attribute{
					"enabled":            resourceschema.BoolAttribute{Computed: true},
					"state":              resourceschema.StringAttribute{Computed: true},
					"desired_generation": resourceschema.Int64Attribute{Computed: true},
					"applied_generation": resourceschema.Int64Attribute{Computed: true},
					"last_error":         resourceschema.StringAttribute{Computed: true},
					"user_jwt":           resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true},
					"user_jwt_version":   resourceschema.Int64Attribute{Optional: true},
					"nkey_seed":          resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true},
					"nkey_seed_version":  resourceschema.Int64Attribute{Optional: true},
				},
			},
			"nkey": resourceschema.SingleNestedBlock{
				MarkdownDescription: "Optional NKEY file authentication feature. Conflicts with `jwt_seed` and `user_credentials`.",
				Attributes: map[string]resourceschema.Attribute{
					"enabled":            resourceschema.BoolAttribute{Computed: true},
					"state":              resourceschema.StringAttribute{Computed: true},
					"desired_generation": resourceschema.Int64Attribute{Computed: true},
					"applied_generation": resourceschema.Int64Attribute{Computed: true},
					"last_error":         resourceschema.StringAttribute{Computed: true},
					"filename":           resourceschema.StringAttribute{Optional: true},
					"content_base64":     resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true},
					"version":            resourceschema.Int64Attribute{Optional: true},
					"checksum_sha256":    resourceschema.StringAttribute{Computed: true},
				},
			},
			"user_credentials": resourceschema.SingleNestedBlock{
				MarkdownDescription: "Optional user credentials authentication feature. Conflicts with `jwt_seed` and `nkey`.",
				Attributes: map[string]resourceschema.Attribute{
					"enabled":            resourceschema.BoolAttribute{Computed: true},
					"state":              resourceschema.StringAttribute{Computed: true},
					"desired_generation": resourceschema.Int64Attribute{Computed: true},
					"applied_generation": resourceschema.Int64Attribute{Computed: true},
					"last_error":         resourceschema.StringAttribute{Computed: true},
					"filename":           resourceschema.StringAttribute{Optional: true},
					"content_base64":     resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true},
					"version":            resourceschema.Int64Attribute{Optional: true},
					"checksum_sha256":    resourceschema.StringAttribute{Computed: true},
				},
			},
		},
	}
}

func (r *natsJetstreamOutletResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	outletClient, ok := req.ProviderData.(natsJetstreamOutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected NATS JetStream outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}
	r.client = outletClient
}

func (r *natsJetstreamOutletResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config natsJetstreamOutletResourceModel
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
		resp.Diagnostics.AddError("Missing NATS JetStream Credentials", "The `credentials` block must set both `urls` and `urls_version`.")
	}

	configuredAuthBlocks := 0
	if !config.JWTSeed.IsNull() && !config.JWTSeed.IsUnknown() {
		configuredAuthBlocks++
		jwtSeed, diags := expandNatsJetstreamJWTSeed(ctx, config.JWTSeed)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if jwtSeed.UserJWT.IsNull() || jwtSeed.UserJWT.IsUnknown() || jwtSeed.UserJWTVersion.IsNull() || jwtSeed.UserJWTVersion.IsUnknown() ||
			jwtSeed.NkeySeed.IsNull() || jwtSeed.NkeySeed.IsUnknown() || jwtSeed.NkeySeedVersion.IsNull() || jwtSeed.NkeySeedVersion.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete NATS JetStream JWT Seed Configuration", "When `jwt_seed` is configured, `user_jwt`, `user_jwt_version`, `nkey_seed`, and `nkey_seed_version` must all be set.")
		}
	}
	if !config.Nkey.IsNull() && !config.Nkey.IsUnknown() {
		configuredAuthBlocks++
		nkey, diags := expandNatsJetstreamNkey(ctx, config.Nkey)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if nkey.Filename.IsNull() || nkey.Filename.IsUnknown() || nkey.ContentBase64.IsNull() || nkey.ContentBase64.IsUnknown() || nkey.Version.IsNull() || nkey.Version.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete NATS JetStream NKEY Configuration", "When `nkey` is configured, `filename`, `content_base64`, and `version` must all be set.")
		}
	}
	if !config.UserCredentials.IsNull() && !config.UserCredentials.IsUnknown() {
		configuredAuthBlocks++
		userCredentials, diags := expandNatsJetstreamUserCredentials(ctx, config.UserCredentials)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if userCredentials.Filename.IsNull() || userCredentials.Filename.IsUnknown() || userCredentials.ContentBase64.IsNull() || userCredentials.ContentBase64.IsUnknown() || userCredentials.Version.IsNull() || userCredentials.Version.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete NATS JetStream User Credentials Configuration", "When `user_credentials` is configured, `filename`, `content_base64`, and `version` must all be set.")
		}
	}

	if configuredAuthBlocks > 1 {
		resp.Diagnostics.AddError("Conflicting NATS JetStream Auth Features", "Only one of `jwt_seed`, `nkey`, or `user_credentials` can be configured at the same time.")
	}
}

func (r *natsJetstreamOutletResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var plan natsJetstreamOutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	credentials, jwtSeed, nkey, userCredentials, ok := r.readConfig(ctx, req.Config, plan, &resp.Diagnostics)
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
		resp.Diagnostics.AddError("Invalid NATS JetStream URLs", err.Error())
		return
	}
	headers, err := stringMapValueOrNil(ctx, plan.Headers)
	if err != nil {
		resp.Diagnostics.AddError("Invalid NATS JetStream Headers", err.Error())
		return
	}
	metadataIncludePatterns, err := stringListValueOrNil(ctx, plan.MetadataIncludePatterns)
	if err != nil {
		resp.Diagnostics.AddError("Invalid NATS JetStream Metadata Include Patterns", err.Error())
		return
	}
	metadataIncludePrefixes, err := stringListValueOrNil(ctx, plan.MetadataIncludePrefixes)
	if err != nil {
		resp.Diagnostics.AddError("Invalid NATS JetStream Metadata Include Prefixes", err.Error())
		return
	}

	outlet, err := r.client.CreateNatsJetstreamOutlet(ctx, client.NatsJetstreamOutletCreateRequest{
		Name:                    plan.Name.ValueString(),
		Description:             stringValueOrEmpty(plan.Description),
		Enabled:                 plan.Enabled.ValueBool(),
		BackoffInitialInterval:  stringPointerValue(plan.BackoffInitialInterval),
		BackoffJitter:           float64PointerValue(plan.BackoffJitter),
		BackoffMaxElapsedTime:   stringPointerValue(plan.BackoffMaxElapsedTime),
		BackoffMaxInterval:      stringPointerValue(plan.BackoffMaxInterval),
		Headers:                 headers,
		InjectTracingMap:        stringPointerValue(plan.InjectTracingMap),
		MaxInFlight:             int64PointerValue(plan.MaxInFlight),
		MaxRetries:              int64PointerValue(plan.MaxRetries),
		MetadataIncludePatterns: metadataIncludePatterns,
		MetadataIncludePrefixes: metadataIncludePrefixes,
		NotifyPolicy:            stringPointerValue(plan.NotifyPolicy),
		Subject:                 plan.Subject.ValueString(),
		TLSEnableRenegotiation:  boolPointerValue(plan.TLSEnableRenegotiation),
		TLSSkipCertVerify:       boolPointerValue(plan.TLSSkipCertVerify),
		URLs:                    urls,
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create NATS JetStream Outlet", err.Error())
		return
	}

	if err := r.applyFeatureConfig(ctx, outlet.ID, types.ObjectNull(natsJetstreamJWTSeedAttributeTypes), types.ObjectNull(natsJetstreamNkeyAttributeTypes), types.ObjectNull(natsJetstreamUserCredentialsAttributeTypes), jwtSeed, nkey, userCredentials); err != nil {
		if deleteErr := r.client.DeleteNatsJetstreamOutlet(ctx, outlet.ID); deleteErr != nil && !client.IsNotFound(deleteErr) {
			resp.Diagnostics.AddError("Unable to Configure NATS JetStream Features", fmt.Sprintf("%s; cleanup failed: %s", err, deleteErr))
			return
		}
		resp.Diagnostics.AddError("Unable to Configure NATS JetStream Features", err.Error())
		return
	}
	if !configuredTransform.IsNull() && !configuredTransform.IsUnknown() {
		if _, err := r.client.SetTypedOutletTransform(ctx, "nats_jetstream", outlet.ID, configuredTransform.ValueString()); err != nil {
			if deleteErr := r.client.DeleteNatsJetstreamOutlet(ctx, outlet.ID); deleteErr != nil && !client.IsNotFound(deleteErr) {
				resp.Diagnostics.AddError("Unable to Configure NATS JetStream Transform", fmt.Sprintf("%s; cleanup failed: %s", err, deleteErr))
				return
			}
			resp.Diagnostics.AddError("Unable to Configure NATS JetStream Transform", err.Error())
			return
		}
	}

	newState, removed, diags := r.buildState(ctx, outlet.ID, plan.Credentials, plan.JWTSeed, plan.Nkey, plan.UserCredentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if removed {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "nats_jetstream", outlet.ID, plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe NATS JetStream Outlet Lifecycle", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *natsJetstreamOutletResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var state natsJetstreamOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	newState, removed, diags := r.buildState(ctx, state.ID.ValueInt64(), state.Credentials, state.JWTSeed, state.Nkey, state.UserCredentials)
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

func (r *natsJetstreamOutletResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var plan, state natsJetstreamOutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	credentials, jwtSeed, nkey, userCredentials, ok := r.readConfig(ctx, req.Config, plan, &resp.Diagnostics)
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
		resp.Diagnostics.AddError("Invalid NATS JetStream URLs", err.Error())
		return
	}
	headers, err := stringMapValueOrNil(ctx, plan.Headers)
	if err != nil {
		resp.Diagnostics.AddError("Invalid NATS JetStream Headers", err.Error())
		return
	}
	metadataIncludePatterns, err := stringListValueOrNil(ctx, plan.MetadataIncludePatterns)
	if err != nil {
		resp.Diagnostics.AddError("Invalid NATS JetStream Metadata Include Patterns", err.Error())
		return
	}
	metadataIncludePrefixes, err := stringListValueOrNil(ctx, plan.MetadataIncludePrefixes)
	if err != nil {
		resp.Diagnostics.AddError("Invalid NATS JetStream Metadata Include Prefixes", err.Error())
		return
	}

	if err := r.client.UpdateNatsJetstreamOutlet(ctx, state.ID.ValueInt64(), client.NatsJetstreamOutletUpdateRequest{
		Name:                    plan.Name.ValueString(),
		Description:             stringValueOrEmpty(plan.Description),
		Enabled:                 plan.Enabled.ValueBool(),
		BackoffInitialInterval:  stringPointerValue(plan.BackoffInitialInterval),
		BackoffJitter:           float64PointerValue(plan.BackoffJitter),
		BackoffMaxElapsedTime:   stringPointerValue(plan.BackoffMaxElapsedTime),
		BackoffMaxInterval:      stringPointerValue(plan.BackoffMaxInterval),
		Headers:                 headers,
		InjectTracingMap:        stringPointerValue(plan.InjectTracingMap),
		MaxInFlight:             int64PointerValue(plan.MaxInFlight),
		MaxRetries:              int64PointerValue(plan.MaxRetries),
		MetadataIncludePatterns: metadataIncludePatterns,
		MetadataIncludePrefixes: metadataIncludePrefixes,
		NotifyPolicy:            stringPointerValue(plan.NotifyPolicy),
		Subject:                 plan.Subject.ValueString(),
		TLSEnableRenegotiation:  boolPointerValue(plan.TLSEnableRenegotiation),
		TLSSkipCertVerify:       boolPointerValue(plan.TLSSkipCertVerify),
		URLs:                    urls,
	}); err != nil {
		resp.Diagnostics.AddError("Unable to Update NATS JetStream Outlet", err.Error())
		return
	}

	fieldUpdater := natsJetstreamOutletFieldUpdater{client: r.client}
	mutations := []struct {
		name string
		err  error
	}{
		{"backoff_initial_interval", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "nats_jetstream.backoff.initial_interval", plan.BackoffInitialInterval, state.BackoffInitialInterval)},
		{"backoff_jitter", syncFloat64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "nats_jetstream.backoff.jitter", plan.BackoffJitter, state.BackoffJitter)},
		{"backoff_max_elapsed_time", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "nats_jetstream.backoff.max_elapsed_time", plan.BackoffMaxElapsedTime, state.BackoffMaxElapsedTime)},
		{"backoff_max_interval", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "nats_jetstream.backoff.max_interval", plan.BackoffMaxInterval, state.BackoffMaxInterval)},
		{"headers", syncStringMapField(ctx, fieldUpdater, state.ID.ValueInt64(), "nats_jetstream.headers", plan.Headers, state.Headers)},
		{"inject_tracing_map", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "nats_jetstream.inject_tracing_map", plan.InjectTracingMap, state.InjectTracingMap)},
		{"max_in_flight", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "nats_jetstream.max_in_flight", plan.MaxInFlight, state.MaxInFlight)},
		{"max_retries", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "nats_jetstream.max_retries", plan.MaxRetries, state.MaxRetries)},
		{"metadata_include_patterns", syncStringListField(ctx, fieldUpdater, state.ID.ValueInt64(), "nats_jetstream.metadata.include_patterns", plan.MetadataIncludePatterns, state.MetadataIncludePatterns)},
		{"metadata_include_prefixes", syncStringListField(ctx, fieldUpdater, state.ID.ValueInt64(), "nats_jetstream.metadata.include_prefixes", plan.MetadataIncludePrefixes, state.MetadataIncludePrefixes)},
		{"notify_policy", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "notify_policy", plan.NotifyPolicy, state.NotifyPolicy)},
		{"tls_enable_renegotiation", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "nats_jetstream.tls.enable_renegotiation", plan.TLSEnableRenegotiation, state.TLSEnableRenegotiation)},
		{"tls_skip_cert_verify", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "nats_jetstream.tls.skip_cert_verify", plan.TLSSkipCertVerify, state.TLSSkipCertVerify)},
	}
	for _, mutation := range mutations {
		if mutation.err != nil {
			resp.Diagnostics.AddError("Unable to Update NATS JetStream Field", fmt.Sprintf("%s: %s", mutation.name, mutation.err))
			return
		}
	}

	if err := r.applyFeatureConfig(ctx, state.ID.ValueInt64(), state.JWTSeed, state.Nkey, state.UserCredentials, jwtSeed, nkey, userCredentials); err != nil {
		resp.Diagnostics.AddError("Unable to Update NATS JetStream Features", err.Error())
		return
	}
	if err := syncTypedOutletTransform(ctx, r.client, "nats_jetstream", state.ID.ValueInt64(), configuredTransform, selectedTransform(state.Transform)); err != nil {
		resp.Diagnostics.AddError("Unable to Update NATS JetStream Transform", err.Error())
		return
	}

	newState, removed, diags := r.buildState(ctx, state.ID.ValueInt64(), plan.Credentials, plan.JWTSeed, plan.Nkey, plan.UserCredentials)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if removed {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "nats_jetstream", state.ID.ValueInt64(), plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe NATS JetStream Outlet Lifecycle", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *natsJetstreamOutletResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var state natsJetstreamOutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	operation, err := r.client.DeleteNatsJetstreamOutletOperation(ctx, state.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Delete NATS JetStream Outlet", err.Error())
		return
	}
	if err == nil && operation.DesiredGeneration != 0 {
		if err := observeTypedOutletLifecycle(ctx, r.client, "nats_jetstream", state.ID.ValueInt64(), operation.DesiredGeneration, false, true); err != nil {
			resp.Diagnostics.AddError("Unable to Observe NATS JetStream Outlet Lifecycle", err.Error())
		}
	}
}

func (r *natsJetstreamOutletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid NATS JetStream Outlet Import Identifier", fmt.Sprintf("Expected numeric outlet ID, got %q: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (r *natsJetstreamOutletResource) readConfig(ctx context.Context, cfg tfsdk.Config, plan natsJetstreamOutletResourceModel, diags *diag.Diagnostics) (amqpURLsCredentialsModel, natsJetstreamJWTSeedModel, natsJetstreamNkeyModel, natsJetstreamUserCredentialsModel, bool) {
	credentials, credentialsDiags := expandAmqpURLsCredentials(ctx, plan.Credentials)
	diags.Append(credentialsDiags...)
	jwtSeed, jwtSeedDiags := expandNatsJetstreamJWTSeed(ctx, plan.JWTSeed)
	diags.Append(jwtSeedDiags...)
	nkey, nkeyDiags := expandNatsJetstreamNkey(ctx, plan.Nkey)
	diags.Append(nkeyDiags...)
	userCredentials, userCredentialsDiags := expandNatsJetstreamUserCredentials(ctx, plan.UserCredentials)
	diags.Append(userCredentialsDiags...)
	if diags.HasError() {
		return credentials, jwtSeed, nkey, userCredentials, false
	}

	credentials, credentialsDiags = enrichAmqpURLsCredentialsFromConfig(ctx, cfg, credentials)
	diags.Append(credentialsDiags...)
	jwtSeed, jwtSeedDiags = enrichNatsJetstreamJWTSeedFromConfig(ctx, cfg, jwtSeed)
	diags.Append(jwtSeedDiags...)
	nkey, nkeyDiags = enrichNatsJetstreamNkeyFromConfig(ctx, cfg, nkey)
	diags.Append(nkeyDiags...)
	userCredentials, userCredentialsDiags = enrichNatsJetstreamUserCredentialsFromConfig(ctx, cfg, userCredentials)
	diags.Append(userCredentialsDiags...)
	if diags.HasError() {
		return credentials, jwtSeed, nkey, userCredentials, false
	}
	return credentials, jwtSeed, nkey, userCredentials, true
}

func (r *natsJetstreamOutletResource) deleteJWTSeed(ctx context.Context, outletID int64) error {
	if err := r.client.DeleteNatsJetstreamJWTSeedFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
		return err
	}
	for _, role := range []string{natsJetstreamJWTUserRole, natsJetstreamNkeySeedRole} {
		if _, err := r.client.DeleteNatsJetstreamCredential(ctx, outletID, role); err != nil && !client.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (r *natsJetstreamOutletResource) deleteNkey(ctx context.Context, outletID int64) error {
	if err := r.client.DeleteNatsJetstreamNkeyFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
		return err
	}
	if _, err := r.client.DeleteNatsJetstreamCredential(ctx, outletID, natsJetstreamNkeyFileRole); err != nil && !client.IsNotFound(err) {
		return err
	}
	return nil
}

func (r *natsJetstreamOutletResource) deleteUserCredentials(ctx context.Context, outletID int64) error {
	if err := r.client.DeleteNatsJetstreamUserCredentialsFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
		return err
	}
	if _, err := r.client.DeleteNatsJetstreamCredential(ctx, outletID, natsJetstreamUserCredsRole); err != nil && !client.IsNotFound(err) {
		return err
	}
	return nil
}

func (r *natsJetstreamOutletResource) applyFeatureConfig(ctx context.Context, outletID int64, stateJWTSeed, stateNkey, stateUserCredentials types.Object, planJWTSeed natsJetstreamJWTSeedModel, planNkey natsJetstreamNkeyModel, planUserCredentials natsJetstreamUserCredentialsModel) error {
	if !stateJWTSeed.IsNull() && !natsFeatureEnabled(stateJWTSeed) {
		stateJWTSeed = types.ObjectNull(natsJetstreamJWTSeedAttributeTypes)
	}
	if !stateNkey.IsNull() && !natsFeatureEnabled(stateNkey) {
		stateNkey = types.ObjectNull(natsJetstreamNkeyAttributeTypes)
	}
	if !stateUserCredentials.IsNull() && !natsFeatureEnabled(stateUserCredentials) {
		stateUserCredentials = types.ObjectNull(natsJetstreamUserCredentialsAttributeTypes)
	}

	jwtConfigured := !planJWTSeed.UserJWTVersion.IsNull()
	nkeyConfigured := !planNkey.Version.IsNull()
	userCredsConfigured := !planUserCredentials.Version.IsNull()

	if jwtConfigured {
		if !stateNkey.IsNull() && !stateNkey.IsUnknown() {
			if err := r.deleteNkey(ctx, outletID); err != nil {
				return err
			}
			stateNkey = types.ObjectNull(natsJetstreamNkeyAttributeTypes)
		}
		if !stateUserCredentials.IsNull() && !stateUserCredentials.IsUnknown() {
			if err := r.deleteUserCredentials(ctx, outletID); err != nil {
				return err
			}
			stateUserCredentials = types.ObjectNull(natsJetstreamUserCredentialsAttributeTypes)
		}
	}
	if nkeyConfigured {
		if !stateJWTSeed.IsNull() && !stateJWTSeed.IsUnknown() {
			if err := r.deleteJWTSeed(ctx, outletID); err != nil {
				return err
			}
			stateJWTSeed = types.ObjectNull(natsJetstreamJWTSeedAttributeTypes)
		}
		if !stateUserCredentials.IsNull() && !stateUserCredentials.IsUnknown() {
			if err := r.deleteUserCredentials(ctx, outletID); err != nil {
				return err
			}
			stateUserCredentials = types.ObjectNull(natsJetstreamUserCredentialsAttributeTypes)
		}
	}
	if userCredsConfigured {
		if !stateJWTSeed.IsNull() && !stateJWTSeed.IsUnknown() {
			if err := r.deleteJWTSeed(ctx, outletID); err != nil {
				return err
			}
			stateJWTSeed = types.ObjectNull(natsJetstreamJWTSeedAttributeTypes)
		}
		if !stateNkey.IsNull() && !stateNkey.IsUnknown() {
			if err := r.deleteNkey(ctx, outletID); err != nil {
				return err
			}
			stateNkey = types.ObjectNull(natsJetstreamNkeyAttributeTypes)
		}
	}

	if !jwtConfigured {
		if !stateJWTSeed.IsNull() && !stateJWTSeed.IsUnknown() {
			if err := r.deleteJWTSeed(ctx, outletID); err != nil {
				return err
			}
		}
	} else if stateJWTSeed.IsNull() || stateJWTSeed.IsUnknown() {
		if err := r.client.ConfigureNatsJetstreamJWTSeedFeature(ctx, outletID, client.NatsJetstreamJWTSeedFeatureRequest{
			UserJWT: planJWTSeed.UserJWT.ValueString(), NkeySeed: planJWTSeed.NkeySeed.ValueString(),
		}); err != nil {
			return err
		}
	} else {
		stateModel, diags := expandNatsJetstreamJWTSeed(ctx, stateJWTSeed)
		if diags.HasError() {
			return errors.New(diags.Errors()[0].Detail())
		}
		if !planJWTSeed.UserJWTVersion.Equal(stateModel.UserJWTVersion) || !planJWTSeed.NkeySeedVersion.Equal(stateModel.NkeySeedVersion) {
			if err := r.client.UpdateNatsJetstreamJWTSeedFeature(ctx, outletID, client.NatsJetstreamJWTSeedFeatureRequest{
				UserJWT: planJWTSeed.UserJWT.ValueString(), NkeySeed: planJWTSeed.NkeySeed.ValueString(),
			}); err != nil {
				return err
			}
		}
	}

	if !nkeyConfigured {
		if !stateNkey.IsNull() && !stateNkey.IsUnknown() {
			if err := r.deleteNkey(ctx, outletID); err != nil {
				return err
			}
		}
	} else if stateNkey.IsNull() || stateNkey.IsUnknown() {
		if err := r.client.ConfigureNatsJetstreamNkeyFeature(ctx, outletID, client.NatsJetstreamNkeyFeatureRequest{
			Filename: planNkey.Filename.ValueString(), ContentBase64: planNkey.ContentBase64.ValueString(),
		}); err != nil {
			return err
		}
	} else {
		stateModel, diags := expandNatsJetstreamNkey(ctx, stateNkey)
		if diags.HasError() {
			return errors.New(diags.Errors()[0].Detail())
		}
		if !planNkey.Version.Equal(stateModel.Version) || !planNkey.Filename.Equal(stateModel.Filename) {
			if err := r.client.UpdateNatsJetstreamNkeyFeature(ctx, outletID, client.NatsJetstreamNkeyFeatureRequest{
				Filename: planNkey.Filename.ValueString(), ContentBase64: planNkey.ContentBase64.ValueString(),
			}); err != nil {
				return err
			}
		}
	}

	if !userCredsConfigured {
		if !stateUserCredentials.IsNull() && !stateUserCredentials.IsUnknown() {
			if err := r.deleteUserCredentials(ctx, outletID); err != nil {
				return err
			}
		}
	} else if stateUserCredentials.IsNull() || stateUserCredentials.IsUnknown() {
		if err := r.client.ConfigureNatsJetstreamUserCredentialsFeature(ctx, outletID, client.NatsJetstreamUserCredentialsFeatureRequest{
			Filename: planUserCredentials.Filename.ValueString(), ContentBase64: planUserCredentials.ContentBase64.ValueString(),
		}); err != nil {
			return err
		}
	} else {
		stateModel, diags := expandNatsJetstreamUserCredentials(ctx, stateUserCredentials)
		if diags.HasError() {
			return errors.New(diags.Errors()[0].Detail())
		}
		if !planUserCredentials.Version.Equal(stateModel.Version) || !planUserCredentials.Filename.Equal(stateModel.Filename) {
			if err := r.client.UpdateNatsJetstreamUserCredentialsFeature(ctx, outletID, client.NatsJetstreamUserCredentialsFeatureRequest{
				Filename: planUserCredentials.Filename.ValueString(), ContentBase64: planUserCredentials.ContentBase64.ValueString(),
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *natsJetstreamOutletResource) buildState(ctx context.Context, outletID int64, previousCredentials, previousJWTSeed, previousNkey, previousUserCredentials types.Object) (natsJetstreamOutletResourceModel, bool, diag.Diagnostics) {
	outlet, err := r.client.GetNatsJetstreamOutlet(ctx, outletID)
	if err != nil {
		if client.IsNotFound(err) {
			return natsJetstreamOutletResourceModel{}, true, nil
		}
		var diags diag.Diagnostics
		diags.AddError("Unable to Read NATS JetStream Outlet", err.Error())
		return natsJetstreamOutletResourceModel{}, false, diags
	}
	credentials, err := r.client.ListNatsJetstreamCredentials(ctx, outletID)
	if err != nil {
		if !client.IsNotFound(err) {
			var diags diag.Diagnostics
			diags.AddError("Unable to Read NATS JetStream Credentials", err.Error())
			return natsJetstreamOutletResourceModel{}, false, diags
		}
		credentials = nil
	}
	jwtSeedFeature, err := r.client.GetNatsJetstreamJWTSeedFeature(ctx, outletID)
	var jwtSeedFeaturePtr *client.OutletFeatureDetail
	if err != nil && !client.IsNotFound(err) {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read NATS JetStream JWT Seed Feature", err.Error())
		return natsJetstreamOutletResourceModel{}, false, diags
	}
	if err == nil && jwtSeedFeature.Feature != "" {
		jwtSeedFeaturePtr = &jwtSeedFeature
	}
	nkeyFeature, err := r.client.GetNatsJetstreamNkeyFeature(ctx, outletID)
	var nkeyFeaturePtr *client.OutletFeatureDetail
	if err != nil && !client.IsNotFound(err) {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read NATS JetStream NKEY Feature", err.Error())
		return natsJetstreamOutletResourceModel{}, false, diags
	}
	if err == nil && nkeyFeature.Feature != "" {
		nkeyFeaturePtr = &nkeyFeature
	}
	userCredentialsFeature, err := r.client.GetNatsJetstreamUserCredentialsFeature(ctx, outletID)
	var userCredentialsFeaturePtr *client.OutletFeatureDetail
	if err != nil && !client.IsNotFound(err) {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read NATS JetStream User Credentials Feature", err.Error())
		return natsJetstreamOutletResourceModel{}, false, diags
	}
	if err == nil && userCredentialsFeature.Feature != "" {
		userCredentialsFeaturePtr = &userCredentialsFeature
	}

	credentialsObject, diags := buildAmqpURLsCredentialsObject(ctx, previousCredentials, credentials, natsJetstreamURLsCredentialRole)
	if diags.HasError() {
		return natsJetstreamOutletResourceModel{}, false, diags
	}
	jwtSeedObject, diags := buildNatsJetstreamJWTSeedObject(ctx, previousJWTSeed, jwtSeedFeaturePtr, credentials)
	if diags.HasError() {
		return natsJetstreamOutletResourceModel{}, false, diags
	}
	nkeyObject, diags := buildNatsJetstreamNkeyObject(ctx, previousNkey, nkeyFeaturePtr, credentials)
	if diags.HasError() {
		return natsJetstreamOutletResourceModel{}, false, diags
	}
	userCredentialsObject, diags := buildNatsJetstreamUserCredentialsObject(ctx, previousUserCredentials, userCredentialsFeaturePtr, credentials)
	if diags.HasError() {
		return natsJetstreamOutletResourceModel{}, false, diags
	}

	transform, err := r.client.GetTypedOutletTransform(ctx, "nats_jetstream", outletID)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read NATS JetStream Transform", err.Error())
		return natsJetstreamOutletResourceModel{}, false, diags
	}
	return natsJetstreamOutletModelFromAPI(outlet, credentialsObject, jwtSeedObject, nkeyObject, userCredentialsObject, transform), false, nil
}

func natsJetstreamOutletModelFromAPI(outlet client.NatsJetstreamOutlet, credentials, jwtSeed, nkey, userCredentials types.Object, transform client.OutletTransformState) natsJetstreamOutletResourceModel {
	return natsJetstreamOutletResourceModel{
		ID:                      types.Int64Value(outlet.ID),
		Name:                    types.StringValue(outlet.Name),
		Description:             nullableString(outlet.Description),
		Enabled:                 types.BoolValue(outlet.Enabled),
		BackoffInitialInterval:  nullableStringPointer(outlet.BackoffInitialInterval),
		BackoffJitter:           nullableFloat64Pointer(outlet.BackoffJitter),
		BackoffMaxElapsedTime:   nullableStringPointer(outlet.BackoffMaxElapsedTime),
		BackoffMaxInterval:      nullableStringPointer(outlet.BackoffMaxInterval),
		Headers:                 nullableStringMap(outlet.Headers),
		InjectTracingMap:        nullableStringPointer(outlet.InjectTracingMap),
		MaxInFlight:             nullableInt64Pointer(outlet.MaxInFlight),
		MaxRetries:              nullableInt64Pointer(outlet.MaxRetries),
		MetadataIncludePatterns: nullableStringList(outlet.MetadataIncludePatterns),
		MetadataIncludePrefixes: nullableStringList(outlet.MetadataIncludePrefixes),
		NotifyPolicy:            nullableStringPointerPreservingEmpty(outlet.NotifyPolicy),
		Subject:                 types.StringValue(outlet.Subject),
		TLSEnableRenegotiation:  nullableBoolPointer(outlet.TLSEnableRenegotiation),
		TLSSkipCertVerify:       nullableBoolPointer(outlet.TLSSkipCertVerify),
		Credentials:             credentials,
		JWTSeed:                 jwtSeed,
		Nkey:                    nkey,
		UserCredentials:         userCredentials,
		ApplyStatus:             nullableString(outlet.ApplyStatus),
		LastApplyError:          nullableString(outlet.LastApplyError),
		ConfigGeneration:        types.Int64Value(outlet.ConfigGeneration),
		AppliedGeneration:       types.Int64Value(outlet.AppliedGeneration),
		Transform:               transformObject(transform),
	}
}

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

var (
	_ datasource.DataSource              = &natsJetstreamOutletDataSource{}
	_ datasource.DataSourceWithConfigure = &natsJetstreamOutletDataSource{}
)

type natsJetstreamOutletDataSource struct{ client natsJetstreamOutletClient }

type natsJetstreamOutletDataSourceModel struct {
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
	JWTSeed                 types.Object  `tfsdk:"jwt_seed"`
	Nkey                    types.Object  `tfsdk:"nkey"`
	UserCredentials         types.Object  `tfsdk:"user_credentials"`
	ApplyStatus             types.String  `tfsdk:"apply_status"`
	LastApplyError          types.String  `tfsdk:"last_apply_error"`
	ConfigGeneration        types.Int64   `tfsdk:"config_generation"`
	AppliedGeneration       types.Int64   `tfsdk:"applied_generation"`
	Transform               types.Object  `tfsdk:"transform"`
}

func NewNatsJetstreamOutletDataSource() datasource.DataSource {
	return &natsJetstreamOutletDataSource{}
}

func (d *natsJetstreamOutletDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nats_jetstream_outlet"
}

func (d *natsJetstreamOutletDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Reads a Noozle NATS JetStream outlet by ID.",
		Attributes: map[string]datasourceschema.Attribute{
			"id":                        datasourceschema.Int64Attribute{Required: true},
			"name":                      datasourceschema.StringAttribute{Computed: true},
			"description":               datasourceschema.StringAttribute{Computed: true},
			"enabled":                   datasourceschema.BoolAttribute{Computed: true},
			"backoff_initial_interval":  datasourceschema.StringAttribute{Computed: true},
			"backoff_jitter":            datasourceschema.Float64Attribute{Computed: true},
			"backoff_max_elapsed_time":  datasourceschema.StringAttribute{Computed: true},
			"backoff_max_interval":      datasourceschema.StringAttribute{Computed: true},
			"headers":                   datasourceschema.MapAttribute{Computed: true, ElementType: types.StringType},
			"inject_tracing_map":        datasourceschema.StringAttribute{Computed: true},
			"max_in_flight":             datasourceschema.Int64Attribute{Computed: true},
			"max_retries":               datasourceschema.Int64Attribute{Computed: true},
			"metadata_include_patterns": datasourceschema.ListAttribute{Computed: true, ElementType: types.StringType},
			"metadata_include_prefixes": datasourceschema.ListAttribute{Computed: true, ElementType: types.StringType},
			"notify_policy":             datasourceschema.StringAttribute{Computed: true},
			"subject":                   datasourceschema.StringAttribute{Computed: true},
			"tls_enable_renegotiation":  datasourceschema.BoolAttribute{Computed: true},
			"tls_skip_cert_verify":      datasourceschema.BoolAttribute{Computed: true},
			"apply_status":              datasourceschema.StringAttribute{Computed: true},
			"last_apply_error":          datasourceschema.StringAttribute{Computed: true},
			"config_generation":         datasourceschema.Int64Attribute{Computed: true},
			"applied_generation":        datasourceschema.Int64Attribute{Computed: true},
			"transform": datasourceschema.SingleNestedAttribute{Computed: true, Attributes: map[string]datasourceschema.Attribute{
				"selected": datasourceschema.StringAttribute{Computed: true}, "effective": datasourceschema.StringAttribute{Computed: true}, "state": datasourceschema.StringAttribute{Computed: true}, "applied": datasourceschema.BoolAttribute{Computed: true}, "generation": datasourceschema.Int64Attribute{Computed: true}, "last_error": datasourceschema.StringAttribute{Computed: true},
			}},
		},
		Blocks: map[string]datasourceschema.Block{
			"jwt_seed": datasourceschema.SingleNestedBlock{Attributes: map[string]datasourceschema.Attribute{
				"enabled": datasourceschema.BoolAttribute{Computed: true}, "state": datasourceschema.StringAttribute{Computed: true}, "desired_generation": datasourceschema.Int64Attribute{Computed: true}, "applied_generation": datasourceschema.Int64Attribute{Computed: true}, "last_error": datasourceschema.StringAttribute{Computed: true}, "user_jwt_version": datasourceschema.Int64Attribute{Computed: true}, "nkey_seed_version": datasourceschema.Int64Attribute{Computed: true},
			}},
			"nkey": datasourceschema.SingleNestedBlock{Attributes: map[string]datasourceschema.Attribute{
				"enabled": datasourceschema.BoolAttribute{Computed: true}, "state": datasourceschema.StringAttribute{Computed: true}, "desired_generation": datasourceschema.Int64Attribute{Computed: true}, "applied_generation": datasourceschema.Int64Attribute{Computed: true}, "last_error": datasourceschema.StringAttribute{Computed: true}, "filename": datasourceschema.StringAttribute{Computed: true}, "checksum_sha256": datasourceschema.StringAttribute{Computed: true},
			}},
			"user_credentials": datasourceschema.SingleNestedBlock{Attributes: map[string]datasourceschema.Attribute{
				"enabled": datasourceschema.BoolAttribute{Computed: true}, "state": datasourceschema.StringAttribute{Computed: true}, "desired_generation": datasourceschema.Int64Attribute{Computed: true}, "applied_generation": datasourceschema.Int64Attribute{Computed: true}, "last_error": datasourceschema.StringAttribute{Computed: true}, "filename": datasourceschema.StringAttribute{Computed: true}, "checksum_sha256": datasourceschema.StringAttribute{Computed: true},
			}},
		},
	}
}

func (d *natsJetstreamOutletDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	outletClient, ok := req.ProviderData.(natsJetstreamOutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected NATS JetStream outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}
	d.client = outletClient
}

func (d *natsJetstreamOutletDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var config natsJetstreamOutletDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	outlet, err := d.client.GetNatsJetstreamOutlet(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read NATS JetStream Outlet", err.Error())
		return
	}
	credentials, err := d.client.ListNatsJetstreamCredentials(ctx, config.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Read NATS JetStream Credentials", err.Error())
		return
	}
	var credentialsList []client.OutletMaterialSummary
	if err == nil {
		credentialsList = credentials
	}
	jwtSeedFeature, err := d.client.GetNatsJetstreamJWTSeedFeature(ctx, config.ID.ValueInt64())
	var jwtSeedFeaturePtr *client.OutletFeatureDetail
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Read NATS JetStream JWT Seed Feature", err.Error())
			return
		}
	} else if jwtSeedFeature.Feature != "" {
		jwtSeedFeaturePtr = &jwtSeedFeature
	}
	nkeyFeature, err := d.client.GetNatsJetstreamNkeyFeature(ctx, config.ID.ValueInt64())
	var nkeyFeaturePtr *client.OutletFeatureDetail
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Read NATS JetStream NKEY Feature", err.Error())
			return
		}
	} else if nkeyFeature.Feature != "" {
		nkeyFeaturePtr = &nkeyFeature
	}
	userCredentialsFeature, err := d.client.GetNatsJetstreamUserCredentialsFeature(ctx, config.ID.ValueInt64())
	var userCredentialsFeaturePtr *client.OutletFeatureDetail
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Read NATS JetStream User Credentials Feature", err.Error())
			return
		}
	} else if userCredentialsFeature.Feature != "" {
		userCredentialsFeaturePtr = &userCredentialsFeature
	}

	jwtSeedObject, diags := buildNatsJetstreamJWTSeedDataSourceObject(jwtSeedFeaturePtr)
	resp.Diagnostics.Append(diags...)
	nkeyObject, diags := buildNatsJetstreamNkeyDataSourceObject(nkeyFeaturePtr, credentialsList)
	resp.Diagnostics.Append(diags...)
	userCredentialsObject, diags := buildNatsJetstreamUserCredentialsDataSourceObject(userCredentialsFeaturePtr, credentialsList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	transform, err := d.client.GetTypedOutletTransform(ctx, "nats_jetstream", config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read NATS JetStream Transform", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &natsJetstreamOutletDataSourceModel{
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
		JWTSeed:                 jwtSeedObject,
		Nkey:                    nkeyObject,
		UserCredentials:         userCredentialsObject,
		ApplyStatus:             nullableString(outlet.ApplyStatus),
		LastApplyError:          nullableString(outlet.LastApplyError),
		ConfigGeneration:        types.Int64Value(outlet.ConfigGeneration),
		AppliedGeneration:       types.Int64Value(outlet.AppliedGeneration),
		Transform:               transformObject(transform),
	})...)
}

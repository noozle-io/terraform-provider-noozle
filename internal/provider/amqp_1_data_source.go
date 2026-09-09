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
	_ datasource.DataSource              = &amqp1OutletDataSource{}
	_ datasource.DataSourceWithConfigure = &amqp1OutletDataSource{}
)

type amqp1OutletDataSource struct{ client amqp1OutletClient }

type amqp1OutletDataSourceModel struct {
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

func NewAmqp1OutletDataSource() datasource.DataSource { return &amqp1OutletDataSource{} }

func (d *amqp1OutletDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_amqp_1_outlet"
}

func (d *amqp1OutletDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Reads a Noozle AMQP 1.0 outlet by ID.",
		Attributes: map[string]datasourceschema.Attribute{
			"id":                         datasourceschema.Int64Attribute{Required: true},
			"name":                       datasourceschema.StringAttribute{Computed: true},
			"description":                datasourceschema.StringAttribute{Computed: true},
			"enabled":                    datasourceschema.BoolAttribute{Computed: true},
			"application_properties_map": datasourceschema.StringAttribute{Computed: true},
			"max_in_flight":              datasourceschema.Int64Attribute{Computed: true},
			"metadata_exclude_prefixes":  datasourceschema.ListAttribute{Computed: true, ElementType: types.StringType},
			"notify_policy":              datasourceschema.StringAttribute{Computed: true},
			"target_address":             datasourceschema.StringAttribute{Computed: true},
			"tls_enable_renegotiation":   datasourceschema.BoolAttribute{Computed: true},
			"tls_skip_cert_verify":       datasourceschema.BoolAttribute{Computed: true},
			"apply_status":               datasourceschema.StringAttribute{Computed: true},
			"last_apply_error":           datasourceschema.StringAttribute{Computed: true},
			"config_generation":          datasourceschema.Int64Attribute{Computed: true},
			"applied_generation":         datasourceschema.Int64Attribute{Computed: true},
			"transform": datasourceschema.SingleNestedAttribute{Computed: true, Attributes: map[string]datasourceschema.Attribute{
				"selected": datasourceschema.StringAttribute{Computed: true}, "effective": datasourceschema.StringAttribute{Computed: true}, "state": datasourceschema.StringAttribute{Computed: true}, "applied": datasourceschema.BoolAttribute{Computed: true}, "generation": datasourceschema.Int64Attribute{Computed: true}, "last_error": datasourceschema.StringAttribute{Computed: true},
			}},
		},
		Blocks: map[string]datasourceschema.Block{
			"tls": datasourceschema.SingleNestedBlock{Attributes: map[string]datasourceschema.Attribute{
				"enabled": datasourceschema.BoolAttribute{Computed: true}, "state": datasourceschema.StringAttribute{Computed: true}, "desired_generation": datasourceschema.Int64Attribute{Computed: true}, "applied_generation": datasourceschema.Int64Attribute{Computed: true}, "last_error": datasourceschema.StringAttribute{Computed: true}, "ca_cert_filename": datasourceschema.StringAttribute{Computed: true}, "ca_cert_checksum_sha256": datasourceschema.StringAttribute{Computed: true},
			}},
			"mtls": datasourceschema.SingleNestedBlock{Attributes: map[string]datasourceschema.Attribute{
				"enabled": datasourceschema.BoolAttribute{Computed: true}, "state": datasourceschema.StringAttribute{Computed: true}, "desired_generation": datasourceschema.Int64Attribute{Computed: true}, "applied_generation": datasourceschema.Int64Attribute{Computed: true}, "last_error": datasourceschema.StringAttribute{Computed: true}, "ca_cert_filename": datasourceschema.StringAttribute{Computed: true}, "ca_cert_checksum_sha256": datasourceschema.StringAttribute{Computed: true}, "client_cert_filename": datasourceschema.StringAttribute{Computed: true}, "client_cert_checksum_sha256": datasourceschema.StringAttribute{Computed: true}, "client_key_filename": datasourceschema.StringAttribute{Computed: true}, "client_key_checksum_sha256": datasourceschema.StringAttribute{Computed: true},
			}},
			"sasl_anonymous": datasourceschema.SingleNestedBlock{Attributes: map[string]datasourceschema.Attribute{
				"enabled": datasourceschema.BoolAttribute{Computed: true}, "state": datasourceschema.StringAttribute{Computed: true}, "desired_generation": datasourceschema.Int64Attribute{Computed: true}, "applied_generation": datasourceschema.Int64Attribute{Computed: true}, "last_error": datasourceschema.StringAttribute{Computed: true},
			}},
			"sasl_plain": datasourceschema.SingleNestedBlock{Attributes: map[string]datasourceschema.Attribute{
				"enabled": datasourceschema.BoolAttribute{Computed: true}, "state": datasourceschema.StringAttribute{Computed: true}, "desired_generation": datasourceschema.Int64Attribute{Computed: true}, "applied_generation": datasourceschema.Int64Attribute{Computed: true}, "last_error": datasourceschema.StringAttribute{Computed: true}, "user": datasourceschema.StringAttribute{Computed: true}, "password_version": datasourceschema.Int64Attribute{Computed: true},
			}},
		},
	}
}

func (d *amqp1OutletDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	outletClient, ok := req.ProviderData.(amqp1OutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected AMQP 1.0 outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}
	d.client = outletClient
}

func (d *amqp1OutletDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var config amqp1OutletDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	outlet, err := d.client.GetAmqp1Outlet(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read AMQP 1.0 Outlet", err.Error())
		return
	}
	credentials, err := d.client.ListAmqp1Credentials(ctx, config.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Read AMQP 1.0 Credentials", err.Error())
		return
	}
	var credentialsList []client.OutletMaterialSummary
	if err == nil {
		credentialsList = credentials
	}
	tlsFeature, err := d.client.GetAmqp1TLSFeature(ctx, config.ID.ValueInt64())
	var tlsFeaturePtr *client.OutletFeatureDetail
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Read AMQP 1.0 TLS Feature", err.Error())
			return
		}
	} else if tlsFeature.Feature != "" {
		tlsFeaturePtr = &tlsFeature
	}
	mtlsFeature, err := d.client.GetAmqp1MTLSFeature(ctx, config.ID.ValueInt64())
	var mtlsFeaturePtr *client.OutletFeatureDetail
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Read AMQP 1.0 mTLS Feature", err.Error())
			return
		}
	} else if mtlsFeature.Feature != "" {
		mtlsFeaturePtr = &mtlsFeature
	}
	saslAnonymousFeature, err := d.client.GetAmqp1SASLAnonymousFeature(ctx, config.ID.ValueInt64())
	var saslAnonymousFeaturePtr *client.OutletFeatureDetail
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Read AMQP 1.0 SASL Anonymous Feature", err.Error())
			return
		}
	} else if saslAnonymousFeature.Feature != "" {
		saslAnonymousFeaturePtr = &saslAnonymousFeature
	}
	saslPlainFeature, err := d.client.GetAmqp1SASLPlainFeature(ctx, config.ID.ValueInt64())
	var saslPlainFeaturePtr *client.OutletFeatureDetail
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Read AMQP 1.0 SASL Plain Feature", err.Error())
			return
		}
	} else if saslPlainFeature.Feature != "" {
		saslPlainFeaturePtr = &saslPlainFeature
	}
	tlsObject, diags := buildAMQPTLSDataSourceObject(tlsFeaturePtr, credentialsList)
	resp.Diagnostics.Append(diags...)
	mtlsObject, diags := buildAMQPMTLSDataSourceObject(mtlsFeaturePtr, credentialsList)
	resp.Diagnostics.Append(diags...)
	saslAnonymousObject, diags := buildAMQPSASLAnonymousDataSourceObject(saslAnonymousFeaturePtr)
	resp.Diagnostics.Append(diags...)
	saslPlainObject, diags := buildAMQPSASLPlainDataSourceObject(saslPlainFeaturePtr, credentialsList)
	resp.Diagnostics.Append(diags...)
	transform, err := d.client.GetTypedOutletTransform(ctx, "amqp_1", config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read AMQP 1.0 Transform", err.Error())
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &amqp1OutletDataSourceModel{
		ID: types.Int64Value(outlet.ID), Name: types.StringValue(outlet.Name), Description: nullableString(outlet.Description), Enabled: types.BoolValue(outlet.Enabled),
		ApplicationPropertiesMap: nullableStringPointer(outlet.ApplicationPropertiesMap), MaxInFlight: nullableInt64Pointer(outlet.MaxInFlight), MetadataExcludePrefixes: nullableStringList(outlet.MetadataExcludePrefixes), NotifyPolicy: nullableStringPointerPreservingEmpty(outlet.NotifyPolicy),
		TargetAddress: types.StringValue(outlet.TargetAddress), TLSEnableRenegotiation: nullableBoolPointer(outlet.TLSEnableRenegotiation), TLSSkipCertVerify: nullableBoolPointer(outlet.TLSSkipCertVerify),
		TLS: tlsObject, MTLS: mtlsObject, SASLAnonymous: saslAnonymousObject, SASLPlain: saslPlainObject, ApplyStatus: nullableString(outlet.ApplyStatus), LastApplyError: nullableString(outlet.LastApplyError), ConfigGeneration: types.Int64Value(outlet.ConfigGeneration), AppliedGeneration: types.Int64Value(outlet.AppliedGeneration),
		Transform: transformObject(transform),
	})...)
}

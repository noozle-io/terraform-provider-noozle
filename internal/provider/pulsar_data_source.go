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
	_ datasource.DataSource              = &pulsarOutletDataSource{}
	_ datasource.DataSourceWithConfigure = &pulsarOutletDataSource{}
)

type pulsarOutletDataSource struct {
	client pulsarOutletClient
}

type pulsarOutletDataSourceModel struct {
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

func NewPulsarOutletDataSource() datasource.DataSource { return &pulsarOutletDataSource{} }

func (d *pulsarOutletDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pulsar_outlet"
}

func (d *pulsarOutletDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Reads a Noozle Pulsar outlet by ID.",
		Attributes: map[string]datasourceschema.Attribute{
			"id":                 datasourceschema.Int64Attribute{MarkdownDescription: "Unique identifier of the outlet.", Required: true},
			"name":               datasourceschema.StringAttribute{MarkdownDescription: "Name presented for the outlet.", Computed: true},
			"description":        datasourceschema.StringAttribute{MarkdownDescription: "Optional note explaining the outlet purpose.", Computed: true},
			"enabled":            datasourceschema.BoolAttribute{MarkdownDescription: "Whether the outlet is active.", Computed: true},
			"key":                datasourceschema.StringAttribute{MarkdownDescription: "Pulsar message key interpolation.", Computed: true},
			"max_in_flight":      datasourceschema.Int64Attribute{MarkdownDescription: "Maximum number of in-flight Pulsar messages.", Computed: true},
			"notify_policy":      datasourceschema.StringAttribute{MarkdownDescription: "Delivery policy applied when notifying this outlet.", Computed: true},
			"ordering_key":       datasourceschema.StringAttribute{MarkdownDescription: "Pulsar ordering key interpolation.", Computed: true},
			"topic":              datasourceschema.StringAttribute{MarkdownDescription: "Pulsar topic.", Computed: true},
			"url":                datasourceschema.StringAttribute{MarkdownDescription: "Pulsar broker URL.", Computed: true},
			"apply_status":       datasourceschema.StringAttribute{MarkdownDescription: "Observed apply/runtime status for the outlet.", Computed: true},
			"last_apply_error":   datasourceschema.StringAttribute{MarkdownDescription: "Last observed apply/runtime error message.", Computed: true},
			"config_generation":  datasourceschema.Int64Attribute{MarkdownDescription: "Desired outlet config generation.", Computed: true},
			"applied_generation": datasourceschema.Int64Attribute{MarkdownDescription: "Last observed applied config generation.", Computed: true},
			"transform": datasourceschema.SingleNestedAttribute{Computed: true, Attributes: map[string]datasourceschema.Attribute{
				"selected": datasourceschema.StringAttribute{Computed: true}, "effective": datasourceschema.StringAttribute{Computed: true}, "state": datasourceschema.StringAttribute{Computed: true}, "applied": datasourceschema.BoolAttribute{Computed: true}, "generation": datasourceschema.Int64Attribute{Computed: true}, "last_error": datasourceschema.StringAttribute{Computed: true},
			}},
		},
		Blocks: map[string]datasourceschema.Block{
			"tls": datasourceschema.SingleNestedBlock{
				MarkdownDescription: "Observed TLS feature state for the Pulsar outlet.",
				Attributes: map[string]datasourceschema.Attribute{
					"enabled":                 datasourceschema.BoolAttribute{MarkdownDescription: "Whether the TLS feature is enabled.", Computed: true},
					"state":                   datasourceschema.StringAttribute{MarkdownDescription: "Observed TLS feature state.", Computed: true},
					"desired_generation":      datasourceschema.Int64Attribute{MarkdownDescription: "Desired TLS feature generation.", Computed: true},
					"applied_generation":      datasourceschema.Int64Attribute{MarkdownDescription: "Applied TLS feature generation.", Computed: true},
					"last_error":              datasourceschema.StringAttribute{MarkdownDescription: "Last observed TLS feature error.", Computed: true},
					"ca_cert_filename":        datasourceschema.StringAttribute{MarkdownDescription: "Observed TLS CA certificate filename.", Computed: true},
					"ca_cert_checksum_sha256": datasourceschema.StringAttribute{MarkdownDescription: "Observed checksum of the TLS CA certificate file.", Computed: true},
				},
			},
			"oauth2": datasourceschema.SingleNestedBlock{
				MarkdownDescription: "Observed OAuth2 feature state for the Pulsar outlet.",
				Attributes: map[string]datasourceschema.Attribute{
					"enabled":                     datasourceschema.BoolAttribute{MarkdownDescription: "Whether the OAuth2 feature is enabled.", Computed: true},
					"state":                       datasourceschema.StringAttribute{MarkdownDescription: "Observed OAuth2 feature state.", Computed: true},
					"desired_generation":          datasourceschema.Int64Attribute{MarkdownDescription: "Desired OAuth2 feature generation.", Computed: true},
					"applied_generation":          datasourceschema.Int64Attribute{MarkdownDescription: "Applied OAuth2 feature generation.", Computed: true},
					"last_error":                  datasourceschema.StringAttribute{MarkdownDescription: "Last observed OAuth2 feature error.", Computed: true},
					"issuer_url":                  datasourceschema.StringAttribute{MarkdownDescription: "Observed OAuth2 issuer URL.", Computed: true},
					"audience":                    datasourceschema.StringAttribute{MarkdownDescription: "Observed OAuth2 audience.", Computed: true},
					"private_key_filename":        datasourceschema.StringAttribute{MarkdownDescription: "Observed OAuth2 private key filename.", Computed: true},
					"private_key_checksum_sha256": datasourceschema.StringAttribute{MarkdownDescription: "Observed checksum of the OAuth2 private key file.", Computed: true},
				},
			},
			"token": datasourceschema.SingleNestedBlock{
				MarkdownDescription: "Observed token feature state for the Pulsar outlet.",
				Attributes: map[string]datasourceschema.Attribute{
					"enabled":              datasourceschema.BoolAttribute{MarkdownDescription: "Whether the token feature is enabled.", Computed: true},
					"state":                datasourceschema.StringAttribute{MarkdownDescription: "Observed token feature state.", Computed: true},
					"desired_generation":   datasourceschema.Int64Attribute{MarkdownDescription: "Desired token feature generation.", Computed: true},
					"applied_generation":   datasourceschema.Int64Attribute{MarkdownDescription: "Applied token feature generation.", Computed: true},
					"last_error":           datasourceschema.StringAttribute{MarkdownDescription: "Last observed token feature error.", Computed: true},
					"access_token_version": datasourceschema.Int64Attribute{MarkdownDescription: "Configured Pulsar access token version.", Computed: true},
				},
			},
		},
	}
}

func (d *pulsarOutletDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	outletClient, ok := req.ProviderData.(pulsarOutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected Pulsar outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}
	d.client = outletClient
}

func (d *pulsarOutletDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var config pulsarOutletDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	outlet, err := d.client.GetPulsarOutlet(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Pulsar Outlet", err.Error())
		return
	}
	credentials, err := d.client.ListPulsarCredentials(ctx, config.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Read Pulsar Credentials", err.Error())
		return
	}
	var credentialsList []client.OutletMaterialSummary
	if err == nil {
		credentialsList = credentials
	}
	tlsFeature, err := d.client.GetPulsarTLSFeature(ctx, config.ID.ValueInt64())
	var tlsFeaturePtr *client.OutletFeatureDetail
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Read Pulsar TLS Feature", err.Error())
			return
		}
	} else if tlsFeature.Feature != "" {
		tlsFeaturePtr = &tlsFeature
	}
	oauth2Feature, err := d.client.GetPulsarOAuth2Feature(ctx, config.ID.ValueInt64())
	var oauth2FeaturePtr *client.OutletFeatureDetail
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Read Pulsar OAuth2 Feature", err.Error())
			return
		}
	} else if oauth2Feature.Feature != "" {
		oauth2FeaturePtr = &oauth2Feature
	}
	tokenFeature, err := d.client.GetPulsarTokenFeature(ctx, config.ID.ValueInt64())
	var tokenFeaturePtr *client.OutletFeatureDetail
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Read Pulsar Token Feature", err.Error())
			return
		}
	} else if tokenFeature.Feature != "" {
		tokenFeaturePtr = &tokenFeature
	}
	tlsObject, diags := buildPulsarTLSDataSourceObject(tlsFeaturePtr, credentialsList)
	resp.Diagnostics.Append(diags...)
	oauth2Object, diags := buildPulsarOAuth2DataSourceObject(oauth2FeaturePtr, credentialsList)
	resp.Diagnostics.Append(diags...)
	tokenObject, diags := buildPulsarTokenDataSourceObject(tokenFeaturePtr, credentialsList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	transform, err := d.client.GetTypedOutletTransform(ctx, "pulsar", config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Pulsar Transform", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &pulsarOutletDataSourceModel{
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
		TLS:               tlsObject,
		OAuth2:            oauth2Object,
		Token:             tokenObject,
		ApplyStatus:       nullableString(outlet.ApplyStatus),
		LastApplyError:    nullableString(outlet.LastApplyError),
		ConfigGeneration:  types.Int64Value(outlet.ConfigGeneration),
		AppliedGeneration: types.Int64Value(outlet.AppliedGeneration),
		Transform:         transformObject(transform),
	})...)
}

func buildPulsarTLSDataSourceObject(feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	credential := findMaterialOrNil(credentials, pulsarTLSCARole)
	attrTypes := map[string]attr.Type{
		"enabled":                 types.BoolType,
		"state":                   types.StringType,
		"desired_generation":      types.Int64Type,
		"applied_generation":      types.Int64Type,
		"last_error":              types.StringType,
		"ca_cert_filename":        types.StringType,
		"ca_cert_checksum_sha256": types.StringType,
	}
	if feature == nil && credential == nil {
		return types.ObjectNull(attrTypes), nil
	}
	values := buildPulsarFeatureBase(feature)
	values["ca_cert_filename"] = nullableString(materialField(credential, func(m *client.OutletMaterialSummary) string { return m.Filename }))
	values["ca_cert_checksum_sha256"] = nullableString(materialField(credential, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 }))
	return types.ObjectValue(attrTypes, values)
}

func buildPulsarOAuth2DataSourceObject(feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	credential := findMaterialOrNil(credentials, pulsarOAuth2PrivateKeyRole)
	attrTypes := map[string]attr.Type{
		"enabled":                     types.BoolType,
		"state":                       types.StringType,
		"desired_generation":          types.Int64Type,
		"applied_generation":          types.Int64Type,
		"last_error":                  types.StringType,
		"issuer_url":                  types.StringType,
		"audience":                    types.StringType,
		"private_key_filename":        types.StringType,
		"private_key_checksum_sha256": types.StringType,
	}
	if feature == nil && credential == nil {
		return types.ObjectNull(attrTypes), nil
	}
	values := buildPulsarFeatureBase(feature)
	values["issuer_url"] = types.StringNull()
	values["audience"] = types.StringNull()
	values["private_key_filename"] = nullableString(materialField(credential, func(m *client.OutletMaterialSummary) string { return m.Filename }))
	values["private_key_checksum_sha256"] = nullableString(materialField(credential, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 }))
	return types.ObjectValue(attrTypes, values)
}

func buildPulsarTokenDataSourceObject(feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	credential := findMaterialOrNil(credentials, pulsarTokenRole)
	attrTypes := map[string]attr.Type{
		"enabled":              types.BoolType,
		"state":                types.StringType,
		"desired_generation":   types.Int64Type,
		"applied_generation":   types.Int64Type,
		"last_error":           types.StringType,
		"access_token_version": types.Int64Type,
	}
	if feature == nil && credential == nil {
		return types.ObjectNull(attrTypes), nil
	}
	values := buildPulsarFeatureBase(feature)
	values["access_token_version"] = types.Int64Null()
	return types.ObjectValue(attrTypes, values)
}

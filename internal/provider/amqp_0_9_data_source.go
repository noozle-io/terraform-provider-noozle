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
	_ datasource.DataSource              = &amqp09OutletDataSource{}
	_ datasource.DataSourceWithConfigure = &amqp09OutletDataSource{}
)

type amqp09OutletDataSource struct{ client amqp09OutletClient }

type amqp09OutletDataSourceModel struct {
	ID                      types.Int64  `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
	AppID                   types.String `tfsdk:"app_id"`
	ContentEncoding         types.String `tfsdk:"content_encoding"`
	ContentType             types.String `tfsdk:"content_type"`
	CorrelationID           types.String `tfsdk:"correlation_id"`
	Exchange                types.String `tfsdk:"exchange"`
	ExchangeDeclareDurable  types.Bool   `tfsdk:"exchange_declare_durable"`
	ExchangeDeclareEnabled  types.Bool   `tfsdk:"exchange_declare_enabled"`
	ExchangeDeclareType     types.String `tfsdk:"exchange_declare_type"`
	Expiration              types.String `tfsdk:"expiration"`
	Immediate               types.Bool   `tfsdk:"immediate"`
	Key                     types.String `tfsdk:"key"`
	Mandatory               types.Bool   `tfsdk:"mandatory"`
	MaxInFlight             types.Int64  `tfsdk:"max_in_flight"`
	MessageID               types.String `tfsdk:"message_id"`
	MetadataExcludePrefixes types.List   `tfsdk:"metadata_exclude_prefixes"`
	NotifyPolicy            types.String `tfsdk:"notify_policy"`
	Persistent              types.Bool   `tfsdk:"persistent"`
	Priority                types.String `tfsdk:"priority"`
	ReplyTo                 types.String `tfsdk:"reply_to"`
	Timeout                 types.String `tfsdk:"timeout"`
	TLSEnableRenegotiation  types.Bool   `tfsdk:"tls_enable_renegotiation"`
	TLSSkipCertVerify       types.Bool   `tfsdk:"tls_skip_cert_verify"`
	AMQPType                types.String `tfsdk:"type"`
	UserID                  types.String `tfsdk:"user_id"`
	TLS                     types.Object `tfsdk:"tls"`
	MTLS                    types.Object `tfsdk:"mtls"`
	ApplyStatus             types.String `tfsdk:"apply_status"`
	LastApplyError          types.String `tfsdk:"last_apply_error"`
	ConfigGeneration        types.Int64  `tfsdk:"config_generation"`
	AppliedGeneration       types.Int64  `tfsdk:"applied_generation"`
	Transform               types.Object `tfsdk:"transform"`
}

func NewAmqp09OutletDataSource() datasource.DataSource { return &amqp09OutletDataSource{} }

func (d *amqp09OutletDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_amqp_0_9_outlet"
}

func (d *amqp09OutletDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Reads a Noozle AMQP 0.9.1 outlet by ID.",
		Attributes: map[string]datasourceschema.Attribute{
			"id":                        datasourceschema.Int64Attribute{Required: true},
			"name":                      datasourceschema.StringAttribute{Computed: true},
			"description":               datasourceschema.StringAttribute{Computed: true},
			"enabled":                   datasourceschema.BoolAttribute{Computed: true},
			"app_id":                    datasourceschema.StringAttribute{Computed: true},
			"content_encoding":          datasourceschema.StringAttribute{Computed: true},
			"content_type":              datasourceschema.StringAttribute{Computed: true},
			"correlation_id":            datasourceschema.StringAttribute{Computed: true},
			"exchange":                  datasourceschema.StringAttribute{Computed: true},
			"exchange_declare_durable":  datasourceschema.BoolAttribute{Computed: true},
			"exchange_declare_enabled":  datasourceschema.BoolAttribute{Computed: true},
			"exchange_declare_type":     datasourceschema.StringAttribute{Computed: true},
			"expiration":                datasourceschema.StringAttribute{Computed: true},
			"immediate":                 datasourceschema.BoolAttribute{Computed: true},
			"key":                       datasourceschema.StringAttribute{Computed: true},
			"mandatory":                 datasourceschema.BoolAttribute{Computed: true},
			"max_in_flight":             datasourceschema.Int64Attribute{Computed: true},
			"message_id":                datasourceschema.StringAttribute{Computed: true},
			"metadata_exclude_prefixes": datasourceschema.ListAttribute{Computed: true, ElementType: types.StringType},
			"notify_policy":             datasourceschema.StringAttribute{Computed: true},
			"persistent":                datasourceschema.BoolAttribute{Computed: true},
			"priority":                  datasourceschema.StringAttribute{Computed: true},
			"reply_to":                  datasourceschema.StringAttribute{Computed: true},
			"timeout":                   datasourceschema.StringAttribute{Computed: true},
			"tls_enable_renegotiation":  datasourceschema.BoolAttribute{Computed: true},
			"tls_skip_cert_verify":      datasourceschema.BoolAttribute{Computed: true},
			"type":                      datasourceschema.StringAttribute{Computed: true},
			"user_id":                   datasourceschema.StringAttribute{Computed: true},
			"apply_status":              datasourceschema.StringAttribute{Computed: true},
			"last_apply_error":          datasourceschema.StringAttribute{Computed: true},
			"config_generation":         datasourceschema.Int64Attribute{Computed: true},
			"applied_generation":        datasourceschema.Int64Attribute{Computed: true},
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
		},
	}
}

func (d *amqp09OutletDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	outletClient, ok := req.ProviderData.(amqp09OutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected AMQP 0.9.1 outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}
	d.client = outletClient
}

func (d *amqp09OutletDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var config amqp09OutletDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	outlet, err := d.client.GetAmqp09Outlet(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read AMQP 0.9.1 Outlet", err.Error())
		return
	}
	credentials, err := d.client.ListAmqp09Credentials(ctx, config.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Read AMQP 0.9.1 Credentials", err.Error())
		return
	}
	var credentialsList []client.OutletMaterialSummary
	if err == nil {
		credentialsList = credentials
	}
	tlsFeature, err := d.client.GetAmqp09TLSFeature(ctx, config.ID.ValueInt64())
	var tlsFeaturePtr *client.OutletFeatureDetail
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Read AMQP 0.9.1 TLS Feature", err.Error())
			return
		}
	} else if tlsFeature.Feature != "" {
		tlsFeaturePtr = &tlsFeature
	}
	mtlsFeature, err := d.client.GetAmqp09MTLSFeature(ctx, config.ID.ValueInt64())
	var mtlsFeaturePtr *client.OutletFeatureDetail
	if err != nil {
		if !client.IsNotFound(err) {
			resp.Diagnostics.AddError("Unable to Read AMQP 0.9.1 mTLS Feature", err.Error())
			return
		}
	} else if mtlsFeature.Feature != "" {
		mtlsFeaturePtr = &mtlsFeature
	}
	tlsObject, diags := buildAMQPTLSDataSourceObject(tlsFeaturePtr, credentialsList)
	resp.Diagnostics.Append(diags...)
	mtlsObject, diags := buildAMQPMTLSDataSourceObject(mtlsFeaturePtr, credentialsList)
	resp.Diagnostics.Append(diags...)
	transform, err := d.client.GetTypedOutletTransform(ctx, "amqp_0_9", config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read AMQP 0.9.1 Transform", err.Error())
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &amqp09OutletDataSourceModel{
		ID: types.Int64Value(outlet.ID), Name: types.StringValue(outlet.Name), Description: nullableString(outlet.Description), Enabled: types.BoolValue(outlet.Enabled),
		AppID: nullableStringPointer(outlet.AppID), ContentEncoding: nullableStringPointer(outlet.ContentEncoding), ContentType: nullableStringPointer(outlet.ContentType), CorrelationID: nullableStringPointer(outlet.CorrelationID),
		Exchange: types.StringValue(outlet.Exchange), ExchangeDeclareDurable: nullableBoolPointer(outlet.ExchangeDeclareDurable), ExchangeDeclareEnabled: nullableBoolPointer(outlet.ExchangeDeclareEnabled), ExchangeDeclareType: nullableStringPointer(outlet.ExchangeDeclareType),
		Expiration: nullableStringPointer(outlet.Expiration), Immediate: nullableBoolPointer(outlet.Immediate), Key: types.StringValue(outlet.Key), Mandatory: nullableBoolPointer(outlet.Mandatory),
		MaxInFlight: nullableInt64Pointer(outlet.MaxInFlight), MessageID: nullableStringPointer(outlet.MessageID), MetadataExcludePrefixes: nullableStringList(outlet.MetadataExcludePrefixes), NotifyPolicy: nullableStringPointerPreservingEmpty(outlet.NotifyPolicy), Persistent: nullableBoolPointer(outlet.Persistent),
		Priority: nullableStringPointer(outlet.Priority), ReplyTo: nullableStringPointer(outlet.ReplyTo), Timeout: nullableStringPointer(outlet.Timeout), TLSEnableRenegotiation: nullableBoolPointer(outlet.TLSEnableRenegotiation), TLSSkipCertVerify: nullableBoolPointer(outlet.TLSSkipCertVerify), AMQPType: nullableStringPointer(outlet.Type), UserID: nullableStringPointer(outlet.UserID),
		TLS: tlsObject, MTLS: mtlsObject, ApplyStatus: nullableString(outlet.ApplyStatus), LastApplyError: nullableString(outlet.LastApplyError), ConfigGeneration: types.Int64Value(outlet.ConfigGeneration), AppliedGeneration: types.Int64Value(outlet.AppliedGeneration),
		Transform: transformObject(transform),
	})...)
}

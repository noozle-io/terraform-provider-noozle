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
	_ resource.Resource                   = &amqp09OutletResource{}
	_ resource.ResourceWithConfigure      = &amqp09OutletResource{}
	_ resource.ResourceWithImportState    = &amqp09OutletResource{}
	_ resource.ResourceWithValidateConfig = &amqp09OutletResource{}
)

type amqp09OutletClient interface {
	CreateAmqp09Outlet(context.Context, client.Amqp09OutletCreateRequest) (client.Amqp09Outlet, error)
	GetAmqp09Outlet(context.Context, int64) (client.Amqp09Outlet, error)
	UpdateAmqp09Outlet(context.Context, int64, client.Amqp09OutletUpdateRequest) error
	DeleteAmqp09Outlet(context.Context, int64) error
	DeleteAmqp09OutletOperation(context.Context, int64) (client.OutletOperation, error)
	GetTypedOutletStatus(context.Context, string, int64, int64) (client.OutletStatus, error)
	SetAmqp09Field(context.Context, int64, string, any) error
	DeleteAmqp09Field(context.Context, int64, string) error
	ListAmqp09Credentials(context.Context, int64) ([]client.OutletMaterialSummary, error)
	DeleteAmqp09Credential(context.Context, int64, string) (client.OutletMaterialDeleteResponse, error)
	GetAmqp09TLSFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureAmqp09TLSFeature(context.Context, int64, client.AmqpTLSFeatureRequest) error
	UpdateAmqp09TLSFeature(context.Context, int64, client.AmqpTLSFeatureRequest) error
	DeleteAmqp09TLSFeature(context.Context, int64) error
	GetAmqp09MTLSFeature(context.Context, int64) (client.OutletFeatureDetail, error)
	ConfigureAmqp09MTLSFeature(context.Context, int64, client.AmqpMTLSFeatureRequest) error
	UpdateAmqp09MTLSFeature(context.Context, int64, client.AmqpMTLSFeatureRequest) error
	DeleteAmqp09MTLSFeature(context.Context, int64) error
	typedOutletTransformClient
}

type amqp09OutletFieldUpdater struct{ client amqp09OutletClient }

func (u amqp09OutletFieldUpdater) SetField(ctx context.Context, outletID int64, field string, value any) error {
	return u.client.SetAmqp09Field(ctx, outletID, field, value)
}

func (u amqp09OutletFieldUpdater) DeleteField(ctx context.Context, outletID int64, field string) error {
	return u.client.DeleteAmqp09Field(ctx, outletID, field)
}

type amqp09OutletResource struct{ client amqp09OutletClient }

type amqp09OutletResourceModel struct {
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
	Credentials             types.Object `tfsdk:"credentials"`
	TLS                     types.Object `tfsdk:"tls"`
	MTLS                    types.Object `tfsdk:"mtls"`
	ApplyStatus             types.String `tfsdk:"apply_status"`
	LastApplyError          types.String `tfsdk:"last_apply_error"`
	ConfigGeneration        types.Int64  `tfsdk:"config_generation"`
	AppliedGeneration       types.Int64  `tfsdk:"applied_generation"`
	Transform               types.Object `tfsdk:"transform"`
}

func NewAmqp09OutletResource() resource.Resource { return &amqp09OutletResource{} }

func (r *amqp09OutletResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_amqp_0_9_outlet"
}

func (r *amqp09OutletResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		MarkdownDescription: "Manages a Noozle AMQP 0.9.1 outlet.",
		Attributes: map[string]resourceschema.Attribute{
			"id":                        resourceschema.Int64Attribute{MarkdownDescription: "Unique identifier of the outlet.", Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}},
			"name":                      resourceschema.StringAttribute{MarkdownDescription: "Name presented for the outlet.", Required: true},
			"description":               resourceschema.StringAttribute{MarkdownDescription: "Optional note explaining the outlet purpose.", Optional: true},
			"enabled":                   resourceschema.BoolAttribute{MarkdownDescription: "Whether the outlet is active.", Required: true},
			"app_id":                    resourceschema.StringAttribute{MarkdownDescription: "AMQP application ID interpolation.", Optional: true},
			"content_encoding":          resourceschema.StringAttribute{MarkdownDescription: "AMQP content encoding interpolation.", Optional: true},
			"content_type":              resourceschema.StringAttribute{MarkdownDescription: "AMQP content type interpolation.", Optional: true},
			"correlation_id":            resourceschema.StringAttribute{MarkdownDescription: "AMQP correlation ID interpolation.", Optional: true},
			"exchange":                  resourceschema.StringAttribute{MarkdownDescription: "AMQP exchange name.", Required: true},
			"exchange_declare_durable":  resourceschema.BoolAttribute{MarkdownDescription: "Declare the AMQP exchange as durable.", Optional: true},
			"exchange_declare_enabled":  resourceschema.BoolAttribute{MarkdownDescription: "Passively declare and verify the AMQP exchange.", Optional: true},
			"exchange_declare_type":     resourceschema.StringAttribute{MarkdownDescription: "AMQP exchange type to declare.", Optional: true},
			"expiration":                resourceschema.StringAttribute{MarkdownDescription: "AMQP per-message TTL interpolation.", Optional: true},
			"immediate":                 resourceschema.BoolAttribute{MarkdownDescription: "Set the AMQP immediate flag.", Optional: true},
			"key":                       resourceschema.StringAttribute{MarkdownDescription: "AMQP routing key.", Required: true},
			"mandatory":                 resourceschema.BoolAttribute{MarkdownDescription: "Set the AMQP mandatory flag.", Optional: true},
			"max_in_flight":             resourceschema.Int64Attribute{MarkdownDescription: "Maximum number of in-flight AMQP messages.", Optional: true},
			"message_id":                resourceschema.StringAttribute{MarkdownDescription: "AMQP message ID interpolation.", Optional: true},
			"metadata_exclude_prefixes": resourceschema.ListAttribute{MarkdownDescription: "Metadata prefixes to exclude from AMQP headers.", Optional: true, ElementType: types.StringType},
			"notify_policy":             resourceschema.StringAttribute{MarkdownDescription: "Delivery policy applied when notifying this outlet.", Optional: true},
			"persistent":                resourceschema.BoolAttribute{MarkdownDescription: "Publish AMQP messages as persistent.", Optional: true},
			"priority":                  resourceschema.StringAttribute{MarkdownDescription: "AMQP priority interpolation.", Optional: true},
			"reply_to":                  resourceschema.StringAttribute{MarkdownDescription: "AMQP reply-to interpolation.", Optional: true},
			"timeout":                   resourceschema.StringAttribute{MarkdownDescription: "AMQP publish timeout.", Optional: true},
			"tls_enable_renegotiation":  resourceschema.BoolAttribute{MarkdownDescription: "Allow TLS renegotiation for AMQP connections.", Optional: true},
			"tls_skip_cert_verify":      resourceschema.BoolAttribute{MarkdownDescription: "Skip AMQP TLS certificate verification.", Optional: true},
			"type":                      resourceschema.StringAttribute{MarkdownDescription: "AMQP message type interpolation.", Optional: true},
			"user_id":                   resourceschema.StringAttribute{MarkdownDescription: "AMQP user ID interpolation.", Optional: true},
			"apply_status":              resourceschema.StringAttribute{MarkdownDescription: "Observed apply/runtime status for the outlet.", Computed: true},
			"last_apply_error":          resourceschema.StringAttribute{MarkdownDescription: "Last observed apply/runtime error message.", Computed: true},
			"config_generation":         resourceschema.Int64Attribute{MarkdownDescription: "Desired outlet config generation.", Computed: true},
			"applied_generation":        resourceschema.Int64Attribute{MarkdownDescription: "Last observed applied config generation.", Computed: true},
			"transform": resourceschema.SingleNestedAttribute{MarkdownDescription: "Optional outlet transform selection and observed runtime state.", Optional: true, Computed: true, Attributes: map[string]resourceschema.Attribute{
				"selected": resourceschema.StringAttribute{MarkdownDescription: "Explicitly selected transform preset.", Optional: true}, "effective": resourceschema.StringAttribute{MarkdownDescription: "Effective transform after defaults are applied.", Computed: true}, "state": resourceschema.StringAttribute{MarkdownDescription: "Observed transform runtime state.", Computed: true}, "applied": resourceschema.BoolAttribute{MarkdownDescription: "Whether the transform has been applied.", Computed: true}, "generation": resourceschema.Int64Attribute{MarkdownDescription: "Observed transform generation.", Computed: true}, "last_error": resourceschema.StringAttribute{MarkdownDescription: "Last observed transform apply error.", Computed: true},
			}},
		},
		Blocks: map[string]resourceschema.Block{
			"credentials": resourceschema.SingleNestedBlock{
				MarkdownDescription: "Connection URLs for this outlet. The URLs are write-only; rotate them by updating the list and incrementing `urls_version`.",
				Attributes: map[string]resourceschema.Attribute{
					"urls":         resourceschema.ListAttribute{MarkdownDescription: "AMQP connection URLs.", Optional: true, Sensitive: true, WriteOnly: true, ElementType: types.StringType},
					"urls_version": resourceschema.Int64Attribute{MarkdownDescription: "Monotonic version used to rotate the connection URLs.", Optional: true},
				},
			},
			"tls": resourceschema.SingleNestedBlock{
				MarkdownDescription: "Optional TLS feature for the outlet. Conflicts with `mtls`.",
				Attributes: map[string]resourceschema.Attribute{
					"enabled":                 resourceschema.BoolAttribute{Computed: true},
					"state":                   resourceschema.StringAttribute{Computed: true},
					"desired_generation":      resourceschema.Int64Attribute{Computed: true},
					"applied_generation":      resourceschema.Int64Attribute{Computed: true},
					"last_error":              resourceschema.StringAttribute{Computed: true},
					"ca_cert_filename":        resourceschema.StringAttribute{Optional: true},
					"ca_cert_content_base64":  resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true},
					"ca_cert_version":         resourceschema.Int64Attribute{Optional: true},
					"ca_cert_checksum_sha256": resourceschema.StringAttribute{Computed: true},
				},
			},
			"mtls": resourceschema.SingleNestedBlock{
				MarkdownDescription: "Optional mutual TLS feature for the outlet. Conflicts with `tls`.",
				Attributes: map[string]resourceschema.Attribute{
					"enabled":                     resourceschema.BoolAttribute{Computed: true},
					"state":                       resourceschema.StringAttribute{Computed: true},
					"desired_generation":          resourceschema.Int64Attribute{Computed: true},
					"applied_generation":          resourceschema.Int64Attribute{Computed: true},
					"last_error":                  resourceschema.StringAttribute{Computed: true},
					"ca_cert_filename":            resourceschema.StringAttribute{Optional: true},
					"ca_cert_content_base64":      resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true},
					"ca_cert_version":             resourceschema.Int64Attribute{Optional: true},
					"ca_cert_checksum_sha256":     resourceschema.StringAttribute{Computed: true},
					"client_cert_filename":        resourceschema.StringAttribute{Optional: true},
					"client_cert_content_base64":  resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true},
					"client_cert_version":         resourceschema.Int64Attribute{Optional: true},
					"client_cert_checksum_sha256": resourceschema.StringAttribute{Computed: true},
					"client_key_filename":         resourceschema.StringAttribute{Optional: true},
					"client_key_content_base64":   resourceschema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true},
					"client_key_version":          resourceschema.Int64Attribute{Optional: true},
					"client_key_checksum_sha256":  resourceschema.StringAttribute{Computed: true},
				},
			},
		},
	}
}

func (r *amqp09OutletResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	outletClient, ok := req.ProviderData.(amqp09OutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected AMQP 0.9.1 outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}
	r.client = outletClient
}

func (r *amqp09OutletResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config amqp09OutletResourceModel
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
		resp.Diagnostics.AddError("Missing AMQP 0.9.1 Credentials", "The `credentials` block must set both `urls` and `urls_version`.")
	}

	if !config.TLS.IsNull() && !config.TLS.IsUnknown() && !config.MTLS.IsNull() && !config.MTLS.IsUnknown() {
		resp.Diagnostics.AddError("Conflicting AMQP 0.9.1 Features", "Only one of `tls` or `mtls` can be configured at the same time.")
	}

	if !config.TLS.IsNull() && !config.TLS.IsUnknown() {
		tls, diags := expandAMQPTLS(ctx, config.TLS)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if tls.CACertFilename.IsNull() || tls.CACertFilename.IsUnknown() || tls.CACertContentBase64.IsNull() || tls.CACertContentBase64.IsUnknown() || tls.CACertVersion.IsNull() || tls.CACertVersion.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete AMQP 0.9.1 TLS Configuration", "When `tls` is configured, `ca_cert_filename`, `ca_cert_content_base64`, and `ca_cert_version` must all be set.")
		}
	}

	if !config.MTLS.IsNull() && !config.MTLS.IsUnknown() {
		mtls, diags := expandAMQPMTLS(ctx, config.MTLS)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if mtls.CACertFilename.IsNull() || mtls.CACertFilename.IsUnknown() ||
			mtls.CACertContentBase64.IsNull() || mtls.CACertContentBase64.IsUnknown() ||
			mtls.CACertVersion.IsNull() || mtls.CACertVersion.IsUnknown() ||
			mtls.ClientCertFilename.IsNull() || mtls.ClientCertFilename.IsUnknown() ||
			mtls.ClientCertContentBase64.IsNull() || mtls.ClientCertContentBase64.IsUnknown() ||
			mtls.ClientCertVersion.IsNull() || mtls.ClientCertVersion.IsUnknown() ||
			mtls.ClientKeyFilename.IsNull() || mtls.ClientKeyFilename.IsUnknown() ||
			mtls.ClientKeyContentBase64.IsNull() || mtls.ClientKeyContentBase64.IsUnknown() ||
			mtls.ClientKeyVersion.IsNull() || mtls.ClientKeyVersion.IsUnknown() {
			resp.Diagnostics.AddError("Incomplete AMQP 0.9.1 mTLS Configuration", "When `mtls` is configured, CA cert, client cert, and client key filename/content/version values must all be set.")
		}
	}
}

func (r *amqp09OutletResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var plan amqp09OutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	credentials, tls, mtls, ok := r.readConfig(ctx, req.Config, plan, &resp.Diagnostics)
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
		resp.Diagnostics.AddError("Invalid AMQP 0.9.1 URLs", err.Error())
		return
	}
	metadataExcludePrefixes, err := stringListValueOrNil(ctx, plan.MetadataExcludePrefixes)
	if err != nil {
		resp.Diagnostics.AddError("Invalid AMQP 0.9.1 Metadata Exclude Prefixes", err.Error())
		return
	}
	outlet, err := r.client.CreateAmqp09Outlet(ctx, client.Amqp09OutletCreateRequest{
		Name: plan.Name.ValueString(), Description: stringValueOrEmpty(plan.Description), Enabled: plan.Enabled.ValueBool(),
		AppID: stringPointerValue(plan.AppID), ContentEncoding: stringPointerValue(plan.ContentEncoding), ContentType: stringPointerValue(plan.ContentType), CorrelationID: stringPointerValue(plan.CorrelationID),
		Exchange: plan.Exchange.ValueString(), ExchangeDeclareDurable: boolPointerValue(plan.ExchangeDeclareDurable), ExchangeDeclareEnabled: boolPointerValue(plan.ExchangeDeclareEnabled), ExchangeDeclareType: stringPointerValue(plan.ExchangeDeclareType),
		Expiration: stringPointerValue(plan.Expiration), Immediate: boolPointerValue(plan.Immediate), Key: plan.Key.ValueString(), Mandatory: boolPointerValue(plan.Mandatory),
		MaxInFlight: int64PointerValue(plan.MaxInFlight), MessageID: stringPointerValue(plan.MessageID), MetadataExcludePrefixes: metadataExcludePrefixes, NotifyPolicy: stringPointerValue(plan.NotifyPolicy), Persistent: boolPointerValue(plan.Persistent),
		Priority: stringPointerValue(plan.Priority), ReplyTo: stringPointerValue(plan.ReplyTo), Timeout: stringPointerValue(plan.Timeout),
		TLSEnableRenegotiation: boolPointerValue(plan.TLSEnableRenegotiation), TLSSkipCertVerify: boolPointerValue(plan.TLSSkipCertVerify), Type: stringPointerValue(plan.AMQPType), URLs: urls, UserID: stringPointerValue(plan.UserID),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create AMQP 0.9.1 Outlet", err.Error())
		return
	}
	if err := r.applyFeatureConfig(ctx, outlet.ID, types.ObjectNull(amqpTLSAttributeTypes), types.ObjectNull(amqpMTLSAttributeTypes), tls, mtls); err != nil {
		if deleteErr := r.client.DeleteAmqp09Outlet(ctx, outlet.ID); deleteErr != nil && !client.IsNotFound(deleteErr) {
			resp.Diagnostics.AddError("Unable to Configure AMQP 0.9.1 Features", fmt.Sprintf("%s; cleanup failed: %s", err, deleteErr))
			return
		}
		resp.Diagnostics.AddError("Unable to Configure AMQP 0.9.1 Features", err.Error())
		return
	}
	if !configuredTransform.IsNull() && !configuredTransform.IsUnknown() {
		if _, err := r.client.SetTypedOutletTransform(ctx, "amqp_0_9", outlet.ID, configuredTransform.ValueString()); err != nil {
			if deleteErr := r.client.DeleteAmqp09Outlet(ctx, outlet.ID); deleteErr != nil && !client.IsNotFound(deleteErr) {
				resp.Diagnostics.AddError("Unable to Configure AMQP 0.9.1 Transform", fmt.Sprintf("%s; cleanup failed: %s", err, deleteErr))
				return
			}
			resp.Diagnostics.AddError("Unable to Configure AMQP 0.9.1 Transform", err.Error())
			return
		}
	}
	newState, removed, diags := r.buildState(ctx, outlet.ID, plan.Credentials, plan.TLS, plan.MTLS)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if removed {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "amqp_0_9", outlet.ID, plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe AMQP 0.9.1 Outlet Lifecycle", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *amqp09OutletResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var state amqp09OutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	newState, removed, diags := r.buildState(ctx, state.ID.ValueInt64(), state.Credentials, state.TLS, state.MTLS)
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

func (r *amqp09OutletResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, recorder := client.WithOutletOperationRecorder(ctx)
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var plan, state amqp09OutletResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	credentials, tls, mtls, ok := r.readConfig(ctx, req.Config, plan, &resp.Diagnostics)
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
		resp.Diagnostics.AddError("Invalid AMQP 0.9.1 URLs", err.Error())
		return
	}
	metadataExcludePrefixes, err := stringListValueOrNil(ctx, plan.MetadataExcludePrefixes)
	if err != nil {
		resp.Diagnostics.AddError("Invalid AMQP 0.9.1 Metadata Exclude Prefixes", err.Error())
		return
	}
	if err := r.client.UpdateAmqp09Outlet(ctx, state.ID.ValueInt64(), client.Amqp09OutletUpdateRequest{
		Name: plan.Name.ValueString(), Description: stringValueOrEmpty(plan.Description), Enabled: plan.Enabled.ValueBool(),
		AppID: stringPointerValue(plan.AppID), ContentEncoding: stringPointerValue(plan.ContentEncoding), ContentType: stringPointerValue(plan.ContentType), CorrelationID: stringPointerValue(plan.CorrelationID),
		Exchange: plan.Exchange.ValueString(), ExchangeDeclareDurable: boolPointerValue(plan.ExchangeDeclareDurable), ExchangeDeclareEnabled: boolPointerValue(plan.ExchangeDeclareEnabled), ExchangeDeclareType: stringPointerValue(plan.ExchangeDeclareType),
		Expiration: stringPointerValue(plan.Expiration), Immediate: boolPointerValue(plan.Immediate), Key: plan.Key.ValueString(), Mandatory: boolPointerValue(plan.Mandatory),
		MaxInFlight: int64PointerValue(plan.MaxInFlight), MessageID: stringPointerValue(plan.MessageID), MetadataExcludePrefixes: metadataExcludePrefixes, NotifyPolicy: stringPointerValue(plan.NotifyPolicy), Persistent: boolPointerValue(plan.Persistent),
		Priority: stringPointerValue(plan.Priority), ReplyTo: stringPointerValue(plan.ReplyTo), Timeout: stringPointerValue(plan.Timeout),
		TLSEnableRenegotiation: boolPointerValue(plan.TLSEnableRenegotiation), TLSSkipCertVerify: boolPointerValue(plan.TLSSkipCertVerify), Type: stringPointerValue(plan.AMQPType), URLs: urls, UserID: stringPointerValue(plan.UserID),
	}); err != nil {
		resp.Diagnostics.AddError("Unable to Update AMQP 0.9.1 Outlet", err.Error())
		return
	}
	fieldUpdater := amqp09OutletFieldUpdater{client: r.client}
	mutations := []struct {
		name string
		err  error
	}{
		{"app_id", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.app_id", plan.AppID, state.AppID)},
		{"content_encoding", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.content_encoding", plan.ContentEncoding, state.ContentEncoding)},
		{"content_type", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.content_type", plan.ContentType, state.ContentType)},
		{"correlation_id", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.correlation_id", plan.CorrelationID, state.CorrelationID)},
		{"exchange_declare_durable", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.exchange_declare.durable", plan.ExchangeDeclareDurable, state.ExchangeDeclareDurable)},
		{"exchange_declare_enabled", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.exchange_declare.enabled", plan.ExchangeDeclareEnabled, state.ExchangeDeclareEnabled)},
		{"exchange_declare_type", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.exchange_declare.type", plan.ExchangeDeclareType, state.ExchangeDeclareType)},
		{"expiration", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.expiration", plan.Expiration, state.Expiration)},
		{"immediate", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.immediate", plan.Immediate, state.Immediate)},
		{"mandatory", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.mandatory", plan.Mandatory, state.Mandatory)},
		{"max_in_flight", syncInt64Field(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.max_in_flight", plan.MaxInFlight, state.MaxInFlight)},
		{"message_id", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.message_id", plan.MessageID, state.MessageID)},
		{"metadata_exclude_prefixes", syncStringListField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.metadata.exclude_prefixes", plan.MetadataExcludePrefixes, state.MetadataExcludePrefixes)},
		{"notify_policy", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "notify_policy", plan.NotifyPolicy, state.NotifyPolicy)},
		{"persistent", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.persistent", plan.Persistent, state.Persistent)},
		{"priority", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.priority", plan.Priority, state.Priority)},
		{"reply_to", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.reply_to", plan.ReplyTo, state.ReplyTo)},
		{"timeout", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.timeout", plan.Timeout, state.Timeout)},
		{"tls_enable_renegotiation", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.tls.enable_renegotiation", plan.TLSEnableRenegotiation, state.TLSEnableRenegotiation)},
		{"tls_skip_cert_verify", syncBoolField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.tls.skip_cert_verify", plan.TLSSkipCertVerify, state.TLSSkipCertVerify)},
		{"type", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.type", plan.AMQPType, state.AMQPType)},
		{"user_id", syncStringField(ctx, fieldUpdater, state.ID.ValueInt64(), "amqp_0_9.user_id", plan.UserID, state.UserID)},
	}
	for _, mutation := range mutations {
		if mutation.err != nil {
			resp.Diagnostics.AddError("Unable to Update AMQP 0.9.1 Field", fmt.Sprintf("%s: %s", mutation.name, mutation.err))
			return
		}
	}
	if err := syncTypedOutletTransform(ctx, r.client, "amqp_0_9", state.ID.ValueInt64(), configuredTransform, selectedTransform(state.Transform)); err != nil {
		resp.Diagnostics.AddError("Unable to Update AMQP 0.9.1 Transform", err.Error())
		return
	}
	if err := r.applyFeatureConfig(ctx, state.ID.ValueInt64(), state.TLS, state.MTLS, tls, mtls); err != nil {
		resp.Diagnostics.AddError("Unable to Update AMQP 0.9.1 Features", err.Error())
		return
	}
	newState, removed, diags := r.buildState(ctx, state.ID.ValueInt64(), plan.Credentials, plan.TLS, plan.MTLS)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if removed {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := observeFinalTypedOutletOperation(ctx, r.client, recorder, "amqp_0_9", state.ID.ValueInt64(), plan.Enabled.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Unable to Observe AMQP 0.9.1 Outlet Lifecycle", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *amqp09OutletResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var state amqp09OutletResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	operation, err := r.client.DeleteAmqp09OutletOperation(ctx, state.ID.ValueInt64())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to Delete AMQP 0.9.1 Outlet", err.Error())
		return
	}
	if err == nil && operation.DesiredGeneration != 0 {
		if err := observeTypedOutletLifecycle(ctx, r.client, "amqp_0_9", state.ID.ValueInt64(), operation.DesiredGeneration, false, true); err != nil {
			resp.Diagnostics.AddError("Unable to Observe AMQP 0.9.1 Outlet Lifecycle", err.Error())
		}
	}
}

func (r *amqp09OutletResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid AMQP 0.9.1 Outlet Import Identifier", fmt.Sprintf("Expected numeric outlet ID, got %q: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (r *amqp09OutletResource) readConfig(ctx context.Context, cfg tfsdk.Config, plan amqp09OutletResourceModel, diags *diag.Diagnostics) (amqpURLsCredentialsModel, amqpTLSModel, amqpMTLSModel, bool) {
	credentials, credentialsDiags := expandAmqpURLsCredentials(ctx, plan.Credentials)
	diags.Append(credentialsDiags...)
	tls, tlsDiags := expandAMQPTLS(ctx, plan.TLS)
	diags.Append(tlsDiags...)
	mtls, mtlsDiags := expandAMQPMTLS(ctx, plan.MTLS)
	diags.Append(mtlsDiags...)
	if diags.HasError() {
		return credentials, tls, mtls, false
	}
	credentials, credentialsDiags = enrichAmqpURLsCredentialsFromConfig(ctx, cfg, credentials)
	diags.Append(credentialsDiags...)
	tls, tlsDiags = enrichAMQPTLSFromConfig(ctx, cfg, "tls", tls)
	diags.Append(tlsDiags...)
	mtls, mtlsDiags = enrichAMQPMTLSFromConfig(ctx, cfg, "mtls", mtls)
	diags.Append(mtlsDiags...)
	if diags.HasError() {
		return credentials, tls, mtls, false
	}
	return credentials, tls, mtls, true
}

func (r *amqp09OutletResource) applyFeatureConfig(ctx context.Context, outletID int64, stateTLS, stateMTLS types.Object, planTLS amqpTLSModel, planMTLS amqpMTLSModel) error {
	if !stateTLS.IsNull() && !amqpFeatureEnabled(stateTLS) {
		stateTLS = types.ObjectNull(amqpTLSAttributeTypes)
	}
	if !stateMTLS.IsNull() && !amqpFeatureEnabled(stateMTLS) {
		stateMTLS = types.ObjectNull(amqpMTLSAttributeTypes)
	}
	if !stateTLS.IsNull() && !stateTLS.IsUnknown() && !planMTLS.ClientKeyVersion.IsNull() {
		if err := r.client.DeleteAmqp09TLSFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
			return err
		}
		if _, err := r.client.DeleteAmqp09Credential(ctx, outletID, amqpCACertRole); err != nil && !client.IsNotFound(err) {
			return err
		}
		stateTLS = types.ObjectNull(amqpTLSAttributeTypes)
	}
	if !stateMTLS.IsNull() && !stateMTLS.IsUnknown() && !planTLS.CACertVersion.IsNull() {
		if err := r.client.DeleteAmqp09MTLSFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
			return err
		}
		for _, role := range []string{amqpCACertRole, amqpClientCertRole, amqpClientKeyRole} {
			if _, err := r.client.DeleteAmqp09Credential(ctx, outletID, role); err != nil && !client.IsNotFound(err) {
				return err
			}
		}
		stateMTLS = types.ObjectNull(amqpMTLSAttributeTypes)
	}
	if planTLS.CACertVersion.IsNull() {
		if !stateTLS.IsNull() && !stateTLS.IsUnknown() {
			if err := r.client.DeleteAmqp09TLSFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
				return err
			}
			if _, err := r.client.DeleteAmqp09Credential(ctx, outletID, amqpCACertRole); err != nil && !client.IsNotFound(err) {
				return err
			}
		}
	} else if stateTLS.IsNull() || stateTLS.IsUnknown() {
		if err := r.client.ConfigureAmqp09TLSFeature(ctx, outletID, client.AmqpTLSFeatureRequest{
			CACertFilename: planTLS.CACertFilename.ValueString(), CACertContentBase64: planTLS.CACertContentBase64.ValueString(),
		}); err != nil {
			return err
		}
	} else {
		stateModel, diags := expandAMQPTLS(ctx, stateTLS)
		if diags.HasError() {
			return errors.New(diags.Errors()[0].Detail())
		}
		if !planTLS.CACertVersion.Equal(stateModel.CACertVersion) || !planTLS.CACertFilename.Equal(stateModel.CACertFilename) {
			if err := r.client.UpdateAmqp09TLSFeature(ctx, outletID, client.AmqpTLSFeatureRequest{
				CACertFilename: planTLS.CACertFilename.ValueString(), CACertContentBase64: planTLS.CACertContentBase64.ValueString(),
			}); err != nil {
				return err
			}
		}
	}
	if planMTLS.ClientKeyVersion.IsNull() {
		if !stateMTLS.IsNull() && !stateMTLS.IsUnknown() {
			if err := r.client.DeleteAmqp09MTLSFeature(ctx, outletID); err != nil && !client.IsNotFound(err) {
				return err
			}
			for _, role := range []string{amqpCACertRole, amqpClientCertRole, amqpClientKeyRole} {
				if _, err := r.client.DeleteAmqp09Credential(ctx, outletID, role); err != nil && !client.IsNotFound(err) {
					return err
				}
			}
		}
	} else if stateMTLS.IsNull() || stateMTLS.IsUnknown() {
		if err := r.client.ConfigureAmqp09MTLSFeature(ctx, outletID, client.AmqpMTLSFeatureRequest{
			CACertFilename: planMTLS.CACertFilename.ValueString(), CACertContentBase64: planMTLS.CACertContentBase64.ValueString(),
			ClientCertFilename: planMTLS.ClientCertFilename.ValueString(), ClientCertContentBase64: planMTLS.ClientCertContentBase64.ValueString(),
			ClientKeyFilename: planMTLS.ClientKeyFilename.ValueString(), ClientKeyContentBase64: planMTLS.ClientKeyContentBase64.ValueString(),
		}); err != nil {
			return err
		}
	} else {
		stateModel, diags := expandAMQPMTLS(ctx, stateMTLS)
		if diags.HasError() {
			return errors.New(diags.Errors()[0].Detail())
		}
		if !planMTLS.CACertVersion.Equal(stateModel.CACertVersion) || !planMTLS.CACertFilename.Equal(stateModel.CACertFilename) ||
			!planMTLS.ClientCertVersion.Equal(stateModel.ClientCertVersion) || !planMTLS.ClientCertFilename.Equal(stateModel.ClientCertFilename) ||
			!planMTLS.ClientKeyVersion.Equal(stateModel.ClientKeyVersion) || !planMTLS.ClientKeyFilename.Equal(stateModel.ClientKeyFilename) {
			if err := r.client.UpdateAmqp09MTLSFeature(ctx, outletID, client.AmqpMTLSFeatureRequest{
				CACertFilename: planMTLS.CACertFilename.ValueString(), CACertContentBase64: planMTLS.CACertContentBase64.ValueString(),
				ClientCertFilename: planMTLS.ClientCertFilename.ValueString(), ClientCertContentBase64: planMTLS.ClientCertContentBase64.ValueString(),
				ClientKeyFilename: planMTLS.ClientKeyFilename.ValueString(), ClientKeyContentBase64: planMTLS.ClientKeyContentBase64.ValueString(),
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *amqp09OutletResource) buildState(ctx context.Context, outletID int64, previousCredentials, previousTLS, previousMTLS types.Object) (amqp09OutletResourceModel, bool, diag.Diagnostics) {
	outlet, err := r.client.GetAmqp09Outlet(ctx, outletID)
	if err != nil {
		if client.IsNotFound(err) {
			return amqp09OutletResourceModel{}, true, nil
		}
		var diags diag.Diagnostics
		diags.AddError("Unable to Read AMQP 0.9.1 Outlet", err.Error())
		return amqp09OutletResourceModel{}, false, diags
	}
	credentials, err := r.client.ListAmqp09Credentials(ctx, outletID)
	if err != nil {
		if !client.IsNotFound(err) {
			var diags diag.Diagnostics
			diags.AddError("Unable to Read AMQP 0.9.1 Credentials", err.Error())
			return amqp09OutletResourceModel{}, false, diags
		}
		credentials = nil
	}
	tlsFeature, err := r.client.GetAmqp09TLSFeature(ctx, outletID)
	var tlsFeaturePtr *client.OutletFeatureDetail
	if err != nil && !client.IsNotFound(err) {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read AMQP 0.9.1 TLS Feature", err.Error())
		return amqp09OutletResourceModel{}, false, diags
	}
	if err == nil && tlsFeature.Feature != "" {
		tlsFeaturePtr = &tlsFeature
	}
	mtlsFeature, err := r.client.GetAmqp09MTLSFeature(ctx, outletID)
	var mtlsFeaturePtr *client.OutletFeatureDetail
	if err != nil && !client.IsNotFound(err) {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read AMQP 0.9.1 mTLS Feature", err.Error())
		return amqp09OutletResourceModel{}, false, diags
	}
	if err == nil && mtlsFeature.Feature != "" {
		mtlsFeaturePtr = &mtlsFeature
	}
	credentialsObject, diags := buildAmqpURLsCredentialsObject(ctx, previousCredentials, credentials, amqp09URLsCredentialRole)
	if diags.HasError() {
		return amqp09OutletResourceModel{}, false, diags
	}
	tlsObject, diags := buildAMQPTLSObject(ctx, previousTLS, tlsFeaturePtr, credentials)
	if diags.HasError() {
		return amqp09OutletResourceModel{}, false, diags
	}
	mtlsObject, diags := buildAMQPMTLSObject(ctx, previousMTLS, mtlsFeaturePtr, credentials)
	if diags.HasError() {
		return amqp09OutletResourceModel{}, false, diags
	}
	transform, err := r.client.GetTypedOutletTransform(ctx, "amqp_0_9", outletID)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Unable to Read AMQP 0.9.1 Transform", err.Error())
		return amqp09OutletResourceModel{}, false, diags
	}
	return amqp09OutletModelFromAPI(outlet, credentialsObject, tlsObject, mtlsObject, transform), false, nil
}

func amqp09OutletModelFromAPI(outlet client.Amqp09Outlet, credentials, tls, mtls types.Object, transform client.OutletTransformState) amqp09OutletResourceModel {
	return amqp09OutletResourceModel{
		ID: types.Int64Value(outlet.ID), Name: types.StringValue(outlet.Name), Description: nullableString(outlet.Description), Enabled: types.BoolValue(outlet.Enabled),
		AppID: nullableStringPointer(outlet.AppID), ContentEncoding: nullableStringPointer(outlet.ContentEncoding), ContentType: nullableStringPointer(outlet.ContentType), CorrelationID: nullableStringPointer(outlet.CorrelationID),
		Exchange: types.StringValue(outlet.Exchange), ExchangeDeclareDurable: nullableBoolPointer(outlet.ExchangeDeclareDurable), ExchangeDeclareEnabled: nullableBoolPointer(outlet.ExchangeDeclareEnabled), ExchangeDeclareType: nullableStringPointer(outlet.ExchangeDeclareType),
		Expiration: nullableStringPointer(outlet.Expiration), Immediate: nullableBoolPointer(outlet.Immediate), Key: types.StringValue(outlet.Key), Mandatory: nullableBoolPointer(outlet.Mandatory),
		MaxInFlight: nullableInt64Pointer(outlet.MaxInFlight), MessageID: nullableStringPointer(outlet.MessageID), MetadataExcludePrefixes: nullableStringList(outlet.MetadataExcludePrefixes), NotifyPolicy: nullableStringPointerPreservingEmpty(outlet.NotifyPolicy), Persistent: nullableBoolPointer(outlet.Persistent),
		Priority: nullableStringPointer(outlet.Priority), ReplyTo: nullableStringPointer(outlet.ReplyTo), Timeout: nullableStringPointer(outlet.Timeout),
		TLSEnableRenegotiation: nullableBoolPointer(outlet.TLSEnableRenegotiation), TLSSkipCertVerify: nullableBoolPointer(outlet.TLSSkipCertVerify), AMQPType: nullableStringPointer(outlet.Type), UserID: nullableStringPointer(outlet.UserID),
		Credentials: credentials, TLS: tls, MTLS: mtls, ApplyStatus: nullableString(outlet.ApplyStatus), LastApplyError: nullableString(outlet.LastApplyError),
		ConfigGeneration: types.Int64Value(outlet.ConfigGeneration), AppliedGeneration: types.Int64Value(outlet.AppliedGeneration),
		Transform: transformObject(transform),
	}
}

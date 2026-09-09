package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &gcpPubsubOutletDataSource{}
	_ datasource.DataSourceWithConfigure = &gcpPubsubOutletDataSource{}
)

type gcpPubsubOutletDataSource struct {
	client gcpPubsubOutletClient
}

type gcpPubsubOutletDataSourceModel struct {
	ID                                types.Int64   `tfsdk:"id"`
	Name                              types.String  `tfsdk:"name"`
	Description                       types.String  `tfsdk:"description"`
	Enabled                           types.Bool    `tfsdk:"enabled"`
	BatchingByteSize                  types.Int64   `tfsdk:"batching_byte_size"`
	BatchingCheck                     types.String  `tfsdk:"batching_check"`
	BatchingCount                     types.Int64   `tfsdk:"batching_count"`
	BatchingJitter                    types.Float64 `tfsdk:"batching_jitter"`
	BatchingPeriod                    types.String  `tfsdk:"batching_period"`
	ByteThreshold                     types.Int64   `tfsdk:"byte_threshold"`
	CountThreshold                    types.Int64   `tfsdk:"count_threshold"`
	DelayThreshold                    types.String  `tfsdk:"delay_threshold"`
	Endpoint                          types.String  `tfsdk:"endpoint"`
	FlowControlLimitExceededBehavior  types.String  `tfsdk:"flow_control_limit_exceeded_behavior"`
	FlowControlMaxOutstandingBytes    types.Int64   `tfsdk:"flow_control_max_outstanding_bytes"`
	FlowControlMaxOutstandingMessages types.Int64   `tfsdk:"flow_control_max_outstanding_messages"`
	MaxInFlight                       types.Int64   `tfsdk:"max_in_flight"`
	MetadataExcludePrefixes           types.List    `tfsdk:"metadata_exclude_prefixes"`
	NotifyPolicy                      types.String  `tfsdk:"notify_policy"`
	OrderingKey                       types.String  `tfsdk:"ordering_key"`
	Project                           types.String  `tfsdk:"project"`
	PublishTimeout                    types.String  `tfsdk:"publish_timeout"`
	Topic                             types.String  `tfsdk:"topic"`
	ApplyStatus                       types.String  `tfsdk:"apply_status"`
	LastApplyError                    types.String  `tfsdk:"last_apply_error"`
	ConfigGeneration                  types.Int64   `tfsdk:"config_generation"`
	AppliedGeneration                 types.Int64   `tfsdk:"applied_generation"`
	Transform                         types.Object  `tfsdk:"transform"`
}

func NewGcpPubsubOutletDataSource() datasource.DataSource { return &gcpPubsubOutletDataSource{} }

func (d *gcpPubsubOutletDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gcp_pubsub_outlet"
}

func (d *gcpPubsubOutletDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Reads a Noozle GCP Pub/Sub outlet by ID.",
		Attributes: map[string]datasourceschema.Attribute{
			"id":                                    datasourceschema.Int64Attribute{MarkdownDescription: "Unique identifier of the outlet.", Required: true},
			"name":                                  datasourceschema.StringAttribute{MarkdownDescription: "Name presented for the outlet.", Computed: true},
			"description":                           datasourceschema.StringAttribute{MarkdownDescription: "Optional note explaining the outlet purpose.", Computed: true},
			"enabled":                               datasourceschema.BoolAttribute{MarkdownDescription: "Whether the outlet is active.", Computed: true},
			"batching_byte_size":                    datasourceschema.Int64Attribute{MarkdownDescription: "Maximum bytes per Pub/Sub batch.", Computed: true},
			"batching_check":                        datasourceschema.StringAttribute{MarkdownDescription: "Bloblang condition for flushing a Pub/Sub batch.", Computed: true},
			"batching_count":                        datasourceschema.Int64Attribute{MarkdownDescription: "Maximum messages per Pub/Sub batch.", Computed: true},
			"batching_jitter":                       datasourceschema.Float64Attribute{MarkdownDescription: "Randomized jitter to apply to the Pub/Sub batching period.", Computed: true},
			"batching_period":                       datasourceschema.StringAttribute{MarkdownDescription: "Maximum wait period before flushing a Pub/Sub batch.", Computed: true},
			"byte_threshold":                        datasourceschema.Int64Attribute{MarkdownDescription: "Publish when the Pub/Sub buffer reaches this many bytes.", Computed: true},
			"count_threshold":                       datasourceschema.Int64Attribute{MarkdownDescription: "Publish when the Pub/Sub buffer reaches this many messages.", Computed: true},
			"delay_threshold":                       datasourceschema.StringAttribute{MarkdownDescription: "Publish a non-empty Pub/Sub buffer after this delay.", Computed: true},
			"endpoint":                              datasourceschema.StringAttribute{MarkdownDescription: "Optional Google Cloud Pub/Sub API endpoint override.", Computed: true},
			"flow_control_limit_exceeded_behavior":  datasourceschema.StringAttribute{MarkdownDescription: "Behavior when the Pub/Sub flow-control buffer is full.", Computed: true},
			"flow_control_max_outstanding_bytes":    datasourceschema.Int64Attribute{MarkdownDescription: "Maximum buffered Pub/Sub bytes waiting to be published.", Computed: true},
			"flow_control_max_outstanding_messages": datasourceschema.Int64Attribute{MarkdownDescription: "Maximum buffered Pub/Sub messages waiting to be published.", Computed: true},
			"max_in_flight":                         datasourceschema.Int64Attribute{MarkdownDescription: "Maximum number of in-flight Pub/Sub messages.", Computed: true},
			"metadata_exclude_prefixes":             datasourceschema.ListAttribute{MarkdownDescription: "Metadata prefixes to exclude from Pub/Sub attributes.", Computed: true, ElementType: types.StringType},
			"notify_policy":                         datasourceschema.StringAttribute{MarkdownDescription: "Delivery policy applied when notifying this outlet.", Computed: true},
			"ordering_key":                          datasourceschema.StringAttribute{MarkdownDescription: "Pub/Sub ordering key interpolation.", Computed: true},
			"project":                               datasourceschema.StringAttribute{MarkdownDescription: "Google Cloud project ID.", Computed: true},
			"publish_timeout":                       datasourceschema.StringAttribute{MarkdownDescription: "Maximum time to wait before abandoning a Pub/Sub publish attempt.", Computed: true},
			"topic":                                 datasourceschema.StringAttribute{MarkdownDescription: "Google Cloud Pub/Sub topic.", Computed: true},
			"apply_status":                          datasourceschema.StringAttribute{MarkdownDescription: "Observed apply/runtime status for the outlet.", Computed: true},
			"last_apply_error":                      datasourceschema.StringAttribute{MarkdownDescription: "Last observed apply/runtime error message.", Computed: true},
			"config_generation":                     datasourceschema.Int64Attribute{MarkdownDescription: "Desired outlet config generation.", Computed: true},
			"applied_generation":                    datasourceschema.Int64Attribute{MarkdownDescription: "Last observed applied config generation.", Computed: true},
			"transform": datasourceschema.SingleNestedAttribute{Computed: true, Attributes: map[string]datasourceschema.Attribute{
				"selected": datasourceschema.StringAttribute{Computed: true}, "effective": datasourceschema.StringAttribute{Computed: true}, "state": datasourceschema.StringAttribute{Computed: true}, "applied": datasourceschema.BoolAttribute{Computed: true}, "generation": datasourceschema.Int64Attribute{Computed: true}, "last_error": datasourceschema.StringAttribute{Computed: true},
			}},
		},
	}
}

func (d *gcpPubsubOutletDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	outletClient, ok := req.ProviderData.(gcpPubsubOutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected GCP Pub/Sub outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}
	d.client = outletClient
}

func (d *gcpPubsubOutletDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var config gcpPubsubOutletDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	outlet, err := d.client.GetGcpPubsubOutlet(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read GCP Pub/Sub Outlet", err.Error())
		return
	}
	transform, err := d.client.GetTypedOutletTransform(ctx, "gcp_pubsub", config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read GCP Pub/Sub Transform", err.Error())
		return
	}
	resourceState := gcpPubsubOutletModelFromAPI(outlet, types.ObjectNull(gcpCredentialsAttributeTypes), transform)
	state := gcpPubsubOutletDataSourceModel{
		ID:                                resourceState.ID,
		Name:                              resourceState.Name,
		Description:                       resourceState.Description,
		Enabled:                           resourceState.Enabled,
		BatchingByteSize:                  resourceState.BatchingByteSize,
		BatchingCheck:                     resourceState.BatchingCheck,
		BatchingCount:                     resourceState.BatchingCount,
		BatchingJitter:                    resourceState.BatchingJitter,
		BatchingPeriod:                    resourceState.BatchingPeriod,
		ByteThreshold:                     resourceState.ByteThreshold,
		CountThreshold:                    resourceState.CountThreshold,
		DelayThreshold:                    resourceState.DelayThreshold,
		Endpoint:                          resourceState.Endpoint,
		FlowControlLimitExceededBehavior:  resourceState.FlowControlLimitExceededBehavior,
		FlowControlMaxOutstandingBytes:    resourceState.FlowControlMaxOutstandingBytes,
		FlowControlMaxOutstandingMessages: resourceState.FlowControlMaxOutstandingMessages,
		MaxInFlight:                       resourceState.MaxInFlight,
		MetadataExcludePrefixes:           resourceState.MetadataExcludePrefixes,
		NotifyPolicy:                      resourceState.NotifyPolicy,
		OrderingKey:                       resourceState.OrderingKey,
		Project:                           resourceState.Project,
		PublishTimeout:                    resourceState.PublishTimeout,
		Topic:                             resourceState.Topic,
		ApplyStatus:                       resourceState.ApplyStatus,
		LastApplyError:                    resourceState.LastApplyError,
		ConfigGeneration:                  resourceState.ConfigGeneration,
		AppliedGeneration:                 resourceState.AppliedGeneration,
		Transform:                         resourceState.Transform,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

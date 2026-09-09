package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &awsSqsOutletDataSource{}
	_ datasource.DataSourceWithConfigure = &awsSqsOutletDataSource{}
)

type awsSqsOutletDataSource struct {
	client awsSqsOutletClient
}

type awsSqsOutletDataSourceModel struct {
	ID                      types.Int64   `tfsdk:"id"`
	Name                    types.String  `tfsdk:"name"`
	Description             types.String  `tfsdk:"description"`
	Enabled                 types.Bool    `tfsdk:"enabled"`
	BackoffInitialInterval  types.String  `tfsdk:"backoff_initial_interval"`
	BackoffMaxElapsedTime   types.String  `tfsdk:"backoff_max_elapsed_time"`
	BackoffMaxInterval      types.String  `tfsdk:"backoff_max_interval"`
	BatchingByteSize        types.Int64   `tfsdk:"batching_byte_size"`
	BatchingCheck           types.String  `tfsdk:"batching_check"`
	BatchingCount           types.Int64   `tfsdk:"batching_count"`
	BatchingJitter          types.Float64 `tfsdk:"batching_jitter"`
	BatchingPeriod          types.String  `tfsdk:"batching_period"`
	DelaySeconds            types.Int64   `tfsdk:"delay_seconds"`
	EndpointURL             types.String  `tfsdk:"endpoint_url"`
	ExpiryWindow            types.String  `tfsdk:"expiry_window"`
	ExternalID              types.String  `tfsdk:"external_id"`
	MaxInFlight             types.Int64   `tfsdk:"max_in_flight"`
	MaxRetries              types.Int64   `tfsdk:"max_retries"`
	MessageDeduplicationID  types.String  `tfsdk:"message_deduplication_id"`
	MessageGroupID          types.String  `tfsdk:"message_group_id"`
	MetadataExcludePrefixes types.List    `tfsdk:"metadata_exclude_prefixes"`
	NotifyPolicy            types.String  `tfsdk:"notify_policy"`
	Profile                 types.String  `tfsdk:"profile"`
	QueueURL                types.String  `tfsdk:"queue_url"`
	Region                  types.String  `tfsdk:"region"`
	RoleARN                 types.String  `tfsdk:"role_arn"`
	UseEC2InstanceRole      types.Bool    `tfsdk:"use_ec2_instance_role"`
	ApplyStatus             types.String  `tfsdk:"apply_status"`
	LastApplyError          types.String  `tfsdk:"last_apply_error"`
	ConfigGeneration        types.Int64   `tfsdk:"config_generation"`
	AppliedGeneration       types.Int64   `tfsdk:"applied_generation"`
}

func NewAwsSqsOutletDataSource() datasource.DataSource {
	return &awsSqsOutletDataSource{}
}

func (d *awsSqsOutletDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_sqs_outlet"
}

func (d *awsSqsOutletDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Reads a Noozle AWS SQS outlet by ID.",
		Attributes: map[string]datasourceschema.Attribute{
			"id":                        datasourceschema.Int64Attribute{MarkdownDescription: "Unique identifier of the outlet.", Required: true},
			"name":                      datasourceschema.StringAttribute{MarkdownDescription: "Name presented for the outlet.", Computed: true},
			"description":               datasourceschema.StringAttribute{MarkdownDescription: "Optional note explaining the outlet purpose.", Computed: true},
			"enabled":                   datasourceschema.BoolAttribute{MarkdownDescription: "Whether the outlet is active.", Computed: true},
			"backoff_initial_interval":  datasourceschema.StringAttribute{MarkdownDescription: "Initial SQS retry backoff interval.", Computed: true},
			"backoff_max_elapsed_time":  datasourceschema.StringAttribute{MarkdownDescription: "Maximum total SQS retry backoff time.", Computed: true},
			"backoff_max_interval":      datasourceschema.StringAttribute{MarkdownDescription: "Maximum SQS retry backoff interval.", Computed: true},
			"batching_byte_size":        datasourceschema.Int64Attribute{MarkdownDescription: "Maximum bytes per SQS batch.", Computed: true},
			"batching_check":            datasourceschema.StringAttribute{MarkdownDescription: "Bloblang condition for flushing an SQS batch.", Computed: true},
			"batching_count":            datasourceschema.Int64Attribute{MarkdownDescription: "Maximum messages per SQS batch.", Computed: true},
			"batching_jitter":           datasourceschema.Float64Attribute{MarkdownDescription: "Randomized jitter to apply to the SQS batching period.", Computed: true},
			"batching_period":           datasourceschema.StringAttribute{MarkdownDescription: "Maximum wait period before flushing an SQS batch.", Computed: true},
			"delay_seconds":             datasourceschema.Int64Attribute{MarkdownDescription: "Delay in seconds before an SQS message becomes visible.", Computed: true},
			"endpoint_url":              datasourceschema.StringAttribute{MarkdownDescription: "AWS service endpoint URL for SQS.", Computed: true},
			"expiry_window":             datasourceschema.StringAttribute{MarkdownDescription: "Expiry window for refreshed SQS credentials.", Computed: true},
			"external_id":               datasourceschema.StringAttribute{MarkdownDescription: "External ID for assuming an AWS IAM role.", Computed: true},
			"max_in_flight":             datasourceschema.Int64Attribute{MarkdownDescription: "Maximum number of in-flight SQS messages.", Computed: true},
			"max_retries":               datasourceschema.Int64Attribute{MarkdownDescription: "Maximum number of SQS send retries.", Computed: true},
			"message_deduplication_id":  datasourceschema.StringAttribute{MarkdownDescription: "AWS FIFO message deduplication ID.", Computed: true},
			"message_group_id":          datasourceschema.StringAttribute{MarkdownDescription: "AWS FIFO message group ID.", Computed: true},
			"metadata_exclude_prefixes": datasourceschema.ListAttribute{MarkdownDescription: "Metadata header prefixes to exclude from SQS messages.", Computed: true, ElementType: types.StringType},
			"notify_policy":             datasourceschema.StringAttribute{MarkdownDescription: "Delivery policy applied when notifying this outlet.", Computed: true},
			"profile":                   datasourceschema.StringAttribute{MarkdownDescription: "AWS shared credential profile name for SQS.", Computed: true},
			"queue_url":                 datasourceschema.StringAttribute{MarkdownDescription: "AWS SQS queue URL.", Computed: true},
			"region":                    datasourceschema.StringAttribute{MarkdownDescription: "AWS region for SQS.", Computed: true},
			"role_arn":                  datasourceschema.StringAttribute{MarkdownDescription: "AWS IAM role ARN to assume for SQS.", Computed: true},
			"use_ec2_instance_role":     datasourceschema.BoolAttribute{MarkdownDescription: "Use the EC2 instance role for SQS credentials.", Computed: true},
			"apply_status":              datasourceschema.StringAttribute{MarkdownDescription: "Observed apply/runtime status for the outlet.", Computed: true},
			"last_apply_error":          datasourceschema.StringAttribute{MarkdownDescription: "Last observed apply/runtime error message.", Computed: true},
			"config_generation":         datasourceschema.Int64Attribute{MarkdownDescription: "Desired outlet config generation.", Computed: true},
			"applied_generation":        datasourceschema.Int64Attribute{MarkdownDescription: "Last observed applied config generation.", Computed: true},
		},
	}
}

func (d *awsSqsOutletDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	outletClient, ok := req.ProviderData.(awsSqsOutletClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data Type", fmt.Sprintf("Expected AWS SQS outlet client, got: %T. Please report this provider bug.", req.ProviderData))
		return
	}
	d.client = outletClient
}

func (d *awsSqsOutletDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}
	var config awsSqsOutletDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	outlet, err := d.client.GetAwsSqsOutlet(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read AWS SQS Outlet", err.Error())
		return
	}
	resourceState := awsSqsOutletModelFromAPI(outlet, types.ObjectNull(awsCredentialsAttributeTypes))
	state := awsSqsOutletDataSourceModel{
		ID:                      resourceState.ID,
		Name:                    resourceState.Name,
		Description:             resourceState.Description,
		Enabled:                 resourceState.Enabled,
		BackoffInitialInterval:  resourceState.BackoffInitialInterval,
		BackoffMaxElapsedTime:   resourceState.BackoffMaxElapsedTime,
		BackoffMaxInterval:      resourceState.BackoffMaxInterval,
		BatchingByteSize:        resourceState.BatchingByteSize,
		BatchingCheck:           resourceState.BatchingCheck,
		BatchingCount:           resourceState.BatchingCount,
		BatchingJitter:          resourceState.BatchingJitter,
		BatchingPeriod:          resourceState.BatchingPeriod,
		DelaySeconds:            resourceState.DelaySeconds,
		EndpointURL:             resourceState.EndpointURL,
		ExpiryWindow:            resourceState.ExpiryWindow,
		ExternalID:              resourceState.ExternalID,
		MaxInFlight:             resourceState.MaxInFlight,
		MaxRetries:              resourceState.MaxRetries,
		MessageDeduplicationID:  resourceState.MessageDeduplicationID,
		MessageGroupID:          resourceState.MessageGroupID,
		MetadataExcludePrefixes: resourceState.MetadataExcludePrefixes,
		NotifyPolicy:            resourceState.NotifyPolicy,
		Profile:                 resourceState.Profile,
		QueueURL:                resourceState.QueueURL,
		Region:                  resourceState.Region,
		RoleARN:                 resourceState.RoleARN,
		UseEC2InstanceRole:      resourceState.UseEC2InstanceRole,
		ApplyStatus:             resourceState.ApplyStatus,
		LastApplyError:          resourceState.LastApplyError,
		ConfigGeneration:        resourceState.ConfigGeneration,
		AppliedGeneration:       resourceState.AppliedGeneration,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

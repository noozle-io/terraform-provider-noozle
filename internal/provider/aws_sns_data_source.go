package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &awsSnsOutletDataSource{}
	_ datasource.DataSourceWithConfigure = &awsSnsOutletDataSource{}
)

type awsSnsOutletDataSource struct {
	client awsSnsOutletClient
}

type awsSnsOutletDataSourceModel struct {
	ID                      types.Int64  `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
	EndpointURL             types.String `tfsdk:"endpoint_url"`
	ExpiryWindow            types.String `tfsdk:"expiry_window"`
	ExternalID              types.String `tfsdk:"external_id"`
	MaxInFlight             types.Int64  `tfsdk:"max_in_flight"`
	MessageDeduplicationID  types.String `tfsdk:"message_deduplication_id"`
	MessageGroupID          types.String `tfsdk:"message_group_id"`
	MetadataExcludePrefixes types.List   `tfsdk:"metadata_exclude_prefixes"`
	NotifyPolicy            types.String `tfsdk:"notify_policy"`
	Profile                 types.String `tfsdk:"profile"`
	Region                  types.String `tfsdk:"region"`
	RoleARN                 types.String `tfsdk:"role_arn"`
	Timeout                 types.String `tfsdk:"timeout"`
	TopicARN                types.String `tfsdk:"topic_arn"`
	UseEC2InstanceRole      types.Bool   `tfsdk:"use_ec2_instance_role"`
	ApplyStatus             types.String `tfsdk:"apply_status"`
	LastApplyError          types.String `tfsdk:"last_apply_error"`
	ConfigGeneration        types.Int64  `tfsdk:"config_generation"`
	AppliedGeneration       types.Int64  `tfsdk:"applied_generation"`
}

func NewAwsSnsOutletDataSource() datasource.DataSource {
	return &awsSnsOutletDataSource{}
}

func (d *awsSnsOutletDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aws_sns_outlet"
}

func (d *awsSnsOutletDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		MarkdownDescription: "Reads a Noozle AWS SNS outlet by ID.",
		Attributes: map[string]datasourceschema.Attribute{
			"id":                        datasourceschema.Int64Attribute{MarkdownDescription: "Unique identifier of the outlet.", Required: true},
			"name":                      datasourceschema.StringAttribute{MarkdownDescription: "Name presented for the outlet.", Computed: true},
			"description":               datasourceschema.StringAttribute{MarkdownDescription: "Optional note explaining the outlet purpose.", Computed: true},
			"enabled":                   datasourceschema.BoolAttribute{MarkdownDescription: "Whether the outlet is active.", Computed: true},
			"endpoint_url":              datasourceschema.StringAttribute{MarkdownDescription: "AWS service endpoint URL for SNS.", Computed: true},
			"expiry_window":             datasourceschema.StringAttribute{MarkdownDescription: "Expiry window for refreshed SNS credentials.", Computed: true},
			"external_id":               datasourceschema.StringAttribute{MarkdownDescription: "External ID for assuming an AWS IAM role.", Computed: true},
			"max_in_flight":             datasourceschema.Int64Attribute{MarkdownDescription: "Maximum number of in-flight SNS messages.", Computed: true},
			"message_deduplication_id":  datasourceschema.StringAttribute{MarkdownDescription: "AWS FIFO message deduplication ID.", Computed: true},
			"message_group_id":          datasourceschema.StringAttribute{MarkdownDescription: "AWS FIFO message group ID.", Computed: true},
			"metadata_exclude_prefixes": datasourceschema.ListAttribute{MarkdownDescription: "Metadata header prefixes to exclude from SNS messages.", Computed: true, ElementType: types.StringType},
			"notify_policy":             datasourceschema.StringAttribute{MarkdownDescription: "Delivery policy applied when notifying this outlet.", Computed: true},
			"profile":                   datasourceschema.StringAttribute{MarkdownDescription: "AWS shared credential profile name for SNS.", Computed: true},
			"region":                    datasourceschema.StringAttribute{MarkdownDescription: "AWS region for SNS.", Computed: true},
			"role_arn":                  datasourceschema.StringAttribute{MarkdownDescription: "AWS IAM role ARN to assume for SNS.", Computed: true},
			"timeout":                   datasourceschema.StringAttribute{MarkdownDescription: "Timeout for SNS publish requests.", Computed: true},
			"topic_arn":                 datasourceschema.StringAttribute{MarkdownDescription: "AWS SNS topic ARN.", Computed: true},
			"use_ec2_instance_role":     datasourceschema.BoolAttribute{MarkdownDescription: "Use the EC2 instance role for SNS credentials.", Computed: true},
			"apply_status":              datasourceschema.StringAttribute{MarkdownDescription: "Observed apply/runtime status for the outlet.", Computed: true},
			"last_apply_error":          datasourceschema.StringAttribute{MarkdownDescription: "Last observed apply/runtime error message.", Computed: true},
			"config_generation":         datasourceschema.Int64Attribute{MarkdownDescription: "Desired outlet config generation.", Computed: true},
			"applied_generation":        datasourceschema.Int64Attribute{MarkdownDescription: "Last observed applied config generation.", Computed: true},
		},
	}
}

func (d *awsSnsOutletDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	outletClient, ok := req.ProviderData.(awsSnsOutletClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected AWS SNS outlet client, got: %T. Please report this provider bug.", req.ProviderData),
		)
		return
	}

	d.client = outletClient
}

func (d *awsSnsOutletDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return
	}

	var config awsSnsOutletDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	outlet, err := d.client.GetAwsSnsOutlet(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read AWS SNS Outlet", err.Error())
		return
	}

	resourceState := awsSnsOutletModelFromAPI(outlet, types.ObjectNull(awsCredentialsAttributeTypes))
	state := awsSnsOutletDataSourceModel{
		ID:                      resourceState.ID,
		Name:                    resourceState.Name,
		Description:             resourceState.Description,
		Enabled:                 resourceState.Enabled,
		EndpointURL:             resourceState.EndpointURL,
		ExpiryWindow:            resourceState.ExpiryWindow,
		ExternalID:              resourceState.ExternalID,
		MaxInFlight:             resourceState.MaxInFlight,
		MessageDeduplicationID:  resourceState.MessageDeduplicationID,
		MessageGroupID:          resourceState.MessageGroupID,
		MetadataExcludePrefixes: resourceState.MetadataExcludePrefixes,
		NotifyPolicy:            resourceState.NotifyPolicy,
		Profile:                 resourceState.Profile,
		Region:                  resourceState.Region,
		RoleARN:                 resourceState.RoleARN,
		Timeout:                 resourceState.Timeout,
		TopicARN:                resourceState.TopicARN,
		UseEC2InstanceRole:      resourceState.UseEC2InstanceRole,
		ApplyStatus:             resourceState.ApplyStatus,
		LastApplyError:          resourceState.LastApplyError,
		ConfigGeneration:        resourceState.ConfigGeneration,
		AppliedGeneration:       resourceState.AppliedGeneration,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

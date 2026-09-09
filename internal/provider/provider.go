package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-noozle/internal/client"
)

// Ensure NoozleProvider satisfies the provider interface.
var _ provider.Provider = &NoozleProvider{}

// NoozleProvider defines the provider implementation.
type NoozleProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// NoozleProviderModel describes the provider data model.
type NoozleProviderModel struct {
	Host           types.String `tfsdk:"host"`
	APIKey         types.String `tfsdk:"api_key"`
	TenantID       types.String `tfsdk:"tenant_id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Debug          types.Bool   `tfsdk:"debug"`
}

var tenantIDPattern = regexp.MustCompile(`^t_[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}$`)

func (p *NoozleProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "noozle"
	resp.Version = p.version
}

func (p *NoozleProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "Noozle API base URL. Can also be set with `NOOZLE_HOST`.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Noozle API key. Can also be set with `NOOZLE_API_KEY`.",
				Optional:            true,
				Sensitive:           true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "Required Tenant context for query, outlet, and outlet-query-link operations. Set `tenant_id` or `NOOZLE_TENANT_ID`.",
				Optional:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Optional organization context used to disambiguate the selected tenant. Can also be set with `NOOZLE_ORGANIZATION_ID`.",
				Optional:            true,
			},
			"debug": schema.BoolAttribute{
				MarkdownDescription: "Enable debug logging for HTTP request and response metadata. Payload contents are never logged.",
				Optional:            true,
			},
		},
	}
}

func (p *NoozleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data NoozleProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Noozle API Host",
			"The provider cannot create the Noozle client because the host value is unknown.",
		)
	}

	if data.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Noozle API Key",
			"The provider cannot create the Noozle client because the API key value is unknown.",
		)
	}

	if data.TenantID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("tenant_id"),
			"Unknown Noozle Tenant ID",
			"The provider cannot create the Noozle client because the tenant_id value is unknown.",
		)
	}

	if data.OrganizationID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("organization_id"),
			"Unknown Noozle Organization ID",
			"The provider cannot create the Noozle client because the organization_id value is unknown.",
		)
	}

	if data.Debug.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("debug"),
			"Unknown Noozle Debug Setting",
			"The provider cannot create the Noozle client because the debug value is unknown.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := resolveConfigString(data.Host, "NOOZLE_HOST")
	apiKey := resolveConfigString(data.APIKey, "NOOZLE_API_KEY")
	tenantID := resolveConfigString(data.TenantID, "NOOZLE_TENANT_ID")
	organizationID := resolveConfigString(data.OrganizationID, "NOOZLE_ORGANIZATION_ID")
	debug := resolveConfigBool(data.Debug)

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Noozle API Host",
			"Set the provider `host` attribute or the `NOOZLE_HOST` environment variable.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Noozle API Key",
			"Set the provider `api_key` attribute or the `NOOZLE_API_KEY` environment variable.",
		)
	}

	if err := validateRequiredTenantContext(tenantID); err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("tenant_id"),
			"Missing Noozle Tenant ID",
			err.Error(),
		)
	} else if !tenantIDPattern.MatchString(tenantID) {
		resp.Diagnostics.AddAttributeError(
			path.Root("tenant_id"),
			"Invalid Noozle Tenant ID",
			"Set tenant_id or NOOZLE_TENANT_ID to a tenant identifier in the documented t_<ULID> format.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	clientOptions := []client.Option{client.WithTenantContext(client.TenantContext{
		TenantID:       tenantID,
		OrganizationID: organizationID,
	})}
	if debug {
		clientOptions = append(clientOptions, client.WithLogger(tflogClientLogger{}))
	}

	noozleClient, err := client.New(host, apiKey, clientOptions...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Noozle Client",
			fmt.Sprintf("An unexpected error occurred while creating the Noozle API client: %s", err),
		)
		return
	}

	resp.DataSourceData = noozleClient
	resp.ResourceData = noozleClient
}

func (p *NoozleProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewQueryResource,
		NewQueryOutletLinkResource,
		NewNtfyOutletResource,
		NewSseOutletResource,
		NewHttpServerOutletResource,
		NewAmqp09OutletResource,
		NewAmqp1OutletResource,
		NewKafkaOutletResource,
		NewNatsJetstreamOutletResource,
		NewPulsarOutletResource,
		NewAwsSnsOutletResource,
		NewAwsSqsOutletResource,
		NewGcpPubsubOutletResource,
	}
}

func (p *NoozleProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewQueryDataSource,
		NewQueryOutletLinkDataSource,
		NewNtfyOutletDataSource,
		NewSseOutletDataSource,
		NewHttpServerOutletDataSource,
		NewAmqp09OutletDataSource,
		NewAmqp1OutletDataSource,
		NewKafkaOutletDataSource,
		NewNatsJetstreamOutletDataSource,
		NewPulsarOutletDataSource,
		NewAwsSnsOutletDataSource,
		NewAwsSqsOutletDataSource,
		NewGcpPubsubOutletDataSource,
	}
}

func resolveConfigString(value types.String, envVar string) string {
	if !value.IsNull() {
		return strings.TrimSpace(value.ValueString())
	}

	return strings.TrimSpace(os.Getenv(envVar))
}

func resolveConfigBool(value types.Bool) bool {
	if value.IsNull() || value.IsUnknown() {
		return false
	}

	return value.ValueBool()
}

func validateRequiredTenantContext(tenantID string) error {
	if tenantID != "" {
		return nil
	}

	return errors.New("set the provider `tenant_id` attribute or the `NOOZLE_TENANT_ID` environment variable")
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NoozleProvider{
			version: version,
		}
	}
}

type tflogClientLogger struct{}

func (tflogClientLogger) Debug(ctx context.Context, message string, fields map[string]any) {
	tflog.Debug(ctx, message, fields)
}

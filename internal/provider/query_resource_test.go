package provider

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"terraform-provider-noozle/internal/client"
)

func TestResolveConfigString(t *testing.T) {
	t.Setenv("NOOZLE_HOST", " https://env.example.test/ ")

	got := resolveConfigString(types.StringNull(), "NOOZLE_HOST")
	if got != "https://env.example.test/" {
		t.Fatalf("unexpected env-resolved host: %q", got)
	}

	got = resolveConfigString(types.StringValue(" https://config.example.test "), "NOOZLE_HOST")
	if got != "https://config.example.test" {
		t.Fatalf("unexpected config-resolved host: %q", got)
	}
}

func TestValidateRequiredTenantContext(t *testing.T) {
	t.Parallel()

	if err := validateRequiredTenantContext(""); err == nil {
		t.Fatal("expected missing tenant context to be rejected")
	}
	if err := validateRequiredTenantContext("t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q"); err != nil {
		t.Fatalf("expected tenant context to be accepted: %v", err)
	}
}

func TestProviderConfigureRequiresTenantContextFromConfigurationOrEnvironment(t *testing.T) {
	providerConfig := func(tenantID *string) provider.ConfigureRequest {
		var tenantValue any
		if tenantID != nil {
			tenantValue = *tenantID
		}
		schema := providerschema.Schema{Attributes: map[string]providerschema.Attribute{
			"host":            providerschema.StringAttribute{Optional: true},
			"api_key":         providerschema.StringAttribute{Optional: true},
			"tenant_id":       providerschema.StringAttribute{Optional: true},
			"organization_id": providerschema.StringAttribute{Optional: true},
			"debug":           providerschema.BoolAttribute{Optional: true},
		}}
		rawType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"host": tftypes.String, "api_key": tftypes.String, "tenant_id": tftypes.String, "organization_id": tftypes.String, "debug": tftypes.Bool,
		}}
		return provider.ConfigureRequest{Config: tfsdk.Config{Schema: schema, Raw: tftypes.NewValue(rawType, map[string]tftypes.Value{
			"host":            tftypes.NewValue(tftypes.String, "https://noozle.example.test"),
			"api_key":         tftypes.NewValue(tftypes.String, "api-key"),
			"tenant_id":       tftypes.NewValue(tftypes.String, tenantValue),
			"organization_id": tftypes.NewValue(tftypes.String, nil),
			"debug":           tftypes.NewValue(tftypes.Bool, nil),
		})}}
	}

	t.Run("missing", func(t *testing.T) {
		resp := &provider.ConfigureResponse{}
		(&NoozleProvider{}).Configure(context.Background(), providerConfig(nil), resp)
		if !resp.Diagnostics.HasError() || resp.Diagnostics[0].Summary() != "Missing Noozle Tenant ID" {
			t.Fatalf("expected missing tenant context diagnostic, got %#v", resp.Diagnostics)
		}
	})

	t.Run("environment", func(t *testing.T) {
		t.Setenv("NOOZLE_TENANT_ID", "t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q")
		resp := &provider.ConfigureResponse{}
		(&NoozleProvider{}).Configure(context.Background(), providerConfig(nil), resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("configure provider with environment tenant context: %#v", resp.Diagnostics)
		}
	})

	t.Run("configuration", func(t *testing.T) {
		tenantID := "t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q"
		resp := &provider.ConfigureResponse{}
		(&NoozleProvider{}).Configure(context.Background(), providerConfig(&tenantID), resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("configure provider with tenant context: %#v", resp.Diagnostics)
		}
	})
}

func TestResolveConfigBool(t *testing.T) {
	t.Parallel()

	if resolveConfigBool(types.BoolNull()) {
		t.Fatal("expected null bool config to resolve to false")
	}
	if !resolveConfigBool(types.BoolValue(true)) {
		t.Fatal("expected true bool config to resolve to true")
	}
}

func TestProviderSchemaIncludesTenantContext(t *testing.T) {
	t.Parallel()

	p := &NoozleProvider{}
	resp := &provider.SchemaResponse{}
	p.Schema(context.Background(), provider.SchemaRequest{}, resp)

	for _, name := range []string{"tenant_id", "organization_id"} {
		attribute, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("provider schema is missing %q", name)
			continue
		}

		stringAttribute, ok := attribute.(providerschema.StringAttribute)
		if !ok || !stringAttribute.Optional {
			t.Errorf("provider schema attribute %q = %#v, want optional string", name, attribute)
		}
	}
}

func TestMissingTenantContextErrorProducesTerraformDiagnostic(t *testing.T) {
	t.Parallel()

	var diags diag.Diagnostics
	if !addMissingTenantContextDiagnosticForError(client.ErrMissingTenantContext, &diags) {
		t.Fatal("expected missing tenant context error to produce a diagnostic")
	}
	if !diags.HasError() {
		t.Fatal("expected missing tenant context diagnostic to be an error")
	}
	if got := diags[0].Summary(); got != "Missing Tenant Context" {
		t.Fatalf("unexpected diagnostic summary: %q", got)
	}
}

func TestQueryResourceModelToQueryRequest(t *testing.T) {
	t.Parallel()

	model := queryResourceModel{
		Name:           types.StringValue("Breaking News"),
		Description:    types.StringValue("desc"),
		FullExpression: types.StringValue("feed:finance"),
		Tags:           types.ListValueMust(types.StringType, []attr.Value{types.StringValue("finance"), types.StringValue("alerts")}),
		IsActive:       types.BoolValue(true),
	}

	var diags diag.Diagnostics
	req := model.toQueryRequest(context.Background(), &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if req.Name != "Breaking News" || req.FullExpression != "feed:finance" || len(req.Tags) != 2 {
		t.Fatalf("unexpected query request: %+v", req)
	}
}

func TestQueryModelFromAPI(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 27, 12, 34, 56, 0, time.UTC)
	state, diags := queryModelFromAPI(context.Background(), client.Query{
		ID:             12,
		Name:           "Breaking News",
		Description:    "desc",
		FullExpression: "feed:finance",
		Tags:           []string{"finance"},
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, queryResourceModel{})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if state.ID.ValueInt64() != 12 || state.CreatedAt.ValueString() != now.Format(time.RFC3339) {
		t.Fatalf("unexpected query state: %+v", state)
	}
}

func TestQueryModelFromAPIPreservesNullTags(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 27, 12, 34, 56, 0, time.UTC)
	state, diags := queryModelFromAPI(context.Background(), client.Query{
		ID:             12,
		Name:           "Breaking News",
		Description:    "desc",
		FullExpression: "feed:finance",
		Tags:           []string{},
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, queryResourceModel{
		Tags: types.ListNull(types.StringType),
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !state.Tags.IsNull() {
		t.Fatalf("expected null tags, got: %#v", state.Tags)
	}
}

func TestImportStateRejectsNonNumericID(t *testing.T) {
	t.Parallel()

	r := &queryResource{}
	resp := &resource.ImportStateResponse{}

	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "abc"}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostics for invalid import id")
	}
}

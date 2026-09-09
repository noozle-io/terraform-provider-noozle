package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

func TestHttpServerOutletModelFromAPI(t *testing.T) {
	t.Parallel()

	address := "127.0.0.1:4195"
	corsEnabled := true

	state := httpServerOutletModelFromAPI(client.HttpServerOutlet{
		ID:                 301,
		Name:               "HTTP Server",
		Description:        "desc",
		Enabled:            true,
		Address:            &address,
		AllowedVerbs:       []string{"POST", "GET"},
		CorsAllowedHeaders: []string{"Authorization"},
		CorsEnabled:        &corsEnabled,
		Path:               "/api/events",
	}, types.ObjectNull(httpServerTLSAttributeTypes), client.OutletTransformState{SelectedTransform: stringPointerForTest("summarize"), EffectiveTransform: stringPointerForTest("summarize"), Applied: true})

	if state.Address.ValueString() != "127.0.0.1:4195" {
		t.Fatalf("expected address to be set, got %s", state.Address)
	}
	if state.Path.ValueString() != "/api/events" {
		t.Fatalf("expected path to be set, got %s", state.Path)
	}
	if selected := selectedTransform(state.Transform); selected.ValueString() != "summarize" || !transformBoolAttribute(t, state.Transform, "applied").ValueBool() {
		t.Fatalf("unexpected transform state: %#v", state.Transform)
	}
}

func TestBuildHttpServerTLSObjectPreservesConfiguredBlockWhenRemoteIsEmpty(t *testing.T) {
	t.Parallel()

	previous, diags := types.ObjectValue(httpServerTLSAttributeTypes, map[string]attr.Value{
		"enabled":                     types.BoolValue(true),
		"state":                       types.StringValue("enabled"),
		"desired_generation":          types.Int64Value(2),
		"applied_generation":          types.Int64Value(2),
		"last_error":                  types.StringNull(),
		"server_cert_filename":        types.StringValue("cert.pem"),
		"server_cert_content_base64":  types.StringNull(),
		"server_cert_version":         types.Int64Value(1),
		"server_cert_checksum_sha256": types.StringNull(),
		"server_key_filename":         types.StringValue("key.pem"),
		"server_key_content_base64":   types.StringNull(),
		"server_key_version":          types.Int64Value(1),
		"server_key_checksum_sha256":  types.StringNull(),
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	result, diags := buildHttpServerTLSObject(context.Background(), previous, nil, nil)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if result.IsNull() {
		t.Fatal("expected TLS block to remain present")
	}
}

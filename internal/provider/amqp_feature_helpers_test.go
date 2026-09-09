package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

func TestBuildAMQPMTLSObjectIgnoresSharedCACertWithoutMTLSFeature(t *testing.T) {
	t.Parallel()

	result, diags := buildAMQPMTLSObject(context.Background(), types.ObjectNull(amqpMTLSAttributeTypes), nil, []client.OutletMaterialSummary{
		{Role: amqpCACertRole, Filename: "ca.pem"},
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !result.IsNull() {
		t.Fatal("expected mtls block to stay absent when only the shared ca_cert material is present")
	}
}

func TestBuildAMQPTLSObjectIgnoresSharedCACertWithoutTLSFeature(t *testing.T) {
	t.Parallel()

	result, diags := buildAMQPTLSObject(context.Background(), types.ObjectNull(amqpTLSAttributeTypes), nil, []client.OutletMaterialSummary{
		{Role: amqpCACertRole, Filename: "ca.pem"},
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !result.IsNull() {
		t.Fatal("expected tls block to stay absent when the feature is absent")
	}
}

func TestBuildAMQPTLSObjectIgnoresDisabledFeatureWithoutPriorBlock(t *testing.T) {
	t.Parallel()

	result, diags := buildAMQPTLSObject(context.Background(), types.ObjectNull(amqpTLSAttributeTypes), &client.OutletFeatureDetail{
		Feature: "amqp_0_9.tls",
		Enabled: false,
	}, nil)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !result.IsNull() {
		t.Fatal("expected tls block to stay absent when the remote feature is disabled and the block was not configured")
	}
}

func TestBuildAMQPMTLSDataSourceObjectIgnoresSharedCACertWithoutMTLSFeature(t *testing.T) {
	t.Parallel()

	result, diags := buildAMQPMTLSDataSourceObject(nil, []client.OutletMaterialSummary{
		{Role: amqpCACertRole, Filename: "ca.pem"},
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !result.IsNull() {
		t.Fatal("expected mtls data source block to stay absent when only the shared ca_cert material is present")
	}
}

func TestBuildAMQPTLSDataSourceObjectIgnoresSharedCACertWithoutTLSFeature(t *testing.T) {
	t.Parallel()

	result, diags := buildAMQPTLSDataSourceObject(nil, []client.OutletMaterialSummary{
		{Role: amqpCACertRole, Filename: "ca.pem"},
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !result.IsNull() {
		t.Fatal("expected tls data source block to stay absent when the feature is absent")
	}
}

func TestBuildAMQPTLSDataSourceObjectIgnoresDisabledFeature(t *testing.T) {
	t.Parallel()

	result, diags := buildAMQPTLSDataSourceObject(&client.OutletFeatureDetail{
		Feature: "amqp_0_9.tls",
		Enabled: false,
	}, []client.OutletMaterialSummary{
		{Role: amqpCACertRole, Filename: "ca.pem"},
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !result.IsNull() {
		t.Fatal("expected tls data source block to stay absent when the remote feature is disabled")
	}
}

func TestBuildAMQPMTLSObjectPreservesConfiguredBlockWhenRemoteFeatureIsLagging(t *testing.T) {
	t.Parallel()

	previous, diags := types.ObjectValue(amqpMTLSAttributeTypes, map[string]attr.Value{
		"enabled":                     types.BoolNull(),
		"state":                       types.StringNull(),
		"desired_generation":          types.Int64Null(),
		"applied_generation":          types.Int64Null(),
		"last_error":                  types.StringNull(),
		"ca_cert_filename":            types.StringValue("ca.pem"),
		"ca_cert_content_base64":      types.StringNull(),
		"ca_cert_version":             types.Int64Value(1),
		"ca_cert_checksum_sha256":     types.StringNull(),
		"client_cert_filename":        types.StringValue("client.pem"),
		"client_cert_content_base64":  types.StringNull(),
		"client_cert_version":         types.Int64Value(1),
		"client_cert_checksum_sha256": types.StringNull(),
		"client_key_filename":         types.StringValue("client-key.pem"),
		"client_key_content_base64":   types.StringNull(),
		"client_key_version":          types.Int64Value(1),
		"client_key_checksum_sha256":  types.StringNull(),
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics building previous object: %v", diags)
	}

	result, diags := buildAMQPMTLSObject(context.Background(), previous, nil, []client.OutletMaterialSummary{
		{Role: amqpCACertRole, Filename: "ca.pem"},
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if result.IsNull() {
		t.Fatal("expected configured mtls block to remain present while remote feature state catches up")
	}
}

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

func TestBuildGcpCredentialsObjectPreservesConfiguredBlockWhenRemoteIsEmpty(t *testing.T) {
	t.Parallel()

	previous, diags := types.ObjectValue(gcpCredentialsAttributeTypes, map[string]attr.Value{
		"filename":        types.StringValue("service-account.json"),
		"content_base64":  types.StringNull(),
		"content_version": types.Int64Value(1),
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics building previous object: %v", diags)
	}

	result, diags := buildGcpCredentialsObject(context.Background(), previous, nil)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics rebuilding credentials object: %v", diags)
	}

	if result.IsNull() {
		t.Fatal("expected credentials block to remain present")
	}

	credentials, diags := expandGcpCredentials(context.Background(), result)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics expanding result: %v", diags)
	}
	if credentials.Filename.ValueString() != "service-account.json" {
		t.Fatalf("expected filename to be preserved, got %q", credentials.Filename.ValueString())
	}
	if credentials.ContentVersion.ValueInt64() != 1 {
		t.Fatalf("expected content_version to be preserved, got %d", credentials.ContentVersion.ValueInt64())
	}
}

func TestBuildAwsCredentialsObjectPreservesConfiguredBlockWhenRemoteIsEmpty(t *testing.T) {
	t.Parallel()

	previous, diags := types.ObjectValue(awsCredentialsAttributeTypes, map[string]attr.Value{
		"access_key_id":             types.StringNull(),
		"access_key_id_version":     types.Int64Value(1),
		"secret_access_key":         types.StringNull(),
		"secret_access_key_version": types.Int64Value(2),
		"session_token":             types.StringNull(),
		"session_token_version":     types.Int64Null(),
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics building previous object: %v", diags)
	}

	result, diags := buildAwsCredentialsObject(context.Background(), previous, nil)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics rebuilding credentials object: %v", diags)
	}

	if result.IsNull() {
		t.Fatal("expected credentials block to remain present")
	}

	credentials, diags := expandAwsCredentials(context.Background(), result)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics expanding result: %v", diags)
	}
	if credentials.AccessKeyIDVersion.ValueInt64() != 1 {
		t.Fatalf("expected access_key_id_version to be preserved, got %d", credentials.AccessKeyIDVersion.ValueInt64())
	}
	if credentials.SecretAccessKeyVersion.ValueInt64() != 2 {
		t.Fatalf("expected secret_access_key_version to be preserved, got %d", credentials.SecretAccessKeyVersion.ValueInt64())
	}
	if !credentials.SessionTokenVersion.IsNull() {
		t.Fatal("expected null session_token_version to remain null")
	}
}

func TestBuildGcpCredentialsObjectUsesRemoteFilenameWhenAvailable(t *testing.T) {
	t.Parallel()

	previous, diags := types.ObjectValue(gcpCredentialsAttributeTypes, map[string]attr.Value{
		"filename":        types.StringValue("old.json"),
		"content_base64":  types.StringNull(),
		"content_version": types.Int64Value(3),
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics building previous object: %v", diags)
	}

	result, diags := buildGcpCredentialsObject(context.Background(), previous, []client.OutletMaterialSummary{
		{Role: gcpPubsubCredentialsRole, Filename: "remote.json"},
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics rebuilding credentials object: %v", diags)
	}

	credentials, diags := expandGcpCredentials(context.Background(), result)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics expanding result: %v", diags)
	}
	if credentials.Filename.ValueString() != "remote.json" {
		t.Fatalf("expected remote filename to win, got %q", credentials.Filename.ValueString())
	}
	if credentials.ContentVersion.ValueInt64() != 3 {
		t.Fatalf("expected content_version to be preserved, got %d", credentials.ContentVersion.ValueInt64())
	}
}

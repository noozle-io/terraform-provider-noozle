package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"terraform-provider-noozle/internal/client"
)

type awsCredentialsModel struct {
	AccessKeyID            types.String `tfsdk:"access_key_id"`
	AccessKeyIDVersion     types.Int64  `tfsdk:"access_key_id_version"`
	SecretAccessKey        types.String `tfsdk:"secret_access_key"`
	SecretAccessKeyVersion types.Int64  `tfsdk:"secret_access_key_version"`
	SessionToken           types.String `tfsdk:"session_token"`
	SessionTokenVersion    types.Int64  `tfsdk:"session_token_version"`
}

type gcpCredentialsModel struct {
	Filename       types.String `tfsdk:"filename"`
	ContentBase64  types.String `tfsdk:"content_base64"`
	ContentVersion types.Int64  `tfsdk:"content_version"`
}

var awsCredentialsAttributeTypes = map[string]attr.Type{
	"access_key_id":             types.StringType,
	"access_key_id_version":     types.Int64Type,
	"secret_access_key":         types.StringType,
	"secret_access_key_version": types.Int64Type,
	"session_token":             types.StringType,
	"session_token_version":     types.Int64Type,
}

var gcpCredentialsAttributeTypes = map[string]attr.Type{
	"filename":        types.StringType,
	"content_base64":  types.StringType,
	"content_version": types.Int64Type,
}

func expandAwsCredentials(ctx context.Context, object types.Object) (awsCredentialsModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return awsCredentialsModel{
			AccessKeyID:            types.StringNull(),
			AccessKeyIDVersion:     types.Int64Null(),
			SecretAccessKey:        types.StringNull(),
			SecretAccessKeyVersion: types.Int64Null(),
			SessionToken:           types.StringNull(),
			SessionTokenVersion:    types.Int64Null(),
		}, nil
	}

	var credentials awsCredentialsModel
	diags := object.As(ctx, &credentials, basetypes.ObjectAsOptions{})
	return credentials, diags
}

func expandGcpCredentials(ctx context.Context, object types.Object) (gcpCredentialsModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return gcpCredentialsModel{
			Filename:       types.StringNull(),
			ContentBase64:  types.StringNull(),
			ContentVersion: types.Int64Null(),
		}, nil
	}

	var credentials gcpCredentialsModel
	diags := object.As(ctx, &credentials, basetypes.ObjectAsOptions{})
	return credentials, diags
}

func buildAwsCredentialsObject(ctx context.Context, previous types.Object, remote []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandAwsCredentials(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(awsCredentialsAttributeTypes), diags
	}

	hasAccessKey := hasCredentialRole(remote, "access_key_id")
	hasSecretAccessKey := hasCredentialRole(remote, "secret_access_key")
	hasSessionToken := hasCredentialRole(remote, "session_token")
	if !hasAccessKey && !hasSecretAccessKey && !hasSessionToken {
		if previous.IsNull() || previous.IsUnknown() {
			return types.ObjectNull(awsCredentialsAttributeTypes), nil
		}

		return types.ObjectValue(awsCredentialsAttributeTypes, map[string]attr.Value{
			"access_key_id":             types.StringNull(),
			"access_key_id_version":     previousModel.AccessKeyIDVersion,
			"secret_access_key":         types.StringNull(),
			"secret_access_key_version": previousModel.SecretAccessKeyVersion,
			"session_token":             types.StringNull(),
			"session_token_version":     previousModel.SessionTokenVersion,
		})
	}

	return types.ObjectValue(awsCredentialsAttributeTypes, map[string]attr.Value{
		"access_key_id":             types.StringNull(),
		"access_key_id_version":     versionOrNull(hasAccessKey, previousModel.AccessKeyIDVersion),
		"secret_access_key":         types.StringNull(),
		"secret_access_key_version": versionOrNull(hasSecretAccessKey, previousModel.SecretAccessKeyVersion),
		"session_token":             types.StringNull(),
		"session_token_version":     versionOrNull(hasSessionToken, previousModel.SessionTokenVersion),
	})
}

func buildGcpCredentialsObject(ctx context.Context, previous types.Object, remote []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandGcpCredentials(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(gcpCredentialsAttributeTypes), diags
	}

	credential := findMaterialOrNil(remote, gcpPubsubCredentialsRole)
	if credential == nil {
		if previous.IsNull() || previous.IsUnknown() {
			return types.ObjectNull(gcpCredentialsAttributeTypes), nil
		}

		return types.ObjectValue(gcpCredentialsAttributeTypes, map[string]attr.Value{
			"filename":        previousModel.Filename,
			"content_base64":  types.StringNull(),
			"content_version": previousModel.ContentVersion,
		})
	}

	filename := previousModel.Filename
	if credential.Filename != "" {
		filename = types.StringValue(credential.Filename)
	}

	return types.ObjectValue(gcpCredentialsAttributeTypes, map[string]attr.Value{
		"filename":        filename,
		"content_base64":  types.StringNull(),
		"content_version": previousModel.ContentVersion,
	})
}

func hasCredentialRole(remote []client.OutletMaterialSummary, role string) bool {
	return findMaterialOrNil(remote, role) != nil
}

func versionOrNull(present bool, version types.Int64) types.Int64 {
	if !present {
		return types.Int64Null()
	}
	if version.IsNull() || version.IsUnknown() {
		return types.Int64Null()
	}

	return version
}

func awsCredentialRolesToDelete(plan, state awsCredentialsModel) []string {
	roles := make([]string, 0, 3)

	if plan.AccessKeyIDVersion.IsNull() && !state.AccessKeyIDVersion.IsNull() {
		roles = append(roles, "access_key_id")
	}
	if plan.SecretAccessKeyVersion.IsNull() && !state.SecretAccessKeyVersion.IsNull() {
		roles = append(roles, "secret_access_key")
	}
	if plan.SessionTokenVersion.IsNull() && !state.SessionTokenVersion.IsNull() {
		roles = append(roles, "session_token")
	}

	return roles
}

func gcpCredentialRemoved(plan, state gcpCredentialsModel) bool {
	return plan.ContentVersion.IsNull() && !state.ContentVersion.IsNull()
}

func enrichAwsCredentialsFromConfig(ctx context.Context, config tfsdk.Config, credentials awsCredentialsModel) (awsCredentialsModel, diag.Diagnostics) {
	if credentials.AccessKeyID.IsNull() || credentials.AccessKeyID.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root("credentials").AtName("access_key_id"), &value)
		if diags.HasError() {
			return credentials, diags
		}
		credentials.AccessKeyID = value
	}

	if credentials.SecretAccessKey.IsNull() || credentials.SecretAccessKey.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root("credentials").AtName("secret_access_key"), &value)
		if diags.HasError() {
			return credentials, diags
		}
		credentials.SecretAccessKey = value
	}

	if credentials.SessionToken.IsNull() || credentials.SessionToken.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root("credentials").AtName("session_token"), &value)
		if diags.HasError() {
			return credentials, diags
		}
		credentials.SessionToken = value
	}

	return credentials, nil
}

func enrichGcpCredentialsFromConfig(ctx context.Context, config tfsdk.Config, credentials gcpCredentialsModel) (gcpCredentialsModel, diag.Diagnostics) {
	if credentials.ContentBase64.IsNull() || credentials.ContentBase64.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root("credentials").AtName("content_base64"), &value)
		if diags.HasError() {
			return credentials, diags
		}
		credentials.ContentBase64 = value
	}

	return credentials, nil
}

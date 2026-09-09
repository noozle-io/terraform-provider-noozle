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

type amqpURLsCredentialsModel struct {
	URLs        types.List  `tfsdk:"urls"`
	URLsVersion types.Int64 `tfsdk:"urls_version"`
}

var amqpURLsCredentialsAttributeTypes = map[string]attr.Type{
	"urls":         types.ListType{ElemType: types.StringType},
	"urls_version": types.Int64Type,
}

func expandAmqpURLsCredentials(ctx context.Context, object types.Object) (amqpURLsCredentialsModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return amqpURLsCredentialsModel{
			URLs:        types.ListNull(types.StringType),
			URLsVersion: types.Int64Null(),
		}, nil
	}

	var credentials amqpURLsCredentialsModel
	diags := object.As(ctx, &credentials, basetypes.ObjectAsOptions{})
	return credentials, diags
}

func buildAmqpURLsCredentialsObject(ctx context.Context, previous types.Object, remote []client.OutletMaterialSummary, role string) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandAmqpURLsCredentials(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(amqpURLsCredentialsAttributeTypes), diags
	}

	credential := findMaterialOrNil(remote, role)
	if credential == nil {
		if previous.IsNull() || previous.IsUnknown() {
			return types.ObjectNull(amqpURLsCredentialsAttributeTypes), nil
		}

		return types.ObjectValue(amqpURLsCredentialsAttributeTypes, map[string]attr.Value{
			"urls":         types.ListNull(types.StringType),
			"urls_version": previousModel.URLsVersion,
		})
	}

	return types.ObjectValue(amqpURLsCredentialsAttributeTypes, map[string]attr.Value{
		"urls":         types.ListNull(types.StringType),
		"urls_version": previousModel.URLsVersion,
	})
}

func enrichAmqpURLsCredentialsFromConfig(ctx context.Context, config tfsdk.Config, credentials amqpURLsCredentialsModel) (amqpURLsCredentialsModel, diag.Diagnostics) {
	if credentials.URLs.IsNull() || credentials.URLs.IsUnknown() {
		var value types.List
		diags := config.GetAttribute(ctx, path.Root("credentials").AtName("urls"), &value)
		if diags.HasError() {
			return credentials, diags
		}
		credentials.URLs = value
	}

	return credentials, nil
}

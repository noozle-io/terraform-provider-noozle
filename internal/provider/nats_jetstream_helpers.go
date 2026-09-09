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

const (
	natsJetstreamURLsCredentialRole = "nats_jetstream_urls"
	natsJetstreamJWTUserRole        = "nats_user_jwt"
	natsJetstreamNkeyFileRole       = "nkey_file"
	natsJetstreamNkeySeedRole       = "nkey_seed"
	natsJetstreamUserCredsRole      = "user_credentials"
)

type natsJetstreamJWTSeedModel struct {
	Enabled           types.Bool   `tfsdk:"enabled"`
	State             types.String `tfsdk:"state"`
	DesiredGeneration types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration types.Int64  `tfsdk:"applied_generation"`
	LastError         types.String `tfsdk:"last_error"`
	UserJWT           types.String `tfsdk:"user_jwt"`
	UserJWTVersion    types.Int64  `tfsdk:"user_jwt_version"`
	NkeySeed          types.String `tfsdk:"nkey_seed"`
	NkeySeedVersion   types.Int64  `tfsdk:"nkey_seed_version"`
}

type natsJetstreamNkeyModel struct {
	Enabled           types.Bool   `tfsdk:"enabled"`
	State             types.String `tfsdk:"state"`
	DesiredGeneration types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration types.Int64  `tfsdk:"applied_generation"`
	LastError         types.String `tfsdk:"last_error"`
	Filename          types.String `tfsdk:"filename"`
	ContentBase64     types.String `tfsdk:"content_base64"`
	Version           types.Int64  `tfsdk:"version"`
	ChecksumSHA256    types.String `tfsdk:"checksum_sha256"`
}

type natsJetstreamUserCredentialsModel struct {
	Enabled           types.Bool   `tfsdk:"enabled"`
	State             types.String `tfsdk:"state"`
	DesiredGeneration types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration types.Int64  `tfsdk:"applied_generation"`
	LastError         types.String `tfsdk:"last_error"`
	Filename          types.String `tfsdk:"filename"`
	ContentBase64     types.String `tfsdk:"content_base64"`
	Version           types.Int64  `tfsdk:"version"`
	ChecksumSHA256    types.String `tfsdk:"checksum_sha256"`
}

var natsJetstreamJWTSeedAttributeTypes = map[string]attr.Type{
	"enabled":            types.BoolType,
	"state":              types.StringType,
	"desired_generation": types.Int64Type,
	"applied_generation": types.Int64Type,
	"last_error":         types.StringType,
	"user_jwt":           types.StringType,
	"user_jwt_version":   types.Int64Type,
	"nkey_seed":          types.StringType,
	"nkey_seed_version":  types.Int64Type,
}

var natsJetstreamNkeyAttributeTypes = map[string]attr.Type{
	"enabled":            types.BoolType,
	"state":              types.StringType,
	"desired_generation": types.Int64Type,
	"applied_generation": types.Int64Type,
	"last_error":         types.StringType,
	"filename":           types.StringType,
	"content_base64":     types.StringType,
	"version":            types.Int64Type,
	"checksum_sha256":    types.StringType,
}

var natsJetstreamUserCredentialsAttributeTypes = map[string]attr.Type{
	"enabled":            types.BoolType,
	"state":              types.StringType,
	"desired_generation": types.Int64Type,
	"applied_generation": types.Int64Type,
	"last_error":         types.StringType,
	"filename":           types.StringType,
	"content_base64":     types.StringType,
	"version":            types.Int64Type,
	"checksum_sha256":    types.StringType,
}

func expandNatsJetstreamJWTSeed(ctx context.Context, object types.Object) (natsJetstreamJWTSeedModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return natsJetstreamJWTSeedModel{
			Enabled:           types.BoolNull(),
			State:             types.StringNull(),
			DesiredGeneration: types.Int64Null(),
			AppliedGeneration: types.Int64Null(),
			LastError:         types.StringNull(),
			UserJWT:           types.StringNull(),
			UserJWTVersion:    types.Int64Null(),
			NkeySeed:          types.StringNull(),
			NkeySeedVersion:   types.Int64Null(),
		}, nil
	}
	var model natsJetstreamJWTSeedModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func expandNatsJetstreamNkey(ctx context.Context, object types.Object) (natsJetstreamNkeyModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return natsJetstreamNkeyModel{
			Enabled:           types.BoolNull(),
			State:             types.StringNull(),
			DesiredGeneration: types.Int64Null(),
			AppliedGeneration: types.Int64Null(),
			LastError:         types.StringNull(),
			Filename:          types.StringNull(),
			ContentBase64:     types.StringNull(),
			Version:           types.Int64Null(),
			ChecksumSHA256:    types.StringNull(),
		}, nil
	}
	var model natsJetstreamNkeyModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func expandNatsJetstreamUserCredentials(ctx context.Context, object types.Object) (natsJetstreamUserCredentialsModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return natsJetstreamUserCredentialsModel{
			Enabled:           types.BoolNull(),
			State:             types.StringNull(),
			DesiredGeneration: types.Int64Null(),
			AppliedGeneration: types.Int64Null(),
			LastError:         types.StringNull(),
			Filename:          types.StringNull(),
			ContentBase64:     types.StringNull(),
			Version:           types.Int64Null(),
			ChecksumSHA256:    types.StringNull(),
		}, nil
	}
	var model natsJetstreamUserCredentialsModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func enrichNatsJetstreamJWTSeedFromConfig(ctx context.Context, config tfsdk.Config, model natsJetstreamJWTSeedModel) (natsJetstreamJWTSeedModel, diag.Diagnostics) {
	if model.UserJWT.IsNull() || model.UserJWT.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root("jwt_seed").AtName("user_jwt"), &value)
		if diags.HasError() {
			return model, diags
		}
		model.UserJWT = value
	}
	if model.NkeySeed.IsNull() || model.NkeySeed.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root("jwt_seed").AtName("nkey_seed"), &value)
		if diags.HasError() {
			return model, diags
		}
		model.NkeySeed = value
	}
	return model, nil
}

func enrichNatsJetstreamNkeyFromConfig(ctx context.Context, config tfsdk.Config, model natsJetstreamNkeyModel) (natsJetstreamNkeyModel, diag.Diagnostics) {
	if model.ContentBase64.IsNull() || model.ContentBase64.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root("nkey").AtName("content_base64"), &value)
		if diags.HasError() {
			return model, diags
		}
		model.ContentBase64 = value
	}
	return model, nil
}

func enrichNatsJetstreamUserCredentialsFromConfig(ctx context.Context, config tfsdk.Config, model natsJetstreamUserCredentialsModel) (natsJetstreamUserCredentialsModel, diag.Diagnostics) {
	if model.ContentBase64.IsNull() || model.ContentBase64.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root("user_credentials").AtName("content_base64"), &value)
		if diags.HasError() {
			return model, diags
		}
		model.ContentBase64 = value
	}
	return model, nil
}

func buildNatsFeatureBase(feature *client.OutletFeatureDetail) map[string]attr.Value {
	enabled := types.BoolNull()
	state := types.StringNull()
	desiredGeneration := types.Int64Null()
	appliedGeneration := types.Int64Null()
	lastError := types.StringNull()
	if feature != nil {
		enabled = types.BoolValue(feature.Enabled)
		state = nullableString(feature.State)
		desiredGeneration = nullableInt64Pointer(feature.DesiredGeneration)
		appliedGeneration = nullableInt64Pointer(feature.AppliedGeneration)
		lastError = nullableString(feature.LastError)
	}
	return map[string]attr.Value{
		"enabled":            enabled,
		"state":              state,
		"desired_generation": desiredGeneration,
		"applied_generation": appliedGeneration,
		"last_error":         lastError,
	}
}

func buildNatsJetstreamJWTSeedObject(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandNatsJetstreamJWTSeed(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(natsJetstreamJWTSeedAttributeTypes), diags
	}
	jwt := findMaterialOrNil(credentials, natsJetstreamJWTUserRole)
	nkeySeed := findMaterialOrNil(credentials, natsJetstreamNkeySeedRole)
	if previous.IsNull() && (feature == nil || !feature.Enabled) {
		return types.ObjectNull(natsJetstreamJWTSeedAttributeTypes), nil
	}
	values := buildNatsFeatureBase(feature)
	if feature == nil && jwt != nil && nkeySeed != nil {
		values["enabled"] = types.BoolValue(true)
	}
	values["user_jwt"] = types.StringNull()
	values["user_jwt_version"] = previousModel.UserJWTVersion
	values["nkey_seed"] = types.StringNull()
	values["nkey_seed_version"] = previousModel.NkeySeedVersion
	return types.ObjectValue(natsJetstreamJWTSeedAttributeTypes, values)
}

func buildNatsJetstreamNkeyObject(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandNatsJetstreamNkey(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(natsJetstreamNkeyAttributeTypes), diags
	}
	credential := findMaterialOrNil(credentials, natsJetstreamNkeyFileRole)
	if previous.IsNull() && (feature == nil || !feature.Enabled) {
		return types.ObjectNull(natsJetstreamNkeyAttributeTypes), nil
	}
	values := buildNatsFeatureBase(feature)
	if feature == nil && credential != nil {
		values["enabled"] = types.BoolValue(true)
	}
	filename := previousModel.Filename
	checksum := previousModel.ChecksumSHA256
	if credential != nil {
		filename = nullableString(credential.Filename)
		checksum = nullableString(credential.ChecksumSHA256)
	}
	values["filename"] = filename
	values["content_base64"] = types.StringNull()
	values["version"] = previousModel.Version
	values["checksum_sha256"] = checksum
	return types.ObjectValue(natsJetstreamNkeyAttributeTypes, values)
}

func buildNatsJetstreamUserCredentialsObject(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandNatsJetstreamUserCredentials(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(natsJetstreamUserCredentialsAttributeTypes), diags
	}
	credential := findMaterialOrNil(credentials, natsJetstreamUserCredsRole)
	if previous.IsNull() && (feature == nil || !feature.Enabled) {
		return types.ObjectNull(natsJetstreamUserCredentialsAttributeTypes), nil
	}
	values := buildNatsFeatureBase(feature)
	if feature == nil && credential != nil {
		values["enabled"] = types.BoolValue(true)
	}
	filename := previousModel.Filename
	checksum := previousModel.ChecksumSHA256
	if credential != nil {
		filename = nullableString(credential.Filename)
		checksum = nullableString(credential.ChecksumSHA256)
	}
	values["filename"] = filename
	values["content_base64"] = types.StringNull()
	values["version"] = previousModel.Version
	values["checksum_sha256"] = checksum
	return types.ObjectValue(natsJetstreamUserCredentialsAttributeTypes, values)
}

func buildNatsJetstreamJWTSeedDataSourceObject(feature *client.OutletFeatureDetail) (types.Object, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"enabled":            types.BoolType,
		"state":              types.StringType,
		"desired_generation": types.Int64Type,
		"applied_generation": types.Int64Type,
		"last_error":         types.StringType,
		"user_jwt_version":   types.Int64Type,
		"nkey_seed_version":  types.Int64Type,
	}
	if feature == nil || !feature.Enabled {
		return types.ObjectNull(attrTypes), nil
	}
	values := buildNatsFeatureBase(feature)
	values["user_jwt_version"] = types.Int64Null()
	values["nkey_seed_version"] = types.Int64Null()
	return types.ObjectValue(attrTypes, values)
}

func buildNatsJetstreamNkeyDataSourceObject(feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	credential := findMaterialOrNil(credentials, natsJetstreamNkeyFileRole)
	attrTypes := map[string]attr.Type{
		"enabled":            types.BoolType,
		"state":              types.StringType,
		"desired_generation": types.Int64Type,
		"applied_generation": types.Int64Type,
		"last_error":         types.StringType,
		"filename":           types.StringType,
		"checksum_sha256":    types.StringType,
	}
	if (feature == nil || !feature.Enabled) && credential == nil {
		return types.ObjectNull(attrTypes), nil
	}
	values := buildNatsFeatureBase(feature)
	values["filename"] = nullableString(credentialField(credential, func(m *client.OutletMaterialSummary) string { return m.Filename }))
	values["checksum_sha256"] = nullableString(credentialField(credential, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 }))
	return types.ObjectValue(attrTypes, values)
}

func buildNatsJetstreamUserCredentialsDataSourceObject(feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	credential := findMaterialOrNil(credentials, natsJetstreamUserCredsRole)
	attrTypes := map[string]attr.Type{
		"enabled":            types.BoolType,
		"state":              types.StringType,
		"desired_generation": types.Int64Type,
		"applied_generation": types.Int64Type,
		"last_error":         types.StringType,
		"filename":           types.StringType,
		"checksum_sha256":    types.StringType,
	}
	if (feature == nil || !feature.Enabled) && credential == nil {
		return types.ObjectNull(attrTypes), nil
	}
	values := buildNatsFeatureBase(feature)
	values["filename"] = nullableString(credentialField(credential, func(m *client.OutletMaterialSummary) string { return m.Filename }))
	values["checksum_sha256"] = nullableString(credentialField(credential, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 }))
	return types.ObjectValue(attrTypes, values)
}

func natsFeatureEnabled(object types.Object) bool {
	if object.IsNull() || object.IsUnknown() {
		return false
	}
	attributes := object.Attributes()
	enabled, ok := attributes["enabled"]
	if !ok {
		return false
	}
	boolValue, ok := enabled.(types.Bool)
	if !ok {
		return false
	}
	return !boolValue.IsNull() && !boolValue.IsUnknown() && boolValue.ValueBool()
}

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
	pulsarTLSFeatureID         = "pulsar.tls"
	pulsarOAuth2FeatureID      = "pulsar.oauth2"
	pulsarTokenFeatureID       = "pulsar.token"
	pulsarTLSCARole            = "ca_cert"
	pulsarOAuth2PrivateKeyRole = "oauth2_private_key"
	pulsarTokenRole            = "pulsar_access_token"
)

type pulsarTLSModel struct {
	Enabled             types.Bool   `tfsdk:"enabled"`
	State               types.String `tfsdk:"state"`
	DesiredGeneration   types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration   types.Int64  `tfsdk:"applied_generation"`
	LastError           types.String `tfsdk:"last_error"`
	CACertFilename      types.String `tfsdk:"ca_cert_filename"`
	CACertContentBase64 types.String `tfsdk:"ca_cert_content_base64"`
	CACertVersion       types.Int64  `tfsdk:"ca_cert_version"`
	CACertChecksumSHA   types.String `tfsdk:"ca_cert_checksum_sha256"`
}

var pulsarTLSAttributeTypes = map[string]attr.Type{
	"enabled":                 types.BoolType,
	"state":                   types.StringType,
	"desired_generation":      types.Int64Type,
	"applied_generation":      types.Int64Type,
	"last_error":              types.StringType,
	"ca_cert_filename":        types.StringType,
	"ca_cert_content_base64":  types.StringType,
	"ca_cert_version":         types.Int64Type,
	"ca_cert_checksum_sha256": types.StringType,
}

type pulsarOAuth2Model struct {
	Enabled                 types.Bool   `tfsdk:"enabled"`
	State                   types.String `tfsdk:"state"`
	DesiredGeneration       types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration       types.Int64  `tfsdk:"applied_generation"`
	LastError               types.String `tfsdk:"last_error"`
	IssuerURL               types.String `tfsdk:"issuer_url"`
	Audience                types.String `tfsdk:"audience"`
	PrivateKeyFilename      types.String `tfsdk:"private_key_filename"`
	PrivateKeyContentBase64 types.String `tfsdk:"private_key_content_base64"`
	PrivateKeyVersion       types.Int64  `tfsdk:"private_key_version"`
	PrivateKeyChecksumSHA   types.String `tfsdk:"private_key_checksum_sha256"`
}

var pulsarOAuth2AttributeTypes = map[string]attr.Type{
	"enabled":                     types.BoolType,
	"state":                       types.StringType,
	"desired_generation":          types.Int64Type,
	"applied_generation":          types.Int64Type,
	"last_error":                  types.StringType,
	"issuer_url":                  types.StringType,
	"audience":                    types.StringType,
	"private_key_filename":        types.StringType,
	"private_key_content_base64":  types.StringType,
	"private_key_version":         types.Int64Type,
	"private_key_checksum_sha256": types.StringType,
}

type pulsarTokenModel struct {
	Enabled            types.Bool   `tfsdk:"enabled"`
	State              types.String `tfsdk:"state"`
	DesiredGeneration  types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration  types.Int64  `tfsdk:"applied_generation"`
	LastError          types.String `tfsdk:"last_error"`
	AccessToken        types.String `tfsdk:"access_token"`
	AccessTokenVersion types.Int64  `tfsdk:"access_token_version"`
}

var pulsarTokenAttributeTypes = map[string]attr.Type{
	"enabled":              types.BoolType,
	"state":                types.StringType,
	"desired_generation":   types.Int64Type,
	"applied_generation":   types.Int64Type,
	"last_error":           types.StringType,
	"access_token":         types.StringType,
	"access_token_version": types.Int64Type,
}

func expandPulsarTLS(ctx context.Context, object types.Object) (pulsarTLSModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return pulsarTLSModel{
			Enabled:             types.BoolNull(),
			State:               types.StringNull(),
			DesiredGeneration:   types.Int64Null(),
			AppliedGeneration:   types.Int64Null(),
			LastError:           types.StringNull(),
			CACertFilename:      types.StringNull(),
			CACertContentBase64: types.StringNull(),
			CACertVersion:       types.Int64Null(),
			CACertChecksumSHA:   types.StringNull(),
		}, nil
	}
	var model pulsarTLSModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func expandPulsarOAuth2(ctx context.Context, object types.Object) (pulsarOAuth2Model, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return pulsarOAuth2Model{
			Enabled:                 types.BoolNull(),
			State:                   types.StringNull(),
			DesiredGeneration:       types.Int64Null(),
			AppliedGeneration:       types.Int64Null(),
			LastError:               types.StringNull(),
			IssuerURL:               types.StringNull(),
			Audience:                types.StringNull(),
			PrivateKeyFilename:      types.StringNull(),
			PrivateKeyContentBase64: types.StringNull(),
			PrivateKeyVersion:       types.Int64Null(),
			PrivateKeyChecksumSHA:   types.StringNull(),
		}, nil
	}
	var model pulsarOAuth2Model
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func expandPulsarToken(ctx context.Context, object types.Object) (pulsarTokenModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return pulsarTokenModel{
			Enabled:            types.BoolNull(),
			State:              types.StringNull(),
			DesiredGeneration:  types.Int64Null(),
			AppliedGeneration:  types.Int64Null(),
			LastError:          types.StringNull(),
			AccessToken:        types.StringNull(),
			AccessTokenVersion: types.Int64Null(),
		}, nil
	}
	var model pulsarTokenModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func enrichPulsarTLSFromConfig(ctx context.Context, config tfsdk.Config, model pulsarTLSModel) (pulsarTLSModel, diag.Diagnostics) {
	if model.CACertContentBase64.IsNull() || model.CACertContentBase64.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root("tls").AtName("ca_cert_content_base64"), &value)
		if diags.HasError() {
			return model, diags
		}
		model.CACertContentBase64 = value
	}
	return model, nil
}

func enrichPulsarOAuth2FromConfig(ctx context.Context, config tfsdk.Config, model pulsarOAuth2Model) (pulsarOAuth2Model, diag.Diagnostics) {
	if model.PrivateKeyContentBase64.IsNull() || model.PrivateKeyContentBase64.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root("oauth2").AtName("private_key_content_base64"), &value)
		if diags.HasError() {
			return model, diags
		}
		model.PrivateKeyContentBase64 = value
	}
	return model, nil
}

func enrichPulsarTokenFromConfig(ctx context.Context, config tfsdk.Config, model pulsarTokenModel) (pulsarTokenModel, diag.Diagnostics) {
	if model.AccessToken.IsNull() || model.AccessToken.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root("token").AtName("access_token"), &value)
		if diags.HasError() {
			return model, diags
		}
		model.AccessToken = value
	}
	return model, nil
}

func buildPulsarFeatureBase(feature *client.OutletFeatureDetail) map[string]attr.Value {
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

func buildPulsarTLSObject(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandPulsarTLS(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(pulsarTLSAttributeTypes), diags
	}
	credential := findMaterialOrNil(credentials, pulsarTLSCARole)
	if previous.IsNull() && feature == nil && credential == nil {
		return types.ObjectNull(pulsarTLSAttributeTypes), nil
	}
	values := buildPulsarFeatureBase(feature)
	filename := previousModel.CACertFilename
	checksum := previousModel.CACertChecksumSHA
	if credential != nil {
		filename = nullableString(credential.Filename)
		checksum = nullableString(credential.ChecksumSHA256)
	}
	values["ca_cert_filename"] = filename
	values["ca_cert_content_base64"] = types.StringNull()
	values["ca_cert_version"] = previousModel.CACertVersion
	values["ca_cert_checksum_sha256"] = checksum
	return types.ObjectValue(pulsarTLSAttributeTypes, values)
}

func buildPulsarOAuth2Object(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandPulsarOAuth2(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(pulsarOAuth2AttributeTypes), diags
	}
	credential := findMaterialOrNil(credentials, pulsarOAuth2PrivateKeyRole)
	if previous.IsNull() && feature == nil && credential == nil {
		return types.ObjectNull(pulsarOAuth2AttributeTypes), nil
	}
	values := buildPulsarFeatureBase(feature)
	values["issuer_url"] = previousModel.IssuerURL
	values["audience"] = previousModel.Audience
	filename := previousModel.PrivateKeyFilename
	checksum := previousModel.PrivateKeyChecksumSHA
	if credential != nil {
		filename = nullableString(credential.Filename)
		checksum = nullableString(credential.ChecksumSHA256)
	}
	values["private_key_filename"] = filename
	values["private_key_content_base64"] = types.StringNull()
	values["private_key_version"] = previousModel.PrivateKeyVersion
	values["private_key_checksum_sha256"] = checksum
	return types.ObjectValue(pulsarOAuth2AttributeTypes, values)
}

func buildPulsarTokenObject(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandPulsarToken(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(pulsarTokenAttributeTypes), diags
	}
	credential := findMaterialOrNil(credentials, pulsarTokenRole)
	if previous.IsNull() && feature == nil && credential == nil {
		return types.ObjectNull(pulsarTokenAttributeTypes), nil
	}
	values := buildPulsarFeatureBase(feature)
	if feature == nil && (values["enabled"].IsNull() || values["enabled"].IsUnknown()) {
		values["enabled"] = types.BoolValue(true)
	}
	values["access_token"] = types.StringNull()
	values["access_token_version"] = previousModel.AccessTokenVersion
	if credential == nil && previous.IsNull() {
		return types.ObjectNull(pulsarTokenAttributeTypes), nil
	}
	return types.ObjectValue(pulsarTokenAttributeTypes, values)
}

func pulsarFeatureEnabled(object types.Object) bool {
	if object.IsNull() || object.IsUnknown() {
		return false
	}
	var enabled types.Bool
	_ = object.As(context.Background(), &struct {
		Enabled *types.Bool `tfsdk:"enabled"`
	}{Enabled: &enabled}, basetypes.ObjectAsOptions{})
	return !enabled.IsNull() && !enabled.IsUnknown() && enabled.ValueBool()
}

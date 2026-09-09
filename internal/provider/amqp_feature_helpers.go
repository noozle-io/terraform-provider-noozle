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
	amqp09URLsCredentialRole = "amqp_0_9_urls"
	amqp1URLsCredentialRole  = "amqp_1_urls"
	amqpCACertRole           = "ca_cert"
	amqpClientCertRole       = "client_cert"
	amqpClientKeyRole        = "client_key"
	amqpSASLPasswordRole     = "sasl_password"
)

type amqpTLSModel struct {
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

type amqpMTLSModel struct {
	Enabled                 types.Bool   `tfsdk:"enabled"`
	State                   types.String `tfsdk:"state"`
	DesiredGeneration       types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration       types.Int64  `tfsdk:"applied_generation"`
	LastError               types.String `tfsdk:"last_error"`
	CACertFilename          types.String `tfsdk:"ca_cert_filename"`
	CACertContentBase64     types.String `tfsdk:"ca_cert_content_base64"`
	CACertVersion           types.Int64  `tfsdk:"ca_cert_version"`
	CACertChecksumSHA       types.String `tfsdk:"ca_cert_checksum_sha256"`
	ClientCertFilename      types.String `tfsdk:"client_cert_filename"`
	ClientCertContentBase64 types.String `tfsdk:"client_cert_content_base64"`
	ClientCertVersion       types.Int64  `tfsdk:"client_cert_version"`
	ClientCertChecksumSHA   types.String `tfsdk:"client_cert_checksum_sha256"`
	ClientKeyFilename       types.String `tfsdk:"client_key_filename"`
	ClientKeyContentBase64  types.String `tfsdk:"client_key_content_base64"`
	ClientKeyVersion        types.Int64  `tfsdk:"client_key_version"`
	ClientKeyChecksumSHA    types.String `tfsdk:"client_key_checksum_sha256"`
}

type amqpSASLPlainModel struct {
	Enabled           types.Bool   `tfsdk:"enabled"`
	State             types.String `tfsdk:"state"`
	DesiredGeneration types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration types.Int64  `tfsdk:"applied_generation"`
	LastError         types.String `tfsdk:"last_error"`
	User              types.String `tfsdk:"user"`
	Password          types.String `tfsdk:"password"`
	PasswordVersion   types.Int64  `tfsdk:"password_version"`
}

type amqpSASLAnonymousModel struct {
	Enabled           types.Bool   `tfsdk:"enabled"`
	State             types.String `tfsdk:"state"`
	DesiredGeneration types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration types.Int64  `tfsdk:"applied_generation"`
	LastError         types.String `tfsdk:"last_error"`
}

var amqpTLSAttributeTypes = map[string]attr.Type{
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

var amqpMTLSAttributeTypes = map[string]attr.Type{
	"enabled":                     types.BoolType,
	"state":                       types.StringType,
	"desired_generation":          types.Int64Type,
	"applied_generation":          types.Int64Type,
	"last_error":                  types.StringType,
	"ca_cert_filename":            types.StringType,
	"ca_cert_content_base64":      types.StringType,
	"ca_cert_version":             types.Int64Type,
	"ca_cert_checksum_sha256":     types.StringType,
	"client_cert_filename":        types.StringType,
	"client_cert_content_base64":  types.StringType,
	"client_cert_version":         types.Int64Type,
	"client_cert_checksum_sha256": types.StringType,
	"client_key_filename":         types.StringType,
	"client_key_content_base64":   types.StringType,
	"client_key_version":          types.Int64Type,
	"client_key_checksum_sha256":  types.StringType,
}

var amqpSASLPlainAttributeTypes = map[string]attr.Type{
	"enabled":            types.BoolType,
	"state":              types.StringType,
	"desired_generation": types.Int64Type,
	"applied_generation": types.Int64Type,
	"last_error":         types.StringType,
	"user":               types.StringType,
	"password":           types.StringType,
	"password_version":   types.Int64Type,
}

var amqpSASLAnonymousAttributeTypes = map[string]attr.Type{
	"enabled":            types.BoolType,
	"state":              types.StringType,
	"desired_generation": types.Int64Type,
	"applied_generation": types.Int64Type,
	"last_error":         types.StringType,
}

func buildAMQPFeatureBase(feature *client.OutletFeatureDetail) map[string]attr.Value {
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

func expandAMQPTLS(ctx context.Context, object types.Object) (amqpTLSModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return amqpTLSModel{
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
	var model amqpTLSModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func expandAMQPMTLS(ctx context.Context, object types.Object) (amqpMTLSModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return amqpMTLSModel{
			Enabled:                 types.BoolNull(),
			State:                   types.StringNull(),
			DesiredGeneration:       types.Int64Null(),
			AppliedGeneration:       types.Int64Null(),
			LastError:               types.StringNull(),
			CACertFilename:          types.StringNull(),
			CACertContentBase64:     types.StringNull(),
			CACertVersion:           types.Int64Null(),
			CACertChecksumSHA:       types.StringNull(),
			ClientCertFilename:      types.StringNull(),
			ClientCertContentBase64: types.StringNull(),
			ClientCertVersion:       types.Int64Null(),
			ClientCertChecksumSHA:   types.StringNull(),
			ClientKeyFilename:       types.StringNull(),
			ClientKeyContentBase64:  types.StringNull(),
			ClientKeyVersion:        types.Int64Null(),
			ClientKeyChecksumSHA:    types.StringNull(),
		}, nil
	}
	var model amqpMTLSModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func expandAMQPSASLPlain(ctx context.Context, object types.Object) (amqpSASLPlainModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return amqpSASLPlainModel{
			Enabled:           types.BoolNull(),
			State:             types.StringNull(),
			DesiredGeneration: types.Int64Null(),
			AppliedGeneration: types.Int64Null(),
			LastError:         types.StringNull(),
			User:              types.StringNull(),
			Password:          types.StringNull(),
			PasswordVersion:   types.Int64Null(),
		}, nil
	}
	var model amqpSASLPlainModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func expandAMQPSASLAnonymous(ctx context.Context, object types.Object) (amqpSASLAnonymousModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return amqpSASLAnonymousModel{
			Enabled:           types.BoolNull(),
			State:             types.StringNull(),
			DesiredGeneration: types.Int64Null(),
			AppliedGeneration: types.Int64Null(),
			LastError:         types.StringNull(),
		}, nil
	}
	var model amqpSASLAnonymousModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func enrichAMQPTLSFromConfig(ctx context.Context, config tfsdk.Config, pathRoot string, model amqpTLSModel) (amqpTLSModel, diag.Diagnostics) {
	if model.CACertContentBase64.IsNull() || model.CACertContentBase64.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root(pathRoot).AtName("ca_cert_content_base64"), &value)
		if diags.HasError() {
			return model, diags
		}
		model.CACertContentBase64 = value
	}
	return model, nil
}

func enrichAMQPMTLSFromConfig(ctx context.Context, config tfsdk.Config, pathRoot string, model amqpMTLSModel) (amqpMTLSModel, diag.Diagnostics) {
	if model.CACertContentBase64.IsNull() || model.CACertContentBase64.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root(pathRoot).AtName("ca_cert_content_base64"), &value)
		if diags.HasError() {
			return model, diags
		}
		model.CACertContentBase64 = value
	}
	if model.ClientCertContentBase64.IsNull() || model.ClientCertContentBase64.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root(pathRoot).AtName("client_cert_content_base64"), &value)
		if diags.HasError() {
			return model, diags
		}
		model.ClientCertContentBase64 = value
	}
	if model.ClientKeyContentBase64.IsNull() || model.ClientKeyContentBase64.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root(pathRoot).AtName("client_key_content_base64"), &value)
		if diags.HasError() {
			return model, diags
		}
		model.ClientKeyContentBase64 = value
	}
	return model, nil
}

func enrichAMQPSASLPlainFromConfig(ctx context.Context, config tfsdk.Config, pathRoot string, model amqpSASLPlainModel) (amqpSASLPlainModel, diag.Diagnostics) {
	if model.Password.IsNull() || model.Password.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root(pathRoot).AtName("password"), &value)
		if diags.HasError() {
			return model, diags
		}
		model.Password = value
	}
	return model, nil
}

func buildAMQPTLSObject(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandAMQPTLS(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(amqpTLSAttributeTypes), diags
	}
	credential := findMaterialOrNil(credentials, amqpCACertRole)
	if previous.IsNull() && (feature == nil || !feature.Enabled) {
		return types.ObjectNull(amqpTLSAttributeTypes), nil
	}
	values := buildAMQPFeatureBase(feature)
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
	return types.ObjectValue(amqpTLSAttributeTypes, values)
}

func buildAMQPMTLSObject(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandAMQPMTLS(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(amqpMTLSAttributeTypes), diags
	}
	caCert := findMaterialOrNil(credentials, amqpCACertRole)
	clientCert := findMaterialOrNil(credentials, amqpClientCertRole)
	clientKey := findMaterialOrNil(credentials, amqpClientKeyRole)
	if previous.IsNull() && (feature == nil || !feature.Enabled) && clientCert == nil && clientKey == nil {
		return types.ObjectNull(amqpMTLSAttributeTypes), nil
	}
	values := buildAMQPFeatureBase(feature)
	values["ca_cert_filename"] = previousModel.CACertFilename
	values["ca_cert_content_base64"] = types.StringNull()
	values["ca_cert_version"] = previousModel.CACertVersion
	values["ca_cert_checksum_sha256"] = previousModel.CACertChecksumSHA
	values["client_cert_filename"] = previousModel.ClientCertFilename
	values["client_cert_content_base64"] = types.StringNull()
	values["client_cert_version"] = previousModel.ClientCertVersion
	values["client_cert_checksum_sha256"] = previousModel.ClientCertChecksumSHA
	values["client_key_filename"] = previousModel.ClientKeyFilename
	values["client_key_content_base64"] = types.StringNull()
	values["client_key_version"] = previousModel.ClientKeyVersion
	values["client_key_checksum_sha256"] = previousModel.ClientKeyChecksumSHA
	if caCert != nil {
		values["ca_cert_filename"] = nullableString(caCert.Filename)
		values["ca_cert_checksum_sha256"] = nullableString(caCert.ChecksumSHA256)
	}
	if clientCert != nil {
		values["client_cert_filename"] = nullableString(clientCert.Filename)
		values["client_cert_checksum_sha256"] = nullableString(clientCert.ChecksumSHA256)
	}
	if clientKey != nil {
		values["client_key_filename"] = nullableString(clientKey.Filename)
		values["client_key_checksum_sha256"] = nullableString(clientKey.ChecksumSHA256)
	}
	return types.ObjectValue(amqpMTLSAttributeTypes, values)
}

func buildAMQPSASLPlainObject(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandAMQPSASLPlain(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(amqpSASLPlainAttributeTypes), diags
	}
	credential := findMaterialOrNil(credentials, amqpSASLPasswordRole)
	if previous.IsNull() && (feature == nil || !feature.Enabled) && credential == nil {
		return types.ObjectNull(amqpSASLPlainAttributeTypes), nil
	}
	values := buildAMQPFeatureBase(feature)
	values["user"] = previousModel.User
	values["password"] = types.StringNull()
	values["password_version"] = previousModel.PasswordVersion
	if feature == nil && credential != nil {
		values["enabled"] = types.BoolValue(true)
	}
	return types.ObjectValue(amqpSASLPlainAttributeTypes, values)
}

func buildAMQPSASLAnonymousObject(feature *client.OutletFeatureDetail) (types.Object, diag.Diagnostics) {
	if feature == nil || !feature.Enabled {
		return types.ObjectNull(amqpSASLAnonymousAttributeTypes), nil
	}
	return types.ObjectValue(amqpSASLAnonymousAttributeTypes, buildAMQPFeatureBase(feature))
}

func buildAMQPTLSDataSourceObject(feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	credential := findMaterialOrNil(credentials, amqpCACertRole)
	attrTypes := map[string]attr.Type{
		"enabled":                 types.BoolType,
		"state":                   types.StringType,
		"desired_generation":      types.Int64Type,
		"applied_generation":      types.Int64Type,
		"last_error":              types.StringType,
		"ca_cert_filename":        types.StringType,
		"ca_cert_checksum_sha256": types.StringType,
	}
	if feature == nil || !feature.Enabled {
		return types.ObjectNull(attrTypes), nil
	}
	values := buildAMQPFeatureBase(feature)
	values["ca_cert_filename"] = nullableString(credentialField(credential, func(m *client.OutletMaterialSummary) string { return m.Filename }))
	values["ca_cert_checksum_sha256"] = nullableString(credentialField(credential, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 }))
	return types.ObjectValue(attrTypes, values)
}

func buildAMQPMTLSDataSourceObject(feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	caCert := findMaterialOrNil(credentials, amqpCACertRole)
	clientCert := findMaterialOrNil(credentials, amqpClientCertRole)
	clientKey := findMaterialOrNil(credentials, amqpClientKeyRole)
	attrTypes := map[string]attr.Type{
		"enabled":                     types.BoolType,
		"state":                       types.StringType,
		"desired_generation":          types.Int64Type,
		"applied_generation":          types.Int64Type,
		"last_error":                  types.StringType,
		"ca_cert_filename":            types.StringType,
		"ca_cert_checksum_sha256":     types.StringType,
		"client_cert_filename":        types.StringType,
		"client_cert_checksum_sha256": types.StringType,
		"client_key_filename":         types.StringType,
		"client_key_checksum_sha256":  types.StringType,
	}
	if (feature == nil || !feature.Enabled) && clientCert == nil && clientKey == nil {
		return types.ObjectNull(attrTypes), nil
	}
	values := buildAMQPFeatureBase(feature)
	values["ca_cert_filename"] = nullableString(credentialField(caCert, func(m *client.OutletMaterialSummary) string { return m.Filename }))
	values["ca_cert_checksum_sha256"] = nullableString(credentialField(caCert, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 }))
	values["client_cert_filename"] = nullableString(credentialField(clientCert, func(m *client.OutletMaterialSummary) string { return m.Filename }))
	values["client_cert_checksum_sha256"] = nullableString(credentialField(clientCert, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 }))
	values["client_key_filename"] = nullableString(credentialField(clientKey, func(m *client.OutletMaterialSummary) string { return m.Filename }))
	values["client_key_checksum_sha256"] = nullableString(credentialField(clientKey, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 }))
	return types.ObjectValue(attrTypes, values)
}

func buildAMQPSASLPlainDataSourceObject(feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	credential := findMaterialOrNil(credentials, amqpSASLPasswordRole)
	attrTypes := map[string]attr.Type{
		"enabled":            types.BoolType,
		"state":              types.StringType,
		"desired_generation": types.Int64Type,
		"applied_generation": types.Int64Type,
		"last_error":         types.StringType,
		"user":               types.StringType,
		"password_version":   types.Int64Type,
	}
	if (feature == nil || !feature.Enabled) && credential == nil {
		return types.ObjectNull(attrTypes), nil
	}
	values := buildAMQPFeatureBase(feature)
	values["user"] = types.StringNull()
	values["password_version"] = types.Int64Null()
	if feature == nil && credential != nil {
		values["enabled"] = types.BoolValue(true)
	}
	return types.ObjectValue(attrTypes, values)
}

func buildAMQPSASLAnonymousDataSourceObject(feature *client.OutletFeatureDetail) (types.Object, diag.Diagnostics) {
	return buildAMQPSASLAnonymousObject(feature)
}

func amqpFeatureEnabled(object types.Object) bool {
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

func credentialField(material *client.OutletMaterialSummary, selector func(*client.OutletMaterialSummary) string) string {
	if material == nil {
		return ""
	}
	return selector(material)
}

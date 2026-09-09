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
	kafkaCACertRole          = "ca_cert"
	kafkaClientCertRole      = "client_cert"
	kafkaClientKeyRole       = "client_key"
	kafkaSASLPasswordRole    = "sasl_password"
	kafkaSASLAccessTokenRole = "kafka_sasl_access_token"
	kafkaAwsSecretRole       = "aws_secret"
	kafkaAwsTokenRole        = "aws_token"
)

type kafkaTLSModel struct {
	Enabled                  types.Bool   `tfsdk:"enabled"`
	State                    types.String `tfsdk:"state"`
	DesiredGeneration        types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration        types.Int64  `tfsdk:"applied_generation"`
	LastError                types.String `tfsdk:"last_error"`
	CACertFilename           types.String `tfsdk:"ca_cert_filename"`
	CACertContentBase64      types.String `tfsdk:"ca_cert_content_base64"`
	CACertVersion            types.Int64  `tfsdk:"ca_cert_version"`
	CACertChecksumSHA256     types.String `tfsdk:"ca_cert_checksum_sha256"`
	ClientCertFilename       types.String `tfsdk:"client_cert_filename"`
	ClientCertContentBase64  types.String `tfsdk:"client_cert_content_base64"`
	ClientCertVersion        types.Int64  `tfsdk:"client_cert_version"`
	ClientCertChecksumSHA256 types.String `tfsdk:"client_cert_checksum_sha256"`
	ClientKeyFilename        types.String `tfsdk:"client_key_filename"`
	ClientKeyContentBase64   types.String `tfsdk:"client_key_content_base64"`
	ClientKeyVersion         types.Int64  `tfsdk:"client_key_version"`
	ClientKeyChecksumSHA256  types.String `tfsdk:"client_key_checksum_sha256"`
}

type kafkaSASLUserPassModel struct {
	Enabled           types.Bool   `tfsdk:"enabled"`
	State             types.String `tfsdk:"state"`
	DesiredGeneration types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration types.Int64  `tfsdk:"applied_generation"`
	LastError         types.String `tfsdk:"last_error"`
	User              types.String `tfsdk:"user"`
	Password          types.String `tfsdk:"password"`
	PasswordVersion   types.Int64  `tfsdk:"password_version"`
}

type kafkaSASLOAuthbearerModel struct {
	Enabled            types.Bool   `tfsdk:"enabled"`
	State              types.String `tfsdk:"state"`
	DesiredGeneration  types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration  types.Int64  `tfsdk:"applied_generation"`
	LastError          types.String `tfsdk:"last_error"`
	AccessToken        types.String `tfsdk:"access_token"`
	AccessTokenVersion types.Int64  `tfsdk:"access_token_version"`
}

type kafkaSASLOAuthbearerCacheModel struct {
	Enabled           types.Bool   `tfsdk:"enabled"`
	State             types.String `tfsdk:"state"`
	DesiredGeneration types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration types.Int64  `tfsdk:"applied_generation"`
	LastError         types.String `tfsdk:"last_error"`
	TokenCache        types.String `tfsdk:"token_cache"`
	TokenKey          types.String `tfsdk:"token_key"`
}

type kafkaSASLAwsMSKIAMModel struct {
	Enabled           types.Bool   `tfsdk:"enabled"`
	State             types.String `tfsdk:"state"`
	DesiredGeneration types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration types.Int64  `tfsdk:"applied_generation"`
	LastError         types.String `tfsdk:"last_error"`
	Region            types.String `tfsdk:"region"`
	Endpoint          types.String `tfsdk:"endpoint"`
	Profile           types.String `tfsdk:"profile"`
	ID                types.String `tfsdk:"id"`
	FromEC2Role       types.String `tfsdk:"from_ec2_role"`
	Role              types.String `tfsdk:"role"`
	RoleExternalID    types.String `tfsdk:"role_external_id"`
	ExpiryWindow      types.String `tfsdk:"expiry_window"`
	AwsSecret         types.String `tfsdk:"aws_secret"`
	AwsSecretVersion  types.Int64  `tfsdk:"aws_secret_version"`
	AwsToken          types.String `tfsdk:"aws_token"`
	AwsTokenVersion   types.Int64  `tfsdk:"aws_token_version"`
}

var kafkaTLSAttributeTypes = map[string]attr.Type{
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

var kafkaSASLUserPassAttributeTypes = map[string]attr.Type{
	"enabled":            types.BoolType,
	"state":              types.StringType,
	"desired_generation": types.Int64Type,
	"applied_generation": types.Int64Type,
	"last_error":         types.StringType,
	"user":               types.StringType,
	"password":           types.StringType,
	"password_version":   types.Int64Type,
}

var kafkaSASLOAuthbearerAttributeTypes = map[string]attr.Type{
	"enabled":              types.BoolType,
	"state":                types.StringType,
	"desired_generation":   types.Int64Type,
	"applied_generation":   types.Int64Type,
	"last_error":           types.StringType,
	"access_token":         types.StringType,
	"access_token_version": types.Int64Type,
}

var kafkaSASLOAuthbearerCacheAttributeTypes = map[string]attr.Type{
	"enabled":            types.BoolType,
	"state":              types.StringType,
	"desired_generation": types.Int64Type,
	"applied_generation": types.Int64Type,
	"last_error":         types.StringType,
	"token_cache":        types.StringType,
	"token_key":          types.StringType,
}

var kafkaSASLAwsMSKIAMAttributeTypes = map[string]attr.Type{
	"enabled":            types.BoolType,
	"state":              types.StringType,
	"desired_generation": types.Int64Type,
	"applied_generation": types.Int64Type,
	"last_error":         types.StringType,
	"region":             types.StringType,
	"endpoint":           types.StringType,
	"profile":            types.StringType,
	"id":                 types.StringType,
	"from_ec2_role":      types.StringType,
	"role":               types.StringType,
	"role_external_id":   types.StringType,
	"expiry_window":      types.StringType,
	"aws_secret":         types.StringType,
	"aws_secret_version": types.Int64Type,
	"aws_token":          types.StringType,
	"aws_token_version":  types.Int64Type,
}

func buildKafkaFeatureBase(feature *client.OutletFeatureDetail, previousConfigured bool) map[string]attr.Value {
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
	} else if previousConfigured {
		enabled = types.BoolValue(true)
	}
	return map[string]attr.Value{
		"enabled":            enabled,
		"state":              state,
		"desired_generation": desiredGeneration,
		"applied_generation": appliedGeneration,
		"last_error":         lastError,
	}
}

func expandKafkaTLS(ctx context.Context, object types.Object) (kafkaTLSModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return kafkaTLSModel{
			Enabled: types.BoolNull(), State: types.StringNull(), DesiredGeneration: types.Int64Null(), AppliedGeneration: types.Int64Null(), LastError: types.StringNull(),
			CACertFilename: types.StringNull(), CACertContentBase64: types.StringNull(), CACertVersion: types.Int64Null(), CACertChecksumSHA256: types.StringNull(),
			ClientCertFilename: types.StringNull(), ClientCertContentBase64: types.StringNull(), ClientCertVersion: types.Int64Null(), ClientCertChecksumSHA256: types.StringNull(),
			ClientKeyFilename: types.StringNull(), ClientKeyContentBase64: types.StringNull(), ClientKeyVersion: types.Int64Null(), ClientKeyChecksumSHA256: types.StringNull(),
		}, nil
	}
	var model kafkaTLSModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func expandKafkaSASLUserPass(ctx context.Context, object types.Object) (kafkaSASLUserPassModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return kafkaSASLUserPassModel{
			Enabled: types.BoolNull(), State: types.StringNull(), DesiredGeneration: types.Int64Null(), AppliedGeneration: types.Int64Null(), LastError: types.StringNull(),
			User: types.StringNull(), Password: types.StringNull(), PasswordVersion: types.Int64Null(),
		}, nil
	}
	var model kafkaSASLUserPassModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func expandKafkaSASLOAuthbearer(ctx context.Context, object types.Object) (kafkaSASLOAuthbearerModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return kafkaSASLOAuthbearerModel{
			Enabled: types.BoolNull(), State: types.StringNull(), DesiredGeneration: types.Int64Null(), AppliedGeneration: types.Int64Null(), LastError: types.StringNull(),
			AccessToken: types.StringNull(), AccessTokenVersion: types.Int64Null(),
		}, nil
	}
	var model kafkaSASLOAuthbearerModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func expandKafkaSASLOAuthbearerCache(ctx context.Context, object types.Object) (kafkaSASLOAuthbearerCacheModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return kafkaSASLOAuthbearerCacheModel{
			Enabled: types.BoolNull(), State: types.StringNull(), DesiredGeneration: types.Int64Null(), AppliedGeneration: types.Int64Null(), LastError: types.StringNull(),
			TokenCache: types.StringNull(), TokenKey: types.StringNull(),
		}, nil
	}
	var model kafkaSASLOAuthbearerCacheModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func expandKafkaSASLAwsMSKIAM(ctx context.Context, object types.Object) (kafkaSASLAwsMSKIAMModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return kafkaSASLAwsMSKIAMModel{
			Enabled: types.BoolNull(), State: types.StringNull(), DesiredGeneration: types.Int64Null(), AppliedGeneration: types.Int64Null(), LastError: types.StringNull(),
			Region: types.StringNull(), Endpoint: types.StringNull(), Profile: types.StringNull(), ID: types.StringNull(), FromEC2Role: types.StringNull(),
			Role: types.StringNull(), RoleExternalID: types.StringNull(), ExpiryWindow: types.StringNull(),
			AwsSecret: types.StringNull(), AwsSecretVersion: types.Int64Null(), AwsToken: types.StringNull(), AwsTokenVersion: types.Int64Null(),
		}, nil
	}
	var model kafkaSASLAwsMSKIAMModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	return model, diags
}

func enrichKafkaTLSFromConfig(ctx context.Context, config tfsdk.Config, pathRoot string, model kafkaTLSModel) (kafkaTLSModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	if model.CACertContentBase64.IsNull() || model.CACertContentBase64.IsUnknown() {
		var value types.String
		diags.Append(config.GetAttribute(ctx, path.Root(pathRoot).AtName("ca_cert_content_base64"), &value)...)
		model.CACertContentBase64 = value
	}
	if model.ClientCertContentBase64.IsNull() || model.ClientCertContentBase64.IsUnknown() {
		var value types.String
		diags.Append(config.GetAttribute(ctx, path.Root(pathRoot).AtName("client_cert_content_base64"), &value)...)
		model.ClientCertContentBase64 = value
	}
	if model.ClientKeyContentBase64.IsNull() || model.ClientKeyContentBase64.IsUnknown() {
		var value types.String
		diags.Append(config.GetAttribute(ctx, path.Root(pathRoot).AtName("client_key_content_base64"), &value)...)
		model.ClientKeyContentBase64 = value
	}
	return model, diags
}

func enrichKafkaSASLUserPassFromConfig(ctx context.Context, config tfsdk.Config, pathRoot string, model kafkaSASLUserPassModel) (kafkaSASLUserPassModel, diag.Diagnostics) {
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

func enrichKafkaSASLOAuthbearerFromConfig(ctx context.Context, config tfsdk.Config, model kafkaSASLOAuthbearerModel) (kafkaSASLOAuthbearerModel, diag.Diagnostics) {
	if model.AccessToken.IsNull() || model.AccessToken.IsUnknown() {
		var value types.String
		diags := config.GetAttribute(ctx, path.Root("sasl_oauthbearer").AtName("access_token"), &value)
		if diags.HasError() {
			return model, diags
		}
		model.AccessToken = value
	}
	return model, nil
}

func enrichKafkaSASLAwsMSKIAMFromConfig(ctx context.Context, config tfsdk.Config, model kafkaSASLAwsMSKIAMModel) (kafkaSASLAwsMSKIAMModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	if model.AwsSecret.IsNull() || model.AwsSecret.IsUnknown() {
		var value types.String
		diags.Append(config.GetAttribute(ctx, path.Root("sasl_aws_msk_iam").AtName("aws_secret"), &value)...)
		model.AwsSecret = value
	}
	if model.AwsToken.IsNull() || model.AwsToken.IsUnknown() {
		var value types.String
		diags.Append(config.GetAttribute(ctx, path.Root("sasl_aws_msk_iam").AtName("aws_token"), &value)...)
		model.AwsToken = value
	}
	return model, diags
}

func buildKafkaTLSObject(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandKafkaTLS(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(kafkaTLSAttributeTypes), diags
	}
	if previous.IsNull() && (feature == nil || !feature.Enabled) {
		return types.ObjectNull(kafkaTLSAttributeTypes), nil
	}
	values := buildKafkaFeatureBase(feature, !previous.IsNull())
	caCert := findMaterialOrNil(credentials, kafkaCACertRole)
	clientCert := findMaterialOrNil(credentials, kafkaClientCertRole)
	clientKey := findMaterialOrNil(credentials, kafkaClientKeyRole)
	values["ca_cert_filename"] = chooseMaterialString(previousModel.CACertFilename, caCert, func(m *client.OutletMaterialSummary) string { return m.Filename })
	values["ca_cert_content_base64"] = types.StringNull()
	values["ca_cert_version"] = previousModel.CACertVersion
	values["ca_cert_checksum_sha256"] = chooseMaterialString(previousModel.CACertChecksumSHA256, caCert, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 })
	values["client_cert_filename"] = chooseMaterialString(previousModel.ClientCertFilename, clientCert, func(m *client.OutletMaterialSummary) string { return m.Filename })
	values["client_cert_content_base64"] = types.StringNull()
	values["client_cert_version"] = previousModel.ClientCertVersion
	values["client_cert_checksum_sha256"] = chooseMaterialString(previousModel.ClientCertChecksumSHA256, clientCert, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 })
	values["client_key_filename"] = chooseMaterialString(previousModel.ClientKeyFilename, clientKey, func(m *client.OutletMaterialSummary) string { return m.Filename })
	values["client_key_content_base64"] = types.StringNull()
	values["client_key_version"] = previousModel.ClientKeyVersion
	values["client_key_checksum_sha256"] = chooseMaterialString(previousModel.ClientKeyChecksumSHA256, clientKey, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 })
	return types.ObjectValue(kafkaTLSAttributeTypes, values)
}

func buildKafkaSASLUserPassObject(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandKafkaSASLUserPass(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(kafkaSASLUserPassAttributeTypes), diags
	}
	if previous.IsNull() && (feature == nil || !feature.Enabled) {
		return types.ObjectNull(kafkaSASLUserPassAttributeTypes), nil
	}
	values := buildKafkaFeatureBase(feature, !previous.IsNull())
	values["user"] = previousModel.User
	values["password"] = types.StringNull()
	values["password_version"] = previousModel.PasswordVersion
	return types.ObjectValue(kafkaSASLUserPassAttributeTypes, values)
}

func buildKafkaSASLOAuthbearerObject(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandKafkaSASLOAuthbearer(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(kafkaSASLOAuthbearerAttributeTypes), diags
	}
	if previous.IsNull() && (feature == nil || !feature.Enabled) {
		return types.ObjectNull(kafkaSASLOAuthbearerAttributeTypes), nil
	}
	values := buildKafkaFeatureBase(feature, !previous.IsNull())
	values["access_token"] = types.StringNull()
	values["access_token_version"] = previousModel.AccessTokenVersion
	return types.ObjectValue(kafkaSASLOAuthbearerAttributeTypes, values)
}

func buildKafkaSASLOAuthbearerCacheObject(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandKafkaSASLOAuthbearerCache(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(kafkaSASLOAuthbearerCacheAttributeTypes), diags
	}
	if previous.IsNull() && (feature == nil || !feature.Enabled) {
		return types.ObjectNull(kafkaSASLOAuthbearerCacheAttributeTypes), nil
	}
	values := buildKafkaFeatureBase(feature, !previous.IsNull())
	values["token_cache"] = previousModel.TokenCache
	values["token_key"] = previousModel.TokenKey
	return types.ObjectValue(kafkaSASLOAuthbearerCacheAttributeTypes, values)
}

func buildKafkaSASLAwsMSKIAMObject(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandKafkaSASLAwsMSKIAM(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(kafkaSASLAwsMSKIAMAttributeTypes), diags
	}
	if previous.IsNull() && (feature == nil || !feature.Enabled) {
		return types.ObjectNull(kafkaSASLAwsMSKIAMAttributeTypes), nil
	}
	values := buildKafkaFeatureBase(feature, !previous.IsNull())
	values["region"] = previousModel.Region
	values["endpoint"] = previousModel.Endpoint
	values["profile"] = previousModel.Profile
	values["id"] = previousModel.ID
	values["from_ec2_role"] = previousModel.FromEC2Role
	values["role"] = previousModel.Role
	values["role_external_id"] = previousModel.RoleExternalID
	values["expiry_window"] = previousModel.ExpiryWindow
	values["aws_secret"] = types.StringNull()
	values["aws_secret_version"] = previousModel.AwsSecretVersion
	values["aws_token"] = types.StringNull()
	values["aws_token_version"] = previousModel.AwsTokenVersion
	return types.ObjectValue(kafkaSASLAwsMSKIAMAttributeTypes, values)
}

func buildKafkaTLSDataSourceObject(feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"enabled": types.BoolType, "state": types.StringType, "desired_generation": types.Int64Type, "applied_generation": types.Int64Type, "last_error": types.StringType,
		"ca_cert_filename": types.StringType, "ca_cert_checksum_sha256": types.StringType,
		"client_cert_filename": types.StringType, "client_cert_checksum_sha256": types.StringType,
		"client_key_filename": types.StringType, "client_key_checksum_sha256": types.StringType,
	}
	if feature == nil || !feature.Enabled {
		return types.ObjectNull(attrTypes), nil
	}
	values := buildKafkaFeatureBase(feature, false)
	caCert := findMaterialOrNil(credentials, kafkaCACertRole)
	clientCert := findMaterialOrNil(credentials, kafkaClientCertRole)
	clientKey := findMaterialOrNil(credentials, kafkaClientKeyRole)
	values["ca_cert_filename"] = nullableString(credentialField(caCert, func(m *client.OutletMaterialSummary) string { return m.Filename }))
	values["ca_cert_checksum_sha256"] = nullableString(credentialField(caCert, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 }))
	values["client_cert_filename"] = nullableString(credentialField(clientCert, func(m *client.OutletMaterialSummary) string { return m.Filename }))
	values["client_cert_checksum_sha256"] = nullableString(credentialField(clientCert, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 }))
	values["client_key_filename"] = nullableString(credentialField(clientKey, func(m *client.OutletMaterialSummary) string { return m.Filename }))
	values["client_key_checksum_sha256"] = nullableString(credentialField(clientKey, func(m *client.OutletMaterialSummary) string { return m.ChecksumSHA256 }))
	return types.ObjectValue(attrTypes, values)
}

func buildKafkaFeatureDataSourceObject(feature *client.OutletFeatureDetail) (types.Object, diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"enabled": types.BoolType, "state": types.StringType, "desired_generation": types.Int64Type, "applied_generation": types.Int64Type, "last_error": types.StringType,
	}
	if feature == nil || !feature.Enabled {
		return types.ObjectNull(attrTypes), nil
	}
	return types.ObjectValue(attrTypes, buildKafkaFeatureBase(feature, false))
}

func chooseMaterialString(previous types.String, material *client.OutletMaterialSummary, selector func(*client.OutletMaterialSummary) string) types.String {
	if material != nil {
		return nullableString(selector(material))
	}
	return previous
}

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"terraform-provider-noozle/internal/client"
)

const (
	httpServerTLSCertRole = "server_cert"
	httpServerTLSKeyRole  = "server_key"
)

type httpServerTLSModel struct {
	Enabled                 types.Bool   `tfsdk:"enabled"`
	State                   types.String `tfsdk:"state"`
	DesiredGeneration       types.Int64  `tfsdk:"desired_generation"`
	AppliedGeneration       types.Int64  `tfsdk:"applied_generation"`
	LastError               types.String `tfsdk:"last_error"`
	ServerCertFilename      types.String `tfsdk:"server_cert_filename"`
	ServerCertContentBase64 types.String `tfsdk:"server_cert_content_base64"`
	ServerCertVersion       types.Int64  `tfsdk:"server_cert_version"`
	ServerCertChecksumSHA   types.String `tfsdk:"server_cert_checksum_sha256"`
	ServerKeyFilename       types.String `tfsdk:"server_key_filename"`
	ServerKeyContentBase64  types.String `tfsdk:"server_key_content_base64"`
	ServerKeyVersion        types.Int64  `tfsdk:"server_key_version"`
	ServerKeyChecksumSHA    types.String `tfsdk:"server_key_checksum_sha256"`
}

var httpServerTLSAttributeTypes = map[string]attr.Type{
	"enabled":                     types.BoolType,
	"state":                       types.StringType,
	"desired_generation":          types.Int64Type,
	"applied_generation":          types.Int64Type,
	"last_error":                  types.StringType,
	"server_cert_filename":        types.StringType,
	"server_cert_content_base64":  types.StringType,
	"server_cert_version":         types.Int64Type,
	"server_cert_checksum_sha256": types.StringType,
	"server_key_filename":         types.StringType,
	"server_key_content_base64":   types.StringType,
	"server_key_version":          types.Int64Type,
	"server_key_checksum_sha256":  types.StringType,
}

func expandHttpServerTLS(ctx context.Context, object types.Object) (httpServerTLSModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return httpServerTLSModel{
			Enabled:                 types.BoolNull(),
			State:                   types.StringNull(),
			DesiredGeneration:       types.Int64Null(),
			AppliedGeneration:       types.Int64Null(),
			LastError:               types.StringNull(),
			ServerCertFilename:      types.StringNull(),
			ServerCertContentBase64: types.StringNull(),
			ServerCertVersion:       types.Int64Null(),
			ServerCertChecksumSHA:   types.StringNull(),
			ServerKeyFilename:       types.StringNull(),
			ServerKeyContentBase64:  types.StringNull(),
			ServerKeyVersion:        types.Int64Null(),
			ServerKeyChecksumSHA:    types.StringNull(),
		}, nil
	}

	var tls httpServerTLSModel
	diags := object.As(ctx, &tls, basetypes.ObjectAsOptions{})
	return tls, diags
}

func buildHttpServerTLSObject(ctx context.Context, previous types.Object, feature *client.OutletFeatureDetail, credentials []client.OutletMaterialSummary) (types.Object, diag.Diagnostics) {
	previousModel, diags := expandHttpServerTLS(ctx, previous)
	if diags.HasError() {
		return types.ObjectNull(httpServerTLSAttributeTypes), diags
	}

	serverCert := findMaterialOrNil(credentials, httpServerTLSCertRole)
	serverKey := findMaterialOrNil(credentials, httpServerTLSKeyRole)
	if previous.IsNull() && feature == nil && serverCert == nil && serverKey == nil {
		return types.ObjectNull(httpServerTLSAttributeTypes), nil
	}

	enabled := previousModel.Enabled
	state := previousModel.State
	desiredGeneration := previousModel.DesiredGeneration
	appliedGeneration := previousModel.AppliedGeneration
	lastError := previousModel.LastError
	if feature != nil {
		enabled = types.BoolValue(feature.Enabled)
		state = nullableString(feature.State)
		desiredGeneration = nullableInt64Pointer(feature.DesiredGeneration)
		appliedGeneration = nullableInt64Pointer(feature.AppliedGeneration)
		lastError = nullableString(feature.LastError)
	} else {
		if enabled.IsNull() || enabled.IsUnknown() {
			enabled = types.BoolValue(true)
		}
		if state.IsNull() || state.IsUnknown() {
			state = types.StringValue("unknown")
		}
	}

	serverCertFilename := previousModel.ServerCertFilename
	serverCertChecksum := previousModel.ServerCertChecksumSHA
	if serverCert != nil {
		serverCertFilename = nullableString(serverCert.Filename)
		serverCertChecksum = nullableString(serverCert.ChecksumSHA256)
	}

	serverKeyFilename := previousModel.ServerKeyFilename
	serverKeyChecksum := previousModel.ServerKeyChecksumSHA
	if serverKey != nil {
		serverKeyFilename = nullableString(serverKey.Filename)
		serverKeyChecksum = nullableString(serverKey.ChecksumSHA256)
	}

	return types.ObjectValue(httpServerTLSAttributeTypes, map[string]attr.Value{
		"enabled":                     enabled,
		"state":                       state,
		"desired_generation":          desiredGeneration,
		"applied_generation":          appliedGeneration,
		"last_error":                  lastError,
		"server_cert_filename":        serverCertFilename,
		"server_cert_content_base64":  types.StringNull(),
		"server_cert_version":         previousModel.ServerCertVersion,
		"server_cert_checksum_sha256": serverCertChecksum,
		"server_key_filename":         serverKeyFilename,
		"server_key_content_base64":   types.StringNull(),
		"server_key_version":          previousModel.ServerKeyVersion,
		"server_key_checksum_sha256":  serverKeyChecksum,
	})
}

func httpServerTLSChanged(plan, state httpServerTLSModel) bool {
	return !plan.ServerCertFilename.Equal(state.ServerCertFilename) ||
		!plan.ServerCertVersion.Equal(state.ServerCertVersion) ||
		!plan.ServerKeyFilename.Equal(state.ServerKeyFilename) ||
		!plan.ServerKeyVersion.Equal(state.ServerKeyVersion)
}

func httpServerTLSEnabled(model httpServerTLSModel) bool {
	return !model.Enabled.IsNull() && !model.Enabled.IsUnknown() && model.Enabled.ValueBool()
}

package provider

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"terraform-provider-noozle/internal/client"
)

func addMissingTenantContextDiagnosticForError(err error, diags *diag.Diagnostics) bool {
	if !errors.Is(err, client.ErrMissingTenantContext) {
		return false
	}

	diags.AddError(
		"Missing Tenant Context",
		"Set provider tenant_id or NOOZLE_TENANT_ID before managing Noozle queries, outlets, or outlet-query links.",
	)
	return true
}

package provider

import "terraform-provider-noozle/internal/client"

const gcpPubsubCredentialsRole = "gcp_credentials_file"

func findMaterialOrNil(materials []client.OutletMaterialSummary, role string) *client.OutletMaterialSummary {
	for _, material := range materials {
		if material.Role == role {
			materialCopy := material
			return &materialCopy
		}
	}

	return nil
}

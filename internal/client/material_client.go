package client

import (
	"context"
)

type OutletMaterialSummary struct {
	Role           string `json:"role"`
	Type           string `json:"type"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	Filename       string `json:"filename"`
	ContentType    string `json:"content_type"`
	ChecksumSHA256 string `json:"checksum_sha256"`
}

type outletMaterialsResponse struct {
	Materials   []OutletMaterialSummary `json:"materials"`
	Credentials []OutletMaterialSummary `json:"credentials"`
}

type OutletMaterialDeleteResponse struct {
	Role               string `json:"role"`
	Deleted            bool   `json:"deleted"`
	ScheduledReconcile bool   `json:"scheduled_reconcile"`
}

func (c *Client) listOutletMaterials(ctx context.Context, path string) ([]OutletMaterialSummary, error) {
	var response outletMaterialsResponse
	if err := c.doJSON(ctx, httpMethodGet, path, nil, &response, statusOK); err != nil {
		return nil, err
	}

	if len(response.Credentials) > 0 {
		return response.Credentials, nil
	}

	return response.Materials, nil
}

func (c *Client) deleteOutletMaterial(ctx context.Context, path string) (OutletMaterialDeleteResponse, error) {
	var response OutletMaterialDeleteResponse
	if err := c.doJSON(ctx, httpMethodDelete, path, nil, &response, statusOK); err != nil {
		return OutletMaterialDeleteResponse{}, err
	}

	return response, nil
}

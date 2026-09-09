package client

import (
	"context"
	"fmt"
)

func (c *Client) ListAwsSnsCredentials(ctx context.Context, outletID int64) ([]OutletMaterialSummary, error) {
	return c.listOutletMaterials(ctx, fmt.Sprintf("/v1/terraform/aws_sns/%d/credentials", outletID))
}

func (c *Client) DeleteAwsSnsMaterial(ctx context.Context, outletID int64, role string) (OutletMaterialDeleteResponse, error) {
	return c.deleteOutletMaterial(ctx, fmt.Sprintf("/v1/terraform/aws_sns/%d/credentials/%s", outletID, role))
}

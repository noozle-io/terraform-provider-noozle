package client

import (
	"context"
	"fmt"
)

func (c *Client) ListAwsSqsCredentials(ctx context.Context, outletID int64) ([]OutletMaterialSummary, error) {
	return c.listOutletMaterials(ctx, fmt.Sprintf("/v1/terraform/aws_sqs/%d/credentials", outletID))
}

func (c *Client) DeleteAwsSqsMaterial(ctx context.Context, outletID int64, role string) (OutletMaterialDeleteResponse, error) {
	return c.deleteOutletMaterial(ctx, fmt.Sprintf("/v1/terraform/aws_sqs/%d/credentials/%s", outletID, role))
}

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

// typedOutletTransformClient hides provider-specific transform routes behind a
// small capability shared by transform-supporting outlet resources.
type typedOutletTransformClient interface {
	GetTypedOutletTransform(context.Context, string, int64) (client.OutletTransformState, error)
	SetTypedOutletTransform(context.Context, string, int64, string) (client.OutletTransformState, error)
	ClearTypedOutletTransform(context.Context, string, int64) (client.OutletTransformState, error)
}

func syncTypedOutletTransform(ctx context.Context, transformClient typedOutletTransformClient, outletType string, outletID int64, plan, state types.String) error {
	if plan.IsUnknown() {
		return nil
	}
	if plan.Equal(state) {
		return nil
	}
	if plan.IsNull() {
		_, err := transformClient.ClearTypedOutletTransform(ctx, outletType, outletID)
		return err
	}
	_, err := transformClient.SetTypedOutletTransform(ctx, outletType, outletID, plan.ValueString())
	return err
}

var outletTransformAttributeTypes = map[string]attr.Type{
	"selected":   types.StringType,
	"effective":  types.StringType,
	"state":      types.StringType,
	"applied":    types.BoolType,
	"generation": types.Int64Type,
	"last_error": types.StringType,
}

func transformObject(transform client.OutletTransformState) types.Object {
	return types.ObjectValueMust(outletTransformAttributeTypes, map[string]attr.Value{
		"selected":   nullableStringPointerPreservingEmpty(transform.SelectedTransform),
		"effective":  nullableStringPointerPreservingEmpty(transform.EffectiveTransform),
		"state":      nullableStringPointerPreservingEmpty(transform.State),
		"applied":    types.BoolValue(transform.Applied),
		"generation": nullableInt64Pointer(transform.Generation),
		"last_error": nullableStringPointerPreservingEmpty(transform.LastError),
	})
}

func selectedTransform(value types.Object) types.String {
	if value.IsNull() {
		return types.StringNull()
	}
	if value.IsUnknown() {
		return types.StringUnknown()
	}
	selected, ok := value.Attributes()["selected"].(types.String)
	if !ok {
		return types.StringNull()
	}
	return selected
}

// selectedTransformFromConfig reads the configured selection rather than the
// Optional+Computed plan value. This lets removing the block clear a remotely
// selected transform instead of retaining the previous runtime state.
func selectedTransformFromConfig(ctx context.Context, config tfsdk.Config) (types.String, diag.Diagnostics) {
	var transform types.Object
	diags := config.GetAttribute(ctx, path.Root("transform"), &transform)
	if diags.HasError() {
		return types.StringNull(), diags
	}
	return selectedTransform(transform), diags
}

package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type outletFieldUpdater interface {
	SetField(context.Context, int64, string, any) error
	DeleteField(context.Context, int64, string) error
}

func stringValueOrEmpty(value types.String) string {
	if value.IsNull() || value.IsUnknown() {
		return ""
	}

	return value.ValueString()
}

func stringPointerValue(value types.String) *string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	stringValue := value.ValueString()
	return &stringValue
}

func int64PointerValue(value types.Int64) *int64 {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	intValue := value.ValueInt64()
	return &intValue
}

func boolPointerValue(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	boolValue := value.ValueBool()
	return &boolValue
}

func float64PointerValue(value types.Float64) *float64 {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	floatValue := value.ValueFloat64()
	return &floatValue
}

func nullableString(value string) types.String {
	if value == "" {
		return types.StringNull()
	}

	return types.StringValue(value)
}

func nullableStringPointer(value *string) types.String {
	if value == nil || *value == "" {
		return types.StringNull()
	}

	return types.StringValue(*value)
}

func nullableStringPointerPreservingEmpty(value *string) types.String {
	if value == nil {
		return types.StringNull()
	}

	return types.StringValue(*value)
}

func nullableInt64Pointer(value *int64) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}

	return types.Int64Value(*value)
}

func nullableBoolPointer(value *bool) types.Bool {
	if value == nil {
		return types.BoolNull()
	}

	return types.BoolValue(*value)
}

func nullableFloat64Pointer(value *float64) types.Float64 {
	if value == nil {
		return types.Float64Null()
	}

	return types.Float64Value(*value)
}

func nullableStringList(values []string) types.List {
	if values == nil {
		return types.ListNull(types.StringType)
	}

	elements := make([]attrValue, 0, len(values))
	for _, value := range values {
		elements = append(elements, types.StringValue(value))
	}

	return types.ListValueMust(types.StringType, toTypeValues(elements))
}

func nullableStringMap(values map[string]string) types.Map {
	if values == nil {
		return types.MapNull(types.StringType)
	}

	elements := make(map[string]attr.Value, len(values))
	for key, value := range values {
		elements[key] = types.StringValue(value)
	}

	return types.MapValueMust(types.StringType, elements)
}

type attrValue = types.String

func toTypeValues(values []attrValue) []attr.Value {
	result := make([]attr.Value, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}

	return result
}

func stringListValueOrNil(ctx context.Context, value types.List) ([]string, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}

	values := make([]string, 0, len(value.Elements()))
	if diags := value.ElementsAs(ctx, &values, false); diags.HasError() {
		return nil, errors.New(diags.Errors()[0].Detail())
	}

	return values, nil
}

func stringMapValueOrNil(ctx context.Context, value types.Map) (map[string]string, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}

	values := make(map[string]string, len(value.Elements()))
	if diags := value.ElementsAs(ctx, &values, false); diags.HasError() {
		return nil, errors.New(diags.Errors()[0].Detail())
	}

	return values, nil
}

func mustStringListOrEmpty(ctx context.Context, value types.List) []string {
	values, err := stringListValueOrNil(ctx, value)
	if err != nil || values == nil {
		return []string{}
	}

	return values
}

func syncStringField(ctx context.Context, client outletFieldUpdater, outletID int64, field string, plan, state types.String) error {
	if plan.Equal(state) {
		return nil
	}
	if plan.IsNull() {
		return client.DeleteField(ctx, outletID, field)
	}

	return client.SetField(ctx, outletID, field, plan.ValueString())
}

func syncInt64Field(ctx context.Context, client outletFieldUpdater, outletID int64, field string, plan, state types.Int64) error {
	if plan.Equal(state) {
		return nil
	}
	if plan.IsNull() {
		return client.DeleteField(ctx, outletID, field)
	}

	return client.SetField(ctx, outletID, field, plan.ValueInt64())
}

func syncBoolField(ctx context.Context, client outletFieldUpdater, outletID int64, field string, plan, state types.Bool) error {
	if plan.Equal(state) {
		return nil
	}
	if plan.IsNull() {
		return client.DeleteField(ctx, outletID, field)
	}

	return client.SetField(ctx, outletID, field, plan.ValueBool())
}

func syncFloat64Field(ctx context.Context, client outletFieldUpdater, outletID int64, field string, plan, state types.Float64) error {
	if plan.Equal(state) {
		return nil
	}
	if plan.IsNull() {
		return client.DeleteField(ctx, outletID, field)
	}

	return client.SetField(ctx, outletID, field, plan.ValueFloat64())
}

func syncStringListField(ctx context.Context, client outletFieldUpdater, outletID int64, field string, plan, state types.List) error {
	if plan.Equal(state) {
		return nil
	}
	if plan.IsNull() {
		return client.DeleteField(ctx, outletID, field)
	}

	values, err := stringListValueOrNil(ctx, plan)
	if err != nil {
		return err
	}

	return client.SetField(ctx, outletID, field, values)
}

func syncStringMapField(ctx context.Context, client outletFieldUpdater, outletID int64, field string, plan, state types.Map) error {
	if plan.Equal(state) {
		return nil
	}
	if plan.IsNull() {
		return client.DeleteField(ctx, outletID, field)
	}

	values, err := stringMapValueOrNil(ctx, plan)
	if err != nil {
		return err
	}

	return client.SetField(ctx, outletID, field, values)
}

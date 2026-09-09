package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type outletLifecycleTestModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type outletLifecycleTestAdapter struct {
	createResult outletLifecycleResult[outletLifecycleTestModel]
	createDiags  diag.Diagnostics
	readResult   outletLifecycleResult[outletLifecycleTestModel]
	updateDiags  diag.Diagnostics
}

func (a outletLifecycleTestAdapter) Create(context.Context, tfsdk.Config, outletLifecycleTestModel) (outletLifecycleResult[outletLifecycleTestModel], diag.Diagnostics) {
	return a.createResult, a.createDiags
}

func (a outletLifecycleTestAdapter) Read(context.Context, outletLifecycleTestModel) (outletLifecycleResult[outletLifecycleTestModel], diag.Diagnostics) {
	return a.readResult, nil
}

func (a outletLifecycleTestAdapter) Update(context.Context, tfsdk.Config, outletLifecycleTestModel, outletLifecycleTestModel) (outletLifecycleResult[outletLifecycleTestModel], diag.Diagnostics) {
	return outletLifecycleResult[outletLifecycleTestModel]{}, a.updateDiags
}

func (outletLifecycleTestAdapter) Delete(context.Context, outletLifecycleTestModel) diag.Diagnostics {
	return nil
}

func TestOutletLifecycleCreatesStateThroughAdapter(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schema := outletLifecycleTestSchema()
	adapter := outletLifecycleTestAdapter{createResult: outletLifecycleResult[outletLifecycleTestModel]{
		State: outletLifecycleTestModel{ID: types.Int64Value(42), Name: types.StringValue("created")},
	}}
	lifecycle := newOutletLifecycle("test Outlet", outletLifecycleTestAdapterFactory(adapter))
	configureOutletLifecycle(t, ctx, &lifecycle, adapter)

	plan := tfsdk.Plan{Schema: schema}
	planDiags := plan.Set(ctx, &outletLifecycleTestModel{ID: types.Int64Unknown(), Name: types.StringValue("planned")})
	if planDiags.HasError() {
		t.Fatalf("set plan: %#v", planDiags)
	}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: schema}}

	lifecycle.Create(ctx, resource.CreateRequest{Plan: plan}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("create diagnostics: %#v", resp.Diagnostics)
	}
	var got outletLifecycleTestModel
	stateDiags := resp.State.Get(ctx, &got)
	if stateDiags.HasError() {
		t.Fatalf("read state: %#v", stateDiags)
	}
	if got.ID.ValueInt64() != 42 || got.Name.ValueString() != "created" {
		t.Fatalf("state = %#v, want adapter state", got)
	}
}

func TestOutletLifecycleRemovesStateWhenAdapterObservesMissingOutlet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schema := outletLifecycleTestSchema()
	adapter := outletLifecycleTestAdapter{readResult: outletLifecycleResult[outletLifecycleTestModel]{Removed: true}}
	lifecycle := newOutletLifecycle("test Outlet", outletLifecycleTestAdapterFactory(adapter))
	configureOutletLifecycle(t, ctx, &lifecycle, adapter)

	state := tfsdk.State{Schema: schema}
	stateDiags := state.Set(ctx, &outletLifecycleTestModel{ID: types.Int64Value(42), Name: types.StringValue("present")})
	if stateDiags.HasError() {
		t.Fatalf("set state: %#v", stateDiags)
	}
	resp := &resource.ReadResponse{State: tfsdk.State{Schema: schema}}

	lifecycle.Read(ctx, resource.ReadRequest{State: state}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("read diagnostics: %#v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Fatalf("state = %#v, want removed", resp.State.Raw)
	}
}

func TestOutletLifecycleDoesNotWriteStateAfterAdapterFailure(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schema := outletLifecycleTestSchema()
	adapter := outletLifecycleTestAdapter{updateDiags: diag.Diagnostics{diag.NewErrorDiagnostic("Update failed", "partial remote change")}}
	lifecycle := newOutletLifecycle("test Outlet", outletLifecycleTestAdapterFactory(adapter))
	configureOutletLifecycle(t, ctx, &lifecycle, adapter)

	plan := tfsdk.Plan{Schema: schema}
	planDiags := plan.Set(ctx, &outletLifecycleTestModel{ID: types.Int64Value(42), Name: types.StringValue("new")})
	if planDiags.HasError() {
		t.Fatalf("set plan: %#v", planDiags)
	}
	state := tfsdk.State{Schema: schema}
	stateDiags := state.Set(ctx, &outletLifecycleTestModel{ID: types.Int64Value(42), Name: types.StringValue("old")})
	if stateDiags.HasError() {
		t.Fatalf("set state: %#v", stateDiags)
	}
	resp := &resource.UpdateResponse{State: state}

	lifecycle.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, resp)

	if !resp.Diagnostics.HasError() || resp.Diagnostics.Errors()[0].Summary() != "Update failed" {
		t.Fatalf("diagnostics = %#v, want adapter failure", resp.Diagnostics)
	}
	var got outletLifecycleTestModel
	resultDiags := resp.State.Get(ctx, &got)
	if resultDiags.HasError() {
		t.Fatalf("read state: %#v", resultDiags)
	}
	if got.Name.ValueString() != "old" {
		t.Fatalf("state = %#v, want prior state after failure", got)
	}
}

func outletLifecycleTestAdapterFactory(adapter outletLifecycleTestAdapter) outletLifecycleAdapterFactory[outletLifecycleTestModel] {
	return func(providerData any) (outletLifecycleAdapter[outletLifecycleTestModel], diag.Diagnostics) {
		if _, ok := providerData.(outletLifecycleTestAdapter); !ok {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Unexpected Provider Data Type", "test adapter required")}
		}
		return adapter, nil
	}
}

func configureOutletLifecycle(t *testing.T, ctx context.Context, lifecycle *outletLifecycle[outletLifecycleTestModel], adapter outletLifecycleTestAdapter) {
	t.Helper()
	resp := &resource.ConfigureResponse{}
	lifecycle.Configure(ctx, resource.ConfigureRequest{ProviderData: adapter}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("configure diagnostics: %#v", resp.Diagnostics)
	}
}

func outletLifecycleTestSchema() resourceschema.Schema {
	return resourceschema.Schema{Attributes: map[string]resourceschema.Attribute{
		"id":   resourceschema.Int64Attribute{Computed: true},
		"name": resourceschema.StringAttribute{Required: true},
	}}
}

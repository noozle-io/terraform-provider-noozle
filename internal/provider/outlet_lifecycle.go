package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// outletLifecycle centralizes the Terraform-facing lifecycle shared by typed
// Outlets. Typed adapters retain their route, codec, feature, and material
// facts while this module owns state and diagnostic sequencing.
type outletLifecycle[M any] struct {
	adapter    outletLifecycleAdapter[M]
	newAdapter outletLifecycleAdapterFactory[M]
	outletName string
}

type outletLifecycleAdapter[M any] interface {
	Create(context.Context, tfsdk.Config, M) (outletLifecycleResult[M], diag.Diagnostics)
	Read(context.Context, M) (outletLifecycleResult[M], diag.Diagnostics)
	Update(context.Context, tfsdk.Config, M, M) (outletLifecycleResult[M], diag.Diagnostics)
	Delete(context.Context, M) diag.Diagnostics
}

type outletLifecycleAdapterFactory[M any] func(any) (outletLifecycleAdapter[M], diag.Diagnostics)

type outletLifecycleResult[M any] struct {
	State   M
	Removed bool
}

func newOutletLifecycle[M any](outletName string, newAdapter outletLifecycleAdapterFactory[M]) outletLifecycle[M] {
	return outletLifecycle[M]{outletName: outletName, newAdapter: newAdapter}
}

func (l *outletLifecycle[M]) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if l.newAdapter == nil {
		resp.Diagnostics.AddError("Unconfigured Outlet Lifecycle", "The typed Outlet lifecycle was not initialized.")
		return
	}

	adapter, diags := l.newAdapter(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	l.adapter = adapter
}

func (l *outletLifecycle[M]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if !l.ready(resp.Diagnostics, &resp.Diagnostics) {
		return
	}

	var plan M
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, diags := l.adapter.Create(ctx, req.Config, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	l.setResult(ctx, result, &resp.State, &resp.Diagnostics)
}

func (l *outletLifecycle[M]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if !l.ready(resp.Diagnostics, &resp.Diagnostics) {
		return
	}

	var state M
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, diags := l.adapter.Read(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	l.setResult(ctx, result, &resp.State, &resp.Diagnostics)
}

func (l *outletLifecycle[M]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if !l.ready(resp.Diagnostics, &resp.Diagnostics) {
		return
	}

	var plan M
	var state M
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, diags := l.adapter.Update(ctx, req.Config, plan, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	l.setResult(ctx, result, &resp.State, &resp.Diagnostics)
}

func (l *outletLifecycle[M]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if !l.ready(resp.Diagnostics, &resp.Diagnostics) {
		return
	}

	var state M
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(l.adapter.Delete(ctx, state)...)
}

func (l *outletLifecycle[M]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid "+l.outletName+" Import Identifier",
			fmt.Sprintf("Expected numeric outlet ID, got %q: %s", req.ID, err),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (l *outletLifecycle[M]) ready(existing diag.Diagnostics, target *diag.Diagnostics) bool {
	if existing.HasError() {
		return false
	}
	if l.adapter == nil {
		target.AddError("Unconfigured Noozle Client", "The provider did not supply a configured Noozle client.")
		return false
	}
	return true
}

func (l *outletLifecycle[M]) setResult(ctx context.Context, result outletLifecycleResult[M], state *tfsdk.State, diags *diag.Diagnostics) {
	if result.Removed {
		state.RemoveResource(ctx)
		return
	}
	diags.Append(state.Set(ctx, &result.State)...)
}

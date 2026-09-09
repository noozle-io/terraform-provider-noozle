package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

func TestParseOutletQueryLinkID(t *testing.T) {
	t.Parallel()

	outletID, queryID, err := parseOutletQueryLinkID("99:42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if outletID != 99 || queryID != 42 {
		t.Fatalf("unexpected parsed values: %d %d", outletID, queryID)
	}
}

func TestParseOutletQueryLinkIDRejectsInvalidValue(t *testing.T) {
	t.Parallel()

	_, _, err := parseOutletQueryLinkID("bad")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOutletHasQuery(t *testing.T) {
	t.Parallel()

	if !outletHasQuery(client.OutletRef{ID: 99, QueryIDs: []int64{12, 42}}, 42) {
		t.Fatal("expected query to be present")
	}
	if outletHasQuery(client.OutletRef{ID: 99, QueryIDs: []int64{12}}, 42) {
		t.Fatal("expected query to be absent")
	}
}

func TestOutletQueryLinkImportStateRejectsInvalidID(t *testing.T) {
	t.Parallel()

	r := &queryOutletLinkResource{}
	resp := &resource.ImportStateResponse{}

	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "abc"}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostics for invalid import id")
	}
}

func TestOutletQueryLinkID(t *testing.T) {
	t.Parallel()

	model := queryOutletLinkResourceModel{
		OutletID: types.Int64Value(99),
		QueryID:  types.Int64Value(42),
	}
	got := outletQueryLinkID(model.OutletID.ValueInt64(), model.QueryID.ValueInt64())
	if got != "99:42" {
		t.Fatalf("unexpected composite id: %q", got)
	}
}

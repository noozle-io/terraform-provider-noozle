package provider

import (
	"context"
	"testing"
	"time"

	"terraform-provider-noozle/internal/client"
)

func TestQueryDataSourceModelFromAPI(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 27, 12, 34, 56, 0, time.UTC)
	state, diags := queryDataSourceModelFromAPI(context.Background(), client.Query{
		ID:             12,
		Name:           "Breaking News",
		Description:    "desc",
		FullExpression: "feed:finance",
		Tags:           []string{"finance"},
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if state.ID.ValueInt64() != 12 || state.UpdatedAt.ValueString() != now.Format(time.RFC3339) {
		t.Fatalf("unexpected query data source state: %+v", state)
	}
}

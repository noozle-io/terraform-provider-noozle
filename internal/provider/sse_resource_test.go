package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"terraform-provider-noozle/internal/client"
)

func TestSseOutletModelFromAPI(t *testing.T) {
	t.Parallel()

	state := sseOutletModelFromAPI(client.SseOutlet{
		ID:                101,
		Name:              "SSE Alerts",
		Description:       "desc",
		Enabled:           true,
		EventType:         "article.match",
		Path:              "/v1/sse/stream/alerts",
		RetryMS:           5000,
		NotifyPolicy:      stringPointerForTest("FIRST_ONLY"),
		StreamURL:         "https://api.example.com/v1/sse/stream/alerts",
		ApplyStatus:       "healthy",
		ConfigGeneration:  3,
		AppliedGeneration: 3,
	})

	if state.ID.ValueInt64() != 101 || state.RetryMS.ValueInt64() != 5000 || state.NotifyPolicy.ValueString() != "FIRST_ONLY" {
		t.Fatalf("unexpected state: %+v", state)
	}
	if state.StreamURL.IsNull() {
		t.Fatalf("expected non-null stream_url")
	}
}

func TestSseOutletImportStateRejectsNonNumericID(t *testing.T) {
	t.Parallel()

	r := &sseOutletResource{}
	resp := &resource.ImportStateResponse{}

	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "abc"}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostics for invalid import id")
	}
}

func TestSseOutletSchemasExposeNotifyPolicy(t *testing.T) {
	t.Parallel()

	resourceResponse := &resource.SchemaResponse{}
	(&sseOutletResource{}).Schema(context.Background(), resource.SchemaRequest{}, resourceResponse)
	resourceAttribute, ok := resourceResponse.Schema.Attributes["notify_policy"].(resourceschema.StringAttribute)
	if !ok || !resourceAttribute.Optional {
		t.Errorf("resource notify_policy = %#v, want optional string", resourceResponse.Schema.Attributes["notify_policy"])
	}

	dataSourceResponse := &datasource.SchemaResponse{}
	(&sseOutletDataSource{}).Schema(context.Background(), datasource.SchemaRequest{}, dataSourceResponse)
	dataSourceAttribute, ok := dataSourceResponse.Schema.Attributes["notify_policy"].(datasourceschema.StringAttribute)
	if !ok || !dataSourceAttribute.Computed {
		t.Errorf("data source notify_policy = %#v, want computed string", dataSourceResponse.Schema.Attributes["notify_policy"])
	}
}

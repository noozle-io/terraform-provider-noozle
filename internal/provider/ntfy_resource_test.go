package provider

import (
	"context"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

func TestNtfyOutletModelFromAPI(t *testing.T) {
	t.Parallel()

	state := ntfyOutletModelFromAPI(client.NtfyOutlet{
		ID:                99,
		Name:              "Alerts",
		Description:       "desc",
		Enabled:           true,
		Topic:             "alerts",
		NotifyPolicy:      stringPointerForTest("ALL"),
		ServerURL:         stringPointerForTest("https://ntfy.example.com"),
		ApplyStatus:       "healthy",
		LastApplyError:    "",
		ConfigGeneration:  4,
		AppliedGeneration: 4,
	}, types.StringValue("alerts"))
	if state.ID.ValueInt64() != 99 || state.Topic.ValueString() != "alerts" || state.EffectiveTopic.ValueString() != "alerts" || state.NotifyPolicy.ValueString() != "ALL" || state.ServerURL.ValueString() != "https://ntfy.example.com" {
		t.Fatalf("unexpected state: %+v", state)
	}
	if !state.LastApplyError.IsNull() {
		t.Fatalf("expected null last_apply_error, got %#v", state.LastApplyError)
	}
}

func TestNtfyOutletModelFromAPIFallsBackToEffectiveTopicForImport(t *testing.T) {
	t.Parallel()

	state := ntfyOutletModelFromAPI(prefixedNtfyOutlet(), types.StringNull())
	if got, want := state.Topic.ValueString(), "u45yu3-trump"; got != want {
		t.Fatalf("imported topic = %q, want effective topic %q", got, want)
	}
	if got, want := state.EffectiveTopic.ValueString(), "u45yu3-trump"; got != want {
		t.Fatalf("imported effective topic = %q, want %q", got, want)
	}
}

func TestNtfyOutletAdapterPreservesConfiguredTopicWhenAPIPrefixesIt(t *testing.T) {
	t.Parallel()

	adapter := ntfyOutletAdapter{client: prefixedNtfyOutletClient{}}
	plan := ntfyOutletResourceModel{
		Name:    types.StringValue("Trump alerts"),
		Enabled: types.BoolValue(true),
		Topic:   types.StringValue("trump"),
	}

	created, diags := adapter.Create(context.Background(), tfsdk.Config{}, plan)
	if diags.HasError() {
		t.Fatalf("create diagnostics: %#v", diags)
	}
	if got, want := created.State.Topic.ValueString(), plan.Topic.ValueString(); got != want {
		t.Fatalf("created topic = %q, want configured topic %q", got, want)
	}
	if got, want := created.State.EffectiveTopic.ValueString(), "u45yu3-trump"; got != want {
		t.Fatalf("created effective topic = %q, want %q", got, want)
	}

	read, diags := adapter.Read(context.Background(), created.State)
	if diags.HasError() {
		t.Fatalf("read diagnostics: %#v", diags)
	}
	if got, want := read.State.Topic.ValueString(), plan.Topic.ValueString(); got != want {
		t.Fatalf("refreshed topic = %q, want configured topic %q", got, want)
	}
	if got, want := read.State.EffectiveTopic.ValueString(), "u45yu3-trump"; got != want {
		t.Fatalf("refreshed effective topic = %q, want %q", got, want)
	}
}

func TestNtfyOutletAdapterObservesCreateUpdateAndDelete(t *testing.T) {
	t.Parallel()

	client := &lifecycleNtfyOutletClient{}
	adapter := ntfyOutletAdapter{client: client}
	plan := ntfyOutletResourceModel{
		ID:      types.Int64Value(99),
		Name:    types.StringValue("Alerts"),
		Enabled: types.BoolValue(true),
		Topic:   types.StringValue("alerts"),
	}

	created, diags := adapter.Create(context.Background(), tfsdk.Config{}, plan)
	if diags.HasError() {
		t.Fatalf("create diagnostics: %#v", diags)
	}
	updated, diags := adapter.Update(context.Background(), tfsdk.Config{}, plan, created.State)
	if diags.HasError() {
		t.Fatalf("update diagnostics: %#v", diags)
	}
	if diags := adapter.Delete(context.Background(), updated.State); diags.HasError() {
		t.Fatalf("delete diagnostics: %#v", diags)
	}

	if got, want := client.observedGenerations, []int64{1, 2, 3}; !slices.Equal(got, want) {
		t.Fatalf("observed generations = %v, want %v", got, want)
	}
}

func TestNtfyOutletImportStateRejectsNonNumericID(t *testing.T) {
	t.Parallel()

	r := &ntfyOutletResource{}
	resp := &resource.ImportStateResponse{}

	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "abc"}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected diagnostics for invalid import id")
	}
}

func TestNtfyOutletSchemasExposeOptionalFields(t *testing.T) {
	t.Parallel()

	resourceResponse := &resource.SchemaResponse{}
	(&ntfyOutletResource{}).Schema(context.Background(), resource.SchemaRequest{}, resourceResponse)
	for _, name := range []string{"notify_policy", "server_url"} {
		attribute, ok := resourceResponse.Schema.Attributes[name].(resourceschema.StringAttribute)
		if !ok || !attribute.Optional {
			t.Errorf("resource attribute %q = %#v, want optional string", name, resourceResponse.Schema.Attributes[name])
		}
	}

	effectiveTopic, ok := resourceResponse.Schema.Attributes["effective_topic"].(resourceschema.StringAttribute)
	if !ok || !effectiveTopic.Computed {
		t.Errorf("resource effective_topic = %#v, want computed string", resourceResponse.Schema.Attributes["effective_topic"])
	}

	dataSourceResponse := &datasource.SchemaResponse{}
	(&ntfyOutletDataSource{}).Schema(context.Background(), datasource.SchemaRequest{}, dataSourceResponse)
	for _, name := range []string{"notify_policy", "server_url"} {
		attribute, ok := dataSourceResponse.Schema.Attributes[name].(datasourceschema.StringAttribute)
		if !ok || !attribute.Computed {
			t.Errorf("data source attribute %q = %#v, want computed string", name, dataSourceResponse.Schema.Attributes[name])
		}
	}
}

func TestStringValueOrEmpty(t *testing.T) {
	t.Parallel()

	if got := stringValueOrEmpty(types.StringNull()); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
	if got := stringValueOrEmpty(types.StringValue("value")); got != "value" {
		t.Fatalf("unexpected string value: %q", got)
	}
}

func stringPointerForTest(value string) *string {
	return &value
}

type prefixedNtfyOutletClient struct{}

type lifecycleNtfyOutletClient struct {
	observedGenerations []int64
}

func (c *lifecycleNtfyOutletClient) CreateNtfyOutlet(ctx context.Context, _ client.NtfyOutletCreateRequest) (client.NtfyOutlet, error) {
	client.RecordOutletOperation(ctx, client.OutletOperation{Kind: "create", DesiredGeneration: 1})
	return client.NtfyOutlet{ID: 99, Name: "Alerts", Enabled: true, Topic: "alerts", ConfigGeneration: 1}, nil
}

func (c *lifecycleNtfyOutletClient) GetNtfyOutlet(context.Context, int64) (client.NtfyOutlet, error) {
	return client.NtfyOutlet{ID: 99, Name: "Alerts", Enabled: true, Topic: "alerts", ConfigGeneration: 2}, nil
}

func (c *lifecycleNtfyOutletClient) UpdateNtfyOutlet(ctx context.Context, _ int64, _ client.NtfyOutletUpdateRequest) error {
	client.RecordOutletOperation(ctx, client.OutletOperation{Kind: "update", DesiredGeneration: 2})
	return nil
}

func (c *lifecycleNtfyOutletClient) DeleteNtfyOutlet(context.Context, int64) error { return nil }

func (c *lifecycleNtfyOutletClient) DeleteNtfyOutletOperation(context.Context, int64) (client.OutletOperation, error) {
	return client.OutletOperation{Kind: "delete", DesiredGeneration: 3}, nil
}

func (c *lifecycleNtfyOutletClient) GetTypedOutletStatus(_ context.Context, _ string, _ int64, generation int64) (client.OutletStatus, error) {
	c.observedGenerations = append(c.observedGenerations, generation)
	if generation == 3 {
		return client.OutletStatus{RequestedGeneration: 3, DesiredStatus: "deleted", RuntimeStatus: "deleted"}, nil
	}
	return client.OutletStatus{RequestedGeneration: generation, DesiredStatus: "enabled", RuntimeStatus: "active", Realization: client.OutletRuntimeRealization{Phase: "running", Ready: true}}, nil
}

func (c *lifecycleNtfyOutletClient) SetNtfyTopic(context.Context, int64, string) error { return nil }

func (c *lifecycleNtfyOutletClient) ClearNtfyNotifyPolicy(context.Context, int64) error { return nil }

func (c *lifecycleNtfyOutletClient) ClearNtfyServerURL(context.Context, int64) error { return nil }

func (prefixedNtfyOutletClient) CreateNtfyOutlet(ctx context.Context, _ client.NtfyOutletCreateRequest) (client.NtfyOutlet, error) {
	client.RecordOutletOperation(ctx, client.OutletOperation{Kind: "create", DesiredGeneration: 4})
	return prefixedNtfyOutlet(), nil
}

func (prefixedNtfyOutletClient) GetNtfyOutlet(context.Context, int64) (client.NtfyOutlet, error) {
	return prefixedNtfyOutlet(), nil
}

func (prefixedNtfyOutletClient) UpdateNtfyOutlet(context.Context, int64, client.NtfyOutletUpdateRequest) error {
	return nil
}

func (prefixedNtfyOutletClient) DeleteNtfyOutlet(context.Context, int64) error {
	return nil
}

func (prefixedNtfyOutletClient) DeleteNtfyOutletOperation(context.Context, int64) (client.OutletOperation, error) {
	return client.OutletOperation{Kind: "delete", DesiredGeneration: 4}, nil
}

func (prefixedNtfyOutletClient) GetTypedOutletStatus(context.Context, string, int64, int64) (client.OutletStatus, error) {
	return client.OutletStatus{
		RequestedGeneration: 4,
		DesiredStatus:       "enabled",
		RuntimeStatus:       "active",
		Realization:         client.OutletRuntimeRealization{Phase: "running", Ready: true},
	}, nil
}

func (prefixedNtfyOutletClient) SetNtfyTopic(context.Context, int64, string) error {
	return nil
}

func (prefixedNtfyOutletClient) ClearNtfyNotifyPolicy(context.Context, int64) error {
	return nil
}

func (prefixedNtfyOutletClient) ClearNtfyServerURL(context.Context, int64) error {
	return nil
}

func prefixedNtfyOutlet() client.NtfyOutlet {
	return client.NtfyOutlet{ID: 99, Name: "Trump alerts", Enabled: true, Topic: "u45yu3-trump", ConfigGeneration: 4}
}

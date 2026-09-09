package provider

import (
	"context"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

func TestTransformSupportingOutletSchemasExposeTransformState(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		resource   resource.Resource
		dataSource datasource.DataSource
	}{
		"amqp_0_9":       {resource: &amqp09OutletResource{}, dataSource: &amqp09OutletDataSource{}},
		"amqp_1":         {resource: &amqp1OutletResource{}, dataSource: &amqp1OutletDataSource{}},
		"gcp_pubsub":     {resource: &gcpPubsubOutletResource{}, dataSource: &gcpPubsubOutletDataSource{}},
		"http_server":    {resource: &httpServerOutletResource{}, dataSource: &httpServerOutletDataSource{}},
		"kafka":          {resource: &kafkaOutletResource{}, dataSource: &kafkaOutletDataSource{}},
		"nats_jetstream": {resource: &natsJetstreamOutletResource{}, dataSource: &natsJetstreamOutletDataSource{}},
		"pulsar":         {resource: &pulsarOutletResource{}, dataSource: &pulsarOutletDataSource{}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			resourceResponse := &resource.SchemaResponse{}
			test.resource.Schema(context.Background(), resource.SchemaRequest{}, resourceResponse)
			attribute, ok := resourceResponse.Schema.Attributes["transform"].(resourceschema.SingleNestedAttribute)
			if !ok || !attribute.Optional || !attribute.Computed {
				t.Errorf("resource transform = %#v, want optional computed nested attribute", resourceResponse.Schema.Attributes["transform"])
			}
			for _, field := range []string{"selected", "effective", "state", "applied", "generation", "last_error"} {
				if _, ok := attribute.Attributes[field]; !ok {
					t.Errorf("resource transform missing %q", field)
				}
			}

			dataSourceResponse := &datasource.SchemaResponse{}
			test.dataSource.Schema(context.Background(), datasource.SchemaRequest{}, dataSourceResponse)
			dataSourceAttribute, ok := dataSourceResponse.Schema.Attributes["transform"].(datasourceschema.SingleNestedAttribute)
			if !ok || !dataSourceAttribute.Computed {
				t.Errorf("data source transform = %#v, want computed nested attribute", dataSourceResponse.Schema.Attributes["transform"])
			}
		})
	}
}

func TestTransformStateValuesPreserveSelectionAndRuntimeStatus(t *testing.T) {
	t.Parallel()

	transform := transformObject(client.OutletTransformState{
		SelectedTransform:  stringPointerForTest("summarize"),
		EffectiveTransform: stringPointerForTest("summarize"),
		State:              stringPointerForTest("enabled"),
		Applied:            true,
		Generation:         int64PointerForTransformTest(4),
		LastError:          stringPointerForTest("apply failed"),
	})
	selected := selectedTransform(transform)
	if selected.ValueString() != "summarize" || transformStringAttribute(t, transform, "effective").ValueString() != "summarize" || transformStringAttribute(t, transform, "state").ValueString() != "enabled" || !transformBoolAttribute(t, transform, "applied").ValueBool() || transformInt64Attribute(t, transform, "generation").ValueInt64() != 4 || transformStringAttribute(t, transform, "last_error").ValueString() != "apply failed" {
		t.Fatalf("unexpected transform values: %#v", transform)
	}
	if selected := selectedTransform(transformObject(client.OutletTransformState{})); !selected.IsNull() {
		t.Fatalf("absent selected transform = %#v, want null", selected)
	}
	if selected := selectedTransform(types.ObjectUnknown(outletTransformAttributeTypes)); !selected.IsUnknown() {
		t.Fatalf("unknown transform object selection = %#v, want unknown", selected)
	}
	unknownSelection := types.ObjectValueMust(outletTransformAttributeTypes, map[string]attr.Value{
		"selected":   types.StringUnknown(),
		"effective":  types.StringNull(),
		"state":      types.StringNull(),
		"applied":    types.BoolNull(),
		"generation": types.Int64Null(),
		"last_error": types.StringNull(),
	})
	if selected := selectedTransform(unknownSelection); !selected.IsUnknown() {
		t.Fatalf("unknown selected transform = %#v, want unknown", selected)
	}
}

func TestSyncTypedOutletTransformSetsAndClearsSelection(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name       string
		outletType string
		plan       types.String
		state      types.String
		wantCall   string
	}{
		{name: "amqp_0_9 set", outletType: "amqp_0_9", plan: types.StringValue("summarize"), state: types.StringNull(), wantCall: "set:amqp_0_9:41:summarize"},
		{name: "amqp_1 clear removed selection", outletType: "amqp_1", plan: types.StringNull(), state: types.StringValue("summarize"), wantCall: "clear:amqp_1:41"},
		{name: "gcp pubsub set", outletType: "gcp_pubsub", plan: types.StringValue("summarize"), state: types.StringNull(), wantCall: "set:gcp_pubsub:41:summarize"},
		{name: "http server clear removed selection", outletType: "http_server", plan: types.StringNull(), state: types.StringValue("summarize"), wantCall: "clear:http_server:41"},
		{name: "kafka set", outletType: "kafka", plan: types.StringValue("summarize"), state: types.StringNull(), wantCall: "set:kafka:41:summarize"},
		{name: "nats jetstream clear removed selection", outletType: "nats_jetstream", plan: types.StringNull(), state: types.StringValue("summarize"), wantCall: "clear:nats_jetstream:41"},
		{name: "pulsar set", outletType: "pulsar", plan: types.StringValue("summarize"), state: types.StringNull(), wantCall: "set:pulsar:41:summarize"},
		{name: "unchanged", outletType: "amqp_1", plan: types.StringValue("summarize"), state: types.StringValue("summarize")},
		{name: "unknown selection", outletType: "amqp_1", plan: types.StringUnknown(), state: types.StringValue("summarize")},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			transformClient := &recordingTransformClient{}
			if err := syncTypedOutletTransform(context.Background(), transformClient, test.outletType, 41, test.plan, test.state); err != nil {
				t.Fatalf("sync transform: %v", err)
			}
			if transformClient.call != test.wantCall {
				t.Fatalf("call = %q, want %q", transformClient.call, test.wantCall)
			}
		})
	}
}

func transformStringAttribute(t *testing.T, transform types.Object, name string) types.String {
	t.Helper()
	value, ok := transform.Attributes()[name].(types.String)
	if !ok {
		t.Fatalf("transform attribute %q = %#v, want string", name, transform.Attributes()[name])
	}
	return value
}

func transformBoolAttribute(t *testing.T, transform types.Object, name string) types.Bool {
	t.Helper()
	value, ok := transform.Attributes()[name].(types.Bool)
	if !ok {
		t.Fatalf("transform attribute %q = %#v, want bool", name, transform.Attributes()[name])
	}
	return value
}

func transformInt64Attribute(t *testing.T, transform types.Object, name string) types.Int64 {
	t.Helper()
	value, ok := transform.Attributes()[name].(types.Int64)
	if !ok {
		t.Fatalf("transform attribute %q = %#v, want int64", name, transform.Attributes()[name])
	}
	return value
}

type recordingTransformClient struct{ call string }

func (c *recordingTransformClient) GetTypedOutletTransform(context.Context, string, int64) (client.OutletTransformState, error) {
	return client.OutletTransformState{}, nil
}

func (c *recordingTransformClient) SetTypedOutletTransform(_ context.Context, outletType string, outletID int64, transform string) (client.OutletTransformState, error) {
	c.call = "set:" + outletType + ":" + strconv.FormatInt(outletID, 10) + ":" + transform
	return client.OutletTransformState{}, nil
}

func (c *recordingTransformClient) ClearTypedOutletTransform(_ context.Context, outletType string, outletID int64) (client.OutletTransformState, error) {
	c.call = "clear:" + outletType + ":" + strconv.FormatInt(outletID, 10)
	return client.OutletTransformState{}, nil
}

func int64PointerForTransformTest(value int64) *int64 { return &value }

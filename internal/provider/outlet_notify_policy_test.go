package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBrokerOutletSchemasExposeNotifyPolicy(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		resource   resource.Resource
		dataSource datasource.DataSource
	}{
		"amqp_0_9":       {resource: &amqp09OutletResource{}, dataSource: &amqp09OutletDataSource{}},
		"amqp_1":         {resource: &amqp1OutletResource{}, dataSource: &amqp1OutletDataSource{}},
		"kafka":          {resource: &kafkaOutletResource{}, dataSource: &kafkaOutletDataSource{}},
		"nats_jetstream": {resource: &natsJetstreamOutletResource{}, dataSource: &natsJetstreamOutletDataSource{}},
		"pulsar":         {resource: &pulsarOutletResource{}, dataSource: &pulsarOutletDataSource{}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			resourceResponse := &resource.SchemaResponse{}
			test.resource.Schema(context.Background(), resource.SchemaRequest{}, resourceResponse)
			resourceAttribute, ok := resourceResponse.Schema.Attributes["notify_policy"].(resourceschema.StringAttribute)
			if !ok || !resourceAttribute.Optional {
				t.Errorf("resource notify_policy = %#v, want optional string", resourceResponse.Schema.Attributes["notify_policy"])
			}

			dataSourceResponse := &datasource.SchemaResponse{}
			test.dataSource.Schema(context.Background(), datasource.SchemaRequest{}, dataSourceResponse)
			dataSourceAttribute, ok := dataSourceResponse.Schema.Attributes["notify_policy"].(datasourceschema.StringAttribute)
			if !ok || !dataSourceAttribute.Computed {
				t.Errorf("data source notify_policy = %#v, want computed string", dataSourceResponse.Schema.Attributes["notify_policy"])
			}
		})
	}
}

func TestNullableStringPointerPreservingEmpty(t *testing.T) {
	t.Parallel()

	empty := ""
	if got := nullableStringPointerPreservingEmpty(&empty); !got.Equal(types.StringValue("")) {
		t.Fatalf("empty policy = %#v, want known empty string", got)
	}
	if got := nullableStringPointerPreservingEmpty(nil); !got.IsNull() {
		t.Fatalf("absent policy = %#v, want null", got)
	}
}

func TestCloudAndHTTPOutletSchemasExposeNotifyPolicy(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		resource   resource.Resource
		dataSource datasource.DataSource
	}{
		"aws_sns":     {resource: &awsSnsOutletResource{}, dataSource: &awsSnsOutletDataSource{}},
		"aws_sqs":     {resource: &awsSqsOutletResource{}, dataSource: &awsSqsOutletDataSource{}},
		"gcp_pubsub":  {resource: &gcpPubsubOutletResource{}, dataSource: &gcpPubsubOutletDataSource{}},
		"http_server": {resource: &httpServerOutletResource{}, dataSource: &httpServerOutletDataSource{}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			resourceResponse := &resource.SchemaResponse{}
			test.resource.Schema(context.Background(), resource.SchemaRequest{}, resourceResponse)
			resourceAttribute, ok := resourceResponse.Schema.Attributes["notify_policy"].(resourceschema.StringAttribute)
			if !ok || !resourceAttribute.Optional {
				t.Errorf("resource notify_policy = %#v, want optional string", resourceResponse.Schema.Attributes["notify_policy"])
			}

			dataSourceResponse := &datasource.SchemaResponse{}
			test.dataSource.Schema(context.Background(), datasource.SchemaRequest{}, dataSourceResponse)
			dataSourceAttribute, ok := dataSourceResponse.Schema.Attributes["notify_policy"].(datasourceschema.StringAttribute)
			if !ok || !dataSourceAttribute.Computed {
				t.Errorf("data source notify_policy = %#v, want computed string", dataSourceResponse.Schema.Attributes["notify_policy"])
			}
		})
	}
}

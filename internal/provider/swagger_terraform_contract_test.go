package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"gopkg.in/yaml.v3"
)

type terraformOpenAPIContract struct {
	Paths map[string]terraformOpenAPIPath `yaml:"paths"`
}

type terraformOpenAPIPath struct {
	Outlet *terraformOutletContract `yaml:"x-noozle-outlet-provider"`
}

type terraformOutletContract struct {
	ID                 string                     `yaml:"id"`
	SupportsTransforms bool                       `yaml:"supports_transforms"`
	Fields             []terraformFieldContract   `yaml:"fields"`
	Features           []terraformFeatureContract `yaml:"features"`
}

type terraformFieldContract struct {
	ID         string `yaml:"id"`
	ConfigPath string `yaml:"config_path"`
	Type       string `yaml:"type"`
	Required   bool   `yaml:"required"`
	Sensitive  bool   `yaml:"sensitive"`
}

type terraformFeatureContract struct {
	ID string `yaml:"id"`
}

type terraformOutletTarget struct {
	resource   resource.Resource
	dataSource datasource.DataSource
}

// TestTerraformOpenAPIOutletSchemasMatchContract keeps the Terraform-facing
// schema aligned with the typed outlet overlay. Sensitive fields are excluded:
// they are intentionally represented by write-only credential or feature
// blocks, rather than returned as ordinary outlet attributes.
func TestTerraformOpenAPIOutletSchemasMatchContract(t *testing.T) {
	t.Parallel()

	contract := loadTerraformOpenAPIContract(t)
	targets := map[string]terraformOutletTarget{
		"amqp_0_9":       {resource: NewAmqp09OutletResource(), dataSource: NewAmqp09OutletDataSource()},
		"amqp_1":         {resource: NewAmqp1OutletResource(), dataSource: NewAmqp1OutletDataSource()},
		"aws_sns":        {resource: NewAwsSnsOutletResource(), dataSource: NewAwsSnsOutletDataSource()},
		"aws_sqs":        {resource: NewAwsSqsOutletResource(), dataSource: NewAwsSqsOutletDataSource()},
		"gcp_pubsub":     {resource: NewGcpPubsubOutletResource(), dataSource: NewGcpPubsubOutletDataSource()},
		"http_server":    {resource: NewHttpServerOutletResource(), dataSource: NewHttpServerOutletDataSource()},
		"kafka":          {resource: NewKafkaOutletResource(), dataSource: NewKafkaOutletDataSource()},
		"nats_jetstream": {resource: NewNatsJetstreamOutletResource(), dataSource: NewNatsJetstreamOutletDataSource()},
		"ntfy":           {resource: NewNtfyOutletResource(), dataSource: NewNtfyOutletDataSource()},
		"pulsar":         {resource: NewPulsarOutletResource(), dataSource: NewPulsarOutletDataSource()},
		"sse":            {resource: NewSseOutletResource(), dataSource: NewSseOutletDataSource()},
	}

	if len(contract.Paths) == 0 {
		t.Fatal("OpenAPI contract has no paths")
	}

	for outletID, target := range targets {
		outlet := contract.outlet(outletID)
		if outlet == nil {
			t.Fatalf("OpenAPI contract has no root typed path for %q", outletID)
		}

		t.Run(outletID, func(t *testing.T) {
			t.Parallel()

			resourceResponse := &resource.SchemaResponse{}
			target.resource.Schema(context.Background(), resource.SchemaRequest{}, resourceResponse)
			dataSourceResponse := &datasource.SchemaResponse{}
			target.dataSource.Schema(context.Background(), datasource.SchemaRequest{}, dataSourceResponse)

			for _, field := range outlet.Fields {
				if field.Sensitive {
					continue
				}

				attributeName := terraformAttributeName(outletID, field.ID)
				resourceAttribute, ok := resourceResponse.Schema.Attributes[attributeName]
				if !ok {
					t.Errorf("resource is missing %q for OpenAPI field %q", attributeName, field.ID)
					continue
				}
				assertResourceFieldContract(t, attributeName, field, resourceAttribute)

				dataSourceAttribute, ok := dataSourceResponse.Schema.Attributes[attributeName]
				if !ok {
					t.Errorf("data source is missing %q for OpenAPI field %q", attributeName, field.ID)
					continue
				}
				assertDataSourceFieldContract(t, attributeName, field, dataSourceAttribute)
			}

			for _, feature := range outlet.Features {
				featureName := terraformFeatureName(outletID, feature.ID)
				if _, ok := resourceResponse.Schema.Blocks[featureName]; !ok {
					t.Errorf("resource is missing feature block %q for OpenAPI feature %q", featureName, feature.ID)
				}
				if _, ok := dataSourceResponse.Schema.Blocks[featureName]; !ok {
					t.Errorf("data source is missing feature block %q for OpenAPI feature %q", featureName, feature.ID)
				}
			}

			_, resourceHasTransform := resourceResponse.Schema.Attributes["transform"]
			_, dataSourceHasTransform := dataSourceResponse.Schema.Attributes["transform"]
			if resourceHasTransform != outlet.SupportsTransforms {
				t.Errorf("resource transform support = %t, want %t", resourceHasTransform, outlet.SupportsTransforms)
			}
			if dataSourceHasTransform != outlet.SupportsTransforms {
				t.Errorf("data source transform support = %t, want %t", dataSourceHasTransform, outlet.SupportsTransforms)
			}
		})
	}

	for path, pathContract := range contract.Paths {
		if pathContract.Outlet == nil || path != "/v1/terraform/"+pathContract.Outlet.ID {
			continue
		}
		if _, ok := targets[pathContract.Outlet.ID]; !ok {
			t.Errorf("OpenAPI contract defines unregistered outlet %q", pathContract.Outlet.ID)
		}
	}
}

func loadTerraformOpenAPIContract(t *testing.T) terraformOpenAPIContract {
	t.Helper()

	contractDirectory := os.Getenv("NOOZLE_OPENAPI_DIR")
	if contractDirectory == "" {
		t.Skip("set NOOZLE_OPENAPI_DIR to run Swagger-derived contract tests")
	}

	contents, err := os.ReadFile(filepath.Join(contractDirectory, "swagger-terraform.yaml"))
	if err != nil {
		t.Fatalf("read Terraform OpenAPI contract: %v", err)
	}

	var contract terraformOpenAPIContract
	if err := yaml.Unmarshal(contents, &contract); err != nil {
		t.Fatalf("parse Terraform OpenAPI contract: %v", err)
	}

	return contract
}

func (c terraformOpenAPIContract) outlet(outletID string) *terraformOutletContract {
	path, ok := c.Paths["/v1/terraform/"+outletID]
	if !ok || path.Outlet == nil || path.Outlet.ID != outletID {
		return nil
	}

	return path.Outlet
}

func terraformAttributeName(outletID, fieldID string) string {
	fieldID = strings.TrimPrefix(fieldID, outletID+".")
	return strings.ReplaceAll(fieldID, ".", "_")
}

func terraformFeatureName(outletID, featureID string) string {
	featureID = strings.TrimPrefix(featureID, outletID+".")
	featureID = strings.TrimPrefix(featureID, "nats.auth.")
	return strings.ReplaceAll(featureID, ".", "_")
}

func assertResourceFieldContract(t *testing.T, name string, field terraformFieldContract, attribute resourceschema.Attribute) {
	t.Helper()
	assertTerraformAttributeType(t, name, field.Type, attribute)

	if required := schemaBool(attribute, "Required"); required != field.Required {
		t.Errorf("resource %q required = %t, want %t", name, required, field.Required)
	}
	if optional := schemaBool(attribute, "Optional"); optional == field.Required {
		t.Errorf("resource %q optional = %t, want %t", name, optional, !field.Required)
	}
}

func assertDataSourceFieldContract(t *testing.T, name string, field terraformFieldContract, attribute datasourceschema.Attribute) {
	t.Helper()
	assertTerraformAttributeType(t, name, field.Type, attribute)

	if !schemaBool(attribute, "Computed") {
		t.Errorf("data source %q computed = false, want true", name)
	}
}

func assertTerraformAttributeType(t *testing.T, name, contractType string, attribute any) {
	t.Helper()

	want := map[string]string{
		"string":      "StringAttribute",
		"boolean":     "BoolAttribute",
		"integer":     "Int64Attribute",
		"number":      "Float64Attribute",
		"float":       "Float64Attribute",
		"string_list": "ListAttribute",
		"string_map":  "MapAttribute",
		"object":      "MapAttribute",
	}[contractType]
	if want == "" {
		t.Errorf("%q has unsupported OpenAPI field type %q", name, contractType)
		return
	}

	got := reflect.TypeOf(attribute).Name()
	if got != want {
		t.Errorf("%q schema type = %s, want %s for OpenAPI type %q", name, got, want, contractType)
	}
}

func schemaBool(attribute any, name string) bool {
	value := reflect.ValueOf(attribute)
	field := value.FieldByName(name)
	if !field.IsValid() || field.Kind() != reflect.Bool {
		panic(fmt.Sprintf("%T has no boolean %s field", attribute, name))
	}
	return field.Bool()
}

package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// swaggerHTTPContract is deliberately limited to the HTTP details exercised by
// the client. Schema-to-Terraform mapping remains covered by the provider
// contract test.
type swaggerHTTPContract struct {
	Paths map[string]swaggerHTTPPath `yaml:"paths"`
}

type swaggerHTTPPath struct {
	Get    swaggerHTTPOperation `yaml:"get"`
	Post   swaggerHTTPOperation `yaml:"post"`
	Put    swaggerHTTPOperation `yaml:"put"`
	Delete swaggerHTTPOperation `yaml:"delete"`
	Patch  swaggerHTTPOperation `yaml:"patch"`
}

type swaggerHTTPOperation struct {
	RequestBody *swaggerHTTPRequestBody `yaml:"requestBody"`
	Responses   map[string]any          `yaml:"responses"`
}

type swaggerHTTPRequestBody struct {
	Required bool           `yaml:"required"`
	Content  map[string]any `yaml:"content"`
}

func TestSwaggerDerivedQueryAndOutletLinkHTTPContracts(t *testing.T) {
	t.Parallel()

	contract := loadSwaggerHTTPContract(t, "swagger.yaml")
	const tenantID = "t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q"

	tests := []struct {
		name       string
		path       string
		method     string
		resourceID int64
		jsonBody   bool
		response   string
		invoke     func(*Client) error
	}{
		{
			name:     "create query",
			path:     "/v1/tenants/{tenant_id}/queries",
			method:   http.MethodPost,
			jsonBody: true,
			response: queryResponseJSON,
			invoke: func(c *Client) error {
				_, err := c.CreateQuery(context.Background(), QueryRequest{Name: "Alerts", FullExpression: "tag:alert"})
				return err
			},
		},
		{
			name:       "get query",
			path:       "/v1/tenants/{tenant_id}/queries/{id}",
			method:     http.MethodGet,
			resourceID: 42,
			response:   queryResponseJSON,
			invoke: func(c *Client) error {
				_, err := c.GetQuery(context.Background(), 42)
				return err
			},
		},
		{
			name:       "update query",
			path:       "/v1/tenants/{tenant_id}/queries/{id}",
			method:     http.MethodPut,
			resourceID: 42,
			jsonBody:   true,
			response:   queryResponseJSON,
			invoke: func(c *Client) error {
				_, err := c.UpdateQuery(context.Background(), 42, QueryRequest{Name: "Alerts", FullExpression: "tag:alert"})
				return err
			},
		},
		{
			name:       "delete query",
			path:       "/v1/tenants/{tenant_id}/queries/{id}",
			method:     http.MethodDelete,
			resourceID: 42,
			invoke: func(c *Client) error {
				return c.DeleteQuery(context.Background(), 42)
			},
		},
		{
			name:       "link outlet query",
			path:       "/v1/tenants/{tenant_id}/outlets/{id}/queries",
			method:     http.MethodPost,
			resourceID: 99,
			jsonBody:   true,
			response:   `{}`,
			invoke: func(c *Client) error {
				return c.LinkOutletQuery(context.Background(), 99, 42)
			},
		},
		{
			name:       "unlink outlet query",
			path:       "/v1/tenants/{tenant_id}/outlets/{id}/queries",
			method:     http.MethodDelete,
			resourceID: 99,
			jsonBody:   true,
			response:   `{}`,
			invoke: func(c *Client) error {
				return c.UnlinkOutletQuery(context.Background(), 99, 42)
			},
		},
		{
			name:       "get outlet",
			path:       "/v1/tenants/{tenant_id}/outlets/{id}",
			method:     http.MethodGet,
			resourceID: 99,
			response:   `{"id":99,"queries":[42]}`,
			invoke: func(c *Client) error {
				_, err := c.GetOutlet(context.Background(), 99)
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			op := contract.operation(t, test.path, test.method)
			assertJSONRequestContract(t, op, test.jsonBody)
			status := op.singleSuccessStatus(t)
			client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
				assertSwaggerRequest(t, r, test.method, swaggerPath(test.path, tenantID, test.resourceID))
				return jsonResponse(status, test.response), nil
			}, WithTenantContext(TenantContext{TenantID: tenantID}))
			if err := test.invoke(client); err != nil {
				t.Fatalf("client operation: %v", err)
			}
		})
	}
}

func TestSwaggerDerivedTypedOutletDeleteHTTPContracts(t *testing.T) {
	t.Parallel()

	contract := loadSwaggerHTTPContract(t, "swagger-terraform.yaml")
	deleteOutlets := typedOutletDeleteClients()

	if len(deleteOutlets) != 11 {
		t.Fatalf("typed outlet clients = %d, want 11", len(deleteOutlets))
	}

	for outletType, deleteOutlet := range deleteOutlets {
		t.Run(outletType, func(t *testing.T) {
			path := "/v1/terraform/" + outletType + "/{id}"
			op := contract.operation(t, path, http.MethodDelete)
			statuses := op.successStatuses(t)
			for _, status := range statuses {
				t.Run(fmt.Sprintf("status_%d", status), func(t *testing.T) {
					client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
						assertSwaggerRequest(t, r, http.MethodDelete, "/v1/terraform/"+outletType+"/99")
						if status == http.StatusNoContent {
							return jsonResponse(status, ""), nil
						}
						return jsonResponse(status, `{"operation":{"kind":"delete","desired_generation":2}}`), nil
					}, WithTenantContext(TenantContext{TenantID: "t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q"}))
					if err := deleteOutlet(client, context.Background(), 99); err != nil {
						t.Fatalf("delete outlet: %v", err)
					}
				})
			}
		})
	}
}

func TestSwaggerDerivedTypedOutletRootRequestContracts(t *testing.T) {
	t.Parallel()

	contract := loadSwaggerHTTPContract(t, "swagger-terraform.yaml")
	for path, item := range contract.Paths {
		outletType, ok := strings.CutPrefix(path, "/v1/terraform/")
		if !ok || strings.Contains(outletType, "/") {
			continue
		}
		t.Run(outletType, func(t *testing.T) {
			op := item.Post
			if len(op.Responses) == 0 {
				t.Errorf("POST %s is missing from the Terraform overlay", path)
				return
			}
			if op.RequestBody == nil || !op.RequestBody.Required || !slices.Contains(mapKeys(op.RequestBody.Content), "multipart/form-data") {
				t.Errorf("POST %s must require multipart/form-data payload", path)
			}
		})
	}
}

func typedOutletDeleteClients() map[string]func(*Client, context.Context, int64) error {
	return map[string]func(*Client, context.Context, int64) error{
		"amqp_0_9":       func(c *Client, ctx context.Context, id int64) error { return c.DeleteAmqp09Outlet(ctx, id) },
		"amqp_1":         func(c *Client, ctx context.Context, id int64) error { return c.DeleteAmqp1Outlet(ctx, id) },
		"aws_sns":        func(c *Client, ctx context.Context, id int64) error { return c.DeleteAwsSnsOutlet(ctx, id) },
		"aws_sqs":        func(c *Client, ctx context.Context, id int64) error { return c.DeleteAwsSqsOutlet(ctx, id) },
		"gcp_pubsub":     func(c *Client, ctx context.Context, id int64) error { return c.DeleteGcpPubsubOutlet(ctx, id) },
		"http_server":    func(c *Client, ctx context.Context, id int64) error { return c.DeleteHttpServerOutlet(ctx, id) },
		"kafka":          func(c *Client, ctx context.Context, id int64) error { return c.DeleteKafkaOutlet(ctx, id) },
		"nats_jetstream": func(c *Client, ctx context.Context, id int64) error { return c.DeleteNatsJetstreamOutlet(ctx, id) },
		"ntfy":           func(c *Client, ctx context.Context, id int64) error { return c.DeleteNtfyOutlet(ctx, id) },
		"pulsar":         func(c *Client, ctx context.Context, id int64) error { return c.DeletePulsarOutlet(ctx, id) },
		"sse":            func(c *Client, ctx context.Context, id int64) error { return c.DeleteSseOutlet(ctx, id) },
	}
}

func loadSwaggerHTTPContract(t *testing.T, filename string) swaggerHTTPContract {
	t.Helper()
	contractDirectory := os.Getenv("NOOZLE_OPENAPI_DIR")
	if contractDirectory == "" {
		t.Skip("set NOOZLE_OPENAPI_DIR to run Swagger-derived contract tests")
	}

	contents, err := os.ReadFile(filepath.Join(contractDirectory, filename))
	if err != nil {
		t.Fatalf("read Swagger contract: %v", err)
	}
	var contract swaggerHTTPContract
	if err := yaml.Unmarshal(contents, &contract); err != nil {
		t.Fatalf("parse Swagger contract: %v", err)
	}
	return contract
}

func (c swaggerHTTPContract) operation(t *testing.T, path, method string) swaggerHTTPOperation {
	t.Helper()
	item, ok := c.Paths[path]
	if !ok {
		t.Fatalf("Swagger contract has no path %s", path)
	}
	var op swaggerHTTPOperation
	switch method {
	case http.MethodGet:
		op = item.Get
	case http.MethodPost:
		op = item.Post
	case http.MethodPut:
		op = item.Put
	case http.MethodDelete:
		op = item.Delete
	case http.MethodPatch:
		op = item.Patch
	default:
		t.Fatalf("unsupported HTTP method %s", method)
	}
	if len(op.Responses) == 0 {
		t.Fatalf("Swagger contract has no %s operation for %s", method, path)
	}
	return op
}

func (op swaggerHTTPOperation) successStatuses(t *testing.T) []int {
	t.Helper()
	statuses := make([]int, 0, len(op.Responses))
	for status := range op.Responses {
		var code int
		if _, err := fmt.Sscanf(status, "%d", &code); err == nil && code >= 200 && code < 300 {
			statuses = append(statuses, code)
		}
	}
	slices.Sort(statuses)
	if len(statuses) == 0 {
		t.Fatal("Swagger operation has no successful response")
	}
	return statuses
}

func (op swaggerHTTPOperation) singleSuccessStatus(t *testing.T) int {
	t.Helper()
	statuses := op.successStatuses(t)
	if len(statuses) != 1 {
		t.Fatalf("Swagger operation success statuses = %v, want exactly one", statuses)
	}
	return statuses[0]
}

func assertJSONRequestContract(t *testing.T, op swaggerHTTPOperation, wantBody bool) {
	t.Helper()
	if !wantBody {
		return
	}
	if op.RequestBody == nil || !op.RequestBody.Required || !slices.Contains(mapKeys(op.RequestBody.Content), "application/json") {
		t.Fatal("Swagger operation must require an application/json request body")
	}
}

func assertSwaggerRequest(t *testing.T, request *http.Request, method, path string) {
	t.Helper()
	if request.Method != method || request.URL.Path != path {
		t.Errorf("request = %s %s, want %s %s", request.Method, request.URL.Path, method, path)
	}
	if request.Header.Get("X-Noozle-Tenant") == "" {
		t.Error("request is missing X-Noozle-Tenant")
	}
}

func swaggerPath(path, tenantID string, resourceID int64) string {
	return strings.NewReplacer("{tenant_id}", tenantID, "{id}", fmt.Sprintf("%d", resourceID)).Replace(path)
}

func mapKeys(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

const queryResponseJSON = `{"id":42,"name":"Alerts","full_expression":"tag:alert","tags":[],"is_active":true,"created_at":"2026-09-03T00:00:00Z","updated_at":"2026-09-03T00:00:00Z"}`

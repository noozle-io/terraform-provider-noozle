package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"slices"
	"testing"
)

func TestBentoOutletDeletesAcceptAsyncIntent(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodDelete {
			t.Fatalf("method = %s, want DELETE", r.Method)
		}
		return jsonResponse(http.StatusAccepted, `{"outlet_id":99,"tenant_id":"t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q","desired_status":"deleting","operation":{"kind":"delete","desired_generation":2,"status_link":"/status?generation=2","realization":{"active_generation":1,"realization_attempt":1,"phase":"draining","ready":false,"degraded":false,"last_error":null}}}`), nil
	}, WithTenantContext(TenantContext{TenantID: "t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q"}))

	deleteOutlets := map[string]func(context.Context, int64) error{
		"amqp_0_9":       client.DeleteAmqp09Outlet,
		"amqp_1":         client.DeleteAmqp1Outlet,
		"aws_sns":        client.DeleteAwsSnsOutlet,
		"aws_sqs":        client.DeleteAwsSqsOutlet,
		"gcp_pubsub":     client.DeleteGcpPubsubOutlet,
		"http_server":    client.DeleteHttpServerOutlet,
		"kafka":          client.DeleteKafkaOutlet,
		"nats_jetstream": client.DeleteNatsJetstreamOutlet,
		"ntfy":           client.DeleteNtfyOutlet,
		"pulsar":         client.DeletePulsarOutlet,
	}

	for outletType, deleteOutlet := range deleteOutlets {
		t.Run(outletType, func(t *testing.T) {
			if err := deleteOutlet(context.Background(), 99); err != nil {
				t.Fatalf("delete outlet: %v", err)
			}
		})
	}
}

func TestNtfyOutletOperationsSendTenantContext(t *testing.T) {
	t.Parallel()

	const (
		tenantID       = "t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q"
		organizationID = "org_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q"
	)

	var requests []struct {
		method string
		path   string
		tenant string
		org    string
	}
	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		requests = append(requests, struct {
			method string
			path   string
			tenant string
			org    string
		}{
			method: r.Method,
			path:   r.URL.Path,
			tenant: r.Header.Get("X-Noozle-Tenant"),
			org:    r.Header.Get("X-Noozle-Organization"),
		})

		if r.Method == http.MethodPost {
			return jsonResponse(http.StatusCreated, `{"id":99,"name":"Alerts","enabled":true,"config":{"ntfy.topic":"alerts"}}`), nil
		}
		return jsonResponse(http.StatusOK, `{"id":99,"name":"Alerts","enabled":true,"config":{"ntfy.topic":"alerts"}}`), nil
	}, WithTenantContext(TenantContext{TenantID: tenantID, OrganizationID: organizationID}))

	if _, err := client.CreateNtfyOutlet(context.Background(), NtfyOutletCreateRequest{Name: "Alerts", Enabled: true, Topic: "alerts"}); err != nil {
		t.Fatalf("create ntfy outlet: %v", err)
	}
	if _, err := client.GetNtfyOutlet(context.Background(), 99); err != nil {
		t.Fatalf("get ntfy outlet: %v", err)
	}

	want := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/v1/terraform/ntfy"},
		{http.MethodGet, "/v1/terraform/ntfy/99"},
	}
	if len(requests) != len(want) {
		t.Fatalf("requests = %d, want %d", len(requests), len(want))
	}
	for i, expected := range want {
		got := requests[i]
		if got.method != expected.method || got.path != expected.path {
			t.Errorf("request %d = %s %s, want %s %s", i, got.method, got.path, expected.method, expected.path)
		}
		if got.tenant != tenantID || got.org != organizationID {
			t.Errorf("request %d tenant context = tenant %q, organization %q", i, got.tenant, got.org)
		}
	}
}

func TestNtfyOutletOperationsRequireTenantContext(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		t.Fatalf("ntfy outlet operation should not make an HTTP request without tenant context: %s %s", r.Method, r.URL.Path)
		return nil, nil
	}, WithTenantContext(TenantContext{}))

	for name, operation := range map[string]func() error{
		"create": func() error {
			_, err := client.CreateNtfyOutlet(context.Background(), NtfyOutletCreateRequest{Name: "Alerts", Enabled: true, Topic: "alerts"})
			return err
		},
		"read": func() error {
			_, err := client.GetNtfyOutlet(context.Background(), 99)
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			if err := operation(); !errors.Is(err, ErrMissingTenantContext) {
				t.Fatalf("expected missing tenant context error, got %v", err)
			}
		})
	}
}

func TestCreateNtfyOutletIncludesProviderFields(t *testing.T) {
	t.Parallel()

	var payload typedOutletCreateRequest

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		parseMultipartPayload(t, r, &payload, nil)

		return jsonResponse(http.StatusCreated, `{
			"id": 99,
			"name": "Alerts",
			"description": "desc",
			"enabled": true,
			"config": {
				"notify_policy": "ALL",
				"ntfy.server_url": "https://ntfy.example.com",
				"ntfy.topic": "alerts"
			}
		}`), nil
	})

	outlet, err := client.CreateNtfyOutlet(context.Background(), NtfyOutletCreateRequest{
		Name:         "Alerts",
		Description:  "desc",
		Enabled:      true,
		Topic:        "alerts",
		NotifyPolicy: stringPointer("ALL"),
		ServerURL:    stringPointer("https://ntfy.example.com"),
	})
	if err != nil {
		t.Fatalf("create ntfy outlet: %v", err)
	}

	if payload.Fields["ntfy.topic"] != "alerts" {
		t.Fatalf("unexpected request fields: %+v", payload.Fields)
	}
	if payload.Fields["notify_policy"] != "ALL" || payload.Fields["ntfy.server_url"] != "https://ntfy.example.com" {
		t.Fatalf("missing optional ntfy fields: %+v", payload.Fields)
	}
	if outlet.Topic != "alerts" || outlet.NotifyPolicy == nil || *outlet.NotifyPolicy != "ALL" || outlet.ServerURL == nil || *outlet.ServerURL != "https://ntfy.example.com" {
		t.Fatalf("unexpected outlet response: %+v", outlet)
	}
}

func TestCreateNtfyOutletOmitsOptionalFields(t *testing.T) {
	t.Parallel()

	var payload typedOutletCreateRequest
	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		parseMultipartPayload(t, r, &payload, nil)
		return jsonResponse(http.StatusCreated, `{
			"id": 99,
			"name": "Alerts",
			"enabled": true,
			"config": {"ntfy.topic": "alerts"}
		}`), nil
	})

	if _, err := client.CreateNtfyOutlet(context.Background(), NtfyOutletCreateRequest{
		Name:    "Alerts",
		Enabled: true,
		Topic:   "alerts",
	}); err != nil {
		t.Fatalf("create ntfy outlet: %v", err)
	}
	if _, ok := payload.Fields["notify_policy"]; ok {
		t.Fatalf("unexpected notify_policy field: %+v", payload.Fields)
	}
	if _, ok := payload.Fields["ntfy.server_url"]; ok {
		t.Fatalf("unexpected ntfy.server_url field: %+v", payload.Fields)
	}
}

func TestUpdateSseOutletUsesMultipartPayload(t *testing.T) {
	t.Parallel()

	var payload typedOutletUpdateRequest
	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		parseMultipartPayload(t, r, &payload, nil)
		return jsonResponse(http.StatusOK, `{"status":"ok"}`), nil
	})

	if err := client.UpdateSseOutlet(context.Background(), 101, SseOutletUpdateRequest{
		Name:         "SSE Alerts",
		Description:  "desc",
		Enabled:      true,
		EventType:    "article.match",
		Path:         "/events",
		RetryMS:      7000,
		NotifyPolicy: stringPointer("FIRST_ONLY"),
	}); err != nil {
		t.Fatalf("update sse outlet: %v", err)
	}

	if payload.Fields["sse.event_type"] != "article.match" {
		t.Fatalf("unexpected request payload: %+v", payload)
	}
	if payload.Fields["notify_policy"] != "FIRST_ONLY" {
		t.Fatalf("missing notify_policy: %+v", payload.Fields)
	}
}

func TestGetSseOutletDecodesTypedConfig(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{
			"id": 101,
			"name": "SSE Alerts",
			"description": "desc",
			"enabled": true,
			"stream_url": "https://api.example.com/v1/sse/stream/alerts",
			"config": {
				"notify_policy": "ALL",
				"sse.event_type": "article.match",
				"sse.path": "/v1/sse/stream/alerts",
				"sse.retry_ms": 5000
			}
		}`), nil
	})

	outlet, err := client.GetSseOutlet(context.Background(), 101)
	if err != nil {
		t.Fatalf("get sse outlet: %v", err)
	}
	if outlet.EventType != "article.match" || outlet.RetryMS != 5000 || outlet.NotifyPolicy == nil || *outlet.NotifyPolicy != "ALL" {
		t.Fatalf("unexpected outlet: %+v", outlet)
	}
}

func TestClearOptionalOutletFieldsUsesTypedFieldRoutes(t *testing.T) {
	t.Parallel()

	var paths []string
	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		paths = append(paths, r.URL.Path)
		if r.Method != http.MethodDelete {
			t.Fatalf("method = %s, want DELETE", r.Method)
		}
		return jsonResponse(http.StatusOK, `{"status":"ok"}`), nil
	})

	if err := client.ClearNtfyNotifyPolicy(context.Background(), 99); err != nil {
		t.Fatalf("clear ntfy notify policy: %v", err)
	}
	if err := client.ClearNtfyServerURL(context.Background(), 99); err != nil {
		t.Fatalf("clear ntfy server URL: %v", err)
	}
	if err := client.ClearSseNotifyPolicy(context.Background(), 101); err != nil {
		t.Fatalf("clear SSE notify policy: %v", err)
	}

	want := []string{
		"/v1/terraform/ntfy/99/fields/notify_policy",
		"/v1/terraform/ntfy/99/fields/ntfy.server_url",
		"/v1/terraform/sse/101/fields/notify_policy",
	}
	if !slices.Equal(paths, want) {
		t.Fatalf("field paths = %v, want %v", paths, want)
	}
}

func TestBrokerOutletUpdatesIncludeNotifyPolicy(t *testing.T) {
	t.Parallel()

	tests := map[string]func(*Client) error{
		"amqp_0_9": func(client *Client) error {
			return client.UpdateAmqp09Outlet(context.Background(), 1, Amqp09OutletUpdateRequest{Enabled: true, Exchange: "alerts", Key: "alerts.critical", NotifyPolicy: stringPointer("FIRST_ONLY")})
		},
		"amqp_1": func(client *Client) error {
			return client.UpdateAmqp1Outlet(context.Background(), 1, Amqp1OutletUpdateRequest{Enabled: true, TargetAddress: "queue://alerts", NotifyPolicy: stringPointer("FIRST_ONLY")})
		},
		"kafka": func(client *Client) error {
			return client.UpdateKafkaOutlet(context.Background(), 1, KafkaOutletUpdateRequest{Enabled: true, Addresses: []string{"broker:9092"}, Topic: "alerts", NotifyPolicy: stringPointer("FIRST_ONLY")})
		},
		"nats_jetstream": func(client *Client) error {
			return client.UpdateNatsJetstreamOutlet(context.Background(), 1, NatsJetstreamOutletUpdateRequest{Enabled: true, Subject: "alerts", NotifyPolicy: stringPointer("FIRST_ONLY")})
		},
		"pulsar": func(client *Client) error {
			return client.UpdatePulsarOutlet(context.Background(), 1, PulsarOutletUpdateRequest{Enabled: true, Topic: "persistent://public/default/alerts", URL: "pulsar://broker:6650", NotifyPolicy: stringPointer("FIRST_ONLY")})
		},
	}

	for name, update := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var payload typedOutletUpdateRequest
			client := newTestClient(t, func(request *http.Request) (*http.Response, error) {
				parseMultipartPayload(t, request, &payload, nil)
				return jsonResponse(http.StatusOK, `{"status":"ok"}`), nil
			})

			if err := update(client); err != nil {
				t.Fatalf("update outlet: %v", err)
			}
			if got := payload.Fields["notify_policy"]; got != "FIRST_ONLY" {
				t.Fatalf("notify_policy = %#v, want FIRST_ONLY", got)
			}
		})
	}
}

func TestCloudAndHTTPOutletUpdatesIncludeNotifyPolicy(t *testing.T) {
	t.Parallel()

	tests := map[string]func(*Client) error{
		"aws_sns": func(client *Client) error {
			return client.UpdateAwsSnsOutlet(context.Background(), 1, AwsSnsOutletUpdateRequest{Enabled: true, TopicARN: "arn:aws:sns:us-east-1:123456789012:alerts", NotifyPolicy: stringPointer("FIRST_ONLY")})
		},
		"aws_sqs": func(client *Client) error {
			return client.UpdateAwsSqsOutlet(context.Background(), 1, AwsSqsOutletUpdateRequest{Enabled: true, QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/alerts", NotifyPolicy: stringPointer("FIRST_ONLY")})
		},
		"gcp_pubsub": func(client *Client) error {
			return client.UpdateGcpPubsubOutlet(context.Background(), 1, GcpPubsubOutletUpdateRequest{Enabled: true, Project: "noozle", Topic: "alerts", NotifyPolicy: stringPointer("FIRST_ONLY")})
		},
		"http_server": func(client *Client) error {
			return client.UpdateHttpServerOutlet(context.Background(), 1, HttpServerOutletUpdateRequest{Enabled: true, Path: "/alerts", NotifyPolicy: stringPointer("FIRST_ONLY")})
		},
	}

	for name, update := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var payload typedOutletUpdateRequest
			client := newTestClient(t, func(request *http.Request) (*http.Response, error) {
				parseMultipartPayload(t, request, &payload, nil)
				return jsonResponse(http.StatusOK, `{"status":"ok"}`), nil
			})

			if err := update(client); err != nil {
				t.Fatalf("update outlet: %v", err)
			}
			if got := payload.Fields["notify_policy"]; got != "FIRST_ONLY" {
				t.Fatalf("notify_policy = %#v, want FIRST_ONLY", got)
			}
		})
	}
}

func TestCloudAndHTTPOutletCreatesIncludeNotifyPolicy(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		response string
		create   func(*Client) error
	}{
		"aws_sns": {
			response: `{"id":1,"name":"alerts","enabled":true,"config":{"notify_policy":"FIRST_ONLY","aws_sns.topic_arn":"arn:aws:sns:us-east-1:123456789012:alerts"}}`,
			create: func(client *Client) error {
				_, err := client.CreateAwsSnsOutlet(context.Background(), AwsSnsOutletCreateRequest{Name: "alerts", Enabled: true, TopicARN: "arn:aws:sns:us-east-1:123456789012:alerts", NotifyPolicy: stringPointer("FIRST_ONLY")})
				return err
			},
		},
		"aws_sqs": {
			response: `{"id":1,"name":"alerts","enabled":true,"config":{"notify_policy":"FIRST_ONLY","aws_sqs.queue_url":"https://sqs.us-east-1.amazonaws.com/123456789012/alerts"}}`,
			create: func(client *Client) error {
				_, err := client.CreateAwsSqsOutlet(context.Background(), AwsSqsOutletCreateRequest{Name: "alerts", Enabled: true, QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/alerts", NotifyPolicy: stringPointer("FIRST_ONLY")})
				return err
			},
		},
		"gcp_pubsub": {
			response: `{"id":1,"name":"alerts","enabled":true,"config":{"notify_policy":"FIRST_ONLY","gcp_pubsub.project":"noozle","gcp_pubsub.topic":"alerts"}}`,
			create: func(client *Client) error {
				_, err := client.CreateGcpPubsubOutlet(context.Background(), GcpPubsubOutletCreateRequest{Name: "alerts", Enabled: true, Project: "noozle", Topic: "alerts", NotifyPolicy: stringPointer("FIRST_ONLY")})
				return err
			},
		},
		"http_server": {
			response: `{"id":1,"name":"alerts","enabled":true,"config":{"notify_policy":"FIRST_ONLY","http_server.path":"/alerts"}}`,
			create: func(client *Client) error {
				_, err := client.CreateHttpServerOutlet(context.Background(), HttpServerOutletCreateRequest{Name: "alerts", Enabled: true, Path: "/alerts", NotifyPolicy: stringPointer("FIRST_ONLY")})
				return err
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var payload typedOutletCreateRequest
			client := newTestClient(t, func(request *http.Request) (*http.Response, error) {
				parseMultipartPayload(t, request, &payload, nil)
				return jsonResponse(http.StatusCreated, test.response), nil
			})

			if err := test.create(client); err != nil {
				t.Fatalf("create outlet: %v", err)
			}
			if got := payload.Fields["notify_policy"]; got != "FIRST_ONLY" {
				t.Fatalf("notify_policy = %#v, want FIRST_ONLY", got)
			}
		})
	}
}

func TestBrokerOutletCreatesIncludeNotifyPolicy(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		path     string
		response string
		create   func(*Client) error
	}{
		"amqp_0_9": {
			path:     "/v1/terraform/amqp_0_9",
			response: `{"id":1,"name":"alerts","enabled":true,"config":{"notify_policy":"FIRST_ONLY","amqp_0_9.exchange":"alerts","amqp_0_9.key":"alerts.critical"}}`,
			create: func(client *Client) error {
				_, err := client.CreateAmqp09Outlet(context.Background(), Amqp09OutletCreateRequest{Name: "alerts", Enabled: true, Exchange: "alerts", Key: "alerts.critical", NotifyPolicy: stringPointer("FIRST_ONLY")})
				return err
			},
		},
		"amqp_1": {
			path:     "/v1/terraform/amqp_1",
			response: `{"id":1,"name":"alerts","enabled":true,"config":{"notify_policy":"FIRST_ONLY","amqp_1.target_address":"queue://alerts"}}`,
			create: func(client *Client) error {
				_, err := client.CreateAmqp1Outlet(context.Background(), Amqp1OutletCreateRequest{Name: "alerts", Enabled: true, TargetAddress: "queue://alerts", NotifyPolicy: stringPointer("FIRST_ONLY")})
				return err
			},
		},
		"kafka": {
			path:     "/v1/terraform/kafka",
			response: `{"id":1,"name":"alerts","enabled":true,"config":{"notify_policy":"FIRST_ONLY","kafka.addresses":["broker:9092"],"kafka.topic":"alerts"}}`,
			create: func(client *Client) error {
				_, err := client.CreateKafkaOutlet(context.Background(), KafkaOutletCreateRequest{Name: "alerts", Enabled: true, Addresses: []string{"broker:9092"}, Topic: "alerts", NotifyPolicy: stringPointer("FIRST_ONLY")})
				return err
			},
		},
		"nats_jetstream": {
			path:     "/v1/terraform/nats_jetstream",
			response: `{"id":1,"name":"alerts","enabled":true,"config":{"notify_policy":"FIRST_ONLY","nats_jetstream.subject":"alerts"}}`,
			create: func(client *Client) error {
				_, err := client.CreateNatsJetstreamOutlet(context.Background(), NatsJetstreamOutletCreateRequest{Name: "alerts", Enabled: true, Subject: "alerts", NotifyPolicy: stringPointer("FIRST_ONLY")})
				return err
			},
		},
		"pulsar": {
			path:     "/v1/terraform/pulsar",
			response: `{"id":1,"name":"alerts","enabled":true,"config":{"notify_policy":"FIRST_ONLY","pulsar.topic":"persistent://public/default/alerts","pulsar.url":"pulsar://broker:6650"}}`,
			create: func(client *Client) error {
				_, err := client.CreatePulsarOutlet(context.Background(), PulsarOutletCreateRequest{Name: "alerts", Enabled: true, Topic: "persistent://public/default/alerts", URL: "pulsar://broker:6650", NotifyPolicy: stringPointer("FIRST_ONLY")})
				return err
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var payload typedOutletCreateRequest
			client := newTestClient(t, func(request *http.Request) (*http.Response, error) {
				if request.URL.Path != test.path {
					t.Fatalf("path = %q, want %q", request.URL.Path, test.path)
				}
				parseMultipartPayload(t, request, &payload, nil)
				return jsonResponse(http.StatusCreated, test.response), nil
			})

			if err := test.create(client); err != nil {
				t.Fatalf("create outlet: %v", err)
			}
			if got := payload.Fields["notify_policy"]; got != "FIRST_ONLY" {
				t.Fatalf("notify_policy = %#v, want FIRST_ONLY", got)
			}
		})
	}
}

func TestStringFieldPointerPreservingEmpty(t *testing.T) {
	t.Parallel()

	empty, err := notifyPolicyPointerPreservingEmpty(map[string]any{"notify_policy": ""})
	if err != nil || empty == nil || *empty != "" {
		t.Fatalf("empty policy = %#v, %v; want non-nil empty string", empty, err)
	}
	absent, err := notifyPolicyPointerPreservingEmpty(map[string]any{})
	if err != nil || absent != nil {
		t.Fatalf("absent policy = %#v, %v; want nil", absent, err)
	}
}

func TestTypedOutletTransformUsesProviderSpecificRoutes(t *testing.T) {
	t.Parallel()

	var methods, paths []string
	client := newTestClient(t, func(request *http.Request) (*http.Response, error) {
		methods = append(methods, request.Method)
		paths = append(paths, request.URL.Path)
		if request.Method == http.MethodPut {
			var payload outletTransformSetRequest
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Fatalf("decode transform request: %v", err)
			}
			if payload.Transform != "summarize" {
				t.Fatalf("transform = %q, want summarize", payload.Transform)
			}
		}
		if request.Method == http.MethodDelete {
			return jsonResponse(http.StatusOK, `{"selected_transform":null,"effective_transform":"raw","state":"enabled","applied":true,"generation":5}`), nil
		}
		return jsonResponse(http.StatusOK, `{"selected_transform":"summarize","effective_transform":"summarize","state":"enabled","applied":true,"generation":4}`), nil
	})

	transform, err := client.GetTypedOutletTransform(context.Background(), "amqp_0_9", 41)
	if err != nil {
		t.Fatalf("get transform: %v", err)
	}
	if transform.SelectedTransform == nil || *transform.SelectedTransform != "summarize" || transform.Generation == nil || *transform.Generation != 4 || !transform.Applied {
		t.Fatalf("unexpected transform state: %+v", transform)
	}
	if _, err := client.SetTypedOutletTransform(context.Background(), "amqp_0_9", 41, "summarize"); err != nil {
		t.Fatalf("set transform: %v", err)
	}
	cleared, err := client.ClearTypedOutletTransform(context.Background(), "amqp_0_9", 41)
	if err != nil {
		t.Fatalf("clear transform: %v", err)
	}
	if cleared.SelectedTransform != nil || cleared.EffectiveTransform == nil || *cleared.EffectiveTransform != "raw" || cleared.Generation == nil || *cleared.Generation != 5 {
		t.Fatalf("unexpected cleared transform state: %+v", cleared)
	}

	if !slices.Equal(methods, []string{http.MethodGet, http.MethodPut, http.MethodDelete}) || !slices.Equal(paths, []string{"/v1/terraform/amqp_0_9/41/transform", "/v1/terraform/amqp_0_9/41/transform", "/v1/terraform/amqp_0_9/41/transform"}) {
		t.Fatalf("requests = %v %v", methods, paths)
	}
}

func TestTypedOutletTransformUsesAMQP1Routes(t *testing.T) {
	t.Parallel()

	var methods, paths []string
	client := newTestClient(t, func(request *http.Request) (*http.Response, error) {
		methods = append(methods, request.Method)
		paths = append(paths, request.URL.Path)
		if request.Method == http.MethodPut {
			var payload outletTransformSetRequest
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Fatalf("decode transform request: %v", err)
			}
			if payload.Transform != "summarize" {
				t.Fatalf("transform = %q, want summarize", payload.Transform)
			}
		}
		return jsonResponse(http.StatusOK, `{"selected_transform":"summarize","effective_transform":"summarize","state":"enabled","applied":true,"generation":4,"last_error":""}`), nil
	})

	if _, err := client.GetTypedOutletTransform(context.Background(), "amqp_1", 41); err != nil {
		t.Fatalf("get transform: %v", err)
	}
	if _, err := client.SetTypedOutletTransform(context.Background(), "amqp_1", 41, "summarize"); err != nil {
		t.Fatalf("set transform: %v", err)
	}
	if _, err := client.ClearTypedOutletTransform(context.Background(), "amqp_1", 41); err != nil {
		t.Fatalf("clear transform: %v", err)
	}

	path := "/v1/terraform/amqp_1/41/transform"
	if !slices.Equal(methods, []string{http.MethodGet, http.MethodPut, http.MethodDelete}) || !slices.Equal(paths, []string{path, path, path}) {
		t.Fatalf("requests = %v %v", methods, paths)
	}
}

func TestTypedOutletTransformUsesGCPAndHTTPServerRoutes(t *testing.T) {
	t.Parallel()

	for _, outletType := range []string{"gcp_pubsub", "http_server"} {
		t.Run(outletType, func(t *testing.T) {
			t.Parallel()

			var methods, paths []string
			client := newTestClient(t, func(request *http.Request) (*http.Response, error) {
				methods = append(methods, request.Method)
				paths = append(paths, request.URL.Path)
				if request.Method == http.MethodPut {
					var payload outletTransformSetRequest
					if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
						t.Fatalf("decode transform request: %v", err)
					}
					if payload.Transform != "summarize" {
						t.Fatalf("transform = %q, want summarize", payload.Transform)
					}
				}
				return jsonResponse(http.StatusOK, `{"selected_transform":"summarize","effective_transform":"summarize","state":"enabled","applied":true,"generation":4}`), nil
			})

			if _, err := client.GetTypedOutletTransform(context.Background(), outletType, 41); err != nil {
				t.Fatalf("get transform: %v", err)
			}
			if _, err := client.SetTypedOutletTransform(context.Background(), outletType, 41, "summarize"); err != nil {
				t.Fatalf("set transform: %v", err)
			}
			if _, err := client.ClearTypedOutletTransform(context.Background(), outletType, 41); err != nil {
				t.Fatalf("clear transform: %v", err)
			}

			path := "/v1/terraform/" + outletType + "/41/transform"
			if !slices.Equal(methods, []string{http.MethodGet, http.MethodPut, http.MethodDelete}) || !slices.Equal(paths, []string{path, path, path}) {
				t.Fatalf("requests = %v %v", methods, paths)
			}
		})
	}
}

func TestTypedOutletTransformUsesKafkaNATSAndPulsarRoutes(t *testing.T) {
	t.Parallel()

	for _, outletType := range []string{"kafka", "nats_jetstream", "pulsar"} {
		t.Run(outletType, func(t *testing.T) {
			t.Parallel()

			var methods, paths []string
			client := newTestClient(t, func(request *http.Request) (*http.Response, error) {
				methods = append(methods, request.Method)
				paths = append(paths, request.URL.Path)
				if request.Method == http.MethodPut {
					var payload outletTransformSetRequest
					if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
						t.Fatalf("decode transform request: %v", err)
					}
					if payload.Transform != "summarize" {
						t.Fatalf("transform = %q, want summarize", payload.Transform)
					}
				}
				return jsonResponse(http.StatusOK, `{"selected_transform":"summarize","effective_transform":"summarize","state":"enabled","applied":true,"generation":4}`), nil
			})

			if _, err := client.GetTypedOutletTransform(context.Background(), outletType, 41); err != nil {
				t.Fatalf("get transform: %v", err)
			}
			if _, err := client.SetTypedOutletTransform(context.Background(), outletType, 41, "summarize"); err != nil {
				t.Fatalf("set transform: %v", err)
			}
			if _, err := client.ClearTypedOutletTransform(context.Background(), outletType, 41); err != nil {
				t.Fatalf("clear transform: %v", err)
			}

			path := "/v1/terraform/" + outletType + "/41/transform"
			if !slices.Equal(methods, []string{http.MethodGet, http.MethodPut, http.MethodDelete}) || !slices.Equal(paths, []string{path, path, path}) {
				t.Fatalf("requests = %v %v", methods, paths)
			}
		})
	}
}

func stringPointer(value string) *string {
	return &value
}

func TestCreateAwsSnsOutletIncludesOptionalFields(t *testing.T) {
	t.Parallel()

	var payload typedOutletCreateRequest

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		parseMultipartPayload(t, r, &payload, nil)

		return jsonResponse(http.StatusCreated, `{
			"id": 201,
			"name": "SNS Alerts",
			"description": "desc",
			"enabled": true,
			"config": {
				"aws_sns.topic_arn": "arn:aws:sns:us-east-1:123456789012:noozle-alerts",
				"aws_sns.region": "us-east-1",
				"aws_sns.max_in_flight": 32,
				"aws_sns.metadata.exclude_prefixes": ["x-internal-"],
				"aws_sns.use_ec2_instance_role": true
			}
		}`), nil
	})

	region := "us-east-1"
	maxInFlight := int64(32)
	useEC2Role := true
	accessKeyID := "AKIA..."
	outlet, err := client.CreateAwsSnsOutlet(context.Background(), AwsSnsOutletCreateRequest{
		Name:                    "SNS Alerts",
		Description:             "desc",
		Enabled:                 true,
		TopicARN:                "arn:aws:sns:us-east-1:123456789012:noozle-alerts",
		Region:                  &region,
		MaxInFlight:             &maxInFlight,
		MetadataExcludePrefixes: []string{"x-internal-"},
		UseEC2InstanceRole:      &useEC2Role,
		AccessKeyID:             &accessKeyID,
	})
	if err != nil {
		t.Fatalf("create aws sns outlet: %v", err)
	}

	if payload.Fields["aws_sns.topic_arn"] != "arn:aws:sns:us-east-1:123456789012:noozle-alerts" {
		t.Fatalf("unexpected request fields: %+v", payload.Fields)
	}
	if payload.Fields["access_key_id"] != "AKIA..." {
		t.Fatalf("unexpected credential payload: %+v", payload.Fields)
	}
	if outlet.Region == nil || *outlet.Region != "us-east-1" {
		t.Fatalf("unexpected outlet response: %+v", outlet)
	}
}

func TestGetAwsSqsOutletDecodesTypedConfig(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{
			"id": 202,
			"name": "SQS Alerts",
			"description": "desc",
			"enabled": true,
			"config": {
				"aws_sqs.queue_url": "https://sqs.us-east-1.amazonaws.com/123456789012/noozle-alerts",
				"aws_sqs.batching.count": 10,
				"aws_sqs.batching.jitter": 0.25,
				"aws_sqs.metadata.exclude_prefixes": ["x-internal-"]
			}
		}`), nil
	})

	outlet, err := client.GetAwsSqsOutlet(context.Background(), 202)
	if err != nil {
		t.Fatalf("get aws sqs outlet: %v", err)
	}
	if outlet.QueueURL != "https://sqs.us-east-1.amazonaws.com/123456789012/noozle-alerts" {
		t.Fatalf("unexpected queue url: %+v", outlet)
	}
	if outlet.BatchingCount == nil || *outlet.BatchingCount != 10 {
		t.Fatalf("unexpected batching count: %+v", outlet)
	}
	if outlet.BatchingJitter == nil || *outlet.BatchingJitter != 0.25 {
		t.Fatalf("unexpected batching jitter: %+v", outlet)
	}
}

func TestCreateGcpPubsubOutletIncludesTypedFieldsAndCredentialFile(t *testing.T) {
	t.Parallel()

	var payload typedOutletCreateRequest
	fileParts := map[string]string{}

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		parseMultipartPayload(t, r, &payload, fileParts)

		return jsonResponse(http.StatusCreated, `{
			"id": 203,
			"name": "PubSub Alerts",
			"description": "desc",
			"enabled": true,
			"config": {
				"gcp_pubsub.project": "example-project",
				"gcp_pubsub.topic": "noozle-alerts",
				"gcp_pubsub.batching.jitter": 0.1
			}
		}`), nil
	})

	batchingJitter := 0.1
	filename := "service-account.json"
	content := "eyJ0eXBlIjoiZ2NwIn0="
	outlet, err := client.CreateGcpPubsubOutlet(context.Background(), GcpPubsubOutletCreateRequest{
		Name:                     "PubSub Alerts",
		Description:              "desc",
		Enabled:                  true,
		Project:                  "example-project",
		Topic:                    "noozle-alerts",
		BatchingJitter:           &batchingJitter,
		CredentialsFilename:      &filename,
		CredentialsContentBase64: &content,
	})
	if err != nil {
		t.Fatalf("create gcp pubsub outlet: %v", err)
	}

	if payload.Fields["gcp_pubsub.project"] != "example-project" || payload.Fields["gcp_pubsub.topic"] != "noozle-alerts" {
		t.Fatalf("unexpected request fields: %+v", payload.Fields)
	}
	fileRef, ok := payload.Fields["gcp_credentials_file"].(map[string]any)
	if !ok || fileRef["$file"] != "gcp_credentials_file_part" {
		t.Fatalf("unexpected credential file ref: %+v", payload.Fields["gcp_credentials_file"])
	}
	if fileParts["gcp_credentials_file_part"] == "" {
		t.Fatalf("expected multipart file content, got none")
	}
	if outlet.BatchingJitter == nil || *outlet.BatchingJitter != 0.1 {
		t.Fatalf("unexpected outlet response: %+v", outlet)
	}
}

func parseMultipartPayload(t *testing.T, r *http.Request, target any, fileParts map[string]string) {
	t.Helper()

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		t.Fatalf("parse media type: %v", err)
	}
	if mediaType != "multipart/form-data" {
		t.Fatalf("expected multipart/form-data, got %s", mediaType)
	}

	reader := multipart.NewReader(r.Body, params["boundary"])
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			return
		}
		if err != nil {
			t.Fatalf("read multipart part: %v", err)
		}

		body, err := io.ReadAll(part)
		if err != nil {
			t.Fatalf("read part body: %v", err)
		}

		if part.FormName() == "payload" {
			if err := json.Unmarshal(body, target); err != nil {
				t.Fatalf("decode payload: %v", err)
			}
			continue
		}

		if fileParts != nil {
			fileParts[part.FormName()] = string(body)
		}
	}
}

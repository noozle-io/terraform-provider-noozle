package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewRejectsInvalidBaseURL(t *testing.T) {
	t.Parallel()

	_, err := New("://bad", "secret")
	if err == nil {
		t.Fatal("expected invalid base URL error")
	}
}

func TestNewRequiresAPIKey(t *testing.T) {
	t.Parallel()

	_, err := New("https://example.com", "")
	if err == nil {
		t.Fatal("expected missing api key error")
	}
}

func TestQueryOperationsUseTenantContextRoutesAndHeaders(t *testing.T) {
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

		switch r.Method {
		case http.MethodPost:
			return jsonResponse(http.StatusCreated, queryJSONResponse), nil
		case http.MethodGet, http.MethodPut:
			return jsonResponse(http.StatusOK, queryJSONResponse), nil
		case http.MethodDelete:
			return jsonResponse(http.StatusNoContent, ""), nil
		default:
			t.Fatalf("unexpected method: %s", r.Method)
			return nil, nil
		}
	}, WithTenantContext(TenantContext{TenantID: tenantID, OrganizationID: organizationID}))

	if _, err := client.CreateQuery(context.Background(), QueryRequest{}); err != nil {
		t.Fatalf("create query: %v", err)
	}
	if _, err := client.GetQuery(context.Background(), 12); err != nil {
		t.Fatalf("get query: %v", err)
	}
	if _, err := client.UpdateQuery(context.Background(), 12, QueryRequest{}); err != nil {
		t.Fatalf("update query: %v", err)
	}
	if err := client.DeleteQuery(context.Background(), 12); err != nil {
		t.Fatalf("delete query: %v", err)
	}

	expected := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/v1/tenants/" + tenantID + "/queries"},
		{http.MethodGet, "/v1/tenants/" + tenantID + "/queries/12"},
		{http.MethodPut, "/v1/tenants/" + tenantID + "/queries/12"},
		{http.MethodDelete, "/v1/tenants/" + tenantID + "/queries/12"},
	}
	if len(requests) != len(expected) {
		t.Fatalf("expected %d requests, got %d", len(expected), len(requests))
	}
	for i, want := range expected {
		got := requests[i]
		if got.method != want.method || got.path != want.path {
			t.Errorf("request %d = %s %s, want %s %s", i, got.method, got.path, want.method, want.path)
		}
		if got.tenant != tenantID || got.org != organizationID {
			t.Errorf("request %d tenant context = tenant %q, organization %q", i, got.tenant, got.org)
		}
	}
}

func TestQueryOperationsRequireTenantContext(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		t.Fatal("query operation should not make an HTTP request without tenant context")
		return nil, nil
	}, WithTenantContext(TenantContext{}))

	_, err := client.CreateQuery(context.Background(), QueryRequest{})
	if !errors.Is(err, ErrMissingTenantContext) {
		t.Fatalf("expected missing tenant context error, got %v", err)
	}
}

func TestOutletQueryLinkOperationsUseCanonicalTenantRoutes(t *testing.T) {
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

		switch r.Method {
		case http.MethodPost, http.MethodDelete:
			return jsonResponse(http.StatusOK, `{"status":"linked"}`), nil
		case http.MethodGet:
			return jsonResponse(http.StatusOK, `{"id":99,"queries":[42]}`), nil
		default:
			t.Fatalf("unexpected method: %s", r.Method)
			return nil, nil
		}
	}, WithTenantContext(TenantContext{TenantID: tenantID, OrganizationID: organizationID}))

	if err := client.LinkOutletQuery(context.Background(), 99, 42); err != nil {
		t.Fatalf("link outlet query: %v", err)
	}
	if _, err := client.GetOutlet(context.Background(), 99); err != nil {
		t.Fatalf("get outlet: %v", err)
	}
	if err := client.UnlinkOutletQuery(context.Background(), 99, 42); err != nil {
		t.Fatalf("unlink outlet query: %v", err)
	}

	expected := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/v1/tenants/" + tenantID + "/outlets/99/queries"},
		{http.MethodGet, "/v1/tenants/" + tenantID + "/outlets/99"},
		{http.MethodDelete, "/v1/tenants/" + tenantID + "/outlets/99/queries"},
	}
	if len(requests) != len(expected) {
		t.Fatalf("expected %d requests, got %d", len(expected), len(requests))
	}
	for i, want := range expected {
		got := requests[i]
		if got.method != want.method || got.path != want.path {
			t.Errorf("request %d = %s %s, want %s %s", i, got.method, got.path, want.method, want.path)
		}
		if got.tenant != tenantID || got.org != organizationID {
			t.Errorf("request %d tenant context = tenant %q, organization %q", i, got.tenant, got.org)
		}
	}
}

func TestOutletQueryLinkOperationsRequireTenantContext(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		t.Fatal("outlet-query link operation should not make an HTTP request without tenant context")
		return nil, nil
	}, WithTenantContext(TenantContext{}))

	if err := client.LinkOutletQuery(context.Background(), 99, 42); !errors.Is(err, ErrMissingTenantContext) {
		t.Fatalf("expected missing tenant context error, got %v", err)
	}
}

func TestGetTenantOutletClassifiesRemoteAbsence(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		if got, want := r.URL.Path, "/v1/tenants/t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q/outlets/99"; got != want {
			t.Fatalf("tenant outlet path = %q, want %q", got, want)
		}
		return jsonResponse(http.StatusNotFound, `{"error":"Resource not found","type":"not_found"}`), nil
	})

	_, err := client.GetOutlet(context.Background(), 99)
	if !IsNotFound(err) {
		t.Fatalf("expected remote absence to be not found, got %v", err)
	}
}

const queryJSONResponse = `{
	"id": 12,
	"name": "Breaking News",
	"description": "desc",
	"full_expression": "feed:finance",
	"tags": ["finance"],
	"is_active": true,
	"created_at": "2024-04-22T09:15:00Z",
	"updated_at": "2024-04-23T10:05:00Z"
}`

func TestCreateQueryBuildsRequestAndDecodesResponse(t *testing.T) {
	t.Parallel()

	var receivedMethod string
	var receivedPath string
	var receivedKey string
	var receivedContentType string
	var requestBody QueryRequest

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		receivedKey = r.Header.Get("X-API-Key")
		receivedContentType = r.Header.Get("Content-Type")

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Fatalf("decode request: %v", err)
		}

		return jsonResponse(http.StatusCreated, `{
			"id": 12,
			"name": "Breaking News",
			"description": "desc",
			"full_expression": "feed:finance",
			"tags": ["finance"],
			"is_active": true,
			"created_at": "2024-04-22T09:15:00Z",
			"updated_at": "2024-04-23T10:05:00Z"
		}`), nil
	})

	query, err := client.CreateQuery(context.Background(), QueryRequest{
		Name:           "Breaking News",
		Description:    "desc",
		FullExpression: "feed:finance",
		Tags:           []string{"finance"},
		IsActive:       true,
	})
	if err != nil {
		t.Fatalf("create query: %v", err)
	}

	if receivedMethod != http.MethodPost {
		t.Fatalf("expected POST, got %s", receivedMethod)
	}
	if receivedPath != "/v1/tenants/t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q/queries" {
		t.Fatalf("expected canonical tenant query path, got %s", receivedPath)
	}
	if receivedKey != "nk2_secret" {
		t.Fatalf("expected api key header, got %q", receivedKey)
	}
	if receivedContentType != "application/json" {
		t.Fatalf("expected application/json, got %q", receivedContentType)
	}
	if requestBody.FullExpression != "feed:finance" {
		t.Fatalf("unexpected request body: %+v", requestBody)
	}
	if query.ID != 12 || query.Name != "Breaking News" || query.CreatedAt.IsZero() {
		t.Fatalf("unexpected query response: %+v", query)
	}
}

func TestGetQueryDecodesDetailResponse(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{
			"id": 12,
			"name": "Breaking News",
			"description": "desc",
			"full_expression": "feed:finance",
			"tags": ["finance"],
			"is_active": true,
			"outlets": [],
			"created_at": "2024-04-22T09:15:00Z",
			"updated_at": "2024-04-23T10:05:00Z",
			"active_trace_session_ids": ["trace-1"]
		}`), nil
	})

	query, err := client.GetQuery(context.Background(), 12)
	if err != nil {
		t.Fatalf("get query: %v", err)
	}

	if query.ID != 12 || query.UpdatedAt.IsZero() {
		t.Fatalf("unexpected query: %+v", query)
	}
}

func TestDeleteQueryUsesNoContent(t *testing.T) {
	t.Parallel()

	var receivedMethod string
	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		receivedMethod = r.Method
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})

	if err := client.DeleteQuery(context.Background(), 12); err != nil {
		t.Fatalf("delete query: %v", err)
	}
	if receivedMethod != http.MethodDelete {
		t.Fatalf("expected DELETE, got %s", receivedMethod)
	}
}

func TestLinkOutletQueryBuildsSingleElementBody(t *testing.T) {
	t.Parallel()

	var requestBody struct {
		QueryIDs []int64 `json:"query_ids"`
	}

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		return jsonResponse(http.StatusOK, `{"status":"linked"}`), nil
	})

	if err := client.LinkOutletQuery(context.Background(), 99, 42); err != nil {
		t.Fatalf("link outlet query: %v", err)
	}
	if len(requestBody.QueryIDs) != 1 || requestBody.QueryIDs[0] != 42 {
		t.Fatalf("unexpected link request body: %+v", requestBody)
	}
}

func TestGetOutletDecodesQueryIDs(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{
			"id": 99,
			"user_id": 456,
			"name": "Breaking Alerts",
			"description": "desc",
			"outlet_type": "ntfy",
			"enabled": true,
			"queries": [12, 42]
		}`), nil
	})

	outlet, err := client.GetOutlet(context.Background(), 99)
	if err != nil {
		t.Fatalf("get outlet: %v", err)
	}

	if outlet.ID != 99 || len(outlet.QueryIDs) != 2 || outlet.QueryIDs[1] != 42 {
		t.Fatalf("unexpected outlet: %+v", outlet)
	}
}

func TestAPIErrorDecodingPreservesMetadata(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusNotFound, `{
			"error": "Resource not found",
			"type": "not_found",
			"correlation_id": "corr-123",
			"details": {"resource":"query","id":12}
		}`), nil
	})

	_, err := client.GetQuery(context.Background(), 12)
	if err == nil {
		t.Fatal("expected not found error")
	}
	if !IsNotFound(err) {
		t.Fatalf("expected not found classification, got %v", err)
	}

	var apiErr *APIError
	if !errorsAs(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.CorrelationID != "corr-123" {
		t.Fatalf("unexpected correlation id: %+v", apiErr)
	}
	if apiErr.Details["resource"] != "query" {
		t.Fatalf("unexpected details: %+v", apiErr.Details)
	}
}

func TestRetryGetOnServerError(t *testing.T) {
	t.Parallel()

	var requests atomic.Int32
	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		current := requests.Add(1)
		if current == 1 {
			return jsonResponse(http.StatusInternalServerError, `{"error":"temporary","type":"internal"}`), nil
		}

		return jsonResponse(http.StatusOK, `{
			"id": 12,
			"name": "Breaking News",
			"description": "desc",
			"full_expression": "feed:finance",
			"tags": [],
			"is_active": true,
			"outlets": [],
			"created_at": "2024-04-22T09:15:00Z",
			"updated_at": "2024-04-23T10:05:00Z",
			"active_trace_session_ids": []
		}`), nil
	}, WithRetryMax(1))

	if _, err := client.GetQuery(context.Background(), 12); err != nil {
		t.Fatalf("get query after retry: %v", err)
	}
	if requests.Load() != 2 {
		t.Fatalf("expected 2 requests, got %d", requests.Load())
	}
}

func TestDoNotRetryMutatingRequests(t *testing.T) {
	t.Parallel()

	var requests atomic.Int32
	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		requests.Add(1)
		return jsonResponse(http.StatusInternalServerError, `{"error":"temporary","type":"internal"}`), nil
	}, WithRetryMax(3))

	_, err := client.CreateQuery(context.Background(), QueryRequest{
		Name:           "Breaking News",
		FullExpression: "feed:finance",
	})
	if err == nil {
		t.Fatal("expected create query error")
	}
	if requests.Load() != 1 {
		t.Fatalf("expected 1 request, got %d", requests.Load())
	}
}

func TestDebugLogsRedactRequestAndResponsePayloads(t *testing.T) {
	t.Parallel()

	logger := &recordingLogger{}

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusCreated, `{
			"id": 12,
			"name": "Breaking News",
			"description": "desc",
			"full_expression": "feed:finance",
			"tags": ["finance"],
			"is_active": true,
			"created_at": "2024-04-22T09:15:00Z",
			"updated_at": "2024-04-23T10:05:00Z"
		}`), nil
	}, WithLogger(logger))

	_, err := client.CreateQuery(context.Background(), QueryRequest{
		Name:           "Breaking News",
		Description:    "desc",
		FullExpression: "feed:finance",
		Tags:           []string{"finance"},
		IsActive:       true,
	})
	if err != nil {
		t.Fatalf("create query: %v", err)
	}

	if len(logger.entries) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(logger.entries))
	}

	requestEntry := logger.entries[0]
	responseEntry := logger.entries[1]

	if requestEntry.message != "noozle client request" || requestEntry.fields["payload_bytes"] == 0 {
		t.Fatalf("unexpected request log entry: %+v", requestEntry)
	}
	if _, found := requestEntry.fields["payload"]; found {
		t.Fatalf("request log contains its payload: %+v", requestEntry)
	}
	if responseEntry.message != "noozle client response" || responseEntry.fields["payload_bytes"] == 0 {
		t.Fatalf("unexpected response log entry: %+v", responseEntry)
	}
	if _, found := responseEntry.fields["payload"]; found {
		t.Fatalf("response log contains its payload: %+v", responseEntry)
	}
}

func TestDebugLogsRedactErrorResponsePayloads(t *testing.T) {
	t.Parallel()

	logger := &recordingLogger{}

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusBadRequest, `{
			"error": "Invalid outlet fields",
			"type": "bad_request",
			"correlation_id": "corr-400"
		}`), nil
	}, WithLogger(logger))

	_, err := client.CreateQuery(context.Background(), QueryRequest{
		Name:           "Breaking News",
		FullExpression: "feed:finance",
	})
	if err == nil {
		t.Fatal("expected create query error")
	}

	if len(logger.entries) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(logger.entries))
	}

	responseEntry := logger.entries[1]
	if responseEntry.message != "noozle client response" || responseEntry.fields["status_code"] != http.StatusBadRequest || responseEntry.fields["payload_bytes"] == 0 {
		t.Fatalf("unexpected error response log entry: %+v", responseEntry)
	}
	if _, found := responseEntry.fields["payload"]; found {
		t.Fatalf("response log contains its payload: %+v", responseEntry)
	}
}

// errorsAs keeps the tests local to stdlib behavior without shadowing errors.As in assertions.
func errorsAs(err error, target any) bool {
	return errors.As(err, target)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := f(r)
	if resp != nil && resp.Request == nil {
		resp.Request = r
	}

	return resp, err
}

func newTestClient(t *testing.T, fn roundTripFunc, opts ...Option) *Client {
	t.Helper()

	httpClient := &http.Client{
		Timeout:   2 * time.Second,
		Transport: fn,
	}

	clientOpts := append([]Option{
		WithHTTPClient(httpClient),
		WithTenantContext(TenantContext{TenantID: "t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q"}),
	}, opts...)
	client, err := New("https://example.test/", "nk2_secret", clientOpts...)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	return client
}

func jsonResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
}

type recordingLogger struct {
	entries []logEntry
}

type logEntry struct {
	message string
	fields  map[string]any
}

func (l *recordingLogger) Debug(_ context.Context, message string, fields map[string]any) {
	copiedFields := make(map[string]any, len(fields))
	for key, value := range fields {
		copiedFields[key] = value
	}

	l.entries = append(l.entries, logEntry{
		message: message,
		fields:  copiedFields,
	})
}

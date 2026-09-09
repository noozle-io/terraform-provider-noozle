package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestCreateHttpServerOutletIncludesTypedFields(t *testing.T) {
	t.Parallel()

	var contentType string
	var requestPayload string

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		contentType = r.Header.Get("Content-Type")
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		requestPayload = string(payload)

		return jsonResponse(http.StatusCreated, `{
			"id": 91,
			"name": "HTTP Server",
			"description": "desc",
			"enabled": true,
			"config": {
				"http_server.path": "/api/events",
				"http_server.allowed_verbs": ["POST"]
			},
			"apply_status": "healthy",
			"last_apply_error": "",
			"config_generation": 1,
			"applied_generation": 1
		}`), nil
	})

	outlet, err := client.CreateHttpServerOutlet(context.Background(), HttpServerOutletCreateRequest{
		Name:         "HTTP Server",
		Description:  "desc",
		Enabled:      true,
		Path:         "/api/events",
		AllowedVerbs: []string{"POST"},
	})
	if err != nil {
		t.Fatalf("create http_server outlet: %v", err)
	}

	if outlet.Path != "/api/events" {
		t.Fatalf("unexpected outlet: %+v", outlet)
	}
	if !strings.Contains(contentType, "multipart/form-data") {
		t.Fatalf("expected multipart content type, got %q", contentType)
	}
	for _, expected := range []string{`"http_server.path":"/api/events"`, `"http_server.allowed_verbs":["POST"]`} {
		if !strings.Contains(requestPayload, expected) {
			t.Fatalf("expected payload to contain %q, got %q", expected, requestPayload)
		}
	}
}

func TestConfigureHttpServerTLSFeatureBuildsMultipartRequest(t *testing.T) {
	t.Parallel()

	var contentType string
	var requestPayload string

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		contentType = r.Header.Get("Content-Type")
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		requestPayload = string(payload)

		return jsonResponse(http.StatusOK, `{"feature":"http_server.tls","state":"enabled","applied":true}`), nil
	})

	err := client.ConfigureHttpServerTLSFeature(context.Background(), 91, HttpServerTLSFeatureRequest{
		ServerCertFilename:      "cert.pem",
		ServerCertContentBase64: "Y2VydA==",
		ServerKeyFilename:       "key.pem",
		ServerKeyContentBase64:  "a2V5",
	})
	if err != nil {
		t.Fatalf("configure http_server tls feature: %v", err)
	}

	if !strings.Contains(contentType, "multipart/form-data") {
		t.Fatalf("expected multipart content type, got %q", contentType)
	}
	for _, expected := range []string{`"server_cert":{"$file":"server_cert_part"}`, `"server_key":{"$file":"server_key_part"}`, `filename="cert.pem"`, `filename="key.pem"`} {
		if !strings.Contains(requestPayload, expected) {
			t.Fatalf("expected payload to contain %q, got %q", expected, requestPayload)
		}
	}
}

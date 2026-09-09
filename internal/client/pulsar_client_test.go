package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestCreatePulsarOutletIncludesTypedFields(t *testing.T) {
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
			"id": 321,
			"name": "Pulsar",
			"description": "desc",
			"enabled": true,
			"config": {
				"pulsar.topic": "persistent://public/default/noozle-alerts",
				"pulsar.url": "pulsar://pulsar.example.com:6650",
				"pulsar.max_in_flight": 64
			},
			"apply_status": "healthy",
			"last_apply_error": "",
			"config_generation": 2,
			"applied_generation": 2
		}`), nil
	})

	outlet, err := client.CreatePulsarOutlet(context.Background(), PulsarOutletCreateRequest{
		Name:        "Pulsar",
		Description: "desc",
		Enabled:     true,
		MaxInFlight: int64Ptr(64),
		Topic:       "persistent://public/default/noozle-alerts",
		URL:         "pulsar://pulsar.example.com:6650",
	})
	if err != nil {
		t.Fatalf("create pulsar outlet: %v", err)
	}

	if outlet.Topic != "persistent://public/default/noozle-alerts" || outlet.URL != "pulsar://pulsar.example.com:6650" {
		t.Fatalf("unexpected outlet: %+v", outlet)
	}
	if !strings.Contains(contentType, "multipart/form-data") {
		t.Fatalf("expected multipart content type, got %q", contentType)
	}
	for _, expected := range []string{
		`"pulsar.topic":"persistent://public/default/noozle-alerts"`,
		`"pulsar.url":"pulsar://pulsar.example.com:6650"`,
		`"pulsar.max_in_flight":64`,
	} {
		if !strings.Contains(requestPayload, expected) {
			t.Fatalf("expected payload to contain %q, got %q", expected, requestPayload)
		}
	}
}

func TestConfigurePulsarOAuth2FeatureBuildsMultipartRequest(t *testing.T) {
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

		return jsonResponse(http.StatusOK, `{"feature":"pulsar.oauth2","state":"enabled","applied":true}`), nil
	})

	err := client.ConfigurePulsarOAuth2Feature(context.Background(), 321, PulsarOAuth2FeatureRequest{
		IssuerURL:               "https://issuer.example.com",
		Audience:                "pulsar",
		PrivateKeyFilename:      "oauth2-private-key.json",
		PrivateKeyContentBase64: "e30=",
	})
	if err != nil {
		t.Fatalf("configure pulsar oauth2 feature: %v", err)
	}

	if !strings.Contains(contentType, "multipart/form-data") {
		t.Fatalf("expected multipart content type, got %q", contentType)
	}
	for _, expected := range []string{
		`"issuer_url":"https://issuer.example.com"`,
		`"audience":"pulsar"`,
		`"oauth2_private_key":{"$file":"oauth2_private_key_part"}`,
		`filename="oauth2-private-key.json"`,
	} {
		if !strings.Contains(requestPayload, expected) {
			t.Fatalf("expected payload to contain %q, got %q", expected, requestPayload)
		}
	}
}

func int64Ptr(value int64) *int64 {
	return &value
}

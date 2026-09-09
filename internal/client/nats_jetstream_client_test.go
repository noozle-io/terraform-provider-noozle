package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestCreateNatsJetstreamOutletIncludesTypedFields(t *testing.T) {
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
			"id": 432,
			"name": "NATS",
			"description": "desc",
			"enabled": true,
			"config": {
				"nats_jetstream.subject": "noozle.alerts",
				"nats_jetstream.max_in_flight": 64,
				"nats_jetstream.headers": {
					"X-Tenant": "${! meta(\"tenant\") }"
				}
			},
			"apply_status": "healthy",
			"last_apply_error": "",
			"config_generation": 2,
			"applied_generation": 2
		}`), nil
	})

	outlet, err := client.CreateNatsJetstreamOutlet(context.Background(), NatsJetstreamOutletCreateRequest{
		Name:        "NATS",
		Description: "desc",
		Enabled:     true,
		Headers: map[string]string{
			"X-Tenant": "${! meta(\"tenant\") }",
		},
		MaxInFlight: int64Ptr(64),
		Subject:     "noozle.alerts",
		URLs:        []string{"nats://nats.example.com:4222"},
	})
	if err != nil {
		t.Fatalf("create nats jetstream outlet: %v", err)
	}

	if outlet.Subject != "noozle.alerts" {
		t.Fatalf("unexpected outlet: %+v", outlet)
	}
	if !strings.Contains(contentType, "multipart/form-data") {
		t.Fatalf("expected multipart content type, got %q", contentType)
	}
	for _, expected := range []string{
		`"nats_jetstream.subject":"noozle.alerts"`,
		`"nats_jetstream.urls":["nats://nats.example.com:4222"]`,
		`"nats_jetstream.max_in_flight":64`,
		`"nats_jetstream.headers":{"X-Tenant":"${! meta(\"tenant\") }"}`,
	} {
		if !strings.Contains(requestPayload, expected) {
			t.Fatalf("expected payload to contain %q, got %q", expected, requestPayload)
		}
	}
}

func TestConfigureNatsJetstreamNkeyFeatureBuildsMultipartRequest(t *testing.T) {
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

		return jsonResponse(http.StatusOK, `{"feature":"nats.auth.nkey","state":"enabled","applied":true}`), nil
	})

	err := client.ConfigureNatsJetstreamNkeyFeature(context.Background(), 432, NatsJetstreamNkeyFeatureRequest{
		Filename:      "nkey.creds",
		ContentBase64: "bmF0cw==",
	})
	if err != nil {
		t.Fatalf("configure nats jetstream nkey feature: %v", err)
	}

	if !strings.Contains(contentType, "multipart/form-data") {
		t.Fatalf("expected multipart content type, got %q", contentType)
	}
	for _, expected := range []string{
		`"nkey_file":{"$file":"nkey_file_part"}`,
		`filename="nkey.creds"`,
	} {
		if !strings.Contains(requestPayload, expected) {
			t.Fatalf("expected payload to contain %q, got %q", expected, requestPayload)
		}
	}
}

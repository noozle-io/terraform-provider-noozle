package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestCreateAmqp09OutletIncludesURLsAndTypedFields(t *testing.T) {
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
			"id": 401,
			"name": "AMQP 0.9.1",
			"description": "desc",
			"enabled": true,
			"config": {
				"amqp_0_9.exchange": "noozle.alerts",
				"amqp_0_9.key": "articles.match",
				"amqp_0_9.max_in_flight": 32
			},
			"apply_status": "healthy",
			"last_apply_error": "",
			"config_generation": 1,
			"applied_generation": 1
		}`), nil
	})

	_, err := client.CreateAmqp09Outlet(context.Background(), Amqp09OutletCreateRequest{
		Name:        "AMQP 0.9.1",
		Description: "desc",
		Enabled:     true,
		Exchange:    "noozle.alerts",
		Key:         "articles.match",
		MaxInFlight: int64Ptr(32),
		URLs:        []string{"amqps://broker.example.com:5671"},
	})
	if err != nil {
		t.Fatalf("create amqp_0_9 outlet: %v", err)
	}

	if !strings.Contains(contentType, "multipart/form-data") {
		t.Fatalf("expected multipart content type, got %q", contentType)
	}
	for _, expected := range []string{
		`"amqp_0_9.exchange":"noozle.alerts"`,
		`"amqp_0_9.key":"articles.match"`,
		`"amqp_0_9.urls":["amqps://broker.example.com:5671"]`,
	} {
		if !strings.Contains(requestPayload, expected) {
			t.Fatalf("expected payload to contain %q, got %q", expected, requestPayload)
		}
	}
}

func TestCreateAmqp1OutletIncludesURLsAndTypedFields(t *testing.T) {
	t.Parallel()

	var requestPayload string

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		requestPayload = string(payload)

		return jsonResponse(http.StatusCreated, `{
			"id": 402,
			"name": "AMQP 1.0",
			"description": "desc",
			"enabled": true,
			"config": {
				"amqp_1.target_address": "alerts/out",
				"amqp_1.max_in_flight": 16
			},
			"apply_status": "healthy",
			"last_apply_error": "",
			"config_generation": 1,
			"applied_generation": 1
		}`), nil
	})

	_, err := client.CreateAmqp1Outlet(context.Background(), Amqp1OutletCreateRequest{
		Name:          "AMQP 1.0",
		Description:   "desc",
		Enabled:       true,
		TargetAddress: "alerts/out",
		MaxInFlight:   int64Ptr(16),
		URLs:          []string{"amqps://broker.example.com:5671"},
	})
	if err != nil {
		t.Fatalf("create amqp_1 outlet: %v", err)
	}

	for _, expected := range []string{
		`"amqp_1.target_address":"alerts/out"`,
		`"amqp_1.max_in_flight":16`,
		`"amqp_1.urls":["amqps://broker.example.com:5671"]`,
	} {
		if !strings.Contains(requestPayload, expected) {
			t.Fatalf("expected payload to contain %q, got %q", expected, requestPayload)
		}
	}
}

func TestConfigureAmqp1SASLPlainFeatureBuildsMultipartRequest(t *testing.T) {
	t.Parallel()

	var requestPayload string

	client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		requestPayload = string(payload)
		return jsonResponse(http.StatusOK, `{"feature":"amqp_1.sasl.plain","state":"enabled","applied":true}`), nil
	})

	err := client.ConfigureAmqp1SASLPlainFeature(context.Background(), 402, Amqp1SASLPlainFeatureRequest{
		User:     "alerts-producer",
		Password: "super-secret",
	})
	if err != nil {
		t.Fatalf("configure amqp_1 sasl plain feature: %v", err)
	}

	for _, expected := range []string{
		`"user":"alerts-producer"`,
		`"sasl_password":"super-secret"`,
	} {
		if !strings.Contains(requestPayload, expected) {
			t.Fatalf("expected payload to contain %q, got %q", expected, requestPayload)
		}
	}
}

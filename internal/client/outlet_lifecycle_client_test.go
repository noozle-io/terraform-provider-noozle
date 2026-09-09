package client

import (
	"context"
	"net/http"
	"testing"
)

func TestGetTypedOutletStatusObservesRequestedGeneration(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", request.Method)
		}
		if request.URL.Path != "/v1/terraform/ntfy/99/status" {
			t.Fatalf("path = %s", request.URL.Path)
		}
		if request.URL.Query().Get("generation") != "7" {
			t.Fatalf("generation = %q, want 7", request.URL.Query().Get("generation"))
		}
		if prefer := request.Header.Get("Prefer"); prefer != "" {
			t.Fatalf("Prefer header = %q, want absent", prefer)
		}
		return jsonResponse(http.StatusOK, `{
			"outlet_id": 99,
			"tenant_id": "t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q",
			"desired_status": "enabled",
			"requested_generation": 7,
			"desired_generation": 7,
			"superseded": false,
			"runtime_status": "active",
			"realization": {
				"active_generation": 7,
				"realization_attempt": 1,
				"phase": "running",
				"ready": true,
				"degraded": false,
				"last_error": null
			},
			"checked_at": "2026-09-02T12:00:00Z"
		}`), nil
	}, WithTenantContext(TenantContext{TenantID: "t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q"}))

	status, err := client.GetTypedOutletStatus(context.Background(), "ntfy", 99, 7)
	if err != nil {
		t.Fatalf("get outlet status: %v", err)
	}
	if status.Superseded || status.RequestedGeneration != 7 || status.RuntimeStatus != "active" || !status.Realization.Ready {
		t.Fatalf("unexpected outlet status: %+v", status)
	}
}

func TestTypedOutletDeleteOperationsDecodeGeneration(t *testing.T) {
	t.Parallel()

	for name, statusCode := range map[string]int{"ok": http.StatusOK, "accepted": http.StatusAccepted, "no_content": http.StatusNoContent} {
		t.Run(name, func(t *testing.T) {
			client := newTestClient(t, func(r *http.Request) (*http.Response, error) {
				if r.Method != http.MethodDelete {
					t.Fatalf("method = %s, want DELETE", r.Method)
				}
				if statusCode == http.StatusNoContent {
					return jsonResponse(statusCode, ""), nil
				}
				return jsonResponse(statusCode, `{
					"outlet_id": 99,
					"tenant_id": "t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q",
					"desired_status": "deleting",
					"operation": {
						"kind": "delete",
						"desired_generation": 12,
						"status_link": "/v1/terraform/ntfy/99/status?generation=12",
						"realization": {"active_generation": 11, "realization_attempt": 1, "phase": "draining", "ready": false, "degraded": false, "last_error": null}
					}
				}`), nil
			})
			deleteOperations := map[string]func(context.Context, int64) (OutletOperation, error){
				"amqp_0_9":       client.DeleteAmqp09OutletOperation,
				"amqp_1":         client.DeleteAmqp1OutletOperation,
				"aws_sns":        client.DeleteAwsSnsOutletOperation,
				"aws_sqs":        client.DeleteAwsSqsOutletOperation,
				"gcp_pubsub":     client.DeleteGcpPubsubOutletOperation,
				"http_server":    client.DeleteHttpServerOutletOperation,
				"kafka":          client.DeleteKafkaOutletOperation,
				"nats_jetstream": client.DeleteNatsJetstreamOutletOperation,
				"ntfy":           client.DeleteNtfyOutletOperation,
				"pulsar":         client.DeletePulsarOutletOperation,
			}
			wantGeneration := int64(12)
			if statusCode == http.StatusNoContent {
				wantGeneration = 0
			}
			for outletType, deleteOperation := range deleteOperations {
				t.Run(outletType, func(t *testing.T) {
					operation, err := deleteOperation(context.Background(), 99)
					if err != nil {
						t.Fatalf("delete outlet: %v", err)
					}
					if operation.DesiredGeneration != wantGeneration {
						t.Fatalf("operation = %#v, want generation %d", operation, wantGeneration)
					}
					if wantGeneration != 0 && operation.Kind != "delete" {
						t.Fatalf("operation = %#v, want delete operation", operation)
					}
				})
			}
		})
	}
}

func TestOutletOperationRecorderKeepsFinalMutation(t *testing.T) {
	t.Parallel()

	responses := []string{
		`{"operation":{"kind":"update","desired_generation":5,"status_link":"/status?generation=5","realization":{"active_generation":4,"realization_attempt":1,"phase":"starting","ready":false,"degraded":false,"last_error":null}}}`,
		`{"operation":{"kind":"update","desired_generation":6,"status_link":"/status?generation=6","realization":{"active_generation":4,"realization_attempt":2,"phase":"starting","ready":false,"degraded":false,"last_error":null}}}`,
		`{"operation":{"kind":"refresh","desired_generation":7,"status_link":"/status?generation=7","realization":{"active_generation":6,"realization_attempt":3,"phase":"running","ready":true,"degraded":false,"last_error":null}}}`,
	}
	client := newTestClient(t, func(_ *http.Request) (*http.Response, error) {
		response := responses[0]
		responses = responses[1:]
		return jsonResponse(http.StatusOK, response), nil
	}, WithTenantContext(TenantContext{TenantID: "t_01JZ8V5Q5B4Z3T2Y1X0W9V8R7Q"}))

	ctx, recorder := WithOutletOperationRecorder(context.Background())
	if err := client.doJSON(ctx, http.MethodPatch, "/v1/terraform/kafka/99/fields/first", nil, nil, statusOK); err != nil {
		t.Fatalf("first mutation: %v", err)
	}
	if err := client.doJSON(ctx, http.MethodPatch, "/v1/terraform/kafka/99/fields/final", nil, nil, statusOK); err != nil {
		t.Fatalf("final mutation: %v", err)
	}
	if err := client.doJSON(ctx, http.MethodGet, "/v1/terraform/kafka/99", nil, nil, statusOK); err != nil {
		t.Fatalf("refresh: %v", err)
	}

	operation, ok := recorder.FinalOperation()
	if !ok || operation.DesiredGeneration != 6 {
		t.Fatalf("final operation = %#v, %t; want generation 6", operation, ok)
	}
}

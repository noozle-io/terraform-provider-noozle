package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-noozle/internal/client"
)

func TestAwsSnsOutletModelFromAPI(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	maxInFlight := int64(32)
	useEC2Role := true

	state := awsSnsOutletModelFromAPI(client.AwsSnsOutlet{
		ID:                      201,
		Name:                    "SNS Alerts",
		Description:             "desc",
		Enabled:                 true,
		Region:                  &region,
		MaxInFlight:             &maxInFlight,
		MetadataExcludePrefixes: []string{"x-internal-"},
		TopicARN:                "arn:aws:sns:us-east-1:123456789012:noozle-alerts",
		UseEC2InstanceRole:      &useEC2Role,
	}, types.ObjectNull(awsCredentialsAttributeTypes))

	if state.Region.ValueString() != "us-east-1" {
		t.Fatalf("expected region to be set, got %s", state.Region)
	}
	if state.MaxInFlight.ValueInt64() != 32 {
		t.Fatalf("expected max_in_flight to be set, got %d", state.MaxInFlight.ValueInt64())
	}
}

func TestAwsSqsOutletModelFromAPI(t *testing.T) {
	t.Parallel()

	batchingJitter := 0.25

	state := awsSqsOutletModelFromAPI(client.AwsSqsOutlet{
		ID:                      202,
		Name:                    "SQS Alerts",
		Description:             "desc",
		Enabled:                 true,
		QueueURL:                "https://sqs.us-east-1.amazonaws.com/123456789012/noozle-alerts",
		BatchingJitter:          &batchingJitter,
		MetadataExcludePrefixes: []string{"x-internal-"},
	}, types.ObjectNull(awsCredentialsAttributeTypes))

	if state.QueueURL.ValueString() == "" {
		t.Fatal("expected queue_url to be set")
	}
	if state.BatchingJitter.ValueFloat64() != 0.25 {
		t.Fatalf("expected batching_jitter to be set, got %f", state.BatchingJitter.ValueFloat64())
	}
}

func TestGcpPubsubOutletModelFromAPI(t *testing.T) {
	t.Parallel()

	batchingCount := int64(100)

	state := gcpPubsubOutletModelFromAPI(client.GcpPubsubOutlet{
		ID:                      203,
		Name:                    "PubSub Alerts",
		Description:             "desc",
		Enabled:                 true,
		Project:                 "example-project",
		Topic:                   "noozle-alerts",
		BatchingCount:           &batchingCount,
		MetadataExcludePrefixes: []string{"x-internal-"},
	}, types.ObjectNull(gcpCredentialsAttributeTypes), client.OutletTransformState{SelectedTransform: stringPointerForTest("summarize"), EffectiveTransform: stringPointerForTest("summarize"), Applied: true})

	if state.Project.ValueString() != "example-project" {
		t.Fatalf("expected project to be set, got %s", state.Project)
	}
	if state.BatchingCount.ValueInt64() != 100 {
		t.Fatalf("expected batching_count to be set, got %d", state.BatchingCount.ValueInt64())
	}
	if selected := selectedTransform(state.Transform); selected.ValueString() != "summarize" || !transformBoolAttribute(t, state.Transform, "applied").ValueBool() {
		t.Fatalf("unexpected transform state: %#v", state.Transform)
	}
}

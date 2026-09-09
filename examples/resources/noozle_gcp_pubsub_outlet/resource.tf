# Copyright (c) HashiCorp, Inc.

resource "noozle_gcp_pubsub_outlet" "alerts" {
  name          = "GCP Pub/Sub alerts"
  enabled       = true
  project       = "example-project"
  topic         = "noozle-alerts"
  notify_policy = "FIRST_ONLY"

  transform {
    selected = "summarize"
  }
}

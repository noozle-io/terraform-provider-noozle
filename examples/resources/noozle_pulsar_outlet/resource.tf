# Copyright (c) HashiCorp, Inc.

resource "noozle_pulsar_outlet" "alerts" {
  name          = "Pulsar alerts"
  enabled       = true
  topic         = "persistent://public/default/noozle-alerts"
  url           = "pulsar://pulsar.example.com:6650"
  notify_policy = "FIRST_ONLY"

  transform {
    selected = "summarize"
  }
}

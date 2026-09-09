# Copyright (c) HashiCorp, Inc.

resource "noozle_kafka_outlet" "alerts" {
  name          = "Kafka alerts"
  enabled       = true
  addresses     = ["broker-1.example.internal:9093"]
  topic         = "events.outbound"
  notify_policy = "FIRST_ONLY"

  transform {
    selected = "summarize"
  }
}

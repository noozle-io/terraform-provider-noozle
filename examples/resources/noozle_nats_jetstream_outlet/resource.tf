# Copyright (c) HashiCorp, Inc.

resource "noozle_nats_jetstream_outlet" "alerts" {
  name          = "NATS JetStream alerts"
  enabled       = true
  subject       = "noozle.alerts"
  notify_policy = "FIRST_ONLY"

  transform {
    selected = "summarize"
  }

  credentials {
    urls = [
      "nats://nats-a.example.com:4222",
    ]
    urls_version = 1
  }
}

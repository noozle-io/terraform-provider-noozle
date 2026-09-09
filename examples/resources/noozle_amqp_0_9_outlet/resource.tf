# Copyright (c) HashiCorp, Inc.

resource "noozle_amqp_0_9_outlet" "alerts" {
  name          = "AMQP 0.9.1 alerts"
  enabled       = true
  exchange      = "noozle.alerts"
  key           = "alerts.critical"
  notify_policy = "FIRST_ONLY"

  transform {
    selected = "summarize"
  }

  credentials {
    urls = [
      "amqps://alerts:replace-me@rabbitmq-a.example.com:5671/%2f",
    ]
    urls_version = 1
  }
}

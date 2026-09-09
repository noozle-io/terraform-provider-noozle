# Copyright (c) HashiCorp, Inc.

resource "noozle_amqp_1_outlet" "alerts" {
  name           = "AMQP 1.0 alerts"
  enabled        = true
  target_address = "queue://noozle.alerts"
  notify_policy  = "FIRST_ONLY"

  transform {
    selected = "summarize"
  }

  credentials {
    urls = [
      "amqps://broker-a.example.com:5671",
    ]
    urls_version = 1
  }
}

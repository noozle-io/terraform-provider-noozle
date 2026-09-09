# Copyright (c) HashiCorp, Inc.

resource "noozle_http_server_outlet" "alerts" {
  name          = "HTTP Server alerts"
  enabled       = true
  path          = "/api/outlets/alerts"
  notify_policy = "FIRST_ONLY"

  transform {
    selected = "summarize"
  }
}

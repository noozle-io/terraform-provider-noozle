# Copyright (c) HashiCorp, Inc.

resource "noozle_sse_outlet" "alerts" {
  name          = "SSE alerts"
  enabled       = true
  event_type    = "article.match"
  path          = "/v1/sse/stream/noozle-alerts"
  retry_ms      = 5000
  notify_policy = "ALL"
}

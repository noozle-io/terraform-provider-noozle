# Copyright (c) HashiCorp, Inc.

resource "noozle_ntfy_outlet" "alerts" {
  name          = "ntfy alerts"
  enabled       = true
  topic         = "noozle-alerts"
  server_url    = "https://ntfy.example.com"
  notify_policy = "FIRST_ONLY"
}

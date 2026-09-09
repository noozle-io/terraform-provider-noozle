# Copyright (c) HashiCorp, Inc.

resource "noozle_aws_sqs_outlet" "alerts" {
  name          = "AWS SQS alerts"
  enabled       = true
  queue_url     = "https://sqs.us-east-1.amazonaws.com/123456789012/noozle-alerts"
  notify_policy = "FIRST_ONLY"
}

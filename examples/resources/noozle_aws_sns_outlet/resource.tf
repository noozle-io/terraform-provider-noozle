# Copyright (c) HashiCorp, Inc.

resource "noozle_aws_sns_outlet" "alerts" {
  name          = "AWS SNS alerts"
  enabled       = true
  topic_arn     = "arn:aws:sns:us-east-1:123456789012:noozle-alerts"
  notify_policy = "FIRST_ONLY"
}

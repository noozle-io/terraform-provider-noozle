# Copyright (c) HashiCorp, Inc.

resource "noozle_query" "terraform_test" {
  name            = "Terraform Test"
  full_expression = "language == \"en\" && about(\"Boeing 777x\", 0.8)"
  is_active       = true
}

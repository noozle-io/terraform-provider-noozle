# Copyright (c) HashiCorp, Inc.

resource "noozle_query_outlet_link" "terraform_example" {
  outlet_id = 123
  query_id  = noozle_query.terraform_test.id
}

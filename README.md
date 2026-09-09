# Terraform Provider Noozle

Terraform provider for managing Noozle queries and outlet configurations.

## Requirements

- Terraform `>= 1.0`
- Go `>= 1.25`

## Provider Configuration

The provider uses API-key authentication.

```hcl
provider "noozle" {
  host            = "https://noozle.io"
  api_key         = "your-api-key"
  tenant_id       = "your-tenant-id"
  organization_id = "your-organization-id"
  debug           = false
}
```

`tenant_id` is required when managing or reading queries, outlets, and
outlet-query links. It selects the tenant used by the tenant-aware API;
`organization_id` is optional and disambiguates tenant resolution.

Environment variables are also supported:

- `NOOZLE_HOST`
- `NOOZLE_API_KEY`
- `NOOZLE_TENANT_ID`
- `NOOZLE_ORGANIZATION_ID`

Set `debug = true` and run Terraform with `TF_LOG_PROVIDER=DEBUG` to see HTTP client request and response payload logs.

## Current Surface

Implemented resources and matching singular data sources include:

- `noozle_query`
- `noozle_query_outlet_link`
- outlet resources for `ntfy`, `sse`, `http_server`, `amqp_0_9`, `amqp_1`, `kafka`, `nats_jetstream`, `pulsar`, `aws_sns`, `aws_sqs`, and `gcp_pubsub`

Examples for every resource and data source live under [`examples/`](examples/). Import snippets are included under each managed resource example directory.

## Repository Layout

- [`internal/provider/`](internal/provider/) Terraform provider, resources, data sources, and tests
- [`internal/client/`](internal/client/) modular Noozle HTTP client
- [`examples/`](examples/) runnable Terraform examples

## Local Development

Common commands:

```sh
make fmt
make lint
make test
make build
make install
make generate
```

Notes:

- `make fmt` runs `gofmt -s -w -e .`
- `make lint` runs `golangci-lint`
- `make testacc` runs acceptance tests with `TF_ACC=1`
- `make install` places the local provider binary in your Go bin path

After any Go change, run `make fmt` and `make lint` before considering the work complete.

## Using The Provider Locally

For local Terraform development, install the provider binary and use a development override so Terraform loads your local build instead of a published registry release.

1. Build and install the provider:

```sh
make install
```

This places `terraform-provider-noozle` in your Go bin directory.

2. Create a CLI config file at `~/.terraformrc` on Linux/macOS or `%APPDATA%/terraform.rc` on Windows:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/noozle/noozle" = "/home/<your-user>/go/bin"
  }

  direct {}
}
```

Adjust the path to match your local Go bin directory.

3. Use the provider in Terraform:

```hcl
terraform {
  required_providers {
    noozle = {
      source = "noozle/noozle"
    }
  }
}

provider "noozle" {
  host      = "https://noozle.io"
  api_key   = "your-api-key"
  tenant_id = "your-tenant-id"
}
```

4. Reinstall the binary after Go code changes:

```sh
make install
```

Then rerun `terraform plan` or `terraform apply` from your example or test configuration.

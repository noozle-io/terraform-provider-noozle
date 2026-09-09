# Examples

This directory contains examples that are mostly used for documentation, but can also be run/tested manually via the Terraform CLI.

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or are testable even if some parts are not relevant for the documentation.

* **provider/provider.tf** example file for the provider index page
* **data-sources/`full data source name`/data-source.tf** example file for the named data source page
* **resources/`full resource name`/resource.tf** example file for the named data source page
Examples in this directory follow the Terraform provider documentation layout:

- `provider/provider.tf` for provider configuration
- `resources/<name>/resource.tf` for managed resources
- `data-sources/<name>/data-source.tf` for data sources

The current examples cover:

- `noozle_query`
- `noozle_query_outlet_link`
- `noozle_ntfy_outlet`
- `noozle_sse_outlet`
- `noozle_http_server_outlet`
- `noozle_pulsar_outlet`
- `noozle_nats_jetstream_outlet`
- `noozle_amqp_0_9_outlet`
- `noozle_amqp_1_outlet`
- `noozle_kafka_outlet`
- `noozle_aws_sns_outlet`
- `noozle_aws_sqs_outlet`
- `noozle_gcp_pubsub_outlet`
- `data.noozle_query`
- `data.noozle_query_outlet_link`
- `data.noozle_ntfy_outlet`
- `data.noozle_sse_outlet`
- `data.noozle_http_server_outlet`
- `data.noozle_pulsar_outlet`
- `data.noozle_nats_jetstream_outlet`
- `data.noozle_amqp_0_9_outlet`
- `data.noozle_amqp_1_outlet`
- `data.noozle_kafka_outlet`
- `data.noozle_aws_sns_outlet`
- `data.noozle_aws_sqs_outlet`
- `data.noozle_gcp_pubsub_outlet`

Managed resource examples also include `import.sh` snippets showing the expected import identifier format.

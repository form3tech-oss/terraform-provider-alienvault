# AlienVault Terraform Provider

[![Build Status](https://travis-ci.org/form3tech-oss/terraform-provider-alienvault.svg?branch=master)](https://travis-ci.org/form3tech-oss/terraform-provider-alienvault)

Terraform Provider for [AlienVault USM Anywhere](https://www.alienvault.com/products/usm-anywhere).

## Example Usage

```hcl
provider "alienvault" {
    fqdn     = "" # fill these in!
    username = ""
    password = ""
}

resource "alienvault_sensor" "main" {
  count   = "${var.provision_alienvault_sensor ? 1 : 0}"
  ip      = "${data.aws_instance.alienvault_sensor.public_ip}"
  name    = "${var.stack_name}-sensor"
}


resource "alienvault_job_aws_cloudwatch" "route53_siem" {
  count         = "${var.provision_alienvault_sensor ? 1 : 0}"
  name          = "Route53 log collection"
  sensor        = "${alienvault_sensor.main.id}"
  schedule      = "hourly"
  # to log route53 events to cloudwatch, we have to use the us-east1 region
  region        = "us-east-1"
  group         = "${aws_cloudwatch_log_group.aws_route53_stack.name}"
  source_format = "raw"
  plugin        = "Route 53 DNS Queries"
}
```
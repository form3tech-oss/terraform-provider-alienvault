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

resource "alienvault_job_aws_bucket" "nginx-logs-bucket-job" {
    name     = "nginx log collection"
    sensor   = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
    schedule = "hourly"
    bucket   = "this-does-not-exist"
    path     = "/something/logs"
    source_format = "raw"
    plugin   = "Nginx"
}

resource "alienvault_job_aws_cloudwatch" "test-e2e-cloudwatch-job" {
    name = "RDS log collection"
    sensor = "6a89f4aa-fa8e-44d4-9ffb-9ba1ae946777"
    schedule = "daily"
    region = "us-east-1"
    group = "test-group"
    stream = "test-stream"
    source_format = "raw"
    plugin = "PostgreSQL"
}
```
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
    schedule = "0 0 0/1 1/1 * ? *"
    bucket   = "this-does-not-exist"
    path     = "/something/logs"
    source_format = "raw"
    plugin   = "PostgreSQL"
}
```
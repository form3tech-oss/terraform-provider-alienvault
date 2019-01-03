resource "alienvault_job_aws_bucket" "example-bucket-job" {
    name = "test-example"
    description = "This is a test, feel free to remove"
    sensor = "6a89f4aa-fa8e-44d4-9ffb-9ba1ae946777"
    schedule = "0 0 0/1 1/1 * ? *"
    bucket = "test-bucket-123"
    source_format = "raw"
    plugin = "PostgreSQL"
}


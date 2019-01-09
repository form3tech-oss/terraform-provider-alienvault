package alienvault

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/form3tech-oss/alienvault"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccJobAWSBucketConfig = `
	resource "alienvault_job_aws_bucket" "test-e2e-bucket-job" {
		name = "%s"
		sensor = "my-sensor"
		schedule = "0 0 0/1 1/1 * ? *"
		bucket = "this-does-not-exist"
		path = "/something/logs"
		source_format = "raw"
		plugin = "PostgreSQL"
	}`

func TestAccResourceJobAWSBucket(t *testing.T) {
	var job alienvault.AWSBucketJob
	jobName := fmt.Sprintf("test-e2e-bucket-%d-%s", time.Now().UnixNano(), acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "alienvault_job_aws_bucket.test-e2e-bucket-job",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckJobAWSBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccJobAWSBucketConfig, jobName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckJobAWSBucketExists("alienvault_job_aws_bucket.test-e2e-bucket-job", &job),
					testAccCheckJobAWSBucketHasPresets("alienvault_job_aws_bucket.test-e2e-bucket-job", &job),
					resource.TestCheckResourceAttr("alienvault_job_aws_bucket.test-e2e-bucket-job", "name", jobName),
					resource.TestCheckResourceAttr("alienvault_job_aws_bucket.test-e2e-bucket-job", "sensor", "my-sensor"),
					resource.TestCheckResourceAttr("alienvault_job_aws_bucket.test-e2e-bucket-job", "schedule", "0 0 0/1 1/1 * ? *"),
					resource.TestCheckResourceAttr("alienvault_job_aws_bucket.test-e2e-bucket-job", "bucket", "this-does-not-exist"),
					resource.TestCheckResourceAttr("alienvault_job_aws_bucket.test-e2e-bucket-job", "path", "/something/logs"),
					resource.TestCheckResourceAttr("alienvault_job_aws_bucket.test-e2e-bucket-job", "source_format", "raw"),
					resource.TestCheckResourceAttr("alienvault_job_aws_bucket.test-e2e-bucket-job", "plugin", "PostgreSQL"),
					resource.TestCheckResourceAttr("alienvault_job_aws_bucket.test-e2e-bucket-job", "disabled", "false"),
				),
			},
		},
	})
}

func testAccCheckJobAWSBucketDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*alienvault.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alienvault_job_aws_bucket" {
			continue
		}

		_, err := client.GetAWSBucketJob(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("job %q still exists", rs.Primary.ID)
		}

		if !strings.Contains(err.Error(), "could not be found") {
			return fmt.Errorf("Unexpected error when checking for existence of job: %s", err)
		}
	}

	return nil
}

func testAccCheckJobAWSBucketHasPresets(n string, res *alienvault.AWSBucketJob) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no job ID is set")
		}

		client := testAccProvider.Meta().(*alienvault.Client)

		job, err := client.GetAWSBucketJob(rs.Primary.ID)
		if err != nil {
			return err
		}

		if job.App != alienvault.JobApplicationAWS {
			return fmt.Errorf("unexpected job application: '%s'", job.App)
		}

		if !job.Custom {
			return fmt.Errorf("unexpected job state - should be flagged as a custom job but is not")
		}

		if job.Action != alienvault.JobActionMonitorBucket {
			return fmt.Errorf("unexpected job action: '%s'", job.Action)
		}

		if job.Type != alienvault.JobTypeCollection {
			return fmt.Errorf("unexpected job type: '%s'", job.Type)
		}
		return nil
	}
}

func testAccCheckJobAWSBucketExists(n string, res *alienvault.AWSBucketJob) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no job ID is set")
		}

		client := testAccProvider.Meta().(*alienvault.Client)

		job, err := client.GetAWSBucketJob(rs.Primary.ID)
		if err != nil {
			return err
		}

		*res = *job
		return nil
	}
}

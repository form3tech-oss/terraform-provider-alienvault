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

const testAccJobAWSCloudWatchConfig = `
	resource "alienvault_job_aws_cloudwatch" "test-e2e-cloudwatch-job" {
		name = "%s"
		sensor = "my-sensor"
		schedule = "0 0 0/1 1/1 * ? *"
		region = "us-east-1"
		group = "test-group"
		stream = "test-stream"
		source_format = "raw"
		plugin = "PostgreSQL"
	}`

func TestAccResourceJobAWSCloudWatch(t *testing.T) {
	var job alienvault.AWSCloudWatchJob
	jobName := fmt.Sprintf("test-e2e-cloudwatch-%d-%s", time.Now().UnixNano(), acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckJobAWSCloudWatchDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccJobAWSCloudWatchConfig, jobName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckJobAWSCloudWatchExists("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", &job),
					testAccCheckJobAWSCloudWatchHasPresets("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", &job),
					resource.TestCheckResourceAttr("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", "name", jobName),
					resource.TestCheckResourceAttr("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", "sensor", "my-sensor"),
					resource.TestCheckResourceAttr("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", "schedule", "0 0 0/1 1/1 * ? *"),
					resource.TestCheckResourceAttr("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", "region", "us-east-1"),
					resource.TestCheckResourceAttr("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", "group", "test-group"),
					resource.TestCheckResourceAttr("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", "stream", "test-stream"),
					resource.TestCheckResourceAttr("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", "source_format", "raw"),
					resource.TestCheckResourceAttr("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", "plugin", "PostgreSQL"),
					resource.TestCheckResourceAttr("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", "disabled", "false"),
				),
			},
		},
	})
}

func testAccCheckJobAWSCloudWatchDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*alienvault.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alienvault_job_aws_cloudWatch" {
			continue
		}

		_, err := client.GetAWSCloudWatchJob(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("job %q still exists", rs.Primary.ID)
		}

		if !strings.Contains(err.Error(), "could not be found") {
			return fmt.Errorf("Unexpected error when checking for existence of job: %s", err)
		}
	}

	return nil
}

func testAccCheckJobAWSCloudWatchHasPresets(n string, res *alienvault.AWSCloudWatchJob) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no job ID is set")
		}

		client := testAccProvider.Meta().(*alienvault.Client)

		job, err := client.GetAWSCloudWatchJob(rs.Primary.ID)
		if err != nil {
			return err
		}

		if job.App != alienvault.JobApplicationAWS {
			return fmt.Errorf("unexpected job application: '%s'", job.App)
		}

		if !job.Custom {
			return fmt.Errorf("unexpected job state - should be flagged as a custom job but is not")
		}

		if job.Action != alienvault.JobActionMonitorCloudWatch {
			return fmt.Errorf("unexpected job action: '%s'", job.Action)
		}

		if job.Type != alienvault.JobTypeCollection {
			return fmt.Errorf("unexpected job type: '%s'", job.Type)
		}
		return nil
	}
}

func testAccCheckJobAWSCloudWatchExists(n string, res *alienvault.AWSCloudWatchJob) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no job ID is set")
		}

		client := testAccProvider.Meta().(*alienvault.Client)

		job, err := client.GetAWSCloudWatchJob(rs.Primary.ID)
		if err != nil {
			return err
		}

		*res = *job
		return nil
	}
}

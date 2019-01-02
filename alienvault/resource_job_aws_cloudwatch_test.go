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
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"gotest.tools/assert"
)

const testAccJobAWSCloudWatchConfig_basic = `
	resource "alienvault_job_aws_cloudwatch" "test-e2e-cloudwatch-job" {
		name = "%s"
		sensor = "6a89f4aa-fa8e-44d4-9ffb-9ba1ae946777"
		schedule = "0 0 0/1 1/1 * ? *"
		region = "us-east-1"
		group = "test-group"
		stream = "test-stream"
		source_format = "raw"
		plugin = "PostgreSQL"
	}`

func TestAccResourceJobAWSCloudWatch(t *testing.T) {
	var job alienvault.Job
	jobName := fmt.Sprintf("test-e2e-cloudwatch-%d-%s", time.Now().UnixNano(), acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckJobAWSCloudWatchDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccJobAWSCloudWatchConfig_basic, jobName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckJobAWSCloudWatchExists("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", &job),
					testAccCheckJobAWSCloudWatchHasPresets("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", &job),
					resource.TestCheckResourceAttr("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", "name", jobName),
					resource.TestCheckResourceAttr("alienvault_job_aws_cloudwatch.test-e2e-cloudwatch-job", "sensor", "6a89f4aa-fa8e-44d4-9ffb-9ba1ae946777"),
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

		_, err := client.GetJob(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("job %q still exists", rs.Primary.ID)
		}

		if !strings.Contains(err.Error(), "could not be found") {
			return fmt.Errorf("Unexpected error when checking for existence of job: %s", err)
		}
	}

	return nil
}

func testAccCheckJobAWSCloudWatchHasPresets(n string, res *alienvault.Job) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no job ID is set")
		}

		client := testAccProvider.Meta().(*alienvault.Client)

		job, err := client.GetJob(rs.Primary.ID)
		if err != nil {
			return err
		}

		if job.App != JobAppAWS {
			return fmt.Errorf("unexpected job application: '%s'", job.App)
		}

		if !job.Custom {
			return fmt.Errorf("unexpected job state - should be flagged as a custom job but is not")
		}

		if job.Action != JobActionMonitorCloudWatch {
			return fmt.Errorf("unexpected job action: '%s'", job.Action)
		}

		if job.Type != JobTypeCollection {
			return fmt.Errorf("unexpected job type: '%s'", job.Type)
		}
		return nil
	}
}

func testAccCheckJobAWSCloudWatchExists(n string, res *alienvault.Job) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no job ID is set")
		}

		client := testAccProvider.Meta().(*alienvault.Client)

		job, err := client.GetJob(rs.Primary.ID)
		if err != nil {
			return err
		}

		*res = *job
		return nil
	}
}

func TestFlattenJobAWSCloudWatch(t *testing.T) {

	resourceLocalData := schema.TestResourceDataRaw(t, resourceJobAWSCloudWatch().Schema, map[string]interface{}{})

	job := &alienvault.Job{
		UUID:        "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		Name:        "cloudWatch-job",
		Description: "A job for retrieving some logs from a cloudWatch",
		SensorID:    "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		Schedule:    "0 0 0/1 1/1 * ? *",
		Disabled:    true,
		Params: map[string]interface{}{
			"regionName": "us-east-1",
			"groupName":  "test-group",
			"streamName": "test-stream",
			"source":     "raw",
			"plugin":     "PostgreSQL",
		},
	}

	flattenJobAWSCloudWatch(job, resourceLocalData)

	assert.Equal(t, job.Name, resourceLocalData.Get("name").(string))
	assert.Equal(t, job.Description, resourceLocalData.Get("description").(string))
	assert.Equal(t, job.UUID, resourceLocalData.Get("uuid").(string))
	assert.Equal(t, job.SensorID, resourceLocalData.Get("sensor").(string))
	assert.Equal(t, job.Schedule, resourceLocalData.Get("schedule").(string))
	assert.Equal(t, job.Disabled, resourceLocalData.Get("disabled").(bool))
	assert.Equal(t, job.Params["regionName"], resourceLocalData.Get("region").(string))
	assert.Equal(t, job.Params["groupName"], resourceLocalData.Get("group").(string))
	assert.Equal(t, job.Params["streamName"], resourceLocalData.Get("stream").(string))
	assert.Equal(t, job.Params["source"], resourceLocalData.Get("source_format").(string))
	assert.Equal(t, job.Params["plugin"], resourceLocalData.Get("plugin").(string))
}

func TestExpandJobAWSCloudWatch(t *testing.T) {

	input := map[string]interface{}{
		"name":          "cloudWatch-job",
		"description":   "A job for retrieving some logs from a cloudWatch",
		"sensor":        "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		"schedule":      "0 0 0/1 1/1 * ? *",
		"disabled":      true,
		"region":        "us-east-1",
		"group":         "test-group",
		"stream":        "test-stream",
		"source_format": JobSourceFormatRaw,
		"plugin":        "PostgreSQL",
	}

	resourceLocalData := schema.TestResourceDataRaw(t, resourceJobAWSCloudWatch().Schema, input)
	resourceLocalData.SetId("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")

	job := expandJobAWSCloudWatch(resourceLocalData)

	assert.Equal(t, job.Name, input["name"])
	assert.Equal(t, job.Description, input["description"])
	assert.Equal(t, job.UUID, resourceLocalData.Id())
	assert.Equal(t, job.SensorID, input["sensor"])
	assert.Equal(t, job.Schedule, input["schedule"])
	assert.Equal(t, job.Disabled, input["disabled"])
	assert.Equal(t, job.Params["regionName"], input["region"])
	assert.Equal(t, job.Params["groupName"], input["group"])
	assert.Equal(t, job.Params["streamName"], input["stream"])
	assert.Equal(t, job.Params["source"], input["source_format"])
	assert.Equal(t, job.Params["plugin"], input["plugin"])
	assert.Equal(t, job.Custom, true)
	assert.Equal(t, job.App, JobAppAWS)
	assert.Equal(t, job.Action, JobActionMonitorCloudWatch)
	assert.Equal(t, job.Type, JobTypeCollection)

}

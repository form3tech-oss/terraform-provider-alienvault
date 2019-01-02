package alienvault

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"gotest.tools/assert"
)

func TestFlattenJobAWSBucket(t *testing.T) {

	resourceLocalData := schema.TestResourceDataRaw(t, resourceJobAWSBucket().Schema, map[string]interface{}{})

	job := &Job{
		UUID:        "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		Name:        "bucket-job",
		Description: "A job for retrieving some logs from a bucket",
		SensorID:    "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		Schedule:    "0 0 0/1 1/1 * ? *",
		Disabled:    true,
		Params: map[string]interface{}{
			"bucketName": "special-bucket-72",
			"path":       "/database-logs/",
			"source":     "raw",
			"plugin":     "PostgreSQL",
		},
	}

	flattenJobAWSBucket(job, resourceLocalData)

	assert.Equal(t, job.Name, resourceLocalData.Get("name").(string))
	assert.Equal(t, job.Description, resourceLocalData.Get("description").(string))
	assert.Equal(t, job.UUID, resourceLocalData.Get("uuid").(string))
	assert.Equal(t, job.SensorID, resourceLocalData.Get("sensor").(string))
	assert.Equal(t, job.Schedule, resourceLocalData.Get("schedule").(string))
	assert.Equal(t, job.Disabled, resourceLocalData.Get("disabled").(bool))
	assert.Equal(t, job.Params["bucketName"], resourceLocalData.Get("bucket").(string))
	assert.Equal(t, job.Params["path"], resourceLocalData.Get("path").(string))
	assert.Equal(t, job.Params["source"], resourceLocalData.Get("source_format").(string))
	assert.Equal(t, job.Params["plugin"], resourceLocalData.Get("plugin").(string))
}

func TestExpandJobAWSBucket(t *testing.T) {

	input := map[string]interface{}{
		"name":          "bucket-job",
		"description":   "A job for retrieving some logs from a bucket",
		"sensor":        "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		"schedule":      "0 0 0/1 1/1 * ? *",
		"disabled":      true,
		"bucket":        "special-bucket-72",
		"path":          "/database-logs/",
		"source_format": JOB_SOURCE_FORMAT_RAW,
		"plugin":        "PostgreSQL",
	}

	resourceLocalData := schema.TestResourceDataRaw(t, resourceJobAWSBucket().Schema, input)
	resourceLocalData.SetId("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")

	job := expandJobAWSBucket(resourceLocalData)

	assert.Equal(t, job.Name, input["name"])
	assert.Equal(t, job.Description, input["description"])
	assert.Equal(t, job.UUID, resourceLocalData.Id())
	assert.Equal(t, job.SensorID, input["sensor"])
	assert.Equal(t, job.Schedule, input["schedule"])
	assert.Equal(t, job.Disabled, input["disabled"])
	assert.Equal(t, job.Params["bucketName"], input["bucket"])
	assert.Equal(t, job.Params["path"], input["path"])
	assert.Equal(t, job.Params["source"], input["source_format"])
	assert.Equal(t, job.Params["plugin"], input["plugin"])
	assert.Equal(t, job.Custom, true)
	assert.Equal(t, job.App, JOB_APP_AWS)
	assert.Equal(t, job.Action, JOB_ACTION_MONITOR_BUCKET)
	assert.Equal(t, job.Type, JOB_TYPE_COLLECTION)

}

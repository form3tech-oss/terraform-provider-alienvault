package alienvault

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceJobAWSBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobAWSBucketCreate,
		Read:   resourceJobAWSBucketRead,
		Update: resourceJobAWSBucketUpdate,
		Delete: resourceJobAWSBucketDelete,

		Schema: map[string]*schema.Schema{
			"uuid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"sensor": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateJobSensor,
			},
			"schedule": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bucket": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_format": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateJobSourceFormat,
			},
			"plugin": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateJobPlugin,
			},
		},
	}
}

func resourceJobAWSBucketCreate(d *schema.ResourceData, m interface{}) error {

	job := expandJobAWSBucket(d)

	if err := m.(*Client).CreateJob(job); err != nil {
		return err
	}

	if job.UUID == "" {
		return fmt.Errorf("Failed to determine UUID of created resource")
	}

	d.SetId(job.UUID)
	return resourceJobAWSBucketRead(d, m)
}

func resourceJobAWSBucketRead(d *schema.ResourceData, m interface{}) error {
	job, err := m.(*Client).GetJob(d.Id())
	if err != nil {
		return err
	}
	flattenJobAWSBucket(job, d)
	return nil
}

func resourceJobAWSBucketUpdate(d *schema.ResourceData, m interface{}) error {

	job := expandJobAWSBucket(d)
	if err := m.(*Client).UpdateJob(job); err != nil {
		return err
	}

	return resourceJobAWSBucketRead(d, m)
}

func resourceJobAWSBucketDelete(d *schema.ResourceData, m interface{}) error {
	return m.(*Client).DeleteJob(d.Id())
}

func flattenJobAWSBucket(job *Job, d *schema.ResourceData) {

	if job.UUID != "" {
		d.SetId(job.UUID)
		d.Set("uuid", job.UUID)
	}

	if job.Description != "" {
		d.Set("description", job.Description)
	}

	if job.Params != nil {
		if bucket, ok := job.Params["bucketName"]; ok {
			if bucketSafe, ok := bucket.(string); ok {
				d.Set("bucket", bucketSafe)
			}
		}

		if source, ok := job.Params["source"]; ok {
			if sourceSafe, ok := source.(string); ok {
				d.Set("source_format", sourceSafe)
			}
		}

		if plugin, ok := job.Params["plugin"]; ok {
			if pluginSafe, ok := plugin.(string); ok {
				d.Set("plugin", pluginSafe)
			}
		}

		if path, ok := job.Params["path"]; ok {
			if pathSafe, ok := path.(string); ok {
				d.Set("path", pathSafe)
			}
		}
	}

	d.Set("name", job.Name)
	d.Set("sensor", job.SensorID)
	d.Set("schedule", job.Schedule)
	d.Set("disabled", job.Disabled)
}

func expandJobAWSBucket(d *schema.ResourceData) *Job {

	job := &Job{}
	job.Name = d.Get("name").(string)
	job.SensorID = d.Get("sensor").(string)
	job.Schedule = d.Get("schedule").(string)
	job.Disabled = d.Get("disabled").(bool)
	job.Description = d.Get("description").(string)

	job.Params = map[string]interface{}{
		"bucketName": d.Get("bucket").(string),
		"source":     d.Get("source_format").(string),
	}

	if plugin, ok := d.GetOk("plugin"); ok {
		job.Params["plugin"] = plugin.(string)
	}

	if path, ok := d.GetOk("path"); ok {
		job.Params["path"] = path.(string)
	}

	// these are locked as specific to this type of job
	job.Custom = true
	job.App = JOB_APP_AWS
	job.Action = JOB_ACTION_MONITOR_BUCKET
	job.Type = JOB_TYPE_COLLECTION

	if d.Id() != "" {
		job.UUID = d.Id()
	}

	return job
}

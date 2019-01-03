package alienvault

import (
	"fmt"

	"github.com/form3tech-oss/alienvault"
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "This uses a non-standard cron format to schedule the job, which AV have not currently documented. For now you can use 'daily' and 'hourly' here, which will be automatically converted to the AV cron format by this provider.",
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

	if err := m.(*alienvault.Client).CreateAWSBucketJob(job); err != nil {
		return err
	}

	if job.UUID == "" {
		return fmt.Errorf("Failed to determine UUID of created resource")
	}

	d.SetId(job.UUID)
	return resourceJobAWSBucketRead(d, m)
}

func resourceJobAWSBucketRead(d *schema.ResourceData, m interface{}) error {
	job, err := m.(*alienvault.Client).GetAWSBucketJob(d.Id())
	if err != nil {
		return err
	}
	flattenJobAWSBucket(job, d)
	return nil
}

func resourceJobAWSBucketUpdate(d *schema.ResourceData, m interface{}) error {

	job := expandJobAWSBucket(d)
	if err := m.(*alienvault.Client).UpdateAWSBucketJob(job); err != nil {
		return err
	}

	return resourceJobAWSBucketRead(d, m)
}

func resourceJobAWSBucketDelete(d *schema.ResourceData, m interface{}) error {
	return m.(*alienvault.Client).DeleteAWSBucketJob(d.Id())
}

func flattenJobAWSBucket(job *alienvault.AWSBucketJob, d *schema.ResourceData) {

	if job.UUID != "" {
		d.SetId(job.UUID)
		d.Set("uuid", job.UUID)
	}

	if job.Description != "" {
		d.Set("description", job.Description)
	}

	d.Set("bucket", job.Params.BucketName)
	d.Set("path", job.Params.Path)
	d.Set("source_format", job.Params.SourceFormat)
	d.Set("plugin", job.Params.Plugin)

	d.Set("name", job.Name)
	d.Set("sensor", job.SensorID)
	d.Set("schedule", translateScheduleToTF(job.Schedule))
	d.Set("disabled", job.Disabled)
}

func expandJobAWSBucket(d *schema.ResourceData) *alienvault.AWSBucketJob {

	job := &alienvault.AWSBucketJob{}
	job.Name = d.Get("name").(string)
	job.SensorID = d.Get("sensor").(string)
	job.Schedule = translateScheduleFromTF(d.Get("schedule").(string))
	job.Disabled = d.Get("disabled").(bool)
	job.Description = d.Get("description").(string)

	job.Params.BucketName = d.Get("bucket").(string)
	job.Params.SourceFormat = alienvault.JobSourceFormat(d.Get("source_format").(string))

	if plugin, ok := d.GetOk("plugin"); ok {
		job.Params.Plugin = plugin.(string)
	}

	if path, ok := d.GetOk("path"); ok {
		job.Params.Path = path.(string)
	}

	if d.Id() != "" {
		job.UUID = d.Id()
	}

	return job
}

var scheduleMap = map[string]alienvault.JobSchedule{
	"hourly": alienvault.JobScheduleHourly,
	"daily":  alienvault.JobScheduleDaily,
}

func translateScheduleFromTF(schedule string) alienvault.JobSchedule {
	if js, ok := scheduleMap[schedule]; ok {
		return js
	}

	return alienvault.JobSchedule(schedule)
}

func translateScheduleToTF(schedule alienvault.JobSchedule) string {
	for tf, js := range scheduleMap {
		if schedule == js {
			return tf
		}
	}

	return string(schedule)
}

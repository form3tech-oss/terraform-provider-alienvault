package alienvault

import (
	"fmt"

	"github.com/form3tech-oss/alienvault"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceJobAWSCloudWatch() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobAWSCloudWatchCreate,
		Read:   resourceJobAWSCloudWatchRead,
		Update: resourceJobAWSCloudWatchUpdate,
		Delete: resourceJobAWSCloudWatchDelete,

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
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "AWS region (e.g. us-west-2) you are collecting CloudWatch Logs from. Use * for all.",
			},
			"group": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "CloudWatch Logs group that contains one or more streams with the same log content",
			},
			"stream": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "CloudWatch Logs stream name (e.g. i-038a3acdf913dd5fa or * to include all streams)",
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

func resourceJobAWSCloudWatchCreate(d *schema.ResourceData, m interface{}) error {

	job := expandJobAWSCloudWatch(d)

	if err := m.(*alienvault.Client).CreateAWSCloudWatchJob(job); err != nil {
		return err
	}

	if job.UUID == "" {
		return fmt.Errorf("Failed to determine UUID of created resource")
	}

	d.SetId(job.UUID)
	return resourceJobAWSCloudWatchRead(d, m)
}

func resourceJobAWSCloudWatchRead(d *schema.ResourceData, m interface{}) error {
	job, err := m.(*alienvault.Client).GetAWSCloudWatchJob(d.Id())
	if err != nil {
		return err
	}
	flattenJobAWSCloudWatch(job, d)
	return nil
}

func resourceJobAWSCloudWatchUpdate(d *schema.ResourceData, m interface{}) error {

	job := expandJobAWSCloudWatch(d)
	if err := m.(*alienvault.Client).UpdateAWSCloudWatchJob(job); err != nil {
		return err
	}

	return resourceJobAWSCloudWatchRead(d, m)
}

func resourceJobAWSCloudWatchDelete(d *schema.ResourceData, m interface{}) error {
	return m.(*alienvault.Client).DeleteAWSCloudWatchJob(d.Id())
}

func flattenJobAWSCloudWatch(job *alienvault.AWSCloudWatchJob, d *schema.ResourceData) {

	if job.UUID != "" {
		d.SetId(job.UUID)
		d.Set("uuid", job.UUID)
	}

	if job.Description != "" {
		d.Set("description", job.Description)
	}

	d.Set("region", job.Params.Region)
	d.Set("group", job.Params.Group)
	d.Set("stream", job.Params.Stream)
	d.Set("source_format", job.Params.SourceFormat)
	d.Set("plugin", job.Params.Plugin)

	d.Set("name", job.Name)
	d.Set("sensor", job.SensorID)
	d.Set("schedule", translateScheduleToTF(job.Schedule))
	d.Set("disabled", job.Disabled)
}

func expandJobAWSCloudWatch(d *schema.ResourceData) *alienvault.AWSCloudWatchJob {

	job := &alienvault.AWSCloudWatchJob{}
	job.Name = d.Get("name").(string)
	job.SensorID = d.Get("sensor").(string)
	job.Schedule = translateScheduleFromTF(d.Get("schedule").(string))
	job.Disabled = d.Get("disabled").(bool)
	job.Description = d.Get("description").(string)

	job.Params.Region = d.Get("region").(string)
	job.Params.Group = d.Get("group").(string)
	job.Params.Stream = d.Get("stream").(string)
	job.Params.SourceFormat = alienvault.JobSourceFormat(d.Get("source_format").(string))

	if plugin, ok := d.GetOk("plugin"); ok {
		job.Params.Plugin = plugin.(string)
	}

	if d.Id() != "" {
		job.UUID = d.Id()
	}

	return job
}
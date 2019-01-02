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

	if err := m.(*alienvault.Client).CreateJob(job); err != nil {
		return err
	}

	if job.UUID == "" {
		return fmt.Errorf("Failed to determine UUID of created resource")
	}

	d.SetId(job.UUID)
	return resourceJobAWSCloudWatchRead(d, m)
}

func resourceJobAWSCloudWatchRead(d *schema.ResourceData, m interface{}) error {
	job, err := m.(*alienvault.Client).GetJob(d.Id())
	if err != nil {
		return err
	}
	flattenJobAWSCloudWatch(job, d)
	return nil
}

func resourceJobAWSCloudWatchUpdate(d *schema.ResourceData, m interface{}) error {

	job := expandJobAWSCloudWatch(d)
	if err := m.(*alienvault.Client).UpdateJob(job); err != nil {
		return err
	}

	return resourceJobAWSCloudWatchRead(d, m)
}

func resourceJobAWSCloudWatchDelete(d *schema.ResourceData, m interface{}) error {
	return m.(*alienvault.Client).DeleteJob(d.Id())
}

func flattenJobAWSCloudWatch(job *alienvault.Job, d *schema.ResourceData) {

	if job.UUID != "" {
		d.SetId(job.UUID)
		d.Set("uuid", job.UUID)
	}

	if job.Description != "" {
		d.Set("description", job.Description)
	}

	if job.Params != nil {
		if region, ok := job.Params["regionName"]; ok {
			if regionSafe, ok := region.(string); ok {
				d.Set("region", regionSafe)
			}
		}

		if group, ok := job.Params["groupName"]; ok {
			if groupSafe, ok := group.(string); ok {
				d.Set("group", groupSafe)
			}
		}

		if stream, ok := job.Params["streamName"]; ok {
			if streamSafe, ok := stream.(string); ok {
				d.Set("stream", streamSafe)
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

	}

	d.Set("name", job.Name)
	d.Set("sensor", job.SensorID)
	d.Set("schedule", job.Schedule)
	d.Set("disabled", job.Disabled)
}

func expandJobAWSCloudWatch(d *schema.ResourceData) *alienvault.Job {

	job := &alienvault.Job{}
	job.Name = d.Get("name").(string)
	job.SensorID = d.Get("sensor").(string)
	job.Schedule = d.Get("schedule").(string)
	job.Disabled = d.Get("disabled").(bool)
	job.Description = d.Get("description").(string)

	job.Params = map[string]interface{}{
		"regionName": d.Get("region").(string),
		"groupName":  d.Get("group").(string),
		"streamName": d.Get("stream").(string),
		"source":     d.Get("source_format").(string),
	}

	if plugin, ok := d.GetOk("plugin"); ok {
		job.Params["plugin"] = plugin.(string)
	}

	// these are locked as specific to this type of job
	job.Custom = true
	job.App = JobAppAWS
	job.Action = JobActionMonitorCloudWatch
	job.Type = JobTypeCollection

	if d.Id() != "" {
		job.UUID = d.Id()
	}

	return job
}

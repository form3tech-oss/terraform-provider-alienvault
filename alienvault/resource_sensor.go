package alienvault

import (
	"context"
	"net"
	"time"

	"github.com/form3tech-oss/alienvault"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSensor() *schema.Resource {

	// create time has to take into account the time for the sensor appliance
	// to be activated and configured, which is usually 20-30m in total...
	createTime := time.Hour

	return &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: &createTime,
		},
		Create: resourceSensorCreate,
		Update: resourceSensorUpdate,
		Read:   resourceSensorRead,
		Delete: resourceSensorDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the sensor",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-readable description of the sensor",
				Default:     "Created by terraform",
			},
			"ip": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The public IP address of the sensor",
				ForceNew:     true, // register new sensor as the appliance has been recreated/changed
				ValidateFunc: validateIP,
			},
		},
	}
}

func resourceSensorCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alienvault.Client)

	sensor := expandSensor(d)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	ip := net.ParseIP(d.Get("ip").(string))
	if ip == nil {
		// this is a panic because it should never happen - the IP field will be ensured to be a valid IP by the schema ValidateFunc
		panic("Failed to parse valid IP")
	}

	if err := client.CreateSensorViaAppliance(ctx, sensor, ip); err != nil {
		return err
	}

	d.SetId(sensor.UUID)
	return resourceSensorRead(d, m)
}

func resourceSensorUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alienvault.Client)
	sensor := expandSensor(d)
	return client.UpdateSensor(sensor)
}

func resourceSensorRead(d *schema.ResourceData, m interface{}) error {
	sensor, err := m.(*alienvault.Client).GetSensor(d.Id())
	if err != nil {
		return err
	}
	flattenSensor(sensor, d)
	return nil
}

func resourceSensorDelete(d *schema.ResourceData, m interface{}) error {
	sensor := expandSensor(d)
	return m.(*alienvault.Client).DeleteSensor(sensor)
}

func flattenSensor(sensor *alienvault.Sensor, d *schema.ResourceData) {
	d.SetId(sensor.UUID)
	d.Set("name", sensor.Name)
	d.Set("description", sensor.Description)
}

func expandSensor(d *schema.ResourceData) *alienvault.Sensor {
	sensor := &alienvault.Sensor{}
	sensor.UUID = d.Id()
	sensor.Name = d.Get("name").(string)
	if description, ok := d.GetOk("description"); ok {
		sensor.Description = description.(string)
	}
	return sensor
}

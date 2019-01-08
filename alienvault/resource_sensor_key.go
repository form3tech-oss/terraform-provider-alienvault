package alienvault

import (
	"fmt"

	"github.com/form3tech-oss/alienvault"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSensorKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSensorKeyCreate,
		Read:   resourceSensorKeyRead,
		Delete: resourceSensorKeyDelete,

		Schema: map[string]*schema.Schema{},
	}
}

func resourceSensorKeyCreate(d *schema.ResourceData, m interface{}) error {

	key, err := m.(*alienvault.Client).CreateSensorKey(false)
	if err != nil {
		return err
	}

	if key.ID == "" {
		return fmt.Errorf("Failed to determine ID of created key resource")
	}

	d.SetId(key.ID)
	return resourceSensorKeyRead(d, m)
}

func resourceSensorKeyRead(d *schema.ResourceData, m interface{}) error {
	job, err := m.(*alienvault.Client).GetSensorKey(d.Id())
	if err != nil {
		return err
	}
	flattenSensorKey(job, d)
	return nil
}

func resourceSensorKeyUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSensorKeyDelete(d *schema.ResourceData, m interface{}) error {
	key := expandSensorKey(d)
	return m.(*alienvault.Client).DeleteSensorKey(key)
}

func flattenSensorKey(key *alienvault.SensorKey, d *schema.ResourceData) {

	if key.ID != "" {
		d.SetId(key.ID)
	}

}

func expandSensorKey(d *schema.ResourceData) *alienvault.SensorKey {
	key := &alienvault.SensorKey{}
	key.ID = d.Id()
	return key
}

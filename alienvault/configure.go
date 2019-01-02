package alienvault

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	client := NewClient(
		d.Get("fqdn").(string),
		Credentials{
			Username: d.Get("username").(string),
			Password: d.Get("password").(string),
		})

	if err := client.Authenticate(); err != nil {
		return nil, err
	}

	return client, nil
}

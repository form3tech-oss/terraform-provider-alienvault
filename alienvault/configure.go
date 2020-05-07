package alienvault

import (
	"fmt"

	"github.com/form3tech-oss/alienvault"
	"github.com/hashicorp/terraform/helper/schema"
)

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := alienvault.New(
		d.Get("fqdn").(string),
		alienvault.Credentials{
			Username: d.Get("username").(string),
			Password: d.Get("password").(string),
		},
		d.Get("skip_tls_verify").(bool),
		d.Get("api_version").(int),
	)

	if err := client.Authenticate(); err != nil {
		return nil, fmt.Errorf("failed in authenticate: %w", err)
	}

	return client, nil
}

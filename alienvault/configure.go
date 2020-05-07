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
		return nil, fmt.Errorf("failed in authenticate (u: %s, p: %s): %w", d.Get("username").(hiddenString), d.Get("password").(hiddenString), err)
	}

	return client, nil
}

type hiddenString interface{}

func (v hiddenString) String() string {
	sv := string(v)
	if len(sv) < 3 {
		return "****"
	}
	return "****" + vs[len(sv)-3:]
}

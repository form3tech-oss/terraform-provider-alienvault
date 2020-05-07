package alienvault

import (
	"fmt"
	"strconv"

	"github.com/form3tech-oss/alienvault"
	"github.com/hashicorp/terraform/helper/schema"
)

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	version, ok := d.Get("api_version").(int)
	if !ok {
		if versionStr, ok := d.Get("api_version").(string); ok && versionStr != "" {
			var err error
			version, err = strconv.Atoi(versionStr)
			if err != nil {
				return nil, fmt.Errorf("invalid version: %s", versionStr)
			}
		}
	}

	client := alienvault.New(
		d.Get("fqdn").(string),
		alienvault.Credentials{
			Username: d.Get("username").(string),
			Password: d.Get("password").(string),
		},
		d.Get("skip_tls_verify").(bool),
		version,
	)

	if err := client.Authenticate(); err != nil {
		return nil, fmt.Errorf("failed in authenticate (u: %s, p: %s): %w", d.Get("username").(hiddenValue), d.Get("password").(hiddenValue), err)
	}

	return client, nil
}

type hiddenValue string

func (v hiddenValue) String() string {
	if len(v) < 3 {
		return "****"
	}
	return "****" + string(v[len(v)-3:])
}

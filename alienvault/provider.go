package alienvault

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider makes the AlienVault provider available
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"fqdn": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The fully qualified domain name for your AlienVault instance e.g. example.alienvault.cloud",
				DefaultFunc: schema.EnvDefaultFunc("ALIENVAULT_FQDN", nil),
				Sensitive:   true,
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "AV username",
				DefaultFunc: schema.EnvDefaultFunc("ALIENVAULT_USERNAME", nil),
				Sensitive:   true,
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "AV password",
				DefaultFunc: schema.EnvDefaultFunc("ALIENVAULT_PASSWORD", nil),
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"alienvault_job_aws_bucket":     resourceJobAWSBucket(),
			"alienvault_job_aws_cloudwatch": resourceJobAWSCloudWatch(),
			"alienvault_sensor":             resourceSensor(),
		},
		ConfigureFunc: providerConfigure,
	}
}

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
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "AV username",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "AV password",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"alienvault_job_aws_bucket": resourceJobAWSBucket(),
			//"alienvault_job_aws_cloudwatch": resourceJobAWSCloudWatch(),
		},
		ConfigureFunc: providerConfigure,
	}
}

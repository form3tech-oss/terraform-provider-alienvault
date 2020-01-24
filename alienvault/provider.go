package alienvault

import (
    "github.com/hashicorp/terraform/helper/schema"
    "os"
)

// Provider makes the AlienVault provider available
func Provider() *schema.Provider {
    return &schema.Provider{
        Schema: map[string]*schema.Schema{
            "fqdn": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "The fully qualified domain name for your AlienVault instance e.g. example.alienvault.cloud",
                DefaultFunc: schema.EnvDefaultFunc("ALIENVAULT_FQDN", nil),
                Sensitive:   true,
            },
            "api_version": {
                Type:        schema.TypeInt,
                Description: "The api version you are using, normally 1 or 2",
                Default: 1,
                Optional:    true,
            },
            "username": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "AV username",
                DefaultFunc: schema.EnvDefaultFunc("ALIENVAULT_USERNAME", nil),
                Sensitive:   true,
            },
            "password": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "AV password",
                DefaultFunc: schema.EnvDefaultFunc("ALIENVAULT_PASSWORD", nil),
                Sensitive:   true,
            },
            "skip_tls_verify": {
                Type:        schema.TypeBool,
                Optional:    true,
                Description: "Skip TLS certificate verification",
                DefaultFunc: func() (interface{}, error) {
                    if v := os.Getenv("ALIENVAULT_SKIP_TLS_VERIFY"); v != "" {
                        return v == "true" || v == "1", nil
                    }
                    return false, nil
                },
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

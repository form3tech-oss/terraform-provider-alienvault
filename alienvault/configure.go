package alienvault

import (
    "github.com/form3tech-oss/alienvault"
    "github.com/hashicorp/terraform/helper/schema"
)

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

    var skipTLS bool
    skipTLSRaw, ok := d.GetOk("skip_tls_verify")
    if ok {
        skipTLS = skipTLSRaw.(bool)
    }

    client := alienvault.New(
        d.Get("fqdn").(string),
        alienvault.Credentials{
            Username: d.Get("username").(string),
            Password: d.Get("password").(string),
        },
        skipTLS,
    )

    if err := client.Authenticate(); err != nil {
        return nil, err
    }

    return client, nil
}

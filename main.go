package main

import (
	"github.com/form3tech-oss/terraform-provider-alienvault/alienvault"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return alienvault.Provider()
		},
	})
}

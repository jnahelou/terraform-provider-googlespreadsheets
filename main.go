package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jnahelou/terraform-provider-googlespreadsheets/googlespreadsheets"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return googlespreadsheets.Provider()
		},
	})
}

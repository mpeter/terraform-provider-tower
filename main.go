package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/mpeter/terraform-provider-tower/tower"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: tower.Provider,
	})
}

package tower

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TOWER_ENDPOINT", nil),
				Description: descriptions["endpoint"],
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TOWER_USERNAME", nil),
				Description: descriptions["username"],
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TOWER_PASSWORD", nil),
				Description: descriptions["password"],
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"tower_organization": resourceOrganization(),
			"tower_inventory":    resourceInventory(),
			"tower_host":         resourceHost(),
			"tower_group":        resourceGroup(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"endpoint": "The API Endpoint used to invoke Ansible Tower",
		"username": "The Ansible Tower API Username",
		"password": "The Ansible Tower API Password",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	config := &Config{
		Endpoint: d.Get("endpoint").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	client, err := config.NewClient()
	if err != nil {
		return nil, fmt.Errorf("Error initializing Tower client: %s", err)
	}

	//client.Login()
	//if err != nil {
	//	return nil, fmt.Errorf("Error getting Me from Tower server: %s", err)
	//}
	return client, nil
}

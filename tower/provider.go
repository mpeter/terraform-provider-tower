package tower

import (

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/mapstructure"
	"log"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TOWER_ENDPOINT", "http://localhost/api/v1"),
				Description: descriptions["endpoint"],
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TOWER_USERNAME", "admin"),
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
			"tower_credential":   resourceCredential(),
			"tower_group":        resourceGroup(),
			"tower_host":         resourceHost(),
			"tower_inventory":    resourceInventory(),
			"tower_job_template": resourceJobTemplate(),
			"tower_organization": resourceOrganization(),
			"tower_project":      resourceProject(),
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

	var config Config
	configRaw := d.Get("").(map[string]interface{})
	if err := mapstructure.Decode(configRaw, &config); err != nil {
		return nil, err
	}

	log.Printf("[INFO] Initializing Tower Client" )
	return config.NewClient()
}

package tower

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mpeter/go-towerapi/towerapi"
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
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"tower_organization": resourceOrganization(),
			"tower_inventory":    resourceInventory(),
			"tower_host":         resourceHost(),
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
	var err error
	e := d.Get("endpoint").(string)
	u := d.Get("username").(string)
	p := d.Get("password").(string)
	client, err := towerapi.NewClient(http.DefaultClient, e, u, p)
	if err != nil {
		return nil, err
	}
	return client, nil
}

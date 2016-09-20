package tower

import (
	"fmt"

	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mpeter/go-towerapi/towerapi"
	"github.com/mpeter/go-towerapi/towerapi/hosts"
)

func resourceHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostCreate,
		Read:   resourceHostRead,
		Update: resourceHostUpdate,
		Delete: resourceHostDelete,

		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"inventory_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},

			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"variables": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "",
				StateFunc: normalizeJsonYaml,
			},
		},
	}
}

func resourceHostCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Hosts

	request, err := buildHost(d, meta)
	if err != nil {
		return err
	}
	i, err := service.Create(request)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(i.ID))
	return resourceHostRead(d, meta)
}

func resourceHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Hosts

	r, err := service.GetByID(d.Id())
	if err != nil {
		return fmt.Errorf("Failed to get inventory from Tower API: %v", err)
	}

	d = setHostResourceData(d, r)

	return nil
}

func resourceHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Hosts
	request, err := buildHost(d, client)
	if err != nil {
		return err
	}
	if _, err := service.Update(request); err != nil {
		return err
	}
	return resourceHostRead(d, meta)
}

func resourceHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Hosts
	if err := service.Delete(d.Id()); err != nil {
		return fmt.Errorf("Failed to delete (%s): %s", d.Id(), err)
	}
	return nil
}

func setHostResourceData(d *schema.ResourceData, r *hosts.Host) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("inventory_id", r.Inventory)
	d.Set("enabled", r.Enabled)
	d.Set("instance_id", r.InstanceID)
	d.Set("variables", r.Variables)
	return d
}

func buildHost(d *schema.ResourceData, meta interface{}) (*hosts.Request, error) {

	request := &hosts.Request{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Inventory:   AtoipOr(d.Get("inventory_id").(string), nil),
		Enabled:     d.Get("enabled").(bool),
		InstanceID:  d.Get("instance_id").(string),
		Variables:   normalizeJsonYaml(d.Get("variables").(string)),
	}
	return request, nil
}

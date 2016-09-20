package tower

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mpeter/go-towerapi/towerapi"
	"github.com/mpeter/go-towerapi/towerapi/inventories"
)

func resourceInventory() *schema.Resource {
	return &schema.Resource{
		Create: resourceInventoryCreate,
		Read:   resourceInventoryRead,
		Update: resourceInventoryUpdate,
		Delete: resourceInventoryDelete,

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

			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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

func resourceInventoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Inventories

	request, err := buildInventory(d, meta)
	if err != nil {
		return err
	}
	i, err := service.Create(request)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(i.ID))
	return resourceInventoryRead(d, meta)
}

func resourceInventoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Inventories
	r, err := service.GetByID(d.Id())
	if err != nil {
		return err
	}
	d = setInventoryResourceData(d, r)
	return nil
}

func resourceInventoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Inventories
	request, err := buildInventory(d, client)
	if err != nil {
		return err
	}
	if _, err := service.Update(request); err != nil {
		return err
	}
	return resourceInventoryRead(d, meta)
}

func resourceInventoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Inventories
	if err := service.Delete(d.Id()); err != nil {
		return err
	}
	return nil
}

func setInventoryResourceData(d *schema.ResourceData, r *inventories.Inventory) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("organization_id", r.Organization)
	d.Set("variables", r.Variables)
	return d
}

func buildInventory(d *schema.ResourceData, meta interface{}) (*inventories.Request, error) {
	request := &inventories.Request{
		Name:         d.Get("name").(string),
		Description:  d.Get("description").(string),
		Organization: AtoipOr(d.Get("organization_id").(string), nil),
		Variables:    normalizeJsonYaml(d.Get("variables").(string)),
	}
	return request, nil
}

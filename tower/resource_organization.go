package tower

import (
	"fmt"

	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mpeter/go-towerapi/towerapi"
	"github.com/mpeter/go-towerapi/towerapi/organizations"
)

func resourceOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceOrganizationCreate,
		Read:   resourceOrganizationRead,
		Update: resourceOrganizationUpdate,
		Delete: resourceOrganizationDelete,

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
		},
	}
}

func resourceOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Organizations

	request, err := buildOrganization(d, meta)
	if err != nil {
		return err
	}
	i, err := service.Create(request)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(i.ID))
	return resourceOrganizationRead(d, meta)
}

func resourceOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Organizations

	r, err := service.GetByID(d.Id())
	if err != nil {
		return fmt.Errorf("Failed to get organization from Tower API: %v", err)
	}

	d = setOrganizationResourceData(d, r)

	return nil
}

func resourceOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Organizations
	request, err := buildOrganization(d, client)
	if err != nil {
		return err
	}
	if _, err := service.Update(request); err != nil {
		return err
	}
	return resourceOrganizationRead(d, meta)
}

func resourceOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Organizations
	if err := service.Delete(d.Id()); err != nil {
		return fmt.Errorf("Failed to delete (%s): %s", d.Id(), err)
	}
	return nil
}

func setOrganizationResourceData(d *schema.ResourceData, r *organizations.Organization) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	return d
}

func buildOrganization(d *schema.ResourceData, meta interface{}) (*organizations.Request, error) {

	request := &organizations.Request{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	return request, nil
}

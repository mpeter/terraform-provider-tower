package tower

import (
	"fmt"

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
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
		},
	}
}

func resourceOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Organizations

	name := d.Get("name").(string)
	if i, _, err := service.GetByName(name); err != nil {
		return fmt.Errorf("Failed to get organization from Tower API: %v", err)
	} else if i != nil {
		return fmt.Errorf("Organization named %s already exists, names must be unique", name)
	}
	request := &organizations.Request{
		Name: name,
	}
	if description, ok := d.GetOk("description"); ok {
		request.Description = description.(string)
	}
	if i, _, err := service.Create(request); err != nil {
		return fmt.Errorf("Failed to create organization from Tower API: %v", err)
	} else {
		id := fmt.Sprintf("%d", i.ID)
		d.SetId(id)
		d.Set("name", i.Name)
		d.Set("description", i.Description)
	}
	return nil
}

func resourceOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Organizations
	id := d.Get("id").(int)

	if i, _, err := service.GetByID(id); err != nil {
		return fmt.Errorf("Failed to get organization from Tower API: %v", err)
	} else {
		id := fmt.Sprintf("%d", i.ID)
		d.SetId(id)
		d.Set("name", i.Name)
		d.Set("id", i.ID)
		d.Set("description", i.Description)
	}
	return nil
}

func resourceOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Organizations
	id := d.Get("id").(int)

	request := &organizations.Request{
		Name: d.Get("name").(string),
	}
	if description, ok := d.GetOk("description"); ok {
		request.Description = description.(string)
	}
	if i, _, err := service.Update(id, request); err != nil {
		return fmt.Errorf("Failed to update organization : %v", err)
	} else {
		d.Set("name", i.Name)
		id := fmt.Sprintf("%d", i.ID)
		d.SetId(id)
		d.Set("id", i.ID)
		d.Set("description", i.Description)
	}
	return nil
}

func resourceOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Organizations
	id := d.Get("id").(int)
	if _, err := service.Delete(id); err != nil {
		return fmt.Errorf("Failed to update organization : %v", err)
	} else {
		d.SetId("")
	}
	return nil
}

package tower

import (
	"fmt"

	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mpeter/go-towerapi/towerapi"
	"github.com/mpeter/go-towerapi/towerapi/groups"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,

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

			"variables": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"variables_json": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"variables_yaml"},
				StateFunc:     normalizeJson,
			},

			"variables_yaml": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				StateFunc: normalizeYaml,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Groups

	request, err := buildGroup(d, meta)
	if err != nil {
		return err
	}
	i, err := service.Create(request)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(i.ID))
	return resourceGroupRead(d, meta)
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Groups

	r, err := service.GetByID(d.Id())
	if err != nil {
		return fmt.Errorf("Failed to get inventory from Tower API: %v", err)
	}

	d = setGroupResourceData(d, r)

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Groups
	request, err := buildGroup(d, meta)
	if err != nil {
		return err
	}
	if _, err := service.Update(request); err != nil {
		return err
	}
	return resourceGroupRead(d, meta)
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Groups
	if err := service.Delete(d.Id()); err != nil {
		return fmt.Errorf("Failed to delete (%s): %s", d.Id(), err)
	}
	return nil
}

func setGroupResourceData(d *schema.ResourceData, r *groups.Group) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("inventory_id", r.Inventory)
	d.Set("variables", r.Variables)
	return d
}

func buildGroup(d *schema.ResourceData, meta interface{}) (*groups.Request, error) {

	inv_id, _ := strconv.Atoi(d.Get("inventory_id").(string))
	request := &groups.Request{
		Name:      d.Get("name").(string),
		Inventory: inv_id,
	}

	if variables_json, ok := d.GetOk("variables_json"); ok {
		if variables_yaml, ok := d.GetOk("variables_yaml"); ok {
			return nil, fmt.Errorf("Both variables_json and variables_yaml are set: %v / %v ", variables_json, variables_yaml)
		}
		request.Variables = normalizeJson(variables_json.(string))
	}
	if variables_yaml, ok := d.GetOk("variables_yaml"); ok {
		if variables_json, ok := d.GetOk("variables_json"); ok {
			return nil, fmt.Errorf("Both variables_yaml and variables_json are set: %v / %v ", variables_yaml, variables_json)
		}
		request.Variables = normalizeYaml(variables_yaml.(string))
	}
	if description, ok := d.GetOk("description"); ok {
		request.Description = description.(string)
	}

	return request, nil
}

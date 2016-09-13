package tower

import (
	"fmt"

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

			"inventory": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},

			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},

			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},

			"variables": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"variables_json": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  false,
				StateFunc: normalizeJson,
			},

			"variables_yaml": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  false,
				StateFunc: normalizeYaml,
			},
		},
	}
}

func resourceHostCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Hosts

	name := d.Get("name").(string)
	if i, _, err := service.GetByName(name); err != nil {
		return fmt.Errorf("Failed to get inventory from Tower API: %v", err)
	} else if i != nil {
		return fmt.Errorf("Host named %s already exists, names must be unique", name)
	}
	inv := d.Get("inventory").(int)
	request := &hosts.Request{
		Name:      name,
		Inventory: inv,
	}
	if description, ok := d.GetOk("description"); ok {
		request.Description = description.(string)
	}
	if enabled, ok := d.GetOk("enabled"); ok {
		request.Enabled = enabled.(bool)
	}
	if instanceid, ok := d.GetOk("instance_id"); ok {
		request.InstanceID = instanceid.(string)
	}
	if variables_json, ok := d.GetOk("variables_json"); ok {
		if variables_yaml, ok := d.GetOk("variables_yaml"); ok {
			return fmt.Errorf("Both variables_json and variables_yaml are set: %v / %v ", variables_json, variables_yaml)
		}
		request.Variables = normalizeJson(variables_json.(string))
	}
	if variables_yaml, ok := d.GetOk("variables_yaml"); ok {
		if variables_json, ok := d.GetOk("variables_json"); ok {
			return fmt.Errorf("Both variables_yaml and variables_json are set: %v / %v ", variables_yaml, variables_json)
		}
		request.Variables = normalizeYaml(variables_yaml.(string))
	}
	if i, _, err := service.Create(request); err != nil {
		return fmt.Errorf("Failed to create inventory from Tower API: %v", err)
	} else {
		id := fmt.Sprintf("%d", i.ID)
		d.SetId(id)
		d.Set("name", i.Name)
		d.Set("description", i.Description)
		d.Set("inventory", i.Inventory)
		d.Set("enabled", i.Enabled)
		d.Set("instance_id", i.InstanceID)
		d.Set("variables", i.Variables)
	}
	return nil
}

func resourceHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	inv := client.Hosts
	id := d.Get("id").(int)

	if i, _, err := inv.GetByID(id); err != nil {
		return fmt.Errorf("Failed to get inventory from Tower API: %v", err)
	} else {
		id := fmt.Sprintf("%d", i.ID)
		d.SetId(id)
		d.Set("name", i.Name)
		d.Set("id", i.ID)
		d.Set("description", i.Description)
		d.Set("inventory", i.Inventory)
		d.Set("enabled", i.Enabled)
		d.Set("instance_id", i.InstanceID)
		d.Set("variables", i.Variables)
	}
	return nil
}

func resourceHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	inv := client.Hosts
	id := d.Get("id").(int)

	request := &hosts.Request{
		Name:      d.Get("name").(string),
		Inventory: d.Get("inventory").(int),
	}
	if description, ok := d.GetOk("description"); ok {
		request.Description = description.(string)
	}
	if enabled, ok := d.GetOk("enabled"); ok {
		request.Enabled = enabled.(bool)
	}
	if instanceid, ok := d.GetOk("instance_id"); ok {
		request.InstanceID = instanceid.(string)
	}
	if variables_json, ok := d.GetOk("variables_json"); ok {
		if variables_yaml, ok := d.GetOk("variables_yaml"); ok {
			return fmt.Errorf("Both variables_json and variables_yaml are set: %v / %v ", variables_json, variables_yaml)
		}
		request.Variables = normalizeJson(variables_json.(string))
	}
	if variables_yaml, ok := d.GetOk("variables_yaml"); ok {
		if variables_json, ok := d.GetOk("variables_json"); ok {
			return fmt.Errorf("Both variables_yaml and variables_json are set: %v / %v ", variables_yaml, variables_json)
		}
		request.Variables = normalizeYaml(variables_yaml.(string))
	}
	if i, _, err := inv.Update(id, request); err != nil {
		return fmt.Errorf("Failed to update inventory : %v", err)
	} else {
		d.Set("name", i.Name)
		id := fmt.Sprintf("%d", i.ID)
		d.SetId(id)
		d.Set("id", i.ID)
		d.Set("description", i.Description)
		d.Set("inventory", i.Inventory)
		d.Set("enabled", i.Enabled)
		d.Set("instance_id", i.InstanceID)
		d.Set("variables", i.Variables)
	}
	return nil
}

func resourceHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	inv := client.Hosts
	id := d.Get("id").(int)
	if _, err := inv.Delete(id); err != nil {
		return fmt.Errorf("Failed to update inventory : %v", err)
	} else {
		d.SetId("")
	}
	return nil
}

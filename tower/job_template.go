package tower

import (
	"fmt"

	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mpeter/go-towerapi/towerapi"
	"github.com/mpeter/go-towerapi/towerapi/job_templates"
)

func resourceJobTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobTemplateCreate,
		Read:   resourceJobTemplateRead,
		Update: resourceJobTemplateUpdate,
		Delete: resourceJobTemplateDelete,

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

func resourceJobTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.JobTemplates

	request, err := buildJobTemplate(d, meta)
	if err != nil {
		return err
	}
	i, err := service.Create(request)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(i.ID))
	return resourceJobTemplateRead(d, meta)
}

func resourceJobTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.JobTemplates

	r, err := service.GetByID(d.Id())
	if err != nil {
		return fmt.Errorf("Failed to get inventory from Tower API: %v", err)
	}

	d = setJobTemplateResourceData(d, r)

	return nil
}

func resourceJobTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.JobTemplates
	request, err := buildJobTemplate(d, client)
	if err != nil {
		return err
	}
	if _, err := service.Update(request); err != nil {
		return err
	}
	return resourceJobTemplateRead(d, meta)
}

func resourceJobTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.JobTemplates
	if err := service.Delete(d.Id()); err != nil {
		return fmt.Errorf("Failed to delete (%s): %s", d.Id(), err)
	}
	return nil
}

func setJobTemplateResourceData(d *schema.ResourceData, r *job_templates.JobTemplate) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("inventory_id", r.Inventory)
	d.Set("enabled", r.Enabled)
	d.Set("instance_id", r.InstanceID)
	d.Set("variables", r.Variables)
	return d
}

func buildJobTemplate(d *schema.ResourceData, meta interface{}) (*job_templates.Request, error) {

	inv_id, _ := strconv.Atoi(d.Get("inventory_id").(string))
	request := &job_templates.Request{
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
	if enabled, ok := d.GetOk("enabled"); ok {
		request.Enabled = enabled.(bool)
	}
	if instance_id, ok := d.GetOk("instance_id"); ok {
		request.InstanceID = instance_id.(string)
	}

	return request, nil
}

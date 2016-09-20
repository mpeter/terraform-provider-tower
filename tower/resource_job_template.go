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
			// Run, Check, Scan
			"job_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "One of: run, check, scan",
			},
			"inventory_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"playbook": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"credential_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cloud_credential_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_credential_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"forks": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"limit": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			//0,1,2,3,4,5
			"verbosity": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "One of 0,1,2,3,4,5",
			},
			"extra_vars": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"job_tags": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"force_handlers": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"skip_tags": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"start_at_task": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"host_config_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ask_variables_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_limit_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_tags_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_skip_tags_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_job_type_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_credential_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"survey_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"allow_simultaneous": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

	if r, err := service.GetByID(d.Id()); err != nil {
		return err
	} else {
		d = setJobTemplateResourceData(d, r)
	}
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
	d.Set("allow_simultaneous", r.AllowSimultaneous)
	d.Set("ask_credential_on_launch", r.AskCredentialOnLaunch)
	d.Set("ask_job_type_on_launch", r.AskJobTypeOnLaunch)
	d.Set("ask_limit_on_launch", r.AskLimitOnLaunch)
	d.Set("ask_skip_tags_on_launch", r.AskSkipTagsOnLaunch)
	d.Set("ask_tags_on_launch", r.AskTagsOnLaunch)
	d.Set("ask_variables_on_launch", r.AskVariablesOnLaunch)
	d.Set("cloud_credential_id", r.CloudCredential)
	d.Set("credential_id", r.Credential)
	d.Set("description", r.Description)
	d.Set("description", r.Description)
	d.Set("extra_vars", r.ExtraVars)
	d.Set("force_handlers", r.ForceHandlers)
	d.Set("forks", r.Forks)
	d.Set("host_config_key", r.HostConfigKey)
	d.Set("inventory_id", r.Inventory)
	d.Set("job_tags", r.JobTags)
	d.Set("job_type", r.JobType)
	d.Set("limit", r.Limit)
	d.Set("name", r.Name)
	d.Set("network_credential_id", r.NetworkCredential)
	d.Set("playbook", r.Playbook)
	d.Set("project_id", r.Project)
	d.Set("skip_tags", r.SkipTags)
	d.Set("start_at_task", r.StartAtTask)
	d.Set("survey_enabled", r.SurveyEnabled)
	d.Set("verbosity", r.Verbosity)
	return d
}

func buildJobTemplate(d *schema.ResourceData, meta interface{}) (*job_templates.Request, error) {

	request := &job_templates.Request{
		AllowSimultaneous:     d.Get("allow_simultaneous").(bool),
		AskCredentialOnLaunch: d.Get("ask_credential_on_launch").(bool),
		AskJobTypeOnLaunch:    d.Get("ask_job_type_on_launch").(bool),
		AskLimitOnLaunch:      d.Get("ask_limit_on_launch").(bool),
		AskSkipTagsOnLaunch:   d.Get("ask_skip_tags_on_launch").(bool),
		AskTagsOnLaunch:       d.Get("ask_tags_on_launch").(bool),
		AskVariablesOnLaunch:  d.Get("ask_variables_on_launch").(bool),
		CloudCredential:       AtoipOr(d.Get("cloud_credential_id").(string), nil),
		Credential:            AtoipOr(d.Get("credential_id").(string), nil),
		Description:           d.Get("description").(string),
		ExtraVars:             d.Get("extra_vars").(string),
		ForceHandlers:         d.Get("force_handlers").(bool),
		Forks:                 d.Get("forks").(int),
		HostConfigKey:         d.Get("host_config_key").(string),
		Inventory:             AtoipOr(d.Get("inventory_id").(string), nil),
		JobTags:               d.Get("job_tags").(string),
		JobType:               d.Get("job_type").(string),
		Limit:                 d.Get("limit").(string),
		Name:                  d.Get("name").(string),
		NetworkCredential:     AtoipOr(d.Get("network_credential_id").(string), nil),
		Playbook:              d.Get("playbook").(string),
		Project:               AtoipOr(d.Get("project_id").(string), nil),
		SkipTags:              d.Get("skip_tags").(string),
		StartAtTask:           d.Get("start_at_task").(string),
		SurveyEnabled:         d.Get("survey_enabled").(bool),
		Verbosity:             d.Get("verbosity").(int),
	}

	return request, nil
}

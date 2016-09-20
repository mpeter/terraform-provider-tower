package tower

import (
	"fmt"

	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mpeter/go-towerapi/towerapi"
	"github.com/mpeter/go-towerapi/towerapi/projects"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Name of this project",
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "Optional description of this project.",
			},

			"local_path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "Local path (relative to PROJECTS_ROOT) containing playbooks and related files for this project.",
			},

			"scm_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "One of \"\" (manual), git, hg, svn",
			},

			"scm_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "",
			},

			"scm_branch": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				Description: "Specific branch, tag or commit to checkout.",
			},
			"scm_clean": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"scm_delete_on_update": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"credential_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"scm_update_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: false,
			},
			"scm_update_cache_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default: 0,
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Projects

	request, err := buildProject(d, meta)
	if err != nil {
		return err
	}
	i, err := service.Create(request)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(i.ID))
	return resourceProjectRead(d, meta)
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Projects

	r, err := service.GetByID(d.Id())
	if err != nil {
		return fmt.Errorf("Failed to get project from Tower API: %v", err)
	}

	d = setProjectResourceData(d, r)

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Projects
	request, err := buildProject(d, client)
	if err != nil {
		return err
	}
	if _, err := service.Update(request); err != nil {
		return err
	}
	return resourceProjectRead(d, meta)
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Projects
	if err := service.Delete(d.Id()); err != nil {
		return fmt.Errorf("Failed to delete (%s): %s", d.Id(), err)
	}
	return nil
}

func setProjectResourceData(d *schema.ResourceData, r *projects.Project) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("local_path", r.LocalPath)
	d.Set("scm_type", r.ScmType)
	d.Set("scm_url", r.ScmURL)
	d.Set("scm_branch", r.ScmBranch)
	d.Set("scm_clean", r.ScmClean)
	d.Set("scm_delete_on_update", r.ScmDeleteOnUpdate)
	d.Set("credential_id", r.Credential)
	d.Set("organization_id", r.Organization)
	d.Set("scm_update_on_launch", r.ScmUpdateOnLaunch)
	d.Set("scm_update_cache_timeout", r.ScmUpdateCacheTimeout)
	return d
}

func buildProject(d *schema.ResourceData, meta interface{}) (*projects.Request, error) {

	request := &projects.Request{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		LocalPath:         d.Get("local_path").(string),
		ScmType:           d.Get("scm_type").(string),
		ScmURL:            d.Get("scm_url").(string),
		ScmBranch:         d.Get("scm_branch").(string),
		ScmClean:          d.Get("scm_clean").(bool),
		ScmDeleteOnUpdate: d.Get("scm_delete_on_update").(bool),
		Credential:        AtoipOr(d.Get("credential_id").(string), nil),
		Organization:      AtoipOr(d.Get("organization_id").(string), nil),
		ScmUpdateOnLaunch: d.Get("scm_update_on_launch").(bool),
		ScmUpdateCacheTimeout: d.Get("scm_update_cache_timeout").(int),
	}

	return request, nil
}

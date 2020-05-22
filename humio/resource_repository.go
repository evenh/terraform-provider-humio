package humio

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	humio "github.com/humio/cli/api"
)

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceRepositoryCreate,
		Read:   resourceRepositoryRead,
		Update: resourceRepositoryUpdate,
		Delete: resourceRepositoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"allow_data_deletion": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"retention": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_size_in_gb": {
							Type:         schema.TypeFloat,
							Optional:     true,
							ValidateFunc: validation.FloatAtLeast(0),
						},
						"ingest_size_in_gb": {
							Type:         schema.TypeFloat,
							Optional:     true,
							ValidateFunc: validation.FloatAtLeast(0),
						},
						"time_in_days": {
							Type:         schema.TypeFloat,
							Optional:     true,
							ValidateFunc: validation.FloatAtLeast(0),
						},
					},
				},
			},
		},
	}
}

func resourceRepositoryCreate(d *schema.ResourceData, client interface{}) error {
	err := client.(*humio.Client).Repositories().Create(
		d.Get("name").(string),
	)
	if err != nil {
		return fmt.Errorf("could not create repository: %v", err)
	}

	err = client.(*humio.Client).Repositories().UpdateDescription(
		d.Get("name").(string),
		d.Get("description").(string),
	)
	if err != nil {
		return fmt.Errorf("could not set description for repository: %v", err)
	}
	retention := d.Get("retention").(*schema.Set).List()[0].(tfMap)
	err = client.(*humio.Client).Repositories().UpdateTimeBasedRetention(
		d.Get("name").(string),
		retention["time_in_days"].(float64),
		d.Get("allow_data_deletion").(bool),
	)
	if err != nil {
		return fmt.Errorf("could not set time based retention for repository: %v", err)
	}
	err = client.(*humio.Client).Repositories().UpdateIngestBasedRetention(
		d.Get("name").(string),
		retention["ingest_size_in_gb"].(float64),
		d.Get("allow_data_deletion").(bool),
	)
	if err != nil {
		return fmt.Errorf("could not set time based retention for repository: %v", err)
	}
	err = client.(*humio.Client).Repositories().UpdateStorageBasedRetention(
		d.Get("name").(string),
		retention["storage_size_in_gb"].(float64),
		d.Get("allow_data_deletion").(bool),
	)
	if err != nil {
		return fmt.Errorf("could not set time based retention for repository: %v", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceRepositoryRead(d, client)
}

func resourceRepositoryRead(d *schema.ResourceData, client interface{}) error {
	repo, err := client.(*humio.Client).Repositories().Get(d.Id())
	if err != nil {
		return fmt.Errorf("could not get repository: %v", err)
	}
	resourceDataFromRepo(&repo, d)
	return nil
}

func resourceDataFromRepo(a *humio.Repository, d *schema.ResourceData) error {
	d.Set("name", a.Name)
	d.Set("description", a.Description)
	if err := d.Set("retention", retentionFromRepository(a)); err != nil {
		return fmt.Errorf("error setting retention settings for resource %s: %s", d.Id(), err)
	}
	return nil
}

func retentionFromRepository(a *humio.Repository) []tfMap {
	s := tfMap{}
	s["time_in_days"] = a.RetentionDays
	s["ingest_size_in_gb"] = a.IngestRetentionSizeGB
	s["storage_size_in_gb"] = a.StorageRetentionSizeGB
	return []tfMap{s}
}

func resourceRepositoryUpdate(d *schema.ResourceData, client interface{}) error {
	err := client.(*humio.Client).Repositories().UpdateDescription(
		d.Get("name").(string),
		d.Get("description").(string),
	)
	if err != nil {
		return fmt.Errorf("could not set description for repository: %v", err)
	}
	retention := d.Get("retention").(*schema.Set).List()[0].(tfMap)
	err = client.(*humio.Client).Repositories().UpdateTimeBasedRetention(
		d.Get("name").(string),
		retention["time_in_days"].(float64),
		d.Get("allow_data_deletion").(bool),
	)
	if err != nil {
		return fmt.Errorf("could not set time based retention for repository: %v", err)
	}
	err = client.(*humio.Client).Repositories().UpdateIngestBasedRetention(
		d.Get("name").(string),
		retention["ingest_size_in_gb"].(float64),
		d.Get("allow_data_deletion").(bool),
	)
	if err != nil {
		return fmt.Errorf("could not set time based retention for repository: %v", err)
	}
	err = client.(*humio.Client).Repositories().UpdateStorageBasedRetention(
		d.Get("name").(string),
		retention["storage_size_in_gb"].(float64),
		d.Get("allow_data_deletion").(bool),
	)
	if err != nil {
		return fmt.Errorf("could not set time based retention for repository: %v", err)
	}

	return resourceRepositoryRead(d, client)
}

func resourceRepositoryDelete(d *schema.ResourceData, client interface{}) error {
	deleteReason := "Deleted by Terraform"
	if err := client.(*humio.Client).Repositories().Delete(
		d.Get("name").(string),
		deleteReason,
		d.Get("allow_data_deletion").(bool),
	); err != nil {
		return fmt.Errorf("could not delete repository: %v", err)
	}
	return nil
}

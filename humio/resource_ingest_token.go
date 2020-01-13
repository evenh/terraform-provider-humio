package humio

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	humio "github.com/humio/cli/api"
)

func resourceIngestToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceIngestTokenCreate,
		Read:   resourceIngestTokenRead,
		Update: resourceIngestTokenUpdate,
		Delete: resourceIngestTokenDelete,

		Schema: map[string]*schema.Schema{
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parser": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIngestTokenCreate(d *schema.ResourceData, client interface{}) error {
	ingestToken, err := client.(*humio.Client).IngestTokens().Add(d.Get("repository").(string), d.Get("name").(string), d.Get("parser").(string))
	if err != nil {
		return fmt.Errorf("could not create ingest token: %v", err)
	}
	d.SetId(ingestToken.Name)
	return resourceIngestTokenRead(d, client)
}

func resourceIngestTokenRead(d *schema.ResourceData, client interface{}) error {
	ingestToken, err := client.(*humio.Client).IngestTokens().Get(d.Get("repository").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("could not get ingest token: %v", err)
	}
	resourceDataFromIngestToken(ingestToken, d)
	return nil
}

func resourceDataFromIngestToken(a *humio.IngestToken, d *schema.ResourceData) error {
	d.Set("name", a.Name)
	d.Set("token", a.Token)
	d.Set("parser", a.AssignedParser)
	d.SetId(d.Id())
	return nil
}

func resourceIngestTokenUpdate(d *schema.ResourceData, client interface{}) error {
	ingestToken, err := client.(*humio.Client).IngestTokens().Add(d.Get("repository").(string), d.Get("name").(string), d.Get("parser").(string))
	if err != nil {
		return fmt.Errorf("could not create ingest token: %v", err)
	}
	d.SetId(ingestToken.Name)
	return resourceAlertRead(d, client)
}

func resourceIngestTokenDelete(d *schema.ResourceData, client interface{}) error {
	if err := client.(*humio.Client).IngestTokens().Remove(d.Get("repository").(string), d.Get("name").(string)); err != nil {
		return fmt.Errorf("could not delete ingest token: %v", err)
	}
	return nil
}

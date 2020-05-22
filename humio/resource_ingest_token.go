// Copyright © 2020 Humio Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceIngestTokenCreate(d *schema.ResourceData, client interface{}) error {
	_, err := client.(*humio.Client).IngestTokens().Add(d.Get("repository").(string), d.Get("name").(string), d.Get("parser").(string))
	if err != nil {
		return fmt.Errorf("could not create ingest token: %v", err)
	}
	d.SetId(fmt.Sprintf("%s+%s", d.Get("repository"), d.Get("name")))

	return resourceIngestTokenRead(d, client)
}

func resourceIngestTokenRead(d *schema.ResourceData, client interface{}) error {
	// If we don't have a repository when importing, we parse it from the ID.
	if _, ok := d.GetOk("repository"); !ok {
		parts := parseRepositoryAndName(d.Id())
		//we check that we have parsed the id into the correct number of segments
		if parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("Error Importing humio_ingest_token. Please make sure the ID is in the form REPOSITORYNAME+INGESTTOKENNAMENAME (i.e. myRepoName+myIngestTokenName")
		}
		d.Set("repository", parts[0])
		d.Set("name", parts[1])
	}

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
	return nil
}

func resourceIngestTokenUpdate(d *schema.ResourceData, client interface{}) error {
	_, err := client.(*humio.Client).IngestTokens().Add(d.Get("repository").(string), d.Get("name").(string), d.Get("parser").(string))
	if err != nil {
		return fmt.Errorf("could not create ingest token: %v", err)
	}
	return resourceAlertRead(d, client)
}

func resourceIngestTokenDelete(d *schema.ResourceData, client interface{}) error {
	if err := client.(*humio.Client).IngestTokens().Remove(d.Get("repository").(string), d.Get("name").(string)); err != nil {
		return fmt.Errorf("could not delete ingest token: %v", err)
	}
	return nil
}

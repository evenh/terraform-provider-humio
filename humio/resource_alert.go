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

func resourceAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlertCreate,
		Read:   resourceAlertRead,
		Update: resourceAlertUpdate,
		Delete: resourceAlertDelete,

		Schema: map[string]*schema.Schema{
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"link_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"silenced": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"throttle_time_millis": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"start": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notifiers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				//ForceNew: true, // TODO(mike): figure out why apply causing an in-place update fails. running apply again works!?
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlertCreate(d *schema.ResourceData, client interface{}) error {
	alert, err := alertFromResourceData(d, client)
	if err != nil {
		return fmt.Errorf("could not obtain alert from resource data: %v", err)
	}

	n, err := client.(*humio.Client).Alerts().Add(d.Get("repository").(string), &alert, false)
	if err != nil {
		return fmt.Errorf("could not create alert: %v", err)
	}
	d.SetId(n.ID)

	return resourceAlertRead(d, client)
}

func resourceAlertRead(d *schema.ResourceData, client interface{}) error {
	alert, err := client.(*humio.Client).Alerts().Get(d.Get("repository").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("could not get alert: %v", err)
	}
	resourceDataFromAlert(alert, d)
	return nil
}

func resourceDataFromAlert(a *humio.Alert, d *schema.ResourceData) error {
	d.Set("name", a.Name)
	d.Set("description", a.Description)
	d.Set("throttle_time_millis", a.ThrottleTimeMillis)
	d.Set("silenced", a.Silenced)
	d.Set("notifiers", a.Notifiers)
	d.Set("link_url", a.LinkURL)
	d.Set("labels", a.Labels)
	d.Set("query", a.Query.QueryString)
	d.Set("start", a.Query.Start)
	d.SetId(d.Id())
	return nil
}

func resourceAlertUpdate(d *schema.ResourceData, client interface{}) error {
	alert, err := alertFromResourceData(d, client)
	if err != nil {
		return fmt.Errorf("could not obtain alert from resource data: %v", err)
	}

	_, err = client.(*humio.Client).Alerts().Add(d.Get("repository").(string), &alert, true)
	if err != nil {
		return fmt.Errorf("could not create alert: %v", err)
	}
	d.SetId(alert.ID)

	// TODO(mike): Updating a resource with a different set of notifiers updates the state to not contain the new notifiers
	return resourceAlertRead(d, client)
}

func resourceAlertDelete(d *schema.ResourceData, client interface{}) error {
	if err := client.(*humio.Client).Alerts().Delete(d.Get("repository").(string), d.Get("name").(string)); err != nil {
		return fmt.Errorf("could not delete alert: %v", err)
	}
	return nil
}

func alertFromResourceData(d *schema.ResourceData, client interface{}) (humio.Alert, error) {
	return humio.Alert{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		ThrottleTimeMillis: d.Get("throttle_time_millis").(int),
		Silenced:           d.Get("silenced").(bool),
		Notifiers:          convertInterfaceListToStringSlice(d.Get("notifiers").([]interface{})),
		LinkURL:            d.Get("link_url").(string),
		Labels:             convertInterfaceListToStringSlice(d.Get("labels").([]interface{})),
		Query: humio.HumioQuery{
			QueryString: d.Get("query").(string),
			Start:       d.Get("start").(string),
			End:         "now",
			IsLive:      true,
		},
	}, nil
}

func convertInterfaceListToStringSlice(s []interface{}) []string {
	var element []string
	for _, item := range s {
		value, _ := item.(string)
		element = append(element, value)
	}
	return element
}

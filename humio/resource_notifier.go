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
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	humio "github.com/humio/cli/api"
)

var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func resourceNotifier() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotifierCreate,
		Read:   resourceNotifierRead,
		Update: resourceNotifierUpdate,
		Delete: resourceNotifierDelete,
		/*
			Importer: &schema.ResourceImporter{
				State: schema.ImportStatePassthrough,
			},
		*/

		Schema: map[string]*schema.Schema{
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entity": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"opsgenie", "pagerduty", "slack", "victorops", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"body_template": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"recipients": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
									v := val.(string)
									if len(v) > 254 || !rxEmail.MatchString(v) {
										errs = append(errs, fmt.Errorf("%q must be a valid email, got: %s", key, v))
									}
									return warns, errs
								},
							},
						},
						"subject_template": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"opsgenie": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "pagerduty", "slack", "victorops", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateURL,
						},
						"genie_key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"pagerduty": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "opsgenie", "slack", "victorops", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"routing_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"critical",
								"error",
								"warning",
								"info",
							}, false),
						},
					},
				},
			},
			"slack": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "opsgenie", "pagerduty", "victorops", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fields": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type:     schema.TypeString,
								Required: true,
							},
						},
						"url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateURL,
						},
					},
				},
			},
			"victorops": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "opsgenie", "pagerduty", "slack", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"notify_url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateURL,
						},
					},
				},
			},
			"webhook": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "opsgenie", "pagerduty", "slack", "victorops"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"body_template": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "{\n  \"repository\": \"{repo_name}\",\n  \"timestamp\": \"{alert_triggered_timestamp}\",\n  \"alert\": {\n    \"name\": \"{alert_name}\",\n    \"description\": \"{alert_description}\",\n    \"query\": {\n      \"queryString\": \"{query_string} \",\n      \"end\": \"{query_time_end}\",\n      \"start\": \"{query_time_start}\"\n    },\n    \"notifierID\": \"{alert_notifier_id}\",\n    \"id\": \"{alert_id}\"  },\n  \"warnings\": \"{warnings}\",\n  \"events\": {events},\n  \"numberOfEvents\": {event_count}\n}",
						},
						"headers": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"method": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "POST",
							ValidateFunc: validation.StringInSlice([]string{
								http.MethodConnect,
								http.MethodDelete,
								http.MethodGet,
								http.MethodHead,
								http.MethodOptions,
								http.MethodPatch,
								http.MethodPost,
								http.MethodPut,
								http.MethodTrace,
							}, false),
						},
						"url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateURL,
						},
					},
				},
			},
		},
	}
}

func resourceNotifierCreate(d *schema.ResourceData, client interface{}) error {
	notifier, err := notifierFromResourceData(d, client)
	if err != nil {
		return fmt.Errorf("could not obtain notifier from resource data: %v", err)
	}

	n, err := client.(*humio.Client).Notifiers().Add(d.Get("repository").(string), &notifier, false)
	if err != nil {
		return fmt.Errorf("could not create notifier: %v", err)
	}
	d.SetId(n.ID)

	return resourceNotifierRead(d, client)
}

func resourceNotifierRead(d *schema.ResourceData, client interface{}) error {
	// TODO: to fix import functionality, we must ensure the user can provide the notifier ID, which we use to look up repository name and alert name
	notifier, err := client.(*humio.Client).Notifiers().Get(d.Get("repository").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("could not get notifier: %v", err)
	}
	resourceDataFromNotifier(notifier, d)
	return nil
}

func resourceDataFromNotifier(n *humio.Notifier, d *schema.ResourceData) error {
	d.Set("name", n.Name)
	d.Set("entity", n.Entity)
	for k, v := range n.Properties {
		d.Set(k, v)
	}
	return nil
}

func resourceNotifierUpdate(d *schema.ResourceData, client interface{}) error {
	notifier, err := notifierFromResourceData(d, client)
	if err != nil {
		return fmt.Errorf("could not obtain notifier from resource data: %v", err)
	}

	_, err = client.(*humio.Client).Notifiers().Add(d.Get("repository").(string), &notifier, true)
	if err != nil {
		return fmt.Errorf("could not create notifier: %v", err)
	}

	return resourceNotifierRead(d, client)
}

func resourceNotifierDelete(d *schema.ResourceData, client interface{}) error {
	if err := client.(*humio.Client).Notifiers().Delete(d.Get("repository").(string), d.Get("name").(string)); err != nil {
		return fmt.Errorf("could not delete notifier: %v", err)
	}
	return nil
}

func notifierFromResourceData(d *schema.ResourceData, client interface{}) (humio.Notifier, error) {
	notifier := humio.Notifier{
		Entity: d.Get("entity").(string),
		Name:   d.Get("name").(string),
	}

	switch d.Get("entity") {
	case humio.NotifierTypeEmail:
		notifier.Properties = notifierPropertiesFromSet(d.Get("email").(*schema.Set))
	case humio.NotifierTypeOpsGenie:
		notifier.Properties = notifierPropertiesFromSet(d.Get("opsgenie").(*schema.Set))
	case humio.NotifierTypePagerDuty:
		notifier.Properties = notifierPropertiesFromSet(d.Get("pagerduty").(*schema.Set))
	case humio.NotifierTypeSlack:
		notifier.Properties = notifierPropertiesFromSet(d.Get("slack").(*schema.Set))
	case humio.NotifierTypeVictorOps:
		notifier.Properties = notifierPropertiesFromSet(d.Get("victorops").(*schema.Set))
	case humio.NotifierTypeWebHook:
		notifier.Properties = notifierPropertiesFromSet(d.Get("webhook").(*schema.Set))
	default:
		return humio.Notifier{}, fmt.Errorf("unsupported notifier entity: %s", d.Get("entity"))
	}
	return notifier, nil
}

func notifierPropertiesFromSet(s *schema.Set) map[string]interface{} {
	if s.Len() == 0 {
		return map[string]interface{}{}
	}
	res := s.List()[0].(tfMap)
	allConfigurations := map[string]interface{}{
		"bodyTemplate":    res["body_template"],
		"recipients":      res["recipients"],
		"subjectTemplate": res["subject_template"],
		"apiUrl":          res["api_url"],
		"genieKey":        res["genie_key"],
		"routingKey":      res["routing_key"],
		"severity":        res["severity"],
		"url":             res["url"],
		"fields":          res["fields"],
		"messageType":     res["message_type"],
		"notifyUrl":       res["notify_url"],
		"method":          res["method"],
		"headers":         res["headers"],
	}

	configurationsSet := map[string]interface{}{}
	for key, value := range allConfigurations {
		if value != "" {
			configurationsSet[key] = value
		}
	}
	return configurationsSet
}

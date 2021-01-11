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
	"context"
	"fmt"
	"net/http"
	"reflect"
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	humio "github.com/humio/cli/api"
)

var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func resourceNotifier() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotifierCreate,
		ReadContext:   resourceNotifierRead,
		UpdateContext: resourceNotifierUpdate,
		DeleteContext: resourceNotifierDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"notifier_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entity": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					humio.NotifierTypeEmail,
					humio.NotifierTypeHumioRepo,
					humio.NotifierTypeOpsGenie,
					humio.NotifierTypePagerDuty,
					humio.NotifierTypeSlack,
					humio.NotifierTypeSlackPostMessage,
					humio.NotifierTypeVictorOps,
					humio.NotifierTypeWebHook,
				}, false)),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"humiorepo", "opsgenie", "pagerduty", "slack", "slackpostmessage", "victorops", "webhook"},
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
								ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
									v := val.(string)
									if len(v) > 254 || !rxEmail.MatchString(v) {
										return diag.FromErr(fmt.Errorf("%q must be a valid email, got: %s", key, v))
									}
									return nil
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
			"humiorepo": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "opsgenie", "pagerduty", "slack", "slackpostmessage", "victorops", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ingest_token": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"opsgenie": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "humiorepo", "pagerduty", "slack", "slackpostmessage", "victorops", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_url": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "https://api.opsgenie.com",
							ValidateDiagFunc: validateURL,
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
				ConflictsWith: []string{"email", "humiorepo", "opsgenie", "slack", "slackpostmessage", "victorops", "webhook"},
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
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
								"critical",
								"error",
								"warning",
								"info",
							}, false)),
						},
					},
				},
			},
			"slack": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "humiorepo", "opsgenie", "pagerduty", "slackpostmessage", "victorops", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fields": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"url": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateURL,
						},
					},
				},
			},
			"slackpostmessage": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "humiorepo", "opsgenie", "pagerduty", "slack", "victorops", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_token": {
							Type:     schema.TypeString,
							Required: true,
						},
						"channels": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"fields": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"use_proxy": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"victorops": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "humiorepo", "opsgenie", "pagerduty", "slack", "slackpostmessage", "webhook"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"notify_url": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateURL,
						},
					},
				},
			},
			"webhook": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				ConflictsWith: []string{"email", "humiorepo", "opsgenie", "pagerduty", "slack", "slackpostmessage", "victorops"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"body_template": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "{\n  \"repository\": \"{repo_name}\",\n  \"timestamp\": \"{alert_triggered_timestamp}\",\n  \"alert\": {\n    \"name\": \"{alert_name}\",\n    \"description\": \"{alert_description}\",\n    \"query\": {\n      \"queryString\": \"{query_string} \",\n      \"end\": \"{query_time_end}\",\n      \"start\": \"{query_time_start}\"\n    },\n    \"notifierID\": \"{alert_notifier_id}\",\n    \"id\": \"{alert_id}\"\n  },\n  \"warnings\": \"{warnings}\",\n  \"events\": {events},\n  \"numberOfEvents\": {event_count}\n  }",
						},
						"headers": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"method": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "POST",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
								http.MethodGet,
								http.MethodPost,
								http.MethodPut,
							}, false)),
						},
						"url": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateURL,
						},
					},
				},
			},
		},
	}
}

func resourceNotifierCreate(ctx context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	notifier, err := notifierFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain notifier from resource data: %s", err)
	}

	n, err := client.(*humio.Client).Notifiers().Add(
		d.Get("repository").(string),
		&notifier,
		false,
	)
	if err != nil {
		return diag.Errorf("could not create notifier: %s", err)
	}
	d.SetId(fmt.Sprintf("%s+%s", d.Get("repository").(string), n.Name))

	return resourceNotifierRead(ctx, d, client)
}

func resourceNotifierRead(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	parts := parseRepositoryAndID(d.Id())
	// If we don't have a repository when importing, we parse it from the ID.
	if _, ok := d.GetOk("repository"); !ok {
		//we check that we have parsed the id into the correct number of segments
		if parts[0] == "" || parts[1] == "" {
			return diag.Errorf("error importing humio_notifier. Please make sure the ID is in the form REPOSITORYNAME+NOTIFIERID (i.e. myRepoName+12345678901234567890123456789012")
		}
		err := d.Set("repository", parts[0])
		if err != nil {
			return diag.Errorf("error setting repository for resource %s: %s", d.Id(), err)
		}
		err = d.Set("name", parts[1])
		if err != nil {
			return diag.Errorf("error setting name for resource %s: %s", d.Id(), err)
		}
	}

	notifier, err := client.(*humio.Client).Notifiers().Get(
		d.Get("repository").(string),
		d.Get("name").(string),
	)
	if err != nil || reflect.DeepEqual(*notifier, humio.Notifier{}) {
		return diag.Errorf("could not get notifier: %s", err)
	}
	return resourceDataFromNotifier(notifier, d)
}

func resourceDataFromNotifier(n *humio.Notifier, d *schema.ResourceData) diag.Diagnostics {
	err := d.Set("notifier_id", n.ID)
	if err != nil {
		return diag.Errorf("could not set notifier_id for notifier: %s", err)
	}
	err = d.Set("name", n.Name)
	if err != nil {
		return diag.Errorf("could not set name for notifier: %s", err)
	}
	err = d.Set("entity", n.Entity)
	if err != nil {
		return diag.Errorf("could not set entity for notifier: %s", err)
	}

	switch n.Entity {
	case humio.NotifierTypeEmail:
		if err := d.Set("email", emailFromNotifier(n)); err != nil {
			return diag.Errorf("error setting email settings for resource %s: %s", d.Id(), err)
		}
	case humio.NotifierTypeHumioRepo:
		if err := d.Set("humiorepo", humiorepoFromNotifier(n)); err != nil {
			return diag.Errorf("error setting humiorepo settings for resource %s: %s", d.Id(), err)
		}
	case humio.NotifierTypeOpsGenie:
		if err := d.Set("opsgenie", opsgenieFromNotifier(n)); err != nil {
			return diag.Errorf("error setting opsgenie settings for resource %s: %s", d.Id(), err)
		}
	case humio.NotifierTypePagerDuty:
		if err := d.Set("pagerduty", pagerdutyFromNotifier(n)); err != nil {
			return diag.Errorf("error setting pagerduty settings for resource %s: %s", d.Id(), err)
		}
	case humio.NotifierTypeSlack:
		if err := d.Set("slack", slackFromNotifier(n)); err != nil {
			return diag.Errorf("error setting slack settings for resource %s: %s", d.Id(), err)
		}
	case humio.NotifierTypeSlackPostMessage:
		if err := d.Set("slackpostmessage", slackpostmessageFromNotifier(n)); err != nil {
			return diag.Errorf("error setting slackpostmessage settings for resource %s: %s", d.Id(), err)
		}
	case humio.NotifierTypeVictorOps:
		if err := d.Set("victorops", victoropsFromNotifier(n)); err != nil {
			return diag.Errorf("error setting victorops settings for resource %s: %s", d.Id(), err)
		}
	case humio.NotifierTypeWebHook:
		if err := d.Set("webhook", webhookFromNotifier(n)); err != nil {
			return diag.Errorf("error setting webhook settings for resource %s: %s", d.Id(), err)
		}
	default:
		return diag.Errorf("unsupported notifier entity: %s", n.Entity)
	}

	return nil
}

func resourceNotifierUpdate(ctx context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	notifier, err := notifierFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain notifier from resource data: %s", err)
	}

	_, err = client.(*humio.Client).Notifiers().Add(
		d.Get("repository").(string),
		&notifier,
		true,
	)
	if err != nil {
		return diag.Errorf("could not update notifier: %s", err)
	}

	return resourceNotifierRead(ctx, d, client)
}

func resourceNotifierDelete(_ context.Context, d *schema.ResourceData, client interface{}) diag.Diagnostics {
	notifier, err := notifierFromResourceData(d)
	if err != nil {
		return diag.Errorf("could not obtain notifier from resource data: %s", err)
	}

	err = client.(*humio.Client).Notifiers().Delete(
		d.Get("repository").(string),
		notifier.Name,
	)
	if err != nil {
		return diag.Errorf("could not delete notifier: %s", err)
	}
	return nil
}

// notifierFromResourceData returns a humio.Notifier based on either the new change or the current state depending on update bool.
func notifierFromResourceData(d *schema.ResourceData) (humio.Notifier, error) {
	notifier := humio.Notifier{
		ID:     d.Get("notifier_id").(string),
		Entity: d.Get("entity").(string),
		Name:   d.Get("name").(string),
	}

	switch d.Get("entity") {
	case humio.NotifierTypeEmail:
		properties := getNotifierPropertiesFromResourceData(d, "email", "recipients")
		notifier.Properties = map[string]interface{}{
			"recipients": properties[0]["recipients"].([]interface{}),
		}
		if properties[0]["body_template"].(string) != "" {
			notifier.Properties["bodyTemplate"] = properties[0]["body_template"].(string)
		}
		if properties[0]["subject_template"].(string) != "" {
			notifier.Properties["subjectTemplate"] = properties[0]["subject_template"].(string)
		}
	case humio.NotifierTypeHumioRepo:
		properties := getNotifierPropertiesFromResourceData(d, "humiorepo", "ingest_token")
		notifier.Properties = map[string]interface{}{
			"ingestToken": properties[0]["ingest_token"].(string),
		}
	case humio.NotifierTypeOpsGenie:
		properties := getNotifierPropertiesFromResourceData(d, "opsgenie", "genie_key")
		notifier.Properties = map[string]interface{}{
			"apiUrl":   properties[0]["api_url"].(string),
			"genieKey": properties[0]["genie_key"].(string),
		}
	case humio.NotifierTypePagerDuty:
		properties := getNotifierPropertiesFromResourceData(d, "pagerduty", "routing_key")
		notifier.Properties = map[string]interface{}{
			"routingKey": properties[0]["routing_key"].(string),
			"severity":   properties[0]["severity"].(string),
		}
	case humio.NotifierTypeSlack:
		properties := getNotifierPropertiesFromResourceData(d, "slack", "url")
		notifier.Properties = map[string]interface{}{
			"url":    properties[0]["url"].(string),
			"fields": properties[0]["fields"].(map[string]interface{}),
		}
	case humio.NotifierTypeSlackPostMessage:
		properties := getNotifierPropertiesFromResourceData(d, "slackpostmessage", "api_token")
		notifier.Properties = map[string]interface{}{
			"apiToken": properties[0]["api_token"].(string),
			"channels": properties[0]["channels"].([]interface{}),
			"fields":   properties[0]["fields"].(map[string]interface{}),
			"useProxy": properties[0]["use_proxy"].(bool),
		}
	case humio.NotifierTypeVictorOps:
		properties := getNotifierPropertiesFromResourceData(d, "victorops", "notify_url")
		notifier.Properties = map[string]interface{}{
			"messageType": properties[0]["message_type"].(string),
			"notifyUrl":   properties[0]["notify_url"].(string),
		}
	case humio.NotifierTypeWebHook:
		properties := getNotifierPropertiesFromResourceData(d, "webhook", "url")
		notifier.Properties = map[string]interface{}{
			"bodyTemplate": properties[0]["body_template"].(string),
			"headers":      properties[0]["headers"].(map[string]interface{}),
			"method":       properties[0]["method"].(string),
			"url":          properties[0]["url"].(string),
		}
	default:
		return humio.Notifier{}, fmt.Errorf("unsupported notifier entity: %s", d.Get("entity"))
	}

	return notifier, nil
}

// getNotifierPropertiesFromResourceData returns the first non-empty set of notifier properties related to a given notifier.
// We do this as a workaround for an issue where we get a list longer than 1 which should not happen given MaxItems is
// set to 1 in the schema definition.
func getNotifierPropertiesFromResourceData(d *schema.ResourceData, notifierName, requiredPropertyName string) []tfMap {
	_, newProperties := d.GetChange(notifierName)
	newPropertiesList := newProperties.(*schema.Set).List()
	if len(newPropertiesList) == 0 {
		properties := d.Get(notifierName).(*schema.Set).List()[0]
		return []tfMap{properties.(tfMap)}
	}
	for idx := range newPropertiesList {
		if newPropertiesList[idx].(tfMap)[requiredPropertyName] != "" {
			return []tfMap{newPropertiesList[idx].(tfMap)}
		}
	}

	return []tfMap{}
}

func emailFromNotifier(n *humio.Notifier) []tfMap {
	s := tfMap{}
	s["recipients"] = n.Properties["recipients"]
	if n.Properties["bodyTemplate"] != nil {
		s["body_template"] = n.Properties["bodyTemplate"]
	}
	if n.Properties["subjectTemplate"] != nil {
		s["subject_template"] = n.Properties["subjectTemplate"]
	}
	return []tfMap{s}
}

func humiorepoFromNotifier(n *humio.Notifier) []tfMap {
	s := tfMap{}
	s["ingest_token"] = n.Properties["ingestToken"]
	return []tfMap{s}
}

func opsgenieFromNotifier(n *humio.Notifier) []tfMap {
	s := tfMap{}
	s["api_url"] = n.Properties["apiUrl"]
	s["genie_key"] = n.Properties["genieKey"]
	return []tfMap{s}
}

func pagerdutyFromNotifier(n *humio.Notifier) []tfMap {
	s := tfMap{}
	s["routing_key"] = n.Properties["routingKey"]
	s["severity"] = n.Properties["severity"]
	return []tfMap{s}
}

func slackFromNotifier(n *humio.Notifier) []tfMap {
	s := tfMap{}
	s["fields"] = n.Properties["fields"]
	s["url"] = n.Properties["url"]
	return []tfMap{s}
}

func slackpostmessageFromNotifier(n *humio.Notifier) []tfMap {
	s := tfMap{}
	s["api_token"] = n.Properties["apiToken"]
	s["channels"] = n.Properties["channels"]
	s["fields"] = n.Properties["fields"]
	s["use_proxy"] = n.Properties["useProxy"]
	return []tfMap{s}
}

func victoropsFromNotifier(n *humio.Notifier) []tfMap {
	s := tfMap{}
	s["message_type"] = n.Properties["messageType"]
	s["notify_url"] = n.Properties["notifyUrl"]
	return []tfMap{s}
}

func webhookFromNotifier(n *humio.Notifier) []tfMap {
	s := tfMap{}
	s["body_template"] = n.Properties["bodyTemplate"]
	s["headers"] = n.Properties["headers"]
	s["method"] = n.Properties["method"]
	s["url"] = n.Properties["url"]
	return []tfMap{s}
}

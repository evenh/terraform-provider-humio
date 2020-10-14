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
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	humio "github.com/humio/cli/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNotifierRequiredFields(t *testing.T) {
	config := notifierEmpty
	accTestCase(t, []resource.TestStep{
		{
			Config:      config,
			ExpectError: regexp.MustCompile(`The argument "repository" is required, but no definition was found.`),
		},
		{
			Config:      config,
			ExpectError: regexp.MustCompile(`The argument "entity" is required, but no definition was found.`),
		},
		{
			Config:      config,
			ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`),
		},
	}, nil)
}

func TestAccNotifierInvalidInputs(t *testing.T) {
	config := notifierInvalidInputs
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "repository"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "entity"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "name"`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "email" is not expected here`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "humiorepo" is not expected here`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "opsgenie" is not expected here`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "pagerduty" is not expected here`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "slack" is not expected here`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "slackpostmessage" is not expected here`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "victorops" is not expected here`)},
		{Config: config, ExpectError: regexp.MustCompile(`An argument named "webhook" is not expected here`)},
	}, nil)
}

func TestAccNotifierInvalidEmailSettings(t *testing.T) {
	config := notifierInvalidEmailSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "body_template"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "recipients"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "subject_template"`)},
	}, nil)
}

func TestAccNotifierInvalidHumioRepoSettings(t *testing.T) {
	config := notifierInvalidHumioRepoSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "ingest_token"`)},
	}, nil)
}

func TestAccNotifierInvalidOpsGenieSettings(t *testing.T) {
	config := notifierInvalidOpsGenieSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "api_url"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "genie_key"`)},
	}, nil)
}

func TestAccNotifierInvalidPagerDutySettings(t *testing.T) {
	config := notifierInvalidPagerDutySettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "routing_key"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "severity"`)},
	}, nil)
}

func TestAccNotifierInvalidSlackSettings(t *testing.T) {
	config := notifierInvalidSlackSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "fields"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "url"`)},
	}, nil)
}

func TestAccNotifierInvalidSlackPostMessageSettings(t *testing.T) {
	config := notifierInvalidSlackPostMessageSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "api_token"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "channels"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "fields"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "use_proxy"`)},
	}, nil)
}

func TestAccNotifierInvalidVictorOpsSettings(t *testing.T) {
	config := notifierInvalidVictorOpsSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "message_type"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "notify_url"`)},
	}, nil)
}

func TestAccNotifierInvalidWebHookSettings(t *testing.T) {
	config := notifierInvalidWebHookSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "body_template"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "headers"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "method"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "url"`)},
	}, nil)
}

func TestAccNotifierEmailBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierEmailBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "EmailNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-email-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.recipients.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.recipients.0", "test@example.org"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.body_template", ""),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.subject_template", ""),

				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierEmailBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierEmailBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "EmailNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-email-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.recipients.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.recipients.0", "test@example.org"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.body_template", ""),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.subject_template", ""),

				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
		{
			Config: notifierEmailFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "EmailNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-email-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.recipients.#", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.recipients.0", "test@example.org"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.recipients.0", "ops@example.org"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.body_template", "this is the body"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.subject_template", "this is the subject"),

				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: notifierEmailFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "EmailNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-email-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.recipients.#", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.recipients.0", "test@example.org"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.recipients.1", "ops@example.org"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.body_template", "this is the body"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.subject_template", "this is the subject"),

				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierEmailFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierEmailFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "EmailNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-email-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.recipients.#", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.recipients.0", "test@example.org"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.recipients.1", "ops@example.org"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.body_template", "this is the body"),
				resource.TestCheckResourceAttr("humio_notifier.test", "email.0.subject_template", "this is the subject"),

				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierHumioRepoFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierHumioRepoFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "HumioRepoNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-humiorepo-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.0.ingest_token", "secrettoken"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierOpsGenieBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierOpsGenieBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "OpsGenieNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-opsgenie-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.0.api_url", "https://api.opsgenie.com"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.0.genie_key", "secretgeniekey"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierOpsGenieBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierOpsGenieBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "OpsGenieNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-opsgenie-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.0.api_url", "https://api.opsgenie.com"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.0.genie_key", "secretgeniekey"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
		{
			Config: notifierOpsGenieFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "OpsGenieNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-opsgenie-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.0.api_url", "https://127.0.0.1/iasjdojaoijdioajd"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.0.genie_key", "secretgeniekey"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: notifierOpsGenieFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "OpsGenieNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-opsgenie-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.0.api_url", "https://127.0.0.1/iasjdojaoijdioajd"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.0.genie_key", "secretgeniekey"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierOpsGenieFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierOpsGenieFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "OpsGenieNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-opsgenie-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.0.api_url", "https://127.0.0.1/iasjdojaoijdioajd"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.0.genie_key", "secretgeniekey"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierPagerDutyFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierPagerDutyFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "PagerDutyNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-pagerduty-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.0.routing_key", "secretroutingkey"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.0.severity", "critical"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierSlackBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierSlackBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "SlackNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-slack-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.%", "3"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.Events String", "{events_str}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.Time Interval", "{query_time_interval}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.url", "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierSlackBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierSlackBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "SlackNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-slack-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.%", "3"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.Events String", "{events_str}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.Time Interval", "{query_time_interval}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.url", "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
		{
			Config: notifierSlackFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "SlackNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-slack-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.%", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.Link", "{url}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.url", "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: notifierSlackFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "SlackNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-slack-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.%", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.Link", "{url}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.url", "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierSlackFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierSlackFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "SlackNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-slack-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.%", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.Link", "{url}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.0.url", "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierSlackPostMessageBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierSlackPostMessageBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "SlackPostMessageNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-slackpostmessage-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.api_token", "secretapitoken"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.#", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.0", "#alerts"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.1", "#ops"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.%", "3"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.Events String", "{events_str}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.Time Interval", "{query_time_interval}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.use_proxy", "true"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierSlackPostMessageBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierSlackPostMessageBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "SlackPostMessageNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-slackpostmessage-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.api_token", "secretapitoken"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.#", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.0", "#alerts"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.1", "#ops"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.%", "3"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.Events String", "{events_str}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.Time Interval", "{query_time_interval}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.use_proxy", "true"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
		{
			Config: notifierSlackPostMessageFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "SlackPostMessageNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-slackpostmessage-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.api_token", "secretapitoken"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.#", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.0", "#alerts"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.1", "#ops"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.%", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.Link", "{url}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.use_proxy", "false"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: notifierSlackPostMessageFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "SlackPostMessageNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-slackpostmessage-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.api_token", "secretapitoken"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.#", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.0", "#alerts"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.1", "#ops"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.%", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.Link", "{url}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.use_proxy", "false"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierSlackPostMessageFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierSlackPostMessageFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "SlackPostMessageNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-slackpostmessage-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.api_token", "secretapitoken"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.#", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.0", "#alerts"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.channels.1", "#ops"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.%", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.Link", "{url}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.fields.Query", "{query_string}"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.0.use_proxy", "false"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierVictorOpsFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierVictorOpsFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "VictorOpsNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-victorops-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.0.message_type", "important"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.0.notify_url", "https://127.0.0.1/iasjdojaoijdioajd"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierWebHookBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierWebHookBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "WebHookNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-webhook-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.body_template", "{\n  \"repository\": \"{repo_name}\",\n  \"timestamp\": \"{alert_triggered_timestamp}\",\n  \"alert\": {\n    \"name\": \"{alert_name}\",\n    \"description\": \"{alert_description}\",\n    \"query\": {\n      \"queryString\": \"{query_string} \",\n      \"end\": \"{query_time_end}\",\n      \"start\": \"{query_time_start}\"\n    },\n    \"notifierID\": \"{alert_notifier_id}\",\n    \"id\": \"{alert_id}\"\n  },\n  \"warnings\": \"{warnings}\",\n  \"events\": {events},\n  \"numberOfEvents\": {event_count}\n  }"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.headers.%", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.headers.Content-Type", "application/json"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.method", "POST"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.url", "https://127.0.0.1/iasjdojaoijdioajd"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierWebHookBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierWebHookBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "WebHookNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-webhook-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.body_template", "{\n  \"repository\": \"{repo_name}\",\n  \"timestamp\": \"{alert_triggered_timestamp}\",\n  \"alert\": {\n    \"name\": \"{alert_name}\",\n    \"description\": \"{alert_description}\",\n    \"query\": {\n      \"queryString\": \"{query_string} \",\n      \"end\": \"{query_time_end}\",\n      \"start\": \"{query_time_start}\"\n    },\n    \"notifierID\": \"{alert_notifier_id}\",\n    \"id\": \"{alert_id}\"\n  },\n  \"warnings\": \"{warnings}\",\n  \"events\": {events},\n  \"numberOfEvents\": {event_count}\n  }"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.headers.%", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.headers.Content-Type", "application/json"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.method", "POST"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.url", "https://127.0.0.1/iasjdojaoijdioajd"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
			),
		},
		{
			Config: notifierWebHookFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "WebHookNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-webhook-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.body_template", "custom body"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.headers.%", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.headers.custom/header1", "this1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.headers.custom2", "this2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.method", "GET"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.url", "https://127.0.0.1/iasjdojaoijdioajd"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: notifierWebHookFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "WebHookNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-webhook-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.body_template", "custom body"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.headers.%", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.headers.custom/header1", "this1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.headers.custom2", "this2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.method", "GET"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.url", "https://127.0.0.1/iasjdojaoijdioajd"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func TestAccNotifierWebHookFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: notifierWebHookFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("humio_notifier.test", "notifier_id"),
				resource.TestCheckResourceAttr("humio_notifier.test", "repository", "sandbox"),
				resource.TestCheckResourceAttr("humio_notifier.test", "entity", "WebHookNotifier"),
				resource.TestCheckResourceAttr("humio_notifier.test", "name", "notifier-webhook-test"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.#", "1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.body_template", "custom body"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.headers.%", "2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.headers.custom/header1", "this1"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.headers.custom2", "this2"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.method", "GET"),
				resource.TestCheckResourceAttr("humio_notifier.test", "webhook.0.url", "https://127.0.0.1/iasjdojaoijdioajd"),

				resource.TestCheckResourceAttr("humio_notifier.test", "email.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "humiorepo.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "opsgenie.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "pagerduty.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slack.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "slackpostmessage.#", "0"),
				resource.TestCheckResourceAttr("humio_notifier.test", "victorops.#", "0"),
			),
		},
	}, testAccCheckNotifierDestroy)
}

func testAccCheckNotifierDestroy(s *terraform.State) error {
	conn := testAccProviders["humio"].Meta().(*humio.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "humio_notifier" {
			continue
		}

		parts := parseRepositoryAndID(rs.Primary.ID)
		resp, err := conn.Notifiers().Get(parts[0], parts[1])
		emptyNotifier := humio.Notifier{}
		if err == nil {
			if !reflect.DeepEqual(*resp, emptyNotifier) {
				return fmt.Errorf("notifier still exist for id %s: %#+v", rs.Primary.ID, *resp)
			}
		}
		if err != nil {
			if strings.HasPrefix(err.Error(), "could not find a notifier") {
				return nil
			}
			return fmt.Errorf("could not validate if notifers have been cleaned up: %s", err)
		}
	}
	return nil
}

const notifierEmpty = `
resource "humio_notifier" "test" {}
`

const notifierInvalidInputs = `
resource "humio_notifier" "test" {
    repository       = ["invalid"]
    entity           = ["invalid"]
    name             = ["invalid"]
    email            = "invalid"
    humiorepo        = "invalid"
    opsgenie         = "invalid"
    pagerduty        = "invalid"
    slack            = "invalid"
    slackpostmessage = "invalid"
    victorops        = "invalid"
    webhook          = "invalid"
}
`

const notifierInvalidEmailSettings = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "EmailNotifer"
    name       = "notifier-invalid-email"
    email {
        body_template    = ["invalid"]
        recipients       = "invalid"
        subject_template = ["invalid"]
    }
}
`

const notifierInvalidHumioRepoSettings = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "HumioRepoNotifer"
    name       = "notifier-invalid-humiorepo"
    humiorepo {
        ingest_token = ["invalid"]
    }
}
`

const notifierInvalidOpsGenieSettings = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "OpsGenieNotifer"
    name       = "notifier-invalid-opsgenie"
    opsgenie {
        api_url   = ["invalid"]
        genie_key = ["invalid"]
    }
}
`

const notifierInvalidPagerDutySettings = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "PagerDutyNotifer"
    name       = "notifier-invalid-pagerduty"
    pagerduty {
        routing_key = ["invalid"]
        severity    = ["invalid"]
    }
}
`

const notifierInvalidSlackSettings = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "SlackNotifer"
    name       = "notifier-invalid-slack"
    slack {
        fields = "invalid"
        url    = ["invalid"]
    }
}
`

const notifierInvalidSlackPostMessageSettings = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "SlackPostMessageNotifer"
    name       = "notifier-invalid-slackpostmessage"
    slackpostmessage {
        api_token = ["invalid"]
        channels  = "invalid"
        fields    = "invalid"
        use_proxy = ["invalid"]
    }
}
`

const notifierInvalidVictorOpsSettings = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "VictorOpsNotifer"
    name       = "notifier-invalid-victorops"
    victorops {
        message_type = ["invalid"]
        notify_url   = ["invalid"]
    }
}
`

const notifierInvalidWebHookSettings = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "WebHookNotifer"
    name       = "notifier-invalid-webhook"
    webhook {
        body_template = ["invalid"]
        headers       = "invalid"
        method        = ["invalid"]
        url           = ["invalid"]
    }
}
`

const notifierEmailBasic = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "EmailNotifier"
    name       = "notifier-email-test"
    email {
        recipients = ["test@example.org"]
    }
}
`

const notifierEmailFull = `
resource "humio_notifier" "test" {
    repository  = "sandbox"
    entity      = "EmailNotifier"
    name        = "notifier-email-test"
    email {
        body_template    = "this is the body"
        recipients       = ["test@example.org", "ops@example.org"]
        subject_template = "this is the subject"
    }
}
`

const notifierHumioRepoFull = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "HumioRepoNotifier"
    name       = "notifier-humiorepo-test"
    humiorepo {
        ingest_token = "secrettoken"
    }
}
`

const notifierOpsGenieBasic = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "OpsGenieNotifier"
    name       = "notifier-opsgenie-test"
    opsgenie {
        genie_key = "secretgeniekey"
    }
}
`

const notifierOpsGenieFull = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "OpsGenieNotifier"
    name       = "notifier-opsgenie-test"
    opsgenie {
        api_url   = "https://127.0.0.1/iasjdojaoijdioajd"
        genie_key = "secretgeniekey"
    }
}
`

const notifierPagerDutyFull = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "PagerDutyNotifier"
    name       = "notifier-pagerduty-test"
    pagerduty {
        routing_key = "secretroutingkey"
        severity    = "critical"
    }
}
`

const notifierSlackBasic = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "SlackNotifier"
    name       = "notifier-slack-test"
    slack {
        fields = {
            "Events String" = "{events_str}"
            "Query"         = "{query_string}"
            "Time Interval" = "{query_time_interval}"
        }
        url = "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"
    }
}
`

const notifierSlackFull = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "SlackNotifier"
    name       = "notifier-slack-test"
    slack {
        fields = {
			"Link" = "{url}"
			"Query" = "{query_string}"
        }
        url = "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"
    }
}
`

const notifierSlackPostMessageBasic = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "SlackPostMessageNotifier"
    name       = "notifier-slackpostmessage-test"
    slackpostmessage {
        api_token = "secretapitoken"
        channels  = ["#alerts","#ops"]
        fields = {
            "Events String" = "{events_str}"
            "Query"         = "{query_string}"
            "Time Interval" = "{query_time_interval}"
        }
    }
}
`

const notifierSlackPostMessageFull = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "SlackPostMessageNotifier"
    name       = "notifier-slackpostmessage-test"
    slackpostmessage {
        api_token = "secretapitoken"
        channels  = ["#alerts","#ops"]
        fields = {
			"Link" = "{url}"
			"Query" = "{query_string}"
        }
        use_proxy = false
    }
}
`

const notifierVictorOpsFull = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "VictorOpsNotifier"
    name       = "notifier-victorops-test"
    victorops {
        message_type = "important"
        notify_url   = "https://127.0.0.1/iasjdojaoijdioajd"
    }
}
`

const notifierWebHookBasic = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "WebHookNotifier"
    name       = "notifier-webhook-test"
    webhook {
        headers = {
            "Content-Type" = "application/json"
        }
        url = "https://127.0.0.1/iasjdojaoijdioajd"
    }
}
`

const notifierWebHookFull = `
resource "humio_notifier" "test" {
    repository = "sandbox"
    entity     = "WebHookNotifier"
    name       = "notifier-webhook-test"
    webhook {
        body_template = "custom body"
        headers       = {
            "custom/header1" = "this1"
            custom2          = "this2"
        }
        method = "GET"
        url    = "https://127.0.0.1/iasjdojaoijdioajd"
    }
}
`

var wantEmailNotifier = humio.Notifier{
	ID:     "",
	Entity: "EmailNotifier",
	Name:   "test-notifier",
	Properties: map[string]interface{}{
		"recipients":      []interface{}{"test@example.org", "ops@example.org"},
		"bodyTemplate":    "this is the subject",
		"subjectTemplate": "this is the body",
	},
}

var wantHumioRepoNotifier = humio.Notifier{
	ID:     "",
	Entity: "HumioRepoNotifier",
	Name:   "test-notifier",
	Properties: map[string]interface{}{
		"ingestToken": "12345678901234567890123456789012",
	},
}

var wantOpsGenieNotifier = humio.Notifier{
	ID:     "",
	Entity: "OpsGenieNotifier",
	Name:   "test-notifier",
	Properties: map[string]interface{}{
		"apiUrl":   "https://example.org",
		"genieKey": "12345678901234567890123456789012",
	},
}

var wantPagerDutyNotifier = humio.Notifier{
	ID:     "",
	Entity: "PagerDutyNotifier",
	Name:   "test-notifier",
	Properties: map[string]interface{}{
		"routingKey": "12345678901234567890123456789012",
		"severity":   "critical",
	},
}

var wantSlackNotifier = humio.Notifier{
	ID:     "",
	Entity: "SlackNotifier",
	Name:   "test-notifier",
	Properties: map[string]interface{}{
		"url": "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ",
		"fields": map[string]interface{}{
			"Link":  "{url}",
			"Query": "{query_string}",
		},
	},
}

var wantSlackPostMessageNotifier = humio.Notifier{
	ID:     "",
	Entity: "SlackPostMessageNotifier",
	Name:   "test-notifier",
	Properties: map[string]interface{}{
		"apiToken": "12345678901234567890123456789012",
		"channels": []interface{}{"#alerts", "ops"},
		"fields": map[string]interface{}{
			"Link":  "{url}",
			"Query": "{query_string}",
		},
		"useProxy": true,
	},
}

var wantVictorOpsNotifier = humio.Notifier{
	ID:     "",
	Entity: "VictorOpsNotifier",
	Name:   "test-notifier",
	Properties: map[string]interface{}{
		"messageType": "12345678901234567890123456789012",
		"notifyUrl":   "https://example.org",
	},
}

var wantWebHookNotifier = humio.Notifier{
	ID:     "",
	Entity: "WebHookNotifier",
	Name:   "test-notifier",
	Properties: map[string]interface{}{
		"bodyTemplate": "12345678901234567890123456789012",
		"headers": map[string]interface{}{
			"token": "abcdefghij123456678",
		},
		"method": "POST",
		"url":    "https://example.org",
	},
}

func TestEncodeDecodeEmailNotifierResource(t *testing.T) {
	res := resourceNotifier()
	data := res.TestResourceData()
	resourceDataFromNotifier(&wantEmailNotifier, data)
	got, err := notifierFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantEmailNotifier, got) {
		t.Error(cmp.Diff(wantEmailNotifier, got))
	}
}

func TestEncodeDecodeHumioRepoNotifierResource(t *testing.T) {
	res := resourceNotifier()
	data := res.TestResourceData()
	resourceDataFromNotifier(&wantHumioRepoNotifier, data)
	got, err := notifierFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantHumioRepoNotifier, got) {
		t.Error(cmp.Diff(wantHumioRepoNotifier, got))
	}
}

func TestEncodeDecodeOpsGenieNotifierResource(t *testing.T) {
	res := resourceNotifier()
	data := res.TestResourceData()
	resourceDataFromNotifier(&wantOpsGenieNotifier, data)
	got, err := notifierFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantOpsGenieNotifier, got) {
		t.Error(cmp.Diff(wantOpsGenieNotifier, got))
	}
}

func TestEncodeDecodePagerDutyNotifierResource(t *testing.T) {
	res := resourceNotifier()
	data := res.TestResourceData()
	resourceDataFromNotifier(&wantPagerDutyNotifier, data)
	got, err := notifierFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantPagerDutyNotifier, got) {
		t.Error(cmp.Diff(wantPagerDutyNotifier, got))
	}
}

func TestEncodeDecodeSlackNotifierResource(t *testing.T) {
	res := resourceNotifier()
	data := res.TestResourceData()
	resourceDataFromNotifier(&wantSlackNotifier, data)
	got, err := notifierFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantSlackNotifier, got) {
		t.Error(cmp.Diff(wantSlackNotifier, got))
	}
}

func TestEncodeDecodeSlackPostMessageNotifierResource(t *testing.T) {
	res := resourceNotifier()
	data := res.TestResourceData()
	resourceDataFromNotifier(&wantSlackPostMessageNotifier, data)
	got, err := notifierFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantSlackPostMessageNotifier, got) {
		t.Error(cmp.Diff(wantSlackPostMessageNotifier, got))
	}
}

func TestEncodeDecodeVictorOpsNotifierResource(t *testing.T) {
	res := resourceNotifier()
	data := res.TestResourceData()
	resourceDataFromNotifier(&wantVictorOpsNotifier, data)
	got, err := notifierFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantVictorOpsNotifier, got) {
		t.Error(cmp.Diff(wantVictorOpsNotifier, got))
	}
}

func TestEncodeDecodeWebHookNotifierResource(t *testing.T) {
	res := resourceNotifier()
	data := res.TestResourceData()
	resourceDataFromNotifier(&wantWebHookNotifier, data)
	got, err := notifierFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantWebHookNotifier, got) {
		t.Error(cmp.Diff(wantWebHookNotifier, got))
	}
}

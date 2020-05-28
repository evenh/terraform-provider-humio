resource "humio_notifier" "example_webhook" {
  repository = "humio"
  name       = "example_webhook"
  entity     = "WebHookNotifier"

  webhook {
    method  = "POST"
    url     = "http://webhook.site/#!/1bb57835-4df2-42a4-84ae-fa5379ee4deb"
    headers = {
      Content-Type = "application/json"
    }
    body_template = <<TEMPLATE
{
  "repository": "{repo_name}",
  "timestamp": "{alert_triggered_timestamp}",
  "alert": {
    "name": "{alert_name}",
    "description": "{alert_description}",
    "query": {
      "queryString": "{query_string} ",
      "end": "{query_time_end}",
      "start": "{query_time_start}"
    },
    "notifierID": "{alert_notifier_id}",
    "id": "{alert_id}"
  },
  "warnings": "{warnings}",
  "events": {events},
  "numberOfEvents": {event_count}
}
TEMPLATE
  }
}

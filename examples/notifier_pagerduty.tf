resource "humio_notifier" "example_pagerduty" {
  repository = "humio"
  name       = "example_pagerduty"
  entity     = "PagerDutyNotifier"

  pagerduty {
    routing_key = "XXXXXXXXXXXXXXX"
    severity    = "critical"
  }
}
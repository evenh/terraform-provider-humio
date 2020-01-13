resource "humio_notifier" "example_victorops" {
  repository = "humio"
  name       = "example_victorops"
  entity     = "VictorOpsNotifier"

  victorops {
    message_type = "critical"
    notify_url   = "https://example.org"
  }
}
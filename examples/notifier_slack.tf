resource "humio_notifier" "example_slack" {
  repository = "humio"
  name       = "example_slack"
  entity     = "SlackNotifier"

  slack {
    url    = "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"
    fields = {
      "Events String" = "{events_str}"
      "Query"         = "{query_string}"
      "Time Interval" = "{query_time_interval}"
    }
  }
}
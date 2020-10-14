resource "humio_notifier" "example_slackpostmessage" {
  repository = "humio"
  name       = "example_slackpostmessage"
  entity     = "SlackPostMessageNotifier"

  slackpostmessage {
    api_key  = "abcdefghij1234567890"
    channels = ["#alerts","#ops"]
    fields = {
      "Events String" = "{events_str}"
      "Query"         = "{query_string}"
      "Time Interval" = "{query_time_interval}"
    }
  }
}
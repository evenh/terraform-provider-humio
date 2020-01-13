resource "humio_notifier" "example_slack" {
  repository = "humio"
  name       = "example_slack"
  entity     = "SlackNotifier"

  slack {
    url    = "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ"
    fields = {
      Link = "{url}"
      Query = "{query_string}"
    }
  }
}
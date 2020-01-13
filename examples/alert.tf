resource "humio_alert" "example_alert_with_labels" {
  repository  = humio_notifier.example_email.repository
  name        = "example_alert_with_labels"

  notifiers   = [humio_notifier.example_email.id]

  labels               = ["terraform", "ops"]
  throttle_time_millis = 300000
  silenced             = true
  query                = "count()"
  start                = "24h"
}

resource "humio_alert" "example_alert_without_labels" {
  repository  = humio_notifier.example_email_body.repository
  name        = "example_alert_without_labels"

  notifiers = [humio_notifier.example_email_body.id]

  throttle_time_millis = 300000
  silenced             = true
  query                = "count()"
  start                = "24h"
}

resource "humio_alert" "example_alert_with_description" {
  repository  = humio_notifier.example_email_body.repository
  name        = "example_alert_with_description"
  description = "lorem ipsum...."

  notifiers = [
    humio_notifier.example_email_body.id,
    humio_notifier.example_email_subject.id,
  ]

  link_url             = "http://localhost:8080/humio/search?query=count()&live=true&start=24h&fullscreen=false"
  labels               = ["terraform", "ops"]
  throttle_time_millis = 300000
  silenced             = true
  query                = "count()"
  start                = "24h"
}
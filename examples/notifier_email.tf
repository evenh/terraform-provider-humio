resource "humio_notifier" "example_email" {
  repository = "sandbox"
  name       = "example_email"
  entity     = "EmailNotifier"

  email {
    recipients = ["ops@example.com"]
  }
}

resource "humio_notifier" "example_email_body" {
  repository = "humio"
  name       = "example_email_body"
  entity     = "EmailNotifier"

  email {
    recipients    = ["ops@example.com"]
    body_template = "{event_count}"
  }
}

resource "humio_notifier" "example_email_subject" {
  repository = "humio"
  name       = "example_email_subject"
  entity     = "EmailNotifier"

  email {
    recipients       = ["ops@example.com"]
    subject_template = "{alert_name}"
  }
}

resource "humio_notifier" "example_email_body_subject" {
  repository = "humio"
  name       = "example_email_body_subject"
  entity     = "EmailNotifier"

  email {
    recipients       = ["ops@example.com"]
    body_template    = "{event_count}"
    subject_template = "{alert_name}"
  }
}

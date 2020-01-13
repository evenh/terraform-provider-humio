resource "humio_ingest_token" "example_ingest_token_without_parser" {
  repository = "humio"
  name       = "example_ingest_token_without_parser"
}

resource "humio_ingest_token" "example_ingest_token_with_accesslog_parser" {
  repository = "humio"
  name       = "example_ingest_token_with_accesslog_parser"
  parser     = "accesslog"
}

output "ingest_token_without_parser" {
  value       = humio_ingest_token.example_ingest_token_without_parser.token
}

output "ingest_token_with_accesslog_parser" {
  value       = humio_ingest_token.example_ingest_token_with_accesslog_parser.token
}
resource "humio_repository" "example_repo_minimal_fields_set" {
  name = "example_repo_minimal_fields_set"

  retention {}
}

resource "humio_repository" "example_repo_all_fields_set" {
  name        = "example_repo_all_fields_set"
  description = "This is an example"

  retention {
    storage_size_in_gb = 5
    ingest_size_in_gb  = 10
    time_in_days       = 30
  }
}
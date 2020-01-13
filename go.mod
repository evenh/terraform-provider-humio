module github.com/humio/terraform-provider-humio

go 1.13

require (
	github.com/hashicorp/go-getter v1.4.2-0.20200106182914-9813cbd4eb02 // indirect
	github.com/hashicorp/go-hclog v0.10.0 // indirect
	github.com/hashicorp/hcl/v2 v2.2.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.4.1
	github.com/humio/cli v0.24.2
	github.com/stretchr/testify v1.4.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

// Until the PR for the API package is merged in to the humio/cli project,
// checkout the branch 'mike/add_alerts_notifiers_etc' in from github.com/humio/cli
// and point to the path here. The pending PR is: https://github.com/humio/cli/pull/12
replace github.com/humio/cli => /Users/mike/git/humio-cli

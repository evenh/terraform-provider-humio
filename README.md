# Terraform provider for Humio

**EXPERIMENTAL: It is still early days for this plugin and thus not ready for prime time yet.**

## Currently tested with

- [Terraform](https://www.terraform.io/downloads.html) v0.12+
- [Go](https://golang.org/doc/install) 1.14 (to build the provider plugin)

## Installing the provider

We do not publish binaries yet, so for now you need to build it yourself.

1. Clone the git repository of the provider

```bash
git clone https://github.com/humio/terraform-provider-humio
cd terraform-provider-humio
```

2. Build the provider plugin

```bash
go build -o terraform-provider-humio
```

3. Install the provider by following the [official documentation](https://www.terraform.io/docs/plugins/basics.html#installing-plugins).

```bash
mkdir -p ~/.terraform.d/plugins
cp terraform-provider-humio ~/.terraform.d/plugins
```

## Using the provider

### Authentication

You can specify the address and API token for Humio directly on the provider like this:

```hcl
provider "humio" {
    addr      = "https://humio.example.com/" # If not specified, the default is: https://cloud.humio.com/
    api_token = "XXXXXXXXXXXXXXXXXXXXXXXXX"
}
```

It is also possible to specify these settings using environment variables `HUMIO_ADDR` and `HUMIO_API_TOKEN`.

In most cases we recommend configuring the Humio address directly on the provider as described above, whereas the API token should be set as an environment variable to keep it out of the code.

### Supported resources and examples

See [examples directory](examples).

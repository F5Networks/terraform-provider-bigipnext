**ARCHIVED:** This repository is archived due to: **Modernizing BIG-IP TMOS and discontinuing BIG-IP Next**.  
For more information, see [K000152956] (https://my.f5.com/manage/s/article/K000152956)

# Terraform Provider BIG-IP Next

* BIG-IP Next Terraform provider helps you managing BIG-IP Next devices through BIG-IP Next Central Manager (CM) API.

* BIG-IP Next uses a combination of BIG-IP Next Central Manager and BIG-IP Next instances to implement application delivery and security. The BIG-IP Next Central Manager manages the BIG-IP Next instances, assuming responsibility for all administrative and management tasks. The BIG-IP Next instances, responsible for data processing, provide robust automation capabilities, scalability, and ease-of-use for organizations running applications on-premise, in the cloud, or out at the edge

  For more information:

  [F5 BIG-IP Next](https://clouddocs.f5.com/bigip-next/latest/)
  
## Requirements

* [Terraform](https://www.terraform.io/downloads) > 1.x
* [Go](https://go.dev/doc/install) >= 1.19
* [GNU Make](https://www.gnu.org/software/make/)
* [golangci-lint](https://golangci-lint.run/usage/install/#local-installation) (optional)

## Using the Provider

This Terraform Provider is available to install automatically via `terraform init`. It is recommended to setup the following Terraform configuration to pin the major version:

```hcl
# Terraform 1.2.x and later
terraform {
  required_providers {
    bigipnext = {
      source  = "f5networks/bigipnext"
      version = "~> X.Y" # where X.Y is the current major version and minor version
    }
  }
}
```

## Documentation, questions and discussions
Official documentation on how to use this provider can be found on the
[Terraform Registry](https://registry.terraform.io/providers/F5Networks/bigipnext/latest/docs).
In case of specific questions or discussions, please use the
HashiCorp [Terraform Providers Discuss forums](https://discuss.hashicorp.com/c/terraform-providers/31),
in accordance with HashiCorp [Community Guidelines](https://www.hashicorp.com/community-guidelines).

We also provide:

* [Support](.github/SUPPORT.md) page for help when using the provider
* [Contributing](.github/CONTRIBUTING.md) guidelines in case you want to help this project

## Compatibility

Compatibility table between this provider, the [Terraform Plugin Protocol](https://www.terraform.io/plugin/how-terraform-works#terraform-plugin-protocol)
version it implements, and Terraform:

| BIG-IP Next Provider |     Terraform Plugin Protocol      | Terraform | BIG-IP Next CM     Version |
|:--------------------:|:----------------------------------:|:---------:|:--------------------------:|
|  `>= 1.0.0`          |                `6`                 | `>= 1.x`  |      `>= 20.1.0`           |

Details can be found querying the [Registry API](https://www.terraform.io/internals/provider-registry-protocol#list-available-versions)
that return all the details about which version are currently available for a particular provider.

## Development

### Building

1. `git clone` this repository and `cd` into its directory
2. `go build` will trigger the Golang build

The provided `GNUmakefile` defines additional commands generally useful during development,
like for running tests, generating documentation, code formatting and linting.
Taking a look at it's content is recommended.

### Testing

In order to test the provider, you can run

* `make test` to run provider unit tests
* `make testacc` to run provider acceptance tests

It's important to note that acceptance tests (`testacc`) will actually spawn real resources, and often cost money to run. Read more about they work on the
[official page](https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests).

### Generating documentation

This provider uses [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs/)
to generate documentation and store it in the `docs/` directory.
Once a release is cut, the Terraform Registry will download the documentation from `docs/`
and associate it with the release version. Read more about how this works on the
[official page](https://www.terraform.io/registry/providers/docs).

Use `make generate` to ensure the documentation is regenerated with any changes.

### Using a development build

If [running tests and acceptance tests](#testing) isn't enough, it's possible to set up a local terraform configuration
to use a development builds of the provider. This can be achieved by leveraging the Terraform CLI
[configuration file development overrides](https://www.terraform.io/cli/config/config-file#development-overrides-for-provider-developers).

First, use `make install` to place a fresh development build of the provider in your
[`${GOBIN}`](https://pkg.go.dev/cmd/go#hdr-Compile_and_install_packages_and_dependencies)
(defaults to `${GOPATH}/bin` or `${HOME}/go/bin` if `${GOPATH}` is not set). Repeat
this every time you make changes to the provider locally.

Then, setup your environment following [these instructions](https://www.terraform.io/plugin/debugging#terraform-cli-development-overrides)
to make your local terraform use your local build.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

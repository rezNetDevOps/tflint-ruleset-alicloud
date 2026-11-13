# TFLint Ruleset for terraform-provider-alicloud
[![Build Status](https://github.com/rezNetDevOps/tflint-ruleset-alicloud/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/rezNetDevOps/tflint-ruleset-alicloud/actions)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](LICENSE)

TFLint ruleset plugin for the Terraform [Alibaba Cloud provider](https://registry.terraform.io/providers/aliyun/alicloud/latest) (`aliyun/alicloud`).

This ruleset focuses on possible errors and best practices for Alibaba Cloud resources. The rules warn about configuration that is valid Terraform but is likely to fail at `terraform apply`, or that follows a practice worth avoiding.

The project is in early development. No release has been published and no rules are available yet.

## Requirements

- TFLint v0.46+
- Go, at the version in `go.mod`, to build from source

## Installation

Once a release is published, you can install the plugin by adding a config to `.tflint.hcl` and running `tflint --init`:

```hcl
plugin "alicloud" {
    enabled = true
    version = "0.1.0"
    source  = "github.com/rezNetDevOps/tflint-ruleset-alicloud"
}
```

Until then, build and install the plugin from source as described below.

## Getting Started

Terraform and the provider check that a configuration is well formed, but many values are checked only by the Alibaba Cloud API at apply time. For example:

```hcl
resource "alicloud_instance" "web" {
  instance_type = "ecs.g7-large" # should be "ecs.g7.large"
  # ...
}
```

`terraform validate` and `terraform plan` both accept this. The provider only checks that the value starts with `ecs.`. The error appears only when the ECS API rejects the value during `terraform apply`.

The goal of this ruleset is to report such errors before they reach `apply`.

## Rules

No rules are released yet. See [docs/rules/README.md](docs/rules/README.md) for the rules on `main`, with the severity and default state of each.

## Building the plugin

Clone the repository locally and run the following command:

```
$ make
```

You can install the built plugin with the following:

```
$ make install
```

If you install the plugin with `make install`, omit the `version` and `source` attributes in `.tflint.hcl`:

```hcl
plugin "alicloud" {
    enabled = true
}
```

Run the tests with:

```
$ make test
```

## Add a new rule

Each rule lives in its own file under `rules/`, with a table-driven test beside it that uses `helper.TestRunner` from [tflint-plugin-sdk](https://github.com/terraform-linters/tflint-plugin-sdk). Rule names follow `alicloud_<resource>_<check>`, where `<resource>` is the resource type without the `alicloud_` prefix, for example `alicloud_instance_invalid_instance_type_format`.

See [Writing Plugins](https://github.com/terraform-linters/tflint/blob/master/docs/developer-guide/plugins.md) in the TFLint documentation for how rules access configuration.

## Contributing

Bug reports, false positive reports and rule requests are welcome as [issues](https://github.com/rezNetDevOps/tflint-ruleset-alicloud/issues). For anything larger than a small fix, open an issue before sending a pull request.

## License

[Mozilla Public License 2.0](LICENSE)

This is a community project and is not affiliated with Alibaba Cloud.

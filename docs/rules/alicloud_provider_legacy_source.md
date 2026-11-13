# alicloud_provider_legacy_source

Disallow `hashicorp/alicloud` as the provider source.

This rule is disabled by default. Enable it in `.tflint.hcl`:

```hcl
rule "alicloud_provider_legacy_source" {
  enabled = true
}
```

## Example

```hcl
terraform {
  required_version = ">= 1.5"

  required_providers {
    alicloud = {
      source  = "hashicorp/alicloud"
      version = "~> 1.200"
    }
  }
}
```

```
$ tflint
1 issue(s) found:

Notice: "hashicorp/alicloud" is not the canonical provider source, use "aliyun/alicloud" (alicloud_provider_legacy_source)

  on main.tf line 6:
   6:       source  = "hashicorp/alicloud"

Reference: https://github.com/rezNetDevOps/tflint-ruleset-alicloud/blob/v0.1.0/docs/rules/alicloud_provider_legacy_source.md
```

## Why

The Alibaba Cloud provider is maintained at `aliyun/alicloud`, and that is the
source its documentation uses. `hashicorp/alicloud` is a mirror that still
receives the same releases, so it installs and works on its own.

Terraform treats the two addresses as different providers. When one module
uses `aliyun/alicloud` and another uses `hashicorp/alicloud`, Terraform
installs both, and the module with the other address does not inherit the
`provider "alicloud"` configuration. It runs with a default, unconfigured
provider, which can pick up the wrong region or credentials.

## How To Fix

Set `source` to `aliyun/alicloud`:

```hcl
terraform {
  required_providers {
    alicloud = {
      source  = "aliyun/alicloud"
      version = "~> 1.200"
    }
  }
}
```

Run `terraform init -upgrade` afterwards. Existing state that records
`registry.terraform.io/hashicorp/alicloud` can be moved with
`terraform state replace-provider hashicorp/alicloud aliyun/alicloud`.

See the [provider documentation](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs)
and [provider requirements](https://developer.hashicorp.com/terraform/language/providers/requirements).

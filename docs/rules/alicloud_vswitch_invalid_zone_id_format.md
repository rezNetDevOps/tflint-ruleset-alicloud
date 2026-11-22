# alicloud_vswitch_invalid_zone_id_format

Disallow vswitch `zone_id` values that are not in the provider region.

## Example

```hcl
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vswitch" "web" {
  vpc_id     = alicloud_vpc.main.id
  cidr_block = "10.0.1.0/24"
  zone_id    = "cn-shanghai-b"
}
```

```
$ tflint
1 issue(s) found:

Error: "cn-shanghai-b" is not a zone in provider region "cn-hangzhou" (alicloud_vswitch_invalid_zone_id_format)

  on main.tf line 8:
   8:   zone_id    = "cn-shanghai-b"

Reference: https://github.com/rezNetDevOps/tflint-ruleset-alicloud/blob/v0.1.0/docs/rules/alicloud_vswitch_invalid_zone_id_format.md
```

## Why

A vswitch is created in the region of its provider, so its zone must be in that
region. The provider does not validate `zone_id`, so `terraform validate` and
`terraform plan` accept a zone from another region, or the region ID itself,
and `terraform apply` fails when the VPC API rejects it.

A zone ID starts with its region ID. Zones in mainland China add a hyphen and a
letter (`cn-hangzhou-h`); zones outside mainland China add only the letter
(`ap-southeast-1a`). The rule reads `region` from the `provider "alicloud"`
block the vswitch uses, including aliased providers. It skips the check when
that block has no `region` or the value is not known, such as a region set
only through environment variables.

## How To Fix

Use a zone in the provider region, or point the vswitch at a provider for the
zone's region:

```hcl
resource "alicloud_vswitch" "web" {
  vpc_id     = alicloud_vpc.main.id
  cidr_block = "10.0.1.0/24"
  zone_id    = "cn-hangzhou-h"
}
```

See [Regions and zones](https://www.alibabacloud.com/help/en/ecs/user-guide/regions-and-zones).

# alicloud_instance_invalid_instance_type_format

Disallow `instance_type` values that do not have the form `ecs.<family>.<size>`.

## Example

```hcl
resource "alicloud_instance" "web" {
  image_id      = "ubuntu_22_04_x64_20G_alibase_20240130.vhd"
  instance_type = "ecs.g7-2xlarge"
}
```

```
$ tflint
1 issue(s) found:

Error: "ecs.g7-2xlarge" is an invalid instance type, expected ecs.<family>.<size> (alicloud_instance_invalid_instance_type_format)

  on main.tf line 3:
   3:   instance_type = "ecs.g7-2xlarge"

Reference: https://github.com/rezNetDevOps/tflint-ruleset-alicloud/blob/v0.1.0/docs/rules/alicloud_instance_invalid_instance_type_format.md
```

## Why

The provider only checks that `instance_type` starts with `ecs.`, so
`terraform validate` and `terraform plan` accept a malformed value. The ECS API
rejects it when `terraform apply` creates the instance.

Every ECS instance type is lowercase and has three parts separated by dots: the
`ecs` prefix, the instance family such as `g7` or `t5-lc1m2`, and the size such
as `large` or `2xlarge`. This rule checks that shape only. It does not check
that the type exists or is sold in your zone.

## How To Fix

Write the instance type as `ecs.<family>.<size>`:

```hcl
resource "alicloud_instance" "web" {
  image_id      = "ubuntu_22_04_x64_20G_alibase_20240130.vhd"
  instance_type = "ecs.g7.2xlarge"
}
```

See [Instance type naming and classification](https://www.alibabacloud.com/help/en/ecs/user-guide/instance-specification-naming-and-classification).

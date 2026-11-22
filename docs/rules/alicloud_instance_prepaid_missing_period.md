# alicloud_instance_prepaid_missing_period

Disallow `instance_charge_type = "PrePaid"` without `period`.

## Example

```hcl
resource "alicloud_instance" "web" {
  image_id             = "ubuntu_22_04_x64_20G_alibase_20240130.vhd"
  instance_type        = "ecs.g7.large"
  instance_charge_type = "PrePaid"
}
```

```
$ tflint
1 issue(s) found:

Error: "period" is required when instance_charge_type is "PrePaid" (alicloud_instance_prepaid_missing_period)

  on main.tf line 4:
   4:   instance_charge_type = "PrePaid"

Reference: https://github.com/rezNetDevOps/tflint-ruleset-alicloud/blob/v0.1.0/docs/rules/alicloud_instance_prepaid_missing_period.md
```

## Why

A subscription instance needs a subscription length. The ECS API requires
`Period` when `InstanceChargeType` is `PrePaid`. The provider has no default for
`period` and sends it only when it is set, so `terraform validate` and
`terraform plan` accept the configuration and `terraform apply` fails when the
API rejects the request.

## How To Fix

Set `period`. The unit is set by `period_unit` and defaults to `Month`:

```hcl
resource "alicloud_instance" "web" {
  image_id             = "ubuntu_22_04_x64_20G_alibase_20240130.vhd"
  instance_type        = "ecs.g7.large"
  instance_charge_type = "PrePaid"
  period               = 1
}
```

See [Subscription](https://www.alibabacloud.com/help/en/ecs/subscription) and
the `Period` parameter of
[RunInstances](https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ecs-2014-05-26-runinstances).

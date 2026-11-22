package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AlicloudInstancePrepaidMissingPeriod(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "prepaid without period",
			Content: `
resource "alicloud_instance" "web" {
  instance_type        = "ecs.g7.large"
  instance_charge_type = "PrePaid"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudInstancePrepaidMissingPeriodRule(),
					Message: `"period" is required when instance_charge_type is "PrePaid"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 26},
						End:      hcl.Pos{Line: 4, Column: 35},
					},
				},
			},
		},
		{
			Name: "prepaid with null period",
			Content: `
resource "alicloud_instance" "web" {
  instance_charge_type = "PrePaid"
  period               = null
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudInstancePrepaidMissingPeriodRule(),
					Message: `"period" is required when instance_charge_type is "PrePaid"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 26},
						End:      hcl.Pos{Line: 3, Column: 35},
					},
				},
			},
		},
		{
			Name: "prepaid from variable default",
			Content: `
variable "charge_type" {
  default = "PrePaid"
}

resource "alicloud_instance" "web" {
  instance_charge_type = var.charge_type
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudInstancePrepaidMissingPeriodRule(),
					Message: `"period" is required when instance_charge_type is "PrePaid"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 26},
						End:      hcl.Pos{Line: 7, Column: 41},
					},
				},
			},
		},
		{
			Name: "prepaid with period",
			Content: `
resource "alicloud_instance" "web" {
  instance_charge_type = "PrePaid"
  period               = 1
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "prepaid with unknown period",
			Content: `
variable "period" {
  type = number
}

resource "alicloud_instance" "web" {
  instance_charge_type = "PrePaid"
  period               = var.period
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "postpaid without period",
			Content: `
resource "alicloud_instance" "web" {
  instance_charge_type = "PostPaid"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no charge type",
			Content: `
resource "alicloud_instance" "web" {
  instance_type = "ecs.g7.large"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAlicloudInstancePrepaidMissingPeriodRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}

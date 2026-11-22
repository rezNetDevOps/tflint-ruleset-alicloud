package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AlicloudInstanceInvalidInstanceTypeFormat(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "family without size",
			Content: `
resource "alicloud_instance" "web" {
  instance_type = "ecs.g6"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudInstanceInvalidInstanceTypeFormatRule(),
					Message: `"ecs.g6" is an invalid instance type, expected ecs.<family>.<size>`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 19},
						End:      hcl.Pos{Line: 3, Column: 27},
					},
				},
			},
		},
		{
			Name: "uppercase",
			Content: `
resource "alicloud_instance" "web" {
  instance_type = "ecs.G6.Large"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudInstanceInvalidInstanceTypeFormatRule(),
					Message: `"ecs.G6.Large" is an invalid instance type, expected ecs.<family>.<size>`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 19},
						End:      hcl.Pos{Line: 3, Column: 33},
					},
				},
			},
		},
		{
			Name: "underscore instead of dot",
			Content: `
resource "alicloud_instance" "web" {
  instance_type = "ecs.g6_large"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudInstanceInvalidInstanceTypeFormatRule(),
					Message: `"ecs.g6_large" is an invalid instance type, expected ecs.<family>.<size>`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 19},
						End:      hcl.Pos{Line: 3, Column: 33},
					},
				},
			},
		},
		{
			Name: "extra segment",
			Content: `
resource "alicloud_instance" "web" {
  instance_type = "ecs.g6.large.2"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudInstanceInvalidInstanceTypeFormatRule(),
					Message: `"ecs.g6.large.2" is an invalid instance type, expected ecs.<family>.<size>`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 19},
						End:      hcl.Pos{Line: 3, Column: 35},
					},
				},
			},
		},
		{
			Name: "valid types",
			Content: `
resource "alicloud_instance" "a" {
  instance_type = "ecs.g6.large"
}

resource "alicloud_instance" "b" {
  instance_type = "ecs.g7.16xlarge"
}

resource "alicloud_instance" "c" {
  instance_type = "ecs.t5-lc2m1.nano"
}

resource "alicloud_instance" "d" {
  instance_type = "ecs.gn7i-c8g1.2xlarge"
}

resource "alicloud_instance" "e" {
  instance_type = "ecs.e-c1m1.large"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "variable with default",
			Content: `
variable "instance_type" {
  type    = string
  default = "ecs.g6"
}

resource "alicloud_instance" "web" {
  instance_type = var.instance_type
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudInstanceInvalidInstanceTypeFormatRule(),
					Message: `"ecs.g6" is an invalid instance type, expected ecs.<family>.<size>`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 8, Column: 19},
						End:      hcl.Pos{Line: 8, Column: 36},
					},
				},
			},
		},
		{
			Name: "other resource",
			Content: `
resource "alicloud_ecs_launch_template" "web" {
  instance_type = "ecs.g6"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAlicloudInstanceInvalidInstanceTypeFormatRule()

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

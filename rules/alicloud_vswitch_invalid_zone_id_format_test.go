package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AlicloudVswitchInvalidZoneIDFormat(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "zone in another region",
			Content: `
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vswitch" "web" {
  zone_id = "cn-shanghai-b"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudVswitchInvalidZoneIDFormatRule(),
					Message: `"cn-shanghai-b" is not a zone in provider region "cn-hangzhou"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 13},
						End:      hcl.Pos{Line: 7, Column: 28},
					},
				},
			},
		},
		{
			Name: "region used as zone",
			Content: `
provider "alicloud" {
  region = "ap-southeast-1"
}

resource "alicloud_vswitch" "web" {
  zone_id = "ap-southeast-1"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudVswitchInvalidZoneIDFormatRule(),
					Message: `"ap-southeast-1" is not a zone in provider region "ap-southeast-1"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 13},
						End:      hcl.Pos{Line: 7, Column: 29},
					},
				},
			},
		},
		{
			Name: "aliased provider",
			Content: `
provider "alicloud" {
  region = "cn-hangzhou"
}

provider "alicloud" {
  alias  = "sg"
  region = "ap-southeast-1"
}

resource "alicloud_vswitch" "default" {
  zone_id = "cn-hangzhou-h"
}

resource "alicloud_vswitch" "sg" {
  provider = alicloud.sg
  zone_id  = "cn-hangzhou-h"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudVswitchInvalidZoneIDFormatRule(),
					Message: `"cn-hangzhou-h" is not a zone in provider region "ap-southeast-1"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 17, Column: 14},
						End:      hcl.Pos{Line: 17, Column: 29},
					},
				},
			},
		},
		{
			Name: "zone from variable default",
			Content: `
provider "alicloud" {
  region = "cn-beijing"
}

variable "zone" {
  default = "cn-shanghai-b"
}

resource "alicloud_vswitch" "web" {
  zone_id = var.zone
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudVswitchInvalidZoneIDFormatRule(),
					Message: `"cn-shanghai-b" is not a zone in provider region "cn-beijing"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 11, Column: 13},
						End:      hcl.Pos{Line: 11, Column: 21},
					},
				},
			},
		},
		{
			Name: "valid zones",
			Content: `
provider "alicloud" {
  region = "cn-hangzhou"
}

provider "alicloud" {
  alias  = "sg"
  region = "ap-southeast-1"
}

resource "alicloud_vswitch" "hz" {
  zone_id = "cn-hangzhou-h"
}

resource "alicloud_vswitch" "sg" {
  provider = alicloud.sg
  zone_id  = "ap-southeast-1a"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no provider block",
			Content: `
resource "alicloud_vswitch" "web" {
  zone_id = "cn-shanghai-b"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "provider without region",
			Content: `
provider "alicloud" {
  profile = "default"
}

resource "alicloud_vswitch" "web" {
  zone_id = "cn-shanghai-b"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "alias without matching provider",
			Content: `
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vswitch" "web" {
  provider = alicloud.other
  zone_id  = "cn-shanghai-b"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAlicloudVswitchInvalidZoneIDFormatRule()

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

func Test_zoneInRegion(t *testing.T) {
	cases := []struct {
		zone   string
		region string
		want   bool
	}{
		{"cn-hangzhou-h", "cn-hangzhou", true},
		{"cn-wuhan-lr-a", "cn-wuhan-lr", true},
		{"ap-southeast-1a", "ap-southeast-1", true},
		{"eu-central-1c", "eu-central-1", true},
		{"cn-hangzhou", "cn-hangzhou", false},
		{"cn-shanghai-b", "cn-hangzhou", false},
		{"ap-southeast-10a", "ap-southeast-1", false},
		{"cn-hangzhou_h", "cn-hangzhou", false},
	}

	for _, tc := range cases {
		if got := zoneInRegion(tc.zone, tc.region); got != tc.want {
			t.Errorf("zoneInRegion(%q, %q) = %t, want %t", tc.zone, tc.region, got, tc.want)
		}
	}
}

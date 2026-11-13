package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AlicloudProviderLegacySource(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "hashicorp source",
			Content: `
terraform {
  required_providers {
    alicloud = {
      source  = "hashicorp/alicloud"
      version = "~> 1.200"
    }
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudProviderLegacySourceRule(),
					Message: `"hashicorp/alicloud" is not the canonical provider source, use "aliyun/alicloud"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 17},
						End:      hcl.Pos{Line: 5, Column: 37},
					},
				},
			},
		},
		{
			Name: "registry host and mixed case",
			Content: `
terraform {
  required_providers {
    ali = {
      source = "registry.terraform.io/HashiCorp/alicloud"
    }
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudProviderLegacySourceRule(),
					Message: `"registry.terraform.io/HashiCorp/alicloud" is not the canonical provider source, use "aliyun/alicloud"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 16},
						End:      hcl.Pos{Line: 5, Column: 58},
					},
				},
			},
		},
		{
			Name: "quoted source key",
			Content: `
terraform {
  required_providers {
    alicloud = {
      "source" = "hashicorp/alicloud"
    }
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewAlicloudProviderLegacySourceRule(),
					Message: `"hashicorp/alicloud" is not the canonical provider source, use "aliyun/alicloud"`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 18},
						End:      hcl.Pos{Line: 5, Column: 38},
					},
				},
			},
		},
		{
			Name: "aliyun source",
			Content: `
terraform {
  required_providers {
    alicloud = {
      source  = "aliyun/alicloud"
      version = "~> 1.200"
    }
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "configuration aliases",
			Content: `
terraform {
  required_providers {
    alicloud = {
      source                = "aliyun/alicloud"
      configuration_aliases = [alicloud.hz]
    }
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "legacy version string",
			Content: `
terraform {
  required_providers {
    alicloud = "~> 1.200"
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "other hashicorp provider",
			Content: `
terraform {
  required_providers {
    random = {
      source = "hashicorp/random"
    }
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no required_providers",
			Content: `
terraform {
  required_version = ">= 1.5"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAlicloudProviderLegacySourceRule()

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

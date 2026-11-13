package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/rezNetDevOps/tflint-ruleset-alicloud/project"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AlicloudProviderLegacySourceRule checks that required_providers does not use the hashicorp/alicloud source.
type AlicloudProviderLegacySourceRule struct {
	tflint.DefaultRule
}

// NewAlicloudProviderLegacySourceRule returns a new rule.
func NewAlicloudProviderLegacySourceRule() *AlicloudProviderLegacySourceRule {
	return &AlicloudProviderLegacySourceRule{}
}

// Name returns the rule name.
func (r *AlicloudProviderLegacySourceRule) Name() string {
	return "alicloud_provider_legacy_source"
}

// Enabled returns whether the rule is enabled by default.
func (r *AlicloudProviderLegacySourceRule) Enabled() bool {
	return false
}

// Severity returns the rule severity.
func (r *AlicloudProviderLegacySourceRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link.
func (r *AlicloudProviderLegacySourceRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check reports required_providers entries whose source is hashicorp/alicloud.
func (r *AlicloudProviderLegacySourceRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "terraform",
				Body: &hclext.BodySchema{
					Blocks: []hclext.BlockSchema{
						{
							Type: "required_providers",
							Body: &hclext.BodySchema{Mode: hclext.SchemaJustAttributesMode},
						},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, terraform := range content.Blocks {
		for _, providers := range terraform.Body.Blocks {
			for _, attr := range providers.Body.Attributes {
				pairs, diags := hcl.ExprMap(attr.Expr)
				if diags.HasErrors() {
					continue
				}
				for _, pair := range pairs {
					if objectKey(pair.Key) != "source" {
						continue
					}
					val, diags := pair.Value.Value(nil)
					if diags.HasErrors() || !val.IsKnown() || val.IsNull() || !val.Type().Equals(cty.String) {
						continue
					}
					if !isLegacySource(val.AsString()) {
						continue
					}
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("%q is not the canonical provider source, use \"aliyun/alicloud\"", val.AsString()),
						pair.Value.Range(),
					); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func objectKey(expr hcl.Expression) string {
	if key := hcl.ExprAsKeyword(expr); key != "" {
		return key
	}
	val, diags := expr.Value(nil)
	if diags.HasErrors() || !val.IsKnown() || val.IsNull() || !val.Type().Equals(cty.String) {
		return ""
	}
	return val.AsString()
}

func isLegacySource(source string) bool {
	return strings.TrimPrefix(strings.ToLower(source), "registry.terraform.io/") == "hashicorp/alicloud"
}

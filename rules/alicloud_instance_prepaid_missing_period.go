package rules

import (
	"github.com/rezNetDevOps/tflint-ruleset-alicloud/project"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AlicloudInstancePrepaidMissingPeriodRule checks that a PrePaid alicloud_instance sets period.
type AlicloudInstancePrepaidMissingPeriodRule struct {
	tflint.DefaultRule
}

// NewAlicloudInstancePrepaidMissingPeriodRule returns a new rule.
func NewAlicloudInstancePrepaidMissingPeriodRule() *AlicloudInstancePrepaidMissingPeriodRule {
	return &AlicloudInstancePrepaidMissingPeriodRule{}
}

// Name returns the rule name.
func (r *AlicloudInstancePrepaidMissingPeriodRule) Name() string {
	return "alicloud_instance_prepaid_missing_period"
}

// Enabled returns whether the rule is enabled by default.
func (r *AlicloudInstancePrepaidMissingPeriodRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AlicloudInstancePrepaidMissingPeriodRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AlicloudInstancePrepaidMissingPeriodRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check reports alicloud_instance resources with instance_charge_type "PrePaid" and no period.
func (r *AlicloudInstancePrepaidMissingPeriodRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("alicloud_instance", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "instance_charge_type"},
			{Name: "period"},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		chargeType, exists := resource.Body.Attributes["instance_charge_type"]
		if !exists {
			continue
		}

		err := runner.EvaluateExpr(chargeType.Expr, func(value string) error {
			if value != "PrePaid" {
				return nil
			}

			if period, exists := resource.Body.Attributes["period"]; exists {
				var val cty.Value
				if err := runner.EvaluateExpr(period.Expr, &val, nil); err != nil {
					return err
				}
				if !val.IsNull() {
					return nil
				}
			}

			return runner.EmitIssue(
				r,
				`"period" is required when instance_charge_type is "PrePaid"`,
				chargeType.Expr.Range(),
			)
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

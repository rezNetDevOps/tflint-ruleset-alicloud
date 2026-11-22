package rules

import (
	"fmt"
	"regexp"

	"github.com/rezNetDevOps/tflint-ruleset-alicloud/project"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var instanceTypeFormat = regexp.MustCompile(`^ecs\.[a-z0-9-]+\.[a-z0-9]+$`)

// AlicloudInstanceInvalidInstanceTypeFormatRule checks that instance_type has the form ecs.<family>.<size>.
type AlicloudInstanceInvalidInstanceTypeFormatRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAlicloudInstanceInvalidInstanceTypeFormatRule returns a new rule.
func NewAlicloudInstanceInvalidInstanceTypeFormatRule() *AlicloudInstanceInvalidInstanceTypeFormatRule {
	return &AlicloudInstanceInvalidInstanceTypeFormatRule{
		resourceType:  "alicloud_instance",
		attributeName: "instance_type",
	}
}

// Name returns the rule name.
func (r *AlicloudInstanceInvalidInstanceTypeFormatRule) Name() string {
	return "alicloud_instance_invalid_instance_type_format"
}

// Enabled returns whether the rule is enabled by default.
func (r *AlicloudInstanceInvalidInstanceTypeFormatRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AlicloudInstanceInvalidInstanceTypeFormatRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AlicloudInstanceInvalidInstanceTypeFormatRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check reports instance_type values that do not have the form ecs.<family>.<size>.
func (r *AlicloudInstanceInvalidInstanceTypeFormatRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{{Name: r.attributeName}},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		err := runner.EvaluateExpr(attribute.Expr, func(instanceType string) error {
			if instanceTypeFormat.MatchString(instanceType) {
				return nil
			}
			return runner.EmitIssue(
				r,
				fmt.Sprintf("%q is an invalid instance type, expected ecs.<family>.<size>", instanceType),
				attribute.Expr.Range(),
			)
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

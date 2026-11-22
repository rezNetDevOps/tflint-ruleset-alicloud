package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/rezNetDevOps/tflint-ruleset-alicloud/project"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AlicloudVswitchInvalidZoneIDFormatRule checks that a vswitch zone_id belongs to the provider region.
type AlicloudVswitchInvalidZoneIDFormatRule struct {
	tflint.DefaultRule
}

// NewAlicloudVswitchInvalidZoneIDFormatRule returns a new rule.
func NewAlicloudVswitchInvalidZoneIDFormatRule() *AlicloudVswitchInvalidZoneIDFormatRule {
	return &AlicloudVswitchInvalidZoneIDFormatRule{}
}

// Name returns the rule name.
func (r *AlicloudVswitchInvalidZoneIDFormatRule) Name() string {
	return "alicloud_vswitch_invalid_zone_id_format"
}

// Enabled returns whether the rule is enabled by default.
func (r *AlicloudVswitchInvalidZoneIDFormatRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AlicloudVswitchInvalidZoneIDFormatRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AlicloudVswitchInvalidZoneIDFormatRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check reports alicloud_vswitch zone_id values that do not start with the provider region.
func (r *AlicloudVswitchInvalidZoneIDFormatRule) Check(runner tflint.Runner) error {
	regions, err := r.providerRegions(runner)
	if err != nil {
		return err
	}
	if len(regions) == 0 {
		return nil
	}

	resources, err := runner.GetResourceContent("alicloud_vswitch", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "zone_id"},
			{Name: "provider"},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		zoneID, exists := resource.Body.Attributes["zone_id"]
		if !exists {
			continue
		}

		alias, ok := providerAlias(resource.Body.Attributes["provider"])
		if !ok {
			continue
		}
		region, ok := regions[alias]
		if !ok {
			continue
		}

		err := runner.EvaluateExpr(zoneID.Expr, func(zone string) error {
			if zoneInRegion(zone, region) {
				return nil
			}
			return runner.EmitIssue(
				r,
				fmt.Sprintf("%q is not a zone in provider region %q", zone, region),
				zoneID.Expr.Range(),
			)
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *AlicloudVswitchInvalidZoneIDFormatRule) providerRegions(runner tflint.Runner) (map[string]string, error) {
	providers, err := runner.GetProviderContent("alicloud", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "alias"},
			{Name: "region"},
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	regions := map[string]string{}
	for _, provider := range providers.Blocks {
		region, exists := provider.Body.Attributes["region"]
		if !exists {
			continue
		}

		alias := ""
		if attr, exists := provider.Body.Attributes["alias"]; exists {
			if err := runner.EvaluateExpr(attr.Expr, &alias, nil); err != nil {
				continue
			}
		}

		err := runner.EvaluateExpr(region.Expr, func(value string) error {
			regions[alias] = value
			return nil
		}, nil)
		if err != nil {
			return nil, err
		}
	}

	return regions, nil
}

func providerAlias(attr *hclext.Attribute) (string, bool) {
	if attr == nil {
		return "", true
	}
	traversal, diags := hcl.AbsTraversalForExpr(attr.Expr)
	if diags.HasErrors() || traversal.RootName() != "alicloud" {
		return "", false
	}
	switch len(traversal) {
	case 1:
		return "", true
	case 2:
		if step, ok := traversal[1].(hcl.TraverseAttr); ok {
			return step.Name, true
		}
	}
	return "", false
}

func zoneInRegion(zone, region string) bool {
	suffix, found := strings.CutPrefix(zone, region)
	if !found || suffix == "" {
		return false
	}
	return suffix[0] == '-' || ('a' <= suffix[0] && suffix[0] <= 'z')
}

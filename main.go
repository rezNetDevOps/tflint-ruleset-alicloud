package main

import (
	"github.com/rezNetDevOps/tflint-ruleset-alicloud/project"
	"github.com/rezNetDevOps/tflint-ruleset-alicloud/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "alicloud",
			Version: project.Version,
			Rules:   rules.Rules,
		},
	})
}

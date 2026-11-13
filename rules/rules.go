package rules

import "github.com/terraform-linters/tflint-plugin-sdk/tflint"

// Rules is the list of all rules in this ruleset.
var Rules = []tflint.Rule{
	NewAlicloudProviderLegacySourceRule(),
}

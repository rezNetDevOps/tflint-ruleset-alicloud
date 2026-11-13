package project

import "fmt"

// Version is the ruleset version.
const Version string = "0.1.0"

// ReferenceLink returns the documentation URL for the named rule.
func ReferenceLink(name string) string {
	return fmt.Sprintf("https://github.com/rezNetDevOps/tflint-ruleset-alicloud/blob/v%s/docs/rules/%s.md", Version, name)
}

package main

import (
	"github.com/richardTowers/tflint-ruleset-workspaces/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "workspaces",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewResourceNamesIncludeWorkspaceRule("aws_s3_bucket", "bucket"),
				rules.NewResourceNamesIncludeWorkspaceRule("aws_security_group", "name"),
			},
		},
	})
}

package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &RuleSet{
			BuiltinRuleSet: tflint.BuiltinRuleSet{
				Name:    "workspaces",
				Version: "0.1.0",
				Rules:   []tflint.Rule{},
			},
		},
	})
}

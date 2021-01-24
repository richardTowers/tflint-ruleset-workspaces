package main

import (
	"fmt"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/richardTowers/tflint-ruleset-workspaces/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type ResourceConfig struct {
	Resource      string `hcl:"resource,label"`
	AttributeName string `hcl:"attribute,optional"`
}

type Config struct {
	Resources         []ResourceConfig `hcl:"resource,block"`
	OverrideWorkspace string           `hcl:"override_workspace,optional"`
	Remain hcl.Body `hcl:",remain"`
}

type RuleSet struct {
	tflint.BuiltinRuleSet
	config *Config
}

func (r *RuleSet) RuleNames() []string {
	result := []string{}
	for _, rule := range r.rules() {
		result = append(result, rule.Name())
	}
	return result
}

func (r *RuleSet) rules() []*rules.ResourceNamesIncludeWorkspace {
	result := []*rules.ResourceNamesIncludeWorkspace{}
	for _, resource := range r.config.Resources {
		rule := rules.NewResourceNamesIncludeWorkspaceRule(
			resource.Resource,
			resource.AttributeName,
			r.config.OverrideWorkspace,
		)
		result = append(result, rule)
	}
	return result
}

func (r *RuleSet) ApplyConfig(config *tflint.Config) error {
	r.ApplyCommonConfig(config)

	cfg := Config{}
	diags := gohcl.DecodeBody(config.Body, nil, &cfg)
	if diags.HasErrors() {
		return diags
	}
	r.config = &cfg
	return nil
}

func (r *RuleSet) Check(runner tflint.Runner) error {
	for _, rule := range r.rules() {
		if err := rule.Check(runner); err != nil {
			return fmt.Errorf("failed to check `%s` rule: %s", rule.Name(), err)
		}
	}
	return nil
}

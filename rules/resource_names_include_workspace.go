package rules

import (
	"fmt"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"os"
	"strings"
)

type ResourceNamesIncludeWorkspace struct{
	resourceType string
	attributeName string
}

type resourceNamesIncludeWorkspaceRuleConfig struct {
	OverrideWorkspace string `hcl:"overrideWorkspace,optional"`
}

func NewResourceNamesIncludeWorkspaceRule(resourceType, attributeName string) *ResourceNamesIncludeWorkspace {
	return &ResourceNamesIncludeWorkspace{
		resourceType: resourceType,
		attributeName: attributeName,
	}
}

func (r *ResourceNamesIncludeWorkspace) Name() string {
	return fmt.Sprintf("%s_name_includes_workspace", r.resourceType)
}

func (r *ResourceNamesIncludeWorkspace) Enabled() bool {
	return true
}

func (r *ResourceNamesIncludeWorkspace) Severity() string {
	return tflint.ERROR
}

func (r *ResourceNamesIncludeWorkspace) Link() string {
	return "https://github.com/richardTowers/tflint-ruleset-workspaces/blob/master/README.md"
}

func (r *ResourceNamesIncludeWorkspace) Check(runner tflint.Runner) error {
	config := resourceNamesIncludeWorkspaceRuleConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
		return err
	}
	originalWorkspace, envWasSet := os.LookupEnv("TERRAFORM_WORKSPACE")
	defer func() {
		if envWasSet {
			_ = os.Setenv("TERRAFORM_WORKSPACE", originalWorkspace)
		} else {
			_ = os.Unsetenv("TERRAFORM_WORKSPACE")
		}
	}()

	tempWorkspace := "default"
	if originalWorkspace != "" {
		tempWorkspace = originalWorkspace
	}
	if config.OverrideWorkspace != "" {
		tempWorkspace = config.OverrideWorkspace
	}
	_ = os.Setenv("TERRAFORM_WORKSPACE", tempWorkspace)
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var name string
		err := runner.EvaluateExpr(attribute.Expr, &name, nil)

		return runner.EnsureNoError(err, func() error {
			if !strings.Contains(name, tempWorkspace) {
				_ = runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf(`%s resource name "%s" does not include the workspace (%s)`, r.resourceType, name, tempWorkspace),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}

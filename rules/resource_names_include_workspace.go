package rules

import (
	"fmt"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"os"
	"strings"
)

type ResourceNamesIncludeWorkspace struct {
	resourceType string
	attributeName string
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
	workspace, success := os.LookupEnv("TF_WORKSPACE")
	if !success {
		workspace = "default"
	}

	attributeName := r.attributeName
	if attributeName == "" {
		attributeName = "name"
	}
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var name string
		err := runner.EvaluateExpr(attribute.Expr, &name, nil)

		return runner.EnsureNoError(err, func() error {
			if !strings.Contains(name, workspace) {
				_ = runner.EmitIssueOnExpr(
					r,
					fmt.Sprintf(`%s resource name "%s" does not include the workspace (%s)`, r.resourceType, name, workspace),
					attribute.Expr,
				)
			}
			return nil
		})
	})
}

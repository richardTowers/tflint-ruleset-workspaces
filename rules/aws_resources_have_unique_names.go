package rules

import (
	"fmt"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"strings"
)

type AwsResourcesHaveUniqueNamesRule struct{
	resourceType string
	attributeName string
}

func NewAwsResourcesHaveUniqueNamesRule(resourceType, attributeName string) *AwsResourcesHaveUniqueNamesRule {
	return &AwsResourcesHaveUniqueNamesRule{
		resourceType: resourceType,
		attributeName: attributeName,
	}
}

func (r *AwsResourcesHaveUniqueNamesRule) Name() string {
	return fmt.Sprintf("aws_%s_resources_have_unique_names", r.resourceType)
}

func (r *AwsResourcesHaveUniqueNamesRule) Enabled() bool {
	return true
}

func (r *AwsResourcesHaveUniqueNamesRule) Severity() string {
	return tflint.ERROR
}

func (r *AwsResourcesHaveUniqueNamesRule) Link() string {
	return "https://github.com/richardTowers/tflint-ruleset-workspaces/blob/master/README.md"
}

func (r *AwsResourcesHaveUniqueNamesRule) Check(runner tflint.Runner) error {
	return runner.WalkResourceAttributes(r.resourceType, r.attributeName, func(attribute *hcl.Attribute) error {
		var name string
		err := runner.EvaluateExpr(attribute.Expr, &name, nil)

		return runner.EnsureNoError(err, func() error {
			if !strings.Contains(name, "default") {
				_ = runner.EmitIssue(
					r,
					fmt.Sprintf(`Resource name "%s" does not include the workspace`, name),
					attribute.Expr.Range(),
				)
			}
			return nil
		})
	})
}

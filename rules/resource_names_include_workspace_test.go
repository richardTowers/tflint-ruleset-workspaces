package rules

import (
	"os"
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketName(t *testing.T) {
	rule := NewResourceNamesIncludeWorkspaceRule("aws_s3_bucket", "bucket")

	content := `
resource "aws_s3_bucket" "good" {
 bucket = "${terraform.workspace}-bucket"
}
resource "aws_s3_bucket" "bad" {
	bucket = "no-workspace-bucket"
}
`

	t.Run("s3_bucket_without_workspace_prefixed_name", func(t *testing.T) {
		_ = os.Setenv("TF_WORKSPACE", "overridden_workspace")
		// TODO remove once https://github.com/terraform-linters/tflint-plugin-sdk/pull/101 is merged:
		_ = os.Setenv("TERRAFORM_WORKSPACE", "overridden_workspace")

		defer func(){
			_ = os.Unsetenv("TF_WORKSPACE")
			_ = os.Unsetenv("TERRAFORM_WORKSPACE")
		}()

		runner := helper.TestRunner(t, map[string]string{"resource.tf": content, ".tflint.hcl": ""})

		err := rule.Check(runner)
		if err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, helper.Issues{
			{
				Rule:    NewResourceNamesIncludeWorkspaceRule("aws_s3_bucket", "bucket"),
				Message: `aws_s3_bucket resource name "no-workspace-bucket" does not include the workspace (overridden_workspace)`,
				Range: hcl.Range{
					Filename: "resource.tf",
					Start:    hcl.Pos{Line: 6, Column: 11},
					End:      hcl.Pos{Line: 6, Column: 32},
				},
			},
		}, runner.Issues)
	})
}

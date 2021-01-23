package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketName(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
		Error    bool
	}{
		{
			Name: "s3_bucket_without_workspace_prefixed_name",
			Content: `
resource "aws_s3_bucket" "good" {
  bucket = "${terraform.workspace}-bucket"
}
resource "aws_s3_bucket" "bad" {
	bucket = "no-workspace-bucket"
}`,
			Config: `
rule "aws_s3_bucket_name_includes_workspace" {
	enabled = true
	overrideWorkspace = "overridden_workspace"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewResourceNamesIncludeWorkspaceRule("aws_s3_bucket", "bucket"),
					Message: `aws_s3_bucket resource name "no-workspace-bucket" does not include the workspace (overridden_workspace)`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 11},
						End:      hcl.Pos{Line: 6, Column: 32},
					},
				},
			},
		},
	}

	rule := NewResourceNamesIncludeWorkspaceRule("aws_s3_bucket", "bucket")

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content, ".tflint.hcl": tc.Config})

			err := rule.Check(runner)
			if err != nil && !tc.Error {
				t.Fatalf("Unexpected error occurred: %s", err)
			}
			if err == nil && tc.Error {
				t.Fatal("Expected error but got none")
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
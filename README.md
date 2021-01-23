# TFLint Ruleset Workspaces
[![Build Status](https://github.com/richardTowers/tflint-ruleset-workspaces/workflows/build/badge.svg?branch=master)](https://github.com/richardTowers/tflint-ruleset-workspaces/actions)

This custom ruleset checks that resource names are uniquely namespaced by workspaces. It's useful if you're using [terraform workspaces](https://www.terraform.io/docs/state/workspaces.html)
to deploy multiple instances of the same configuration.

## Requirements

- TFLint v0.23+
- Go v1.15

## Installation

Download the plugin and place it in `~/.tflint.d/plugins/tflint-ruleset-template` (or `./.tflint.d/plugins/tflint-ruleset-template`). When using the plugin, configure as follows in `.tflint.hcl`:

```hcl
plugin "workspaces" {
    enabled = true
}
```

## Configuration

Unlike most tflint plugins, this one takes all of its config in the plugin block.

This allows each project to configure which resources should be namespaced by `${terraform.workspace}`. This entirely depends on how the project is using workspaces.

```hcl
plugin "workspaces" {
  enabled            = true
  override_workspace = "xxxx"

  resource "aws_s3_bucket" {
    attribute = "bucket"
  }
  resource "aws_security_group" {
    attribute = "name"
  }
}
```

One rule will be created for each configured `resource`.

`attribute` specifies the attribute which needs to include the workspace - e.g. for S3 buckets this should be "bucket", for security groups it should be "name" (assuming you have multiple workspaces in the same VPC, otherwise namespacing your security groups is not needed). Defaults to `name`.

`override_workspace` forces tflint to use a non-default workspace for this plugin. This is useful if you treat the default workspace as a special case (e.g. with `name = terrafrom.workspace == "default" ? "blah" : "${terraform.workspace}_blah"`)

## Building the plugin

Clone the repository locally and run the following command:

```
$ make
```

You can easily install the built plugin with the following:

```
$ make install
```

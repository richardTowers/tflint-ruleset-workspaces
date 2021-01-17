# TFLint Ruleset Workspaces
[![Build Status](https://github.com/terraform-linters/tflint-ruleset-template/workflows/build/badge.svg?branch=master)](https://github.com/terraform-linters/tflint-ruleset-template/actions)

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

## Rules

|Name|Description|Severity|Enabled|Link|
| --- | --- | --- | --- | --- |
|todo|todo|ERROR|✔||

## Building the plugin

Clone the repository locally and run the following command:

```
$ make
```

You can easily install the built plugin with the following:

```
$ make install
```

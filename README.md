# MyApps

**MyApps** is a universal command line tool for managing manually installed applications.

> **Disclaimer** 
> 
> I wrote this tool over two long nights while preparing my reports and studying for college. I don't have the strength or time to write nice documentation. Just check the built-in help.
> 
> _Linter, linter of the Go, code's not prettiest of them all._

[![asciicast](https://asciinema.org/a/451017.svg)](https://asciinema.org/a/451017)


## Installation

```bash
# Let's say that you downloaded the latest release asset "myapps.tgz"
tar -xzf myapps.tgz
sudo ./install.sh
```

## Usage

Just run the `myapps` command and check usage information.

## Repositories

Use subcommands of `myapps repo` to manage the repositories.

<details>
<summary><strong>Example repository structure</strong></summary>

```yaml
name: default
description: Default repository
maintainer: MyApps
contents: []
  # - name: applicationName
  #   description: Application description
  #   install_scenario:
  #     - some_commands
  #     - some_other_commands
  #   update_scenario:
  #     - some_commands
  #     - some_other_commands
  #   uninstall_scenario:
  #     - some_commands
  #     - some_other_commands
  #   newest_version_check:
  #     - some_commands
  #     - some_other_commands
  #   current_version_check:
  #     - some_commands
  #     - some_other_commands

# Those environment variables are available inside the scenarios:
# - OS (in all scenarios)
# - ARCH (in all scenarios)
# - LATEST_VERSION (in install & update scenarios)
# - CURRENT_VERSION (in update & uninstall scenarios)
# - TEMP - path to temporary directory that will be removed after scenario (in install & update scenarios)
```
</details>

The repository schema is available in `repository.schema.json` file.

## Configuration

Configuration file is located in `/usr/share/myapps/config.yaml` and can be edited by the user.

<details>
<summary><strong>Configuration file structure</strong></summary>

```yaml
# Name of default repository
default_repo: default

# GitHub repository loader configuration
github_loader:
  # Enable automatic detection of default repository branch
  detect_default_branch: true
  # If the automatic detection is disabled, use following branch name as default (when not explicitly specified)
  default_branch_name: master
```

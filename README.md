

# MyApps

**MyApps** is a universal command line tool for managing manually installed applications.

> **Disclaimer** 
> 
> I wrote this tool over two long nights while preparing my reports and studying for college. I don't have the strength or time to write nice documentation. Just check the built-in help.
> 
> _Linter, linter of the Go, code's not prettiest of them all._

## Installation

```bash
# Let's say that you downloaded the latest release asset "myapps.tgz"
tar -xzf myapps.tgz
sudo ./install.sh
```

## Usage

Just run the `myapps` command and check usage information.

## Repositories

Use subcommands of `myapps repos` to manage the repositories.

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
  #   latest_version_check:
  #     - some_commands
  #     - some_other_commands

# Those environment variables are available inside the scenarios:
# - LATEST_VERSION (in install & update scenarios)
# - CURRENT_VERSION (in update & uninstall scenario)
```
</details>
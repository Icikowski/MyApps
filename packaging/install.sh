#!/bin/bash

echo "MyApps installer/updater"
echo "Run this as root or it won't work properly!"

echo " -> Copying application"
rm -rf /usr/bin/myapps
cp myapps /usr/bin/myapps
chmod a+x /usr/bin/myapps

echo " -> Creating directories"
if [ ! -d /usr/share/myapps ]; then
  mkdir /usr/share/myapps /usr/share/myapps/repos
fi
if [ ! -d /usr/share/myapps/repos ]; then
  mkdir /usr/share/myapps/repos
fi

echo " -> Creating files..."
if [ ! -e /usr/share/myapps/config.yaml ]; then
  cat >/usr/share/myapps/config.yaml <<EOF
default_repo: default
EOF
fi

if [ ! -e /usr/share/myapps/deployments.yaml ]; then
  cat >/usr/share/myapps/deployments.yaml <<EOF
[]
EOF
fi

if [ ! -e /usr/share/myapps/repos/default.yaml ]; then
  cat >/usr/share/myapps/repos/default.yaml <<EOF
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
# - LATEST_VERSION (in install & update scenarios)
# - CURRENT_VERSION (in update & uninstall scenario)
EOF
fi

echo " -> Enabling bash completion"
cat >/etc/bash_completion.d/myapps <<'EOF'
#! /bin/bash

_cli_bash_autocomplete() {
  if [[ "${COMP_WORDS[0]}" != "source" ]]; then
    local cur opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ "$cur" == "-"* ]]; then
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} ${cur} --generate-bash-completion )
    else
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    fi
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
  fi
}

complete -o bashdefault -o default -o nospace -F _cli_bash_autocomplete myapps
EOF
chmod a+x /etc/bash_completion.d/myapps

echo " -> DONE"
echo ""
echo "You can change the default repository by editing file:"
echo "  /usr/share/myapps/config.yaml"
echo "Check for more pre-made repositories:"
echo "  https://github.com/Icikowski/MyApps/tree/master/repositories"
echo ""
echo "Have fun!"
#!/bin/bash

echo "MyApps installer"
echo "Run this as root or it won't work properly!"

echo " -> Copying application"
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
cat >/usr/share/myapps/config.yaml <<EOF
default_repo: default
EOF
cat >/usr/share/myapps/deployments.yaml <<EOF
[]
EOF
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
  #   latest_version_check:
  #     - some_commands
  #     - some_other_commands

# Those environment variables are available inside the scenarios:
# - LATEST_VERSION (in install & update scenarios)
# - CURRENT_VERSION (in update & uninstall scenario)
EOF
cat >/usr/share/myapps/repos/icikowski.yaml <<'EOF'
name: icikowski
description: Repository for my personal projects
maintainer: Piotr Icikowski
contents:
  - name: gpts
    description: General Purpose Test Service
    install_scenario: 
      - rm -rf /tmp/gpts-install
      - mkdir /tmp/gpts-install
      - wget -q -O /tmp/gpts-install/gpts https://github.com/Icikowski/GPTS/releases/download/v${LATEST_VERSION}/gpts-${LATEST_VERSION}-linux-amd64
      - mv /tmp/gpts-install/gpts /usr/local/bin/gpts
      - rm -rf /tmp/gpts-install
      - chmod a+x /usr/local/bin/gpts
    update_scenario:
      - rm -f /usr/local/bin/gpts
      - rm -rf /tmp/gpts-install
      - mkdir /tmp/gpts-install
      - wget -q -O /tmp/gpts-install/gpts https://github.com/Icikowski/GPTS/releases/download/v${LATEST_VERSION}/gpts-${LATEST_VERSION}-linux-amd64
      - mv /tmp/gpts-install/gpts /usr/local/bin/gpts
      - rm -rf /tmp/gpts-install
      - chmod a+x /usr/local/bin/gpts
    uninstall_scenario:
      - rm -f /usr/local/bin/gpts
    newest_version_check:
      - curl -s https://charts.icikowski.pl/index.json | jq -r ".entries.gpts[0].version"
    current_version_check:
      - gpts --service-port -1 | head -n 1 | jq -r ".version"
EOF

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
echo "You can check the repository structure in files:"
echo "  /usr/share/myapps/repos/default.yaml"
echo "  /usr/share/myapps/repos/icikowski.yaml"
echo ""
echo "Have fun!"
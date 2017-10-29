#!/usr/bin/env zsh
#
# create a Go workspace
# Author: edoz90
# Usage: `source goenv.zsh [-u|--update]`
#

_GO_PROJECT_DIR="${_:a:h}"

if [[ $(command -v go) ]]; then
  # Default folders for golang workspace
  mkdir -p ${_GO_PROJECT_DIR}/{src,pkg,bin}
  touch requirements.txt

  TEMP_GOPATH="${GOPATH##*:}"
  GOPATH="${_GO_PROJECT_DIR}"
  typeset -U PATH="${GOPATH}/bin:${PATH}"

  # Update all local (and only local) pkgs
  if [[ $1 = "--update" || $1 = "-u" ]]; then
    go get -u all
  fi

  # Add priority to the workspace directories
  GOPATH="${_GO_PROJECT_DIR}:${TEMP_GOPATH}"

  # set PS1 bash
  PS1="$(basename ${_GO_PROJECT_DIR}):${PS1}"

  # Remove all settings
  alias deactivate="deactivate"
else
  echo "Please install Go first"
fi

function deactivate {
  PATH="$(${PATH} | sed -e "s|${_GO_PROJECT_DIR}/bin:||g")"
  unset _GO_PROJECT_DIR
}

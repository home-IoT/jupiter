#
# Common utility functions for other scripts
#
# Copyright (c) 2016-2018 Roozbeh Farahbod
# Distributed under the MIT License.
#

CHILD_SCRIPT=""
CUR_DIR="$PWD"

printError() {
  (>&2 echo "ERROR: $*")
}

showUsage() {
  MSG=$(grep -e '^# Usage' "$1" | sed 's/^#\ //')
  (>&2 echo "$MSG")
}

ensureCleanDevelopBranch() {
  REPO_DIR="$1"
  ensureCleanDevelopBranch_CUR_DIR="$PWD"

  cd "${REPO_DIR}"

  git checkout -q develop
  if [ "$?" != 0 ]
  then
    printError "Cannot swtich to branch 'develop'. You are most likely not on branch 'develop' and have uncommitted files."
    exit 1
  fi

  CHANGED_FILES=$(git status --porcelain)

  if [ ! -z "${CHANGED_FILES}" ]
  then
    printError "There are still files to be committed."
    exit 1
  fi

  cd "${ensureCleanDevelopBranch_CUR_DIR}"
}

# finds and/or validates the repository directory
getRepoDir() {
  grd_SCRIPT="$1"
  REPO_DIR="$2"
  grd_SCRIPT_DIR="$( cd "$( dirname "${grd_SCRIPT}" )" && pwd )"
  REPO_DIR=$("${grd_SCRIPT_DIR}/check-repo-dir" "${REPO_DIR}")
  grd_CHECK_DIR="$?"
  if [ "$grd_CHECK_DIR" != "0" ]
  then
    showUsage ${grd_SCRIPT}
    exit 1
  fi
}

if [ "$1" == "--sub-script" ]
then
  SUB_SCRIPT="true"
  shift
fi

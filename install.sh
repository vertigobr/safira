#!/usr/bin/env bash

# Copyright © Vertigo Tecnologia
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# shellcheck disable=SC2223
: ${BINARY_NAME:="safira"}
: ${REPO_URL:="https://github.com/kyfelipe/safira"}
: ${USE_SUDO:="true"}
: ${SAFIRA_INSTALL_DIR:="/usr/local/bin"}

# initArch discovers the architecture for this system.
initArch() {
  ARCH=$(uname -m)
  case $ARCH in
    armv5*) ARCH="armv5";;
    armv6*) ARCH="armv6";;
    armv7*) ARCH="arm";;
    aarch64) ARCH="arm64";;
    x86) ARCH="386";;
    x86_64) ARCH="amd64";;
    i686) ARCH="386";;
    i386) ARCH="386";;
  esac
}

# initOS discovers the operating system for this system.
initOS() {
  OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

  case "$OS" in
    # Minimalist GNU for Windows
    mingw*) OS='windows';;
  esac
}

# runs the given command as root (detects if we are root already)
runAsRoot() {
  local CMD="$*"

  if [ $EUID -ne 0 -a $USE_SUDO = "true" ]; then
    CMD="sudo $CMD"
  fi

  $CMD
}

# verifySupported checks that the os/arch combination is supported for binary builds.
verifySupported() {
  local supported="darwin-amd64\nlinux-386\nlinux-amd64\nlinux-arm\nlinux-arm64"
  if ! echo "${supported}" | grep -q "${OS}-${ARCH}"; then
    echo "No prebuilt binary for ${OS}-${ARCH}."
    echo "To build from source, go to $REPO_URL"
    exit 1
  fi

  if ! type "curl" > /dev/null && ! type "wget" > /dev/null; then
    echo "Either curl or wget is required"
    exit 1
  fi
}

# checkSafiraInstalledVersion checks which version of safira is installed and
# if it needs to be changed.
checkSafiraInstalledVersion() {
  if [[ -f "${SAFIRA_INSTALL_DIR}/${BINARY_NAME}" ]]; then
    local version=$(safira --version | cut -d " " -f3)
    if [[ "$version" == "$TAG" ]]; then
      echo "Safira ${version} is already ${DESIRED_VERSION:-latest}"
      return 0
    else
      echo "Safira ${TAG} is available. Changing from version ${version}."
      return 1
    fi
  else
    return 1
  fi
}

# checkTagProvided checks whether TAG has provided as an environment variable so we can skip checkLatestVersion.
checkTagProvided() {
  [[ ! -z "$TAG" ]]
}

# checkLatestVersion grabs the latest version string from the releases
checkLatestVersion() {
  local latest_release_url="$REPO_URL/releases/latest"
  if type "curl" > /dev/null; then
    TAG=$(curl -Ls -o /dev/null -w %{url_effective} $latest_release_url | grep -oE "[^/]+$" )
  elif type "wget" > /dev/null; then
    TAG=$(wget $latest_release_url --server-response -O /dev/null 2>&1 | awk '/^  Location: /{DEST=$2} END{ print DEST}' | grep -oE "[^/]+$")
  fi
}

# downloadFile downloads the latest binary package and also the checksum
# for that binary.
downloadFile() {
  SAFIRA_DIST="safira-$TAG-$OS-$ARCH.tar.gz"
  DOWNLOAD_URL="$REPO_URL/releases/download/$TAG/$SAFIRA_DIST"
  CHECKSUM_URL="$DOWNLOAD_URL.md5"
  SAFIRA_TMP_ROOT="$(mktemp -dt safira-installer-XXXXXX)"
  SAFIRA_TMP_FILE="$SAFIRA_TMP_ROOT/$SAFIRA_DIST"
  SAFIRA_SUM_FILE="$SAFIRA_TMP_ROOT/$SAFIRA_DIST.md5"
  echo "Downloading $DOWNLOAD_URL"
  if type "curl" > /dev/null; then
    curl -SsL "$CHECKSUM_URL" -o "$SAFIRA_SUM_FILE"
  elif type "wget" > /dev/null; then
    wget -q -O "$SAFIRA_SUM_FILE" "$CHECKSUM_URL"
  fi
  if type "curl" > /dev/null; then
    curl -SsL "$DOWNLOAD_URL" -o "$SAFIRA_TMP_FILE"
  elif type "wget" > /dev/null; then
    wget -q -O "$SAFIRA_TMP_FILE" "$DOWNLOAD_URL"
  fi
}

# Added hostnames for develop in local environment
addHost() {
  HOSTS_FILE="/etc/hosts"
  HOST_NAMES="127.0.0.1       ipaas.localdomain konga.localdomain gateway.ipaas.localdomain"
  FOUND=false
  while IFS="" read -r p || [ -n "$p" ]
  do
    if [[ "$p" == "$HOST_NAMES" ]]; then
      FOUND=true
      break
    fi
  done < $HOSTS_FILE

  if [ "$FOUND" = false ]; then
    echo "$HOST_NAMES" >> $HOSTS_FILE
  fi
}

#addPath() {
#  PROFILE_FILE="$HOME/.profile"
#  FOLDER_SAFIRA="PATH=\$PATH:$HOME/.safira/bin"
#  FOUND=false
#  while IFS="" read -r p || [ -n "$p" ]
#  do
#    if [[ "$p" == "$FOLDER_SAFIRA" ]]; then
#      FOUND=true
#      break
#    fi
#  done < "$PROFILE_FILE"
#
#  if [ "$FOUND" = false ]; then
#    echo "$FOLDER_SAFIRA" >> "$PROFILE_FILE"
#  fi
#}

# installFile verifies the MD5 for the file, then unpacks and
# installs it.
installFile() {
  SAFIRA_TMP="$SAFIRA_TMP_ROOT/$BINARY_NAME"
  local sum=$(openssl md5 ${SAFIRA_TMP_FILE} | awk '{print $2}')
  local expected_sum=$(cat ${SAFIRA_SUM_FILE})
  if [ "$sum" != "$expected_sum" ]; then
    echo "MD5 sum of ${SAFIRA_TMP_FILE} does not match. Aborting."
    exit 1
  fi

  mkdir -p "$HOME/.safira/bin/"
  mkdir -p "$SAFIRA_TMP"
  tar xf "$SAFIRA_TMP_FILE" -C "$SAFIRA_TMP"
  SAFIRA_TMP_BIN="$SAFIRA_TMP/safira"
  echo "Preparing to install $BINARY_NAME into ${SAFIRA_INSTALL_DIR}"
  runAsRoot cp "$SAFIRA_TMP_BIN" "$SAFIRA_INSTALL_DIR/$BINARY_NAME"
  addHost
#  addPath
  echo "$BINARY_NAME installed into $SAFIRA_INSTALL_DIR/$BINARY_NAME"
}

# fail_trap is executed if an error occurs.
fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    if [[ -n "$INPUT_ARGUMENTS" ]]; then
      echo "Failed to install $BINARY_NAME with the arguments provided: $INPUT_ARGUMENTS"
    else
      echo "Failed to install $BINARY_NAME"
    fi
    echo -e "\tFor support, go to $REPO_URL"
  fi
  cleanup
  exit $result
}

# testVersion tests the installed client to make sure it is working.
testVersion() {
  set +e
  SAFIRA="$(command -v $BINARY_NAME)"
  if [ "$?" = "1" ]; then
    echo "$BINARY_NAME not found. Is $SAFIRA_INSTALL_DIR on your "'$PATH?'
    exit 1
  fi
  set -e
}

cleanup() {
  if [[ -d "${SAFIRA_TMP_ROOT:-}" ]]; then
    rm -rf "$SAFIRA_TMP_ROOT"
  fi
}

# Execution

#Stop execution on any error
trap "fail_trap" EXIT
set -e

# Parsing input arguments (if any)
export INPUT_ARGUMENTS="${@}"
set -u
while [[ $# -gt 0 ]]; do
  case $1 in
    '--version'|-v)
       shift
       if [[ $# -ne 0 ]]; then
           export DESIRED_VERSION="${1}"
       else
           echo -e "Please provide the desired version. e.g. --version v3.0.0 or -v canary"
           exit 0
       fi
       ;;
    '--no-sudo')
       USE_SUDO="false"
       ;;
    '--help'|-h)
       exit 0
       ;;
    *) exit 1
       ;;
  esac
  shift
done
set +u

initArch
initOS
verifySupported
checkTagProvided || checkLatestVersion
if ! checkSafiraInstalledVersion; then
  downloadFile
  installFile
fi
testVersion
cleanup
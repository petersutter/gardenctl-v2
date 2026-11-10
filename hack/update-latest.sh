#!/usr/bin/env bash

# SPDX-FileCopyrightText: Contributors to the Gardener project
#
# SPDX-License-Identifier: Apache-2.0

# Usage: update-latest.sh <released-version> [<latest-file>]
#
# Points the LATEST file (default: LATEST in the repository root) to <released-version>, but only
# if <released-version> is a stable release (vMAJOR.MINOR.PATCH, without pre-release or build
# suffix) that is newer than the version currently referenced by the file. Otherwise the file is
# left untouched.
#
# This allows propagating LATEST to the default branch after a release was created from a hotfix
# branch, without ever moving LATEST backwards (e.g. for a patch release of an older minor version)
# or to a pre-release.

set -o errexit
set -o nounset
set -o pipefail

if [[ $# -lt 1 || $# -gt 2 ]]; then
  echo "usage: $0 <released-version> [<latest-file>]" >&2
  exit 2
fi

released_version="$1"
latest_file="${2:-"$(dirname "$0")/../LATEST"}"

stable_version_regex='^v([0-9]+)\.([0-9]+)\.([0-9]+)$'

if [[ ! "${released_version}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+([-+].*)?$ ]]; then
  echo "error: released version '${released_version}' is not a valid version" >&2
  exit 1
fi

if [[ ! "${released_version}" =~ ${stable_version_regex} ]]; then
  echo "${released_version} is not a stable release, leaving ${latest_file} untouched"
  exit 0
fi

released=("${BASH_REMATCH[1]}" "${BASH_REMATCH[2]}" "${BASH_REMATCH[3]}")

current_version="$(tr -d '[:space:]' < "${latest_file}")"

if [[ ! "${current_version}" =~ ${stable_version_regex} ]]; then
  echo "error: ${latest_file} contains '${current_version}', which is not a stable version" >&2
  exit 1
fi

current=("${BASH_REMATCH[1]}" "${BASH_REMATCH[2]}" "${BASH_REMATCH[3]}")

for i in 0 1 2; do
  # force base 10, so that components are never interpreted as octal numbers
  if (( 10#${released[i]} > 10#${current[i]} )); then
    printf '%s\n' "${released_version}" > "${latest_file}"
    echo "updated ${latest_file} from ${current_version} to ${released_version}"
    exit 0
  fi

  if (( 10#${released[i]} < 10#${current[i]} )); then
    break
  fi
done

echo "${latest_file} already references ${current_version} (>= ${released_version}), leaving it untouched"

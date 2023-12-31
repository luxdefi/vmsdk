#!/usr/bin/env bash
# Copyright (C) 2023-2024, Lux Partners Limited. All rights reserved.
# See the file LICENSE for licensing terms.

set -e


if ! [[ "$0" =~ scripts/mock.gen.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

VMSDK_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  cd .. && pwd
)
source "$VMSDK_PATH"/scripts/constants.sh

if ! command -v mockgen &> /dev/null
then
  echo "mockgen not found, installing..."
  # https://github.com/golang/mock
  go install -v github.com/golang/mock/mockgen@v1.6.0
fi

# tuples of (source interface import path, comma-separated interface names, output file path)
input="scripts/mocks.mockgen.txt"
while IFS= read -r line
do
  IFS='=' read src_import_path interface_name output_path <<< "${line}"
  package_name=$(basename $(dirname $output_path))
  echo "Generating ${output_path}..."
  mockgen -copyright_file=./LICENSE.header -package=${package_name} -destination=${output_path} ${src_import_path} ${interface_name}
done < "$input"

echo "SUCCESS"

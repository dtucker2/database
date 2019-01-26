#!/bin/bash

cd "$(dirname $0)"
DIRS=". query"
set -e
for subdir in $DIRS; do
  pushd $subdir
  go vet
  popd
done

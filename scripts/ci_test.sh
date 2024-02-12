#!/bin/sh
set -e


CI=${CI:-""}
CHECKSUM_FILE=./checksum
CHECKSUM_NEW=$(find . -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat $CHECKSUM_FILE || true)

if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  go test --count=1 -cover -coverprofile cover.out  \
    ./configuration/... \
    ./database/... \
    ./datastore/... \
    ./gateway/... \
    ./infrastructure/... \
    ./internal/... \
    ./services/...
  find . -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > ./checksum
      
  # to view coverage percentage on default browser, uncomment the line bellow
  if [ "$CI" = "" ];
  then
    go tool cover -html=cover.out
  fi

fi

#!/bin/sh
set -e

OPENAPI_DIR=openapi
CHECKSUM_FILE="$OPENAPI_DIR/checksum"
SCANNING_DIR=services

PORTAL_DIR=$SCANNING_DIR/portal/entrypoint/api
SDK_DIR=$SCANNING_DIR/sdk/entrypoint/api

CHECKSUM_NEW=$(find $SCANNING_DIR -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1)
CHECKSUM_OLD=$(cat $CHECKSUM_FILE || true)
if [ "$CHECKSUM_NEW" != "$CHECKSUM_OLD" ];
then
  echo "generating portal ...";
  rm -rf "$OPENAPI_DIR/Portal_*"
  swag init -q --instanceName Portal -d $PORTAL_DIR -g openapi.go -o $OPENAPI_DIR --parseDependency --parseInternal --requiredByDefault;

  echo "generating sdk ...";
  rm -rf "$OPENAPI_DIR/Sdk_*"
  swag init -q --instanceName Sdk -d $SDK_DIR -g openapi.go -o $OPENAPI_DIR --parseDependency --parseInternal --requiredByDefault;

  echo "generating checksum ...";
  find $SCANNING_DIR -type f -name '*.go' -exec sha256sum {} \; | sort -k 2 | sha256sum | cut -d  ' ' -f1 > $CHECKSUM_FILE;
fi
 
#!/bin/sh
set -e

SHOULD_BUILD=${SHOULD_BUILD:-""}
MODE=${MODE:-"START"}

if [ "$MODE" != "START" ];
then
  docker compose -f docker-compose.yaml -f docker-compose.localdev.yaml stop storage dispatcher scheduler sdk portal startup
  exit 0
fi

make generator migrate-up

if [ "$SHOULD_BUILD" != "" ];
then
  docker compose -f docker-compose.localdev.yaml build 
fi

docker compose -f docker-compose.yaml -f docker-compose.localdev.yaml up -d 


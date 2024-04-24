#!/bin/sh
set -e

SHOULD_BUILD=${SHOULD_BUILD:-""}
SHOULD_MONITOR=${SHOULD_MONITOR:-""}
MODE=${MODE:-"START"}

if [ "$MODE" != "START" ];
then
  docker compose -f docker-compose.yaml -f docker-compose.localdev.yaml stop storage dispatcher scheduler sdk portal startup
  docker compose -f docker-compose.yaml -f docker-compose.localdev.yaml rm -f storage dispatcher scheduler sdk portal startup
  
  docker compose -f docker-compose.monitoring.yaml stop otelcol 
  docker compose -f docker-compose.monitoring.yaml rm -f otelcol 
  exit 0
fi

make generator migrate-up

if [ "$SHOULD_BUILD" != "" ];
then
  docker compose -f docker-compose.yaml -f docker-compose.monitoring.yaml up -d otelcol 
fi

if [ "$SHOULD_BUILD" != "" ];
then
  docker compose -f docker-compose.localdev.yaml build 
fi

docker compose -f docker-compose.yaml -f docker-compose.localdev.yaml up -d 


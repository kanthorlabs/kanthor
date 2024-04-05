.PHONY: openapi

generator: openapi ioc
migrate-up: up-db up-ds
migrate-down: down-db down-ds

ioc:
	./scripts/gen_ioc.sh

openapi:
	./scripts/gen_openapi.sh

up-db:
	go run cmd/data/main.go migrate up -s file://data/migration/database/postgres -d "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable&x-migrations-table=kanthor_db_migration"

up-ds:
	go run cmd/data/main.go migrate up -s file://data/migration/datastore/postgres -d "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable&x-migrations-table=kanthor_ds_migration"

down-db:
	go run cmd/data/main.go migrate down -s file://data/migration/database/postgres -d "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable&x-migrations-table=kanthor_db_migration"

down-ds:
	go run cmd/data/main.go migrate down -s file://data/migration/datastore/postgres -d "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable&x-migrations-table=kanthor_ds_migration"


gk-encode:
	cat data/gatekeeper/definitions.json | jq -c . | base64 -w 0

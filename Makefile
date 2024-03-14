all: swagger ioc

ioc:
	./scripts/gen_ioc.sh

swagger:
	./scripts/gen_swagger.sh

db-up:
	migrate -source file://data/migration/database/postgres -database "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable&x-migrations-table=kanthor_db_migration" up 1

db-down:
	migrate -source file://data/migration/database/postgres -database "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable&x-migrations-table=kanthor_db_migration" down 1

ds-up:
	migrate -source file://data/migration/datastore/postgres -database "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable&x-migrations-table=kanthor_ds_migration" up 1

ds-down:
	migrate -source file://data/migration/datastore/postgres -database "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable&x-migrations-table=kanthor_ds_migration" down 1

gk-encode:
	cat data/gatekeeper/definitions.json | jq -c . | base64 -w 0

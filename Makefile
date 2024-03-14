all: swagger ioc

ioc:
	./scripts/gen_ioc.sh

swagger:
	./scripts/gen_swagger.sh

mdb-up:
	migrate -source file://data/migration/database/postgres -database "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable" up 1

mdb-down:
	migrate -source file://data/migration/database/postgres -database "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable" down 1

mds-up:
	migrate -source file://data/migration/datastore/postgres -database "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable" up 1

mds-down:
	migrate -source file://data/migration/datastore/postgres -database "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable" down 1

gk-encode:
	cat data/gatekeeper/definitions.json | jq -c . | base64 -w 0

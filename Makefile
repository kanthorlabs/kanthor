all: swagger ioc

ioc:
	./scripts/gen_ioc.sh

swagger:
	./scripts/gen_swagger.sh

mpg-up:
	migrate -source file://data/migration/database/postgres -database "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable" up 1

mpg-down:
	migrate -source file://data/migration/database/postgres -database "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable" down 1

gk-encode:
	cat data/gatekeeper/definitions.json | jq -c . | base64 -w 0
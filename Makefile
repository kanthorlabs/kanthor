all: ioc swagger

ioc:
	./scripts/gen_ioc.sh

swagger:
	./scripts/gen_swagger.sh

migration-postgres-up:
	migrate -source file://data/migration/database/postgres -database "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable" up 1
migration-postgres-down:
	migrate -source file://data/migration/database/postgres -database "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable" down 1
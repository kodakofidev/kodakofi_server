include .env

# DB_SOURCE="postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:6543/${DB_NAME}?search_path=public"
DB_SOURCE=${DB_URL}
MIGRATIONS_DIR=./migration

# make migrate-init name="table_name"
migrate-init:
	migrate create -dir ${MIGRATIONS_DIR} -ext sql ${name}

# make migrate-up
migrate-up:
	migrate -path ${MIGRATIONS_DIR} -database ${DB_SOURCE} -verbose up

# make migrate-up VERSION=fileNumber
# command for migrate up but only 1 file up
# fileNumber asign with number type

# make migrate-down
migrate-down:
	migrate -path ${MIGRATIONS_DIR} -database ${DB_SOURCE} -verbose down

# make migrate-fix
migrate-fix:
	migrate -path ${MIGRATIONS_DIR} -database ${DB_SOURCE} force 0

# make run
run:
	go run ./cmd/main.go

# make seed
seed:
	go run ./cmd/seeder/seed.main.go

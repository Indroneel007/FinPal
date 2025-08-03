postgres:
	docker run --name some-postgres -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18beta1-alpine

redis:
	docker run -d --name redis -p 6379:6379 redis:8.0

createdb:
	docker exec -it some-postgres createdb --username=root --owner=root simple_bank_2

dropdb:
	docker exec -it some-postgres dropdb simple_bank_2

migrateup:
	migrate -source file://db/migration -database "postgresql://finpal_postgres_user:wMJTxyATm6dtr2NGq29Vm7Eala082iEZ@dpg-d27efo6uk2gs73e30sh0-a/finpal_postgres" -verbose up

migratedown:
	migrate -source file://db/migration -database "postgresql://finpal_postgres_user:wMJTxyATm6dtr2NGq29Vm7Eala082iEZ@dpg-d27efo6uk2gs73e30sh0-a/finpal_postgres" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

tidy:
	go mod tidy

build:
	go build -o bin/app.exe ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test tidy build server redis
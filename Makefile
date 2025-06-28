postgres:
	docker run --name some-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18beta1-alpine

createdb:
	docker exec -it some-postgres createdb --username=root --owner=root simple_bank_2

dropdb:
	docker exec -it some-postgres dropdb simple_bank_2

migrateup:
	migrate -source file://db/migration -database "postgresql://root:secret@localhost:5432/simple_bank_2?sslmode=disable" -verbose up

migratedown:
	migrate -source file://db/migration -database "postgresql://root:secret@localhost:5432/simple_bank_2?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc
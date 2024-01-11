postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=root123 -e POSTGRES_USER=root -d postgres:16-alpine
createDb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank
dropDb:
	docker exec -it postgres dropdb simple_bank

migrateUp:
	migrate -path db/migration -database "postgres://root:root123@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateDown:
	migrate -path db/migration -database "postgres://root:root123@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

server:
	go run .

test:
	go test -v -cover ./...


.PHONY:createDb postgres dropDb migrateUp migrateDown sqlc test server

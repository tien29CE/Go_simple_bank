postgres:
	sudo docker run --name postgres18.3 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:18.3-alpine

createdb:
	sudo docker exec -it postgres18.3 createdb --username=root --owner=root simple_bank

dropdb:
	sudo docker exec -it postgres18.3 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createdb dropdb postgres migrateup migratedown sqlc server

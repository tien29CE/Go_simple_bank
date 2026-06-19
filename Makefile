DB_URL=postgresql://root:root@localhost:5432/simple_bank?sslmode=disable

postgres:
	sudo docker run --name postgres18.3 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:18.3-alpine

createdb:
	sudo docker exec -it postgres18.3 createdb --username=root --owner=root simple_bank

dropdb:
	sudo docker exec -it postgres18.3 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/tien29CE/Go_simple_bank.git/db/sqlc Store

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --experimental_allow_proto3_optional --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

redis:
	sudo docker run --name redis8 --network bank-network -p 6379:6379 -d redis:8-alpine

.PHONY: createdb dropdb postgres migrateup migratedown migrateup1 migratedown1 sqlc server proto evans redis

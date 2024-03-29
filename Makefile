postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simplebank

dropdb:
	docker exec -it postgres16 dropdb simplebank

migrateup:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simplebank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simplebank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simplebank?sslmode=disable" -verbose down 1

sqlc:
	docker run --rm -v ".:/src" -w /src sqlc/sqlc generate
	
server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/viennn/docker-postgres-go/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server mock
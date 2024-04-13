DB_SOURCE=postgresql://root:wIDf9SfNRYQbkQjd5Gof@simplebank.cbtafi1aarj8.us-east-1.rds.amazonaws.com:5432/simplebank?sslmode=disable
DOCKER_NETWORK=bank-network

runpostgres:
	docker run --name postgre1 --network $(DOCKER_NETWORK) -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecret -p 5432:5432 -d postgres:14.5-alpine

createdb:
	docker exec -it postgre1 createdb --user=root --owner=root simplebank

dropdb:
	docker exec -it postgre1 dropdb simplebank

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

#Migrate up last migration
migrateup1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

#Migrate down last migration
migratedown1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run .

mock:
	mockgen -destination db/mock/store.go -package mockdb simplebank/db/sqlc Store

.PHONY: runpostgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock
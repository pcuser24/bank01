DOCKER_NETWORK=bank-network

runpostgres:
	docker run --name postgre1 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecret -p 5432:5432 -d postgres:14.5-alpine

createdb:
	docker exec -it postgre1 createdb --user=root --owner=root simplebank

dropdb:
	docker exec -it postgre1 dropdb simplebank

migrateup:
	migrate -path db/migration -database "${DB_SOURCE}" -verbose up

#Migrate up last migration
migrateup1:
	migrate -path db/migration -database "${DB_SOURCE}" -verbose up 1

migratedown:
	migrate -path db/migration -database "${DB_SOURCE}" -verbose down

#Migrate down last migration
migratedown1:
	migrate -path db/migration -database "${DB_SOURCE}" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run .

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/user2410/simplebank/db/sqlc Store; \
	mockgen -package storage -destination storage/storage_mock.go github.com/user2410/simplebank/storage Storage; 

# Localstack: for development and testing only
# Start S3
BUCKET_NAME=simplebank
ls_s3:
	localstack start -d && \
	sleep 10 && \
	awslocal s3api create-bucket --bucket $(BUCKET_NAME) --region ap-southeast-1 --create-bucket-configuration LocationConstraint=ap-southeast-1

.PHONY: runpostgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock

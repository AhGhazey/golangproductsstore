postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secert -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root inventory

dropdb:
	docker exec -it postgres12 dropdb inventory

migrateup:
	migrate -path ./pkg/storage/postgres/migration -database "postgresql://root:secert@localhost:5432/inventory?sslmode=disable" -verbose up

migratedown:
	migrate -path ./pkg/storage/postgres/migration -database "postgresql://root:secert@localhost:5432/inventory?sslmode=disable" -verbose down

buildserver:
	go build cmd/ims.server/main.go

runserver:
	go run cmd/ims.server/main.go

.PHONY: postgres createdb dropdb migrateup migratedown buildserver runserver
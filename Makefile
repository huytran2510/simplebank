postgres:
	docker run -d -p 5432:5432 --name my_postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 postgres:15.7
createdb: 
	docker exec -it my_postgres createdb -U root -O root simple_bank
dropdb: 
	docker exec -it my_postgres dropdb -U root simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	SQLC_GENERATE_CMD := $(shell yq r sqlc.yml sqlc)
	$(SQLC_GENERATE_CMD) 
.PHONY:
	postgres createdb dropdb migrateup migratedown

test:
	go test -v -cover ./...

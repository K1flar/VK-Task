PG_USER=postgres
PG_PASSWORD=postgres
PG_DB=film_library	

build:
	docker-compose build server

run: 
	docker-compose up server

stop:
	docker-compose down

swag:
	swag init -g cmd/main.go

test:
	go test -v --count=1 ./cmd/... ./internal/... ./pkg/...

cover:
	go test -v --count=1 ./cmd/... ./internal/... ./pkg/... -coverprofile=coverage
	go tool cover -html=coverage

db-start:
	@echo "Starting the database..."
	docker-compose up db -d 

db-stop:
	docker stop vktask-db-1

migrate-down:
	@docker exec -it vktask-db-1 psql -U $(PG_USER) -d $(PG_DB) -f /migrations/migrate.down.sql

migrate-up:
	@docker exec -it vktask-db-1 psql -U $(PG_USER) -d $(PG_DB) -f /migrations/migrate.up.sql  

migrate-reset: migrate-down migrate-up

.PHONY: testdata
testdata:
	@docker exec -it vktask-db-1 psql -U $(PG_USER) -d $(PG_DB) -f /testdata/testdata.sql  
	
clean:
	rm -rf .database
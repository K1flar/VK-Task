PG_USER=postgres
PG_PASSWORD=postgres
PG_DB=film_library	

swag:
	swag init -g cmd/main.go

test:
	go test -v --count=1 ./...

cover:
	go test -v --count=1 ./cmd/... ./internal/... ./pkg/... -coverprofile=coverage
	go tool cover -html=coverage

db-start:
	@echo "Starting the database..."
	@mkdir -p ./store/postgres
	docker run --name postgres --rm -d -p 5432:5432 \
		-v ./testdata:/testdata \
		-v ./store/postgres:/var/lib/postgresql/data \
		-v ./migrations:/migrations \
		-e POSTGRES_USER=$(PG_USER) \
		-e POSTGRES_PASSWORD=$(PG_PASSWORD) \
		-e POSTGRES_DB=$(PG_DB) \
		postgres

db-stop:
	docker stop postgres

migrate-down:
	@docker exec -it postgres psql -U $(PG_USER) -d $(PG_DB) -f /migrations/migrate.down.sql

migrate-up:
	@docker exec -it postgres psql -U $(PG_USER) -d $(PG_DB) -f /migrations/migrate.up.sql  

migrate-reset: migrate-down migrate-up

.PHONY: testdata
testdata:
	@docker exec -it postgres psql -U $(PG_USER) -d $(PG_DB) -f /testdata/testdata.sql  
	
clean:
	rm -rf coverage
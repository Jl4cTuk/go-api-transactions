GOOSE_DRIVER := postgres
GOOSE_DBSTRING := postgres://myuser:mypassword@localhost:5432/mydatabase

run:
	@go run -C . ./cmd/qual

build:
	@go build -C . ./cmd/qual

goose-status:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=migrations status

goose-up:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=migrations up

goose-down:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=migrations down

psql-up:
	@docker compose up -d --build

psql-down:
	@docker compose down

psql-enter:
	@docker exec -it postgres_db psql -U myuser -d mydatabase
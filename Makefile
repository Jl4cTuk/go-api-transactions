GOOSE_DRIVER := postgres
GOOSE_DBSTRING := postgres://myuser:mypassword@localhost:5432/mydatabase

CONFIG_PATH_DEV := ./config/dev.yml
CONFIG_PATH_PROD := ./config/prod.yml
CONFIG_PATH_DEV_DOCKER := ./config/dev.yml

run-dev:
	@CONFIG_PATH=$(CONFIG_PATH_DEV) go run -C . ./cmd/qual

run-prod:
	@CONFIG_PATH=$(CONFIG_PATH_PROD) go run -C . ./cmd/qual

build:
	@CONFIG_PATH=$(CONFIG_PATH) go build -C . ./cmd/qual

goose-status:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=migrations status

goose-up:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=migrations up

goose-down:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=migrations down

psql-up:
	@docker compose -f docker-compose.db.yml up -d --build

psql-down:
	@docker compose -f docker-compose.db.yml down

psql-enter:
	@docker exec -it postgres_db psql -U myuser -d mydatabase

docker-run:
	@docker run --rm --network infotex_infotex -e CONFIG_PATH=$(CONFIG_PATH_DEV_DOCKER) -p 8082:8080 qual:latest

docker-build:
	@docker build -t qual:latest .
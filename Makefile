.PHONY: docker-build

APP_NAME := graphqlshop
REPOSITORY := fanyshu
VERSION := $(if $(TAG),$(TAG),$(if $(BRANCH_NAME),$(BRANCH_NAME),$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)))
NOCACHE := $(if $(NOCACHE),"--no-cache")
PSQL_DEV_DSN := "host=localhost port=5432 dbname=shop_db user=shop_db password=shop_db sslmode=disable connect_timeout=5 binary_parameters=yes"

docker-build:
	@docker build ${NOCACHE} --pull -f ./Dockerfile -t ${REPOSITORY}/${APP_NAME}:${VERSION} --ssh default .

migration-up:
	@goose -dir ./migrations "postgres" $(PSQL_DEV_DSN) up
migration-down:
	@goose -dir ./migrations "postgres" $(PSQL_DEV_DSN) down

gen-install:
	@go get github.com/99designs/gqlgen

gen:
	@gqlgen generate

mocks:
	@echo "Regenerate mocks..."
	@rm -rf ./test/mocks/
	@go generate ./...

psql-exec:
	PGPASSWORD=shop_db docker-compose exec psql \
	psql -h localhost -p 5432 -U shop_db -d shop_db

run:
	@docker-compose up -d psql redis-cluster
	@sleep 5
	@docker-compose up -d --build app

stop:
	@docker-compose down
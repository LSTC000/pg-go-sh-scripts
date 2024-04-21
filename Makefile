GO=go
GO_TEST=$(GO) test
GO_GENERATE=$(GO) generate

MOCK_GEN=mockgen

APP_CMD_NAME=main.go
APP_CMD_PATH=cmd/app
APP_BUILD_NAME=main
APP_BUILD_PATH=build/app

MIGRATION_DIR=migration
MIGRATION_PG_CRED="user=superuser password=qwerty host=0.0.0.0 port=5437 dbname=scripts sslmode=disable"

DC=docker compose
DC_LOCAL=$(DC) -f ./deploy/docker-compose.local.yaml
DC_PROD=$(DC) -f ./deploy/docker-compose.prod.yaml


# --================ App ================--
.PHONY: app-run
app-run:
	$(GO) run ./$(APP_CMD_PATH)/$(APP_CMD_NAME)

.PHONY: app-build
app-build:
	$(GO) build -C ./$(APP_CMD_PATH) -o ../../$(APP_BUILD_PATH)/$(APP_BUILD_NAME)

.PHONY: app-start
app-start:
	./$(APP_BUILD_PATH)/$(APP_BUILD_NAME)


# --================ Code Style ================--
.PHONY: import-run
import-run:
	goimports -w -l ./..

.PHONY: fmt-run
fmt-run:
	gofmt -w -l ./..

.PHONY: lint-run
lint-run:
	golangci-lint run


# --================ Swagger ================--
.PHONY: swag-gen
swag-gen:
	swag init -g ./cmd/app/main.go


# --================ Generate ================--
.PHONY: gen-run
gen-run:
	$(GO_GENERATE) ./...


# --================ Test ================--
.PHONY: test-run
test-run:
	$(GO_TEST) -v -count=1 ./...

.PHONY: test-cover
test-cover:
	$(GO_TEST) -short -count=1 -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out
	rm coverage.out


# --================ Migration Goose ================--
.PHONY: goose-create
goose-create:
	goose -dir $(MIGRATION_DIR) create $(MIGRATION_NAME) sql

.PHONY: goose-pg-up
goose-pg-up:
	goose -dir $(MIGRATION_DIR) postgres $(MIGRATION_PG_CRED) up

.PHONY: goose-pg-down
goose-pg-down:
	goose -dir $(MIGRATION_DIR) postgres $(MIGRATION_PG_CRED) down


# --================ Docker Local ================--
.PHONY: docker-local
docker-local:
	$(DC_LOCAL) up -d --build

.PHONY: docker-local-build
docker-local-build:
	$(DC_LOCAL) build

.PHONY: docker-local-up
docker-local-up:
	$(DC_LOCAL) up -d

.PHONY: docker-local-stop
docker-local-stop:
	$(DC_LOCAL) stop

.PHONY: docker-local-start
docker-local-start:
	$(DC_LOCAL) start

.PHONY: docker-local-down
docker-local-down:
	$(DC_LOCAL) down

.PHONY: docker-local-down-v
docker-local-down-v:
	$(DC_LOCAL) down -v

.PHONY: docker-local-logs
docker-local-logs:
	$(DC_LOCAL) logs

.PHONY: docker-local-logs-f
docker-local-logs-f:
	$(DC_LOCAL) logs -f


# --================ Docker Prod ================--
.PHONY: docker-prod
docker-prod:
	$(DC_PROD) up -d --build

.PHONY: docker-prod-build
docker-prod-build:
	$(DC_PROD) build

.PHONY: docker-prod-up
docker-prod-up:
	$(DC_PROD) up -d

.PHONY: docker-prod-stop
docker-prod-stop:
	$(DC_PROD) stop

.PHONY: docker-prod-start
docker-prod-start:
	$(DC_PROD) start

.PHONY: docker-prod-down
docker-prod-down:
	$(DC_PROD) down

.PHONY: docker-prod-down-v
docker-prod-down-v:
	$(DC_PROD) down -v

.PHONY: docker-prod-logs
docker-prod-logs:
	$(DC_PROD) logs

.PHONY: docker-prod-logs-f
docker-prod-logs-f:
	$(DC_PROD) logs -f

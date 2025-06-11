DEP_COMPOSE_FILES = \
	-f deploy/compose/rdb.yml \
	-f deploy/compose/redis.yml \
	-f deploy/compose/vault-dev.yml

.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fix-lint
fix-lint:
	golangci-lint run --fix --timeout=10m

.PHONY: run-api
run-api:
	@sed -i 's#^VAULT_RDB_ROLE=.*#VAULT_RDB_ROLE=creds/postgres-app-role#' .env
	@SERVICE=api ./scripts/export-vaul-env.sh
	go run cmd/api/main.go

.PHONY: migrate-up
migrate-up:
	@sed -i 's#^VAULT_RDB_ROLE=.*#VAULT_RDB_ROLE=static-creds/postgres-migrate-role#' .env
	@SERVICE=migration ./scripts/export-vaul-env.sh
	@go run cmd/migrate/main.go --type up

.PHONY: migrate-down
migrate-down:
	@sed -i 's#^VAULT_RDB_ROLE=.*#VAULT_RDB_ROLE=static-creds/postgres-migrate-role#' .env
	@SERVICE=migration ./scripts/export-vaul-env.sh
	@go run cmd/migrate/main.go --type down

.PHONY: migrate-step
migrate-step:
	@sed -i 's#^VAULT_RDB_ROLE=.*#VAULT_RDB_ROLE=static-creds/postgres-migrate-role#' .env
	@SERVICE=migration ./scripts/export-vaul-env.sh
	@read -p "Module name: " module; \
	read -p "Step (integer): " step; \
	go run cmd/migrate/main.go --type step --module $$module --step $$step

.PHONY: gen-migration
gen-migration:
	@read -p "Module: " module; \
	read -p "Description: " desc; \
	migrate create -ext sql -digits 14 -dir ./migrations/$$module $$desc

.PHONY: run-deps
run-deps:
	docker compose $(DEP_COMPOSE_FILES) up -d
	
.PHONY: down-deps
down-deps:
	docker compose $(DEP_COMPOSE_FILES) down --volumes

.PHONY: run-kc
run-kc:
	docker compose -f deploy/compose/keycloak.yml up -d

.PHONY: stop-kc
stop-kc:
	docker compose -f deploy/compose/keycloak.yml stop

.PHONY: run-mdb
run-mdb:
	docker compose -f deploy/compose/mongo.yml up -d

.PHONY: down-mdb
down-mdb:
	docker compose -f deploy/compose/mongo.yml down --volumes

.PHONY: export-realm
export-realm:
	scripts/export-realm.sh

.PHONY: run-monitoring
run-monitoring:
	docker compose -f deploy/compose/monitoring.yml up -d

.PHONY: down-monitoring
down-monitoring:
	docker compose -f deploy/compose/monitoring.yml down --volumes

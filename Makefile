.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fix-lint
fix-lint:
	golangci-lint run --fix --timeout=10m

.PHONY: migrate-up
migrate-up:
	go run cmd/migrate/main.go --type up

.PHONY: migrate-down
migrate-down:
	go run cmd/migrate/main.go --type down

.PHONY: migrate-step
migrate-step:
	@read -p "Module name: " module; \
	read -p "Step (integer): " step; \
	go run cmd/migrate/main.go --type step --module $$module --step $$step

.PHONY: gen-migration
gen-migration:
	@read -p "Module: " module; \
	read -p "Description: " desc; \
	migrate create -ext sql -digits 14 -dir ./migrations/$$module $$desc

.PHONY: run-rdb
run-rdb:
	docker compose -f deploy/compose/rdb.yaml up -d

.PHONY: stop-rdb
stop-rdb:
	docker compose -f deploy/compose/rdb.yaml stop

.PHONY: clean-rdb
clean-rdb:
	docker compose -f deploy/compose/rdb.yaml down --volumes

.PHONY: run-kc
run-kc:
	docker compose -f deploy/compose/keycloak.yaml up -d

.PHONY: stop-kc
stop-kc:
	docker compose -f deploy/compose/keycloak.yaml stop

.PHONY: run-mdb
run-mdb:
	docker compose -f deploy/compose/mongo.yaml up -d

.PHONY: stop-mdb
stop-mdb:
	docker compose -f deploy/compose/mongo.yaml stop

.PHONY: clean-mdb
clean-mdb:
	docker compose -f deploy/compose/mongo.yaml down --volumes

.PHONY: run-pgadmin
run-pgadmin:
	docker compose -f deploy/compose/pgadmin.yaml up -d

.PHONY: stop-pgadmin
stop-pgadmin:
	docker compose -f deploy/compose/pgadmin.yaml down

.PHONY: export-realm
export-realm:
	scripts/export-realm.sh

.PHONY: run-vault
run-vault:
	docker compose -f deploy/compose/vault-dev.yaml up -d

.PHONY: stop-vault
stop-vault:
	docker compose -f deploy/compose/vault-dev.yaml stop

.PHONY: down-vault
down-vault:
	docker compose -f deploy/compose/vault-dev.yaml down --volumes

.PHONY: export-vault-env
export-vault-env:
	./scripts/export-vaul-env.sh

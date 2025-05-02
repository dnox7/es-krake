.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: run-migration
run-migration:
	go run cmd/migrate/main.go

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

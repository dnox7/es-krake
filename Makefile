.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: run-kc
run-kc:
	docker compose -f deploy/docker-compose.yaml up -d keycloak

.PHONY: stop-kc
stop-kc:
	docker compose -f deploy/docker-compose.yaml stop keycloak

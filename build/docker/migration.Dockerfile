ARG BASE 
FROM ${BASE}

WORKDIR ${APP_PATH}

COPY ./config ./config
COPY ./cmd/migrate ./cmd/migrate
COPY ./internal/domain/shared ./internal/domain/shared
COPY ./internal/infrastructure/rdb ./internal/infrastructure/rdb
COPY ./pkg ./pkg

RUN --mount=type=cache,target=/gocache \
    go mod vendor && \
    go build \
    -mod=vendor \
    -tags netgo \
    -ldflags="-linkmode external -extldflags \"-static\"" \
    -v \
    -o /migration-tool \
    ./cmd/migrate/main.go

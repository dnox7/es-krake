FROM golang:1.24.3-alpine3.21

ENV GOPATH=/go \
    GOCACHE=/gocache \
    GOMODCACHE=/gocache \
    CGO_ENABLED=1

ENV PKG_NAME=es-krake 
ENV APP_PATH=/usr/local/src/${PKG_NAME}

RUN apk --no-cache add musl-dev gcc make ca-certificates

WORKDIR ${APP_PATH}

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN --mount=type=cache,target=/gocache \
    go mod download

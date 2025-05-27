FROM golang:1.24.3-alpine3.21 as builder

ENV GOPATH=/go
ENV GOCACHE=/gocache
ENV GOMODCACHE=/gocache
ENV CGO_ENABLED=1

ENV PKG_NAME=es-krake
ENV APP_PATH=/usr/local/src/${PKG_NAME}

RUN apk --no-cache add musl-dev gcc make ca-certificates

COPY ./go.mod ${APP_PATH}/go.mod
COPY ./go.sum ${APP_PATH}/go.sum

WORKDIR ${APP_PATH}

RUN --mount=type=cache,target=/gocache go mod download

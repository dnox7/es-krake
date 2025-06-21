package pf

import (
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/validator"
	"github.com/graphql-go/graphql"
)

const (
	PostPlatformLogin = "platform/post_login.json"
)

type PlatformHandler struct {
	logger    *log.Logger
	graphql   graphql.Schema
	validator *validator.JsonSchemaValidator
	debug     bool
}

func NewPlatformHandler(schema graphql.Schema, debug bool, validator *validator.JsonSchemaValidator) *PlatformHandler {
	return &PlatformHandler{
		logger:    log.With("handler", "platform_handler"),
		debug:     debug,
		graphql:   schema,
		validator: validator,
	}
}

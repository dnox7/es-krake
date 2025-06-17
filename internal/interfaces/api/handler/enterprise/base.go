package ent

import (
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/validator"
	"github.com/graphql-go/graphql"
)

const (
	PostEnterpriseLogin = "enterprise/post_enterprise_login.json"
)

type EnterpriseHandler struct {
	logger    *log.Logger
	graphql   graphql.Schema
	validator *validator.JsonSchemaValidator
	debug     bool
}

func NewEnterpriseHandler(schema graphql.Schema, debug bool, validator *validator.JsonSchemaValidator) *EnterpriseHandler {
	return &EnterpriseHandler{
		logger:    log.With("handler", "enterprise_handler"),
		debug:     debug,
		graphql:   schema,
		validator: validator,
	}
}

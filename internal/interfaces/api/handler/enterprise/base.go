package ent

import (
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/graphql-go/graphql"
)

const (
	PostLogin = "enterprise/post_login.json"
)

type EnterpriseHandler struct {
	logger  *log.Logger
	debug   bool
	graphql graphql.Schema
}

func NewEnterpriseHandler(schema graphql.Schema, debug bool) *EnterpriseHandler {
	return &EnterpriseHandler{
		logger:  log.With("handler", "enterprise_handler"),
		debug:   debug,
		graphql: schema,
	}
}

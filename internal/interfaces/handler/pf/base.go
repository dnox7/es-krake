package pf

import (
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/graphql-go/graphql"
)

type PlatformHandler struct {
	logger  *log.Logger
	graphql graphql.Schema
}

func NewPlatformHandler(schema graphql.Schema) *PlatformHandler {
	return &PlatformHandler{
		logger:  log.With("handler", "platform_handler"),
		graphql: schema,
	}
}

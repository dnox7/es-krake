package handler

import (
	"github.com/dpe27/es-krake/internal/interfaces/api/handler/pf"
	"github.com/graphql-go/graphql"
)

type HTTPHandler struct {
	PF *pf.PlatformHandler
}

func NewHTTPHandler(schema graphql.Schema, debug bool) HTTPHandler {
	return HTTPHandler{
		PF: pf.NewPlatformHandler(schema, debug),
	}
}

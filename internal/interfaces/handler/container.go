package handler

import (
	"github.com/dpe27/es-krake/internal/interfaces/handler/pf"
	"github.com/graphql-go/graphql"
)

type HTTPHandler struct {
	PF *pf.PlatformHandler
}

func NewHTTPHandler(schema graphql.Schema) HTTPHandler {
	return HTTPHandler{
		PF: pf.NewPlatformHandler(schema),
	}
}

package handler

import (
	ent "github.com/dpe27/es-krake/internal/interfaces/api/handler/enterprise"
	pf "github.com/dpe27/es-krake/internal/interfaces/api/handler/platform"
	"github.com/graphql-go/graphql"
)

type HTTPHandler struct {
	PF  *pf.PlatformHandler
	Ent *ent.EnterpriseHandler
}

func NewHTTPHandler(schema graphql.Schema, debug bool) HTTPHandler {
	return HTTPHandler{
		PF:  pf.NewPlatformHandler(schema, debug),
		Ent: ent.NewEnterpriseHandler(schema, debug),
	}
}

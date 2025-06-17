package handler

import (
	ent "github.com/dpe27/es-krake/internal/interfaces/api/handler/enterprise"
	pf "github.com/dpe27/es-krake/internal/interfaces/api/handler/platform"
	"github.com/dpe27/es-krake/pkg/validator"
	"github.com/graphql-go/graphql"
)

type HTTPHandler struct {
	Pf  *pf.PlatformHandler
	Ent *ent.EnterpriseHandler
}

func NewHTTPHandler(schema graphql.Schema, debug bool, inputValidator *validator.JsonSchemaValidator) HTTPHandler {
	return HTTPHandler{
		Pf:  pf.NewPlatformHandler(schema, debug, inputValidator),
		Ent: ent.NewEnterpriseHandler(schema, debug, inputValidator),
	}
}

package handler

import (
	buyer "github.com/dpe27/es-krake/internal/interfaces/api/handler/buyer"
	ent "github.com/dpe27/es-krake/internal/interfaces/api/handler/enterprise"
	pf "github.com/dpe27/es-krake/internal/interfaces/api/handler/platform"
	"github.com/dpe27/es-krake/pkg/validator"
	"github.com/graphql-go/graphql"
)

type HTTPHandler struct {
	Pf    *pf.PlatformHandler
	Ent   *ent.EnterpriseHandler
	Buyer *buyer.BuyerHandler
}

func NewHTTPHandler(schema graphql.Schema, debug bool, inputValidator *validator.JsonSchemaValidator) HTTPHandler {
	return HTTPHandler{
		Pf:    pf.NewPlatformHandler(schema, debug, inputValidator),
		Ent:   ent.NewEnterpriseHandler(schema, debug, inputValidator),
		Buyer: buyer.NewBuyerHandler(schema, debug, inputValidator),
	}
}

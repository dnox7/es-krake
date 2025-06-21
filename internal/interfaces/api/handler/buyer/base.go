package buyer

import (
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/validator"
	"github.com/graphql-go/graphql"
)

const (
	PostBuyerLogin  = "buyer/post_login.json"
	PostBuyerSignup = "buyer/post_signup.json"
)

type BuyerHandler struct {
	logger    *log.Logger
	graphql   graphql.Schema
	validator *validator.JsonSchemaValidator
	debug     bool
}

func NewBuyerHandler(schema graphql.Schema, debug bool, validator *validator.JsonSchemaValidator) *BuyerHandler {
	return &BuyerHandler{
		logger:    log.With("handler", "buyer_handler"),
		debug:     debug,
		graphql:   schema,
		validator: validator,
	}
}

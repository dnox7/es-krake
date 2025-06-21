package buyer

import (
	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *BuyerHandler) PostBuyerSignup(c *gin.Context) {
	input, err := nethttp.GetInputAsMap(c)
	if err != nil {
		nethttp.SetGenericErrorResponse(c, err, h.debug)
		return
	}

	validationRes, err := h.validator.Validate(PostBuyerSignup, input)
	if err != nil {
		nethttp.SetGenericErrorResponse(c, err, h.debug)
		return
	}

	if validationRes != nil {
		nethttp.SetJSONValidationErrorResponse(c, h.validator, validationRes)
		return
	}

	res := graphql.Do(graphql.Params{
		Context:    c,
		Schema:     h.graphql,
		RootObject: input,
		RequestString: `
			mutation {
				post_signup_buyer() 
			}
		`,
	})

	if res.HasErrors() {
		nethttp.SetGenericErrorResponse(c, res.Errors[0], h.debug)
		return
	}

	nethttp.SetNoContentResponse(c)
}

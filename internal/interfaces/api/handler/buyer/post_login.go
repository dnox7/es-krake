package buyer

import (
	"net/http"

	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *BuyerHandler) PostBuyerLogin(c *gin.Context) {
	input, err := nethttp.GetInputAsMap(c)
	if err != nil {
		nethttp.SetGenericErrorResponse(c, err, h.debug)
		return
	}

	validationRes, err := h.validator.Validate(PostBuyerLogin, input)
	if err != nil {
		nethttp.SetGenericErrorResponse(c, err, h.debug)
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
				post_login_buyer {
					access_token
					refresh_token
					refresh_expires_in
					realm_name
				}
			}
		`,
	})

	if res.HasErrors() {
		nethttp.SetGenericErrorResponse(c, res.Errors[0], h.debug)
		return
	}

	resData := utils.GetSubMap(res.Data, "post_login_buyer")
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(
		resData["realm_name"].(string),
		resData["access_token"].(string),
		resData["refresh_expires_in"].(int),
		"/buyer/auth",
		"",
		true,
		true,
	)

	delete(resData, "realm_name")
	delete(resData, "access_token")
	delete(resData, "refresh_expires_in")

	resData["link"] = map[string]interface{}{
		"refresh": "/buyer/auth/refresh",
		"logout":  "/buyer/auth/logout",
	}
	nethttp.SetOKResponse(c, resData)
}

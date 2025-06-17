package pf

import (
	"net/http"

	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *PlatformHandler) PostPlatformLogin(c *gin.Context) {
	input, err := nethttp.GetInputAsMap(c)
	if err != nil {
		nethttp.SetGenericErrorResponse(c, err, h.debug)
		return
	}

	validationRes, err := h.validator.Validate(PostPlatformLogin, input)
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
				post_login_platform {
					access_token
					refresh_token
					refresh_expires_in
					realm_name
					permissions {
						id
						name
					}
				}
			}
		`,
	})

	if res.HasErrors() {
		nethttp.SetGenericErrorResponse(c, res.Errors[0], h.debug)
		return
	}

	resData := utils.GetSubMap(res.Data, "post_login_platform")
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(
		resData["realm_name"].(string),
		resData["access_token"].(string),
		resData["refresh_expires_in"].(int),
		"/pf/auth",
		"",
		true,
		true,
	)

	delete(resData, "realm_name")
	delete(resData, "access_token")
	delete(resData, "refresh_expires_in")
	
	resData["link"] = map[string]interface{}{
		"refresh": "/pf/auth/refresh",
		"logout":  "/pf/auth/logout",
	}
	nethttp.SetOKReponse(c, resData)
}

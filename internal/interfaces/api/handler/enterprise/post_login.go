package ent

import (
	"net/http"
	"strconv"

	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *EnterpriseHandler) PostLogin(c *gin.Context) {
	enterpriseID := c.Param("enterprise_id")
	_, err := strconv.Atoi(enterpriseID)
	if err != nil {
		nethttp.SetBadRequestResponse(c, map[string]interface{}{
			"enterprise_id": utils.ErrorInputFail,
		}, nil, nil)
		return
	}

	input, err := nethttp.GetInputAsMap(c)
	if err != nil {
		nethttp.SetGenericErrorResponse(c, err, h.debug)
		return
	}

	validationRes, err := h.validator.Validate(PostEnterpriseLogin, input)
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
			mutation ($enterprise_id: AnyInt!) {
				post_login_enterprise(enterprise_id: $enterprise_id) {
					access_token
					refresh_token
					refresh_expires_in
					realm_name
					enterprise_account {
						id
					}
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

	resData := utils.GetSubMap(res.Data, "post_login_enterprise")
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(
		resData["realm_name"].(string),
		resData["refresh_token"].(string),
		resData["refresh_expires_in"].(int),
		"/ent/auth",
		"",
		true,
		true,
	)

	delete(resData, "realm_name")
	delete(resData, "refresh_token")
	delete(resData, "refresh_expires_in")

	nethttp.SetOKResponse(c, resData)
}

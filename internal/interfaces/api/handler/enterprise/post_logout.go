package ent

import (
	"encoding/json"
	"strconv"

	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *EnterpriseHandler) PostLogout(c *gin.Context) {
	enterpriseID, err := strconv.Atoi(c.Param("enterprise_id"))
	if err != nil {
		nethttp.SetBadRequestResponse(c, map[string]interface{}{
			"enterprise_id": utils.ErrorInputFail,
		}, nil, nil)
		return
	}

	kcUserID := c.GetString("kc_user_id")
	if kcUserID == "" {
		nethttp.SetBadRequestResponse(c, map[string]interface{}{
			"access_token": utils.ErrorInputFail,
		}, nil, nil)
		return
	}

	cookiesMap := make(map[string]interface{})
	cookies := c.Request.Cookies()
	for _, cookie := range cookies {
		cookiesMap[cookie.Name] = cookie.Value
	}

	jsonCookies, err := json.Marshal(cookiesMap)
	if err != nil {
		nethttp.SetGenericErrorResponse(c, err, h.debug)
		return
	}

	res := graphql.Do(graphql.Params{
		Context: c,
		Schema:  h.graphql,
		VariableValues: map[string]interface{}{
			"cookies":       string(jsonCookies),
			"enterprise_id": enterpriseID,
			"kc_user_id":    kcUserID,
		},
		RequestString: `
			mutation ($cookies: String!, $enterprise_id: AnyInt!, $kc_user_id: String!) {
				post_logout_enterprise(
					cookies: $cookies
					enterprise_id: $enterprise_id
					kc_user_id: $kc_user_id
				)
			}
		`,
	})

	if res.HasErrors() {
		nethttp.SetGenericErrorResponse(c, res.Errors[0], h.debug)
		return
	}

	nethttp.SetNoContentResponse(c)
}

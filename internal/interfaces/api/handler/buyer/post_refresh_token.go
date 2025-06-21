package buyer

import (
	"encoding/json"
	"net/http"

	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *BuyerHandler) PostRefreshTokenBuyer(c *gin.Context) {
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
			"cookies": string(jsonCookies),
		},
		RequestString: `
			mutation ($cookies: String!) {
				post_refresh_token_buyer(cookies: $cookies) {
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

	resData := utils.GetSubMap(res.Data, "post_refresh_token_buyer")
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(
		resData["realm_name"].(string),
		resData["refresh_token"].(string),
		resData["refresh_expires_in"].(int),
		"/buyer/auth",
		"",
		true,
		true,
	)

	delete(resData, "realm_name")
	delete(resData, "refresh_token")
	delete(resData, "refresh_expires_in")

	resData["link"] = map[string]interface{}{
		"refresh": "/buyer/auth/refresh",
		"logout":  "/buyer/auth/logout",
	}
	nethttp.SetOKResponse(c, resData)
}

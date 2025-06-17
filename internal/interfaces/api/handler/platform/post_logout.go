package pf

import (
	"encoding/json"

	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *PlatformHandler) PostLogout(c *gin.Context) {
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
				post_logout_platform(cookies: $cookies) 
			}
		`,
	})

	if res.HasErrors() {
		nethttp.SetGenericErrorResponse(c, res.Errors[0], h.debug)
		return
	}

	nethttp.SetNoContentResponse(c)
}

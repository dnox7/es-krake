package pf

import (
	"strconv"

	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *PlatformHandler) GetRoleType(c *gin.Context) {
	roleTypeID, err := strconv.Atoi(c.Param("role_type_id"))
	if err != nil {
		nethttp.SetBadRequestResponse(c, map[string]interface{}{
			"role_type_id": "poweirjgforipewjg",
		}, nil, nil)
	}

	res := graphql.Do(graphql.Params{
		Context: c,
		Schema:  h.graphql,
		VariableValues: map[string]interface{}{
			"roleTypeID": roleTypeID,
		},
		RequestString: `
            query ($roleTypeID: AnyInt!) {
                role_type: get_role_type (
                    role_type_id: $roleTypeID
                ) {
                    id
                    name
                    created_at
                    updated_at
                }
            }
        `,
	})

	if res.HasErrors() {
		nethttp.SetGenericErrorResponse(c, res.Errors[0], true)
		return
	}
	nethttp.SetOKReponse(c, res.Data)
}

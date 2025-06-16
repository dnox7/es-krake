package ent

import (
	"strconv"

	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/gin-gonic/gin"
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

	// input, err := nethttp.GetInputAsMap(c)
	// if err != nil {
	// 	nethttp.SetGenericErrorResponse(c, err, h.debug)
	// 	return
	// }
}

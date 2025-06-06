package nethttp

import (
	"net/http"

	"github.com/dpe27/es-krake/pkg/wraperror"
	"github.com/gin-gonic/gin"
)

func GetInputAsMap(c *gin.Context) (map[string]interface{}, error) {
	contentType := c.ContentType()
	if contentType != MIMEApplicationJSON {
		return nil, wraperror.NewAPIError(
			http.StatusBadRequest,
			"Expected 'application/json' for Content-Type, got '"+contentType+"'",
			nil,
		)
	}

	input := make(map[string]interface{})
	err := c.ShouldBindJSON(&input)
	return input, err
}

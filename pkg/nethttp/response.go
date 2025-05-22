package nethttp

import (
	"archive/zip"
	"bytes"
	"net/http"

	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/gin-gonic/gin"
)

type BaseSuccessResponse struct {
	Data interface{} `json:"data"`
}

type BaseErrorResponse struct {
	Error *ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Message   interface{} `json:"message"`
	Details   interface{} `json:"details"`
	DebugInfo interface{} `json:"debug_information"`
}

func NewErrorResponse(message, details, debugInfo interface{}) *BaseErrorResponse {
	return &BaseErrorResponse{
		Error: &ErrorResponse{
			Message:   message,
			Details:   details,
			DebugInfo: debugInfo,
		},
	}
}

func AbortWithBadRequestResponse(c *gin.Context, msg, detail, debugInfo interface{}) {
	c.AbortWithStatusJSON(
		http.StatusBadRequest,
		NewErrorResponse(msg, detail, debugInfo),
	)
}

func AbortWithForbiddenResponse(c *gin.Context, msg, detail, debugInfo interface{}) {
	c.AbortWithStatusJSON(
		http.StatusForbidden,
		NewErrorResponse(msg, detail, debugInfo),
	)
}

func AbortWithUnauthorizedResponse(c *gin.Context, msg, detail, debugInfo interface{}) {
	c.AbortWithStatusJSON(
		http.StatusUnauthorized,
		NewErrorResponse(msg, detail, debugInfo),
	)
}

func AbortWithInternalServerErrorResponse(c *gin.Context, msg, detail, debugInfo interface{}) {
	c.AbortWithStatusJSON(
		http.StatusInternalServerError,
		NewErrorResponse(msg, detail, debugInfo),
	)
}

func AbortWithRequestTimeoutResponse(c *gin.Context, msg, detail, debugInfo interface{}) {
	c.AbortWithStatusJSON(
		http.StatusRequestTimeout,
		NewErrorResponse(msg, detail, debugInfo),
	)
}

func SetOKReponse(c *gin.Context, data interface{}) {
	c.JSON(
		http.StatusOK,
		&BaseSuccessResponse{Data: data},
	)
}

func SetNoContentResponse(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func SetNotFoundResponse(c *gin.Context, msg, detail, debugInfo interface{}) {
	c.JSON(
		http.StatusNotFound,
		NewErrorResponse(msg, debugInfo, debugInfo),
	)
}

func SetBadRequestResponse(c *gin.Context, msg, detail, debugInfo interface{}) {
	c.JSON(
		http.StatusBadRequest,
		NewErrorResponse(msg, detail, debugInfo),
	)
}

func SetInternalServerErrorResponse(c *gin.Context, msg, detail, debugInfo interface{}) {
	c.JSON(
		http.StatusInternalServerError,
		NewErrorResponse(msg, detail, debugInfo),
	)
}

func ResponseCSV(c *gin.Context, statusCode int, fileName string, data []byte) {
	c.Writer.Header().Set(HeaderContentDescription, "File Transfer")
	c.Writer.Header().Set(HeaderContentDisposition, "attachment;filename="+fileName)
	c.Data(statusCode, MIMETextCSV, data)
}

func ResponseZIP(
	c *gin.Context,
	statusCode int,
	fileName string,
	mapContentFile map[string]bytes.Buffer,
) error {
	c.Writer.Header().Set(HeaderContentType, MIMEApplicationZIP)
	c.Writer.Header().Set(HeaderContentDisposition, "attachment;filename="+fileName)

	w := zip.NewWriter(c.Writer)
	defer func() {
		if err := w.Close(); err != nil {
			log.With().Error(c, utils.ErrorCloseWriter, "error", err.Error())
		}
	}()

	for key, contentFile := range mapContentFile {
		f, err := w.Create(key)
		if err != nil {
			return err
		}

		_, err = f.Write(contentFile.Bytes())
		if err != nil {
			return err
		}
	}

	return nil
}

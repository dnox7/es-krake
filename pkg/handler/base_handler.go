package handler

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"pech/es-krake/pkg/dto"
	"pech/es-krake/pkg/utils"
	"pech/es-krake/pkg/validator"
	wraperror "pech/es-krake/pkg/wrap-error"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/xeipuuv/gojsonschema"
)

type BaseHTTPHandler struct {
	Validator *validator.JsonSchemaValidator
}

func (h *BaseHTTPHandler) ResponseCSV(
	c *gin.Context,
	statusCode int,
	fileName string,
	data []byte,
) {
	c.Writer.Header().Set("Content-Description", "File Transfer")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename="+fileName)
	c.Data(statusCode, "text/csv", data)
}

func (h *BaseHTTPHandler) ResponseZIP(
	c *gin.Context,
	statusCode int,
	fileName string,
	mapContentFile map[string]bytes.Buffer,
) error {
	c.Writer.Header().Set("Content-Type", "application/zip")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename"+fileName)

	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	for key, contentFile := range mapContentFile {
		f, err := zipWriter.Create(key)
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

func (h *BaseHTTPHandler) GetInputAsMap(c *gin.Context) (map[string]interface{}, error) {
	contentType := c.ContentType()
	if contentType != "application/json" {
		return nil, wraperror.NewApiDisplayableError(
			http.StatusBadRequest,
			"Expected 'application/json' for Content-Type, got '"+contentType+"'",
			nil,
		)
	}

	input := make(map[string]interface{})
	err := c.ShouldBindJSON(&input)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (h *BaseHTTPHandler) SetGenericErrorResponse(c *gin.Context, finalErr error) {
	originalErr := finalErr
	if _, ok := originalErr.(gqlerrors.FormattedError); ok {
		err := originalErr.(gqlerrors.FormattedError).OriginalError()
		if err != nil {
			originalErr = err
		}

		if _, ok := originalErr.(*gqlerrors.Error); ok {
			err := originalErr.(*gqlerrors.Error).OriginalError
			if err != nil {
				originalErr = err
			}
		}
	}

	apiErr := &wraperror.ApiDisplayableError{}
	jsonErr := &json.SyntaxError{}

	switch {
	case errors.As(originalErr, &apiErr):
		var debugInfo interface{} = nil
		if os.Getenv("PE_DEBUG") == "true" {
			debugInfo = finalErr
		}

		data := dto.BaseErrorResponse{
			Error: &dto.ErrorResponse{
				Message:   apiErr.Message(),
				DebugInfo: debugInfo,
			},
		}
		c.JSON(apiErr.HttpStatus(), data)

	// case errors.Is(originalErr, gorm.ErrRecordNotFound) || originalErr.Error() == gorm.ErrRecordNotFound.Error():
	// 	data := dto.BaseErrorResponse{
	// 		Error: &dto.ErrorResponse{
	// 			Message: originalErr.Error(),
	// 		},
	// 	}
	// 	c.JSON(http.StatusNotFound, data)

	case errors.As(originalErr, &jsonErr):
		data := &dto.BaseErrorResponse{
			Error: &dto.ErrorResponse{
				Message: "Invalid json",
				Details: map[string]interface{}{
					"offset": jsonErr.Offset,
					"error":  jsonErr.Error(),
				},
			},
		}
		c.JSON(http.StatusBadRequest, data)

	default:

	}

	return
}

func (h *BaseHTTPHandler) SetValidationErrorResponse(c *gin.Context, err error) {
	data := &dto.BaseErrorResponse{
		Error: &dto.ErrorResponse{
			Message: err,
		},
	}
	c.JSON(http.StatusBadRequest, data)
}

func (h *BaseHTTPHandler) SetJSONValidationErrorResponse(
	c *gin.Context,
	validationRes *gojsonschema.Result,
) {
	h.SetJSONValidationWithCustomErrorResponse(
		c,
		validationRes,
		func(res gojsonschema.ResultError) string {
			return ""
		},
	)
}

func (h *BaseHTTPHandler) SetJSONValidationWithCustomErrorResponse(
	c *gin.Context,
	validationRes *gojsonschema.Result,
	getErr func(res gojsonschema.ResultError) string,
) {
	messages := map[string]string{}
	details := make([]map[string]interface{}, 0)
	for _, validationErr := range validationRes.Errors() {
		field := h.Validator.GetErrorField(validationErr)
		det := h.Validator.GetErrorDetails(validationErr)

		msg := getErr(validationErr)
		if msg == "" {
			msg = h.Validator.GetCustomErrorMessage(validationErr)
		}

		messages[field] = msg
		details = append(details, det)
	}

	data := &dto.BaseErrorResponse{
		Error: &dto.ErrorResponse{
			Message: messages,
			Details: details,
		},
	}
	c.JSON(http.StatusBadRequest, data)
}

func (h *BaseHTTPHandler) SetBadRequestErrorResponse(c *gin.Context, messages interface{}) {
	data := &dto.BaseErrorResponse{
		Error: &dto.ErrorResponse{
			Message: messages,
		},
	}
	c.JSON(http.StatusBadRequest, data)
}

func (h *BaseHTTPHandler) SetCustomErrorAndDetailResponse(c *gin.Context, err error, details interface{}) {
	data := &dto.BaseErrorResponse{
		Error: &dto.ErrorResponse{
			Message: err.Error(),
			Details: details,
		},
	}
	c.JSON(http.StatusInternalServerError, data)
}

func (h *BaseHTTPHandler) SetInternalErrorResponse(c *gin.Context, err error) {
	var debugInfo interface{}
	if os.Getenv("PE_DEBUG") == "true" {
		debugInfo = err.Error()
	} else {
		slog.ErrorContext(c, "Unexpected error", "error", err)

		data := &dto.BaseErrorResponse{
			Error: &dto.ErrorResponse{
				Message:   utils.ErrInternalServerMsg,
				DebugInfo: debugInfo,
			},
		}
		c.JSON(http.StatusInternalServerError, data)
	}
}

func (h *BaseHTTPHandler) SetCookie(c *gin.Context, cookieRes []map[string]interface{}, inputMaxAge *int) error {
	for _, cookie := range cookieRes {
		name := cookie["name"].(string)
		value := cookie["value"].(string)
		path := cookie["path"].(string)
		domain := cookie["domain"].(string)
		secure := cookie["secure"].(bool)
		httpOnly := cookie["http_only"].(bool)

		maxAge := *inputMaxAge
		if v, ok := cookie["max_age"].(int); ok && v != 0 {
			maxAge = v
		}

		c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
	}

	return nil
}

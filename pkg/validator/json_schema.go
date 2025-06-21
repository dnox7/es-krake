package validator

import (
	"fmt"
	"os"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/xeipuuv/gojsonschema"
)

var cachedSchemaInstances map[string]*gojsonschema.Schema

type JsonSchemaValidator struct {
	basePath string
	schemas  map[string]*gojsonschema.Schema
}

func NewJsonSchemaValidator(cfg *config.Config) (*JsonSchemaValidator, error) {
	validator := &JsonSchemaValidator{
		basePath: cfg.App.SchemasPath + "/" + cfg.App.ServiceName,
		schemas:  make(map[string]*gojsonschema.Schema),
	}

	gojsonschema.FormatCheckers.Add("date-time", NonStandardDatetimeFormatChecker{})
	gojsonschema.FormatCheckers.Add("strong-password", StrongPassswordChecker{})
	gojsonschema.FormatCheckers.Add("domain", DomainChecker{})
	gojsonschema.FormatCheckers.Add("url", UrlChecker{})
	gojsonschema.FormatCheckers.Add("id-sns", IDSnSChecker{})
	gojsonschema.FormatCheckers.Add("string_with_max_length", MaxLengthChecker{})

	err := validator.loadDirSchemas("")
	if err != nil {
		return nil, err
	}

	return validator, nil
}

func (validator *JsonSchemaValidator) Validate(
	schemaFile string,
	data interface{},
) (*gojsonschema.Result, error) {
	if schemaFile[0] != '/' {
		schemaFile = "/" + schemaFile
	}

	schema, exists := validator.schemas[schemaFile]
	if !exists {
		return nil, fmt.Errorf("the schema '%v' was not found for json validation", schemaFile)
	}

	dataLoader := gojsonschema.NewGoLoader(data)
	res, err := schema.Validate(dataLoader)
	if err != nil {
		return nil, err
	}

	if len(res.Errors()) == 0 {
		return nil, nil
	}

	return res, nil
}

func (validator *JsonSchemaValidator) loadDirSchemas(path string) error {
	schemaFiles, err := os.ReadDir(validator.basePath + path)
	if err != nil {
		return err
	}

	for _, file := range schemaFiles {
		if file.Name() == ".gitkeep" {
			continue
		}

		schemaPath := path + "/" + file.Name()
		if file.IsDir() {
			if err := validator.loadDirSchemas(schemaPath); err != nil {
				return err
			}
			continue
		}

		schemaSource := "file://" + validator.basePath + schemaPath
		var schema *gojsonschema.Schema
		var exists bool
		if schema, exists = cachedSchemaInstances[schemaSource]; !exists {
			schemaLoader := gojsonschema.NewReferenceLoader(schemaSource)
			schema, err := gojsonschema.NewSchema(schemaLoader)
			if err != nil {
				return err
			}

			cachedSchemaInstances[schemaSource] = schema
		}

		validator.schemas[schemaSource] = schema
	}

	return nil
}

func (validator *JsonSchemaValidator) GetErrorDetails(
	res gojsonschema.ResultError,
) map[string]interface{} {
	return map[string]interface{}{
		"context":     res.Context(),
		"description": res.Description(),
		"details":     res.Details(),
		"field":       res.Field(),
		"type":        res.Type(),
		"value":       res.Value(),
	}
}

func (validator *JsonSchemaValidator) GetErrorField(
	res gojsonschema.ResultError,
) string {
	field := res.Field()
	errorDetails := res.Details()
	if property, exists := errorDetails["property"]; exists {
		if propertyStr, ok := property.(string); ok {
			field = field + "." + propertyStr
		}
	}
	return field
}

func (validator *JsonSchemaValidator) GetCustomErrorMessage(
	res gojsonschema.ResultError,
) string {
	details := res.Details()
	format, formatExists := details["format"]

	if res.Type() == "format" && formatExists {
		switch format {
		case "email", "idn-email":
			return utils.ErrorEmailFail
		case "password", "strong-password":
			return utils.ErrorPasswordFail
		case "string_with_max_length":
			return utils.ErrorCheckMaxLengthUnder50Characters
		case "domain":
			return utils.ErrorInvalidDomain
		}
	}

	if res.Type() == "required" {
		return utils.ErrorInputRequired
	}

	minVal, minValExists := details["min"]
	if res.Type() == "string_gte" && minValExists {
		if minVal == 1 {
			return utils.ErrorInputFail
		}
	}

	_, maxValExists := details["max"]
	if res.Type() == "string_lte" && maxValExists {
		return utils.ErrorInputCharacterLimit
	}

	if res.Type() == "max_length_byte" {
		return utils.ErrorInputByteLimit
	}

	return utils.ErrorInputFail
}

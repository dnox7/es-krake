package validator

import (
	"os"

	"github.com/xeipuuv/gojsonschema"
)

var cachedSchemaInstances map[string]*gojsonschema.Schema

type JsonSchemaValidator struct {
	basePath string
	schemas  map[string]*gojsonschema.Schema
}

func NewJsonSchemaValidator() (*JsonSchemaValidator, error) {
	validator := &JsonSchemaValidator{
		basePath: os.Getenv("PE_SCHEMAS_PATH") + "/" + os.Getenv("PE_SERVICE_NAME"),
		schemas:  make(map[string]*gojsonschema.Schema),
	}

	return nil, nil
}

func (validator *JsonSchemaValidator) Validate(
	schemaFile string,
	data interface{},
) (*gojsonschema.Result, error) {
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

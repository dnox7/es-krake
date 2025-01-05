package mount

import (
	"fmt"
	"pech/es-krake/internal/services/main-service/pkg/container"
	"pech/es-krake/pkg/shared/validator"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func MountAll(
	repositories *container.RepositoryContainers,
	services *container.ServiceContainers,
	ginServer *gin.Engine,
	db *gorm.DB,
	logger *logrus.Logger,
) error {
	validator, err := validator.NewJsonSchemaValidator()
	if err != nil {
		return fmt.Errorf("Failed to create a JSON Schema Validator: %w", err)
	}

	graphql, err := container.NewGraphQLSchema(repositories, services, db, logger)
	if err != nil {
		return fmt.Errorf("Failed to create a new GrapQL Schema: %w", err)
	}

	return nil
}

package container

import (
	// baseRepo "pech/es-krake/pkg/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepositoryContainers struct {
}

func NewRepositoryContainers(db *gorm.DB, logger *logrus.Logger) (*RepositoryContainers, error) {
	// base := baseRepo.NewBaseRepository(logger)

	return &RepositoryContainers{}, nil
}

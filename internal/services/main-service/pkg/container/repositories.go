package container

import (
	// baseRepo "pech/es-krake/pkg/repository"

	"gorm.io/gorm"
)

type RepositoryContainers struct {
}

func NewRepositoryContainers(db *gorm.DB) (*RepositoryContainers, error) {
	// base := baseRepo.NewBaseRepository(logger)

	return &RepositoryContainers{}, nil
}

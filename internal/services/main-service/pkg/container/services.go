package container

import (
	"pech/es-krake/pkg/log"

	"gorm.io/gorm"
)

type ServiceContainers struct {
}

func NewServiceContainers(repositoryContainers *RepositoryContainers, db *gorm.DB, logger *log.Logger) *ServiceContainers {
	return &ServiceContainers{}
}

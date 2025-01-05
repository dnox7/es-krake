package container

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServiceContainers struct {
}

func NewServiceContainers(repositoryContainers *RepositoryContainers, db *gorm.DB, logger *logrus.Logger) *ServiceContainers {
	return &ServiceContainers{}
}

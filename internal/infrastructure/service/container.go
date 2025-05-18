package service

import (
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	authService "github.com/dpe27/es-krake/internal/infrastructure/service/auth"
)

type ServicesContainer struct {
	AuthContainer authService.ServiceContainer
}

func NewServicesContainer(repos repository.RepositoriesContainer) *ServicesContainer {
	return &ServicesContainer{
		AuthContainer: authService.NewServiceContainer(repos.AuthContainer),
	}
}

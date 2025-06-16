package service

import (
	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/infrastructure/keycloak"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	authService "github.com/dpe27/es-krake/internal/infrastructure/service/auth"
)

type ServicesContainer struct {
	AuthContainer     authService.ServiceContainer
	KeycloakContainer keycloak.ServiceContainer
}

func NewServicesContainer(
	cfg *config.Config,
	repos *repository.RepositoriesContainer,
	cache redis.RedisRepository,
) *ServicesContainer {
	return &ServicesContainer{
		AuthContainer:     authService.NewServiceContainer(repos.AuthContainer, cache),
		KeycloakContainer: keycloak.NewServiceContainer(cfg),
	}
}

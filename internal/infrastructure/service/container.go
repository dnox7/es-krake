package service

import (
	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/infrastructure/aws"
	"github.com/dpe27/es-krake/internal/infrastructure/keycloak"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	authService "github.com/dpe27/es-krake/internal/infrastructure/service/auth"
	buyerService "github.com/dpe27/es-krake/internal/infrastructure/service/buyer"
	mailService "github.com/dpe27/es-krake/internal/infrastructure/service/mail"
)

type ServicesContainer struct {
	MailService       mailService.MailService
	AuthContainer     authService.ServiceContainer
	KeycloakContainer keycloak.ServiceContainer
	BuyerContainer    buyerService.ServiceContainer
}

func NewServicesContainer(
	cfg *config.Config,
	ses aws.SesService,
	repos *repository.RepositoriesContainer,
	cache redis.RedisRepository,
) *ServicesContainer {
	return &ServicesContainer{
		AuthContainer:     authService.NewServiceContainer(repos.AuthContainer, cache),
		KeycloakContainer: keycloak.NewServiceContainer(cfg),
		MailService:       mailService.NewMailService(ses),
		BuyerContainer:    buyerService.NewServiceContainer(repos.BuyerContainer, cache),
	}
}

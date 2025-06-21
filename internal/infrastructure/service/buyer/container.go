package service

import (
	domainService "github.com/dpe27/es-krake/internal/domain/buyer/service"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	repository "github.com/dpe27/es-krake/internal/infrastructure/repository/buyer"
)

type ServiceContainer struct {
	BuyerService domainService.BuyerService
}

func NewServiceContainer(
	repos repository.RepositoryContainer,
	redisRepo redis.RedisRepository,
) ServiceContainer {
	return ServiceContainer{
		BuyerService: NewBuyerService(),
	}
}

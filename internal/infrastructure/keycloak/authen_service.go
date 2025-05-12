package keycloak

import "github.com/dpe27/es-krake/pkg/log"

type KcAuthenticationService interface {
}

type authenService struct {
	BaseKcService
	logger *log.Logger
}

func NewKcAuthenticationService(base BaseKcService) KcAuthenticationService {
	return &authenService{
		BaseKcService: base,
		logger:        log.With("service", "keycloak_authentication_service"),
	}
}

package keycloak

import (
	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/infrastructure/httpclient"
)

type ServiceContainer struct {
	ClientService           KcClientService
	IdentityProviderService KcIdentityProviderService
	RealmService            KcRealmService
	UserService             KcUserService
}

func NewServiceContainer(cfg *config.Config) ServiceContainer {
	args := httpclient.ClientOptBuilder().
		ServiceName("keycloak_api").
		Build()
	client := httpclient.NewHttpClient(args)
	base := NewBaseKcSevice(cfg, client)

	return ServiceContainer{
		ClientService:           NewKcClientService(base),
		IdentityProviderService: NewKcIdentityProviderService(base),
		RealmService:            NewKcRealmService(base),
		UserService:             NewKcUserService(base),
	}
}

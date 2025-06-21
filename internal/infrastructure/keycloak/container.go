package keycloak

import (
	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/domain/auth/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/httpclient"
)

type ServiceContainer struct {
	AttackDetectionService  KcAttackDetectionService
	AuthenticationService   KcAuthenticationService
	ClientService           KcClientService
	IdentityProviderService KcIdentityProviderService
	KeyService              KcKeyService
	RealmService            KcRealmService
	TokenService            KcTokenService
	UserService             KcUserService
	SnsProviderService      SnsProviderService
}

func NewServiceContainer(cfg *config.Config) ServiceContainer {
	args := httpclient.ClientOptBuilder().
		ServiceName("keycloak_api").
		Build()
	client := httpclient.NewHttpClient(args)
	base := NewBaseKcSevice(cfg, client)

	googleProviderCli := httpclient.NewHttpClient(
		httpclient.ClientOptBuilder().
			ServiceName("google_provider_api").
			Build(),
	)
	providersCli := map[int]httpclient.HttpClient{
		int(repository.GoogleProvider): googleProviderCli,
	}

	return ServiceContainer{
		AttackDetectionService:  NewKcAttackDetectionService(base),
		AuthenticationService:   NewKcAuthenticationService(base),
		ClientService:           NewKcClientService(base),
		IdentityProviderService: NewKcIdentityProviderService(base),
		KeyService:              NewKcKeyService(base),
		RealmService:            NewKcRealmService(base),
		TokenService:            NewKcTokenService(base),
		UserService:             NewKcUserService(base),
		SnsProviderService:      NewSnsProviderService(providersCli),
	}
}

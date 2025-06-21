package keycloak

import (
	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/infrastructure/httpclient"
)

const (
	openIdConnectPath = "/protocol/openid-connect"
	tokenPath         = openIdConnectPath + "/token"
	pubKeyPath        = openIdConnectPath + "/certs"
	logoutPath        = openIdConnectPath + "/logout"

	adminRealmPath = "/auth/admin/realms"
	realmPath      = "/auth/realms"

	masterRealm = "master"
)

type BaseKcService interface {
	Client() httpclient.HttpClient

	AccessTokenLifespan() uint
	RefreshTokenLifespan() uint

	masterRealmClientID() string
	masterRealmClientSecret() string

	AdminApiBaseUrl() string
	AdminRealmUrl(realm string) string
	PublicKeyUrl(issuer string) string
	TokenUrl(realmUrl string) string
	LogoutUrl(realm string) string
}

type baseKcService struct {
	client               httpclient.HttpClient
	domain               string
	clientID             string
	clientSecret         string
	accessTokenLifespan  uint
	refreshTokenLifespan uint
}

func NewBaseKcSevice(cfg *config.Config, cli httpclient.HttpClient) BaseKcService {
	return &baseKcService{
		client:               cli,
		domain:               cfg.Keycloak.Domain,
		clientID:             cfg.Keycloak.ClientID,
		clientSecret:         cfg.Keycloak.ClientSecret,
		accessTokenLifespan:  cfg.Keycloak.AccessTokenLifespan,
		refreshTokenLifespan: cfg.Keycloak.RefreshTokenLifespan,
	}
}

// masterRealmClientID implements BaseKcService.
func (b *baseKcService) masterRealmClientID() string {
	return b.clientID
}

// masterRealmClientSecret implements BaseKcService.

func (b *baseKcService) masterRealmClientSecret() string {
	return b.clientSecret
}

// Client implements BaseKcService.
func (b *baseKcService) Client() httpclient.HttpClient {
	return b.client
}

// AccessTokenLifespan implements BaseKcService.
func (b *baseKcService) AccessTokenLifespan() uint {
	return b.accessTokenLifespan
}

// RefreshTokenLifespan implements BaseKcService.
func (b *baseKcService) RefreshTokenLifespan() uint {
	return b.refreshTokenLifespan
}

// AdminRealmUrl implements BaseKcService.
func (b *baseKcService) AdminApiBaseUrl() string {
	return b.domain + adminRealmPath
}

// RealmUrl implements BaseKcService.
func (b *baseKcService) AdminRealmUrl(realm string) string {
	return b.domain + adminRealmPath + "/" + realm
}

// PublicKeyUrl implements BaseKcService.
func (b *baseKcService) PublicKeyUrl(issuer string) string {
	return issuer + pubKeyPath
}

// TokenUrl implements BaseKcService.
func (b *baseKcService) TokenUrl(realm string) string {
	return b.domain + realmPath + "/" + realm + tokenPath
}

// LogoutUrl implements BaseKcService.
func (b *baseKcService) LogoutUrl(realm string) string {
	return b.domain + realmPath + "/" + realm + logoutPath
}

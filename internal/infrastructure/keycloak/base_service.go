package keycloak

import (
	"fmt"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/infrastructure/httpclient"
)

const (
	openIdConnectPath = "/protocol/openid-connect"
	adminRealmUrlFmt  = "%s/auth/admin/realms/%s"
	tokenUrlFmt       = "%s/auth/realms/%s" + openIdConnectPath + "/token"
	publicKeyUrlFmt   = "%s/auth/realms/%s" + openIdConnectPath + "/certs"
)

type BaseKcService interface {
	Client() httpclient.HttpClient

	ClientID() string
	ClientSecret() string
	AccessTokenLifespan() uint
	RefreshTokenLifespan() uint

	AdminRealmUrl(realm string) string
	PublicKeyUrl(realm string) string
	TokenUrl(realm string) string
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

// Client implements BaseKcService.
func (b *baseKcService) Client() httpclient.HttpClient {
	return b.client
}

// ClientID implements BaseKcService.
func (b *baseKcService) ClientID() string {
	return b.clientID
}

// ClientSecret implements BaseKcService.
func (b *baseKcService) ClientSecret() string {
	return b.clientSecret
}

// AccessTokenLifespan implements BaseKcService.
func (b *baseKcService) AccessTokenLifespan() uint {
	return b.accessTokenLifespan
}

// RefreshTokenLifespan implements BaseKcService.
func (b *baseKcService) RefreshTokenLifespan() uint {
	return b.refreshTokenLifespan
}

// RealmUrl implements BaseKcService.
func (b *baseKcService) AdminRealmUrl(realm string) string {
	return fmt.Sprintf(adminRealmUrlFmt, b.domain, realm)
}

// PublicKeyUrl implements BaseKcService.
func (b *baseKcService) PublicKeyUrl(realm string) string {
	return fmt.Sprintf(publicKeyUrlFmt, b.domain, realm)
}

// TokenUrl implements BaseKcService.
func (b *baseKcService) TokenUrl(realm string) string {
	return fmt.Sprintf(tokenUrlFmt, b.domain, realm)
}

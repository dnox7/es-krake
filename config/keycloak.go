package config

type keycloak struct {
	Domain               string `env:"KEYCLOAK_DOMAIN"                 env-required:"true"`
	ClientID             string `env:"KEYCLOAK_CLIENT_ID"              env-required:"true"`
	ClientSecret         string `env:"KEYCLOAK_CLIENT_SECRET"          env-required:"true"`
	AccessTokenLifespan  uint   `env:"KEYCLOAK_ACCESS_TOKEN_LIFESPAN"  env-required:"true"`
	RefreshTokenLifespan uint   `env:"KEYCLOAK_REFRESH_TOKEN_LIFESPAN" env-required:"true"`
}

package config

type keycloak struct {
	Domain               string `yaml:"domain" env:"KEYCLOAK_DOMAIN" env-required:"true"`
	ClientID             string `yaml:"client_id" env:"KEYCLOAK_CLIENT_ID" env-required:"true"`
	ClientSecret         string `yaml:"client_secret" env:"KEYCLOAK_CLIENT_SECRET" env-required:"true"`
	AccessTokenLifespan  uint   `yaml:"access_token_lifespan" env:"KEYCLOAK_ACCESS_TOKEN_LIFESPAN" env-required:"true"`
	RefreshTokenLifespan uint   `yaml:"refresh_token_lifespan" env:"KEYCLOAK_REFRESH_TOKEN_LIFESPAN" env-required:"true"`
}

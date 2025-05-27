package config

type vault struct {
	address  string `env:"VAULT_ADDRESS"   env-required:"true"`
	roleID   string `env:"VAULT_ROLE_ID"   env-required:"true"`
	secretID string `env:"VAULT_SECRET_ID" env-required:"true"`
}

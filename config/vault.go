package config

type vault struct {
	Address      string `env:"VAULT_ADDRESS"        env-required:"true"`
	RoleID       string `env:"VAULT_ROLE_ID"        env-required:"true"`
	SecretIDFile string `env:"VAULT_SECRET_ID_FILE" env-required:"true"`

	RdbCredentialsPath   string `env:"VAULT_RDB_CREDENTIALS_PATH"   env-required:"true"`
	RedisCredentialsPath string `env:"VAULT_REDIS_CREDENTIALS_PATH" env-required:"true"`
}

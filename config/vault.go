package config

type vault struct {
	Address      string `env:"VAULT_ADDRESS"        env-required:"true"`
	RoleID       string `env:"VAULT_ROLE_ID"        env-required:"true"`
	SecretIDFile string `env:"VAULT_SECRET_ID_FILE" env-required:"true"`

	ESCredentialsPath    string `env:"VAULT_ELASTICSEARCH_CREDENTIALS_PATH" env-required:"true"`
	RdbCredentialsPath   string `env:"VAULT_RDB_CREDENTIALS_PATH"           env-required:"true"`
	MongoCredentialsPath string `env:"VAULT_MONGO_CREDENTIALS_PATH"         env-required:"true"`
	RedisCredentialsPath string `env:"VAULT_REDIS_CREDENTIALS_PATH"         env-required:"true"`
	RedisUsernameKey     string `env:"VAULT_REDIS_USERNAME_KEY"             env-required:"true"`
	RedisPasswordKey     string `env:"VAULT_REDIS_PASSWORD_KEY"             env-required:"true"`
}

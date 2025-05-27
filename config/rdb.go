package config

type rdb struct {
	Driver   string `env:"DB_DRIVER"   env-required:"true"`
	Host     string `env:"DB_HOST"     env-required:"true"`
	Port     string `env:"DB_PORT"     env-required:"true"`
	Username string `env:"DB_USERNAME" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	Name     string `env:"DB_NAME"     env-required:"true"`
	SSLMode  string `env:"DB_SSLMODE"  env-default:"disable"`

	MaxOpenConns int `env:"DB_MAX_OPEN_CONNS" env-required:"true"`
	MaxIdleConns int `env:"DB_MAX_IDLE_CONNS" env-required:"true"`
	MaxLifeTime  int `env:"DB_MAX_LIFE_TIME"  env-required:"true"`
	MaxIdleTime  int `env:"DB_MAX_IDLE_TIME"  env-required:"true"`
	ConnTimeout  int `env:"DB_CONN_TIMEOUT"   env-default:"10000"`
	ConnAttempts int `env:"DB_CONN_ATTEMPTS"  env-default:"10"`

	MigrationsPath string `env:"DB_MIGRATIONS_PATH" env-required:"true"`
}

type RdbCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

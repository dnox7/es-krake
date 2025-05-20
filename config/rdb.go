package config

type rdb struct {
	Driver   string `yaml:"driver" env:"DB_DRIVER" env-required:"true"`
	Host     string `yaml:"host" env:"DB_HOST" env-required:"true"`
	Port     string `yaml:"port" env:"DB_PORT" env-required:"true"`
	Username string `yaml:"username" env:"DB_USERNAME" env-required:"true"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Name     string `yaml:"name" env:"DB_NAME" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSLMODE" env-default:"disable"`

	MaxOpenConns int `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS" env-required:"true"`
	MaxIdleConns int `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS" env-required:"true"`
	MaxLifeTime  int `yaml:"max_life_time" env:"DB_MAX_LIFE_TIME" env-required:"true"`
	MaxIdleTime  int `yaml:"max_idle_time" env:"DB_MAX_IDLE_TIME" env-required:"true"`
	ConnTimeout  int `yaml:"conn_timeout" env:"DB_CONN_TIMEOUT" env-default:"10000"`
	ConnAttempts int `yaml:"conn_attempts" env:"DB_CONN_ATTEMPTS" env-default:"10"`

	MigrationsPath string `yaml:"migrations_path" env:"DB_MIGRATIONS_PATH" env-required:"true"`
}

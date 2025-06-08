package config

type mongo struct {
	Hostname     string `env:"MONGO_HOST"          env-default:"localhost"`
	Port         string `env:"MONGO_PORT"          env-default:"27017"`
	Database     string `env:"MONGO_DATABASE"      env-required:"true"`
	Username     string `env:"MONGO_USERNAME"      env-required:"true"`
	Password     string `env:"MONGO_PASSWORD"      env-required:"true"`
	AuthSource   string `env:"MONGO_AUTH_DB"       env-required:"true"`
	Timeout      int    `env:"MONGO_TIMEOUT"       env-required:"true"`
	ConnTimeout  int    `env:"MONGO_CONN_TIMEOUT"  env-default:"30000"`
	PoolSize     int    `env:"MONGO_POOL_SIZE"     env-required:"true"`
	MaxIdleTime  int    `env:"MONGO_MAX_IDLE_TIME" env-required:"true"`
	ConnAttempts int    `env:"MONGO_CONN_ATTEMPTS" env-default:"3"`
}

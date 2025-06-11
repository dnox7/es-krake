package config

type mongo struct {
	Hostname     string `env:"MONGO_HOST"          env-default:"localhost"`
	Port         string `env:"MONGO_PORT"          env-default:"27017"`
	Database     string `env:"MONGO_DATABASE"      env-required:"true"`
	AuthSource   string `env:"MONGO_AUTH_DB"       env-required:"true"`
	Timeout      int    `env:"MONGO_TIMEOUT"       env-default:"30000"`
	ConnTimeout  int    `env:"MONGO_CONN_TIMEOUT"  env-default:"30000"`
	PoolSize     int    `env:"MONGO_POOL_SIZE"     env-default:"10"`
	MaxIdleTime  int    `env:"MONGO_MAX_IDLE_TIME" env-default:"300000"`
	ConnAttempts int    `env:"MONGO_CONN_ATTEMPTS" env-default:"3"`
}

type MongoCredentials struct {
	Username string
	Password string
}

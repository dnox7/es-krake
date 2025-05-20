package config

type mongo struct {
	Hostname     string `yaml:"host" env:"MONGO_HOST" env-default:"localhost"`
	Port         string `yaml:"port" env:"MONGO_PORT" env-default:"27017"`
	Database     string `yaml:"database" env:"MONGO_DATABASE" env-required:"true"`
	Username     string `yaml:"username" env:"MONGO_USERNAME" env-required:"true"`
	Password     string `yaml:"password" env:"MONGO_PASSWORD" env-required:"true"`
	AuthSource   string `yaml:"auth_db" env:"MONGO_AUTH_DB" env-required:"true"`
	Timeout      int    `yaml:"timeout" env:"MONGO_TIMEOUT" env-required:"true"`
	ConnTimeout  int    `yaml:"conn_timeout" env:"MONGO_CONN_TIMEOUT" env-default:"30000"`
	PoolSize     int    `yaml:"pool_size" env:"MONGO_POOL_SIZE" env-required:"true"`
	MaxIdleTime  int    `yaml:"max_idle_time" env:"MONGO_MAX_IDLE_TIME" env-required:"true"`
	ConnAttempts int    `yaml:"conn_attempts" env:"MONGO_CONN_ATTEMPTS" env-default:"3"`
}

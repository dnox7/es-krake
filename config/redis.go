package config

//nolint:tagliatelle
type Redis struct {
	Host       string `yaml:"host" env:"REDIS_HOST" env-required:"true"`
	Port       string `yaml:"port" env:"REDIS_PORT" env-required:"true"`
	ClientName string `yaml:"client-name" env:"REDIS_CLIENT_NAME" env-required:"true"`
	Username   string `yaml:"username" env:"REDIS_USERNAME" env-required:"true"`
	Password   string `yaml:"password" env:"REDIS_PASSWORD" env-required:"true"`

	MaxRetries     int ` yaml:"max_retries" env:"REDIS_MAX_RETRIES" env-default:"3"`
	PoolSize       int `yaml:"pool_size" env:"REDIS_POOL_SIZE" env-default:"10"`
	MaxIdleConns   int `yaml:"max_idle_conns" env:"REDIS_MAX_IDLE_CONNS" env-default:"0"`
	MaxActiveConns int `yaml:"max_active_conns" env:"REIDS_MAX_ACTIVE_CONNS" env-default:"10"`
	MaxIdleTime    int `yaml:"max_idle_time" env:"REIDS_MAX_IDLE_TIME" env-default:"30"`
	MaxLifeTime    int `yaml:"max_life_time" env:"REDIS_MAX_LIFE_TIME" env-default:"10"`
}

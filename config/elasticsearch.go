package config

type elasticseach struct {
	Addresses []string `env:"ES_ADDRESSES" env-required:"true"`
	Username  string   `env:"ES_USERNAME"  env-required:"true"`
	Password  string   `env:"ES_PASSWORD"  env-required:"true"`

	MaxRetries    int  `env:"ES_MAX_RETRIES"    env-default:"3"`
	EnableMetrics bool `env:"ES_ENABLE_METRICS" env-default:"false"`
	Debug         bool `env:"ES_DEBUG"          env-default:"true"`
}

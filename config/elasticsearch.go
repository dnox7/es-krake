package config

type elasticseach struct {
	Addresses     []string `env:"ES_ADDRESSES"      env-required:"true"`
	MaxRetries    int      `env:"ES_MAX_RETRIES"    env-default:"3"`
	EnableMetrics bool     `env:"ES_ENABLE_METRICS" env-default:"false"`
	Debug         bool     `env:"ES_DEBUG"          env-default:"true"`
}

type ElasticSearchCredentials struct {
	Username string
	Password string
}

package config

type elasticseach struct {
	Addresses []string `yaml:"addresses" env:"ES_ADDRESSES" env-required:"true"`
	Username  string   `yaml:"username" env:"ES_USERNAME" env-required:"true"`
	Password  string   `yaml:"password" env:"ES_PASSWORD" env-required:"true"`

	MaxRetries    int  `yaml:"max_retries" env:"ES_MAX_RETRIES" env-default:"3"`
	EnableMetrics bool `yaml:"enable_metrics" env:"ES_ENABLE_METRICS" env-default:"false"`
	Debug         bool `yaml:"debug" env:"ES_DEBUG" env-default:"true"`
}

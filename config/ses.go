package config

type ses struct {
	Region               string `env:"SES_REGION"`
	CredentialsID        string `env:"SES_CREDENTIALS_ID"`
	CredentialsSecret    string `env:"SES_CREDENTIALS_SECRET"`
	CredentialsToken     string `env:"SES_CREDENTIALS_TOKEN"`
	ConfigurationSetName string `env:"SES_CONFIGURATION_SET"`
}

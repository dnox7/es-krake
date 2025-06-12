package config

type s3 struct {
	Region            string `env:"S3_REGION" env-default:"ap-northeast-1"`
	Endpoint          string `env:"S3_ENDPOINT" env-default:"localhost:9000"`
	DisableSSL        bool   `env:"S3_DISABLE_SSL" env-default:"true"`
	ForcePathStype    bool   `env:"S3_FORCE_PATH_STYLE" env-default:"true"`
	CredentialsID     string `env:"S3_CREDENTIALS_ID" env-default:"minioadmin"`
	CredentialsSecret string `env:"S3_CREDENTIALS_SECRET" env-default:"minioadmin"`
	CredentialsToken  string `env:"S3_CREDENTIALS_TOKEN" env-default:""`
}

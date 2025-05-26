package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

const cfgFilePath = ".env"

type (
	Config struct {
		App      *app      `yaml:"application"`
		Log      *log      `yaml:"logger"`
		Server   *server   `yaml:"server"`
		RDB      *rdb      `yaml:"rdb"`
		MDB      *mongo    `yaml:"mdb"`
		Keycloak *keycloak `yaml:"keycloak"`
		Redis    *Redis    `yaml:"redis"`
	}

	app struct {
		Name    string `yaml:"name" env:"APP_NAME" env-required:"true"`
		Version string `yaml:"version" env:"APP_VERSION" env-required:"true"`
		Env     string `yaml:"env" env:"APP_ENV" env-required:"true"`
	}

	server struct {
		Port string `yaml:"port" env:"SERVER_PORT" env-required:"true"`
	}

	log struct {
		Level string `yaml:"level" env:"LOG_LEVEL" env-required:"true"`
	}
)

func NewConfig() *Config {
	cfg := &Config{}
	root := projectRoot()
	configFilePath := root + cfgFilePath

	err := loadCfgFile(configFilePath, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func loadCfgFile(cfgFilePath string, cfg *Config) error {
	envFileExists := checkFileExists(cfgFilePath)
	if envFileExists {
		err := cleanenv.ReadConfig(cfgFilePath, cfg)
		if err != nil {
			return fmt.Errorf("config error: %w", err)
		}
	} else {
		err := cleanenv.ReadEnv(cfg)
		if err != nil {
			if _, statErr := os.Stat(cfgFilePath); statErr != nil {
				return fmt.Errorf("missing environment variable: %w", err)
			}
			return err
		}
	}
	return nil
}

func checkFileExists(fileName string) bool {
	exist := false
	if _, err := os.Stat(fileName); err == nil {
		exist = true
	}
	return exist
}

func projectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	cwd := filepath.Dir(b)
	return cwd + "/../"
}

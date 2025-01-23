package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	defaultConfigFileType = "yaml"
)

type (
	Config struct {
		App `yaml:"app"`
		Log `yaml:"logger"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	Log struct {
		Level  string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
		Format string `env-required:"true" yaml:"format" env:"LOG_FORMAT"`
	}
)

func NewConfig() *Config {
	cfg := &Config{}
	cwd := projectRoot()
	envFilePath := cwd + ".env"

	err := readEnv(envFilePath, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func readEnv(envFilePath string, cfg *Config) error {
	envFileExists := checkEnvFileExists(envFilePath)
	if envFileExists {
		err := cleanenv.ReadConfig(envFilePath, cfg)
		if err != nil {
			return fmt.Errorf("config error: %w", err)
		}
	} else {
		err := cleanenv.ReadEnv(cfg)
		if err != nil {
			if _, statErr := os.Stat(envFilePath + "." + defaultConfigFileType); statErr != nil {
				return fmt.Errorf("missing environment variables: %w", err)
			}
			return err
		}
	}
	return nil
}

func checkEnvFileExists(fileName string) bool {
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

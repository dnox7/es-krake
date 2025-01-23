package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

const defaultCfgFilePath = "config/default.yml"

type (
	Config struct {
		App    `yaml:"application"`
		Log    `yaml:"logger"`
		Server `yaml:"server"`
		RDB    `yaml:"rdb"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	Server struct {
		Port string `env-required:"true" yaml:"port" env:"SERVER_PORT"`
	}

	Log struct {
		Level  string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
		Format string `env-required:"true" yaml:"format" env:"LOG_FORMAT"`
	}

	RDB struct {
		Driver   string `env-required:"true" yaml:"driver" env:"DB_DRIVER"`
		Host     string `env-required:"true" yaml:"host" env:"DB_HOST"`
		Port     string `env-required:"true" yaml:"port" env:"DB_PORT"`
		Username string `env-required:"true" yaml:"username" env:"DB_USERNAME"`
		Password string `env-required:"true" yaml:"password" env:"DB_PASSWORD"`
		Name     string `env-required:"true" yaml:"name" env:"DB_NAME"`
		SSLMode  string `yaml:"ssl_mode" env:"DB_SSLMODE" env-default:"disable"`

		MaxOpenConns int `env-required:"true" yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS"`
		MaxIdleConns int `env-required:"true" yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
		MaxLifeTime  int `env-required:"true" yaml:"max_life_time" env:"DB_MAX_LIFE_TIME"`
		MaxIdleTime  int `env-required:"true" yaml:"max_idle_time" env:"DB_MAX_IDLE_TIME"`
		ConnTimeout  int `yaml:"conn_timeout" env:"DB_CONN_TIMEOUT" env-default:"1000"`
		ConnAttempts int `yaml:"conn_attempts" env:"DB_CONN_ATTEMPTS" env-default:"10"`
	}
)

func NewConfig() *Config {
	cfg := &Config{}
	root := projectRoot()
	configFilePath := root + defaultCfgFilePath

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
		fmt.Println("ok")
		if err != nil {
			return fmt.Errorf("config error: %w", err)
		}
	} else {
		err := cleanenv.ReadEnv(cfg)
		if err != nil {
			if _, statErr := os.Stat(cfgFilePath); statErr != nil {
				return fmt.Errorf("missing environment variables: %w", err)
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

package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	vaultcli "github.com/dpe27/es-krake/internal/infrastructure/vault"
	"github.com/ilyakaznacheev/cleanenv"
)

const cfgFilePath = ".env"

type (
	Config struct {
		App      *app
		Log      *log
		Server   *server
		RDB      *rdb
		MDB      *mongo
		Keycloak *keycloak
		Vault    *vault
	}

	app struct {
		Name    string `env:"APP_NAME"    env-required:"true"`
		Version string `env:"APP_VERSION" env-required:"true"`
		Env     string `env:"APP_ENV"     env-required:"true"`
	}

	server struct {
		Port string `env:"SERVER_PORT" env-required:"true"`
	}

	log struct {
		Level string `env:"LOG_LEVEL" env-required:"true"`
	}
)

func NewConfig(ctx context.Context) *Config {
	cfg := &Config{}
	root := projectRoot()
	configFilePath := root + cfgFilePath

	err := loadCfg(ctx, configFilePath, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func loadCfg(ctx context.Context, cfgFilePath string, cfg *Config) error {
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

	vcli, token, err := vaultcli.NewVaultAppRoleClient(ctx, vaultcli.VaultParams{
		Address:            cfg.Vault.Address,
		RoleID:             cfg.Vault.RoleID,
		SecretIDFile:       cfg.Vault.SecretIDFile,
		RdbCredentialsPath: cfg.Vault.RdbCredentialsPath,
	})
	if err != nil {
		return err
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

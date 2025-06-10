package initializer

import (
	"context"
	"fmt"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb/migration"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	vaultcli "github.com/dpe27/es-krake/internal/infrastructure/vault"
	vault "github.com/hashicorp/vault/api"
)

func InitVault(ctx context.Context, cfg *config.Config) (*vaultcli.Vault, *vault.Secret, error) {
	v, token, err := vaultcli.NewVaultAppRoleClient(ctx, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize vault connect [address: %s]: %w", cfg.Vault.Address, err)
	}
	return v, token, nil
}

func InitPostgres(
	v *vaultcli.Vault,
	cfg *config.Config,
) (
	pg *rdb.PostgreSQL,
	credLease *vault.Secret,
	stopLogging context.CancelFunc,
	err error,
) {
	ctx := context.Background()
	var cred *config.RdbCredentials
	cred, credLease, err = v.GetRdbCredentials(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to retrieve database credentials from vault: %w", err)
	}

	pg = rdb.NewOrGetSingleton(cfg, cred)

	var loggingPoolSizeCtx context.Context
	loggingPoolSizeCtx, stopLogging = context.WithCancel(ctx)
	pg.LoggingPoolSize(loggingPoolSizeCtx)

	err = pg.Ping(ctx)
	if err != nil {
		stopLogging()
		return nil, nil, nil, fmt.Errorf("database ping failed: %w", err)
	}

	err = migration.CheckAll(cfg, pg.Conn())
	if err != nil {
		stopLogging()
		return nil, nil, nil, fmt.Errorf("the database is not up-to-date: %w", err)
	}

	return pg, credLease, stopLogging, nil
}

func InitRedis(v *vaultcli.Vault, cfg *config.Config) (*redis.RedisRepo, error) {
	ctx := context.Background()
	cred, err := v.GetRedisCredentials(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve redis credentials from vault: %w", err)
	}

	redisRepo := redis.NewRedisRespository(ctx, cfg, cred)
	return redisRepo, nil
}

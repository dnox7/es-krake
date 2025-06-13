package initializer

import (
	"context"
	"fmt"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/infrastructure/aws"
	es "github.com/dpe27/es-krake/internal/infrastructure/elasticsearch"
	"github.com/dpe27/es-krake/internal/infrastructure/httpclient"
	mdb "github.com/dpe27/es-krake/internal/infrastructure/mongodb"
	"github.com/dpe27/es-krake/internal/infrastructure/notify"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb/migration"
	"github.com/dpe27/es-krake/internal/infrastructure/redis"
	vaultcli "github.com/dpe27/es-krake/internal/infrastructure/vault"
	"github.com/dpe27/es-krake/pkg/log"
	vault "github.com/hashicorp/vault/api"
)

func InitVault(ctx context.Context, cfg *config.Config) (*vaultcli.Vault, *vault.Secret, error) {
	v, token, err := vaultcli.NewVaultAppRoleClient(ctx, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize vault connect [address: %s]: %w", cfg.Vault.Address, err)
	}
	return v, token, nil
}

func InitPostgres(v *vaultcli.Vault, cfg *config.Config) (
	pg *rdb.PostgreSQL,
	credLease *vault.Secret,
	stopLogging context.CancelFunc,
	err error,
) {
	ctx := context.Background()
	var cred *config.RdbCredentials
	cred, credLease, err = v.GetRdbCredentials(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to retrieve postgres credentials from vault: %w", err)
	}

	pg = rdb.NewOrGetSingleton(cfg, cred)

	var loggingPoolSizeCtx context.Context
	loggingPoolSizeCtx, stopLogging = context.WithCancel(ctx)
	pg.LoggingPoolSize(loggingPoolSizeCtx)

	err = pg.Ping(ctx)
	if err != nil {
		stopLogging()
		return nil, nil, nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	err = migration.CheckAll(cfg, pg.Conn())
	if err != nil {
		stopLogging()
		return nil, nil, nil, fmt.Errorf("the database is not up-to-date: %w", err)
	}

	return pg, credLease, stopLogging, nil
}

func InitMongo(v *vaultcli.Vault, cfg *config.Config) (m *mdb.Mongo, credLease *vault.Secret, err error) {
	ctx := context.Background()
	var cred *config.MongoCredentials
	cred, credLease, err = v.GetMongoCredentials(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to retrieve mongodb credentials from vault: %w", err)
	}

	m = mdb.NewOrGetSingleton(ctx, cfg, cred)
	err = m.Ping(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("mongodb ping failed: %w", err)
	}

	return m, credLease, nil
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

func InitElasticSearch(v *vaultcli.Vault, cfg *config.Config) (*es.ElasticSearch, *vault.Secret, error) {
	ctx := context.Background()
	cred, lease, err := v.GetElasticSearchCredentials(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to retrieve elasticsearch credentials from vault: %w", err)
	}

	esRepo := es.NewElasticSearchRepository(cfg, cred)
	ok, err := esRepo.Ping(ctx)
	if err != nil || !ok {
		return nil, nil, fmt.Errorf("elasticsearch ping failed: %w", err)
	}

	return esRepo, lease, nil
}

func InitS3Repository(cfg *config.Config) (*aws.S3Repo, error) {
	s3Repo, err := aws.NewS3Repository(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init s3 repository: %w", err)
	}
	log.Debug(context.Background(), "init s3 successfully")
	return s3Repo, nil
}

func InitDiscordNotifier(cfg *config.Config) notify.DiscordNotifier {
	opts := httpclient.ClientOptBuilder().
		ServiceName("discord_notifier").
		Build()
	cli := httpclient.NewHttpClient(opts)
	return notify.NewDiscordNotifier(cfg.Discord.WebhookUrl, cli)
}

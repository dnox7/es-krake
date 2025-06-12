package vaultcli

import (
	"context"
	"time"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/utils"
	vault "github.com/hashicorp/vault/api"
)

type leaseType string

const (
	authTokeLeaseType        leaseType = "auth_token"
	rdbCredentialsLeaseType  leaseType = "rdb_credentials"
	mongoCredetialsLeaseType leaseType = "mongo_credentials"
)

func (v *Vault) PeriodicallyRenewLeases(
	ctx context.Context, authToken *vault.Secret,
	rdbCredLease *vault.Secret, rdbReconnFunc func(cred *config.RdbCredentials) error,
	mongoCredLease *vault.Secret, mongoReconnFunc func(ctx context.Context, cred *config.MongoCredentials) error,
) {
	v.logger.Info(ctx, "starting lease renewal watchers")

	force := utils.NewBroadcaster()

	go v.monitorLease(ctx, authTokeLeaseType, authToken, func() (*vault.Secret, error) {
		authToken, err := v.login(ctx)
		if err != nil {
			return nil, err
		}
		force.Broadcast()
		return authToken, nil
	}, nil)

	go v.monitorLease(ctx, rdbCredentialsLeaseType, rdbCredLease, func() (*vault.Secret, error) {
		cred, lease, err := v.GetRdbCredentials(ctx)
		if err != nil {
			return nil, err
		}

		err = rdbReconnFunc(cred)
		if err != nil {
			v.logger.Error(ctx, "failed to reconnect rdb", err, err.Error())
			return nil, err
		}
		return lease, nil
	}, force)

	go v.monitorLease(ctx, mongoCredetialsLeaseType, mongoCredLease, func() (*vault.Secret, error) {
		cred, lease, err := v.GetMongoCredentials(ctx)
		if err != nil {
			return nil, err
		}

		err = mongoReconnFunc(ctx, cred)
		if err != nil {
			v.logger.Error(ctx, "failed to reconnect mongodb", err, err.Error())
			return nil, err
		}
		return lease, nil
	}, force)
}

func (v *Vault) monitorLease(
	ctx context.Context,
	leaseName leaseType,
	secret *vault.Secret,
	secretFunc func() (*vault.Secret, error),
	force *utils.Broadcaster,
) {
	var (
		err     error
		forceCh chan struct{}
	)
	firstTime := true
	if force != nil {
		forceCh = force.Subscribe()
		defer force.Unsubscribe(forceCh)
	}

	for {
		if ctx.Err() != nil {
			return
		}

		if !firstTime {
			secret, err = secretFunc()
			if err != nil {
				v.logger.Error(ctx, "failed to fetch secret", "lease", string(leaseName), "error", err.Error())
				time.Sleep(5 * time.Second)
				continue
			}
		}

		watcher, err := v.client.NewLifetimeWatcher(&vault.LifetimeWatcherInput{Secret: secret})
		if err != nil {
			v.logger.Error(ctx, "failed to create watcher", "lease", string(leaseName), "error", err.Error())
			time.Sleep(5 * time.Second)
			firstTime = false
			continue
		}

		go watcher.Start()
		watching := true
		for watching {
			select {
			case <-ctx.Done():
				watcher.Stop()
				return

			case info := <-watcher.RenewCh():
				leaseDuration := info.Secret.LeaseDuration
				if info.Secret.Auth != nil {
					leaseDuration = info.Secret.Auth.LeaseDuration
				}
				v.logger.Info(ctx, "lease renewed", "lease", string(leaseName), "duration", leaseDuration)

			case err := <-watcher.DoneCh():
				watcher.Stop()
				v.logger.Warn(ctx, "lease expired or cannot be renewed", "lease", leaseName, "error", err)
				watching = false

			case <-forceCh:
				v.logger.Debug(ctx, "force get credentials", "lease", leaseName)
				watcher.Stop()
				watching = false
			}
		}
		firstTime = false
	}
}

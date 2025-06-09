package vaultcli

import (
	"context"
	"fmt"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"
	vault "github.com/hashicorp/vault/api"
)

type (
	renewResult    uint
	leaseEventType uint
	leaseType      string
)

const (
	authTokeLeaseType         leaseType = "auth_token"
	rdbCredentialsLeaseType   leaseType = "rdb_credentials"
	redisCredentialsLeaseType leaseType = "redis_credentials"

	renewType leaseEventType = iota
	doneType

	renewError renewResult = iota
	exitRequested
	expiringAuthToken
	expiringRdbCredentials
	expiringRedisCredentials
)

func (v *Vault) PeriodicallyRenewLeases(
	ctx context.Context,
	authToken *vault.Secret,
	rdbCredentialsLease *vault.Secret,
	rdbReconnectFunc func(cred *config.RdbCredentials) error,
	redisCredentialsLease *vault.Secret,
	redisReconnectFunc func(ctx context.Context, cred *config.RedisCredentials) error,
) {
	v.logger.Info(ctx, "renew/recreate secrets loops: bein")
	defer v.logger.Info(ctx, "renew/recreate secret loops: end")

	for {
		renewed, err := v.renewLeases(ctx, map[leaseType]*vault.Secret{
			authTokeLeaseType:       authToken,
			rdbCredentialsLeaseType: rdbCredentialsLease,
		})
		if err != nil {
			v.logger.Error(ctx, "failed to renew leases", "error", err.Error())
		}

		switch renewed {
		case exitRequested:
			return
		case expiringAuthToken:
			v.logger.Info(ctx, "auth token: can no longer be renewed; will login in again")
			newToken, err := v.login(ctx)
			if err != nil {
				log.Fatal(ctx, "failed to login into vault server", "error", err.Error())
			}
			authToken = newToken
		case expiringRdbCredentials:
			v.logger.Info(ctx, "rdb credentials: can no longer be renewed; will fetch new credentials and reconnect")
			newCred, newSecret, err := v.GetRdbCredentials(ctx)
			if err != nil {
				log.Fatal(ctx, "failed to fetch rdb credentials", "error", err.Error())
			}

			err = rdbReconnectFunc(newCred)
			if err != nil {
				log.Fatal(ctx, "failed to reconnect rdb", "error", err.Error())
			}
			rdbCredentialsLease = newSecret
		case expiringRedisCredentials:
			v.logger.Info(ctx, "redis credentials: can no longer be renewed; will fetch new credentials and reconnect")
			newCred, newSecret, err := v.GetRedisCredentials(ctx)

			err = redisReconnectFunc(ctx, newCred)
			if err != nil {
				log.Fatal(ctx, "failed to fetch redis credentials", "error", err.Error())
			}
			redisCredentialsLease = newSecret
		}
	}
}

type leaseEvent struct {
	name      leaseType
	eventType leaseEventType
	duration  int
	err       error
}

func (v *Vault) renewLeases(ctx context.Context, leases map[leaseType]*vault.Secret) (renewResult, error) {
	v.logger.Info(ctx, "renew cycle: begin")
	defer v.logger.Info(ctx, "renew cycle: end")

	eventCh := make(chan leaseEvent)
	defer close(eventCh)

	for name, secret := range leases {
		w, err := v.client.NewLifetimeWatcher(&vault.LifetimeWatcherInput{
			Secret: secret,
		})
		if err != nil {
			return renewError, fmt.Errorf("unable to init watcher for %s: %w", name, err)
		}

		go func(name leaseType, w *vault.LifetimeWatcher) {
			w.Start()
			defer w.Stop()

			for {
				select {
				case err := <-w.DoneCh():
					eventCh <- leaseEvent{
						name:      name,
						eventType: doneType,
						duration:  0,
						err:       err,
					}
					return
				case info := <-w.RenewCh():
					var dur int
					if info.Secret.Auth != nil {
						dur = info.Secret.Auth.LeaseDuration
					} else {
						dur = info.Secret.LeaseDuration
					}
					eventCh <- leaseEvent{
						name:      name,
						eventType: renewType,
						duration:  dur,
						err:       nil,
					}
				}
			}
		}(name, w)
	}

	for {
		select {
		case <-ctx.Done():
			return exitRequested, nil
		case evt := <-eventCh:
			if evt.eventType == doneType {
				return v.wrapErr(evt.name), evt.err
			}

			if evt.eventType == renewType {
				v.logger.Info(ctx, "successfully renewed", "lease", string(evt.name), "duration", evt.duration)
			}
		}
	}
}

func (v *Vault) wrapErr(lt leaseType) renewResult {
	switch lt {
	case authTokeLeaseType:
		return expiringAuthToken
	case rdbCredentialsLeaseType:
		return expiringRdbCredentials
	case redisCredentialsLeaseType:
		return expiringRedisCredentials
	default:
		return renewError
	}
}

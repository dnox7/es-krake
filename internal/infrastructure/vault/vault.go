package vaultcli

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/hashicorp/go-cleanhttp"
	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
)

type vaultParams struct {
	// connection parameters
	address      string
	roleID       string
	secretIDFile string

	rdbCredentialsPath   string
	redisCredentialsPath string
}

type Vault struct {
	client *vault.Client
	params vaultParams
	logger *log.Logger
}

func NewVaultAppRoleClient(ctx context.Context, cfg *config.Config) (*Vault, *vault.Secret, error) {
	logger := log.With("service", "vault")
	params := vaultParams{
		address:              cfg.Vault.Address,
		roleID:               cfg.Vault.RoleID,
		secretIDFile:         cfg.Vault.SecretIDFile,
		rdbCredentialsPath:   cfg.Vault.RdbCredentialsPath,
		redisCredentialsPath: cfg.Vault.RedisCredentialsPath,
	}

	logger.Info(ctx, "connecting to vault @", "address", params.address)

	conf := vault.DefaultConfig()
	conf.Address = params.address
	if cfg.App.Env == utils.DevEnv {
		tp := cleanhttp.DefaultPooledTransport()
		tp.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		conf.HttpClient.Transport = tp
	}

	client, err := vault.NewClient(conf)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize vault client: %w", err)
	}

	v := &Vault{
		client: client,
		params: params,
		logger: logger,
	}

	token, err := v.login(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("vault login error: %w", err)
	}

	logger.Info(ctx, "connecting to vault: success!")
	return v, token, nil
}

func (v *Vault) login(ctx context.Context) (*vault.Secret, error) {
	v.logger.Info(ctx, "logging in to vault with approle auth", "role id", v.params.roleID)
	approleSecretID := &approle.SecretID{
		FromFile: v.params.secretIDFile,
	}

	approleAuth, err := approle.NewAppRoleAuth(
		v.params.roleID,
		approleSecretID,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize approle authentication method: %w", err)
	}

	authInfo, err := v.client.Auth().Login(ctx, approleAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login using approle auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no approle info was retutned after login")
	}

	v.logger.Info(ctx, "logging in to vault with approle auth: success!")
	return authInfo, nil
}

func (v *Vault) GetRdbCredentials(ctx context.Context) (*config.RdbCredentials, *vault.Secret, error) {
	v.logger.Info(ctx, "getting database credentials from vault")
	lease, err := v.client.Logical().ReadWithContext(ctx, v.params.rdbCredentialsPath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read rdb secret: %w", err)
	}

	bytes, err := json.Marshal(lease.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("malformed rdb credentials returned: %w", err)
	}

	credentials := &config.RdbCredentials{}
	if err := json.Unmarshal(bytes, credentials); err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal rdb credentials: %w", err)
	}

	v.logger.Info(ctx, "getting rdb credentials from vault: success")
	return credentials, lease, nil
}

func (v *Vault) GetRedisCredentials(ctx context.Context) (*config.RedisCredentials, *vault.Secret, error) {
	v.logger.Info(ctx, "getting redis credentials from vault")
	lease, err := v.client.Logical().ReadWithContext(ctx, v.params.redisCredentialsPath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read rdb secret: %w", err)
	}

	bytes, err := json.Marshal(lease.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("malformed redis credentials returned: %w", err)
	}

	credentials := &config.RedisCredentials{}
	if err := json.Unmarshal(bytes, credentials); err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal redis credentials: %w", err)
	}

	v.logger.Info(ctx, "getting redis credentials from vault: success")
	return credentials, lease, nil
}

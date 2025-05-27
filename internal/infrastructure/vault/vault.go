package vault

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"
	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
)

type VaultParams struct {
	// connection parameters
	Address             string
	ApproleRoleID       string
	ApproleSecretIDFile string

	APIKeyPath              string
	APIKeyMountPath         string
	APIKeyField             string
	DatabaseCredentialsPath string
}

type Vault struct {
	client *vault.Client
	params VaultParams
	logger *log.Logger
}

func NewVaultAppRoleClient(ctx context.Context, params VaultParams) (*Vault, *vault.Secret, error) {
	logger := log.With("service", "vault")
	logger.Info(ctx, "connecting to vault @ %s", params.Address)

	conf := vault.DefaultConfig()
	conf.Address = params.Address

	client, err := vault.NewClient(conf)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to intialize vault client: %w", err)
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
	v.logger.Info(ctx, "logging in to vault with approle auth", "role id", v.params.ApproleRoleID)
	approleSecretID := &approle.SecretID{
		FromFile: v.params.ApproleSecretIDFile,
	}

	approleAuth, err := approle.NewAppRoleAuth(
		v.params.ApproleRoleID,
		approleSecretID,
		approle.WithWrappingToken(),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to intialize approle authentication method: %w", err)
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

func (v *Vault) GetDatabaseCredentials(ctx context.Context) (*config.RdbCredentials, *vault.Secret, error) {
	v.logger.Info(ctx, "getting database credentials from vault")
	lease, err := v.client.Logical().ReadWithContext(ctx, v.params.DatabaseCredentialsPath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read secret: %w", err)
	}

	bytes, err := json.Marshal(lease.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("malformed credentials returned: %w", err)
	}

	credentials := &config.RdbCredentials{}
	if err := json.Unmarshal(bytes, credentials); err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal credentials: %w", err)
	}

	v.logger.Info(ctx, "getting database credentials from vault: success!")
	return credentials, lease, nil
}

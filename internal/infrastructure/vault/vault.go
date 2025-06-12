package vaultcli

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/hashicorp/go-cleanhttp"
	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
)

const (
	secretMountPath               = "secret"
	errUnexpectedSecretKeyTypeFmt = "unexpected secret key type for %q field"
)

type vaultParams struct {
	// connection parameters
	address      string
	roleID       string
	secretIDFile string

	elasticsearchCredentialsPath string
	rdbCredentialsPath           string
	mongoCredetialsPath          string
	redisCredentialsPath         string
	redisUsernameKey             string
	redisPasswordKey             string
}

type Vault struct {
	client *vault.Client
	params vaultParams
	logger *log.Logger
}

func NewVaultAppRoleClient(ctx context.Context, cfg *config.Config) (*Vault, *vault.Secret, error) {
	logger := log.With("service", "vault")
	params := vaultParams{
		address:                      cfg.Vault.Address,
		roleID:                       cfg.Vault.RoleID,
		secretIDFile:                 cfg.Vault.SecretIDFile,
		elasticsearchCredentialsPath: cfg.Vault.ESCredentialsPath,
		rdbCredentialsPath:           cfg.Vault.RdbCredentialsPath,
		mongoCredetialsPath:          cfg.Vault.MongoCredentialsPath,
		redisCredentialsPath:         cfg.Vault.RedisCredentialsPath,
		redisUsernameKey:             cfg.Vault.RedisUsernameKey,
		redisPasswordKey:             cfg.Vault.RedisPasswordKey,
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
	v.logger.Info(ctx, "getting rdb credentials from vault")
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

func (v *Vault) GetMongoCredentials(ctx context.Context) (*config.MongoCredentials, *vault.Secret, error) {
	v.logger.Info(ctx, "getting mongodb credentials from vault")
	lease, err := v.client.Logical().ReadWithContext(ctx, v.params.mongoCredetialsPath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read mongodb secret: %w", err)
	}

	bytes, err := json.Marshal(lease.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("malformed mongodb credentials returned: %w", err)
	}

	credentials := &config.MongoCredentials{}
	if err := json.Unmarshal(bytes, credentials); err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal mongodb credentials: %w", err)
	}

	v.logger.Info(ctx, "getting mongodb credentials from vault: success")
	return credentials, lease, nil
}

func (v *Vault) GetRedisCredentials(ctx context.Context) (*config.RedisCredentials, error) {
	v.logger.Info(ctx, "getting redis credentials from vault")
	secret, err := v.client.KVv2(secretMountPath).Get(ctx, v.params.redisCredentialsPath)
	if err != nil {
		v.logger.Error(ctx, "can not get from kv-v2")
		return nil, fmt.Errorf("unable to read rdb secret: %w", err)
	}

	username, ok := secret.Data[v.params.redisUsernameKey]
	if !ok {
		return nil, errors.New("unable to get redis username")
	}

	usernameStr, ok := username.(string)
	if !ok {
		return nil, fmt.Errorf(errUnexpectedSecretKeyTypeFmt, v.params.redisUsernameKey)
	}

	password, ok := secret.Data[v.params.redisPasswordKey]
	if !ok {
		return nil, errors.New("unable to get redis password")
	}

	passwordStr, ok := password.(string)
	if !ok {
		return nil, fmt.Errorf(errUnexpectedSecretKeyTypeFmt, v.params.redisPasswordKey)
	}

	return &config.RedisCredentials{
		Username: usernameStr,
		Password: passwordStr,
	}, nil
}

func (v *Vault) GetElasticSearchCredentials(ctx context.Context) (*config.ElasticSearchCredentials, *vault.Secret, error) {
	v.logger.Info(ctx, "getting elasticsearch credentials from vault")
	lease, err := v.client.Logical().ReadWithContext(ctx, v.params.elasticsearchCredentialsPath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read elasticsearch secret: %w", err)
	}

	bytes, err := json.Marshal(lease.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("malformed elasticsearch credentials returned: %w", err)
	}

	credentials := &config.ElasticSearchCredentials{}
	if err := json.Unmarshal(bytes, credentials); err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal elasticsearch credentials: %w", err)
	}

	v.logger.Info(ctx, "getting elasticsearch credentials from vault: success")
	return credentials, lease, nil
}

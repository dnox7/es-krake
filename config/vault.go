package config

import (
	"context"
	"fmt"
	"os"
)

const (
	vaultRoleIDEnv   = "VAULT_ROLE_ID"
	vaultSecretIDEnv = "VAULT_SECRET_ID"
	vaultAddressEnv  = "VAULT_ADDRESS"
)

func initVaultConf(ctx context.Context) error {
	roleID := os.Getenv(vaultRoleIDEnv)
	secretID := os.Getenv(vaultSecretIDEnv)
	address := os.Getenv(vaultAddressEnv)

	if roleID == "" || secretID == "" || address == "" {
		var missing []string
		if roleID == "" {
			missing = append(missing, vaultRoleIDEnv)
		}
		if secretID == "" {
			missing = append(missing, vaultSecretIDEnv)
		}
		return fmt.Errorf("missing required environments variables: %v", missing)
	}

}

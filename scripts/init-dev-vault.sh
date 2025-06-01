#!/bin/sh
set -e

vault server -dev -dev-root-token-id=root -dev-tls &
VAULT_PID=$!

while ! vault status >/dev/null 2>&1; do
    echo "Waiting for Vault to be ready..."
    sleep 2
done

echo "Vault server is up!"

echo ">>> Logging in to Vault..."
vault login root >/dev/null

vault auth enable approle || true

for file in /vault/policies/*.hcl; do
    name=$(basename "$file" .hcl)
    echo "Applying policy: $name"
    vault policy write "$name" "$file"
done

echo ">>> Creating approle..."
vault write auth/approle/role/es-krake \
    token_type=service \
    secret_id_ttl=0m \
    secret_id_num_uses=0 \
    token_ttl=60m \
    token_max_ttl=120m \
    token_renewable=true \
    token_policies="read-secret,default"

vault read auth/approle/role/es-krake/role-id | awk '/role_id/ {print $2}' >/vault/config/role_id
vault write -f auth/approle/role/es-krake/secret-id | grep '^secret_id[[:space:]]' | awk '{ print $2 }' >/vault/config/secret_id

vault secrets enable database || true

vault write database/config/my-rds-database \
    plugin_name=postgresql-database-plugin \
    allowed_roles="my-role" \
    connection_url="postgresql://{{username}}:{{password}}@$ESK_RDB_ENDPOINT:5432/dbname?sslmode=require" \
    username="admin-user" \
    password="admin-password"

wait $VAULT_PID

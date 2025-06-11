#!/bin/sh
set -e

vault server -dev -dev-root-token-id=root -dev-tls &
VAULT_PID=$!

while ! vault status >/dev/null 2>&1; do
    echo "Waiting for Vault to be ready..."
    sleep 2
done

echo "Vault server is up!"

vault login root >/dev/null

vault auth enable approle || true

for file in /vault/policies/*.hcl; do
    name=$(basename "$file" .hcl)
    echo "Applying policy: $name"
    vault policy write "$name" "$file"
done

vault write auth/approle/role/es-krake-migration \
    token_type=service \
    secret_id_ttl=0m \
    secret_id_num_uses=0 \
    token_ttl=10m \
    token_max_ttl=12m \
    token_renewable=true \
    token_policies="read-migration-secret,default"

vault read auth/approle/role/es-krake-migration/role-id | awk '/role_id/ {print $2}' >/vault/config/role_id_migration
vault write -f auth/approle/role/es-krake-migration/secret-id | grep '^secret_id[[:space:]]' | awk '{ print $2 }' >/vault/config/secret_id_migration

vault write auth/approle/role/es-krake-api \
    token_type=service \
    secret_id_ttl=0m \
    secret_id_num_uses=0 \
    token_ttl=1h \
    token_max_ttl=2h \
    token_renewable=true \
    token_policies="read-app-secret,default"

vault read auth/approle/role/es-krake-api/role-id | awk '/role_id/ {print $2}' >/vault/config/role_id_api
vault write -f auth/approle/role/es-krake-api/secret-id | grep '^secret_id[[:space:]]' | awk '{ print $2 }' >/vault/config/secret_id_api

vault secrets enable database || true

vault write database/config/esk-rdb \
    plugin_name=postgresql-database-plugin \
    allowed_roles="postgres-app-role,postgres-migrate-role" \
    connection_url="postgres://{{username}}:{{password}}@esk-rdb:5432/$ESK_RDB_NAME?sslmode=disable" \
    username="$ESK_RDB_MASTER_USERNAME" \
    password="$ESK_RDB_MASTER_PASSWORD"

# for api and batch
vault write database/roles/postgres-app-role \
    db_name="esk-rdb" \
    creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; \
        GRANT CONNECT ON DATABASE $ESK_RDB_NAME TO \"{{name}}\";
        GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO \"{{name}}\";" \
    revocation_statements="
        REVOKE ALL PRIVILEGES ON DATABASE esk_dev_1 FROM \"{{name}}\";
        REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM \"{{name}}\";
        DROP ROLE \"{{name}}\";" \
    default_ttl="10m" \
    max_ttl="20m"

# for migration
vault write database/static-roles/postgres-migrate-role \
    db_name="esk-rdb" \
    rotation_statements="ALTER ROLE \"esk_dev_migrator\" WITH PASSWORD '{{password}}';" \
    username="esk_dev_migrator" \
    rotation_period="20m"

vault write database/config/esk-mdb \
    plugin_name=mongodb-database-plugin \
    allowed_roles="mongo-app-role" \
    connection_url="mongodb://{{username}}:{{password}}@esk-mdb:27017/admin?tls=false" \
    username="$ESK_MDB_MASTER_USERNAME" \
    password="$ESK_MDB_MASTER_PASSWORD"

# for mongo
vault write database/roles/mongo-app-role \
    db_name="esk-mdb" \
    creation_statements="{ \
        \"db\": \"$ESK_MDB_NAME\", \
        \"roles\": [ \
            { \
                \"role\": \"readWrite\", \
                \"db\": \"$ESK_MDB_NAME\" \
            } \
        ] \
    }" \
    default_ttl="10s" \
    max_ttl="20s"

# for redis
vault kv put secret/redis \
    app_user="$REDIS_APP_USER" \
    app_user_password="$REDIS_APP_PASSWORD" \
    admin_user="$REDIS_ADMIN_USER" \
    admin_user_password="$REDIS_ADMIN_PASSWORD" >/dev/null 2>&1

wait $VAULT_PID

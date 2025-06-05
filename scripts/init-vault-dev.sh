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

vault write database/config/esk-rdb \
    plugin_name=postgresql-database-plugin \
    allowed_roles="postgres-app-role" \
    connection_url="$ESK_RDB_CONN_URL" \
    username="$ESK_RDB_MASTER_USERNAME" \
    password="$ESK_RDB_MASTER_PASSWORD"

# for api-server
# vault write database/roles/postgres-app-role \
#     db_name="esk-rdb" \
#     creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; \
#     GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO \"{{name}}\";" \
#     default_ttl="1h" \
#     max_ttl="24h"

# for migration
vault write database/roles/postgres-app-role \
    db_name="esk-rdb" \
    creation_statements="
        CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';
        GRANT CONNECT ON DATABASE $ESK_RDB_NAME TO \"{{name}}\";
        GRANT USAGE, CREATE ON SCHEMA public TO \"{{name}}\";
        GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO \"{{name}}\";
        GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO \"{{name}}\";
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO \"{{name}}\";
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO \"{{name}}\";

        -- Transfer ownership of existing tables
        DO \$\$
        DECLARE
            r RECORD;
        BEGIN
            FOR r IN SELECT tablename FROM pg_tables WHERE schemaname = 'public' LOOP
                EXECUTE format('ALTER TABLE public.%I OWNER TO \"{{name}}\"', r.tablename);
            END LOOP;
        END
        \$\$;
    " \
    default_ttl="1h" \
    max_ttl="24h"

wait $VAULT_PID

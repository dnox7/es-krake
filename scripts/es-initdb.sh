#!bin/sh
set -e

# Wait for Elasticsearch to be ready
until curl -s "$ES_URL" >/dev/null; do
    echo "Waiting for Elasticsearch at $ES_URL..."
    sleep 2
done

echo "Elasticsearch is ready."

curl -X POST $ES_URL/_security/role/vault-admin \
    -u $ES_ADMIN:$ES_ADMIN_PASSWORD \
    -H "Content-Type: application/json" \
    -d '{"cluster": [ "manage_security" ]}'

curl -X POST $ES_URL/_security/user/$ES_VAULT_USER \
    -u $ES_ADMIN:$ES_ADMIN_PASSWORD \
    -H "Content-Type: application/json" \
    -d "{
        \"password\": \"$ES_VAULT_PASSWORD\",
        \"roles\": [\"vault-admin\"],
        \"full_name\": \"ESK Vault\",
        \"metadata\": {
            \"plugin_name\": \"Vault Plugin Database Elasticsearch\",
            \"plugin_url\": \"https://github.com/hashicorp/vault-plugin-database-elasticsearch\"
        }
    }"

#!/bin/bash
set -e

CONTAINER_NAME="vault-dev"
ENV_FILE_PATH="./.env"
SECRET_ID_FILE_PATH="./vault_secret"

ROLE_ID=$(docker exec "$CONTAINER_NAME" cat /vault/config/role_id_$SERVICE)
SECRET_ID=$(docker exec "$CONTAINER_NAME" cat /vault/config/secret_id_$SERVICE)

if [ -z "$ROLE_ID" ] || [ -z "$SECRET_ID" ]; then
    echo "failed to get role_id and secret_id"
    exit 1
fi

if [ ! -f "$ENV_FILE_PATH" ]; then
    touch "$ENV_FILE_PATH"
fi

if grep -q '^VAULT_ROLE_ID=' $ENV_FILE_PATH; then
    sed -i "s#^VAULT_ROLE_ID=.*#VAULT_ROLE_ID=$ROLE_ID#" "$ENV_FILE_PATH"
else
    echo "VAULT_ROLE_ID=$ROLE_ID" >>"$ENV_FILE_PATH"
fi

echo "$SECRET_ID" >"$SECRET_ID_FILE_PATH"

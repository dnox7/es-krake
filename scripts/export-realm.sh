#!/bin/bash

CONTAINER_NAME=keycloak
REALM_NAME=krake
EXPORT_DIR=/opt/keycloak/data/export
EXPORT_FILE=$REALM_NAME-realm.json
HOST_EXPORT_DIR=./deploy/compose

if [ ! -d "${HOST_EXPORT_DIR}" ] || [ ! -w "${HOST_EXPORT_DIR}" ]; then
    echo "Error: HOST_EXPORT_DIR '${HOST_EXPORT_DIR}' is not a valid writable directory."
    exit 1
fi

docker exec -it "${CONTAINER_NAME}" /opt/keycloak/bin/kc.sh export \
    --realm "${REALM_NAME}" \
    --dir "${EXPORT_DIR}"

if [ $? -ne 0 ]; then
    echo "Error: Failed to export realm."
    exit 1
fi

docker cp "${CONTAINER_NAME}:${EXPORT_DIR}/${EXPORT_FILE}" "${HOST_EXPORT_DIR}/"

if [ $? -ne 0 ]; then
    echo "Error: Failed to copy the JSON file from container to host."
    exit 1
fi

echo "Realm '${REALM_NAME}' export successful. File saved at: ${HOST_EXPORT_DIR}/${EXPORT_FILE}"

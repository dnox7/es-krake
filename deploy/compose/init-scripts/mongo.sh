#!/bin/bash

MONGO_USER="krake"
MONGO_PASSWORD="123456789"

# Prepare the user credentials for MongoDB
q_MONGO_USER=$(jq --arg v "$MONGO_USER" -n '$v')
q_MONGO_PASSWORD=$(jq --arg v "$MONGO_PASSWORD" -n '$v')

# Run MongoDB commands
mongo -u "$MONGO_INITDB_ROOT_USERNAME" -p "$MONGO_INITDB_ROOT_PASSWORD" admin <<EOF
    use es_krake;
    db.createUser({
        user: $q_MONGO_USER,
        pwd: $q_MONGO_PASSWORD,
        roles: ["readWrite"],
    });
EOF

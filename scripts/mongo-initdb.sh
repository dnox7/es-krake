#!/bin/bash

until mongosh -u "$MONGO_INITDB_ROOT_USERNAME" -p "$MONGO_INITDB_ROOT_PASSWORD" --eval "db.adminCommand('ping')" > /dev/null 2>&1; do
  echo "Waiting for MongoDB to be ready..."
  sleep 2
done

mongosh -u "$MONGO_INITDB_ROOT_USERNAME" -p "$MONGO_INITDB_ROOT_PASSWORD" admin <<EOF
    db.getSiblingDB("$MDB_NAME").createUser({
        user: "$MONGO_ADMIN_USERNAME",
        pwd: "$MONGO_ADMIN_PASSWORD",
        roles: [
            { role: "userAdmin", db: "$MDB_NAME" },
            { role: "readWrite", db: "$MDB_NAME" }
        ]
    });
EOF

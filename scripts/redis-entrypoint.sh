#!/bin/bash
set -e

if
    [ -z "$REDIS_ADMIN_USER" ] || [ -z "$REDIS_ADMIN_PASSWORD" ] || [ -z "$REDIS_APP_USER" ] || [ -z "$REDIS_APP_PASSWORD" ]
then
    echo "ERROR: missing environment variables"
    exit 1
fi

ACL_FILE="/usr/local/etc/redis/users.acl"
CONFIG_FILE="/usr/local/etc/redis/redis.conf"

cat <<EOF >"$ACL_FILE"
user default off
user $REDIS_APP_USER on >"$REDIS_APP_PASSWORD" allcommands ~* -@dangerous
user $REDIS_ADMIN_USER on >"$REDIS_ADMIN_PASSWORD" allcommands allkeys
EOF

exec redis-server "$CONFIG_FILE" "$@"

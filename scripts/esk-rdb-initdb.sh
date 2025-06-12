#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" <<-EOSQL
    CREATE ROLE esk_dev_migrator WITH LOGIN PASSWORD '${MIGRATOR_INITIAL_PASSWORD}';
    GRANT CONNECT ON DATABASE esk_dev_1 TO esk_dev_migrator;
    GRANT USAGE, CREATE ON SCHEMA public TO esk_dev_migrator;

    CREATE OR REPLACE FUNCTION __tmp_create_user() returns void as \$\$
    BEGIN
      IF NOT EXISTS (
              SELECT                       -- SELECT list can stay empty for this
              FROM   pg_catalog.pg_user
              WHERE  usename = '${EXPORTER_USER}') THEN
        CREATE USER ${EXPORTER_USER};
      END IF;
    END;
    \$\$ language plpgsql;

    SELECT __tmp_create_user();
    DROP FUNCTION __tmp_create_user();
    
    ALTER USER ${EXPORTER_USER} WITH PASSWORD '${EXPORTER_PASSWORD}';
    ALTER USER ${EXPORTER_USER} SET SEARCH_PATH TO ${EXPORTER_USER}, pg_catalog;
    GRANT CONNECT ON DATABASE ${POSTGRES_DB} TO ${EXPORTER_USER};
    GRANT ALL ON FUNCTION pg_catalog.pg_ls_waldir() TO ${EXPORTER_USER};
EOSQL

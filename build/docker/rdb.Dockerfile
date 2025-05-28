ARG BASE_MIGRATION
FROM ${BASE_MIGRATION} as migration_builder

FROM postgres:17-alpine3.21

COPY --from=migration_builder --chown=root:root /migration-tool  /migration-tool
COPY ./migrations ./migrations
COPY ./scripts/run-migration.sh /usr/local/bin/run-migration.sh
RUN chmod +x /usr/local/bin/run-migration.sh

path "database/creds/postgres-app-role" {
    capabilities = ["read"]
}

path "database/creds/mongo-app-role" {
    capabilities = ["read"]
}

path "database/creds/esdb-app-role" {
    capabilities = ["read"]
}

path "secret/data/redis" {
    capabilities = ["read"]
}

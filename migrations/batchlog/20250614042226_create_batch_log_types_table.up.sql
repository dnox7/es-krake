CREATE TABLE IF NOT EXISTS batch_log_types (
    id serial PRIMARY KEY,
    type varchar(255) NOT NULL UNIQUE,
    name varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
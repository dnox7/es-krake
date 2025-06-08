CREATE TABLE IF NOT EXISTS access_requirements (
    id serial PRIMARY KEY,
    code varchar(25) NOT NULL UNIQUE,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

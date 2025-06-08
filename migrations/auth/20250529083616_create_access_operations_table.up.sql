CREATE TABLE IF NOT EXISTS access_operations (
    id serial PRIMARY KEY,
    name varchar(50) NOT NULL,
    code varchar(50) NOT NULL UNIQUE,
    description varchar(255),
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

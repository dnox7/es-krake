CREATE TABLE IF NOT EXISTS role_types (
    id smallserial PRIMARY KEY,
    name varchar(50) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

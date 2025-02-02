CREATE TABLE IF NOT EXISTS attribute_types (
    id smallserial PRIMARY KEY,
    name varchar(20) UNIQUE NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS otps (
    id serial PRIMARY KEY,
    kc_user_id varchar(255) NOT NULL,
    token varchar(255) NOT NULL,
    expired_at datetime NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
);
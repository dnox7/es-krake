CREATE TABLE IF NOT EXISTS login_histories (
    id serial PRIMARY KEY,
    kc_user_id varchar(255) NOT NULL,
    logged_in_at timestamp NOT NULL,
    logged_device varchar(255) NOT NULL,
    logged_ip_address varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
);
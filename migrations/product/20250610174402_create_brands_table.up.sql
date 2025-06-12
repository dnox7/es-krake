CREATE TABLE IF NOT EXISTS brands (
    id serial PRIMARY KEY, 
    name varchar(50) NOT NULL,
    is_active boolean NOT NULL DEFAULT TRUE,
    description text,
    thumbnail_path varchar(255),
    website_path varchar(255),
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

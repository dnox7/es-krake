CREATE TABLE IF NOT EXISTS categories (
    id serial PRIMARY KEY,
    name varchar(20) NOT NULL,
    slug varchar(50) NOT NULL UNIQUE,
    description varchar(50),
    is_published boolean NOT NULL,
    display_order integer NOT NULL,
    meta_description text,
    thumbnail_url text,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
)

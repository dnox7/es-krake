CREATE TABLE IF NOT EXISTS products (
    id serial PRIMARY KEY,
    name varchar(50) NOT NULL,
    sku varchar(20) NOT NULL,
    description varchar(50) NOT NULL,
    price decimal(12, 5) NOT NULL,
    has_options boolean NOT NULL DEFAULT FALSE,
    is_allowed_to_order boolean NOT NULL DEFAULT FALSE,
    is_featured boolean NOT NULL DEFAULT FALSE,
    is_visible_individually boolean NOT NULL DEFAULT FALSE,
    stock_tracking_enabled boolean NOT NULL DEFAULT FALSE,
    stock_quantity  bigint NOT NULL,
    tax_class_id smallint NOT NULL,
    meta_title text,
    meta_keyword text,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

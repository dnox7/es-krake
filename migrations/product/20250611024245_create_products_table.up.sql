CREATE TABLE IF NOT EXISTS products (
    id serial PRIMARY KEY, 
    name varchar(255) NOT NULL,
    sku varchar(50) NOT NULL,
    description varchar(255),
    price double precision NOT NULL,
    has_options boolean NOT NULL DEFAULT FALSE,
    is_allowed_to_order boolean NOT NULL DEFAULT FALSE,
    is_published boolean NOT NULL DEFAULT FALSE,
    stock_tracking_enabled boolean NOT NULL DEFAULT TRUE,
    stock_quantity bigint NOT NULL DEFAULT 0,
    thumbnail_url varchar(255),
    brand_id bigint NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_brand
            FOREIGN KEY (brand_id) REFERENCES brands(id)
                ON DELETE RESTRICT
                ON UPDATE RESTRICT
);

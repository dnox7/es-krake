CREATE TABLE IF NOT EXISTS product_categories (
    id serial PRIMARY KEY, 
    product_id bigint NOT NULL,
    category_id bigint NOT NULL,
    display_order integer NOT NULL,
    is_featured_product boolean NOT NULL DEFAULT FALSE,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at integer NOT NULL DEFAULT 0,
    CONSTRAINT fk_category
            FOREIGN KEY (category_id) REFERENCES categories(id)
                ON DELETE RESTRICT
                ON UPDATE RESTRICT,
    CONSTRAINT fk_product
            FOREIGN KEY (product_id) REFERENCES products(id)
                ON DELETE RESTRICT
                ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS product_options (
    id serial PRIMARY KEY, 
    product_id bigint NOT NULL,
    name varchar(50) NOT NULL,
    description varchar(255),
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at integer NOT NULL DEFAULT 0,
    CONSTRAINT fk_product
            FOREIGN KEY (product_id) REFERENCES products(id)
                ON DELETE RESTRICT
                ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS option_attribute_values (
    id serial PRIMARY KEY, 
    attribute_id bigint NOT NULL,
    product_option_id bigint NOT NULL,
    value varchar(50) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at integer NOT NULL DEFAULT 0,
    CONSTRAINT fk_attribute
            FOREIGN KEY (attribute_id) REFERENCES attributes(id)
                ON DELETE RESTRICT
                ON UPDATE RESTRICT,
    CONSTRAINT fk_product_option
            FOREIGN KEY (product_option_id) REFERENCES product_options(id)
                ON DELETE RESTRICT
                ON UPDATE RESTRICT
);

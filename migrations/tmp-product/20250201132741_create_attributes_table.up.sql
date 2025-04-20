CREATE TABLE IF NOT EXISTS attributes (
    id serial PRIMARY KEY, 
    name varchar(10) NOT NULL,
    description text,
    attribute_type_id smallint NOT NULL,
    display_order integer NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_attribute_type
        FOREIGN KEY (attribute_type_id) REFERENCES attribute_types(id)
            ON DELETE RESTRICT
            ON UPDATE RESTRICT
);

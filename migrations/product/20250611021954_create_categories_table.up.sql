CREATE TABLE IF NOT EXISTS categories (
    id serial PRIMARY KEY, 
    name varchar(50) NOT NULL,
    slug varchar(100) NOT NULL,
    is_active bool NOT NULL DEFAULT TRUE,
    description varchar(255),
    is_published boolean NOT NULL DEFAULT FALSE,
    display_order smallint NOT NULL,
    meta_description text,
    thumbnail_url varchar(255),
    parent_id integer,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_category_parent
            FOREIGN KEY (parent_id) REFERENCES categories(id)
                ON DELETE SET NULL
                ON UPDATE RESTRICT
);

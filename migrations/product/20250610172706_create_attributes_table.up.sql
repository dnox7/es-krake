CREATE TABLE IF NOT EXISTS attributes (
    id serial PRIMARY KEY, 
    name varchar(10) NOT NULL,
    description text,
    display_order integer NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

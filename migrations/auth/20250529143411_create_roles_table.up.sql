CREATE TABLE IF NOT EXISTS roles (
    id serial PRIMARY KEY,
    name varchar(50) NOT NULL,
    display_order integer NOT NULL,
    role_type_id smallint NOT NULL DEFAULT 1,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at integer NOT NULL DEFAULT 0,
    CONSTRAINT fk_role_type
            FOREIGN KEY (role_type_id) REFERENCES role_types(id)
                ON DELETE RESTRICT
                ON UPDATE RESTRICT
);

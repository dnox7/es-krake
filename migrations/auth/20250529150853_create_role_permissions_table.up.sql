CREATE TABLE IF NOT EXISTS role_permissions (
    id serial PRIMARY KEY,
    role_id integer NOT NULL,
    permission_id integer NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at integer NOT NULL DEFAULT 0,
    CONSTRAINT fk_role
            FOREIGN KEY (role_id) REFERENCES roles(id)
                ON DELETE RESTRICT
                ON UPDATE RESTRICT
    CONSTRAINT fk_permission
            FOREIGN KEY (permission_id) REFERENCES permission_id(id)
                ON DELETE RESTRICT
                ON UPDATE RESTRICT
);

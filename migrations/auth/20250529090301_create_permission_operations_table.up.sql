CREATE TABLE IF NOT EXISTS permission_operations (
    id serial PRIMARY KEY,
    access_operation_id integer NOT NULL,
    permission_id integer NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_access_operation
            FOREIGN KEY (access_operation_id) REFERENCES access_operations(id)
                ON DELETE RESTRICT
                ON UPDATE RESTRICT,
    CONSTRAINT fk_permission
            FOREIGN KEY (permission_id) REFERENCES permissions(id)
                ON DELETE RESTRICT
                ON UPDATE RESTRICT
);

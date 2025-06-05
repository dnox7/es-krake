CREATE TABLE IF NOT EXISTS access_requirement_operations (
    id serial PRIMARY KEY,
    access_operation_id integer NOT NULL,
    access_requirement_id integer NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_access_operation
            FOREIGN KEY (access_operation_id) REFERENCES access_operations(id)
                ON DELETE RESTRICT
                ON UPDATE RESTRICT,
    CONSTRAINT fk_access_requirement
            FOREIGN KEY (access_requirement_id) REFERENCES access_requirements(id)
                ON DELETE RESTRICT
                ON UPDATE RESTRICT
);

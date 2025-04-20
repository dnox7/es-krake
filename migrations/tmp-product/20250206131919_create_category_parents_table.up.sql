CREATE TABLE IF NOT EXISTS category_parents (
    child_id integer PRIMARY KEY,
    parent_id integer PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_child FOREIGN KEY (child_id)
        REFERENCES categories(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_parent FOREIGN KEY (parent_id)
        REFERENCES categories(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
);

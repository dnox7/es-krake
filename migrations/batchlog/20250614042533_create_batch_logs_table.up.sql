    CREATE TABLE IF NOT EXISTS batch_log_types (
    id serial PRIMARY KEY,
    batch_log_type_id integer NOT NULL,
    event_id integer NOT NULL,
    started_at timestamp NOT NULL,
    ended_at timestamp,
    arguments jsonb,
    success boolean,
    error_message text,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
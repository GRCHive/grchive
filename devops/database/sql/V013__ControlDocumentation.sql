CREATE TABLE process_flow_control_documentation_categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL UNIQUE,
    description TEXT,
    control_id BIGINT NOT NULL REFERENCES process_flow_controls(id) ON DELETE CASCADE,
    UNIQUE(name, control_id)
);

CREATE TABLE process_flow_control_documentation_file (
    id BIGSERIAL PRIMARY KEY,
    storage_id VARCHAR(200) NOT NULL,
    storage_name TEXT NOT NULL,
    relevant_time TIMESTAMPTZ NOT NULL,
    upload_time TIMESTAMPTZ NOT NULL,
    category_id BIGINT NOT NULL REFERENCES process_flow_control_documentation_categories(id) ON DELETE CASCADE
);

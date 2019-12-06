CREATE TABLE process_flow_control_documentation_categories (
    id BIGSERIAL,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    org_id INTEGER NOT NULL REFERENCES organizations(id),
    PRIMARY KEY(id, org_id),
    UNIQUE(name, org_id)
);

CREATE TABLE process_flow_control_documentation_file (
    id BIGSERIAL,
    bucket_id VARCHAR(200),
    storage_id VARCHAR(200),
    storage_name TEXT NOT NULL,
    relevant_time TIMESTAMPTZ NOT NULL,
    upload_time TIMESTAMPTZ NOT NULL,
    category_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id),
    PRIMARY KEY(id, category_id, org_id),
    CONSTRAINT cat_org_fkey
        FOREIGN KEY(category_id, org_id)
        REFERENCES process_flow_control_documentation_categories(id, org_id)
        ON DELETE CASCADE
);

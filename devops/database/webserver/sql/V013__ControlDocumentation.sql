CREATE TABLE process_flow_control_documentation_categories (
    id BIGSERIAL,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    PRIMARY KEY(id, org_id),
    UNIQUE(name, org_id)
);

CREATE TABLE file_metadata (
    id BIGSERIAL,
    relevant_time TIMESTAMPTZ NOT NULL,
    category_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    alt_name VARCHAR(256),
    description TEXT,
    PRIMARY KEY(id, org_id),
    CONSTRAINT cat_org_fkey
        FOREIGN KEY(category_id, org_id)
        REFERENCES process_flow_control_documentation_categories(id, org_id)
        ON DELETE CASCADE
);

CREATE INDEX file_metadata_cat_index ON file_metadata(category_id);

CREATE TABLE file_storage (
    id BIGSERIAL,
    metadata_id BIGINT NOT NULL,
    storage_name TEXT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    bucket_id VARCHAR(200) NOT NULL,
    storage_id VARCHAR(200) NOT NULL,
    upload_time TIMESTAMPTZ NOT NULL,
    upload_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    PRIMARY KEY(id, org_id),
    FOREIGN KEY(metadata_id, org_id) REFERENCES file_metadata(id, org_id) ON DELETE CASCADE,
    FOREIGN KEY(upload_user_id, org_id) REFERENCES user_orgs(user_id, org_id) ON DELETE RESTRICT,
    UNIQUE(bucket_id, storage_id)
);

CREATE INDEX file_storage_metadata_index ON file_storage(metadata_id);

CREATE TABLE database_refresh (
    id BIGSERIAL,
    db_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    refresh_time TIMESTAMPTZ,
    refresh_success BOOLEAN,
    refresh_errors TEXT,
    PRIMARY KEY(org_id, id),
    FOREIGN KEY(db_id, org_id) REFERENCES database_resources(id, org_id) ON DELETE CASCADE
);
CREATE INDEX ON database_refresh(org_id, db_id);

CREATE TABLE database_schemas (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    refresh_id BIGINT NOT NULL,
    schema_name TEXT,
    PRIMARY KEY(org_id, id),
    FOREIGN KEY(org_id, refresh_id) REFERENCES database_refresh(org_id, id) ON DELETE CASCADE
);
CREATE INDEX ON database_schemas(org_id, refresh_id);

CREATE TABLE database_tables (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    schema_id BIGINT NOT NULL,
    table_name TEXT,
    PRIMARY KEY(org_id, id),
    FOREIGN KEY(org_id, schema_id) REFERENCES database_schemas(org_id, id) ON DELETE CASCADE
);
CREATE INDEX ON database_tables(org_id, schema_id);

CREATE TABLE database_columns (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    table_id BIGINT NOT NULL,
    column_name TEXT,
    column_type TEXT,
    PRIMARY KEY(org_id, id),
    FOREIGN KEY(org_id, table_id) REFERENCES database_tables(org_id, id) ON DELETE CASCADE
);
CREATE INDEX ON database_columns(org_id, table_id);

CREATE TABLE database_sql_metadata (
    id BIGSERIAL,
    db_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    PRIMARY KEY(org_id, id),
    FOREIGN KEY(db_id, org_id) REFERENCES database_resources(id, org_id) ON DELETE CASCADE
);
CREATE INDEX ON database_sql_metadata(org_id, db_id);

CREATE TABLE database_sql_queries (
    id BIGSERIAL,
    metadata_id BIGINT NOT NULL,
    version_number INTEGER NOT NULL,
    upload_time TIMESTAMPTZ NOT NULL,
    upload_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    PRIMARY KEY(org_id, id),
    FOREIGN KEY(upload_user_id, org_id) REFERENCES user_orgs(user_id, org_id) ON DELETE RESTRICT,
    FOREIGN KEY(org_id, metadata_id) REFERENCES database_sql_metadata(org_id, id) ON DELETE CASCADE
);
CREATE INDEX ON database_sql_queries(org_id, metadata_id);

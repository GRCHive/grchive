CREATE TABLE data_source_options (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL UNIQUE,
    kotlin_class TEXT NOT NULL
);

INSERT INTO data_source_options (name, kotlin_class)
VALUES
    ('Root.GRCHive', 'grchive.core.data.sources.GrchiveDataSource'),
    ('Root.Database.PostgreSQL', 'grchive.core.data.sources.databases.PostgresDataSource');

CREATE TABLE client_data_source_link (
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    data_id BIGINT NOT NULL UNIQUE,
    source_id BIGINT NOT NULL REFERENCES data_source_options(id) ON DELETE RESTRICT,
    source_target JSONB NOT NULL,
    FOREIGN KEY(data_id, org_id) REFERENCES client_data(id, org_id) ON DELETE CASCADE
);

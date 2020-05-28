CREATE TABLE shell_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL
);

INSERT INTO shell_types (id, name)
VALUES
    (1, 'Bash'),
    (2, 'PowerShell');

CREATE TABLE shell_scripts (
    id BIGSERIAL PRIMARY KEY,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    type_id INTEGER NOT NULL REFERENCES shell_types(id) ON DELETE RESTRICT,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    bucket_id VARCHAR(200) NOT NULL,
    storage_id VARCHAR(200) NOT NULL,
    UNIQUE(org_id, name)
);

CREATE INDEX ON shell_scripts(org_id, type_id);

CREATE TABLE shell_script_versions (
    shell_id BIGINT NOT NULL REFERENCES shell_scripts(id) ON DELETE CASCADE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    upload_time TIMESTAMPTZ NOT NULL,
    upload_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    gcs_generation BIGINT NOT NULL,
    FOREIGN KEY(upload_user_id, org_id) REFERENCES user_orgs(user_id, org_id) ON DELETE RESTRICT
);

CREATE INDEX ON shell_script_versions(org_id, shell_id);
CREATE INDEX ON shell_script_versions(upload_user_id);

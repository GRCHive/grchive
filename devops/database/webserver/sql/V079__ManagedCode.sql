CREATE TABLE managed_kotlin_code (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    git_hash VARCHAR(40) NOT NULL,
    action_time TIMESTAMPTZ NOT NULL,
    PRIMARY KEY(id, org_id)
);

CREATE INDEX ON managed_kotlin_code(git_hash);
DROP TABLE client_data_versions;

CREATE TABLE code_to_client_data_link (
    code_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    data_id BIGINT NOT NULL,
    FOREIGN KEY(code_id, org_id) REFERENCES managed_kotlin_code(id, org_id),
    FOREIGN KEY(data_id, org_id) REFERENCES client_data(id, org_id),
    UNIQUE(code_id, org_id, data_id)
);
CREATE INDEX ON code_to_client_data_link(data_id);

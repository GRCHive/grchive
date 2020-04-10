CREATE TABLE client_scripts (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    UNIQUE(name, org_id),
    PRIMARY KEY(id, org_id)
);
CREATE INDEX ON client_scripts(org_id);

CREATE TABLE code_to_client_scripts_link (
    id BIGSERIAL PRIMARY KEY,
    code_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    script_id BIGINT NOT NULL,
    FOREIGN KEY(code_id, org_id) REFERENCES managed_code(id, org_id) ON DELETE CASCADE,
    FOREIGN KEY(script_id, org_id) REFERENCES client_scripts(id, org_id) ON DELETE CASCADE,
    UNIQUE(code_id, org_id, script_id)
);
CREATE INDEX ON code_to_client_scripts_link(script_id);

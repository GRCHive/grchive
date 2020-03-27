CREATE TABLE client_data (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    UNIQUE(name, org_id),
    PRIMARY KEY(id, org_id)
);
CREATE INDEX ON client_data(org_id);

CREATE TABLE client_data_versions (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    data_id BIGINT NOT NULL UNIQUE,
    version INTEGER NOT NULL,
    kotlin TEXT,
    FOREIGN KEY(data_id, org_id) REFERENCES client_data(id, org_id) ON DELETE CASCADE,
    PRIMARY KEY(id, org_id)
);

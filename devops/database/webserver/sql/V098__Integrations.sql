CREATE TABLE supported_integrations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL
);

INSERT INTO supported_integrations (id, name)
VALUES (1, 'SAP ERP');

CREATE TABLE integrations (
    id BIGSERIAL PRIMARY KEY,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    type INTEGER NOT NULL REFERENCES supported_integrations(id) ON DELETE RESTRICT,
    name VARCHAR(256) NOT NULL,
    description TEXT
);

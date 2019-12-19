CREATE TABLE container_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256),
    description TEXT
);

CREATE TABLE containers (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    type_id INTEGER NOT NULL REFERENCES container_types(id) ON DELETE RESTRICT,
    os_id BIGINT,
    image_name VARCHAR(256) NOT NULL,
    registry TEXT,
    description TEXT,
    FOREIGN KEY(os_id, org_id) REFERENCES operating_systems(id, org_id) ON DELETE NO ACTION
);

CREATE TABLE container_orchestration_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256),
    description TEXT
);

CREATE TABLE container_orchestration (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    type_id INTEGER NOT NULL REFERENCES container_orchestration_types(id) ON DELETE RESTRICT
);

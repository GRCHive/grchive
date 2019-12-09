CREATE TABLE supported_databases (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL
);

CREATE TABLE database_resources (
    id BIGSERIAL,
    name TEXT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    type_id INTEGER NOT NULL REFERENCES supported_databases(id) ON DELETE NO ACTION,
    other_type TEXT,
    version TEXT NOT NULL,
    PRIMARY KEY(id, org_id)
);

CREATE TABLE database_connection_info (
    id BIGSERIAL,
    db_id BIGINT NOT NULL UNIQUE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    connection_string TEXT NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    salt TEXT NOT NULL,
    PRIMARY KEY(id, db_id, org_id),
    FOREIGN KEY(db_id, org_id) REFERENCES database_resources(id, org_id) ON DELETE CASCADE
);

CREATE TABLE database_system_link (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    db_id BIGINT NOT NULL,
    system_id BIGINT NOT NULL,
    FOREIGN KEY(db_id, org_id) REFERENCES database_resources(id, org_id) ON DELETE CASCADE,
    FOREIGN KEY(system_id, org_id) REFERENCES systems(id, org_id) ON DELETE CASCADE,
    PRIMARY KEY(id, org_id)
);

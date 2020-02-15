CREATE TABLE database_functions (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    schema_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    src TEXT NOT NULL,
    PRIMARY KEY(org_id, id),
    FOREIGN KEY(org_id, schema_id) REFERENCES database_schemas(org_id, id) ON DELETE CASCADE
);
CREATE INDEX ON database_functions(org_id);

CREATE TABLE database_function_returns (
    fn_id BIGINT,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    ret_type TEXT,
    FOREIGN KEY(org_id, fn_id) REFERENCES database_functions(org_id, id) ON DELETE CASCADE
);

CREATE INDEX ON database_function_returns(org_id, fn_id);

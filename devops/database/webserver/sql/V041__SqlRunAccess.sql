CREATE TABLE database_sql_query_requests (
    id BIGSERIAL,
    query_id BIGINT NOT NULL,
    upload_time TIMESTAMPTZ NOT NULL,
    upload_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    PRIMARY KEY(org_id, id),
    FOREIGN KEY(upload_user_id, org_id) REFERENCES user_orgs(user_id, org_id) ON DELETE RESTRICT,
    FOREIGN KEY(org_id, query_id) REFERENCES database_sql_queries(org_id, id) ON DELETE CASCADE,
    UNIQUE(org_id, name)
);
CREATE INDEX ON database_sql_query_requests(org_id, query_id);

CREATE TABLE resource_db_sql_requests_access (
    FOREIGN KEY(org_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE,
    FOREIGN KEY(role_id, org_id)
        REFERENCES organization_available_roles(id, org_id)
        ON DELETE CASCADE
) INHERITS (_base_resource_access);

CREATE UNIQUE INDEX ON resource_db_sql_requests_access(role_id);

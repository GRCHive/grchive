CREATE TABLE database_sql_query_requests_approvals (
    request_id BIGINT NOT NULL UNIQUE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    response_time TIMESTAMPTZ NOT NULL,
    responder_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    response BOOLEAN NOT NULL,
    reason TEXT,
    FOREIGN KEY(org_id, request_id) REFERENCES database_sql_query_requests(org_id, id) ON DELETE CASCADE
);

CREATE INDEX ON database_sql_query_requests_approvals(org_id, request_id);
CREATE INDEX ON database_sql_query_requests_approvals(org_id, response);

CREATE TABLE database_sql_query_run_codes(
    request_id BIGINT NOT NULL UNIQUE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    expiration_time TIMESTAMPTZ NOT NULL,
    used_time TIMESTAMPTZ,
    hashed_code VARCHAR(256) NOT NULL,
    salt TEXT NOT NULL,
    FOREIGN KEY(org_id, request_id) REFERENCES database_sql_query_requests(org_id, id) ON DELETE CASCADE
);

CREATE INDEX ON database_sql_query_run_codes(org_id, request_id, hashed_code);

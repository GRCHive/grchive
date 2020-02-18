CREATE TABLE sql_request_comments (
    sql_request_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    comment_id BIGINT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY(org_id, sql_request_id) REFERENCES database_sql_query_requests(org_id, id) ON DELETE CASCADE
);

CREATE INDEX ON sql_request_comments(sql_request_id, org_id);

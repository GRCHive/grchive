CREATE TABLE access_requests (
    id VARCHAR(36) PRIMARY KEY,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    control_id INTEGER NOT NULL REFERENCES process_flow_controls(id) ON DELETE CASCADE,
    from_email VARCHAR(320) NOT NULL,
    request_time TIMESTAMPTZ NOT NULL,
    granted_time TIMESTAMPTZ,
    grant_expiration_time TIMESTAMPTZ
);

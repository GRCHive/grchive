CREATE TABLE audit_resource_modifications (
    event_id BIGINT REFERENCES global_audit_event_history(id) ON DELETE CASCADE,
    data JSONB
);

CREATE UNIQUE INDEX ON audit_resource_modifications(event_id);

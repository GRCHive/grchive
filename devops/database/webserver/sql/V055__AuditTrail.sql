CREATE TABLE global_audit_event_history (
    id BIGSERIAL PRIMARY KEY,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    resource_type VARCHAR(64) NOT NULL,
    resource_id BIGINT NOT NULL,
    resource_extra_data JSONB NOT NULL,
    action VARCHAR(64) NOT NULL,
    performed_at TIMESTAMPTZ NOT NULL,
    pgrole_id oid NOT NULL
);

CREATE INDEX ON global_audit_event_history(resource_type, resource_id, org_id, performed_at, action);

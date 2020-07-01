CREATE TABLE org_pbc_notification_record (
    cadence_id BIGINT NOT NULL REFERENCES org_pbc_notification_cadence_settings(id) ON DELETE CASCADE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    request_id BIGINT NOT NULL,
    FOREIGN KEY(request_id, org_id) REFERENCES document_requests(id, org_id) ON DELETE CASCADE
);

CREATE INDEX ON org_pbc_notification_record(org_id);

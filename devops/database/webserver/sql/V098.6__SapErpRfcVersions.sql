CREATE TABLE sap_erp_rfc_versions (
    id BIGSERIAL PRIMARY KEY,
    rfc_id BIGINT NOT NULL REFERENCES sap_erp_rfc(id) ON DELETE CASCADE,
    created_time TIMESTAMPTZ NOT NULL,
    finished_time TIMESTAMPTZ,
    data JSONB
);

CREATE INDEX ON sap_erp_rfc_versions(rfc_id, id);

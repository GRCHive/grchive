CREATE TABLE sap_erp_rfc (
    id BIGSERIAL PRIMARY KEY,
    integration_id BIGINT NOT NULL REFERENCES integrations(id) ON DELETE CASCADE,
    function_name VARCHAR(512) NOT NULL,
    UNIQUE(function_name, integration_id)
);
CREATE INDEX ON sap_erp_rfc(integration_id);

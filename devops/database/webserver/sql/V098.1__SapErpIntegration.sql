CREATE TABLE sap_erp_integration_info (
    integration_id BIGINT NOT NULL REFERENCES integrations(id) ON DELETE CASCADE,
    client TEXT NOT NULL,
    sysnr TEXT NOT NULL,
    host TEXT NOT NULL,
    real_hostname TEXT,
    username TEXT NOT NULL,
    password TEXT NOT NULL
);

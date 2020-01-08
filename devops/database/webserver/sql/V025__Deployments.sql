CREATE TABLE deployments (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    deployment_type INTEGER NOT NULL,
    PRIMARY KEY (id, org_id)
);

CREATE TABLE vendor_deployments (
    deployment_id BIGINT,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    vendor_name TEXT,
    vendor_product TEXT,
    FOREIGN KEY (deployment_id, org_id)
        REFERENCES deployments(id, org_id)
        ON DELETE CASCADE,
    UNIQUE(deployment_id, org_id)
);

CREATE TABLE vendor_soc_reports (
    deployment_id BIGINT,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    soc_report_file_id BIGINT NOT NULL,
    soc_report_cat_id BIGINT NOT NULL,
    FOREIGN KEY (deployment_id, org_id)
        REFERENCES deployments(id, org_id)
        ON DELETE CASCADE,
    FOREIGN KEY (soc_report_file_id, soc_report_cat_id, org_id)
        REFERENCES process_flow_control_documentation_file(id, category_id, org_id)
        ON DELETE NO ACTION,
    UNIQUE(deployment_id, org_id, soc_report_file_id)
);

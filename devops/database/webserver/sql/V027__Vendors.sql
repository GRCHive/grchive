CREATE TABLE vendors (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    url TEXT,
    PRIMARY KEY(id, org_id),
    UNIQUE(org_id, name)
);

CREATE TABLE vendor_products (
    id BIGSERIAL,
    vendor_id BIGINT NOT NULL,
    product_name VARCHAR(256) NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    description TEXT,
    url TEXT,
    PRIMARY KEY(id, org_id),
    UNIQUE(vendor_id, product_name, org_id),
    FOREIGN KEY (vendor_id, org_id)
        REFERENCES vendors(id, org_id)
        ON DELETE CASCADE
);

CREATE TABLE vendor_documentation_category_link (
    vendor_id BIGINT NOT NULL UNIQUE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    doc_cat_id BIGINT NOT NULL UNIQUE,
    FOREIGN KEY (vendor_id, org_id)
        REFERENCES vendors(id, org_id)
        ON DELETE CASCADE,
    FOREIGN KEY (doc_cat_id, org_id)
        REFERENCES process_flow_control_documentation_categories(id, org_id)
        ON DELETE CASCADE
);

CREATE TABLE vendor_product_soc_reports (
    product_id BIGINT,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    file_id BIGINT NOT NULL,
    cat_id BIGINT NOT NULL,
    FOREIGN KEY (product_id, org_id)
        REFERENCES vendor_products(id, org_id)
        ON DELETE CASCADE,
    FOREIGN KEY (file_id, cat_id, org_id)
        REFERENCES process_flow_control_documentation_file(id, category_id, org_id)
        ON DELETE CASCADE,
    UNIQUE(product_id, file_id)
);

CREATE TABLE vendor_deployments (
    deployment_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    vendor_product_id BIGINT NOT NULL,
    FOREIGN KEY (deployment_id, org_id)
        REFERENCES deployments(id, org_id)
        ON DELETE CASCADE,
    FOREIGN KEY (vendor_product_id, org_id)
        REFERENCES vendor_products(id, org_id)
        ON DELETE CASCADE,
    UNIQUE(deployment_id, vendor_product_id)
);

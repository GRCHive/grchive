CREATE TABLE vendor_soc_request_link (
    vendor_product_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    request_id BIGINT NOT NULL,
    FOREIGN KEY (vendor_product_id, org_id)
        REFERENCES vendor_products(id, org_id)
        ON DELETE CASCADE,
    FOREIGN KEY (request_id, org_id)
        REFERENCES document_requests(id, org_id)
        ON DELETE CASCADE,
    UNIQUE(vendor_product_id, request_id)
);

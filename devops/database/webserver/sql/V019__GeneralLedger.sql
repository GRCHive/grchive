CREATE TABLE general_ledger_categories (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id),
    parent_category_id BIGINT,
    name TEXT NOT NULL,
    description TEXT,
    UNIQUE(org_id, name),
    PRIMARY KEY(id, org_id),
    FOREIGN KEY(parent_category_id, org_id) REFERENCES general_ledger_categories(id, org_id) ON DELETE CASCADE
);

CREATE TABLE general_ledger_accounts (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id),
    parent_category_id BIGINT NOT NULL,
    account_identifier TEXT NOT NULL,
    account_name TEXT NOT NULL,
    account_description TEXT,
    financially_relevant BOOLEAN NOT NULL,
    PRIMARY KEY(id, org_id),
    FOREIGN KEY(parent_category_id, org_id) REFERENCES general_ledger_categories(id, org_id) ON DELETE CASCADE,
    UNIQUE(org_id, account_name),
    UNIQUE(org_id, account_identifier)
);

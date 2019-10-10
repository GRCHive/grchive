CREATE TABLE general_ledgers (
    id SERIAL PRIMARY KEY,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE RESTRICT
);

CREATE TABLE general_ledger_category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    ledger_id INTEGER NOT NULL REFERENCES general_ledgers(id) ON DELETE RESTRICT,
    parent_category_id INTEGER REFERENCES general_ledger_category(id) ON DELETE RESTRICT
);

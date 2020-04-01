CREATE TABLE available_features (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL UNIQUE
);

INSERT INTO available_features (id, name)
VALUES (1, 'Automation');

CREATE TABLE organization_enabled_features (
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    feature_id BIGINT NOT NULL REFERENCES available_features(id) ON DELETE CASCADE,
    requested TIMESTAMPTZ NOT NULL,
    fulfilled TIMESTAMPTZ,
    UNIQUE(org_id, feature_id)
);

CREATE INDEX ON organization_enabled_features (org_id);

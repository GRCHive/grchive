CREATE TABLE managed_code_drone_ci (
    code_id BIGINT NOT NULL UNIQUE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    commit_hash VARCHAR(40) NOT NULL UNIQUE,
    time_start TIMESTAMPTZ NOT NULL,
    time_end TIMESTAMPTZ,
    success BOOLEAN DEFAULT false,
    logs TEXT,
    jar TEXT,
    FOREIGN KEY(code_id, org_id) REFERENCES managed_code(id, org_id) ON DELETE CASCADE
);

CREATE INDEX ON managed_code_drone_ci(commit_hash);

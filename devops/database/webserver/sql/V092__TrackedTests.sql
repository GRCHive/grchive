CREATE TABLE test_sources (
    id BIGSERIAL PRIMARY KEY,
    run_id BIGINT NOT NULL REFERENCES script_runs(id) ON DELETE CASCADE,
    data_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    src TEXT NOT NULL,
    FOREIGN KEY(data_id, org_id) REFERENCES client_data(id, org_id) ON DELETE CASCADE
);
CREATE INDEX ON test_sources(run_id);

CREATE TABLE test_data (
    id BIGSERIAL PRIMARY KEY,
    source_id BIGINT NOT NULL REFERENCES test_sources(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);
CREATE INDEX ON test_data(source_id);

CREATE TABLE test_tests (
    id BIGSERIAL PRIMARY KEY,
    data_a_id BIGINT NOT NULL REFERENCES test_data(id) ON DELETE CASCADE,
    data_b_id BIGINT NOT NULL REFERENCES test_data(id) ON DELETE CASCADE,
    ok BOOLEAN NOT NULL,
    action VARCHAR(64) NOT NULL
);
CREATE INDEX ON test_tests(ok);

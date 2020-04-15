DROP TABLE script_runs;
CREATE TABLE script_runs (
    id BIGSERIAL PRIMARY KEY,
    link_id BIGINT NOT NULL REFERENCES code_to_client_scripts_link(id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    build_log TEXT,
    run_log TEXT
);

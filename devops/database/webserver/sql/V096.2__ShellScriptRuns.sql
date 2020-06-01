CREATE TABLE shell_script_runs (
    id BIGSERIAL PRIMARY KEY,
    script_version_id BIGINT NOT NULL REFERENCES shell_script_versions(id) ON DELETE CASCADE,
    run_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    create_time TIMESTAMPTZ NOT NULL,
    run_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ
);

CREATE TABLE shell_script_run_servers (
    run_id BIGINT NOT NULL REFERENCES shell_script_runs(id) ON DELETE CASCADE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    server_id BIGINT NOT NULL,
    encrypted_log TEXT,
    run_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    FOREIGN KEY(server_id, org_id) REFERENCES infrastructure_servers(id, org_id) ON DELETE CASCADE,
    UNIQUE(run_id, server_id)
);
CREATE INDEX ON shell_script_run_servers(server_id);

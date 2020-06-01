CREATE TABLE request_to_shell_run_link (
    request_id BIGINT NOT NULL REFERENCES generic_requests(id) ON DELETE CASCADE,
    run_id BIGINT NOT NULL REFERENCES shell_script_runs(id) ON DELETE CASCADE,
    UNIQUE(request_id, run_id)
);

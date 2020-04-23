-- Need this table too because we can't just rely on runs since we need the link to client_scripts even
-- when no runs have been made yet.
CREATE TABLE scheduled_task_script_links (
    event_id BIGINT NOT NULL UNIQUE REFERENCES scheduled_tasks(id) ON DELETE CASCADE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    script_id BIGINT NOT NULL,
    FOREIGN KEY(script_id, org_id) REFERENCES client_scripts(id, org_id) ON DELETE CASCADE
);
CREATE INDEX ON scheduled_task_script_links(script_id);
CREATE INDEX ON scheduled_task_script_links(org_id);

CREATE TABLE scheduled_task_script_run_links (
    event_id BIGINT NOT NULL REFERENCES scheduled_tasks(id) ON DELETE CASCADE,
    run_id BIGINT NOT NULL UNIQUE REFERENCES script_runs(id) ON DELETE CASCADE
);
CREATE INDEX ON scheduled_task_script_run_links(event_id);

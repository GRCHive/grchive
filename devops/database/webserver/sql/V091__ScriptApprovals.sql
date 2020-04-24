CREATE TABLE generic_requests(
    id BIGSERIAL PRIMARY KEY,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    upload_time TIMESTAMPTZ NOT NULL,
    upload_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    name TEXT NOT NULL,
    assignee BIGINT REFERENCES users(id) ON DELETE NO ACTION,
    due_date TIMESTAMPTZ,
    description TEXT
);
CREATE INDEX ON generic_requests(org_id);

CREATE TABLE generic_approval (
    id BIGINT PRIMARY KEY,
    request_id BIGINT NOT NULL REFERENCES generic_requests(id) ON DELETE CASCADE,
    response_time TIMESTAMPTZ NOT NULL,
    responder_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    response BOOLEAN NOT NULL,
    reason TEXT
);

CREATE TABLE request_to_script_run_link (
    request_id BIGINT NOT NULL REFERENCES generic_requests(id) ON DELETE CASCADE,
    run_id BIGINT NOT NULL REFERENCES script_runs(id) ON DELETE CASCADE,
    UNIQUE(request_id, run_id)
);

CREATE TABLE request_to_scheduled_task_link (
    request_id BIGINT NOT NULL REFERENCES generic_requests(id) ON DELETE CASCADE,
    task_id BIGINT NOT NULL REFERENCES scheduled_tasks(id) ON DELETE CASCADE,
    UNIQUE(request_id, task_id)
);

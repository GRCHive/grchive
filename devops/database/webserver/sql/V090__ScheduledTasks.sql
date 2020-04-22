CREATE TABLE scheduled_tasks (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    description TEXT DEFAULT '' NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE one_time_tasks (
    event_id BIGINT REFERENCES scheduled_tasks(id) ON DELETE CASCADE,
    event_date_time TIMESTAMPTZ NOT NULL
);
CREATE INDEX ON one_time_tasks(event_id);

CREATE TABLE recurring_tasks (
    event_id BIGINT REFERENCES scheduled_tasks(id) ON DELETE CASCADE,
    start_date_time TIMESTAMPTZ NOT NULL,
    rrule TEXT NOT NULL
);

CREATE INDEX ON recurring_tasks(event_id);

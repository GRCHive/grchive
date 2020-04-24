DROP TABLE scheduled_task_script_links;
CREATE TABLE scheduled_task_script_links (
    event_id BIGINT NOT NULL UNIQUE REFERENCES scheduled_tasks(id) ON DELETE CASCADE,
    link_id BIGINT NOT NULL REFERENCES code_to_client_scripts_link(id) ON DELETE CASCADE
);
CREATE INDEX ON scheduled_task_script_links(link_id);

CREATE TABLE request_to_shell_script_link (
    id BIGSERIAL PRIMARY KEY,
    request_id BIGINT NOT NULL REFERENCES generic_requests(id) ON DELETE CASCADE UNIQUE,
    shell_version_id BIGINT NOT NULL REFERENCES shell_script_versions(id) ON DELETE CASCADE
);
CREATE INDEX ON request_to_shell_script_link(shell_version_id);

CREATE TABLE request_shell_script_to_servers_link (
    link_id BIGINT NOT NULL REFERENCES request_to_shell_script_link(id) ON DELETE CASCADE,
    server_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY(server_id, org_id) REFERENCES infrastructure_servers(id, org_id),
    UNIQUE(link_id, server_id)
);

CREATE INDEX ON request_shell_script_to_servers_link(server_id);
CREATE INDEX ON request_shell_script_to_servers_link(link_id);

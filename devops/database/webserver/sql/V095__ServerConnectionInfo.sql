CREATE TABLE server_username_password_connection (
    id BIGSERIAL PRIMARY KEY,
    server_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    FOREIGN KEY(server_id, org_id) REFERENCES infrastructure_servers(id, org_id) ON DELETE CASCADE,
    UNIQUE(server_id)
);

CREATE TABLE server_ssh_key_connection (
    id BIGSERIAL PRIMARY KEY,
    server_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    username TEXT NOT NULL,
    private_key TEXT NOT NULL,
    FOREIGN KEY(server_id, org_id) REFERENCES infrastructure_servers(id, org_id) ON DELETE CASCADE,
    UNIQUE(server_id)
);

-- Script is linked to commits, commits are linked to data and parameters.
CREATE TABLE client_script_code_data_sources (
    link_id BIGINT NOT NULL REFERENCES code_to_client_scripts_link(id) ON DELETE CASCADE,
    data_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY(data_id, org_id) REFERENCES client_data(id, org_id) ON DELETE CASCADE,
    UNIQUE(link_id, data_id, org_id)
);
CREATE INDEX ON client_script_code_data_sources(link_id);

CREATE TABLE supported_code_parameter_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE,
    golang_type TEXT NOT NULL,
    kotlin_type TEXT NOT NULL
);

CREATE TABLE client_script_code_parameters (
    link_id BIGINT NOT NULL REFERENCES code_to_client_scripts_link(id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL,
    param_type INTEGER NOT NULL REFERENCES supported_code_parameter_types(id) ON DELETE RESTRICT,
    UNIQUE(link_id, name)
);
CREATE INDEX ON client_script_code_parameters(link_id);

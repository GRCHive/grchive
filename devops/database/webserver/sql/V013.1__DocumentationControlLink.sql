CREATE TABLE _base_link_controls_control_documentation (
    category_id BIGINT NOT NULL,
    org_id INT NOT NULL,
    control_id BIGINT NOT NULL
);

CREATE TABLE controls_input_documentation (
    CONSTRAINT cat_org_fkey
        FOREIGN KEY(category_id, org_id)
        REFERENCES process_flow_control_documentation_categories(id, org_id)
        ON DELETE CASCADE,
    CONSTRAINT control_org_fkey
        FOREIGN KEY(control_id, org_id)
        REFERENCES process_flow_controls(id, org_id)
        ON DELETE CASCADE,
    UNIQUE(category_id, org_id, control_id)
) INHERITS (_base_link_controls_control_documentation);
CREATE TABLE controls_output_documentation (
    CONSTRAINT cat_org_fkey
        FOREIGN KEY(category_id, org_id)
        REFERENCES process_flow_control_documentation_categories(id, org_id)
        ON DELETE CASCADE,
    CONSTRAINT control_org_fkey
        FOREIGN KEY(control_id, org_id)
        REFERENCES process_flow_controls(id, org_id)
        ON DELETE CASCADE,
    UNIQUE(category_id, org_id, control_id)
) INHERITS (_base_link_controls_control_documentation);

INSERT INTO process_flow_node_types (name, description)
VALUES
    ('Activity (Manual)', ''),
    ('Activity (Automated)', ''),
    ('Decision', ''),
    ('Start', ''),
    ('General Ledger Entry', ''),
    ('System', '');

INSERT INTO process_flow_input_output_type(name)
VALUES
    ('Execution'),
    ('Data');
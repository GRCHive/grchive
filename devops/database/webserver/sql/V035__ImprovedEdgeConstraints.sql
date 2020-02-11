DROP TRIGGER IF EXISTS ensure_edges_valid_on_input_io_update ON process_flow_node_inputs;
DROP TRIGGER IF EXISTS ensure_edges_valid_on_output_io_update ON process_flow_node_outputs;
DROP FUNCTION IF EXISTS ensure_edges_valid_on_input_update;
DROP FUNCTION IF EXISTS ensure_edges_valid_on_output_update;
DROP FUNCTION IF EXISTS check_edge_valid;

DROP VIEW IF EXISTS full_process_flow_edges_view;

CREATE VIEW full_process_flow_edges_view AS
    SELECT 
        input_node.id AS input_node_id,
        input.io_type_id AS input_io_type_id,
        input.id AS input_io_id,
        output_node.id AS output_node_id,
        output.io_type_id AS output_io_type_id,
        output.id AS output_io_id
    FROM process_flow_edges AS edge
    INNER JOIN process_flow_node_inputs AS input
        ON edge.input_id = input.id
    INNER JOIN process_flow_node_outputs AS output
        ON edge.output_id = output.id
    INNER JOIN process_flow_nodes AS input_node
        ON input.parent_node_id = input_node.id
    INNER JOIN process_flow_nodes AS output_node
        ON output.parent_node_id = output_node.id;

CREATE OR REPLACE FUNCTION check_edge_valid(edge full_process_flow_edges_view)
    RETURNS void AS
$$
    BEGIN
        IF edge.input_node_id = edge.output_node_id THEN
            RAISE EXCEPTION 'Edge connects a node output to an input on the same node.';
        END IF;

        IF edge.input_io_type_id != edge.output_io_type_id THEN
            RAISE EXCEPTION 'Connecting two IOs of different types.';
        END IF;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION ensure_edges_valid_on_input_update()
    RETURNS trigger AS
$ensure_edges_valid_on_input_update$
    DECLARE
        tmp RECORD;
    BEGIN
        SELECT check_edge_valid(v.*) INTO tmp
        FROM full_process_flow_edges_view AS v
        WHERE v.input_io_id = OLD.id;

        RETURN NEW;
    END;
$ensure_edges_valid_on_input_update$ LANGUAGE plpgsql;

CREATE TRIGGER ensure_edges_valid_on_input_io_update
    AFTER UPDATE OF io_type_id ON process_flow_node_inputs 
    FOR EACH ROW
    EXECUTE FUNCTION ensure_edges_valid_on_input_update();

CREATE OR REPLACE FUNCTION ensure_edges_valid_on_output_update()
    RETURNS trigger AS
$ensure_edges_valid_on_output_update$
    DECLARE
        tmp RECORD;
    BEGIN
        SELECT check_edge_valid(v.*) INTO tmp
        FROM full_process_flow_edges_view AS v
        WHERE v.output_io_id = OLD.id;

        RETURN NEW;
    END;
$ensure_edges_valid_on_output_update$ LANGUAGE plpgsql;

CREATE TRIGGER ensure_edges_valid_on_output_io_update
    AFTER UPDATE OF io_type_id ON process_flow_node_outputs
    FOR EACH ROW
    EXECUTE FUNCTION ensure_edges_valid_on_output_update();

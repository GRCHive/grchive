ALTER TABLE process_flow_node_inputs
ADD COLUMN io_order INTEGER;

ALTER TABLE process_flow_node_outputs
ADD COLUMN io_order INTEGER;

CREATE TEMPORARY SEQUENCE seq_node_inputs;

CREATE OR REPLACE FUNCTION fill_in_node_input_output_order(node_id BIGINT, type_id INTEGER)
    RETURNS void AS
$$
    BEGIN
        PERFORM setval('seq_node_inputs', 1, false);
        UPDATE process_flow_node_inputs
            SET io_order = nextval('seq_node_inputs')
            WHERE parent_node_id=node_id
                AND io_type_id=type_id;

        PERFORM setval('seq_node_inputs', 1, false);
        UPDATE process_flow_node_outputs
            SET io_order = nextval('seq_node_inputs')
            WHERE parent_node_id=node_id
                AND io_type_id=type_id;
    END;
$$ LANGUAGE plpgsql;

DO $$ BEGIN
    PERFORM fill_in_node_input_output_order(node.id, typ.id)
    FROM process_flow_nodes AS node
    CROSS JOIN process_flow_input_output_type AS typ;
END $$ LANGUAGE plpgsql;

ALTER TABLE process_flow_node_inputs
ALTER COLUMN io_order SET NOT NULL;

ALTER TABLE process_flow_node_outputs
ALTER COLUMN io_order SET NOT NULL;

DROP FUNCTION IF EXISTS fill_in_node_input_output_order;

package database

import (
	"errors"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func DeleteEdgeFromId(edgeId int64) error {
	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM process_flow_edges
		WHERE id = $1
	`, edgeId)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func FindAllEdgesForProcessFlow(flowId int64) ([]*core.ProcessFlowEdge, error) {
	edges := []*core.ProcessFlowEdge{}
	err := dbConn.Select(&edges, `
		SELECT DISTINCT edge.*
		FROM process_flow_edges as edge
		INNER JOIN process_flow_node_inputs AS input
			ON edge.input_id = input.id
		INNER JOIN process_flow_node_outputs AS output
			ON edge.output_id = output.id
		INNER JOIN process_flow_nodes AS node
			ON input.parent_node_id = node.id
				OR output.parent_node_id = node.id
		WHERE node.process_flow_id = $1
	`, flowId)
	if err != nil {
		return nil, err
	}
	return edges, nil
}

func CreateNewProcessFlowEdge(edge *core.ProcessFlowEdge) (*core.ProcessFlowEdge, error) {
	if edge.InputIoId == edge.OutputIoId {
		return nil, errors.New("Can not create an edge from a node to itself.")
	}

	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		WITH edge AS (
			INSERT INTO process_flow_edges (input_id, output_id)
			VALUES (:input_id, :output_id)
			RETURNING *
		)
		SELECT 
			edge.id AS "edge.id",
			edge.input_id AS "edge.input_id",
			edge.output_id AS "edge.output_id",
			input.parent_node_id AS "input.parent_node_id",
			input.io_type_id AS "input.io_type_id",
			output.parent_node_id AS "output.parent_node_id",
			output.io_type_id AS "output.io_type_id"
		FROM edge 
			INNER JOIN process_flow_node_inputs AS input
				ON edge.input_id = input.id
			INNER JOIN process_flow_node_outputs AS output
				ON edge.output_id = output.id
	`, edge)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	type QueryResult struct {
		Edge   core.ProcessFlowEdge        `db:"edge"`
		Input  core.ProcessFlowInputOutput `db:"input"`
		Output core.ProcessFlowInputOutput `db:"output"`
	}

	result := QueryResult{}
	rows.Next()
	err = rows.StructScan(&result)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Sanity checks to make sure that the edge is valid.
	// 1) Two types must be the same.
	// 2) Must be an input to an output (or vice versa) [this probably
	// 	  gets resolved by the SQL query?]
	// 3) The edge must not be from the node to the same node.
	if result.Input.TypeId != result.Output.TypeId {
		tx.Rollback()
		return nil, errors.New("Edge Type ID mismatch.")
	}

	if result.Input.ParentNodeId == result.Output.ParentNodeId {
		tx.Rollback()
		return nil, errors.New("Node loop edge.")
	}

	rows.Close()
	err = tx.Commit()
	return &result.Edge, err
}

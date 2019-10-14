package database

import (
	"errors"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"time"
)

func FindProcessFlowWithId(id int64) (*core.ProcessFlow, error) {
	rows, err := dbConn.Queryx(`
		SELECT
			pf.id,
			pf.name,
			org.id AS "org.id",
			org.org_group_id AS "org.org_group_id",
			org.org_group_name AS "org.org_group_name",
			org.org_name AS "org.org_name",
			pf.description,
			pf.created_time,
			pf.last_updated_time
		FROM process_flows AS pf
		INNER JOIN organizations AS org
			ON pf.org_id = org.id
		WHERE pf.id = $1
	`, id)
	if err != nil {
		return nil, err
	}

	var flow *core.ProcessFlow = new(core.ProcessFlow)
	flow.Org = new(core.Organization)

	rows.Next()
	err = rows.StructScan(flow)
	if err != nil {
		return nil, err
	}

	rows.Close()
	return flow, nil
}

func UpdateProcessFlow(flow *core.ProcessFlow) error {
	flow.LastUpdatedTime = time.Now().UTC()
	tx := dbConn.MustBegin()
	_, err := tx.NamedExec(`
		UPDATE process_flows
		SET name = :name, description = :description
		WHERE id = :id
	`, flow)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func InsertNewProcessFlow(flow *core.ProcessFlow) error {
	var err error

	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		INSERT INTO process_flows (name, org_id, description, created_time, last_updated_time)
		VALUES (:name, :org.id, :description, :created_time, :last_updated_time)
		RETURNING id
	`, flow)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&flow.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()

	err = tx.Commit()
	return err
}

func FindOrganizationProcessFlows(org *core.Organization) ([]*core.ProcessFlow, error) {
	result := []*core.ProcessFlow{}

	err := dbConn.Select(&result, `
		SELECT id, name, description, created_time, last_updated_time FROM process_flows WHERE org_id = $1 ORDER BY name ASC
	`, org.Id)

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(result); i++ {
		result[i].Org = org
	}
	return result, nil
}

// Finds all process flows and finds the result index of the specified process flow
func FindOrganizationProcessFlowsWithIndex(org *core.Organization, processFlowId int64) ([]*core.ProcessFlow, int, error) {
	result, err := FindOrganizationProcessFlows(org)
	if err != nil {
		return nil, 0, err
	}

	// TODO: speed this up somehow if/when necessary? Hopefully # of flows stays reasonable.
	// Can't assume the result is stored in a way where we can binary search.
	var resultIndex int = 0
	var found bool = false
	for i := 0; i < len(result); i++ {
		if result[i].Id == processFlowId {
			resultIndex = i
			found = true
			break
		}
	}

	if !found {
		return nil, 0, errors.New("Failed to find requested process flow id in results.")
	}

	return result, resultIndex, nil
}

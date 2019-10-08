package database

import (
	"errors"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

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

	if !rows.Next() {
		return rows.Err()
	}
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
		SELECT id, name, description, created_time, last_updated_time FROM process_flows WHERE org_id = $1
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
func FindOrganizationProcessFlowsWithIndex(org *core.Organization, processFlowId uint32) ([]*core.ProcessFlow, uint32, error) {
	result, err := FindOrganizationProcessFlows(org)
	if err != nil {
		return nil, 0, err
	}

	// TODO: speed this up somehow if/when necessary? Hopefully # of flows stays reasonable.
	// Can't assume the result is stored in a way where we can binary search.
	var resultIndex uint32 = 0
	var found bool = false
	for i := 0; i < len(result); i++ {
		if result[i].Id == processFlowId {
			resultIndex = uint32(i)
			found = true
			break
		}
	}

	if !found {
		return nil, 0, errors.New("Failed to find requested process flow id in results.")
	}

	return result, resultIndex, nil
}

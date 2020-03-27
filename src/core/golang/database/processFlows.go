package database

import (
	"errors"
	"gitlab.com/grchive/grchive/core"
	"strconv"
	"time"
)

func FindProcessFlowWithId(id int64, role *core.Role) (*core.ProcessFlow, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

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
	defer rows.Close()

	var flow *core.ProcessFlow = new(core.ProcessFlow)
	flow.Org = new(core.Organization)

	rows.Next()
	err = rows.StructScan(flow)
	if err != nil {
		return nil, err
	}

	return flow, LogAuditSelect(flow.Org.Id, core.ResourceIdProcessFlow, strconv.FormatInt(flow.Id, 10), role)
}

func UpdateProcessFlow(flow *core.ProcessFlow, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	flow.LastUpdatedTime = time.Now().UTC()
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
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

func InsertNewProcessFlow(flow *core.ProcessFlow, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	var err error

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

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
		rows.Close()
		tx.Rollback()
		return err
	}
	rows.Close()

	err = tx.Commit()
	return err
}

func FindOrganizationProcessFlows(org *core.Organization, role *core.Role) ([]*core.ProcessFlow, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

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

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, p := range result {
		err = LogAuditSelectWithTx(org.Id, core.ResourceIdProcessFlow, strconv.FormatInt(p.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return result, tx.Commit()
}

// Finds all process flows and finds the result index of the specified process flow
func FindOrganizationProcessFlowsWithIndex(org *core.Organization, processFlowId int64, role *core.Role) ([]*core.ProcessFlow, int, error) {
	result, err := FindOrganizationProcessFlows(org, role)
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

func DeleteProcessFlow(flowId int64, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM process_flows
		WHERE id = $1
	`, flowId)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

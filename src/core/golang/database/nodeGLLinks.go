package database

import (
	"gitlab.com/grchive/grchive/core"
)

func NewNodeGLLink(nodeId int64, accountId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		INSERT INTO node_gl_link (node_id, gl_account_id, org_id)
		VALUES ($1, $2, $3)
	`, nodeId, accountId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func DeleteNodeGLLink(nodeId int64, accountId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM node_gl_link
		WHERE node_id = $1
			AND gl_account_id = $2
			AND org_id = $3
	`, nodeId, accountId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func AllGLLinkedToNode(nodeId int64, orgId int32, role *core.Role) ([]*core.GeneralLedgerAccount, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	accounts := make([]*core.GeneralLedgerAccount, 0)
	err := dbConn.Select(&accounts, `
		SELECT acc.*
		FROM general_ledger_accounts AS acc
		INNER JOIN node_gl_link AS link
			ON link.gl_account_id = acc.id
		WHERE link.node_id = $1 AND link.org_id = $2
	`, nodeId, orgId)
	return accounts, err
}

func AllFlowsRelatedToGL(accountId int64, orgId int32, role *core.Role) ([]*core.ProcessFlow, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	flows := make([]*core.ProcessFlow, 0)
	err := dbConn.Select(&flows, `
		SELECT DISTINCT
			flow.id,
			flow.name,
			org.id AS "org.id",
			org.org_group_id AS "org.org_group_id",
			org.org_group_name AS "org.org_group_name",
			org.org_name AS "org.org_name",
			flow.description,
			flow.created_time,
			flow.last_updated_time
		FROM process_flows AS flow
		INNER JOIN process_flow_nodes AS node
			ON node.process_flow_id = flow.id
		INNER JOIN node_gl_link AS link
			ON link.node_id = node.id
		INNER JOIN organizations AS org
			ON flow.org_id = org.id
		WHERE link.gl_account_id = $1 AND org.id = $2
	`, accountId, orgId)
	return flows, err
}

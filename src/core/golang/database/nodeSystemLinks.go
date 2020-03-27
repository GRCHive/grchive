package database

import (
	"gitlab.com/grchive/grchive/core"
	"strconv"
)

func NewNodeSystemLink(nodeId int64, systemId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceSystems, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		INSERT INTO node_system_link (node_id, system_id, org_id)
		VALUES ($1, $2, $3)
	`, nodeId, systemId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func DeleteNodeSystemLink(nodeId int64, systemId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceSystems, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM node_system_link
		WHERE node_id = $1
			AND system_id = $2
			AND org_id = $3
	`, nodeId, systemId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func AllSystemsLinkedToNode(nodeId int64, orgId int32, role *core.Role) ([]*core.System, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceSystems, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	systems := make([]*core.System, 0)
	err := dbConn.Select(&systems, `
		SELECT sys.*
		FROM systems AS sys
		INNER JOIN node_system_link AS link
			ON link.system_id = sys.id
		WHERE link.node_id = $1 AND link.org_id = $2
	`, nodeId, orgId)

	if err != nil {
		return nil, err
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, s := range systems {
		err = LogAuditSelectWithTx(orgId, core.ResourceIdSystem, strconv.FormatInt(s.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return systems, tx.Commit()
}

func AllFlowsRelatedToSystem(systemId int64, orgId int32, role *core.Role) ([]*core.ProcessFlow, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceSystems, core.AccessView) {
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
		INNER JOIN node_system_link AS link
			ON link.node_id = node.id
		INNER JOIN organizations AS org
			ON flow.org_id = org.id
		WHERE link.system_id = $1 AND org.id = $2
	`, systemId, orgId)

	if err != nil {
		return nil, err
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, f := range flows {
		err = LogAuditSelectWithTx(orgId, core.ResourceIdProcessFlow, strconv.FormatInt(f.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return flows, tx.Commit()
}

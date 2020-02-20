package database

import (
	"gitlab.com/grchive/grchive/core"
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
	return systems, err
}

package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func CreateNewSystem(system *core.System, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceSystems, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		INSERT INTO systems(org_id, name, purpose)
		VALUES (:org_id, :name, :purpose)
		RETURNING id
	`, system)

	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&system.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()

	return tx.Commit()
}

func GetAllSystemsForOrg(orgId int32, role *core.Role) ([]*core.System, error) {
	if !role.Permissions.HasAccess(core.ResourceSystems, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	systems := make([]*core.System, 0)
	err := dbConn.Select(&systems, `
		SELECT *
		FROM systems
		WHERE org_id = $1
	`, orgId)
	return systems, err
}

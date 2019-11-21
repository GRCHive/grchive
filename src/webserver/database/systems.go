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
		INSERT INTO systems(org_id, name, purpose, description)
		VALUES (:org_id, :name, :purpose, :description)
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

func GetSystem(sysId int64, orgId int32, role *core.Role) (*core.System, error) {
	if !role.Permissions.HasAccess(core.ResourceSystems, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	sys := core.System{}
	err := dbConn.Get(&sys, `
		SELECT *
		FROM systems
		WHERE id = $1
			AND org_id = $2
	`, sysId, orgId)

	return &sys, err
}

func EditSystem(sys *core.System, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceSystems, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	tx := dbConn.MustBegin()
	_, err := tx.NamedExec(`
		UPDATE systems
		SET name = :name,
			purpose = :purpose,
			description = :description
		WHERE id = :id
			AND org_id = :org_id
	`, sys)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func DeleteSystem(sysId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceSystems, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM systems
		WHERE id = $1
			AND org_id = $2
	`, sysId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

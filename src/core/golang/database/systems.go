package database

import (
	"gitlab.com/grchive/grchive/core"
)

func CreateNewSystem(system *core.System, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceSystems, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

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
		rows.Close()
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

	if err != nil {
		return nil, err
	}

	return systems, nil
}

func GetAllSystemsForOrgWithDeployment(orgId int32, deploymentType int32, role *core.Role) ([]*core.System, error) {
	if !role.Permissions.HasAccess(core.ResourceSystems, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	systems := make([]*core.System, 0)
	err := dbConn.Select(&systems, `
		SELECT sys.*
		FROM systems AS sys
		INNER JOIN deployment_system_link AS link
			ON link.system_id = sys.id
		INNER JOIN deployments AS dp
			ON dp.id = link.deployment_id
		WHERE sys.org_id = $1
			AND dp.deployment_type = $2
	`, orgId, deploymentType)

	if err != nil {
		return nil, err
	}

	return systems, nil
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

	if err != nil {
		return nil, err
	}

	return &sys, nil
}

func EditSystem(sys *core.System, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceSystems, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}
	_, err = tx.NamedExec(`
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

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
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

func FindSystemIdsForDatabase(dbId int64, orgId int32, role *core.Role) ([]int64, error) {
	if !role.Permissions.HasAccess(core.ResourceSystems, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	ids := make([]int64, 0)
	err := dbConn.Select(&ids, `
		SELECT system_id
		FROM database_system_link
		WHERE db_id = $1
			AND org_id = $2
	`, dbId, orgId)
	return ids, err
}

func LinkDatabasesToSystem(sysId int64, orgId int32, dbIds []int64, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceSystems, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceDatabases, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	for _, dbId := range dbIds {
		_, err := tx.Exec(`
			INSERT INTO database_system_link (db_id, org_id, system_id)
			VALUES ($1, $2, $3)
		`, dbId, orgId, sysId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func DeleteDatabaseSystemLink(sysId int64, dbId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceSystems, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceDatabases, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM database_system_link
		WHERE system_id = $1
			AND db_id = $2
			AND org_id = $3
	`, sysId, dbId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func GetAllSupportedDatabaseTypes(role *core.Role) ([]*core.DatabaseType, error) {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	types := make([]*core.DatabaseType, 0)
	err := dbConn.Select(&types, `
		SELECT *
		FROM supported_databases
	`)
	return types, err
}

func InsertNewDatabase(db *core.Database, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		INSERT INTO database_resources(name, org_id, type_id, other_type, version)
		VALUES (:name, :org_id, :type_id, :other_type, :version)
		RETURNING id
	`, db)

	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&db.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()

	return tx.Commit()
}

func GetAllDatabasesForOrg(orgId int32, role *core.Role) ([]*core.Database, error) {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	dbs := make([]*core.Database, 0)
	err := dbConn.Select(&dbs, `
		SELECT *
		FROM database_resources
		WHERE org_id = $1
	`, orgId)

	return dbs, err
}

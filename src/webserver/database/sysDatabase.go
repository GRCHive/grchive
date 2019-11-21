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

func GetDb(dbId int64, orgId int32, role *core.Role) (*core.Database, error) {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	db := core.Database{}
	err := dbConn.Get(&db, `
		SELECT *
		FROM database_resources
		WHERE id = $1
			AND org_id = $2
	`, dbId, orgId)

	return &db, err
}

func EditDb(db *core.Database, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	tx := dbConn.MustBegin()
	_, err := tx.NamedExec(`
		UPDATE database_resources
		SET name = :name,
			type_id = :type_id,
			other_type = :other_type,
			version = :version
		WHERE id = :id
			AND org_id = :org_id
	`, db)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func DeleteDb(dbId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessManage) {
		return core.ErrorUnauthorized
	}
	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM database_resources
		WHERE id = $1
			AND org_id = $2
	`, dbId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()

}

package database

import (
	"gitlab.com/grchive/grchive/core"
)

func GetAllTablesForSchema(schemaId int64, orgId int32, role *core.Role) ([]*core.DbTable, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	data := make([]*core.DbTable, 0)
	err := dbConn.Select(&data, `
		SELECT *
		FROM database_tables
		WHERE schema_id = $1 AND org_id = $2
		ORDER BY table_name ASC
	`, schemaId, orgId)
	return data, err
}

func GetAllColumnsForTable(tableId int64, orgId int32, role *core.Role) ([]*core.DbColumn, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	data := make([]*core.DbColumn, 0)
	err := dbConn.Select(&data, `
		SELECT *
		FROM database_columns
		WHERE table_id = $1 AND org_id = $2
		ORDER BY column_name ASC
	`, tableId, orgId)
	return data, err
}

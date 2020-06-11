package database

import (
	"fmt"
	"github.com/jmoiron/sqlx/types"
	"gitlab.com/grchive/grchive/core"
	"strings"
)

func GetAllTablesForSchema(
	schemaId int64,
	orgId int32,
	start int64,
	limit int64,
	filter string,
	role *core.Role,
) ([]*core.DbTable, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	filterQuery := ""
	if filter != "" {
		filterQuery = fmt.Sprintf(" AND table_name LIKE '%%%s%%'", filter)
	}

	query := strings.Builder{}
	query.WriteString(fmt.Sprintf(`
		SELECT id, org_id, schema_id, table_name, columns
		FROM database_tables
		WHERE schema_id = $1 AND org_id = $2 %s
		ORDER BY table_name ASC
	`, filterQuery))

	if limit >= 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", limit))
	}

	if start >= 0 {
		query.WriteString(fmt.Sprintf(" OFFSET %d", start))
	}

	rows, err := dbConn.Queryx(query.String(), schemaId, orgId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data := make([]*core.DbTable, 0)
	for rows.Next() {
		tbl := core.DbTable{}
		jsonData := types.JSONText{}

		err = rows.Scan(
			&tbl.Id,
			&tbl.OrgId,
			&tbl.SchemaId,
			&tbl.TableName,
			&jsonData,
		)
		if err != nil {
			return nil, err
		}

		err = jsonData.Unmarshal(&tbl.Columns)
		if err != nil {
			return nil, err
		}

		data = append(data, &tbl)
	}

	return data, nil
}

func GetAllFunctionsForSchema(schemaId int64, orgId int32, role *core.Role) ([]*core.DbFunction, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSql, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	data := make([]*core.DbFunction, 0)
	err := dbConn.Select(&data, `
		SELECT *
		FROM database_functions
		WHERE schema_id = $1 AND org_id = $2
		ORDER BY name ASC
	`, schemaId, orgId)
	return data, err
}

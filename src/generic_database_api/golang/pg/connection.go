package pg_api

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"strings"
)

type PgRole struct {
	RolSuper      bool `db:"rolsuper"`
	RolCreateRole bool `db:"rolcreaterole"`
	RolCreateDb   bool `db:"rolcreatedb"`
	RolBypassRLS  bool `db:"rolbypassrls"`
}

type PgGrant struct {
	Select     bool
	Insert     bool
	Update     bool
	Delete     bool
	Truncate   bool
	References bool
	Trigger    bool
}

type PgTableGrants struct {
	tableGrant   PgGrant
	columnGrants map[string]PgGrant
}

// Table Name is Key
type PgSchemaGrants map[string]PgTableGrants

func (g PgGrant) IsReadOnly() bool {
	return g.Select &&
		!g.Insert &&
		!g.Update &&
		!g.Delete &&
		!g.Truncate &&
		!g.References &&
		!g.Trigger
}

func CreateGrantFromRawPrivileges(priv []string) PgGrant {
	g := PgGrant{}
	for _, p := range priv {
		switch p {
		case "SELECT":
			g.Select = true
		case "INSERT":
			g.Insert = true
		case "UPDATE":
			g.Update = true
		case "DELETE":
			g.Delete = true
		case "TRUNCATE":
			g.Truncate = true
		case "REFERENCES":
			g.References = true
		case "TRIGGER":
			g.Trigger = true
		default:
			core.Warning("Unknown privilege: " + p)
		}
	}
	return g
}

func (g PgTableGrants) IsReadOnly() bool {
	if !g.tableGrant.IsReadOnly() {
		return false
	}

	for _, c := range g.columnGrants {
		if !c.IsReadOnly() {
			return false
		}
	}

	return true
}

func (g PgSchemaGrants) IsReadOnly() bool {
	for _, t := range g {
		if !t.IsReadOnly() {
			return false
		}
	}
	return true
}

func retrieveRole(conn *sqlx.DB, username string) (*PgRole, error) {
	rows, err := conn.Queryx(`
		SELECT
			rolsuper,
			rolcreaterole,
			rolcreatedb,
			rolbypassrls
		FROM pg_roles
		WHERE rolname = $1
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rows.Next()
	result := PgRole{}
	err = rows.StructScan(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func retrieveColumnGrants(conn *sqlx.DB, username string, dbName string, schemaName string, tableName string) (map[string]PgGrant, error) {
	grants := map[string]PgGrant{}

	rows, err := conn.Queryx(`
		SELECT
			column_name,
			array_to_string(array_agg(privilege_type), ',')
		FROM information_schema.role_column_grants
		WHERE grantee = $1
			AND table_catalog = $2
			AND table_schema = $3
			AND table_name = $4
		GROUP BY column_name
	`, username, dbName, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type RawSchemaGrant struct {
		ColumnName string `db:"column_name"`
		Privileges string `db:"array_to_string"`
	}

	for rows.Next() {
		rawGrant := RawSchemaGrant{}
		err = rows.StructScan(&rawGrant)
		if err != nil {
			return nil, err
		}

		grants[rawGrant.ColumnName] = CreateGrantFromRawPrivileges(strings.Split(rawGrant.Privileges, ","))
	}

	return grants, nil
}

func retrieveSchemaGrants(conn *sqlx.DB, username string, dbName string) (*PgSchemaGrants, error) {
	grants := PgSchemaGrants{}

	rows, err := conn.Queryx(`
		SELECT
			table_schema,
			table_name,
			array_to_string(array_agg(privilege_type), ',')
		FROM information_schema.role_table_grants
		WHERE grantee = $1
			AND table_catalog = $2
		GROUP BY table_schema, table_name
	`, username, dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type RawSchemaGrant struct {
		TableSchema string `db:"table_schema"`
		TableName   string `db:"table_name"`
		Privileges  string `db:"array_to_string"`
	}

	for rows.Next() {
		rawGrant := RawSchemaGrant{}
		err = rows.StructScan(&rawGrant)
		if err != nil {
			return nil, err
		}

		tb := PgTableGrants{
			tableGrant: CreateGrantFromRawPrivileges(strings.Split(rawGrant.Privileges, ",")),
		}

		tb.columnGrants, err = retrieveColumnGrants(conn, username, dbName, rawGrant.TableSchema, rawGrant.TableName)
		if err != nil {
			return nil, err
		}

		grants[fmt.Sprintf("%s.%s", rawGrant.TableSchema, rawGrant.TableName)] = tb
	}

	return &grants, nil
}

func (pg *PgDriver) ConnectionReadOnly() bool {
	if pg.currentRole.RolSuper ||
		pg.currentRole.RolCreateRole ||
		pg.currentRole.RolCreateDb ||
		pg.currentRole.RolBypassRLS {
		return false
	}

	return pg.grants.IsReadOnly()
}

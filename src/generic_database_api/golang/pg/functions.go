package pg_api

import (
	"gitlab.com/grchive/grchive/core"
)

func (pg *PgDriver) GetFunctions(schema *core.DbSchema) ([]*core.DbFunction, error) {
	rows, err := pg.connection.Queryx(`
		SELECT
			proc.proname AS "name",
			pg_get_functiondef(proc.oid) AS "src",
			pg_get_function_result(proc.oid) AS "ret_type"
		FROM pg_proc AS proc
		INNER JOIN pg_namespace AS nm
			ON proc.pronamespace = nm.oid
		WHERE nm.nspname = $1
	`, schema.SchemaName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fns := make([]*core.DbFunction, 0)
	for rows.Next() {
		f := core.DbFunction{}
		err = rows.StructScan(&f)
		if err != nil {
			return nil, err
		}
		fns = append(fns, &f)
	}
	return fns, nil
}

package pg_api

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"strings"
	"time"
)

func (pg *PgDriver) GetState() *core.FullDbState {
	return pg.state
}

func (pg *PgDriver) ConstructState() error {
	if pg.state != nil {
		return nil
	}
	start := time.Now()
	pg.state = &core.FullDbState{}

	core.Debug("\tQuery Schemas")
	// Schemas
	{
		rows, err := pg.connection.QueryContext(context.Background(), `
			SELECT schema_name
			FROM information_schema.schemata
			WHERE catalog_name = $1
				AND schema_name != 'information_schema'
				AND schema_name != 'pg_catalog'
				AND schema_name NOT LIKE 'pg_toast%'
				AND schema_name NOT LIKE 'pg_temp%'
		`, pg.connInfo.DbName)
		if err != nil {
			return err
		}
		defer rows.Close()

		core.Debug("\t\tParse Schemas")
		for rows.Next() {
			s := core.DbSchema{}
			err = rows.Scan(&s.SchemaName)
			if err != nil {
				return err
			}

			// Ignore system schemas.
			if s.SchemaName == "SYS" ||
				s.SchemaName == "SYSTEM" ||
				strings.HasPrefix(s.SchemaName, "_SYS") {
				continue
			}

			pg.state.AddSchema(&s)
		}
	}

	core.Debug("\tQuery Functions")
	// Functions
	{
		rows, err := pg.connection.QueryContext(context.Background(), `
			SELECT
				proc.proname AS "name",
				pg_get_functiondef(proc.oid) AS "src",
				pg_get_function_result(proc.oid) AS "ret_type",
				nm.nspname
			FROM pg_proc AS proc
			INNER JOIN pg_namespace AS nm
				ON proc.pronamespace = nm.oid
			WHERE nm.nspname != 'information_schema'
				AND nm.nspname != 'pg_catalog'
				AND nm.nspname NOT LIKE 'pg_toast%'
				AND nm.nspname NOT LIKE 'pg_temp%'
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		schemaName := ""

		core.Debug("\t\tParse Functions")
		for rows.Next() {
			s := core.DbFunction{}
			err = rows.Scan(&s.Name, &s.Src, &s.RetType, &schemaName)
			if err != nil {
				return err
			}

			sch := pg.state.GetSchema(schemaName)
			if sch != nil {
				pg.state.AddFunction(&sch.DbSchema, &s)
			}
		}
	}

	core.Debug("\tQuery Tables")
	// Tables
	{
		rows, err := pg.connection.QueryContext(context.Background(), `
			SELECT
				table_name,
				table_schema,
				jsonb_agg(
					jsonb_build_object(
						'Name',
						column_name,
						'Type',
						data_type
					)
				)
			FROM information_schema.columns
			WHERE table_catalog = $1
				AND table_schema != 'information_schema'
				AND table_schema != 'information_schema'
				AND table_schema != 'pg_catalog'
				AND table_schema NOT LIKE 'pg_toast%'
				AND table_schema NOT LIKE 'pg_temp%'
			GROUP BY table_name, table_schema
		`, pg.connInfo.DbName)
		if err != nil {
			return err
		}
		defer rows.Close()

		columns := ""
		schemaName := ""

		core.Debug("\t\tParse Tables")
		for rows.Next() {
			s := core.DbTable{}
			err = rows.Scan(&s.TableName, &schemaName, &columns)
			if err != nil {
				return err
			}

			s.Columns = make([]*core.RawDbColumn, 0)
			err = json.Unmarshal([]byte(columns), &s.Columns)
			if err != nil {
				return err
			}

			sch := pg.state.GetSchema(schemaName)
			if sch != nil {
				sch.AddTable(&s)
			}
		}
	}

	elapsed := time.Since(start)
	core.Debug(fmt.Sprintf(
		"\tFinish Constructing State: %f seconds",
		float64(elapsed.Milliseconds())/1000.0,
	))
	return nil
}

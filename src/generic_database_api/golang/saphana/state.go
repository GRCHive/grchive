package saphana_api

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"strings"
	"time"
)

func (hdb *SapHanaDriver) GetState() *core.FullDbState {
	return hdb.state
}

func (hdb *SapHanaDriver) ConstructState() error {
	if hdb.state != nil {
		return nil
	}

	start := time.Now()
	hdb.state = &core.FullDbState{}

	core.Debug("\tQuery Schemas")
	// Schemas
	{
		rows, err := hdb.connection.QueryContext(context.Background(), `
			SELECT SCHEMA_NAME
			FROM PUBLIC.SCHEMAS
		`)
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

			hdb.state.AddSchema(&s)
		}
	}

	core.Debug("\tQuery Functions")
	// Functions
	{
		rows, err := hdb.connection.QueryContext(context.Background(), `
			SELECT
				PROCEDURE_NAME,
				TO_NVARCHAR(DEFINITION),
				SCHEMA_NAME
			FROM PUBLIC.PROCEDURES
			UNION
			SELECT
				FUNCTION_NAME,
				TO_NVARCHAR(DEFINITION),
				SCHEMA_NAME
			FROM PUBLIC.FUNCTIONS
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		schemaName := ""

		core.Debug("\t\tParse Functions")
		for rows.Next() {
			s := core.DbFunction{}
			err = rows.Scan(&s.Name, &s.Src, &schemaName)
			if err != nil {
				return err
			}

			sch := hdb.state.GetSchema(schemaName)
			if sch != nil {
				hdb.state.AddFunction(&sch.DbSchema, &s)
			}
		}
	}

	core.Debug("\tQuery Tables")
	// Tables
	{
		rows, err := hdb.connection.QueryContext(context.Background(), `
			SELECT 
				COL.TABLE_NAME,
				COL.SCHEMA_NAME,
				'[' || STRING_AGG('{"Name": "' || COL.COLUMN_NAME || '", "Type": "' || COL.DATA_TYPE_NAME || '"}', ',') || ']'
			FROM PUBLIC.TABLE_COLUMNS AS "COL"
			INNER JOIN PUBLIC.TABLES AS "TBL"
				ON COL.TABLE_NAME = TBL.TABLE_NAME
			WHERE TBL.IS_SYSTEM_TABLE = 'FALSE'
				AND COL.SCHEMA_NAME != 'SYS'
				AND COL.SCHEMA_NAME != 'SYSTEM'
				AND COL.SCHEMA_NAME NOT LIKE '_SYS%'
			GROUP BY COL.TABLE_NAME, COL.SCHEMA_NAME
		`)
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

			sch := hdb.state.GetSchema(schemaName)
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

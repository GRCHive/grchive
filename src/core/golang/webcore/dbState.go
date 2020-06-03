package webcore

import (
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
)

func GetDbStateFromRefresh(refreshId int64, orgId int32) (*core.FullDbState, error) {
	state := core.FullDbState{}

	schemas, err := database.GetAllDatabaseSchemaForRefresh(refreshId, orgId, core.ServerRole)
	if err != nil {
		return nil, err
	}

	for _, sch := range schemas {
		state.AddSchema(sch)

		fns, err := database.GetAllFunctionsForSchema(sch.Id, sch.OrgId, core.ServerRole)
		if err != nil {
			return nil, err
		}

		for _, fn := range fns {
			state.AddFunction(sch, fn)
		}

		tables, err := database.GetAllTablesForSchema(sch.Id, sch.OrgId, core.ServerRole)
		if err != nil {
			return nil, err
		}

		for _, tbl := range tables {
			state.AddTable(sch, tbl)

			columns, err := database.GetAllColumnsForTable(tbl.Id, tbl.OrgId, core.ServerRole)
			if err != nil {
				return nil, err
			}

			for _, col := range columns {
				state.AddColumn(sch, tbl, col)
			}
		}
	}

	return &state, nil
}

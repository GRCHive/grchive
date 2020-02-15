package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type AllDatabaseSchemasInputs struct {
	RefreshId int64 `webcore:"refreshId"`
	OrgId     int32 `webcore:"orgId"`
}

func allDatabaseSchemas(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllDatabaseSchemasInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, err := database.GetAllDatabaseSchemaForRefresh(inputs.RefreshId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get schemas: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(data)
}

type GetDatabaseSchemasInput struct {
	SchemaId int64 `webcore:"schemaId"`
	OrgId    int32 `webcore:"orgId"`
	FnMode   bool  `webcore:"fnMode"`
}

func getDatabaseSchema(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetDatabaseSchemasInput{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type ColumnMap map[int64][]*core.DbColumn
	type SchemaOutput struct {
		Tables  []*core.DbTable
		Columns ColumnMap
	}

	retStruct := struct {
		Schema    *SchemaOutput
		Functions *[]*core.DbFunction
	}{}

	if inputs.FnMode {
		functions, err := database.GetAllFunctionsForSchema(inputs.SchemaId, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get functions: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		retStruct.Functions = &functions
	} else {
		tables, err := database.GetAllTablesForSchema(inputs.SchemaId, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get tables: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		allColumns := ColumnMap{}

		for _, tbl := range tables {
			allColumns[tbl.Id], err = database.GetAllColumnsForTable(tbl.Id, inputs.OrgId, role)
			if err != nil {
				core.Warning("Failed to get columns: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		retStruct.Schema = &SchemaOutput{
			Tables:  tables,
			Columns: allColumns,
		}
	}

	jsonWriter.Encode(retStruct)
}

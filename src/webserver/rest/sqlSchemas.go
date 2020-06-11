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
	SchemaId int64           `webcore:"schemaId"`
	OrgId    int32           `webcore:"orgId"`
	FnMode   bool            `webcore:"fnMode"`
	Start    int64           `webcore:"start"`
	Limit    int64           `webcore:"limit"`
	Filter   core.NullString `webcore:"filter,optional"`
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

	type SchemaOutput struct {
		Tables []*core.DbTable
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
		tables, err := database.GetAllTablesForSchema(
			inputs.SchemaId,
			inputs.OrgId,
			inputs.Start,
			inputs.Limit,
			inputs.Filter.NullString.String,
			role)
		if err != nil {
			core.Warning("Failed to get tables: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		retStruct.Schema = &SchemaOutput{
			Tables: tables,
		}
	}

	jsonWriter.Encode(retStruct)
}

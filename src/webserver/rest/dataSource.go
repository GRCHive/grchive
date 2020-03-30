package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"math"
	"net/http"
)

func allDataSourceOptions(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	// The only thing we can do for metadata requests is to make sure
	// a valid API key is found.
	_, err := webcore.FindApiKeyInContext(r.Context())
	if err != nil {
		core.Warning("Can't find API key: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	options, err := database.AllDataSourceOptions()
	if err != nil {
		core.Warning("Failed to get all data source options: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(options)
}

type GetDataSourceInput struct {
	Source core.DataSourceLink `webcore:"source"`
}

func getDataSource(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetDataSourceInput{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = webcore.GetCurrentRequestRole(r, inputs.Source.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch inputs.Source.SourceId {
	case core.SourceGrchive:
		jsonWriter.Encode(core.ResourceHandle{
			DisplayText: "GRCHive",
		})
	case core.SourceDbPostgres:
		dbId := int64(math.Round(inputs.Source.SourceTarget["id"].(float64)))
		hd, err := webcore.GetResourceHandle(core.ResourceIdDatabase, dbId, inputs.Source.OrgId)
		if err != nil {
			core.Warning("Failed to get database handle: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonWriter.Encode(hd)
	default:
		core.Warning("Unrecognized source Id.")
		w.WriteHeader(http.StatusBadRequest)
	}
}

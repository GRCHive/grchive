package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"net/http"
)

func getAllProcessFlowNodeTypes(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	types, err := database.GetAllProcessFlowNodeTypes()
	if err != nil {
		core.Warning("Can't get types: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jsonWriter.Encode(struct{}{})
		return
	}
	jsonWriter.Encode(types)
}

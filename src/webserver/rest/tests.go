package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type GetCodeRunTestInput struct {
	OrgId   int32 `webcore:"orgId"`
	RunId   int64 `webcore:"runId"`
	Summary bool  `webcore:"summary"`
}

func getCodeRunTest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetCodeRunTestInput{}
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

	if inputs.Summary {
		summary, err := database.GetCodeRunTestSummary(inputs.RunId, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get test summary: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonWriter.Encode(summary)
	} else {
		core.Warning("Non summary not supported yet.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

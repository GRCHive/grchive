package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

type NewDeploymentInputs struct {
	OrgId    int32          `json:"orgId"`
	SystemId core.NullInt64 `json:"systemId"`
	DbId     core.NullInt64 `json:"dbId"`
}

func newDeployment(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewDeploymentInputs{}
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

	tx := database.CreateTx()
	deployment, err := database.CreateNewEmptyDeploymentWithTx(inputs.OrgId, role, tx)
	if err != nil {
		tx.Rollback()
		core.Warning("Failed to create new deployment: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if inputs.SystemId.NullInt64.Valid {
		err = database.LinkDeploymentWithSystemWithTx(deployment, inputs.SystemId.NullInt64.Int64, role, tx)
	} else if inputs.DbId.NullInt64.Valid {
		err = database.LinkDeploymentWithDatabaseWithTx(deployment, inputs.DbId.NullInt64.Int64, role, tx)
	}

	if err != nil {
		tx.Rollback()
		core.Warning("Failed to link deployment: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if tx.Commit() != nil {
		tx.Rollback()
		core.Warning("Failed to commit new deployment: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(deployment)
}

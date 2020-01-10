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

	deployment, err = database.GetDeploymentFromId(deployment.Id, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get full deployment: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(deployment)
}

type UpdateDeploymentInputs struct {
	Deployment *core.StrippedFullDeployment `json:"deployment"`
}

func updateDeployment(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := UpdateDeploymentInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.Deployment.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = database.UpdateDeployment(inputs.Deployment, role)
	if err != nil {
		core.Warning("Failed to update deployment: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fullDeployment, err := database.GetDeploymentFromId(inputs.Deployment.Id, inputs.Deployment.OrgId, role)
	if err != nil {
		core.Warning("Failed to retrieve full deployment: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(fullDeployment)
}

type NewDeploymentServerLinkInputs struct {
	SystemId []int64 `json:"systemId"`
	DbId     []int64 `json:"dbId"`
	ServerId int64   `json:"serverId"`
	OrgId    int32   `json:"orgId"`
}

func newDeploymentServerLink(w http.ResponseWriter, r *http.Request) {
	inputs := NewDeploymentServerLinkInputs{}
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

	deploymentIds := make([]int64, 0)

	for _, sysId := range inputs.SystemId {
		id, err := database.GetSystemDeploymentId(sysId, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get system deployment id: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		deploymentIds = append(deploymentIds, id)
	}

	for _, dbId := range inputs.DbId {
		id, err := database.GetDatabaseDeploymentId(dbId, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get database deployment id: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		deploymentIds = append(deploymentIds, id)
	}

	err = database.LinkDeploymentsWithServer(deploymentIds, inputs.ServerId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Link deployment with server failed: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type DeleteDeploymentServerLinkInputs struct {
	SystemId core.NullInt64 `json:"systemId"`
	DbId     core.NullInt64 `json:"dbId"`
	ServerId int64          `json:"serverId"`
	OrgId    int32          `json:"orgId"`
}

func deleteDeploymentServerLink(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteDeploymentServerLinkInputs{}
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

	if inputs.SystemId.NullInt64.Valid {
		deploymentId, err := database.GetSystemDeploymentId(inputs.SystemId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get deployment id: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = database.DeleteDeploymentServerLink(deploymentId, inputs.ServerId, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to delete server link: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if inputs.DbId.NullInt64.Valid {
		deploymentId, err := database.GetDatabaseDeploymentId(inputs.DbId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get deployment id: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = database.DeleteDeploymentServerLink(deploymentId, inputs.ServerId, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to delete server link: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type GetLogsInput struct {
	OrgId      int32           `webcore:"orgId"`
	CommitHash core.NullString `webcore:"commitHash,optional"`
	RunId      core.NullInt64  `webcore:"runId,optional"`
	RunLog     core.NullBool   `webcore:"runLog,optional"`
}

func getLogs(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetLogsInput{}
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

	var encryptedLogs string
	if inputs.CommitHash.NullString.Valid {
		encryptedLogs, err = database.GetBuildLogs(inputs.CommitHash.NullString.String, inputs.OrgId, role)
	} else if inputs.RunId.NullInt64.Valid {
		runId := inputs.RunId.NullInt64.Int64
		if inputs.RunLog.NullBool.Valid && inputs.RunLog.NullBool.Bool {
			encryptedLogs, err = database.GetRunLogsForRun(runId, role)
		} else {
			encryptedLogs, err = database.GetBuildLogsForRun(runId, role)
		}
	} else {
		core.Warning("Invalid inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		core.Warning("Failed to get logs: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if encryptedLogs != "" {
		decrypt, err := vault.TransitDecrypt(core.EnvConfig.LogEncryptionPath, []byte(encryptedLogs))
		if err != nil {
			core.Warning("Failed to decrypt result: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonWriter.Encode(string(decrypt))
	} else {
		jsonWriter.Encode("Logs not found.")
	}
}

package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type AllFileVersionsInputs struct {
	FileId int64 `webcore:"fileId"`
	OrgId  int32 `webcore:"orgId"`
}

func allFileVersions(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllFileVersionsInputs{}
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

	versions, err := database.AllFileVersions(inputs.FileId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get file versions: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jsonWriter.Encode(versions)
}

type GetFileVersionInputs struct {
	FileId  int64 `webcore:"fileId"`
	OrgId   int32 `webcore:"orgId"`
	Version int32 `webcore:"version"`
}

func getFileVersion(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetFileVersionInputs{}
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

	data, err := database.GetFileVersionStorageData(inputs.FileId, inputs.OrgId, inputs.Version, role)
	if err != nil {
		core.Warning("Failed to get file version storage data: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	previewStorage, err := database.GetPreviewFileVersionStorageData(inputs.FileId, inputs.OrgId, inputs.Version, role)
	if err != nil {
		core.Warning("Failed to get preview file version storage data: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jsonWriter.Encode(struct {
		Storage    *core.FileStorageData
		HasPreview bool
	}{
		Storage:    data,
		HasPreview: (previewStorage != nil),
	})
}

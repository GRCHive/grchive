package rest

import (
	"encoding/json"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gcloud_api"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
)

type AllShellsInput struct {
	ShellType int32 `webcore:"shellType"`
}

func allShells(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllShellsInput{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Can't find org: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shells, err := database.GetShellScriptsOfTypeForOrganization(inputs.ShellType, org.Id)
	if err != nil {
		core.Warning("Failed to get shell scripts: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(shells)
}

type AllShellVersionsInput struct {
	ShellId int64 `webcore:"shellId"`
}

func allShellVersions(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllShellVersionsInput{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Can't find org: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	versions, err := database.AllShellScriptVersions(inputs.ShellId, org.Id)
	if err != nil {
		core.Warning("Failed to get shell scripts: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(versions)
}

type NewShellInput struct {
	ShellType   int32  `json:"shellType"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Script      string `json:"script"`
}

func newShell(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewShellInput{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Can't find org: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.FindRoleInContext(r.Context())
	if err != nil {
		core.Warning("Can't find role: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	script := core.ShellScript{
		OrgId:       org.Id,
		TypeId:      inputs.ShellType,
		Name:        inputs.Name,
		Description: inputs.Description,
	}

	var gcsGeneration int64

	tx := database.CreateTx()
	database.WrapTx(tx, func() error {
		err = database.NewShellScriptWithTx(tx, &script)
		if err != nil {
			core.Warning("Failed to create new shell script: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		return nil
	}, func() error {
		script.BucketId = core.EnvConfig.Gcloud.ShellBucket
		script.StorageId = fmt.Sprintf("shellscript-%d", script.Id)

		encryptedScript, err := vault.TransitEncrypt("shell", []byte(inputs.Script))
		if err != nil {
			core.Warning("Failed to encrypt script: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		storage := gcloud.DefaultGCloudApi.GetStorageApi()
		gcsGeneration, err = storage.Upload(script.BucketId, script.StorageId, encryptedScript, core.EnvConfig.HmacKey)
		if err != nil {
			core.Warning("Failed to upload script: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		return nil
	}, func() error {
		err := database.UpdateShellScriptGCSStorageWithTx(tx, script.Id, script.BucketId, script.StorageId)
		if err != nil {
			core.Warning("Failed to update shell script GCS: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		return err
	}, func() error {
		err := database.CreateShellScriptVersionWithTx(tx, script.Id, script.OrgId, role.UserId, gcsGeneration)
		if err != nil {
			core.Warning("Failed to create shell script version: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		return err
	})

	jsonWriter.Encode(script)
}

func deleteShellScript(w http.ResponseWriter, r *http.Request) {
	script, err := webcore.FindShellScriptInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get shell script in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	versions, err := database.AllShellScriptVersions(script.Id, script.OrgId)
	if err != nil {
		core.Warning("Failed to get shell script versions: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		err := database.DeleteShellScriptFromIdWithTx(tx, script.Id)
		if err != nil {
			return err
		}
		return nil
	}, func() error {
		storage := gcloud.DefaultGCloudApi.GetStorageApi()
		for _, v := range versions {
			err = storage.DeleteVersioned(script.BucketId, script.StorageId, v.GcsGeneration)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		core.Warning("Failed to delete shell script: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

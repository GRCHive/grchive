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
	"time"
)

const ShellEncryptionPath = "shell"

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

func allShellVersions(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	script, err := webcore.FindShellScriptInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get shell script in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Can't find org: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	versions, err := database.AllShellScriptVersions(script.Id, org.Id)
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

		encryptedScript, err := vault.TransitEncrypt(ShellEncryptionPath, []byte(inputs.Script))
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

func getShellScript(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	script, err := webcore.FindShellScriptInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get shell script in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonWriter.Encode(script)
}

type EditShellInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func editShellScript(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	script, err := webcore.FindShellScriptInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get shell script in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := EditShellInput{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	script.Name = inputs.Name
	script.Description = inputs.Description

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.EditShellScriptWithTx(tx, script)
	})

	if err != nil {
		core.Warning("Failed to edit shell script: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(script)
}

func getShellVersion(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	script, err := webcore.FindShellScriptInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get shell script in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	version, err := webcore.FindShellScriptVersionInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get shell script version in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storage := gcloud.DefaultGCloudApi.GetStorageApi()
	rawData, err := storage.DownloadVersioned(
		script.BucketId,
		script.StorageId,
		version.GcsGeneration,
		core.EnvConfig.HmacKey,
	)

	if err != nil {
		core.Warning("Failed to download version from GCS: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	decryptedScript, err := vault.TransitDecrypt(ShellEncryptionPath, rawData)
	if err != nil {
		core.Warning("Failed to encrypt script: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(string(decryptedScript))
}

type NewShellVersionInputs struct {
	Script string `json:"script"`
}

func newShellVersion(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	script, err := webcore.FindShellScriptInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get shell script in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.FindRoleInContext(r.Context())
	if err != nil {
		core.Warning("Can't find role: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputs := NewShellVersionInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storage := gcloud.DefaultGCloudApi.GetStorageApi()

	encryptedScript, err := vault.TransitEncrypt(ShellEncryptionPath, []byte(inputs.Script))
	if err != nil {
		core.Warning("Failed to encrypt script: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	gcsGeneration, err := storage.Upload(script.BucketId, script.StorageId, encryptedScript, core.EnvConfig.HmacKey)
	if err != nil {
		core.Warning("Failed to upload script: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newVersion := core.ShellScriptVersion{
		ShellId:       script.Id,
		OrgId:         script.OrgId,
		UploadUserId:  role.UserId,
		UploadTime:    time.Now().UTC(),
		GcsGeneration: gcsGeneration,
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.NewShellScriptVersionWithTx(tx, &newVersion)
	})

	if err != nil {
		core.Warning("Failed to create script version: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(newVersion)
}

type RunShellVersionInputs struct {
	Servers []int64 `json:"servers"`
}

func runShellVersion(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := RunShellVersionInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	script, err := webcore.FindShellScriptInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get shell script in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	version, err := webcore.FindShellScriptVersionInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get shell script version in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.FindRoleInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get shell script version in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := webcore.FindOrganizationInContext(r.Context())
	if err != nil {
		core.Warning("Failed to get shell script version in context: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// First create the shell run to store the information
	// about on what servers we want to run this shell script.
	// Next, create a corresponding request for approval. This shell script
	// will only be run upon approval. We do this all in a single transaction
	// to make sure we don't have a dangling shell script run that can't be
	// approved and vice versa.
	run := core.ShellScriptRun{
		ScriptVersionId: version.Id,
		RunUserId:       role.UserId,
		CreateTime:      time.Now().UTC(),
	}

	request := core.GenericRequest{
		OrgId:        org.Id,
		UploadTime:   time.Now().UTC(),
		UploadUserId: role.UserId,
		Name:         fmt.Sprintf("Shell Script Run Request: %s", script.Name),
	}

	result := struct {
		RunId     int64
		RequestId core.NullInt64
	}{}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		err := database.NewShellRunWithTx(tx, &run)
		if err != nil {
			return err
		}
		result.RunId = run.Id
		return nil
	}, func() error {
		for _, sid := range inputs.Servers {
			serverRun := core.ShellScriptRunPerServer{
				RunId:    run.Id,
				OrgId:    org.Id,
				ServerId: sid,
			}

			err := database.NewShellServerRunWithTx(tx, &serverRun)
			if err != nil {
				return err
			}
		}
		return nil
	}, func() error {
		err := database.CreateGenericRequestWithTx(tx, &request)
		if err != nil {
			return err
		}
		result.RequestId = core.CreateNullInt64(request.Id)
		return nil
	}, func() error {
		return database.LinkShellRunToRequestWithTx(tx, result.RunId, result.RequestId.NullInt64.Int64)
	})

	if err != nil {
		core.Warning("Failed to create new shell run and approval: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(result)
}

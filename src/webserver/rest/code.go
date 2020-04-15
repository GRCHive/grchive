package rest

import (
	"encoding/json"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"time"
)

type SaveCodeInput struct {
	OrgId      int32          `json:"orgId"`
	Code       string         `json:"code"`
	DataId     core.NullInt64 `json:"dataId"`
	ScriptId   core.NullInt64 `json:"scriptId"`
	ScriptData *struct {
		Params       []*core.CodeParameter `json:"params"`
		ClientDataId []int64               `json:"clientDataId"`
	} `json:"scriptData"`
}

func saveCode(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := SaveCodeInput{}
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

	if !inputs.DataId.NullInt64.Valid && !inputs.ScriptId.NullInt64.Valid {
		core.Warning("Invalid combination of inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	managedCode := core.ManagedCode{
		OrgId:      inputs.OrgId,
		ActionTime: time.Now().UTC(),
	}

	// Determine GitPath from whether this is a data or script.
	if inputs.DataId.NullInt64.Valid {
		clientData, err := database.GetClientDataFromId(inputs.DataId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get client data: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// For now, assume Kotlin always. If we want to support more in the future we'll have to
		// somehow get this information from the user or something.
		managedCode.GitPath = fmt.Sprintf("src/main/kotlin/data/%s", clientData.Data.Filename("kt"))
	} else if inputs.ScriptId.NullInt64.Valid {
		script, err := database.GetClientScriptFromId(inputs.ScriptId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get client script: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		managedCode.GitPath = fmt.Sprintf("src/main/kotlin/scripts/%s", script.Filename("kt"))
		metadataGitPath := fmt.Sprintf("src/main/resources/scripts/%s", script.MetadataFilename())

		// Hack the StoreManagedCodeToGitea to store a file in an easy way but not keep track of it in the DB.
		tmpCode := core.ManagedCode{
			OrgId:   managedCode.OrgId,
			GitPath: metadataGitPath,
		}

		metadata, err := webcore.GenerateScriptMetadataYaml(inputs.ScriptData.Params, inputs.ScriptData.ClientDataId)
		if err != nil {
			core.Warning("Failed to generate metadata code: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = webcore.StoreManagedCodeToGitea(
			&tmpCode,
			metadata,
			nil,
			"[CI SKIP] Update Metadata: "+metadataGitPath,
		)
		if err != nil {
			core.Warning("Failed to store metadata : " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// There is a possibility here that the link will fail after the storing to Gitea succeeds.
	// Do we care in that case? We can probably survive just losing a link since storing an
	// extra commit in Gitea won't hurt us.
	err = webcore.StoreManagedCodeToGitea(&managedCode, inputs.Code, role, "Update: "+managedCode.GitPath)
	if err != nil {
		core.Warning("Failed to store managed code: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if inputs.DataId.NullInt64.Valid {
		err = database.LinkCodeToData(managedCode.Id, inputs.DataId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to link code to data: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if inputs.ScriptId.NullInt64.Valid {
		tx, err := database.CreateAuditTrailTx(role)
		if err != nil {
			core.Warning("Failed to create tx: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = database.WrapTx(tx, func() error {
			err := database.LinkCodeToScriptWithTx(managedCode.Id, inputs.ScriptId.NullInt64.Int64, inputs.OrgId, role, tx)
			if err != nil {
				return err
			}

			for _, p := range inputs.ScriptData.Params {
				err = database.LinkScriptToParameterWithTx(
					inputs.ScriptId.NullInt64.Int64,
					managedCode.Id,
					inputs.OrgId,
					p.Name,
					p.ParamId,
					role,
					tx,
				)

				if err != nil {
					return err
				}
			}

			for _, d := range inputs.ScriptData.ClientDataId {
				err = database.LinkScriptToDataSourceWithTx(
					inputs.ScriptId.NullInt64.Int64,
					managedCode.Id,
					inputs.OrgId,
					d,
					role,
					tx,
				)

				if err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			core.Warning("Failed to link script: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	jsonWriter.Encode(managedCode)
}

type GetCodeInput struct {
	OrgId    int32          `webcore:"orgId"`
	CodeId   int64          `webcore:"codeId"`
	DataId   core.NullInt64 `webcore:"dataId,optional"`
	ScriptId core.NullInt64 `webcore:"scriptId,optional"`
}

func getCode(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetCodeInput{}
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

	// Need to do a check that the user actually has access to the resource
	// that wraps the code.
	if inputs.DataId.NullInt64.Valid {
		ok, err := database.CheckValidCodeDataLink(inputs.CodeId, inputs.DataId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil || !ok {
			core.Warning("Invalid code data link: " + core.ErrorString(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else if inputs.ScriptId.NullInt64.Valid {
		ok, err := database.CheckValidCodeScriptLink(inputs.CodeId, inputs.ScriptId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil || !ok {
			core.Warning("Invalid code script link: " + core.ErrorString(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		core.Warning("Invalid combination of inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code, err := webcore.GetManagedCodeFromGitea(inputs.CodeId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get code: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type RetScriptData struct {
		Params     []*core.CodeParameter
		ClientData []*core.FullClientDataWithLink
	}

	ret := struct {
		Code       string
		ScriptData *RetScriptData
	}{
		Code:       code,
		ScriptData: nil,
	}

	if inputs.ScriptId.NullInt64.Valid {
		ret.ScriptData = &RetScriptData{}

		ret.ScriptData.Params, err = database.GetLinkedParametersToScriptCode(
			inputs.ScriptId.NullInt64.Int64,
			inputs.CodeId,
			inputs.OrgId,
			role,
		)

		if err != nil {
			core.Warning("Failed to get linked parameters: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ret.ScriptData.ClientData, err = database.GetLinkedDataSourceToScriptCode(
			inputs.ScriptId.NullInt64.Int64,
			inputs.CodeId,
			inputs.OrgId,
			role,
		)

		if err != nil {
			core.Warning("Failed to get linked client data: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	jsonWriter.Encode(ret)
}

type AllCodeInput struct {
	OrgId    int32          `webcore:"orgId"`
	DataId   core.NullInt64 `webcore:"dataId,optional"`
	ScriptId core.NullInt64 `webcore:"scriptId,optional"`
}

func allCode(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllCodeInput{}
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

	var code []*core.ManagedCode

	if inputs.DataId.NullInt64.Valid {
		code, err = database.AllManagedCodeForDataId(inputs.DataId.NullInt64.Int64, inputs.OrgId, role)
	} else if inputs.ScriptId.NullInt64.Valid {
		code, err = database.AllManagedCodeForScriptId(inputs.ScriptId.NullInt64.Int64, inputs.OrgId, role)
	} else {
		core.Warning("Invalid combination of inputs.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		core.Warning("Failed to get managed code: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(code)
}

type GetCodeBuildStatusInput struct {
	OrgId      int32  `webcore:"orgId"`
	CommitHash string `webcore:"commitHash"`
}

func getCodeBuildStatus(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetCodeBuildStatusInput{}
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

	status, err := database.GetCodeBuildStatus(inputs.CommitHash, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get build status: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(status)
}

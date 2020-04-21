package rest

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	drone "gitlab.com/grchive/grchive/drone_api"
	gitea "gitlab.com/grchive/grchive/gitea_api"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"strconv"
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
		managedCode.GitPath = clientData.Data.Filename("kt")
	} else if inputs.ScriptId.NullInt64.Valid {
		script, err := database.GetClientScriptFromId(inputs.ScriptId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get client script: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		managedCode.GitPath = script.Filename("kt")
		metadataGitPath := script.MetadataFilename()

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

		// This needs to happen before the code so that running the code at a certain revision will also pick up
		// the metadata changes.
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
	OrgId      int32           `webcore:"orgId"`
	CodeId     core.NullInt64  `webcore:"codeId,optional"`
	CodeCommit core.NullString `webcore:"codeCommit,optional"`
	DataId     core.NullInt64  `webcore:"dataId,optional"`
	ScriptId   core.NullInt64  `webcore:"scriptId,optional"`
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

	var fullCode *core.ManagedCode
	if inputs.CodeId.NullInt64.Valid {
		fullCode, err = database.GetCode(inputs.CodeId.NullInt64.Int64, inputs.OrgId, role)
	} else if inputs.CodeCommit.NullString.Valid {
		fullCode, err = database.GetCodeFromCommit(inputs.CodeCommit.NullString.String, inputs.OrgId, role)
	} else {
		core.Warning("Invalid combination of inputs to pull code.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		core.Warning("Failed to pull code: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Need to do a check that the user actually has access to the resource
	// that wraps the code. If user doesn't supply data id or script id assume
	// that they don't actually want the code.
	var sendCode bool = false
	if inputs.DataId.NullInt64.Valid {
		sendCode, err = database.CheckValidCodeDataLink(fullCode.Id, inputs.DataId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil || !sendCode {
			core.Warning("Invalid code data link: " + core.ErrorString(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else if inputs.ScriptId.NullInt64.Valid {
		sendCode, err = database.CheckValidCodeScriptLink(fullCode.Id, inputs.ScriptId.NullInt64.Int64, inputs.OrgId, role)
		if err != nil || !sendCode {
			core.Warning("Invalid code script link: " + core.ErrorString(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var code string = ""
	if sendCode {
		code, err = webcore.GetManagedCodeFromGitea(fullCode.Id, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get code: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	type RetScriptData struct {
		Params     []*core.CodeParameter
		ClientData []*core.FullClientDataWithLink
	}

	ret := struct {
		Code       string
		ScriptData *RetScriptData
		Full       *core.ManagedCode
	}{
		Code:       code,
		ScriptData: nil,
		Full:       fullCode,
	}

	if inputs.ScriptId.NullInt64.Valid {
		ret.ScriptData = &RetScriptData{}

		ret.ScriptData.Params, err = database.GetLinkedParametersToScriptCode(
			inputs.ScriptId.NullInt64.Int64,
			fullCode.Id,
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
			fullCode.Id,
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
		code, err = database.AllManagedCodeForOrgId(inputs.OrgId, role)
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

type RunCodeInput struct {
	OrgId  int32                  `json:"orgId"`
	CodeId int64                  `json:"codeId"`
	Latest bool                   `json:"latest"`
	Params map[string]interface{} `json:"params"`
}

func runCode(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := RunCodeInput{}
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

	code, err := database.GetCode(inputs.CodeId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to find code: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Need to make sure this code is linked to a script - otherwise
	// "running" it makes no sense.
	script, err := database.GetScriptForCode(inputs.CodeId, inputs.OrgId, role)
	if err != nil || script == nil {
		core.Warning("Failed to find script: " + core.ErrorString(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !inputs.Latest {
		// In the case where we're not trying to run the latest code, we need to
		// make sure that the version that the client requested to run has actually
		// compiled successfully.
		status, err := database.GetCodeBuildStatus(code.GitHash, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get status: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !status.Success {
			core.Warning("Failed to run a script that failed to compile.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Create a DB entry to track the run. Return this to the user.
	// We don't need to roll this back in case of an error later on as
	// ideally any later stages will log those changes and just let the user
	// know there in the logs stored in the DB.
	run, err := database.CreateScriptRun(code.Id, inputs.OrgId, script.Id, inputs.Latest, inputs.Params, role)
	if err != nil {
		core.Warning("Failed to create script run: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if inputs.Latest {
		repo, err := database.GetLinkedGiteaRepository(inputs.OrgId)
		if err != nil {
			core.Warning("Failed to get linked Gitea repository: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// In this case, we need to fire off a Drone CI job to compile the latest code + the current script revision.
		// We must specify branch/commit here due to a Gitea issue with the /repos/{owner}/{repo}/commits/{ref} API endpoint
		// that Drone uses to find the latest commit of the branch. This endpoint doesn't work in Gitea so we need to specify the
		// commit directly.
		commitSha, err := gitea.GlobalGiteaApi.RepositoryGitGetRefSha(
			gitea.GiteaRepository{
				Owner: repo.GiteaOrg,
				Name:  repo.GiteaRepo,
			},
			"refs/heads/master",
		)

		if err != nil {
			core.Warning("Failed to get latest Git commit: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = drone.GlobalDroneApi.BuildCreate(repo.GiteaOrg, repo.GiteaRepo, map[string]string{
			"branch":     "master",
			"commit":     commitSha,
			"SCRIPT_RUN": strconv.FormatInt(run.Id, 10),
		})

		if err != nil {
			core.Warning("Failed to create Drone build: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		// Grab JAR path from drone CI since that's the only place we store it.
		// This is hitting the same DB table as the GetCodeBuildStatus call earlier, can we merge it somehow?
		jar, err := database.GetCodeJar(code.Id, code.OrgId, role)
		if err != nil {
			core.Warning("Failed to get JAR for code: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// In this case, we can directly send off a request to make the script runner run this script.
		webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
			Exchange: webcore.DEFAULT_EXCHANGE,
			Queue:    webcore.SCRIPT_RUNNER_QUEUE,
			Body: webcore.ScriptRunnerMessage{
				RunId: run.Id,
				Jar:   jar,
			},
		})
	}

	jsonWriter.Encode(run.Id)
}

type AllCodeRunsInput struct {
	OrgId    int32          `webcore:"orgId"`
	ScriptId core.NullInt64 `webcore:"scriptId,optional"`
}

func allCodeRuns(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllCodeRunsInput{}
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

	var runs []*core.ScriptRun
	if inputs.ScriptId.NullInt64.Valid {
		runs, err = database.GetAllScriptRunsForScriptId(inputs.ScriptId.NullInt64.Int64, inputs.OrgId, role)
	} else {
		runs, err = database.GetAllScriptRunsForOrgId(inputs.OrgId, role)
	}

	if err != nil {
		core.Warning("Failed to get script runs: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(runs)
}

type GetCodeRunInput struct {
	OrgId int32 `webcore:"orgId"`
	RunId int64 `webcore:"runId"`
}

func getCodeRun(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetCodeRunInput{}
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

	run, err := database.GetScriptRun(inputs.RunId, role)
	if err != nil {
		core.Warning("Failed to get script run: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(run)
}

type GetScriptCodeLinkInput struct {
	OrgId  int32 `webcore:"orgId"`
	LinkId int64 `webcore:"linkId"`
}

func getScriptCodeLink(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetScriptCodeLinkInput{}
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

	script, err := database.GetScriptFromScriptCodeLink(inputs.LinkId, role)
	if err != nil {
		core.Warning("Failed to get linked script: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	code, err := database.GetCodeFromScriptCodeLink(inputs.LinkId, role)
	if err != nil {
		core.Warning("Failed to get linked code: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Script *core.ClientScript
		Code   *core.ManagedCode
	}{
		Script: script,
		Code:   code,
	})
}

type GetCodeLinkInput struct {
	OrgId  int32 `webcore:"orgId"`
	CodeId int64 `webcore:"codeId"`
}

func getCodeLink(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetCodeLinkInput{}
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

	script, err := database.GetScriptForCode(inputs.CodeId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get script for code: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := database.GetClientDataForCode(inputs.CodeId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get data for code: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Data   *core.ClientData
		Script *core.ClientScript
	}{
		Data:   data,
		Script: script,
	})
}

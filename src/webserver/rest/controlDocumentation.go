package rest

import (
	"bytes"
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gcloud_api"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
	"io"
	"net/http"
	"time"
)

type NewControlDocCatInputs struct {
	Name        string `webcore:"name"`
	Description string `webcore:"description"`
	OrgId       int32  `webcore:"orgId"`
}

type EditControlDocCatInputs struct {
	CatId       int64  `webcore:"catId"`
	Name        string `webcore:"name"`
	Description string `webcore:"description"`
	OrgId       int32  `webcore:"orgId"`
}

type DeleteControlDocCatInputs struct {
	CatId int64 `webcore:"catId"`
	OrgId int32 `webcore:"orgId"`
}

type UploadControlDocInputs struct {
	CatId              int64          `webcore:"catId"`
	OrgId              int32          `webcore:"orgId"`
	RelevantTime       time.Time      `webcore:"relevantTime"`
	AltName            string         `webcore:"altName"`
	Description        string         `webcore:"description"`
	UploadUserId       int64          `webcore:"uploadUserId"`
	FulfilledRequestId core.NullInt64 `webcore:"fulfilledRequestId,optional"`
}

type AllControlDocInputs struct {
	CatId int64 `webcore:"catId"`
	OrgId int32 `webcore:"orgId"`
}

type GetControlDocInputs struct {
	FileId int64 `webcore:"fileId"`
	OrgId  int32 `webcore:"orgId"`
}

type DeleteControlDocInputs struct {
	OrgId   int32   `webcore:"orgId"`
	FileIds []int64 `webcore:"fileIds"`
}

type DownloadControlDocInputs struct {
	FileId  int64 `webcore:"fileId"`
	OrgId   int32 `webcore:"orgId"`
	Version int32 `webcore:"version"`
	Preview bool  `webcore:"preview"`
}

type AllControlDocCatInputs struct {
	OrgId int32 `webcore:"orgId"`
}

type GetControlDocCatInputs struct {
	OrgId int32 `webcore:"orgId"`
	CatId int64 `webcore:"catId"`
	Lean  bool  `webcore:"lean"`
}

type EditControlDocInputs struct {
	FileId       int64     `json:"fileId"`
	OrgId        int32     `json:"orgId"`
	RelevantTime time.Time `json:"relevantTime"`
	AltName      string    `json:"altName"`
	Description  string    `json:"description"`
}

func newControlDocumentationCategory(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewControlDocCatInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newCat := core.ControlDocumentationCategory{
		Name:        inputs.Name,
		Description: inputs.Description,
		OrgId:       org.Id,
	}

	err = database.NewControlDocumentationCategory(&newCat, role)
	if err != nil {
		core.Warning("Failed to create doc cat: " + err.Error())
		if database.IsDuplicateDBEntry(err) {
			w.WriteHeader(http.StatusBadRequest)
			jsonWriter.Encode(database.DuplicateEntryJson)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	jsonWriter.Encode(newCat)
}

func editControlDocumentationCategory(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := EditControlDocCatInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	editCat := core.ControlDocumentationCategory{
		Id:          inputs.CatId,
		Name:        inputs.Name,
		Description: inputs.Description,
		OrgId:       inputs.OrgId,
	}

	err = database.EditControlDocumentationCategory(&editCat, role)
	if err != nil {
		core.Warning("Failed to edit doc cat: " + err.Error())
		if database.IsDuplicateDBEntry(err) {
			w.WriteHeader(http.StatusBadRequest)
			jsonWriter.Encode(database.DuplicateEntryJson)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	jsonWriter.Encode(editCat)
}

func deleteControlDocumentationCategory(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := DeleteControlDocCatInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("No organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = database.DeleteControlDocumentationCategory(inputs.CatId, org.Id, role)
	if err != nil {
		core.Warning("Failed to delete doc cat: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct{}{})
}

func uploadControlDocumentation(w http.ResponseWriter, r *http.Request) {
	// TODO: Do we want to offload this onto a message queue and have a
	// separate process handle this?
	// Steps for uploading control documentation
	// 	1) Create temporary entry in database for the file
	// 	2) Create new transit engine key for the file.
	// 	3) Use transit engine to encrypt file.
	// 	4) Send file to GCloud.
	//  5) Finalize details of the file in the database
	// 	6) Return confirmation of file upload to requester.
	storage := gcloud.DefaultGCloudApi.GetStorageApi()

	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := UploadControlDocInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("Can't find organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil || !role.Permissions.HasAccess(core.ResourceControlDocumentation, core.AccessManage) {
		core.Warning("Bad access: " + core.ErrorString(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		core.Warning("Can't find uploaded file: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	if fileHeader.Size > webcore.MaxFileSizeBytes {
		core.Warning("File too large.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, file)
	if err != nil {
		core.Warning("Could not read file: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	internalFile := core.ControlDocumentationFile{
		RelevantTime: inputs.RelevantTime,
		CategoryId:   inputs.CatId,
		OrgId:        org.Id,
		AltName:      inputs.AltName,
		Description:  inputs.Description,
	}

	bucket := core.EnvConfig.Gcloud.DocBucket

	tx := database.CreateTx()
	storageData, err := webcore.UploadNewFileWithTx(
		storage,
		&internalFile,
		bucket,
		nil,
		fileHeader.Filename,
		buffer.Bytes(),
		role,
		inputs.UploadUserId,
		tx,
		org,
		false, // useExistingMetadata
		true)  // addToFileVersion
	if err != nil {
		core.Warning("Failed to upload new file: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	storageFilename := internalFile.StorageFilename(org)

	// At this point we know we can put in a request to generate a preview.
	webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
		Exchange: webcore.DEFAULT_EXCHANGE,
		Queue:    webcore.FILE_PREVIEW_QUEUE,
		Body: webcore.FilePreviewMessage{
			File:    internalFile,
			Storage: *storageData,
		},
	})

	if inputs.FulfilledRequestId.NullInt64.Valid {
		err = database.FulfillDocumentRequestWithTx(
			inputs.FulfilledRequestId.NullInt64.Int64,
			internalFile.Id,
			inputs.OrgId,
			role,
			tx)
		if err != nil {
			tx.Rollback()
			storage.Delete(bucket, storageFilename)

			core.Warning("Failed to fulfill request in db: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		core.Warning("Failed to commit file: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(internalFile)
}

func allControlDocumentation(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllControlDocInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("Can't find organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type DataOutput struct {
		Files []*core.ControlDocumentationFile
	}
	output := DataOutput{}

	output.Files, err = database.GetControlDocumentationForCategory(inputs.CatId, org.Id, role)
	if err != nil {
		core.Warning("Can't get files: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(output)
}

func deleteControlDocumentation(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := DeleteControlDocInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromId(inputs.OrgId)
	if err != nil {
		core.Warning("Can't find organization: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, org.Id)
	if err != nil || !role.Permissions.HasAccess(core.ResourceControlDocumentation, core.AccessManage) {
		core.Warning("Bad access: " + core.ErrorString(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	storage := gcloud.DefaultGCloudApi.GetStorageApi()

	// TODO: We'll need to develop a utility that runs periodically to
	// clean up things that failed here.
	for _, id := range inputs.FileIds {
		versions, err := database.GetAllVersionsFileStorage(id, org.Id, role)
		if err != nil {
			core.Warning("Failed to get versions: " + err.Error())
			continue
		}

		for _, v := range versions {
			err = storage.Delete(v.BucketId, v.StorageId)
			if err != nil {
				core.Warning("Failed to delete control documentation: " + err.Error())
				continue
			}

			previewStorage, err := database.GetPreviewFileVersionStorageDataFromStorageData(v, role)
			if err != nil {
				core.Warning("Failed to get preview storage data: " + err.Error())
				continue
			}

			if previewStorage == nil {
				continue
			}

			err = storage.Delete(previewStorage.BucketId, previewStorage.StorageId)
			if err != nil {
				core.Warning("Failed to delete control documentation preview: " + err.Error())
				continue
			}
		}
	}

	err = database.DeleteBatchControlDocumentation(inputs.FileIds, org.Id, role)
	if err != nil {
		core.Warning("Can't delete database files: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct{}{})
}

func downloadControlDocumentation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")

	inputs := DownloadControlDocInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil || !role.Permissions.HasAccess(core.ResourceControlDocumentation, core.AccessView) {
		core.Warning("Bad access: " + core.ErrorString(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dbFile, err := database.GetControlDocumentation(inputs.FileId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Can't get file db data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var version *core.FileStorageData

	if inputs.Preview {
		version, err = database.GetPreviewFileVersionStorageData(dbFile.Id, dbFile.OrgId, inputs.Version, role)
	} else {
		version, err = database.GetFileVersionStorageData(dbFile.Id, dbFile.OrgId, inputs.Version, role)
	}

	if err != nil {
		core.Warning("Can't get file version data: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	storage := gcloud.DefaultGCloudApi.GetStorageApi()
	encryptedBytes, err := storage.Download(version.BucketId, version.StorageId)
	if err != nil {
		core.Warning("Can't get file from Backblaze: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	decryptedBytes, err := vault.TransitDecrypt(dbFile.UniqueKey(), encryptedBytes)
	if err != nil {
		core.Warning("Can't decrypt file: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(decryptedBytes)
}

func allControlDocumentationCategories(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllControlDocCatInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil || !role.Permissions.HasAccess(core.ResourceControlDocumentation, core.AccessView) {
		core.Warning("Bad access: " + core.ErrorString(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cats, err := database.GetAllDocumentationCategoriesForOrg(inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get all doc cats: " + core.ErrorString(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(cats)
}

func getControlDocumentationCategory(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetControlDocCatInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil || !role.Permissions.HasAccess(core.ResourceControlDocumentation, core.AccessView) {
		core.Warning("Bad access: " + core.ErrorString(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cat, err := database.GetDocumentationCategory(inputs.CatId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get doc cat: " + core.ErrorString(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var inputControls []*core.Control
	var outputControls []*core.Control

	if !inputs.Lean {
		inputControls, err = database.GetControlsWithInputDocumentationCategory(inputs.CatId, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get all input controls: " + core.ErrorString(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		outputControls, err = database.GetControlsWithOutputDocumentationCategory(inputs.CatId, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get all output controls: " + core.ErrorString(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	jsonWriter.Encode(struct {
		Cat       *core.ControlDocumentationCategory
		InputFor  []*core.Control
		OutputFor []*core.Control
	}{
		Cat:       cat,
		InputFor:  inputControls,
		OutputFor: outputControls,
	})
}

func getControlDocumentation(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetControlDocInputs{}
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

	file, err := database.GetControlDocumentation(inputs.FileId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get doc file: " + core.ErrorString(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	category, err := database.GetDocumentationCategory(file.CategoryId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get doc cat: " + core.ErrorString(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	versions, err := database.AllFileVersions(inputs.FileId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get file versions: " + core.ErrorString(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		File     *core.ControlDocumentationFile
		Category *core.ControlDocumentationCategory
		Versions []core.FileVersion
	}{
		File:     file,
		Category: category,
		Versions: versions,
	})
}

func editControlDocumentation(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := EditControlDocInputs{}
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

	file, err := database.GetControlDocumentation(inputs.FileId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get control documentation: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file.RelevantTime = inputs.RelevantTime
	file.AltName = inputs.AltName
	file.Description = inputs.Description

	err = database.UpdateControlDocumentation(file, role)
	if err != nil {
		core.Warning("Failed to edit control documentation: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	category, err := database.GetDocumentationCategory(file.CategoryId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get doc cat: " + core.ErrorString(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		File     *core.ControlDocumentationFile
		Category *core.ControlDocumentationCategory
	}{
		File:     file,
		Category: category,
	})
}

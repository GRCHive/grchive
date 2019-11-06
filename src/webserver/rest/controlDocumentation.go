package rest

import (
	"bytes"
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/backblaze_api"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/security"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"io"
	"net/http"
	"time"
)

type NewControlDocCatInputs struct {
	ControlId   int64  `webcore:"controlId"`
	Name        string `webcore:"name"`
	Description string `webcore:"description"`
}

type EditControlDocCatInputs struct {
	CatId       int64  `webcore:"catId"`
	ControlId   int64  `webcore:"controlId"`
	Name        string `webcore:"name"`
	Description string `webcore:"description"`
}

type DeleteControlDocCatInputs struct {
	CatId int64 `webcore:"catId"`
}

type UploadControlDocInputs struct {
	CatId        int64     `webcore:"catId"`
	RelevantTime time.Time `webcore:"relevantTime"`
}

type GetControlDocInputs struct {
	CatId     int64 `webcore:"catId"`
	Page      int   `webcore:"page"`
	NeedPages bool  `webcore:"needPages"`
}

type DeleteControlDocInputs struct {
	FileIds      []int64 `webcore:"fileIds"`
	OrgGroupName string  `webcore:"orgGroupName"`
}

type DownloadControlDocInputs struct {
	FileId int64 `webcore:"fileId"`
	OrgId  int32 `webcore:"orgId"`
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

	org, err := database.FindOrganizationFromControlId(inputs.ControlId, core.ServerRole)
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
		ControlId:   inputs.ControlId,
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

	org, err := database.FindOrganizationFromControlId(inputs.ControlId, core.ServerRole)
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
		ControlId:   inputs.ControlId,
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

	org, err := database.FindOrganizationFromDocCatId(inputs.CatId, core.ServerRole)
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

	err = database.DeleteControlDocumentationCategory(inputs.CatId, role)
	if err != nil {
		core.Warning("Failed to delete doc cat: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct{}{})
}

func uploadControlDocumentation(w http.ResponseWriter, r *http.Request) {
	b2Auth, err := backblaze.B2Auth(core.EnvConfig.Backblaze.Key)
	if err != nil {
		core.Warning("Could not auth with Backblaze: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Do we want to offload this onto a message queue and have a
	// separate process handle this?
	// Steps for uploading control documentation
	// 	1) Create temporary entry in database for the file
	// 	2) Create new transit engine key for the file.
	// 	3) Use transit engine to encrypt file.
	// 	4) Send file to Backblaze.
	//  5) Finalize details of the file in the database
	// 	6) Return confirmation of file upload to requester.

	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := UploadControlDocInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	org, err := database.FindOrganizationFromDocCatId(inputs.CatId, core.ServerRole)
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
		StorageName:  fileHeader.Filename,
		RelevantTime: inputs.RelevantTime,
		UploadTime:   time.Now().UTC(),
		CategoryId:   inputs.CatId,
	}

	tx := database.CreateTx()
	err = database.CreateControlDocumentationFileWithTx(&internalFile, tx, role)
	if err != nil {
		tx.Rollback()
		core.Warning("Could not create file: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	transitKey := internalFile.UniqueKey()
	err = security.TransitCreateNewEngineKey(transitKey)
	if err != nil {
		tx.Rollback()
		core.Warning("Could not create transit key: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	encryptedFile, err := security.TransitEncrypt(transitKey, buffer.Bytes())
	if err != nil {
		tx.Rollback()
		core.Warning("Could not encrypt file: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b2Filename := internalFile.StorageFilename(org)

	b2File, err := backblaze.UploadFile(b2Auth,
		core.EnvConfig.Backblaze.ControlDocBucketId,
		b2Filename,
		encryptedFile)
	if err != nil {
		tx.Rollback()
		core.Warning("Could not upload file: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	internalFile.BucketId = b2File.BucketId
	internalFile.StorageId = b2File.FileId
	err = database.UpdateControlDocumentation(&internalFile, tx, role)
	if err != nil {
		tx.Rollback()
		backblaze.DeleteFile(b2Auth, internalFile.StorageFilename(org), b2File)

		core.Warning("Failed to update control documentation: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		core.Warning("Failed to commit file: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(internalFile)
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

	org, err := database.FindOrganizationFromDocCatId(inputs.CatId, core.ServerRole)
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
		Files       []*core.ControlDocumentationFile
		TotalPages  int
		CurrentPage int
	}
	output := DataOutput{
		CurrentPage: inputs.Page,
	}

	const controlDocPageSize int = 10
	controlDocPageOffset := controlDocPageSize * inputs.Page

	output.Files, err = database.GetControlDocumentationForCategory(inputs.CatId, controlDocPageSize, controlDocPageOffset, role)
	if err != nil {
		core.Warning("Can't get files: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if inputs.NeedPages {
		output.TotalPages, err = database.GetTotalControlDocumentationPages(inputs.CatId, controlDocPageSize, role)
		if err != nil {
			core.Warning("Can't total pages: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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

	org, err := database.FindOrganizationFromGroupName(inputs.OrgGroupName)
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

	// If we fail in deleting anything from Backblaze just log and continue.
	// TODO: We'll need to develop a utility that runs periodically to
	// clean up things that failed here.
	b2Auth, err := backblaze.B2Auth(core.EnvConfig.Backblaze.Key)
	if err != nil {
		core.Warning("Could not auth with Backblaze: " + err.Error())
	} else {
		for _, id := range inputs.FileIds {
			file, err := database.GetControlDocumentation(id, org.Id, role)
			if err != nil {
				core.Warning("Failed to find control documentation: " + err.Error())
				continue
			}

			// Need to store actual storage filename on the database
			err = backblaze.DeleteFile(b2Auth, file.StorageFilename(org), backblaze.B2File{
				BucketId: file.BucketId,
				FileId:   file.StorageId,
			})

			if err != nil {
				core.Warning("Failed to delete control documentation: " + err.Error())
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

	b2Auth, err := backblaze.B2Auth(core.EnvConfig.Backblaze.Key)
	if err != nil {
		core.Warning("Could not auth with Backblaze: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inputs := DownloadControlDocInputs{}
	err = webcore.UnmarshalRequestForm(r, &inputs)
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

	encryptedBytes, err := backblaze.DownloadFile(b2Auth, backblaze.B2File{
		BucketId: dbFile.BucketId,
		FileId:   dbFile.StorageId,
	})
	if err != nil {
		core.Warning("Can't get file from Backblaze: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	decryptedBytes, err := security.TransitDecrypt(dbFile.UniqueKey(), encryptedBytes)
	if err != nil {
		core.Warning("Can't decrypt file: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(decryptedBytes)
}

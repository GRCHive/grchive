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

	newCat := core.ControlDocumentationCategory{
		Name:        inputs.Name,
		Description: inputs.Description,
		ControlId:   inputs.ControlId,
	}

	err = database.NewControlDocumentationCategory(&newCat)
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

	editCat := core.ControlDocumentationCategory{
		Id:          inputs.CatId,
		Name:        inputs.Name,
		Description: inputs.Description,
		ControlId:   inputs.ControlId,
	}

	err = database.EditControlDocumentationCategory(&editCat)
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

	err = database.DeleteControlDocumentationCategory(inputs.CatId)
	if err != nil {
		core.Warning("Failed to delete doc cat: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct{}{})
}

func uploadControlDocumentation(w http.ResponseWriter, r *http.Request) {
	parsedUserData, err := webcore.FindSessionParsedDataInContext(r.Context())
	if err != nil {
		core.Warning("Failed to find parsed user data: " + core.ErrorString(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		core.Warning("Can't find uploaded file: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	if fileHeader.Size > webcore.MaxFileSizeBytes {
		core.Warning("File too large." + err.Error())
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
	err = database.CreateControlDocumentationFileWithTx(&internalFile, tx)
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

	b2File, err := backblaze.UploadFile(b2Auth,
		core.EnvConfig.Backblaze.ControlDocBucketId,
		parsedUserData.Org.OktaGroupName+"/"+transitKey,
		encryptedFile)
	if err != nil {
		tx.Rollback()
		core.Warning("Could not upload file: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	internalFile.BucketId = b2File.BucketId
	internalFile.StorageId = b2File.FileId
	err = database.UpdateControlDocumentation(&internalFile, tx)
	if err != nil {
		tx.Rollback()
		backblaze.DeleteFile(b2File)

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

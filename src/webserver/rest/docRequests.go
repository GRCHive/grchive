package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
	"time"
)

type NewDocumentRequestInputs struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	CatId           int64  `json:"catId"`
	OrgId           int32  `json:"orgId"`
	RequestedUserId int64  `json:"requestedUserId"`
}

type GetDocumentRequestInputs struct {
	RequestId int64 `webcore:"requestId"`
	OrgId     int32 `webcore:"orgId"`
}

type DeleteDocumentRequestInputs struct {
	RequestId int64 `json:"requestId"`
	OrgId     int32 `json:"orgId"`
}

type CompleteDocumentRequestInputs struct {
	RequestId int64 `json:"requestId"`
	OrgId     int32 `json:"orgId"`
	Complete  bool  `json:"complete"`
}

type AllDocumentRequestsInputs struct {
	OrgId int32          `webcore:"orgId"`
	CatId core.NullInt64 `webcore:"catId,optional"`
}

type NewDocumentRequestCommentInputs struct {
	UserId    int64  `json:"userId"`
	Text      string `json:"text"`
	CatId     int64  `json:"catId"`
	OrgId     int32  `json:"orgId"`
	RequestId int64  `json:"requestId"`
}

type AllDocumentRequestCommentsInputs struct {
	RequestId int64 `webcore:"requestId"`
	CatId     int64 `webcore:"catId"`
	OrgId     int32 `webcore:"orgId"`
}

func newDocumentRequest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewDocumentRequestInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request := core.DocumentRequest{
		Name:            inputs.Name,
		Description:     inputs.Description,
		CatId:           inputs.CatId,
		OrgId:           inputs.OrgId,
		RequestedUserId: inputs.RequestedUserId,
		RequestTime:     time.Now().UTC(),
	}

	err = database.CreateNewDocumentRequest(&request, role)
	if err != nil {
		core.Warning("Failed to create new doc request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(request)
}

func getDocumentRequest(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := GetDocumentRequestInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	req, err := database.GetDocumentRequest(inputs.RequestId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get single doc request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cat, err := database.GetDocumentationCategory(req.CatId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get doc category: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileIds, err := database.GetFulfilledFileIdsForDocRequest(inputs.RequestId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get relevant file IDs for request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	files := make([]*core.ControlDocumentationFile, len(fileIds))
	for i, id := range fileIds {
		files[i], err = database.GetControlDocumentation(id, inputs.OrgId, role)
		if err != nil {
			core.Warning("Failed to get file metadata: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	jsonWriter.Encode(struct {
		Request  *core.DocumentRequest
		Files    []*core.ControlDocumentationFile
		Category *core.ControlDocumentationCategory
	}{
		Request:  req,
		Files:    files,
		Category: cat,
	})
}

func allDocumentRequests(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllDocumentRequestsInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var reqs []*core.DocumentRequest
	if inputs.CatId.NullInt64.Valid {
		reqs, err = database.GetAllDocumentRequestsForDocCat(inputs.CatId.NullInt64.Int64, inputs.OrgId, role)
	} else {
		reqs, err = database.GetAllDocumentRequestsForOrganization(inputs.OrgId, role)
	}

	if err != nil {
		core.Warning("Failed to get all doc requests for org: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(reqs)
}

func newDocumentRequestComment(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewDocumentRequestCommentInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	comment := core.DocumentRequestComment{
		UserId:    inputs.UserId,
		Text:      inputs.Text,
		PostTime:  time.Now().UTC(),
		CatId:     inputs.CatId,
		OrgId:     inputs.OrgId,
		RequestId: inputs.RequestId,
	}

	err = database.AddDocRequestComment(&comment, role)
	if err != nil {
		core.Warning("Failed to add doc request comment: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(comment)
}

func deleteDocumentRequest(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteDocumentRequestInputs{}
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

	err = database.DeleteDocumentRequest(inputs.RequestId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to delete doc request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func completeDocumentRequest(w http.ResponseWriter, r *http.Request) {
	inputs := CompleteDocumentRequestInputs{}
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

	err = database.CompleteDocumentRequest(inputs.RequestId, inputs.OrgId, inputs.Complete, role)
	if err != nil {
		core.Warning("Failed to complete/reopen doc request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func allDocumentRequestComments(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllDocumentRequestCommentsInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonWriter.Encode(struct{}{})
		return
	}

	role, err := webcore.GetCurrentRequestRole(r, inputs.OrgId)
	if err != nil {
		core.Warning("Bad access: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	comments, err := database.GetAllDocumentRequestComments(inputs.RequestId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get all doc request comments: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(comments)
}

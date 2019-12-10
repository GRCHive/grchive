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
	CatId     int64 `webcore:"catId"`
	OrgId     int32 `webcore:"orgId"`
}

type AllDocumentRequestsInputs struct {
	OrgId int32          `webcore:"orgId"`
	CatId core.NullInt64 `webcore:"catId"`
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

func NewDocumentRequest(w http.ResponseWriter, r *http.Request) {
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

func GetDocumentRequest(w http.ResponseWriter, r *http.Request) {
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

	req, err := database.GetDocumentRequest(inputs.RequestId, inputs.CatId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get single doc request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(struct {
		Request *core.DocumentRequest
	}{
		Request: req,
	})
}

func AllDocumentRequests(w http.ResponseWriter, r *http.Request) {
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

func NewDocumentRequestComment(w http.ResponseWriter, r *http.Request) {
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

func AllDocumentRequestComments(w http.ResponseWriter, r *http.Request) {
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

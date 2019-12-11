package rest

import (
	"encoding/json"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
	"time"
)

type GenericNewCommentInputs struct {
	UserId  int64  `json:"userId"`
	Content string `json:"content"`
}

type DocumentRequestNewCommentInputs struct {
	Comment   GenericNewCommentInputs `json:"comment"`
	RequestId int64                   `json:"requestId"`
	CatId     int64                   `json:"catId"`
	OrgId     int32                   `json:"orgId"`
}

type DocumentRequestAllCommentInputs struct {
	RequestId int64 `webcore:"requestId"`
	OrgId     int32 `webcore:"orgId"`
}

func commentFromInputs(inp GenericNewCommentInputs) *core.Comment {
	return &core.Comment{
		UserId:   inp.UserId,
		PostTime: time.Now().UTC(),
		Content:  inp.Content,
	}
}

func newDocumentRequestComment(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := DocumentRequestNewCommentInputs{}
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

	comment := commentFromInputs(inputs.Comment)
	err = database.InsertDocumentRequestComment(inputs.RequestId, inputs.CatId, inputs.OrgId, comment, role)
	if err != nil {
		core.Warning("Failed to insert doc request comments: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(comment)
}

func allDocumentRequestComments(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := DocumentRequestAllCommentInputs{}
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

	comments, err := database.GetDocumentRequestComments(inputs.RequestId, inputs.OrgId, role)
	if err != nil {
		core.Warning("Failed to get doc request comments: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(comments)
}

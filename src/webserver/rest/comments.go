package rest

import (
	"encoding/json"
	"errors"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"time"
)

type GenericNewCommentInputs struct {
	UserId  int64  `json:"userId"`
	Content string `json:"content"`
}

type NewCommentInputs struct {
	Comment GenericNewCommentInputs `json:"comment"`
	// Doc Request
	RequestId core.NullInt64 `json:"requestId"`
	CatId     core.NullInt64 `json:"catId"`
	// Document
	FileId core.NullInt64 `json:"fileId"`
	OrgId  int32          `json:"orgId"`
}

type AllCommentInputs struct {
	RequestId core.NullInt64 `webcore:"requestId,optional"`
	FileId    core.NullInt64 `webcore:"fileId,optional"`
	OrgId     int32          `webcore:"orgId"`
}

func commentFromInputs(inp GenericNewCommentInputs) *core.Comment {
	return &core.Comment{
		UserId:   inp.UserId,
		PostTime: time.Now().UTC(),
		Content:  inp.Content,
	}
}

func newComment(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := NewCommentInputs{}
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

	if inputs.RequestId.NullInt64.Valid && inputs.CatId.NullInt64.Valid {
		err = database.InsertDocumentRequestComment(
			inputs.RequestId.NullInt64.Int64,
			inputs.CatId.NullInt64.Int64,
			inputs.OrgId,
			comment,
			role)
	} else if inputs.FileId.NullInt64.Valid {
		err = database.InsertDocumentComment(
			inputs.FileId.NullInt64.Int64,
			inputs.OrgId,
			comment,
			role)
	} else {
		err = errors.New("Invalid comment type.")
	}

	if err != nil {
		core.Warning("Failed to insert comments: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(comment)
}

func allComments(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := AllCommentInputs{}
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

	var comments []*core.Comment
	if inputs.RequestId.NullInt64.Valid {
		comments, err = database.GetDocumentRequestComments(inputs.RequestId.NullInt64.Int64, inputs.OrgId, role)
	} else if inputs.FileId.NullInt64.Valid {
		comments, err = database.GetDocumentComments(inputs.FileId.NullInt64.Int64, inputs.OrgId, role)
	} else {
		err = errors.New("Invalid comment type.")
	}

	if err != nil {
		core.Warning("Failed to get comments: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(comments)
}

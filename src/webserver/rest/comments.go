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
	SqlRequestId core.NullInt64 `json:"sqlRequestId"`
	// Sql Request
	RequestId core.NullInt64 `json:"requestId"`
	// Document
	FileId core.NullInt64 `json:"fileId"`
	OrgId  int32          `json:"orgId"`

	CatId core.NullInt64 `json:"catId"`
}

type AllCommentInputs struct {
	SqlRequestId core.NullInt64 `webcore:"sqlRequestId,optional"`
	RequestId    core.NullInt64 `webcore:"requestId,optional"`
	FileId       core.NullInt64 `webcore:"fileId,optional"`
	OrgId        int32          `webcore:"orgId"`
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

	user, err := database.FindUserFromId(role.UserId)
	if err != nil {
		core.Warning("Failed to get user: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	comment := commentFromInputs(inputs.Comment)

	var object interface{} = nil
	if inputs.SqlRequestId.NullInt64.Valid {
		err = database.InsertSqlRequestComment(
			inputs.SqlRequestId.NullInt64.Int64,
			inputs.OrgId,
			comment,
			role)

		if err != nil {
			core.Warning("Failed to insert SQL request comments: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		object, err = database.GetSqlRequest(
			inputs.SqlRequestId.NullInt64.Int64,
			inputs.OrgId,
			role)
		if err != nil {
			core.Warning("Failed to get SQL request: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if inputs.RequestId.NullInt64.Valid {
		err = database.InsertDocumentRequestComment(
			inputs.RequestId.NullInt64.Int64,
			inputs.OrgId,
			comment,
			role)

		if err != nil {
			core.Warning("Failed to insert doc request comments: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		object, err = database.GetDocumentRequest(
			inputs.RequestId.NullInt64.Int64,
			inputs.OrgId,
			role)
		if err != nil {
			core.Warning("Failed to get doc request: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else if inputs.FileId.NullInt64.Valid {
		err = database.InsertDocumentComment(
			inputs.FileId.NullInt64.Int64,
			inputs.OrgId,
			comment,
			role)

		if err != nil {
			core.Warning("Failed to insert document comments: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		object, err = database.GetControlDocumentation(
			inputs.FileId.NullInt64.Int64,
			inputs.OrgId,
			role)
		if err != nil {
			core.Warning("Failed to get doc request: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else {
		core.Warning("Invalid combination of inputs for new comment")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
		Exchange: webcore.EVENT_EXCHANGE,
		Queue:    "",
		Body: webcore.EventMessage{
			Event: core.Event{
				Subject:        *user,
				Verb:           core.VerbComment,
				Object:         object,
				IndirectObject: nil,
				Timestamp:      time.Now().UTC(),
			},
		},
	})

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
	if inputs.SqlRequestId.NullInt64.Valid {
		comments, err = database.GetSqlRequestComments(inputs.SqlRequestId.NullInt64.Int64, inputs.OrgId, role)
	} else if inputs.RequestId.NullInt64.Valid {
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

type UpdateCommentInputs struct {
	CommentId int64  `json:"commentId"`
	Content   string `json:"content"`
}

func updateComment(w http.ResponseWriter, r *http.Request) {
	jsonWriter := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	inputs := UpdateCommentInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := webcore.GetUserIdFromApiRequestContext(r)
	if err != nil {
		core.Warning("Failed to obtain key user id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tx := database.CreateTx()

	comment := core.Comment{
		Id:       inputs.CommentId,
		UserId:   userId,
		PostTime: time.Now().UTC(),
		Content:  inputs.Content,
	}

	if err = database.UpdateCommentWithTx(&comment, tx); err != nil {
		tx.Rollback()
		core.Warning("Failed to edit comment: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		core.Warning("Failed to commit comment edit: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonWriter.Encode(comment)
}

type DeleteCommentInputs struct {
	CommentId int64 `json:"commentId"`
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	inputs := DeleteCommentInputs{}
	err := webcore.UnmarshalRequestForm(r, &inputs)
	if err != nil {
		core.Warning("Can't parse inputs: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId, err := webcore.GetUserIdFromApiRequestContext(r)
	if err != nil {
		core.Warning("Failed to obtain key user id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = database.DeleteComment(inputs.CommentId, userId)
	if err != nil {
		core.Warning("Failed to delete comment: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

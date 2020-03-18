package main

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"time"
)

type ControlEvent struct {
	User    core.User
	Control core.Control
}

func onNotifyControlOwnerChange(data string) error {
	parsedData := ControlEvent{}
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return err
	}

	if !parsedData.Control.OwnerId.NullInt64.Valid {
		return nil
	}

	assignedToUser, err := database.FindUserFromId(parsedData.Control.OwnerId.NullInt64.Int64)
	if err != nil {
		return err
	}

	webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
		Exchange: webcore.EVENT_EXCHANGE,
		Queue:    "",
		Body: webcore.EventMessage{
			Event: core.Event{
				Subject:        parsedData.User,
				Verb:           core.VerbAssign,
				Object:         assignedToUser,
				IndirectObject: parsedData.Control,
				Timestamp:      time.Now().UTC(),
			},
		},
	})

	return nil
}

type DocRequestEvent struct {
	User    core.User
	Request core.DocumentRequest
}

func onNotifyDocRequestAssigneeChange(data string) error {
	parsedData := DocRequestEvent{}
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return err
	}

	if !parsedData.Request.AssigneeUserId.NullInt64.Valid {
		return nil
	}

	assignedToUser, err := database.FindUserFromId(parsedData.Request.AssigneeUserId.NullInt64.Int64)
	if err != nil {
		return err
	}

	webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
		Exchange: webcore.EVENT_EXCHANGE,
		Queue:    "",
		Body: webcore.EventMessage{
			Event: core.Event{
				Subject:        parsedData.User,
				Verb:           core.VerbAssign,
				Object:         assignedToUser,
				IndirectObject: parsedData.Request,
				Timestamp:      time.Now().UTC(),
			},
		},
	})

	return nil
}

func onNotifyDocRequestStatusChange(data string) error {
	parsedData := DocRequestEvent{}
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return err
	}

	if parsedData.Request.CompletionTime.NullTime.Valid {
		webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
			Exchange: webcore.EVENT_EXCHANGE,
			Queue:    "",
			Body: webcore.EventMessage{
				Event: core.Event{
					Subject:        parsedData.User,
					Verb:           core.VerbComplete,
					Object:         parsedData.Request,
					IndirectObject: nil,
					Timestamp:      time.Now().UTC(),
				},
			},
		})
	} else {
		webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
			Exchange: webcore.EVENT_EXCHANGE,
			Queue:    "",
			Body: webcore.EventMessage{
				Event: core.Event{
					Subject:        parsedData.User,
					Verb:           core.VerbReopen,
					Object:         parsedData.Request,
					IndirectObject: nil,
					Timestamp:      time.Now().UTC(),
				},
			},
		})
	}

	return nil
}

type SqlRequestEvent struct {
	User    core.User
	Request core.DbSqlQueryRequest
}

func onNotifySqlRequestAssigneeChange(data string) error {
	parsedData := SqlRequestEvent{}
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return err
	}

	if !parsedData.Request.AssigneeUserId.NullInt64.Valid {
		return nil
	}

	assignedToUser, err := database.FindUserFromId(parsedData.Request.AssigneeUserId.NullInt64.Int64)
	if err != nil {
		return err
	}

	webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
		Exchange: webcore.EVENT_EXCHANGE,
		Queue:    "",
		Body: webcore.EventMessage{
			Event: core.Event{
				Subject:        parsedData.User,
				Verb:           core.VerbAssign,
				Object:         assignedToUser,
				IndirectObject: parsedData.Request,
				Timestamp:      time.Now().UTC(),
			},
		},
	})

	return nil
}

type SqlRequestApprovalEvent struct {
	User     core.User
	Approval core.DbSqlQueryRequestApproval
}

func onNotifySqlRequestApprovalChange(data string) error {
	parsedData := SqlRequestApprovalEvent{}
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return err
	}

	request, err := database.GetSqlRequest(parsedData.Approval.RequestId, parsedData.Approval.OrgId, core.ServerRole)
	if err != nil {
		return err
	}

	var verb core.EventVerb
	if parsedData.Approval.Response {
		verb = core.VerbApprove
	} else {
		verb = core.VerbReject
	}

	webcore.DefaultRabbitMQ.SendMessage(webcore.PublishMessage{
		Exchange: webcore.EVENT_EXCHANGE,
		Queue:    "",
		Body: webcore.EventMessage{
			Event: core.Event{
				Subject:        parsedData.User,
				Verb:           verb,
				Object:         request,
				IndirectObject: nil,
				Timestamp:      time.Now().UTC(),
			},
		},
	})

	return nil
}

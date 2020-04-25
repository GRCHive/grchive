package main

import (
	"encoding/json"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"time"
)

type ControlEvent struct {
	User       core.User
	Control    core.Control
	OldControl core.Control
}

func onNotifyControlOwnerChange(data string) error {
	parsedData := ControlEvent{}
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", parsedData)

	if parsedData.Control.OwnerId.NullInt64.Valid {
		assignedToUser, err := database.FindUserFromId(parsedData.Control.OwnerId.NullInt64.Int64)
		if err != nil {
			return err
		}

		webcore.SendEventToRabbitMQ(core.Event{
			Subject:        parsedData.User,
			Verb:           core.VerbAssign,
			Object:         parsedData.Control,
			IndirectObject: assignedToUser,
			Timestamp:      time.Now().UTC(),
		})
	}

	if parsedData.OldControl.OwnerId.NullInt64.Valid && parsedData.Control.OwnerId != parsedData.OldControl.OwnerId {
		assignedToUser, err := database.FindUserFromId(parsedData.OldControl.OwnerId.NullInt64.Int64)
		if err != nil {
			return err
		}

		webcore.SendEventToRabbitMQ(core.Event{
			Subject:        parsedData.User,
			Verb:           core.VerbUnassign,
			Object:         parsedData.Control,
			IndirectObject: assignedToUser,
			Timestamp:      time.Now().UTC(),
		})
	}

	return nil
}

type DocRequestEvent struct {
	User       core.User
	Request    core.DocumentRequest
	OldRequest core.DocumentRequest
}

func onNotifyDocRequestAssigneeChange(data string) error {
	parsedData := DocRequestEvent{}
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return err
	}

	if parsedData.Request.AssigneeUserId.NullInt64.Valid {
		assignedToUser, err := database.FindUserFromId(parsedData.Request.AssigneeUserId.NullInt64.Int64)
		if err != nil {
			return err
		}

		webcore.SendEventToRabbitMQ(core.Event{
			Subject:        parsedData.User,
			Verb:           core.VerbAssign,
			Object:         parsedData.Request,
			IndirectObject: assignedToUser,
			Timestamp:      time.Now().UTC(),
		})
	}

	if parsedData.OldRequest.AssigneeUserId.NullInt64.Valid && parsedData.Request.AssigneeUserId != parsedData.OldRequest.AssigneeUserId {
		assignedToUser, err := database.FindUserFromId(parsedData.OldRequest.AssigneeUserId.NullInt64.Int64)
		if err != nil {
			return err
		}

		webcore.SendEventToRabbitMQ(core.Event{
			Subject:        parsedData.User,
			Verb:           core.VerbUnassign,
			Object:         parsedData.Request,
			IndirectObject: assignedToUser,
			Timestamp:      time.Now().UTC(),
		})
	}

	return nil
}

func onNotifyDocRequestStatusChange(data string) error {
	parsedData := DocRequestEvent{}
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return err
	}

	if parsedData.Request.CompletionTime.NullTime.Valid {
		webcore.SendEventToRabbitMQ(core.Event{
			Subject:        parsedData.User,
			Verb:           core.VerbComplete,
			Object:         parsedData.Request,
			IndirectObject: nil,
			Timestamp:      time.Now().UTC(),
		})
	} else {
		webcore.SendEventToRabbitMQ(core.Event{
			Subject:        parsedData.User,
			Verb:           core.VerbReopen,
			Object:         parsedData.Request,
			IndirectObject: nil,
			Timestamp:      time.Now().UTC(),
		})
	}

	return nil
}

type SqlRequestEvent struct {
	User       core.User
	Request    core.DbSqlQueryRequest
	OldRequest core.DbSqlQueryRequest
}

func onNotifySqlRequestAssigneeChange(data string) error {
	parsedData := SqlRequestEvent{}
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return err
	}

	if parsedData.Request.AssigneeUserId.NullInt64.Valid {
		assignedToUser, err := database.FindUserFromId(parsedData.Request.AssigneeUserId.NullInt64.Int64)
		if err != nil {
			return err
		}

		webcore.SendEventToRabbitMQ(core.Event{
			Subject:        parsedData.User,
			Verb:           core.VerbAssign,
			Object:         parsedData.Request,
			IndirectObject: assignedToUser,
			Timestamp:      time.Now().UTC(),
		})
	}

	if parsedData.OldRequest.AssigneeUserId.NullInt64.Valid && parsedData.Request.AssigneeUserId != parsedData.OldRequest.AssigneeUserId {
		assignedToUser, err := database.FindUserFromId(parsedData.OldRequest.AssigneeUserId.NullInt64.Int64)
		if err != nil {
			return err
		}

		webcore.SendEventToRabbitMQ(core.Event{
			Subject:        parsedData.User,
			Verb:           core.VerbUnassign,
			Object:         parsedData.Request,
			IndirectObject: assignedToUser,
			Timestamp:      time.Now().UTC(),
		})
	}

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

	webcore.SendEventToRabbitMQ(core.Event{
		Subject:        parsedData.User,
		Verb:           verb,
		Object:         request,
		IndirectObject: nil,
		Timestamp:      time.Now().UTC(),
	})

	return nil
}

type GenericRequestEvent struct {
	User       core.User
	Request    core.GenericRequest
	OldRequest core.GenericRequest
}

func onNotifyGenericRequestAssigneeChange(data string) error {
	parsedData := GenericRequestEvent{}
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return err
	}

	if parsedData.Request.Assignee.NullInt64.Valid {
		assignedToUser, err := database.FindUserFromId(parsedData.Request.Assignee.NullInt64.Int64)
		if err != nil {
			return err
		}

		webcore.SendEventToRabbitMQ(core.Event{
			Subject:        parsedData.User,
			Verb:           core.VerbAssign,
			Object:         parsedData.Request,
			IndirectObject: assignedToUser,
			Timestamp:      time.Now().UTC(),
		})
	}

	if parsedData.OldRequest.Assignee.NullInt64.Valid && parsedData.Request.Assignee != parsedData.OldRequest.Assignee {
		assignedToUser, err := database.FindUserFromId(parsedData.OldRequest.Assignee.NullInt64.Int64)
		if err != nil {
			return err
		}

		webcore.SendEventToRabbitMQ(core.Event{
			Subject:        parsedData.User,
			Verb:           core.VerbUnassign,
			Object:         parsedData.Request,
			IndirectObject: assignedToUser,
			Timestamp:      time.Now().UTC(),
		})
	}

	return nil
}

type GenericRequestApprovalEvent struct {
	User     core.User
	Approval core.GenericApproval
}

func onNotifyGenericRequestApprovalChange(data string) error {
	parsedData := GenericRequestApprovalEvent{}
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return err
	}

	request, err := database.GetGenericRequestFromId(parsedData.Approval.RequestId)
	if err != nil {
		return err
	}

	var verb core.EventVerb
	if parsedData.Approval.Response {
		verb = core.VerbApprove
	} else {
		verb = core.VerbReject
	}

	webcore.SendEventToRabbitMQ(core.Event{
		Subject:        parsedData.User,
		Verb:           verb,
		Object:         request,
		IndirectObject: nil,
		Timestamp:      time.Now().UTC(),
	})

	return nil
}

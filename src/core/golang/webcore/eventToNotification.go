package webcore

import (
	"gitlab.com/grchive/grchive/core"
)

func CreateNotificationFromEvent(event *core.Event) (*core.Notification, error) {
	objectType, objectId, err := core.GetResourceTypeId(event.Object)
	if err != nil {
		return nil, err
	}

	indirectObjectType, indirectObjectId, err := core.GetResourceTypeId(event.IndirectObject)
	if err != nil {
		return nil, err
	}

	orgId, err := GetOrgIdFromResource(event.Object)
	if err != nil {
		return nil, err
	}

	notification := core.Notification{
		OrgId:              orgId,
		Time:               event.Timestamp,
		SubjectType:        core.ResourceIdUser,
		SubjectId:          event.Subject.Id,
		Verb:               string(event.Verb),
		ObjectType:         objectType,
		ObjectId:           objectId,
		IndirectObjectType: indirectObjectType,
		IndirectObjectId:   indirectObjectId,
	}

	return &notification, nil
}

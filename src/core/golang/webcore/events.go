package webcore

import (
	"gitlab.com/grchive/grchive/core"
)

func FindRelevantUsersForEvent(event *core.Event) ([]*core.User, error) {
	users := make([]*core.User, 0)
	userFound := map[int64]bool{}

	addUser := func(u *core.User) {
		if u == nil {
			return
		}

		_, ok := userFound[u.Id]
		if ok {
			return
		}

		users = append(users, u)
		userFound[u.Id] = true
	}

	// Assign: USER assigned USER to RESOURCE
	// Complete: USER completed RESOURCE
	// Reopen: USER reopened RESOURCE
	// Comment: USER commented on RESOURCE
	// In each case, if the object/indirect object is a USER then that
	// user is relevant. If the object/indirect object is not a USER then
	// we need to add all "relevant" users.
	isComment := (event.Verb == core.VerbComment)

	objectUsers, err := FindRelevantUsersForResource(event.Object, isComment)
	if err != nil {
		return nil, err
	}

	for _, u := range objectUsers {
		addUser(u)
	}

	indirectUsers, err := FindRelevantUsersForResource(event.IndirectObject, isComment)
	if err != nil {
		return nil, err
	}

	for _, u := range indirectUsers {
		addUser(u)
	}

	return users, nil
}

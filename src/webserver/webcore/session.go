package webcore

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
)

func ExtractParsedDataFromSession(session *core.UserSession) (*core.UserSessionParsedData, error) {
	user, org, err := database.FindUserFromIdWithOrganization(session.UserId)
	if err != nil {
		return nil, err
	}

	data := &core.UserSessionParsedData{
		Org:         org,
		CurrentUser: user,
	}
	return data, nil
}

package webcore

import (
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
)

func ExtractParsedDataFromSession(session *core.UserSession) (*core.UserSessionParsedData, error) {
	user, err := database.FindUserFromId(session.UserId)
	if err != nil {
		return nil, err
	}

	veri, err := database.IsUserVerified(user.Id)
	if err != nil {
		return nil, err
	}

	accessibleOrgIds, err := database.FindAccessibleOrganizationIdsForUser(user)
	if err != nil {
		return nil, err
	}

	data := &core.UserSessionParsedData{
		CurrentUser:    user,
		AccessibleOrgs: accessibleOrgIds,
		VerifiedEmail:  veri,
	}
	return data, nil
}

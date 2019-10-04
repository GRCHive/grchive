package webcore

import (
	"errors"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
)

func ExtractParsedDataFromSession(session *core.UserSession) (*core.UserSessionParsedData, error) {
	// We should do all the JWT parsing here so that we don't have to keep
	// decoding and parsing the JWT.
	idJwt, err := ReadRawJWTFromString(session.IdToken)
	if err != nil {
		return nil, err
	}

	if len(idJwt.Payload.Groups) != 1 {
		// Having more than one group is weird..having 0 groups is an error.
		if len(idJwt.Payload.Groups) > 1 {
			core.Warning("ID Token with more than one group: " + session.SessionId)
		} else {
			return nil, errors.New("Zero groups")
		}
	}

	groupName := idJwt.Payload.Groups[0]
	org, err := database.FindOrganizationFromGroupName(groupName)
	if err != nil {
		return nil, err
	}

	user := &core.User{
		FirstName: idJwt.Payload.FirstName,
		LastName:  idJwt.Payload.LastName,
		FullName:  idJwt.Payload.Name,
		Email:     idJwt.Payload.Email,
		ParentOrg: org,
	}

	data := &core.UserSessionParsedData{
		Org:         org,
		CurrentUser: user,
	}
	return data, nil
}

package webcore

import (
	"errors"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
)

func GetOrgIdFromResource(in interface{}) (int32, error) {
	if in == nil {
		return -1, nil
	}

	switch v := in.(type) {
	case core.Database:
		return v.OrgId, nil
	case core.DatabaseConnection:
		return v.OrgId, nil
	case core.DbSqlQueryMetadata:
		return v.OrgId, nil
	case core.DbSqlQuery:
		return v.OrgId, nil
	case core.DbSqlQueryRequest:
		return v.OrgId, nil
	case core.ControlDocumentationFile:
		return v.OrgId, nil
	case core.Vendor:
		return v.OrgId, nil
	case core.VendorProduct:
		return v.OrgId, nil
	case core.FileStorageData:
		return v.OrgId, nil
	case core.GeneralLedgerCategory:
		return v.OrgId, nil
	case core.GeneralLedgerAccount:
		return v.OrgId, nil
	case core.Server:
		return v.OrgId, nil
	case core.System:
		return v.OrgId, nil
	case core.DocumentRequest:
		return v.OrgId, nil
	case core.ControlDocumentationCategory:
		return v.OrgId, nil
	case core.Risk:
		return v.OrgId, nil
	case core.Control:
		return v.OrgId, nil
	}

	return 0, errors.New("Unsupported resource (GetOrgIdFromResource).")
}

func FindRelevantUsersForResource(in interface{}, commentThread bool) ([]*core.User, error) {
	var err error
	users := make([]*core.User, 0)

	if commentThread {
		threadId := int64(-1)

		switch v := in.(type) {
		case core.ControlDocumentationFile:
			threadId, err = database.GetDocumentCommentThreadId(v.Id, v.OrgId, core.ServerRole)
			if err != nil {
				return nil, err
			}
		case core.DocumentRequest:
			threadId, err = database.GetDocumentRequestCommentThreadId(v.Id, v.OrgId, core.ServerRole)
			if err != nil {
				return nil, err
			}
		case core.DbSqlQueryRequest:
			threadId, err = database.GetSqlRequestCommentThreadId(v.Id, v.OrgId, core.ServerRole)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("Resource does not support coments.")
		}

		users, err = database.FindUsersInCommentThread(threadId)
		if err != nil {
			return nil, err
		}

	} else {
		switch v := in.(type) {
		case core.User:
			users = append(users, &v)
		case core.Control:
			if v.OwnerId.NullInt64.Valid {
				u, err := database.FindUserFromId(v.OwnerId.NullInt64.Int64)
				if err != nil {
					return nil, err
				}
				users = append(users, u)
			}
		case core.DocumentRequest:
			u1, err := database.FindUserFromId(v.RequestedUserId)
			if err != nil {
				return nil, err
			}
			users = append(users, u1)

			if v.AssigneeUserId.NullInt64.Valid {
				u2, err := database.FindUserFromId(v.AssigneeUserId.NullInt64.Int64)
				if err != nil {
					return nil, err
				}
				users = append(users, u2)
			}
		case core.DbSqlQueryRequest:
			u1, err := database.FindUserFromId(v.UploadUserId)
			if err != nil {
				return nil, err
			}
			users = append(users, u1)

			if v.AssigneeUserId.NullInt64.Valid {
				u2, err := database.FindUserFromId(v.AssigneeUserId.NullInt64.Int64)
				if err != nil {
					return nil, err
				}
				users = append(users, u2)
			}
		default:
			return make([]*core.User, 0), nil
		}
	}

	return users, nil
}

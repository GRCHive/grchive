package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func InsertInviteCode(code *core.InviteCode, role *core.Role) (string, error) {
	if !role.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessManage) {
		return "", core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()

	rows, err := tx.NamedQuery(`
		INSERT INTO invitation_codes (from_user_id, from_org_id, to_email, sent_time)
		VALUES (:from_user_id, :from_org_id, :to_email, :sent_time)
		RETURNING id
	`, code)

	if err != nil {
		tx.Rollback()
		return "", err
	}

	rows.Next()
	err = rows.Scan(&code.Id)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	rows.Close()

	hash, err := core.HashId(code.Id)
	if err != nil {
		return "", err
	}

	return hash, tx.Commit()
}

func FindInviteCodeFromHash(hash string, role *core.Role) (*core.InviteCode, error) {
	if !role.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	id, err := core.ReverseHashId(hash)
	if err != nil {
		return nil, err
	}

	code := core.InviteCode{}
	err = dbConn.Get(&code, `
		SELECT *
		FROM invitation_codes
		WHERE id = $1
	`, id)

	return &code, nil
}

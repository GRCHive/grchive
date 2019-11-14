package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"time"
)

func InsertInviteCodeWithTx(code *core.InviteCode, role *core.Role, tx *sqlx.Tx) (string, error) {
	if !role.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessManage) {
		return "", core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO invitation_codes (from_user_id, from_org_id, to_email, sent_time)
		VALUES (:from_user_id, :from_org_id, :to_email, :sent_time)
		RETURNING id
	`, code)

	if err != nil {
		return "", err
	}

	rows.Next()
	err = rows.Scan(&code.Id)
	if err != nil {
		return "", err
	}
	rows.Close()

	hash, err := core.HashId(code.Id)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func FindInviteCodeFromHash(hash string, email string, role *core.Role) (*core.InviteCode, error) {
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
			AND to_email = $2
			AND used_time IS NULL
	`, id, email)

	if err != nil {
		return nil, err
	}

	return &code, nil
}

func MarkInviteAsUsedWithTx(code *core.InviteCode, tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		UPDATE invitation_codes
		SET used_time = $1
		WHERE id = $2
	`, time.Now().UTC(), code.Id)
	if err != nil {
		return err
	}
	return nil
}

func MarkInviteAsUsed(code *core.InviteCode) error {
	tx := dbConn.MustBegin()
	err := MarkInviteAsUsedWithTx(code, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func FindAllUsersInOrganization(orgId int32, role *core.Role) ([]*core.UserWithRole, error) {
	if !role.Permissions.HasAccess(core.ResourceOrgUsers, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	users := make([]*core.UserWithRole, 0)

	err := dbConn.Select(&users, `
		SELECT 
			u.*,
			ur.role_id as "role_id",
			$1 as "org_id"
		FROM users AS u
		INNER JOIN user_orgs AS uo
			ON u.id = uo.user_id
		INNER JOIN user_roles AS ur
			ON ur.user_id = u.id
				AND ur.org_id = $1
		WHERE uo.org_id = $1
	`, orgId)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func CreateNewUserWithTx(user *core.User, tx *sqlx.Tx) error {
	rows, err := tx.NamedQuery(`
		INSERT INTO users (first_name, last_name, email)
		VALUES (:first_name, :last_name, :email)
		RETURNING id
	`, user)
	if err != nil {
		return err
	}

	rows.Next()
	err = rows.Scan(&user.Id)
	if err != nil {
		return err
	}
	rows.Close()

	return nil
}

func CreateNewUser(user *core.User) error {
	tx := dbConn.MustBegin()
	err := CreateNewUserWithTx(user, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func FindUserFromId(id int64) (*core.User, error) {
	user := core.User{}

	err := dbConn.Get(&user, `
		SELECT * FROM users
		WHERE id = $1
	`, id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUserFromEmail(email string) (*core.User, error) {
	user := core.User{}

	err := dbConn.Get(&user, `
		SELECT * FROM users
		WHERE email = $1
	`, email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUserFromEmail(user *core.User) error {
	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		UPDATE users
		SET first_name = :first_name,
			last_name = :last_name
		WHERE email = :email
		RETURNING *
	`, user)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.StructScan(user)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Close()
	return tx.Commit()
}

func AddUserToOrganizationWithTx(user *core.User, org *core.Organization, tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		INSERT INTO user_orgs (user_id, org_id)
		VALUES ($1, $2)
	`, user.Id, org.Id)

	if err != nil {
		return err
	}
	return nil
}

func AddUserToOrganization(user *core.User, org *core.Organization) error {
	tx := dbConn.MustBegin()
	err := AddUserToOrganizationWithTx(user, org, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func FindAccessibleOrganizationIdsForUser(user *core.User) ([]int32, error) {
	orgIds := make([]int32, 0)

	err := dbConn.Select(&orgIds, `
		SELECT org_id
		FROM user_orgs
		WHERE user_id = $1
	`, user.Id)

	if err != nil {
		return nil, err
	}

	return orgIds, nil
}

func FindAccessibleOrganizationsForUser(userId int64) ([]*core.Organization, error) {
	orgs := make([]*core.Organization, 0)

	err := dbConn.Select(&orgs, `
		SELECT org.*
		FROM user_orgs AS uo
		INNER JOIN organizations AS org
			ON org.id = uo.org_id
		WHERE uo.user_id = $1
	`, userId)

	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func IsUserEmailInOrganization(email string, orgId int32) (bool, error) {
	rows, err := dbConn.Queryx(`
		SELECT *
		FROM users AS u
		INNER JOIN user_orgs AS uo
			ON u.id = uo.user_id
		WHERE u.email = $1
			AND uo.org_id = $2
	`, email, orgId)

	if err != nil {
		return false, nil
	}

	defer rows.Close()
	return rows.Next(), nil
}

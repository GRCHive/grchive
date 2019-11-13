package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func FindAllUsersInOrganization(orgId int32) ([]*core.User, error) {
	users := make([]*core.User, 0)

	err := dbConn.Select(&users, `
		SELECT u.*
		FROM users AS u
		INNER JOIN user_orgs AS uo
			ON u.id = uo.user_id
		WHERE uo.org_id = $1
	`, orgId)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func FindUserFromOktaId(oktaId string) (*core.User, error) {
	user := core.User{}

	err := dbConn.Get(&user, `
		SELECT * FROM users
		WHERE okta_id = $1
	`, oktaId)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateNewUser(user *core.User) error {
	var err error

	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		INSERT INTO users (okta_id, first_name, last_name, email)
		VALUES (:okta_id, :first_name, :last_name, :email)
		RETURNING id
	`, user)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&user.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()

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

func AddUserToOrganization(user *core.User, org *core.Organization) error {
	tx := dbConn.MustBegin()

	_, err := tx.Exec(`
		INSERT INTO user_orgs (user_id, org_id)
		VALUES ($1, $2)
	`, user.Id, org.Id)

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

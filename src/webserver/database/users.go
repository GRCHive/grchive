package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func FindAllUsersInOrganization(orgId int32) ([]*core.User, error) {
	users := make([]*core.User, 0)

	err := dbConn.Select(&users, `
		SELECT * FROM users
		WHERE org_id = $1
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
		INSERT INTO users (okta_id, org_id, first_name, last_name, email)
		VALUES (:okta_id, :org_id, :first_name, :last_name, :email)
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

func FindUserFromIdWithOrganization(id int64) (*core.User, *core.Organization, error) {
	type JointResult struct {
		User core.User         `db:"user"`
		Org  core.Organization `db:"org"`
	}

	result := JointResult{}
	err := dbConn.Get(&result, `
		SELECT 
			users.id AS "user.id",
			users.okta_id AS "user.okta_id",
			users.org_id AS "user.org_id",
			users.first_name AS "user.first_name",
			users.last_name AS "user.last_name",
			users.email AS "user.email",
			org.id AS "org.id",
			org.org_group_id AS "org.org_group_id",
			org.org_group_name AS "org.org_group_name",
			org.org_name AS "org.org_name",
			org.saml_iden AS "org.saml_iden"
		FROM users
		INNER JOIN organizations AS org
			ON org.id = users.org_id
		WHERE users.id = $1
	`, id)

	if err != nil {
		return nil, nil, err
	}

	return &result.User, &result.Org, nil
}

func UpdateUser(user *core.User) error {
	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		UPDATE users
		SET first_name = :first_name,
			last_name = :last_name
		WHERE id = :id
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

	return tx.Commit()
}

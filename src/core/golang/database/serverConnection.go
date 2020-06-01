package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func NewServerSSHPasswordConnectionWithTx(tx *sqlx.Tx, conn *core.ServerSSHPasswordConnection) error {
	rows, err := tx.NamedQuery(`
		INSERT INTO server_username_password_connection (server_id, org_id, username, password)
		VALUES (:server_id, :org_id, :username, :password)
		RETURNING id
	`, conn)

	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	return rows.Scan(&conn.Id)
}

func NewServerSSHKeyConnectionWithTx(tx *sqlx.Tx, conn *core.ServerSSHKeyConnection) error {
	rows, err := tx.NamedQuery(`
		INSERT INTO server_ssh_key_connection (server_id, org_id, username, private_key)
		VALUES (:server_id, :org_id, :username, :private_key)
		RETURNING id
	`, conn)

	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	return rows.Scan(&conn.Id)
}

func EditServerSSHPasswordConnectionWithTx(tx *sqlx.Tx, conn *core.ServerSSHPasswordConnection) error {
	_, err := tx.NamedExec(`
		UPDATE server_username_password_connection
		SET username = :username,
			password = :password
		WHERE id = :id
	`, conn)
	return err
}

func EditServerSSHKeyConnectionWithTx(tx *sqlx.Tx, conn *core.ServerSSHKeyConnection) error {
	_, err := tx.NamedExec(`
		UPDATE server_ssh_key_connection
		SET username = :username,
			private_key = :private_key
		WHERE id = :id
	`, conn)
	return err
}

func GetSSHPasswordConnectionForServer(serverId int64, orgId int32) (*core.ServerSSHPasswordConnection, error) {
	rows, err := dbConn.Queryx(`
		SELECT *
		FROM server_username_password_connection
		WHERE server_id = $1 AND org_id = $2
	`, serverId, orgId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	conn := core.ServerSSHPasswordConnection{}
	err = rows.StructScan(&conn)
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

func GetSSHKeyConnectionForServer(serverId int64, orgId int32) (*core.ServerSSHKeyConnection, error) {
	rows, err := dbConn.Queryx(`
		SELECT *
		FROM server_ssh_key_connection
		WHERE server_id = $1 AND org_id = $2
	`, serverId, orgId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	conn := core.ServerSSHKeyConnection{}
	err = rows.StructScan(&conn)
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

func GetSSHPasswordConnection(id int64) (*core.ServerSSHPasswordConnection, error) {
	conn := core.ServerSSHPasswordConnection{}
	err := dbConn.Get(&conn, `
		SELECT *
		FROM server_username_password_connection
		WHERE id = $1
	`, id)
	return &conn, err
}

func GetSSHKeyConnection(id int64) (*core.ServerSSHKeyConnection, error) {
	conn := core.ServerSSHKeyConnection{}
	err := dbConn.Get(&conn, `
		SELECT *
		FROM server_ssh_key_connection 
		WHERE id = $1
	`, id)
	return &conn, err
}

func DeleteSSHPasswordConnectionWithTx(tx *sqlx.Tx, id int64) error {
	_, err := tx.Exec(`
		DELETE FROM server_username_password_connection
		WHERE id = $1
	`, id)
	return err
}

func DeleteSSHKeyConnectionWithTx(tx *sqlx.Tx, id int64) error {
	_, err := tx.Exec(`
		DELETE FROM server_ssh_key_connection
		WHERE id = $1
	`, id)
	return err
}

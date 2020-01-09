package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func NewServer(server *core.Server, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceServers, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		INSERT INTO infrastructure_servers (org_id, name, description, ip_address, operating_system, location)
		VALUES (:org_id, :name, :description, :ip_address, :operating_system, :location)
		RETURNING id
	`, server)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&server.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows.Close()

	return tx.Commit()
}

func UpdateServer(server *core.Server, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceServers, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.NamedExec(`
		UPDATE infrastructure_servers
		SET name = :name,
			description = :description,
			ip_address = :ip_address,
			operating_system = :operating_system,
			location = :location
		WHERE id = :id AND org_id = :org_id
	`, server)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func DeleteServer(serverId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceServers, core.AccessManage) {
		return core.ErrorUnauthorized
	}
	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM infrastructure_servers
		WHERE id = $1 AND org_id = $2
	`, serverId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func AllServersForOrganization(orgId int32, role *core.Role) ([]*core.Server, error) {
	if !role.Permissions.HasAccess(core.ResourceServers, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	servers := make([]*core.Server, 0)
	err := dbConn.Select(&servers, `
		SELECT *
		FROM infrastructure_servers
		WHERE org_id = $1
	`, orgId)
	return servers, err
}

func GetServer(serverId int64, orgId int32, role *core.Role) (*core.Server, error) {
	if !role.Permissions.HasAccess(core.ResourceServers, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	server := core.Server{}
	err := dbConn.Get(&server, `
		SELECT *
		FROM infrastructure_servers
		WHERE id = $1 AND org_id = $2
	`, serverId, orgId)
	return &server, err
}

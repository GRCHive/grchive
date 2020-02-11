package database

import (
	"gitlab.com/grchive/grchive/core"
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
		rows.Close()
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

func AllServersForDeployment(deployId int64, orgId int32, role *core.Role) ([]*core.Server, error) {
	if !role.Permissions.HasAccess(core.ResourceServers, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	servers := make([]*core.Server, 0)
	err := dbConn.Select(&servers, `
		SELECT srv.*
		FROM infrastructure_servers AS srv
		INNER JOIN deployment_server_link AS link
			ON link.server_id = srv.id
		WHERE srv.org_id = $1
			AND link.deployment_id = $2
	`, orgId, deployId)
	return servers, err
}

func GetSystemsLinkedToServer(serverId int64, orgId int32, role *core.Role) ([]*core.System, error) {
	if !role.Permissions.HasAccess(core.ResourceServers, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceSystems, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	systems := make([]*core.System, 0)
	err := dbConn.Select(&systems, `
		SELECT sys.*
		FROM systems AS sys
		INNER JOIN deployment_system_link AS dlink
			ON dlink.system_id = sys.id
		INNER JOIN deployments AS deploy
			ON dlink.deployment_id = deploy.id
		INNER JOIN deployment_server_link AS slink
			ON dlink.deployment_id = slink.deployment_id
		WHERE slink.server_id = $1
			AND sys.org_id = $2
			AND deploy.deployment_type = $3
	`, serverId, orgId, core.KSelfDeployment)
	return systems, err
}

func GetDatabasesLinkedToServer(serverId int64, orgId int32, role *core.Role) ([]*core.Database, error) {
	if !role.Permissions.HasAccess(core.ResourceServers, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceDatabases, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	dbs := make([]*core.Database, 0)
	err := dbConn.Select(&dbs, `
		SELECT db.*
		FROM database_resources AS db
		INNER JOIN deployment_db_link AS dlink
			ON dlink.db_id = db.id
		INNER JOIN deployments AS deploy
			ON dlink.deployment_id = deploy.id
		INNER JOIN deployment_server_link AS slink
			ON dlink.deployment_id = slink.deployment_id
		WHERE slink.server_id = $1
			AND db.org_id = $2
			AND deploy.deployment_type = $3
	`, serverId, orgId, core.KSelfDeployment)
	return dbs, err
}

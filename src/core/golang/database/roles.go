package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

var resourceToDatabaseMap = map[core.ResourceType]string{
	core.ResourceOrgUsers:                     "resource_organization_users_access",
	core.ResourceOrgRoles:                     "resource_organization_roles_access",
	core.ResourceProcessFlows:                 "resource_process_flows_access",
	core.ResourceControls:                     "resource_controls_access",
	core.ResourceControlDocumentation:         "resource_control_documentation_access",
	core.ResourceControlDocumentationMetadata: "resource_control_documentation_metadata_access",
	core.ResourceRisks:                        "resource_risks_access",
	core.ResourceGeneralLedger:                "resource_gl_access",
	core.ResourceSystems:                      "resource_systems_access",
	core.ResourceDatabases:                    "resource_database_access",
	core.ResourceDbConnections:                "resource_db_conn_access",
	core.ResourceDocRequests:                  "resource_doc_request_access",
	core.ResourceDeployments:                  "resource_deployment_access",
	core.ResourceServers:                      "resource_server_access",
	core.ResourceVendors:                      "resource_vendor_access",
	core.ResourceDbSql:                        "resource_db_sql_access",
	core.ResourceDbSqlQuery:                   "resource_db_sql_query_access",
	core.ResourceDbSqlRequest:                 "resource_db_sql_requests_access",
	core.ResourceClientData:                   "resource_client_data_access",
	core.ResourceManagedCode:                  "resource_managed_code_access",
	core.ResourceClientScripts:                "resource_client_scripts_access",
	core.ResourceScriptRun:                    "resource_script_run_access",
	core.ResourceBuildLog:                     "resource_build_log_access",
	core.ResourceShell:                        "resource_shell_scripts_access",
	core.ResourceShellRun:                     "resource_shell_script_runs_access",
}

var resourceToColumnName = map[core.ResourceType]string{
	core.ResourceOrgUsers:                     "permissions.users_access",
	core.ResourceOrgRoles:                     "permissions.roles_access",
	core.ResourceProcessFlows:                 "permissions.flow_access",
	core.ResourceControls:                     "permissions.control_access",
	core.ResourceControlDocumentation:         "permissions.doc_access",
	core.ResourceControlDocumentationMetadata: "permissions.doc_meta_access",
	core.ResourceRisks:                        "permissions.risk_access",
	core.ResourceGeneralLedger:                "permissions.gl_access",
	core.ResourceSystems:                      "permissions.system_access",
	core.ResourceDatabases:                    "permissions.db_access",
	core.ResourceDbConnections:                "permissions.db_conn_access",
	core.ResourceDocRequests:                  "permissions.doc_request_access",
	core.ResourceDeployments:                  "permissions.deployment_access",
	core.ResourceServers:                      "permissions.server_access",
	core.ResourceVendors:                      "permissions.vendor_access",
	core.ResourceDbSql:                        "permissions.db_sql_access",
	core.ResourceDbSqlQuery:                   "permissions.db_sql_query_access",
	core.ResourceDbSqlRequest:                 "permissions.db_sql_requests_access",
	core.ResourceClientData:                   "permissions.client_data_access",
	core.ResourceManagedCode:                  "permissions.managed_code_access",
	core.ResourceClientScripts:                "permissions.client_script_access",
	core.ResourceScriptRun:                    "permissions.script_run_access",
	core.ResourceBuildLog:                     "permissions.build_log_access",
	core.ResourceShell:                        "permissions.shell_script_access",
	core.ResourceShellRun:                     "permissions.shell_script_runs_access",
}

func createRoleSql(cond string) string {
	selectColumns := `
		role.id AS "role.id",
		role.name AS "role.name",
		role.description AS "role.description",
		role.is_default_role AS "role.is_default_role",
		role.is_admin_role AS "role.is_admin_role",
		role.org_id AS "role.org_id"
	`

	joinStmts := ""

	for _, res := range core.AvailableResources {
		table := resourceToDatabaseMap[res]
		col := resourceToColumnName[res]

		selectColumns = selectColumns + fmt.Sprintf(`
			, %s.access_type AS "%s"
		`, table, col)

		joinStmts = joinStmts + fmt.Sprintf(`
			LEFT JOIN %s
				ON role.id = %s.role_id
		`, table, table)
	}

	return fmt.Sprintf(`
		SELECT
		%s
		FROM organization_available_roles AS role
		%s	
		%s
		`, selectColumns, joinStmts, cond)
}

func FindUserRoleFromStmt(stmt *sqlx.Stmt, args ...interface{}) (*core.Role, error) {
	rows, err := stmt.Queryx(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	role := &core.Role{}
	if !rows.Next() {
		return nil, nil
	}

	err = rows.StructScan(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// This returns nil, nil if no permissions was found for the user.
func FindUserRoleForOrg(userId int64, orgId int32, actionRole *core.Role) (*core.Role, error) {
	if !actionRole.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	stmt, err := dbConn.Preparex(createRoleSql(`
	INNER JOIN user_roles AS ur
		ON ur.role_id = role.id
	WHERE ur.user_id = $1
		AND ur.org_id = $2
	`))

	if err != nil {
		return nil, err
	}

	return FindUserRoleFromStmt(stmt, userId, orgId)
}

// This returns nil, nil if no permissions was found for the org.
func FindDefaultRoleForOrg(orgId int32, actionRole *core.Role) (*core.Role, error) {
	if !actionRole.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	stmt, err := dbConn.Preparex(createRoleSql(`
	WHERE role.org_id = $1
		AND role.is_default_role = 'true'
	`))
	if err != nil {
		return nil, err
	}

	return FindUserRoleFromStmt(stmt, orgId)
}

func FindAdminRoleForOrg(orgId int32, actionRole *core.Role) (*core.Role, error) {
	if !actionRole.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	stmt, err := dbConn.Preparex(createRoleSql(`
	WHERE role.org_id = $1
		AND role.is_admin_role = 'true'
	`))
	if err != nil {
		return nil, err
	}
	return FindUserRoleFromStmt(stmt, orgId)
}

func FindRoleFromId(roleId int64, orgId int32, actionRole *core.Role) (*core.Role, error) {
	if !actionRole.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	stmt, err := dbConn.Preparex(createRoleSql(`
	WHERE role.org_id = $1
		AND role.id = $2
	`))
	if err != nil {
		return nil, err
	}

	return FindUserRoleFromStmt(stmt, orgId, roleId)
}

func InsertOrgRole(metadata *core.RoleMetadata, role *core.Role, actionRole *core.Role) error {
	if !actionRole.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessManage) {
		return core.ErrorUnauthorized
	}
	tx := dbConn.MustBegin()

	rows, err := tx.NamedQuery(`
		INSERT INTO organization_available_roles ( is_default_role, is_admin_role, name, description, org_id)
		VALUES (
			:is_default_role,
			:is_admin_role,
			:name,
			:description,
			:org_id
		)
		RETURNING id
	`, metadata)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&role.Id)
	if err != nil {
		rows.Close()
		tx.Rollback()
		return err
	}
	rows.Close()

	for _, resource := range core.AvailableResources {
		_, err = tx.Exec(fmt.Sprintf(`
			INSERT INTO %s (role_id, org_id, access_type)
			VALUES ( $1, $2, $3)
		`, resourceToDatabaseMap[resource]), role.Id, metadata.OrgId, role.Permissions.GetAccessType(resource))

		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func InsertUserRoleForOrgWithTx(userId int64, orgId int32, roleId int64, actionRole *core.Role, tx *sqlx.Tx) error {
	if !actionRole.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	_, err := tx.Exec(`
		INSERT INTO user_roles (role_id, user_id, org_id)
		VALUES ( $1, $2, $3)
	`, roleId, userId, orgId)

	if err != nil {
		return err
	}

	return nil
}

func InsertUserRoleForOrg(userId int64, orgId int32, roleId int64, actionRole *core.Role) error {
	tx := dbConn.MustBegin()
	err := InsertUserRoleForOrgWithTx(userId, orgId, roleId, actionRole, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func FindRolesForOrg(orgId int32, actionRole *core.Role) ([]*core.RoleMetadata, error) {
	if !actionRole.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	stmt, err := dbConn.Preparex(createRoleSql(`
	WHERE role.org_id = $1
	`))

	if err != nil {
		return nil, err
	}

	roles := make([]*core.Role, 0)
	err = stmt.Select(&roles, orgId)
	if err != nil {
		return nil, err
	}

	metadata := make([]*core.RoleMetadata, len(roles))
	for i, r := range roles {
		metadata[i] = &r.RoleMetadata
	}
	return metadata, err
}

func FindUserIdsWithRole(roleId int64, orgId int32, actionRole *core.Role) ([]int64, error) {
	if !actionRole.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	userIds := make([]int64, 0)

	err := dbConn.Select(&userIds, `
		SELECT ur.user_id
		FROM user_roles AS ur
		WHERE ur.role_id = $1
			AND ur.org_id = $2
	`, roleId, orgId)
	return userIds, err
}

func UpdateRole(role *core.Role, actionRole *core.Role) error {
	if !actionRole.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	// We only need to update the metadata in the RETURNING
	// part because we assume that whatever the user passes in
	// as the permissions should be what gets spit out.
	tx := dbConn.MustBegin()
	rows, err := tx.NamedQuery(`
		UPDATE organization_available_roles
		SET name = :name,
			description = :description
		WHERE id = :id AND org_id = :org_id
		RETURNING *
	`, role.RoleMetadata)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.StructScan(&role.RoleMetadata)
	if err != nil {
		rows.Close()
		tx.Rollback()
		return err
	}

	rows.Close()

	for _, resource := range core.AvailableResources {
		_, err = tx.Exec(fmt.Sprintf(`
			UPDATE %s
			SET access_type = $1
			WHERE role_id = $2
		`, resourceToDatabaseMap[resource]), role.Permissions.GetAccessType(resource), role.Id)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func DeleteRoleMetadata(orgId int32, roleId int64, actionRole *core.Role) error {
	if !actionRole.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()

	_, err := tx.Exec(`
		DELETE FROM organization_available_roles AS role
		WHERE role.org_id = $1
			AND role.id = $2
			AND role.is_default_role = 'false'
			AND role.is_admin_role = 'false'
	`, orgId, roleId)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func AddUsersToRoleWithTx(userIds []int64, orgId int32, roleId int64, actionRole *core.Role, tx *sqlx.Tx) error {
	if !actionRole.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	for _, userId := range userIds {
		_, err := tx.Exec(`
			UPDATE user_roles
			SET role_id = $1
			WHERE user_id = $2
				AND org_id = $3
		`, roleId, userId, orgId)

		if err != nil {
			return err
		}
	}

	return nil
}

func AddUsersToRole(userIds []int64, orgId int32, roleId int64, actionRole *core.Role) error {
	tx := dbConn.MustBegin()
	err := AddUsersToRoleWithTx(userIds, orgId, roleId, actionRole, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

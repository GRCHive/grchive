package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

var resourceToDatabaseMap = map[core.ResourceType]string{
	core.ResourceOrgUsers:                     "resource_organization_users_access",
	core.ResourceOrgRoles:                     "resource_organization_roles_access",
	core.ResourceProcessFlows:                 "resource_process_flows_access",
	core.ResourceControls:                     "resource_controls_access",
	core.ResourceControlDocumentation:         "resource_control_documentation_access",
	core.ResourceControlDocumentationMetadata: "resource_control_documentation_metadata_access",
	core.ResourceRisks:                        "resource_risks_access",
}

func createRoleSql(cond string) string {
	return fmt.Sprintf(`
		SELECT
			role.id AS "role.id",
			role.name AS "role.name",
			role.description AS "role.description",
			role.is_default_role AS "role.is_default_role",
			role.is_admin_role AS "role.is_admin_role",
			role.org_id AS "role.org_id",
			ruser.access_type AS "permissions.users_access",
			rrole.access_type AS "permissions.roles_access",
			rpf.access_type AS "permissions.flow_access",
			rc.access_type AS "permissions.control_access",
			rcd.access_type AS "permissions.doc_access",
			rcdm.access_type AS "permissions.doc_meta_access",
			rr.access_type AS "permissions.risk_access"
		FROM organization_available_roles AS role
		INNER JOIN resource_organization_users_access AS ruser
			ON role.id = ruser.role_id
		INNER JOIN resource_organization_roles_access AS rrole
			ON role.id = rrole.role_id
		INNER JOIN resource_process_flows_access AS rpf
			ON role.id = rpf.role_id
		INNER JOIN resource_controls_access AS rc
			ON role.id = rc.role_id
		INNER JOIN resource_control_documentation_access AS rcd
			ON role.id = rcd.role_id
		INNER JOIN resource_control_documentation_metadata_access AS rcdm
			ON role.id = rcdm.role_id
		INNER JOIN resource_risks_access AS rr
			ON role.id = rr.role_id
		%s
		`, cond)
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
		tx.Rollback()
		return err
	}
	rows.Close()

	for _, resource := range core.AvailableResources {
		_, err = tx.Exec(fmt.Sprintf(`
			INSERT INTO %s (role_id, access_type)
			VALUES ( $1, $2 )
		`, resourceToDatabaseMap[resource]), role.Id, role.Permissions.GetAccessType(resource))

		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func InsertUserRoleForOrg(userId int64, orgId int32, role *core.Role, actionRole *core.Role) error {
	if !actionRole.Permissions.HasAccess(core.ResourceOrgRoles, core.AccessManage) {
		return core.ErrorUnauthorized
	}
	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		INSERT INTO user_roles (role_id, user_id, org_id)
		VALUES ( $1, $2, $3)
	`, role.Id, userId, orgId)

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

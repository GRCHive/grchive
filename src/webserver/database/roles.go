package database

import (
	"fmt"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

var resourceToDatabaseMap = map[core.ResourceType]string{
	core.ResourceOrgRoles:             "resource_organizations_access",
	core.ResourceProcessFlows:         "resource_process_flows_access",
	core.ResourceControls:             "resource_controls_access",
	core.ResourceControlDocumentation: "resource_control_documentation_access",
	core.ResourceRisks:                "resource_risks_access",
}

// This returns nil, nil if no permissions was found for the user.
func FindUserRoleForOrg(userId int64, orgId int32) (*core.Role, error) {
	rows, err := dbConn.Queryx(`
		SELECT
			rp.role_id AS id,
			rp.permissions AS permissions
		FROM user_roles AS ur
		INNER JOIN role_permissions AS rp
			ON ur.role_id = rp.role_id
		WHERE ur.user_id = $1
			AND ur.org_id = $2
	`, userId, orgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	role := core.Role{}
	err = rows.StructScan(&role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// This returns nil, nil if no permissions was found for the org.
func FindDefaultRoleForOrg(orgId int32) (*core.Role, error) {
	rows, err := dbConn.Queryx(`
		SELECT
			rp.role_id AS id,
			rp.permissions AS permissions
		FROM organization_available_roles AS org
		INNER JOIN role_permissions AS rp
			ON org.id = rp.role_id
		WHERE org.org_id = $1
			AND org.is_default_role = 'true'
	`, orgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	role := core.Role{}
	err = rows.StructScan(&role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func InsertOrgRole(metadata *core.RoleMetadata, role *core.Role) error {
	tx := dbConn.MustBegin()

	rows, err := tx.NamedQuery(`
		INSERT INTO organization_available_roles ( is_default_role, name, description, org_id)
		VALUES (
			:is_default_role,
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

func InsertUserRoleForOrg(userId int64, orgId int32, role *core.Role) error {
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

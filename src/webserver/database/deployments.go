package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func CreateNewEmptyDeploymentWithTx(orgId int32, role *core.Role, tx *sqlx.Tx) (*core.FullDeployment, error) {
	if !role.Permissions.HasAccess(core.ResourceDeployments, core.AccessManage) {
		return nil, core.ErrorUnauthorized
	}

	deployment := core.FullDeployment{
		OrgId: orgId,
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO deployments (org_id)
		VALUES (:org_id)
		RETURNING id
	`, deployment)
	if err != nil {
		return nil, err
	}

	rows.Next()
	err = rows.Scan(&deployment.Id)
	if err != nil {
		return nil, err
	}
	rows.Close()

	return &deployment, nil
}

func LinkDeploymentWithSystemWithTx(deployment *core.FullDeployment, systemId int64, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceSystems, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.Exec(`
		INSERT INTO deployment_system_link (system_id, deployment_id, org_id)
		VALUES ($1, $2, $3)
	`, systemId, deployment.Id, deployment.OrgId)
	return err
}

func LinkDeploymentWithDatabaseWithTx(deployment *core.FullDeployment, dbId int64, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	_, err := tx.Exec(`
		INSERT INTO deployment_db_link (db_id, deployment_id, org_id)
		VALUES ($1, $2, $3)
	`, dbId, deployment.Id, deployment.OrgId)
	return err
}

func getLinkedDeployment(resourceId int64, linkTable string, linkResource string, orgId int32, role *core.Role) (*core.FullDeployment, error) {
	if !role.Permissions.HasAccess(core.ResourceDeployments, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	deployment := core.FullDeployment{}
	rows, err := dbConn.Queryx(fmt.Sprintf(`
		SELECT d.*
		FROM deployments AS d
		INNER JOIN %s AS lnk
			ON lnk.deployment_id = d.id
				AND lnk.org_id = d.org_id
		WHERE lnk.%s = $1
			AND lnk.org_id = $2
	`, linkTable, linkResource), resourceId, orgId)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, nil
	}

	err = rows.StructScan(&deployment)
	return &deployment, err
}

func GetSystemDeployment(systemId int64, orgId int32, role *core.Role) (*core.FullDeployment, error) {
	return getLinkedDeployment(systemId, "deployment_system_link", "system_id", orgId, role)
}

func GetDatabaseDeployment(dbId int64, orgId int32, role *core.Role) (*core.FullDeployment, error) {
	return getLinkedDeployment(dbId, "deployment_db_link", "db_id", orgId, role)
}

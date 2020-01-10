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
		OrgId:          orgId,
		DeploymentType: core.KNoDeployment,
	}

	rows, err := tx.Queryx(`
		INSERT INTO deployments (org_id, deployment_type)
		VALUES ($1, $2)
		RETURNING id
	`, deployment.OrgId, core.KNoDeployment)
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

func getSelfDeploymentHelper(id int64, orgId int32, role *core.Role) (*core.SelfDeployment, error) {
	servers, err := AllServersForDeployment(id, orgId, role)
	if err != nil {
		return nil, err
	}

	self := core.SelfDeployment{
		Servers: servers,
	}
	return &self, nil
}

func getVendorDeploymentHelper(id int64, orgId int32, role *core.Role) (*core.VendorDeployment, error) {
	vendor := core.VendorDeployment{
		SocFiles: make([]*core.ControlDocumentationFile, 0),
	}

	rows, err := dbConn.Queryx(`
		SELECT vd.vendor_name, vd.vendor_product
		FROM vendor_deployments AS vd
		WHERE deployment_id = $1 AND org_id = $2
	`, id, orgId)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return &vendor, nil
	}

	err = rows.StructScan(&vendor)
	if err != nil {
		return nil, err
	}

	vendor.SocFiles, err = GetSocDocumentationForDeployment(id, orgId, role)
	if err != nil {
		return nil, err
	}

	return &vendor, nil
}

func getDeploymentHelper(condition string, lite bool, role *core.Role, args ...interface{}) (*core.FullDeployment, error) {
	if !role.Permissions.HasAccess(core.ResourceDeployments, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	deployment := core.FullDeployment{}
	rows, err := dbConn.Queryx(fmt.Sprintf(`
		SELECT d.*
		FROM deployments AS d
		%s
	`, condition), args...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, nil
	}

	err = rows.StructScan(&deployment)
	if err != nil {
		return nil, err
	}

	if lite {
		return &deployment, nil
	}

	deployment.SelfDeployment, err = getSelfDeploymentHelper(deployment.Id, deployment.OrgId, role)
	if err != nil {
		return nil, err
	}

	deployment.VendorDeployment, err = getVendorDeploymentHelper(deployment.Id, deployment.OrgId, role)
	if err != nil {
		return nil, err
	}

	return &deployment, err
}

func getLinkedDeployment(resourceId int64, linkTable string, linkResource string, orgId int32, lite bool, role *core.Role) (*core.FullDeployment, error) {
	return getDeploymentHelper(fmt.Sprintf(`
		INNER JOIN %s AS lnk
			ON lnk.deployment_id = d.id
				AND lnk.org_id = d.org_id
		WHERE lnk.%s = $1
			AND lnk.org_id = $2
	`, linkTable, linkResource), lite, role, resourceId, orgId)
}

func GetDeploymentFromId(id int64, orgId int32, role *core.Role) (*core.FullDeployment, error) {
	return getDeploymentHelper("WHERE d.id = $1 AND d.org_id = $2", false, role, id, orgId)
}

func GetSystemDeployment(systemId int64, orgId int32, role *core.Role) (*core.FullDeployment, error) {
	return getLinkedDeployment(systemId, "deployment_system_link", "system_id", orgId, false, role)
}

func GetDatabaseDeployment(dbId int64, orgId int32, role *core.Role) (*core.FullDeployment, error) {
	return getLinkedDeployment(dbId, "deployment_db_link", "db_id", orgId, false, role)
}

func GetSystemDeploymentId(systemId int64, orgId int32, role *core.Role) (int64, error) {
	deploy, err := getLinkedDeployment(systemId, "deployment_system_link", "system_id", orgId, true, role)
	if err != nil {
		return -1, err
	}
	return deploy.Id, nil
}

func GetDatabaseDeploymentId(dbId int64, orgId int32, role *core.Role) (int64, error) {
	deploy, err := getLinkedDeployment(dbId, "deployment_db_link", "db_id", orgId, true, role)
	if err != nil {
		return -1, err
	}
	return deploy.Id, nil
}

func updateSelfDeploymentWithTx(id int64, orgId int32, selfDeploy *core.StrippedSelfDeployment, tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		DELETE FROM deployment_server_link
		WHERE deployment_id = $1 AND org_id = $2
	`, id, orgId)

	if err != nil {
		return err
	}

	for _, handle := range selfDeploy.Servers {
		_, err = tx.Exec(`
			INSERT INTO deployment_server_link (server_id, deployment_id, org_id)
			VALUES ($1, $2, $3)
		`, handle.Id, id, orgId)

		if err != nil {
			return err
		}
	}
	return nil
}

func updateVendorDeploymentWithTx(id int64, orgId int32, vendorDeploy *core.StrippedVendorDeployment, tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		INSERT INTO vendor_deployments (deployment_id, org_id, vendor_name, vendor_product)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (deployment_id, org_id) DO UPDATE SET
			vendor_name = EXCLUDED.vendor_name,
			vendor_product = EXCLUDED.vendor_product	
	`, id, orgId, vendorDeploy.VendorName, vendorDeploy.VendorProduct)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM vendor_soc_reports
		WHERE deployment_id = $1 AND org_id = $2
	`, id, orgId)

	if err != nil {
		return err
	}

	for _, handle := range vendorDeploy.SocFiles {
		_, err := tx.Exec(`
			INSERT INTO vendor_soc_reports (deployment_id, org_id, soc_report_file_id, soc_report_cat_id)
			VALUES ($1, $2, $3, $4)
		`, id, orgId, handle.Id, handle.CategoryId)

		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateDeployment(deployment *core.StrippedFullDeployment, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDeployments, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	var err error
	tx := dbConn.MustBegin()

	_, err = tx.Exec(`
		UPDATE deployments
		SET deployment_type = $3
		WHERE id = $1 AND org_id = $2
	`, deployment.Id, deployment.OrgId, deployment.DeploymentType)
	if err != nil {
		tx.Rollback()
		return err
	}

	if deployment.SelfDeployment != nil {
		err = updateSelfDeploymentWithTx(deployment.Id, deployment.OrgId, deployment.SelfDeployment, tx)
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	if deployment.VendorDeployment != nil {
		err = updateVendorDeploymentWithTx(deployment.Id, deployment.OrgId, deployment.VendorDeployment, tx)
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func LinkDeploymentsWithServer(deploymentIds []int64, serverId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDeployments, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceServers, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()

	for _, deployId := range deploymentIds {
		_, err := tx.Exec(`
			INSERT INTO deployment_server_link (deployment_id, server_id, org_id)
			VALUES ($1, $2, $3)
		`, deployId, serverId, orgId)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func DeleteDeploymentServerLink(deploymentId int64, serverId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDeployments, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceServers, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()

	_, err := tx.Exec(`
		DELETE FROM deployment_server_link
		WHERE deployment_id = $1
			AND server_id = $2
			AND org_id = $3
	`, deploymentId, serverId, orgId)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

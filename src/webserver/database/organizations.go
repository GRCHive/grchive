package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func FindOrganizationFromGroupName(groupName string) (*core.Organization, error) {
	rows, err := dbConn.Queryx(`
		SELECT * FROM organizations WHERE org_group_name = $1
	`, groupName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var org *core.Organization = new(core.Organization)
	rows.Next()
	err = rows.StructScan(org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func FindOrganizationFromGroupId(groupId string) (*core.Organization, error) {
	rows, err := dbConn.Queryx(`
		SELECT * FROM organizations WHERE org_group_id = $1
	`, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var org *core.Organization = new(core.Organization)
	rows.Next()
	err = rows.StructScan(org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func FindOrganizationFromSamlIdP(samlIdp string) (*core.Organization, error) {
	rows, err := dbConn.Queryx(`
		SELECT * FROM organizations WHERE saml_iden = $1
	`, samlIdp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var org *core.Organization = new(core.Organization)
	rows.Next()
	err = rows.StructScan(org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func FindOrganizationFromProcessFlowId(flowId int64) (*core.Organization, error) {
	org := core.Organization{}
	err := dbConn.Get(&org, `
		SELECT org.*
		FROM process_flows AS pf
		INNER JOIN organizations AS org
			ON org.id = pf.org_id
		WHERE pf.id = $1
	`, flowId)
	return &org, err
}

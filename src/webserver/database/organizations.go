package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func FindOrganizationFromGroupId(groupId string) (*core.Organization, error) {
	rows, err := dbConn.Queryx(`
		SELECT org_group_id, org_group_name, org_name FROM organizations WHERE org_group_name = $1
	`, groupId)
	if err != nil {
		return nil, err
	}

	var org *core.Organization = new(core.Organization)
	rows.Next()
	err = rows.StructScan(org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

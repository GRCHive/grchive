package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func FindOrganizationFromGroupName(groupId string) (*core.Organization, error) {
	rows, err := dbConn.Queryx(`
		SELECT * FROM organizations WHERE org_group_name = $1
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

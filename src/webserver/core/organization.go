package core

type Organization struct {
	OktaGroupId   string `db:"org_group_id"`
	OktaGroupName string `db:"org_group_name"`
	Name          string `db:"org_name"`
}

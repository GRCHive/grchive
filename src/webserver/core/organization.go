package core

type Organization struct {
	Id            int32  `db:"id"`
	OktaGroupId   string `db:"org_group_id"`
	OktaGroupName string `db:"org_group_name"`
	Name          string `db:"org_name"`
}

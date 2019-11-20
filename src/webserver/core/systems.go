package core

type System struct {
	Id      int64  `db:"id"`
	OrgId   int32  `db:"org_id"`
	Name    string `db:"name"`
	Purpose string `db:"purpose"`
}

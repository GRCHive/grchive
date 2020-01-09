package core

type Server struct {
	Id              int64  `db:"id"`
	OrgId           int32  `db:"org_id"`
	Name            string `db:"name"`
	Description     string `db:"description"`
	OperatingSystem string `db:"operating_system"`
	Location        string `db:"location"`
	IpAddress       string `db:"ip_address"`
}

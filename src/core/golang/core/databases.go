package core

type DatabaseType struct {
	Id   int32  `db:"id"`
	Name string `db:"name"`
}

type Database struct {
	Id        int64  `db:"id"`
	Name      string `db:"name"`
	OrgId     int32  `db:"org_id"`
	TypeId    int32  `db:"type_id"`
	OtherType string `db:"other_type"`
	Version   string `db:"version"`
}

type DatabaseConnection struct {
	Id         int64  `db:"id"`
	DbId       int64  `db:"db_id"`
	OrgId      int32  `db:"org_id"`
	ConnString string `db:"connection_string"`
	Username   string `db:"username"`
	Password   string `db:"password" json:"-"`
	Salt       string `db:"salt" json:"-"`
}

package core

type User struct {
	Id         int64  `db:"id"`
	FirstName  string `db:"first_name"`
	LastName   string `db:"last_name"`
	Email      string `db:"email"`
	OktaUserId string `db:"okta_id"`
	OrgId      int32  `db:"org_id"`
}

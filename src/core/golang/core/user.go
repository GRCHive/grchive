package core

type User struct {
	Id        int64  `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

func (u User) FullName() string {
	return u.FirstName + " " + u.LastName
}

type UserWithRole struct {
	User
	RoleId int64 `db:"role_id"`
	OrgId  int32 `db:"org_id"`
}

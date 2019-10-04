package core

type User struct {
	FirstName string
	LastName  string
	FullName  string
	Email     string
	ParentOrg *Organization
}

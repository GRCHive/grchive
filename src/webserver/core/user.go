package core

type User struct {
	FirstName   string
	LastName    string
	FullName    string
	Email       string
	DisplayName string
	ParentOrg   *Organization
}

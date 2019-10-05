package core

type User struct {
	FirstName   string
	LastName    string
	Email       string
	DisplayName string
	OktaUserId  string
	ParentOrg   *Organization
}

package core

type User struct {
	// We don't store the full name of the user from Okta because
	// apparently after you do a user update and request another ID token,
	// the first/last name will get updated but not the full name. So
	// piece together their full name ourselves.
	FirstName   string
	LastName    string
	Email       string
	DisplayName string
	OktaUserId  string
	ParentOrg   *Organization
}

package database

//import (
//	"github.com/lib/pq"
//	"gitlab.com/b3h47pte/audit-stuff/core"
//)

// string: Returns the SAML identifier. Empty string if not found.
// error: If not nil, indicates that something went wrong in the query.
func FindSAMLIdPFromDomain(domain string) (string, error) {
	var err error

	_, err = dbConn.Query(`
		SELECT idpIdenOkta FROM saml_idp WHERE domain = ?
	`, domain)

	if err != nil {
		return "", err
	}

	return "", nil
}

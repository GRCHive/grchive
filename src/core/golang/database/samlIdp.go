package database

import (
	"gitlab.com/b3h47pte/audit-stuff/core"
)

// string: Returns the SAML identifier. Empty string if not found.
// error: If not nil, indicates that something went wrong in the query.
func FindSAMLIdPFromDomain(domain string) (string, error) {
	rows, err := dbConn.Queryx(`
		SELECT idpIdenOkta FROM saml_idp WHERE domain = $1
	`, domain)

	if err != nil {
		return "", err
	}
	defer rows.Close()

	var iden string
	if !rows.Next() {
		return "", nil
	}
	err = rows.Scan(&iden)

	// If err is not nil then that means that the query found no rows.
	if err != nil {
		core.Info(err.Error())
		return "", err
	}

	return iden, nil
}

package sap_api

import (
	"gitlab.com/grchive/grchive/core"
)

type SapRfcConnectionOptions struct {
	Client       string          `db:"client"`
	SysNr        string          `db:"sysnr"`
	Host         string          `db:"host"`
	RealHostname core.NullString `db:"real_hostname"`
	Username     string          `db:"username"`
	Password     string          `db:"password"`
}

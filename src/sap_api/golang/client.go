package sap_api

import (
	"github.com/sap/gorfc/gorfc"
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

type sapClient struct {
	conn *gorfc.Connection
}

func CreateSapClient(opt SapRfcConnectionOptions) (*sapClient, error) {
	hostname := opt.Host
	if opt.RealHostname.NullString.Valid {
		hostname = opt.RealHostname.NullString.String
	}

	params := gorfc.ConnectionParameters{
		"client": opt.Client,
		"user":   opt.Username,
		"passwd": opt.Password,
		"lang":   "EN",
		"ashost": hostname,
		"sysnr":  opt.SysNr,
	}

	client := sapClient{}

	var err error
	client.conn, err = gorfc.ConnectionFromParams(params)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

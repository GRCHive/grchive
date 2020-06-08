package main

import (
	"flag"
	"github.com/sap/gorfc/gorfc"
	"gitlab.com/grchive/grchive/core"
)

func main() {
	client := flag.String("client", "", "SAP Client.")
	sys := flag.String("sys", "", "SAP System Number.")
	host := flag.String("host", "", "SAP Hostname.")
	username := flag.String("username", "", "SAP Username.")
	password := flag.String("password", "", "SAP Password.")
	flag.Parse()

	system := gorfc.ConnectionParameters{
		"client": *client,
		"user":   *username,
		"passwd": *password,
		"lang":   "EN",
		"ashost": *host,
		"sysnr":  *sys,
	}

	core.Info("Create Connection")
	conn, err := gorfc.ConnectionFromParams(system)
	if err != nil {
		core.Error(err.Error())
	}
	defer conn.Close()

	fnDesc, err := conn.GetFunctionDescription("GRCHIVE_FM1")
	if err != nil {
		core.Error(err.Error())
	}
	core.Info(fnDesc.Name)
	for _, param := range fnDesc.Parameters {
		core.Info("\tParam: ", param.Name, param.ParameterType, param.Direction)
	}

	params := map[string]interface{}{}

	results, err := conn.Call("GRCHIVE_FM1", params)
	if err != nil {
		core.Error(err.Error())
	}

	for k, v := range results {
		core.Info("RESULT: ", k)
		core.Info(v)
	}

	core.Info("OK")
}

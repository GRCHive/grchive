package main

import (
	"flag"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
)

func main() {
	core.Info("STARTING INTEGRATION RUNNER...")

	core.Init()
	database.Init()
	webcore.InitializeWebcore()
	vault.Initialize(vault.VaultConfig{
		Url:      core.EnvConfig.Vault.Url,
		Username: core.EnvConfig.Vault.Username,
		Password: core.EnvConfig.Vault.Password,
	}, core.EnvConfig.Tls.Config())

	sapErpRfc := flag.Int64("sapErpRfc", -1, "SAP ERP RFC to process.")
	sapErpRfcVersion := flag.Int64("sapErpRfcVersion", -1, "SAP ERP RFC version to process.")
	flag.Parse()

	if *sapErpRfcVersion != -1 && *sapErpRfc != -1 {
		err := handleSapErpVersionWrapper(*sapErpRfcVersion, *sapErpRfc)
		if err != nil {
			core.Error("Failed to run SAP ERP: " + err.Error())
		}
	} else {
	}
}

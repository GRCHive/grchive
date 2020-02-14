package main

import (
	"flag"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
)

func main() {
	core.Init()
	database.Init()
	webcore.InitializeWebcore()
	vault.Initialize(vault.VaultConfig{
		Url:      core.EnvConfig.Vault.Url,
		Username: core.EnvConfig.Vault.Username,
		Password: core.EnvConfig.Vault.Password,
	})

	queryId := flag.Int64("queryId", -1, "Refresh ID to retrieve data for. Will not read from gRPC if specified.")
	orgId := flag.Int64("orgId", -1, "Org ID to retrieve data for. Will not read from gRPC if specified.")
	versionNum := flag.Int64("version", -1, "Version to retrieve data for. Will not read from gRPC if specified.")
	flag.Parse()

	if *queryId >= 0 && *orgId >= 0 && *versionNum >= 0 {
		result, err := runQuery(*queryId, int32(*orgId), int32(*versionNum))
		if err != nil {
			core.Error("Failed to run query: " + err.Error())
		}
		core.Info(result.Columns)
		core.Info(result.CsvText)
	} else {
	}
}

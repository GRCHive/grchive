package main

import (
	"flag"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gcloud_api"
	"gitlab.com/grchive/grchive/vault_api"
	//	"os"
)

func main() {
	core.Init()
	database.Init()
	vault.Initialize(vault.VaultConfig{
		Url:      core.EnvConfig.Vault.Url,
		Username: core.EnvConfig.Vault.Username,
		Password: core.EnvConfig.Vault.Password,
	}, core.EnvConfig.Tls.Config())
	gcloud.DefaultGCloudApi.InitFromJson(core.EnvConfig.Gcloud.AuthFilename)

	//scriptId := flag.Int64("scriptId", -1, "Script to retrieve.")
	//jarId := flag.Int64("jarId", -1, "JAR to retrieve.")
	scriptFname := flag.String("scriptFname", "", "Filename of script to compile.")
	jarFname := flag.String("jarFname", "", "Filename of JAR to run.")
	orgId := flag.Int("orgId", -1, "Org ID to retrieve data for.")
	roleId := flag.Int64("roleId", -1, "Role to run script as.")
	immediate := flag.Bool("immediate", false, "Whether to run immediately.")
	local := flag.Bool("local", false, "Whether to filepath inputs instead of retrieving from DB.")
	flag.Parse()

	if *immediate {
		var err error
		if *local {
			dirName, err := createHostWorkspaceDirectory(*scriptFname, *jarFname)
			//defer os.RemoveAll(dirName)
			if err != nil {
				core.Error("Failed to setup workspace directory: " + err.Error())
			}

			out, err := handleLocal(
				dirName,
				int32(*orgId),
				*roleId,
				core.ScriptRunSettings{
					CpuAllocation:          1.0,
					MemBytesAllocation:     core.Gigabytes(1.0),
					DiskSizeBytes:          core.Gigabytes(1.0),
					KotlinContainerVersion: "latest",
				},
			)

			if err != nil {
				core.Error("Failed to run container: " + err.Error())
			}

			core.Info("JAR Path: ", out.CompiledJarFname)
			core.Info("Logs: ", out.Logs)
		} else {
			core.Error("Currently unsupported.")
		}

		if err != nil {
			core.Error("Failed to compile/run: " + err.Error())
		}
	} else {
		core.Error("Currently unsupported.")
	}
}

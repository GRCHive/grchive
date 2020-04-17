package main

import (
	"flag"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gitea_api"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
	"os"
	"strconv"
)

// A multi-part function that's responsible for
// 	1) Compiling the client's Kotlin code into a JAR
// 	2) Deploying that JAR into Artifactory
// 	3) Tracking the logs of steps #1 and #2 and storing that in the database associated with the commit
// 	4) Tracking the JAR path from step #2 and storing that in the database associated with the commit
// 	5) Tracking whether both steps #1 and #2 were successful and storing that in the in database associated with the commit
func run(dir string) {
	// DRONE_COMMIT: https://docs.drone.io/pipeline/environment/reference/drone-commit/
	commitSha := os.Getenv("DRONE_COMMIT")

	tracker := Tracker{
		workDir: dir,
		commit:  commitSha,
		version: fmt.Sprintf("0.0-%s", commitSha),
	}

	// SCRIPT_RUN: A custom variable that can be set if we need to compile specifically for running a script.
	// In that case, behavior is modified to create a snapshot specifically for this run and this information
	// will be passed directly to the runner instead of being stored back in the database.
	scriptRunId, runSet := os.LookupEnv("SCRIPT_RUN")
	if runSet {
		id, err := strconv.ParseInt(scriptRunId, 10, 64)
		if err != nil {
			core.Error("Failed to get script run: " + err.Error())
		}
		tracker.scriptRunId = core.CreateNullInt64(id)

		// Run ID should be unique.
		tracker.version = fmt.Sprintf("0.0-%d-SNAPSHOT", id)
	}

	tracker.Start()

	compileAndDeploy(&tracker)

	tracker.End()
}

func main() {
	core.Init()
	database.Init()
	webcore.InitializeWebcore()

	vault.Initialize(vault.VaultConfig{
		Url:      core.EnvConfig.Vault.Url,
		Username: core.EnvConfig.Vault.Username,
		Password: core.EnvConfig.Vault.Password,
	}, core.EnvConfig.Tls.Config())

	giteaToken, err := vault.GetSecretWithKey(core.EnvConfig.Gitea.Token, "token")
	if err != nil {
		core.Error("Failed to get Gitea token: " + err.Error())
	}

	gitea.GlobalGiteaApi.MustInitialize(gitea.GiteaConfig{
		Protocol: core.EnvConfig.Gitea.Protocol,
		Host:     core.EnvConfig.Gitea.Host,
		Port:     core.EnvConfig.Gitea.Port,
		Token:    giteaToken,
	})

	core.Debug("RabbitMQ Init")
	webcore.DefaultRabbitMQ.Connect(
		*core.EnvConfig.RabbitMQ,
		webcore.MQClientConfig{},
		core.EnvConfig.Tls)
	defer webcore.DefaultRabbitMQ.Cleanup()

	repoDir := flag.String("dir", "", "Directory containing source code to build.")
	flag.Parse()
	run(*repoDir)

	// Ensure all rabbitmq messages got sent to ensure that script run messages get sent before exiting.
	webcore.DefaultRabbitMQ.WaitForAllMessagesToBeSent()
}

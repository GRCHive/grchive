package main

import (
	"flag"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/vault_api"
	"os"
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

	tracker.Start()

	compileAndDeploy(&tracker)

	tracker.End()
}

func main() {
	core.Init()
	database.Init()
	vault.Initialize(vault.VaultConfig{
		Url:      core.EnvConfig.Vault.Url,
		Username: core.EnvConfig.Vault.Username,
		Password: core.EnvConfig.Vault.Password,
	}, core.EnvConfig.Tls.Config())

	repoDir := flag.String("dir", "", "Directory containing source code to build.")
	flag.Parse()
	run(*repoDir)
}

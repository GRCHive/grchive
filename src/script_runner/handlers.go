package main

import (
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/script_runner/worker"
	"os"
	"strings"
)

func handleRunTracker(tracker *Tracker, runId int64, jar string) error {
	mavenDep := strings.Split(jar, ":")
	// Get Client Script information.
	run, err := database.GetScriptRun(tracker.runId, core.ServerRole)
	if err != nil {
		return err
	}

	script, err := database.GetScriptFromScriptCodeLink(run.LinkId, core.ServerRole)
	if err != nil {
		return err
	}

	// Determine expected package and function name.
	// 	- Package: com.grchive.web.client
	// 	- Class: ${Package}.${Upper Camel-case script name}_${Script Id}Kt
	// 	- Function: Lower Camel-case script name.
	className := fmt.Sprintf(
		"com.grchive.web.client.%s_%dKt",
		strcase.ToCamel(script.Name),
		script.Id,
	)
	functionName := strcase.ToLowerCamel(script.Name)
	metadataName := strings.TrimPrefix(script.MetadataFilename(), "src/main/resources")

	options := worker.WorkerOptions{
		ClientLibGroupId:    mavenDep[0],
		ClientLibArtifactId: mavenDep[1],
		ClientLibVersion:    mavenDep[2],
		RunClassName:        className,
		RunFunctionName:     functionName,
		RunMetadataName:     metadataName,
		RunId:               runId,
	}

	worker, err := tracker.factory.CreateWorker(options)
	if err != nil {
		return err
	}
	defer worker.Cleanup()

	retCode, err := worker.Run()
	if err != nil {
		return err
	}

	logs, err := worker.Logs()
	if err != nil {
		return err
	}

	tracker.Log(logs, tracker.stdout)
	if retCode == 0 {
		tracker.MarkSuccess()
	} else {
		tracker.MarkError(errors.New(fmt.Sprintf("Process exited with run code: %d", retCode)))
	}

	return nil
}

func handleRun(runId int64, jar string, mavenDir string, stdout bool) error {
	tracker := Tracker{
		runId:        runId,
		mavenRootDir: mavenDir,
		stdout:       stdout,
	}

	if _, ok := os.LookupEnv("WORKER_K8S"); ok {
		tracker.factory = worker.KubeFactory{}
	} else {
		worker.InitDocker()
		tracker.factory = worker.DockerFactory{}
	}

	err := tracker.Start()

	if err == nil {
		err = handleRunTracker(&tracker, runId, jar)
		if err != nil {
			tracker.MarkError(err)
		}
	} else {
		tracker.MarkError(err)
	}
	return tracker.End()
}

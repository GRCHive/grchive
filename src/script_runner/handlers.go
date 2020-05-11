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

	//// Create folder with the project and mount this as a folder in the worker image.
	//// The worker image will compile and run using the code inside this folder via Maven.
	//workDir, err := ioutil.TempDir("", "script-runner")
	//if err != nil {
	//	return err
	//}
	//tracker.Log("WORK DIR: "+workDir, true)
	//defer os.RemoveAll(workDir)

	//// Copy over the template while replacing all the .tmpl files with an automatically generate
	//// file using certain predetermine variables. This code could probably be shared with the webserver
	//// which does something similar for Gitea repository template generation?
	//templateParams := map[string]string{
	//}

	//templateDirItems, err := ioutil.ReadDir(templateDir)
	//if err != nil {
	//	return err
	//}

	//for _, f := range templateDirItems {
	//	err = handleDirectoryFileTemplateGen(f, templateDir, workDir, templateParams)
	//	if err != nil {
	//		return err
	//	}
	//}

	//// Kick off Docker container to handle the job.
	//containerName := fmt.Sprintf("script-runner-%d", runId)
	//err = createKotlinContainer(workDir, containerName, tracker.mavenRootDir, className, functionName, metadataName, strconv.FormatInt(runId, 10))
	//if err != nil {
	//	return err
	//}

	//// Examine Docker exit code/logs to see whether or not the job succeeded.
	//retCode, err := runKotlinContainer(containerName)
	//if err != nil {
	//	removeKotlinContainer(containerName)
	//	return err
	//}

	//logs, err := readLogsFromContainer(containerName)
	//if err != nil {
	//	removeKotlinContainer(containerName)
	//	return err
	//}

	//tracker.Log(logs, tracker.stdout)
	//if retCode == 0 {
	//	tracker.MarkSuccess()
	//} else {
	//	tracker.MarkError(errors.New(fmt.Sprintf("Process exited with run code: %d", retCode)))
	//}

	//return removeKotlinContainer(containerName)
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

package main

import (
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

const templateDir string = "src/script_runner/template"

func handleDirectoryFileTemplateGen(info os.FileInfo, currentTemplateDir string, workDir string, params map[string]string) error {
	fullTemplatePath := filepath.Join(currentTemplateDir, info.Name())
	fullWorkDirPath := filepath.Join(workDir, info.Name())

	if info.IsDir() {
		templateDirItems, err := ioutil.ReadDir(fullTemplatePath)
		if err != nil {
			return err
		}

		for _, f := range templateDirItems {
			err = handleDirectoryFileTemplateGen(f, fullTemplatePath, fullWorkDirPath, params)
			if err != nil {
				return err
			}
		}
	} else {
		templateData, err := ioutil.ReadFile(fullTemplatePath)
		if err != nil {
			return err
		}

		finalData := string(templateData)
		if strings.HasSuffix(fullTemplatePath, ".tmpl") {
			tmpl, err := template.New(info.Name()).Parse(string(templateData))
			if err != nil {
				return err
			}

			finalData, err = core.TextTemplateToString(tmpl, params)
			fullTemplatePath = strings.TrimSuffix(fullTemplatePath, ".tmpl")
			fullWorkDirPath = strings.TrimSuffix(fullWorkDirPath, ".tmpl")
		}

		err = os.MkdirAll(filepath.Dir(fullWorkDirPath), os.FileMode(0755))
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(fullWorkDirPath, []byte(finalData), os.FileMode(0755))
		if err != nil {
			return err
		}
	}
	return nil
}

func handleRunTracker(tracker *Tracker, runId int64, jar string) error {
	// Create folder with the project and mount this as a folder in the worker image.
	// The worker image will compile and run using the code inside this folder via Maven.
	workDir, err := ioutil.TempDir("", "script-runner")
	if err != nil {
		return err
	}
	tracker.Log("WORK DIR: "+workDir, true)
	defer os.RemoveAll(workDir)

	mavenDep := strings.Split(jar, ":")
	// Copy over the template while replacing all the .tmpl files with an automatically generate
	// file using certain predetermine variables. This code could probably be shared with the webserver
	// which does something similar for Gitea repository template generation?
	templateParams := map[string]string{
		"ARTIFACTORY_HOST":   core.EnvConfig.Artifactory.Host,
		"ARTIFACTORY_PORT":   strconv.FormatInt(int64(core.EnvConfig.Artifactory.Port), 10),
		"CLIENT_GROUP_ID":    mavenDep[0],
		"CLIENT_ARTIFACT_ID": mavenDep[1],
		"CLIENT_VERSION":     mavenDep[2],
	}

	templateDirItems, err := ioutil.ReadDir(templateDir)
	if err != nil {
		return err
	}

	for _, f := range templateDirItems {
		err = handleDirectoryFileTemplateGen(f, templateDir, workDir, templateParams)
		if err != nil {
			return err
		}
	}

	// Get Client Script information.
	run, err := database.GetScriptRun(tracker.runId)
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

	// Kick off Docker container to handle the job.
	containerName := fmt.Sprintf("script-runner-%d", runId)
	err = createKotlinContainer(workDir, containerName, className, functionName, metadataName, strconv.FormatInt(runId, 10))
	if err != nil {
		return err
	}

	// Examine Docker exit code/logs to see whether or not the job succeeded.
	retCode, err := runKotlinContainer(containerName)
	if err != nil {
		removeKotlinContainer(containerName)
		return err
	}

	logs, err := readLogsFromContainer(containerName)
	if err != nil {
		removeKotlinContainer(containerName)
		return err
	}

	tracker.Log(logs, false)
	if retCode == 0 {
		tracker.MarkSuccess()
	} else {
		tracker.MarkError(errors.New(fmt.Sprintf("Process exited with run code: %d", retCode)))
	}

	return removeKotlinContainer(containerName)
}

func handleRun(runId int64, jar string) error {
	tracker := Tracker{
		runId: runId,
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

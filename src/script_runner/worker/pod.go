package main

import (
	"flag"
	"gitlab.com/grchive/grchive/core"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

const templateDir string = "src/script_runner/template"
const workDir string = "/data"

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

func main() {
	core.Init()

	group := flag.String("group", "", "Client Maven group.")
	artifact := flag.String("artifact", "", "Client Maven artifact.")
	version := flag.String("vers", "", "Client Maven version.")
	className := flag.String("class", "", "Class to run.")
	fnName := flag.String("fn", "", "Function to run.")
	metadataName := flag.String("meta", "", "Metadata file.")
	runId := flag.Int64("runId", -1, "Run ID in the database.")
	flag.Parse()

	templateDirItems, err := ioutil.ReadDir(templateDir)
	if err != nil {
		core.Error("Failed to read template items: " + err.Error())
	}

	templateParams := map[string]string{
		"ARTIFACTORY_HOST":              core.EnvConfig.Artifactory.Host,
		"ARTIFACTORY_PORT":              strconv.FormatInt(int64(core.EnvConfig.Artifactory.Port), 10),
		"KOTLIN_CORE_LIB_MAJOR_VERSION": strconv.FormatInt(int64(core.EnvConfig.Kotlin.MajorVersion), 10),
		"KOTLIN_CORE_LIB_MINOR_VERSION": strconv.FormatInt(int64(core.EnvConfig.Kotlin.MinorVersion), 10),
		"KOTLIN_CORE_LIB_GROUP_ID":      core.EnvConfig.Kotlin.GroupId,
		"KOTLIN_CORE_LIB_ARTIFACT_ID":   core.EnvConfig.Kotlin.ArtifactId,
		"CLIENT_GROUP_ID":               *group,
		"CLIENT_ARTIFACT_ID":            *artifact,
		"CLIENT_VERSION":                *version,
	}

	for _, f := range templateDirItems {
		err = handleDirectoryFileTemplateGen(f, templateDir, workDir, templateParams)
		if err != nil {
			core.Error("Failed to generate template item: " + err.Error())
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		core.Error("Failed to get cwd: " + err.Error())
	}

	// This is mainly to support development environments (I think) where
	// the directory where we are running in just contains a symbolic link to the
	// real file. So in reality we want to mount the directory of the real file
	// because the symbolic link won't work in the Docker container.
	configPath, err := core.FindAbsolutePathThroughSymbolicLink(filepath.Join(wd, "src", "webserver", "config", "config.toml"))
	if err != nil {
		core.Error("Failed to get config path: " + err.Error())
	}

	cmd := exec.Command(
		"/data/run.sh",
		*className,
		*fnName,
		*metadataName,
		strconv.FormatInt(*runId, 10),
		configPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "/data"

	err = cmd.Run()
	if err != nil {
		core.Error("Failed to run: " + err.Error())
	}
}

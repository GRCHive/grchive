package main

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"io/ioutil"
	"os"
	"path/filepath"
)

type KotlinOutput struct {
	CompiledJarFname string
	Logs             string
}

func getExpectedScriptFname(dir string) string {
	return filepath.Join(dir, "script.kt")
}

func getExpectedJarFname(dir string) string {
	return filepath.Join(dir, "script.jar")
}

func getExpectedCompiledJarFname(dir string) string {
	return filepath.Join(dir, "script-compile.jar")
}

func createHostWorkspaceDirectory(scriptFname, jarFname string) (string, error) {
	// Create temporary directory to store the inputs to the Docker container.
	// At the end we'd expect this folder to contain 3 files:
	// 	1) script.kt 				, The copied Kotlin script provided by the user.
	// 	2) script.jar 				, The copied JAR file provided by the user (if any).
	// 	3) script-compile.jar 		, The copiled JAR file created by us.
	tempDir, err := ioutil.TempDir(os.TempDir(), "script-runner-grchive")
	if err != nil {
		return "", err
	}

	// Otherwise the user in the docker container can't read it.
	err = os.Chmod(tempDir, os.FileMode(0755))
	if err != nil {
		return tempDir, err
	}

	tmpScriptFname := getExpectedScriptFname(tempDir)
	tmpJarFname := getExpectedJarFname(tempDir)

	// Script filename must exist.
	err = core.CopyFile(scriptFname, tmpScriptFname)
	if err != nil {
		return tempDir, err
	}

	// JAR file may or may not be used. In the case where it's not used
	// we will use the JAR that gets compiled.
	if jarFname != "" {
		err = core.CopyFile(jarFname, tmpJarFname)
		if err != nil {
			return tempDir, err
		}
	}
	return tempDir, err
}

func handleWeb(scriptId int64, orgId int32, roleId int64, jarId int64) error {
	return nil
}

func handleLocal(dirName string, orgId int32, roleId int64, runSettings core.ScriptRunSettings) (*KotlinOutput, error) {
	containerName := fmt.Sprintf("grchive-%s", core.DefaultUuidGen.GenStr())
	workspaceVolumeName := fmt.Sprintf("data-%s", containerName)

	core.Info("CONTAINER: ", containerName)
	core.Info("VOLUME: ", workspaceVolumeName)
	core.Info("WORKSPACE: ", dirName)

	err := createKotlinContainer(dirName, containerName, workspaceVolumeName, runSettings)
	if err != nil {
		return nil, err
	}

	err = runKotlinContainer(containerName, runSettings)
	if err != nil {
		return nil, err
	}

	logs, err := readLogsFromContainer(containerName)
	if err != nil {
		return nil, err
	}

	outputJarFname := getExpectedCompiledJarFname(dirName)
	err = copyDataFromContainer(containerName, getExpectedCompiledJarFname(dockerOutputDir), outputJarFname)
	if err != nil {
		return nil, err
	}

	err = removeKotlinContainer(containerName, workspaceVolumeName)
	if err != nil {
		return nil, err
	}

	return &KotlinOutput{
		CompiledJarFname: outputJarFname,
		Logs:             logs,
	}, nil
}

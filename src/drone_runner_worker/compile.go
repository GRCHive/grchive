package main

import (
	"bytes"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gitea_api"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func trackedRunCmd(tracker *Tracker, args string) error {
	tracker.Log(fmt.Sprintf("!!! Running command %s", args))
	cargs := strings.Split(args, " ")

	cmd := exec.Command(cargs[0], cargs[1:]...)
	cmd.Dir = tracker.workDir

	out, err := cmd.CombinedOutput()
	tracker.Log(string(out) + "\n")

	if err != nil {
		return err
	}

	return nil
}

func computeJarPathFromMvn(tracker *Tracker) (string, error) {
	// For simplicity, this calls our pre-built script in every repository.
	cmd := exec.Command("bash", filepath.Join(tracker.workDir, "get_jar_from_mvn.sh"))
	cmd.Dir = tracker.workDir

	stdout := bytes.Buffer{}
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil
}

func checkoutScriptToRevision(tracker *Tracker) error {
	run, err := database.GetScriptRun(tracker.scriptRunId.NullInt64.Int64, core.ServerRole)
	if err != nil {
		return err
	}

	script, err := database.GetScriptFromScriptCodeLink(run.LinkId, core.ServerRole)
	if err != nil {
		return err
	}

	code, err := database.GetCodeFromScriptCodeLink(run.LinkId, core.ServerRole)
	if err != nil {
		return err
	}

	repo, err := database.GetLinkedGiteaRepository(script.OrgId)
	if err != nil {
		return err
	}

	gitRepo := gitea.GiteaRepository{
		Owner: repo.GiteaOrg,
		Name:  repo.GiteaRepo,
	}

	// Need to get script and script metadata at the specified version.
	{
		scriptFname := script.Filename("kt")
		tracker.Log(fmt.Sprintf("!!! Checking out %s to %s", scriptFname, code.GitHash))
		scriptData, _, err := gitea.GlobalGiteaApi.RepositoryGetFile(gitRepo, scriptFname, code.GitHash)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filepath.Join(tracker.workDir, scriptFname), []byte(scriptData), os.FileMode(0755))
		if err != nil {
			return err
		}
	}

	{
		metadataFname := script.MetadataFilename()
		tracker.Log(fmt.Sprintf("!!! Checking out %s to %s", metadataFname, code.GitHash))
		metadataData, _, err := gitea.GlobalGiteaApi.RepositoryGetFile(gitRepo, metadataFname, code.GitHash)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filepath.Join(tracker.workDir, metadataFname), []byte(metadataData), os.FileMode(0755))
		if err != nil {
			return err
		}
	}

	return nil
}

func compileAndDeploy(tracker *Tracker) {
	// 0. Determine if we're doing a script run compilation. If so, check out the script to the proper version.
	if tracker.IsScriptRunCompile() {
		err := checkoutScriptToRevision(tracker)
		if err != nil {
			tracker.MarkFailure(err)
			return
		}
	}

	// 1. Set version of the pom.xml properly to what's stored in the tracker.
	err := trackedRunCmd(tracker, fmt.Sprintf("mvn versions:set -DnewVersion=%s", tracker.version))
	if err != nil {
		tracker.MarkFailure(err)
		return
	}

	// 2. Compile & Deploy using Maven.
	err = trackedRunCmd(tracker, "mvn deploy")
	if err != nil {
		tracker.MarkFailure(err)
		return
	}

	// 3. Get fullly qualified JAR name from MVN.
	tracker.Log("!!! Obtaining JAR name")
	jarPath, err := computeJarPathFromMvn(tracker)
	if err != nil {
		tracker.MarkFailure(err)
		return
	}

	tracker.MarkSuccess(jarPath)
}

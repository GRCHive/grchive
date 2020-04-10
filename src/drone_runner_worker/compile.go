package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"strings"
)

func trackedRunCmd(tracker *Tracker, args string) error {
	tracker.Log(fmt.Sprintf("!!! Running command %s", args))
	cargs := strings.Split(args, " ")

	cmd := exec.Command(cargs[0], cargs[1:]...)
	cmd.Dir = tracker.workDir

	out, err := cmd.CombinedOutput()
	tracker.logs.WriteString(string(out) + "\n")

	if err != nil {
		return err
	}

	return nil
}

func computeJarPathFromMvn(tracker *Tracker) (string, error) {
	// For simplicity, this calls our pre-built script in every repository.
	cmd := exec.Command("bash", path.Join(tracker.workDir, "get_jar_from_mvn.sh"))
	cmd.Dir = tracker.workDir

	stdout := bytes.Buffer{}
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil
}

func compileAndDeploy(tracker *Tracker) {
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

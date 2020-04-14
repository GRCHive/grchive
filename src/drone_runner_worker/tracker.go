package main

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/vault_api"
	"strings"
	"time"
)

// Handles tracking and storing the result of the Compile/Deploy stages.
type Tracker struct {
	workDir string
	commit  string
	version string

	// Stats
	startTime time.Time

	// Results
	success bool
	logs    strings.Builder
	jar     string
}

func (t *Tracker) Start() {
	t.Log("START CI JOB")
	err := database.StartDroneCIJob(t.commit)
	if err != nil {
		core.Error("Failed to record start in DB:" + err.Error())
	}

	t.startTime = time.Now().UTC()
}

func (t *Tracker) Log(msg string) {
	core.Info(msg)
	t.logs.WriteString(msg + "\n")
}

func (t *Tracker) MarkFailure(err error) {
	t.logs.WriteString("^^^^ FAILURE: " + err.Error() + "\n")
	t.success = false
}

func (t *Tracker) MarkSuccess(jarPath string) {
	t.jar = jarPath
	t.success = true
}

func (t *Tracker) End() {
	encryptedLogs, err := vault.TransitEncrypt(core.EnvConfig.LogEncryptionPath, []byte(t.logs.String()))
	if err != nil {
		core.Error("Failed to encrypt logs:" + err.Error())
	}

	err = database.FinishDroneCIJob(t.commit, t.success, string(encryptedLogs), t.jar)
	if err != nil {
		core.Error("Failed to record start in DB:" + err.Error())
	}

	now := time.Now().UTC()
	elapsed := now.Sub(t.startTime)
	durationSeconds := float64(elapsed.Milliseconds()) / 1000.0
	t.Log(fmt.Sprintf("END CI JOB [Elapsed %f seconds]", durationSeconds))
}

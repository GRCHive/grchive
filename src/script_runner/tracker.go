package main

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/vault_api"
	"strings"
	"time"
)

type Tracker struct {
	runId        int64
	mavenRootDir string
	stdout       bool

	logs    strings.Builder
	success bool

	// Stats
	startTime time.Time
}

func (t *Tracker) Start() error {
	t.startTime = time.Now().UTC()
	err := database.StartExecuteScriptRun(t.runId)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tracker) MarkError(err error) {
	t.Log("Script Failed: "+err.Error(), true)
	t.success = false
}

func (t *Tracker) MarkSuccess() {
	t.Log("Script Success!", true)
	t.success = true
}

func (t *Tracker) Log(msg string, stdout bool) {
	if stdout {
		core.Info(msg)
	}

	t.logs.WriteString(msg + "\n")
}

func (t *Tracker) End() error {
	now := time.Now().UTC()
	elapsed := now.Sub(t.startTime)
	durationSeconds := float64(elapsed.Milliseconds()) / 1000.0
	t.Log(fmt.Sprintf("END RUN JOB [Elapsed %f seconds]", durationSeconds), true)

	encryptedLogs, err := vault.TransitEncrypt(core.EnvConfig.LogEncryptionPath, []byte(t.logs.String()))
	if err != nil {
		return err
	}

	return database.FinishExecuteScriptRun(t.runId, t.success, string(encryptedLogs))
}

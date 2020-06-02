package main

import (
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/vault_api"
	"strings"
	"time"
)

type PerServerTracker struct {
	stdout     bool
	log        strings.Builder
	run        *core.ShellScriptRunPerServer
	success    bool
	startTime  time.Time
	finishTime time.Time
}

func (t *PerServerTracker) MarkStart() error {
	t.startTime = time.Now().UTC()

	tx := database.CreateTx()
	return database.WrapTx(tx, func() error {
		return database.MarkShellScriptRunForServerStartWithTx(tx, t.run.RunId, t.run.ServerId, t.startTime)
	})
}

func (t *PerServerTracker) MarkSuccessFailure(logs string, success bool) error {
	if t.stdout {
		core.Info(logs)
	}
	t.log.WriteString(logs)
	t.success = success
	t.finishTime = time.Now().UTC()
	return nil
}

func (t *PerServerTracker) Finish() {
	var logStr string
	encryptedLogs, err := vault.TransitEncrypt(core.EnvConfig.LogEncryptionPath, []byte(t.log.String()))
	if err != nil {
		logStr = "Failed to encrypt logs : " + err.Error()
		core.Warning(logStr)
	} else {
		logStr = string(encryptedLogs)
	}

	tx := database.CreateTx()
	err = database.WrapTx(tx, func() error {
		return database.MarkShellScriptRunForServerFinishWithTx(
			tx,
			t.run.RunId,
			t.run.ServerId,
			t.finishTime,
			t.success,
			logStr,
		)
	})

	if err != nil {
		core.Warning("Failed to commit script run for server: " + err.Error())
	}
}

type Tracker struct {
	runId  int64
	stdout bool

	startTime time.Time
	perServer map[int64]*PerServerTracker
}

func (t *Tracker) MarkStart() error {
	t.startTime = time.Now().UTC()

	tx := database.CreateTx()
	return database.WrapTx(tx, func() error {
		return database.MarkShellScriptRunStartWithTx(tx, t.runId, t.startTime)
	})
}

func CreateTracker(runId int64, stdout bool) Tracker {
	return Tracker{
		runId:     runId,
		stdout:    stdout,
		perServer: map[int64]*PerServerTracker{},
	}
}

func (t *Tracker) StartTrackingServer(serverId int64, runPerServer *core.ShellScriptRunPerServer) {
	t.perServer[serverId] = &PerServerTracker{
		log:     strings.Builder{},
		run:     runPerServer,
		success: false,
		stdout:  t.stdout,
	}
}

func (t *Tracker) GetTrackerForServer(serverId int64) *PerServerTracker {
	return t.perServer[serverId]
}

func (t *Tracker) Finish() {
	tx := database.CreateTx()
	err := database.WrapTx(tx, func() error {
		return database.MarkShellScriptRunEndWithTx(tx, t.runId, time.Now().UTC())
	})

	if err != nil {
		core.Warning("Failed to commit tracker data: " + err.Error())
	}
}

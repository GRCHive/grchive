package main

import (
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gcloud_api"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
)

type PerServerRunData struct {
	ConnectionChoices core.ServerConnectionOptions
}

type RunData struct {
	servers   []*core.Server
	perServer map[int64]*PerServerRunData

	shellRun     *core.ShellScriptRun
	runPerServer map[int64]*core.ShellScriptRunPerServer

	shell        *core.ShellScript
	shellVersion *core.ShellScriptVersion
	scriptText   string
}

func loadPerServerRunData(serverId int64, orgId int32, t *Tracker) (*PerServerRunData, error) {
	var err error
	data := PerServerRunData{}

	data.ConnectionChoices.SshPassword, err = database.GetSSHPasswordConnectionForServer(serverId, orgId)
	if err != nil {
		return nil, err
	}

	if data.ConnectionChoices.SshPassword != nil {
		data.ConnectionChoices.SshPassword.Password, err = webcore.DecryptEncryptedPassword(data.ConnectionChoices.SshPassword.Password)
		if err != nil {
			return nil, err
		}
	}

	data.ConnectionChoices.SshKey, err = database.GetSSHKeyConnectionForServer(serverId, orgId)
	if err != nil {
		return nil, err
	}

	if data.ConnectionChoices.SshKey != nil {
		data.ConnectionChoices.SshKey.PrivateKey, err = webcore.DecryptEncryptedPassword(data.ConnectionChoices.SshKey.PrivateKey)
		if err != nil {
			return nil, err
		}
	}
	return &data, nil
}

const ShellEncryptionPath = "shell"

func loadRunData(t *Tracker) (*RunData, error) {
	var err error
	runData := RunData{}

	runData.shellRun, err = database.GetShellRun(t.runId)
	if err != nil {
		return nil, err
	}

	runData.servers, err = database.GetShellRunServers(t.runId)
	if err != nil {
		return nil, err
	}

	runData.perServer = map[int64]*PerServerRunData{}
	runData.runPerServer = map[int64]*core.ShellScriptRunPerServer{}
	for _, s := range runData.servers {
		runData.runPerServer[s.Id], err = database.GetShellRunForServer(t.runId, s.Id)
		if err != nil {
			return nil, err
		}

		runData.perServer[s.Id], err = loadPerServerRunData(s.Id, runData.runPerServer[s.Id].OrgId, t)
		if err != nil {
			return nil, err
		}

		t.StartTrackingServer(s.Id, runData.runPerServer[s.Id])
	}

	runData.shellVersion, err = database.GetShellScriptVersionFromId(runData.shellRun.ScriptVersionId)
	if err != nil {
		return nil, err
	}

	runData.shell, err = database.GetShellScriptFromId(runData.shellVersion.ShellId)
	if err != nil {
		return nil, err
	}

	storage := gcloud.DefaultGCloudApi.GetStorageApi()
	rawData, err := storage.DownloadVersioned(
		runData.shell.BucketId,
		runData.shell.StorageId,
		runData.shellVersion.GcsGeneration,
		core.EnvConfig.HmacKey,
	)

	if err != nil {
		return nil, err
	}

	decryptedScript, err := vault.TransitDecrypt(ShellEncryptionPath, rawData)
	if err != nil {
		return nil, err
	}

	runData.scriptText = string(decryptedScript)

	return &runData, nil
}

func run(t *Tracker, data *RunData) error {
	for _, s := range data.servers {
		err := runScriptOnServer(
			data,
			s,
			t.GetTrackerForServer(s.Id),
		)

		if err != nil {
			return err
		}
	}
	return nil
}

func handleRun(runId int64, stdout bool) error {
	tracker := CreateTracker(runId, stdout)
	defer tracker.Finish()

	err := tracker.MarkStart()
	if err != nil {
		core.Warning("Failed to mark run start: " + err.Error())
		return err
	}

	runData, err := loadRunData(&tracker)
	if err != nil {
		core.Warning("Failed to load run data: " + err.Error())
		return err
	}

	err = run(&tracker, runData)
	if err != nil {
		core.Warning("Failed to run: " + err.Error())
		return err
	}

	return nil
}

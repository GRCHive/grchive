package main

import (
	"encoding/json"
	"flag"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gcloud_api"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
)

func processShellRunnerMessages(data []byte) *webcore.RabbitMQError {
	msg := webcore.ShellRunnerMessage{}
	core.Info("RUN SHELL: " + string(data))
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return &webcore.RabbitMQError{err, false}
	}

	err = handleRun(msg.RunId, false)
	if err != nil {
		return &webcore.RabbitMQError{
			err,
			false,
		}
	}

	return nil
}

func main() {
	core.Info("STARTING SHELL RUNNER...")

	core.Init()
	database.Init()
	webcore.InitializeWebcore()
	vault.Initialize(vault.VaultConfig{
		Url:      core.EnvConfig.Vault.Url,
		Username: core.EnvConfig.Vault.Username,
		Password: core.EnvConfig.Vault.Password,
	}, core.EnvConfig.Tls.Config())
	gcloud.DefaultGCloudApi.InitFromJson(core.EnvConfig.Gcloud.AuthFilename)

	runId := flag.Int64("runId", -1, "Run ID in the database.")
	stdout := flag.Bool("stdout", false, "Print to stdout (local only).")
	flag.Parse()

	if *runId != -1 {
		err := handleRun(*runId, *stdout)
		if err != nil {
			core.Error("Failed to run: " + err.Error())
		}
	} else {
		webcore.DefaultRabbitMQ.Connect(*core.EnvConfig.RabbitMQ, webcore.MQClientConfig{
			ConsumerQos: 3,
		}, core.EnvConfig.Tls)
		defer webcore.DefaultRabbitMQ.Cleanup()

		forever := make(chan bool)

		webcore.DefaultRabbitMQ.ReceiveMessages(webcore.SHELL_RUNNER_QUEUE, processShellRunnerMessages)

		<-forever
	}
}

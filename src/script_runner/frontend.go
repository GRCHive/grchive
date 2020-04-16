package main

import (
	"encoding/json"
	"flag"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
)

func processScriptRunnerMessages(data []byte) *webcore.RabbitMQError {
	msg := webcore.ScriptRunnerMessage{}
	core.Info("RUN SCRIPT: " + string(data))
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return &webcore.RabbitMQError{err, false}
	}
	return &webcore.RabbitMQError{
		handleRun(msg.RunId, msg.Jar),
		false,
	}
}

func main() {
	core.Init()
	database.Init()
	webcore.InitializeWebcore()
	vault.Initialize(vault.VaultConfig{
		Url:      core.EnvConfig.Vault.Url,
		Username: core.EnvConfig.Vault.Username,
		Password: core.EnvConfig.Vault.Password,
	}, core.EnvConfig.Tls.Config())

	runId := flag.Int64("runId", -1, "Run ID in the database.")
	jar := flag.String("jar", "", "Client JAR to run.")
	local := flag.Bool("local", false, "Whether to use the runId and jar flags to run locally instead of listening using RabbitMQ.")
	flag.Parse()

	{
		err := pullKotlinImage(core.EnvConfig.Drone.RunnerImage)
		if err != nil {
			core.Error(err)
		}
	}

	if *local {
		err := handleRun(*runId, *jar)
		if err != nil {
			core.Error("Failed to run: " + err.Error())
		}
	} else {
		// This "1" is the number of jobs this worker can run at a time.
		webcore.DefaultRabbitMQ.Connect(*core.EnvConfig.RabbitMQ, webcore.MQClientConfig{
			ConsumerQos: 1,
		}, core.EnvConfig.Tls)
		defer webcore.DefaultRabbitMQ.Cleanup()

		forever := make(chan bool)

		webcore.DefaultRabbitMQ.ReceiveMessages(webcore.SCRIPT_RUNNER_QUEUE, processScriptRunnerMessages)

		<-forever
	}
}

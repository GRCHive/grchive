package main

import (
	"encoding/json"
	"flag"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/vault_api"
	"gitlab.com/grchive/grchive/webcore"
)

func processRefreshRequestMessage(data []byte) *webcore.RabbitMQError {
	msg := webcore.DatabaseRefreshMessage{}
	core.Info("REFRESH: " + string(data))
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return &webcore.RabbitMQError{err, false}
	}
	return processRefreshRequest(&msg.Refresh)
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

	refreshId := flag.Int64("refreshId", -1, "Refresh ID to retrieve data for. Will not read from RabbitMQ if specified.")
	orgId := flag.Int64("orgId", -1, "Org ID to retrieve data for. Will not read from RabbitMQ if specified.")
	flag.Parse()

	if *refreshId >= 0 && *orgId >= 0 {
		refresh, err := database.GetDatabaseRefresh(*refreshId, int32(*orgId), core.ServerRole)
		if err != nil {
			core.Error("Failed to get database refresh: " + err.Error())
		}

		rerr := processRefreshRequest(refresh)
		if rerr != nil {
			core.Error("Failed to process refresh: " + rerr.Err.Error())
		}
	} else {
		webcore.DefaultRabbitMQ.Connect(*core.EnvConfig.RabbitMQ, core.EnvConfig.Tls)
		defer webcore.DefaultRabbitMQ.Cleanup()

		forever := make(chan bool)

		webcore.DefaultRabbitMQ.ReceiveMessages(webcore.DATABASE_REFRESH_QUEUE, processRefreshRequestMessage)

		<-forever
	}
}

package main

import (
	"flag"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
)

func main() {
	immediate := flag.Bool("immediate", false, "Run all jobs immediately.")
	mul := flag.Int64("mul", -1, "If used, uses a sped up clock instead of the wall clock.")
	log := flag.Bool("log", false, "Extra debug logs in stdout.")
	flag.Parse()

	core.Init()
	database.Init()
	webcore.InitializeWebcore()

	var c core.Clock
	if *mul == -1 {
		c = core.RealClock{}
	} else {
		c = core.CreateMultiplierClock(*mul)
	}

	scheduler := CreateScheduler(c)
	scheduler.Log = *log

	webcore.DefaultRabbitMQ.Connect(*core.EnvConfig.RabbitMQ, webcore.MQClientConfig{
		ConsumerQos: 5,
	}, core.EnvConfig.Tls)
	defer webcore.DefaultRabbitMQ.Cleanup()
	webcore.DefaultRabbitMQ.ReceiveMessages(webcore.TASK_MANAGER_QUEUE, scheduler.handleRabbitMQMessage)

	// Load existing tasks from the database here as the
	// database listener will only tell us of changes to
	// these tasks.
	tasks, err := database.GetAllScheduledTasks(core.ServerRole)
	if err != nil {
		core.Error("Failed to grab initial tasks: " + err.Error())
	}

	for _, t := range tasks {
		j, err := createJob(t.Metadata, t.OneTime, t.Recurring, scheduler.Clock)
		if err != nil {
			core.Error("Failed to create job: " + err.Error())
		}

		err = scheduler.AddJob(j)
		if err != nil {
			core.Error("Failed to add job: " + err.Error())
		}
	}

	if *immediate {
		scheduler.RunImmediate(true)
	} else {
		scheduler.SyncRun()
	}
}

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
)

type ScheduledTaskChangeEvent struct {
	Task      core.ScheduledTaskMetadata
	OneTime   *core.ScheduledTaskOneTime
	Recurring *core.ScheduledTaskRecurrence
}

func createJob(t core.ScheduledTaskMetadata, oneTime *core.ScheduledTaskOneTime, recurring *core.ScheduledTaskRecurrence, c core.Clock) (*Job, error) {
	var schedule *Schedule
	var err error
	if oneTime != nil {
		schedule, err = CreateOneTimeJobSchedule(oneTime, c)
	} else if recurring != nil {
		schedule, err = CreateRecurringJobSchedule(recurring, c)
	} else {
		return nil, errors.New("No job schedule found.")
	}

	if err != nil {
		return nil, err
	}

	if schedule == nil {
		return nil, nil
	}

	job, err := CreateJobFromTaskMetadata(t, schedule)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func createOnNotifyScheduledTaskChange(s *Scheduler) database.ListenHandler {
	return func(data string) error {
		event := ScheduledTaskChangeEvent{}
		err := json.Unmarshal([]byte(data), &event)
		if err != nil {
			return err
		}

		job, err := createJob(event.Task, event.OneTime, event.Recurring, s.Clock)
		if err != nil {
			return err
		}

		return s.AddJob(job)
	}
}

type ScheduledTaskDeleteEvent struct {
	Task core.ScheduledTaskMetadata
}

func createOnNotifyScheduledTaskDelete(s *Scheduler) database.ListenHandler {
	return func(data string) error {
		event := ScheduledTaskDeleteEvent{}
		err := json.Unmarshal([]byte(data), &event)
		if err != nil {
			return err
		}

		s.RemoveJob(event.Task.Id)
		return nil
	}
}

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

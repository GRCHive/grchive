package webcore

import (
	"errors"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
)

func CreateScheduledTaskFromRawInputs(input *core.ScheduledTaskRawInput, task core.TaskType, data interface{}, userId int64, orgId int32) error {
	metadata := input.GenerateTaskMetadata(userId, orgId, task, data)

	if input.Repeat {
		var task *core.ScheduledTaskRecurrence
		var err error

		if input.Daily != nil {
			task, err = input.Daily.GenerateRecurringTasks()
		} else if input.Weekly != nil {
			task, err = input.Weekly.GenerateRecurringTasks()
		} else if input.Monthly != nil {
			task, err = input.Monthly.GenerateRecurringTasks()
		} else {
			return errors.New("No recurring task.")
		}

		if err != nil {
			return err
		}
		return database.CreateRecurringTask(metadata, task)
	} else {
		return database.CreateOneTimeTask(metadata, input.GenerateOneTimeTask())
	}
}

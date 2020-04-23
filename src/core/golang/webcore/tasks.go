package webcore

import (
	"errors"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"time"
)

type TaskLinkOptions struct {
	ScriptId core.NullInt64
}

func CreateScheduledTaskFromRawInputs(input *core.ScheduledTaskRawInput, task core.TaskType, data interface{}, userId int64, orgId int32, opts TaskLinkOptions) error {
	metadata := input.GenerateTaskMetadata(userId, orgId, task, data)

	tz, err := time.LoadLocation(input.Timezone)
	if err != nil {
		return err
	}

	tx := database.CreateTx()
	return database.WrapTx(tx, func() error {
		if input.Repeat {
			var task *core.ScheduledTaskRecurrence
			var err error

			if input.Daily != nil {
				task, err = input.Daily.GenerateRecurringTasks(tz)
			} else if input.Weekly != nil {
				task, err = input.Weekly.GenerateRecurringTasks(tz)
			} else if input.Monthly != nil {
				task, err = input.Monthly.GenerateRecurringTasks(tz)
			} else {
				return errors.New("No recurring task.")
			}

			if err != nil {
				return err
			}
			return database.CreateRecurringTaskWithTx(tx, metadata, task)
		} else {
			return database.CreateOneTimeTaskWithTx(tx, metadata, input.GenerateOneTimeTask(tz))
		}
	}, func() error {
		if opts.ScriptId.NullInt64.Valid {
			return database.LinkTaskToScriptWithTx(tx, metadata.Id, opts.ScriptId.NullInt64.Int64, orgId)
		}
		return nil
	})
}

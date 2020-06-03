package webcore

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"time"
)

type TaskLinkOptions struct {
	LinkId    core.NullInt64
	RequestId core.NullInt64
	DbId      core.NullInt64
}

func CreateScheduledTaskFromRawInputs(tx *sqlx.Tx, input *core.ScheduledTaskRawInput, task core.TaskType, data interface{}, userId int64, orgId int32, opts TaskLinkOptions) (*core.ScheduledTaskMetadata, error) {
	metadata := input.GenerateTaskMetadata(userId, orgId, task, data)

	tz, err := time.LoadLocation(input.Timezone)
	if err != nil {
		return nil, err
	}

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
			return nil, errors.New("No recurring task.")
		}

		if err != nil {
			return nil, err
		}

		err = database.CreateRecurringTaskWithTx(tx, metadata, task)
		if err != nil {
			return nil, err
		}
	} else {
		err = database.CreateOneTimeTaskWithTx(tx, metadata, input.GenerateOneTimeTask(tz))
		if err != nil {
			return nil, err
		}
	}

	if opts.LinkId.NullInt64.Valid {
		err = database.LinkTaskToScriptLinkWithTx(tx, metadata.Id, opts.LinkId.NullInt64.Int64)
		if err != nil {
			return nil, err
		}
	}

	if opts.RequestId.NullInt64.Valid {
		err = database.LinkScheduledTaskToRequestWithTx(tx, metadata.Id, opts.RequestId.NullInt64.Int64)
		if err != nil {
			return nil, err
		}
	}

	if opts.DbId.NullInt64.Valid {
		err = database.LinkDatabaseSettingToScheduledTaskWithTx(tx, opts.DbId.NullInt64.Int64, metadata.Id)
		if err != nil {
			return nil, err
		}
	}

	return metadata, nil
}

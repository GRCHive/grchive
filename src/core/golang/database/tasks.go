package database

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/teambition/rrule-go"
	"gitlab.com/grchive/grchive/core"
	"time"
)

func createScheduledTaskMetadataAdder(tx *sqlx.Tx, meta *core.ScheduledTaskMetadata) func() error {
	return func() error {
		rawData, err := json.Marshal(meta.TaskData)
		if err != nil {
			return err
		}

		rows, err := tx.Queryx(`
			INSERT INTO scheduled_tasks (name, description, org_id, user_id, task_type, task_data, scheduled_time)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
			RETURNING id
		`, meta.Name, meta.Description, meta.OrgId, meta.UserId, meta.TaskType, string(rawData))

		if err != nil {
			return err
		}

		defer rows.Close()
		rows.Next()
		return rows.Scan(&meta.Id)
	}
}

func CreateOneTimeTaskWithTx(tx *sqlx.Tx, meta *core.ScheduledTaskMetadata, data *core.ScheduledTaskOneTime) error {
	err := createScheduledTaskMetadataAdder(tx, meta)()

	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO one_time_tasks (event_id, event_date_time)
		VALUES ($1, $2)
	`, meta.Id, data.EventDateTime.UTC())
	return err
}

func CreateRecurringTaskWithTx(tx *sqlx.Tx, meta *core.ScheduledTaskMetadata, data *core.ScheduledTaskRecurrence) error {
	err := createScheduledTaskMetadataAdder(tx, meta)()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO recurring_tasks (event_id, start_date_time, rrule, timezone)
		VALUES ($1, $2, $3, $4)
	`, meta.Id, data.StartDateTime.UTC(), data.RRule.String(), data.Timezone)
	return err
}

func LinkTaskToScriptWithTx(tx *sqlx.Tx, taskId int64, scriptId int64, orgId int32) error {
	_, err := tx.Exec(`
		INSERT INTO scheduled_task_script_links (event_id, org_id, script_id)
		VALUES ($1, $2, $3)
	`, taskId, orgId, scriptId)
	return err
}

func getScheduledTasksHelper(role *core.Role, condition string, args ...interface{}) ([]*core.FullScheduledTask, error) {
	rows, err := dbConn.Queryx(fmt.Sprintf(`
		SELECT
			t.id,
			t.name,
			t.description,
			t.org_id,
			t.user_id,
			t.task_type,
			t.task_data,
			t.scheduled_time,
			ot.event_date_time,
			rt.start_date_time,
			rt.rrule,
			rt.timezone
		FROM scheduled_tasks AS t
		LEFT JOIN one_time_tasks AS ot
			ON ot.event_id = t.id
		LEFT JOIN recurring_tasks AS rt
			ON rt.event_id = t.id
		%s
	`, condition), args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*core.FullScheduledTask, 0)
	for rows.Next() {
		tk := core.FullScheduledTask{}

		taskData := types.JSONText{}

		// Temporary nullable elements to scan into for the one-time data and recurring task data.
		var oneTimeTime core.NullTime
		var recurringTime core.NullTime
		var recurringRule core.NullString
		var recurringTimezone core.NullString

		err = rows.Scan(
			&tk.Metadata.Id,
			&tk.Metadata.Name,
			&tk.Metadata.Description,
			&tk.Metadata.OrgId,
			&tk.Metadata.UserId,
			&tk.Metadata.TaskType,
			&taskData,
			&tk.Metadata.ScheduledTime,
			&oneTimeTime,
			&recurringTime,
			&recurringRule,
			&recurringTimezone,
		)

		if err != nil {
			return nil, err
		}

		err = taskData.Unmarshal(&tk.Metadata.TaskData)
		if err != nil {
			return nil, err
		}

		if oneTimeTime.NullTime.Valid {
			tk.OneTime = &core.ScheduledTaskOneTime{
				EventId:       tk.Metadata.Id,
				EventDateTime: oneTimeTime.NullTime.Time,
			}
		} else {
			rule, err := rrule.StrToRRuleSet(recurringRule.NullString.String)
			if err != nil {
				return nil, err
			}

			loc, err := time.LoadLocation(recurringTimezone.NullString.String)
			if err != nil {
				return nil, err
			}

			// Recreate the ruleset because we need to make sure the DTStart gets
			// set properly to the correct timezone.
			//
			// Doing the naive:
			// 	rule.DTStart(rule.GetDTStart().In(loc))
			// seem to cause crashes later on...
			rule, err = core.RebuildRRuleSet(*rule, func(opt *rrule.ROption) {
				opt.Dtstart = opt.Dtstart.In(loc)
				opt.RFC = false
			})

			tk.Recurring = &core.ScheduledTaskRecurrence{
				EventId:       tk.Metadata.Id,
				StartDateTime: recurringTime.NullTime.Time,
				RRule:         *rule,
				Timezone:      recurringTimezone.NullString.String,
			}
		}

		tasks = append(tasks, &tk)
	}

	return tasks, nil

}

func GetAllScheduledTasks(role *core.Role) ([]*core.FullScheduledTask, error) {
	return getScheduledTasksHelper(role, "")
}

func GetAllScheduledTasksForOrgId(orgId int32, role *core.Role) ([]*core.FullScheduledTask, error) {
	return getScheduledTasksHelper(role, "WHERE t.org_id = $1", orgId)
}

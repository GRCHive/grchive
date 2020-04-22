package database

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func createScheduledTaskMetadataAdder(tx *sqlx.Tx, meta *core.ScheduledTaskMetadata) func() error {
	return func() error {
		rawData, err := json.Marshal(meta.TaskData)
		if err != nil {
			return err
		}

		rows, err := tx.Queryx(`
			INSERT INTO scheduled_tasks (name, description, org_id, user_id, task_type, task_data)
			VALUES ($1, $2, $3, $4, $5, $6)
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

func CreateOneTimeTask(meta *core.ScheduledTaskMetadata, data *core.ScheduledTaskOneTime) error {
	tx := CreateTx()
	return WrapTx(tx, createScheduledTaskMetadataAdder(tx, meta), func() error {
		_, err := tx.Exec(`
			INSERT INTO one_time_tasks (event_id, event_date_time)
			VALUES ($1, $2)
		`, meta.Id, data.EventDateTime.UTC())
		return err
	})
}

func CreateRecurringTask(meta *core.ScheduledTaskMetadata, data *core.ScheduledTaskRecurrence) error {
	tx := CreateTx()
	return WrapTx(tx, createScheduledTaskMetadataAdder(tx, meta), func() error {
		_, err := tx.Exec(`
			INSERT INTO recurring_tasks (event_id, start_date_time, rrule)
			VALUES ($1, $2, $3)
		`, meta.Id, data.StartDateTime.UTC(), data.RRule.String())
		return err
	})
}

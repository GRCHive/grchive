package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"strings"
)

func InsertNotificationWithTx(notification *core.Notification, tx *sqlx.Tx) error {
	rows, err := tx.NamedQuery(`
		INSERT INTO notifications (
			org_id,
			time,
			subject_type,
			subject_id,
			verb,
			object_type,
			object_id,
			indirect_object_type,
			indirect_object_id
		)
		VALUES (
			:org_id,
			:time,
			:subject_type,
			:subject_id,
			:verb,
			:object_type,
			:object_id,
			:indirect_object_type,
			:indirect_object_id
		)
		RETURNING id
	`, notification)

	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()

	err = rows.Scan(&notification.Id)
	if err != nil {
		return err
	}

	return nil
}

func LinkNotificationToUsersWithTx(notificationId int64, orgId int32, users []*core.User, tx *sqlx.Tx) error {
	if len(users) == 0 {
		return nil
	}

	builder := strings.Builder{}
	builder.WriteString(`
		INSERT INTO user_notifications (user_id, org_id, notification_id)
		VALUES
	`)

	for _, u := range users {
		builder.WriteString(fmt.Sprintf(`
			($1, $2, %d)
		`, u.Id))
	}

	_, err := tx.Exec(builder.String(), notificationId, orgId)
	return err
}

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
		INSERT INTO user_notifications (notification_id, org_id, user_id)
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

func AllNotificationsForUserId(userId int64) ([]*core.NotificationWrapper, error) {
	notifications := make([]*core.NotificationWrapper, 0)
	err := dbConn.Select(&notifications, `
		SELECT
			n.id AS "notif.id",
			n.org_id AS "notif.org_id",
			n.time AS "notif.time",
			n.subject_type AS "notif.subject_type",
			n.subject_id AS "notif.subject_id",
			n.verb AS "notif.verb",
			n.object_type AS "notif.object_type",
			n.object_id AS "notif.object_id",
			n.indirect_object_type AS "notif.indirect_object_type",
			n.indirect_object_id AS "notif.indirect_object_id",
			o.org_group_name AS "org_group_name",
			un.read IS NOT NULL AS "read"
		FROM notifications AS n
		INNER JOIN user_notifications AS un
			ON un.notification_id = n.id
		INNER JOIN organizations AS o
			ON o.id = un.org_id
		WHERE un.user_id = $1
		ORDER BY n.time DESC
	`, userId)
	return notifications, err
}

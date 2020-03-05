package database

import (
	"github.com/jmoiron/sqlx/types"
	"gitlab.com/grchive/grchive/core"
)

func AllFilteredAuditEvents(orgId int32, role *core.Role) ([]*core.AuditEvent, error) {
	rows, err := dbConn.Queryx(`
		SELECT
			hist.id,
			hist.org_id,
			hist.resource_type,
			hist.resource_id,
			hist.resource_extra_data,
			hist.action,
			hist.performed_at,
			u.id AS "user_id"
		FROM global_audit_event_history AS hist
		LEFT JOIN postgres_oid_to_users AS lnk
			ON lnk.pg_oid = hist.pgrole_id
		LEFT JOIN users AS u
			ON u.id = lnk.user_id
		WHERE hist.org_id = $1
		ORDER BY hist.performed_at DESC
	`, orgId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*core.AuditEvent, 0)
	for rows.Next() {
		e := core.AuditEvent{}
		extraData := types.JSONText{}

		err = rows.Scan(
			&e.Id,
			&e.OrgId,
			&e.ResourceType,
			&e.ResourceId,
			&extraData,
			&e.Action,
			&e.PerformedAt,
			&e.UserId,
		)

		if err != nil {
			return nil, err
		}

		err = extraData.Unmarshal(&e.ResourceExtraData)
		if err != nil {
			return nil, err
		}

		events = append(events, &e)
	}

	return events, err
}

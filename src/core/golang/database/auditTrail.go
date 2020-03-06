package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"gitlab.com/grchive/grchive/core"
)

func commonAuditEventRetrievalQuery(role *core.Role, cond string, limit string, args ...interface{}) ([]*core.AuditEvent, error) {
	rows, err := dbConn.Queryx(fmt.Sprintf(`
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
		%s
		ORDER BY hist.id DESC
		LIMIT %s
	`, cond, limit), args...)

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

func AllFilteredAuditEvents(orgId int32, role *core.Role) ([]*core.AuditEvent, error) {
	return commonAuditEventRetrievalQuery(role, "WHERE hist.org_id = $1", "ALL", orgId)
}

func GetAuditEvent(eventId int64, role *core.Role) (*core.AuditEvent, error) {
	events, err := commonAuditEventRetrievalQuery(role, "WHERE hist.id = $1", "1", eventId)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, nil
	}

	return events[0], nil
}

func GetLatestAuditEvent(resourceType string, resourceId string, role *core.Role) (*core.AuditEvent, error) {
	events, err := commonAuditEventRetrievalQuery(role, `
		WHERE hist.resource_type = $1 AND hist.resource_id = $2
 	`, "1", resourceType, resourceId)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, nil
	}

	return events[0], nil
}

func getAuditModificationFromStmt(stmt *sqlx.Stmt, args ...interface{}) (map[string]interface{}, error) {
	rows, err := stmt.Queryx(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	data := types.JSONText{}
	err = rows.Scan(&data)
	if err != nil {
		return nil, err
	}

	ret := map[string]interface{}{}
	err = data.Unmarshal(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func GetAuditModificationHistoryData(eventId int64, role *core.Role) (map[string]interface{}, error) {
	stmt, err := dbConn.Preparex(`
		SELECT data
		FROM audit_resource_modifications
		WHERE event_id = $1
	`)

	if err != nil {
		return nil, err
	}

	return getAuditModificationFromStmt(stmt, eventId)
}

func GetLatestAuditModificationHistoryData(resourceType string, resourceId string, role *core.Role) (map[string]interface{}, error) {
	stmt, err := dbConn.Preparex(`
		SELECT mod.data
		FROM audit_resource_modifications AS mod
		INNER JOIN global_audit_event_history AS his
			ON his.id = mod.event_id
		WHERE his.resource_type = $1
			AND his.resource_id = $2
		ORDER BY mod.event_id DESC
		LIMIT 1
	`)

	if err != nil {
		return nil, err
	}

	return getAuditModificationFromStmt(stmt, resourceType, resourceId)
}

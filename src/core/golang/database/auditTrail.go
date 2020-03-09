package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"gitlab.com/grchive/grchive/core"
	"strconv"
	"strings"
)

func commonAuditEventRetrievalQuery(
	role *core.Role,
	retrieve core.AuditTrailRetrievalParams,
	sort core.AuditTrailSortParams) ([]*core.AuditEvent, error) {

	selectParams := strings.Builder{}
	args := make([]interface{}, 0)

	if retrieve.OrgId.NullInt32.Valid {
		selectParams.WriteString("WHERE hist.org_id = $1\n")
		args = []interface{}{retrieve.OrgId.NullInt32.Int32}
	} else if retrieve.EventId.NullInt64.Valid {
		selectParams.WriteString("WHERE hist.id = $1\n")
		args = []interface{}{retrieve.EventId.NullInt64.Int64}
	} else if retrieve.ResourceType.NullString.Valid && retrieve.ResourceId.NullString.Valid {
		selectParams.WriteString("WHERE hist.resource_type = $1 AND hist.resource_id = $2\n")
		args = []interface{}{retrieve.ResourceType.NullString.String, retrieve.ResourceId.NullString.String}
	} else {
		return nil, errors.New("Invalid retrieval parameters.")
	}

	if len(sort.SortColumns) > 0 && sort.SortDirection.NullString.Valid {
		selectParams.WriteString("ORDER BY\n")
		for idx, c := range sort.SortColumns {
			selectParams.WriteString(fmt.Sprintf("%s %s", c, sort.SortDirection.NullString.String))
			if idx != len(sort.SortColumns)-1 {
				selectParams.WriteString(",")
			}
			selectParams.WriteString("\n")
		}
	} else {
		selectParams.WriteString("ORDER BY hist.id DESC\n")
	}

	if sort.Limit.NullInt32.Valid {
		selectParams.WriteString(fmt.Sprintf("LIMIT %d\n", sort.Limit.NullInt32.Int32))

		if sort.Page.NullInt32.Valid {
			selectParams.WriteString(fmt.Sprintf("OFFSET %d\n",
				int64(sort.Limit.NullInt32.Int32)*int64(sort.Page.NullInt32.Int32)))
		}
	}

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
	`, selectParams.String()), args...)

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

func AllFilteredAuditEvents(orgId int32, sort core.AuditTrailSortParams, role *core.Role) ([]*core.AuditEvent, error) {
	return commonAuditEventRetrievalQuery(
		role,
		core.AuditTrailRetrievalParams{
			OrgId: core.CreateNullInt32(orgId),
		},
		sort)
}

func CountFilteredAuditEvents(orgId int32, role *core.Role) (int, error) {
	rows, err := dbConn.Queryx(`
		SELECT COUNT(*)
		FROM global_audit_event_history
		WHERE org_id = $1
	`, orgId)

	if err != nil {
		return -1, err
	}
	defer rows.Close()

	val := int(0)
	rows.Next()
	err = rows.Scan(&val)
	if err != nil {
		return -1, err
	}
	return val, nil
}

func GetAuditEvent(eventId int64, role *core.Role) (*core.AuditEvent, error) {
	events, err := commonAuditEventRetrievalQuery(
		role,
		core.AuditTrailRetrievalParams{
			EventId: core.CreateNullInt64(eventId),
		},
		core.AuditTrailSortParams{
			Limit: core.CreateNullInt32(1),
		})
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, nil
	}

	return events[0], nil
}

func GetLatestAuditEvent(resourceType string, resourceId string, role *core.Role) (*core.AuditEvent, error) {
	events, err := commonAuditEventRetrievalQuery(
		role,
		core.AuditTrailRetrievalParams{
			ResourceType: core.CreateNullString(resourceType),
			ResourceId:   core.CreateNullString(resourceId),
		},
		core.AuditTrailSortParams{
			Limit: core.CreateNullInt32(1),
		})
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

func retrieveResourceExtraData(resourceType string, resourceId string) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	var rows *sqlx.Rows
	var err error

	switch resourceType {
	case core.ResourceFileStorage:
		storageId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT f.category_id AS "cat_id", s.metadata_id AS "file_id"
			FROM file_storage AS s
			INNER JOIN file_metadata AS f
				ON f.id = s.metadata_id
			WHERE s.id = $1
		`, storageId)
	case core.ResourceDocMetadata:
		fileId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT r.category_id AS "cat_id"
			FROM file_metadata AS r
			WHERE r.id = $1
		`, fileId)
	case core.ResourceVendorProduct:
		productId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT p.vendor_id AS "vendor_id"
			FROM vendor_products AS p
			WHERE p.id = $1
		`, productId)
	case core.ResourceGLCat:
		catId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT c.parent_category_id AS "gl_parent_cat_id"
			FROM general_ledger_categories AS c
			WHERE c.id = $1
		`, catId)
	case core.ResourceGLAcc:
		accId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT a.parent_category_id AS "gl_parent_cat_id"
			FROM general_ledger_accounts AS a
			WHERE a.id = $1
		`, accId)
	case core.ResourceFlowNode:
		nodeId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT n.process_flow_id AS "process_flow_id"
			FROM process_flow_nodes AS n
			WHERE n.id = $1
		`, nodeId)
	case core.ResourceFlowNodeInput:
		ioId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT n.process_flow_id AS "process_flow_id", io.parent_node_id AS "node_id"
			FROM process_flow_node_inputs AS io
			INNER JOIN process_flow_nodes AS n
				ON n.id = io.parent_node_id
			WHERE id = $1
		`, ioId)
	case core.ResourceFlowNodeOutput:
		ioId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT n.process_flow_id AS "process_flow_id", io.parent_node_id AS "node_id"
			FROM process_flow_node_outputs AS io
			INNER JOIN process_flow_nodes AS n
				ON n.id = io.parent_node_id
			WHERE id = $1
		`, ioId)
	case core.ResourceSqlQueryMetadata:
		queryId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT r.db_id AS "db_id"
			FROM database_sql_metadata AS r
			WHERE r.id = $1
		`, queryId)
	case core.ResourceSqlQuery:
		queryId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT r.metadata_id AS "sql_metadata_id", m.db_id AS "db_id"
			FROM database_sql_metadata AS r
			INNER JOIN database_sql_metadata AS m
				ON m.id = r.metadata_id
			WHERE r.id = $1
		`, queryId)
	case core.ResourceDatabaseConn:
		connId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT r.db_id AS "db_id"
			FROM database_connection_info AS r
			WHERE r.id = $1
		`, connId)
	case core.ResourceSqlQueryRequest:
		reqId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT r.query_id AS "query_id", q.metadata_id AS "sql_metadata_id"
			FROM database_sql_query_requests AS r
			INNER JOIN database_sql_queries AS q
				ON q.id = r.query_id
			WHERE r.id = $1
		`, reqId)
	}

	if err != nil {
		return nil, err
	}

	if rows != nil {
		defer rows.Close()
		rows.Next()

		err = rows.MapScan(data)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func LogAuditSelectWithTx(orgId int32, resourceType string, resourceId string, role *core.Role, tx *sqlx.Tx) error {
	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
	}

	extraData, err := retrieveResourceExtraData(resourceType, resourceId)
	if err != nil {
		return err
	}

	rawExtraData, err := json.Marshal(extraData)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO global_audit_event_history(
			org_id,
			resource_type,
			resource_id,
			resource_extra_data,
			action,
			performed_at,
			pgrole_id
		)
		SELECT 
			$1,
			$2,
			$3,
			$4,
			'SELECT',
			NOW(),
			pg.oid
		FROM pg_roles AS pg
		WHERE pg.rolname = current_user
	`, orgId, resourceType, resourceId, string(rawExtraData))
	if err != nil {
		return err
	}

	return nil
}

func LogAuditSelect(orgId int32, resourceType string, resourceId string, role *core.Role) error {
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	err = LogAuditSelectWithTx(orgId, resourceType, resourceId, role, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

package database

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"gitlab.com/grchive/grchive/core"
	"strconv"
	"strings"
)

func buildHistResourceFilter(types []string, ids []string) string {
	filter := strings.Builder{}
	filter.WriteString("(")

	for i, typ := range types {
		idStr := ""
		if ids[i] == "*" {
			idStr = "IS NOT NULL"
		} else {
			idStr = fmt.Sprintf("= '%s'", ids[i])
		}

		filter.WriteString(fmt.Sprintf(`
			(hist.resource_type = '%s' AND hist.resource_id %s)
		`,
			typ,
			idStr,
		))

		subTyp := ""
		switch typ {
		case core.ResourceIdDocMetadata:
			subTyp = "file_id"
		case core.ResourceIdDocCat:
			subTyp = "cat_id"
		case core.ResourceIdVendor:
			subTyp = "vendor_id"
		case core.ResourceIdGLCat:
			subTyp = "gl_parent_cat_id"
		case core.ResourceIdProcessFlow:
			subTyp = "process_flow_id"
		case core.ResourceIdFlowNode:
			subTyp = "node_id"
		case core.ResourceIdDatabase:
			subTyp = "db_id"
		case core.ResourceIdSqlQueryMetadata:
			subTyp = "sql_metadata_id"
		case core.ResourceIdSqlQuery:
			subTyp = "query_id"
		case core.ResourceIdClientData:
			subTyp = "client_data_id"
		case core.ResourceIdClientScripts:
			subTyp = "client_script_id"
		}

		if len(subTyp) > 0 {
			filter.WriteString(fmt.Sprintf(`
				OR (hist.resource_extra_data ->> '%s' %s)
			`, subTyp, idStr))
		}

		if i == len(types)-1 {
			filter.WriteString("\n")
		} else {
			filter.WriteString("OR\n")
		}
	}

	filter.WriteString(")")
	return filter.String()
}

func commonAuditEventRetrievalQuery(
	role *core.Role,
	retrieve core.AuditTrailRetrievalParams,
	sort core.AuditTrailSortParams,
	filter core.AuditTrailFilterData) ([]*core.AuditEvent, error) {

	selectParams := strings.Builder{}
	args := make([]interface{}, 0)

	if retrieve.OrgId.NullInt32.Valid {
		selectParams.WriteString("WHERE hist.org_id = $1\n")
		args = []interface{}{retrieve.OrgId.NullInt32.Int32}
	} else if retrieve.EventId.NullInt64.Valid {
		selectParams.WriteString("WHERE hist.id = $1\n")
		args = []interface{}{retrieve.EventId.NullInt64.Int64}
	} else if len(retrieve.ResourceType) > 0 && len(retrieve.ResourceId) > 0 {
		selectParams.WriteString(fmt.Sprintf("WHERE %s\n", buildHistResourceFilter(retrieve.ResourceType, retrieve.ResourceId)))
	} else {
		return nil, errors.New("Invalid retrieval parameters.")
	}

	selectParams.WriteString(fmt.Sprintf("AND %s\n", buildStringFilter("audit_resource_type_to_human_name(hist.resource_type)", filter.ResourceTypeFilter)))
	selectParams.WriteString(fmt.Sprintf("AND %s\n", buildStringFilter("hist.action", filter.ActionFilter)))
	selectParams.WriteString(fmt.Sprintf("AND %s\n", buildStringFilter("COALESCE(user_to_human_name(u.*::users), 'No User')", filter.UserFilter)))
	selectParams.WriteString(fmt.Sprintf("AND %s\n", buildTimeRangeFilter("hist.performed_at", filter.TimeRangeFilter)))

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

	queryString := fmt.Sprintf(`
		SELECT DISTINCT
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
	`, selectParams.String())
	rows, err := dbConn.Queryx(queryString, args...)

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

func AllFilteredAuditEvents(orgId int32, sort core.AuditTrailSortParams, filter core.AuditTrailFilterData, role *core.Role) ([]*core.AuditEvent, error) {
	return commonAuditEventRetrievalQuery(
		role,
		core.AuditTrailRetrievalParams{
			OrgId: core.CreateNullInt32(orgId),
		},
		sort,
		filter)
}

func AllFilteredAuditEventsForResource(resourceType []string, resourceId []string, sort core.AuditTrailSortParams, filter core.AuditTrailFilterData, role *core.Role) ([]*core.AuditEvent, error) {
	return commonAuditEventRetrievalQuery(
		role,
		core.AuditTrailRetrievalParams{
			ResourceType: resourceType,
			ResourceId:   resourceId,
		},
		sort,
		filter)
}

func CountFilteredAuditEvents(orgId int32, filter core.AuditTrailFilterData, role *core.Role) (int, error) {
	// How to prevent duplicate code here with the above function?
	rows, err := dbConn.Queryx(fmt.Sprintf(`
		SELECT COUNT(hist.*)
		FROM global_audit_event_history AS hist
		LEFT JOIN postgres_oid_to_users AS lnk
			ON lnk.pg_oid = hist.pgrole_id
		LEFT JOIN users AS u
			ON u.id = lnk.user_id
		WHERE org_id = $1
			AND %s
			AND %s
			AND %s
			AND %s
	`,
		buildStringFilter("audit_resource_type_to_human_name(hist.resource_type)", filter.ResourceTypeFilter),
		buildStringFilter("hist.action", filter.ActionFilter),
		buildStringFilter("user_to_human_name(u.*::users)", filter.UserFilter),
		buildTimeRangeFilter("hist.performed_at", filter.TimeRangeFilter)),
		orgId)

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

func CountFilteredAuditEventsForResource(resourceType []string, resourceId []string, filter core.AuditTrailFilterData, role *core.Role) (int, error) {
	// How to prevent duplicate code here with the above function?
	rows, err := dbConn.Queryx(fmt.Sprintf(`
		SELECT COUNT(DISTINCT hist.*)
		FROM global_audit_event_history AS hist
		LEFT JOIN postgres_oid_to_users AS lnk
			ON lnk.pg_oid = hist.pgrole_id
		LEFT JOIN users AS u
			ON u.id = lnk.user_id
		WHERE %s
			AND %s
			AND %s
			AND %s
			AND %s
	`,
		buildHistResourceFilter(resourceType, resourceId),
		buildStringFilter("audit_resource_type_to_human_name(hist.resource_type)", filter.ResourceTypeFilter),
		buildStringFilter("hist.action", filter.ActionFilter),
		buildStringFilter("user_to_human_name(u.*::users)", filter.UserFilter),
		buildTimeRangeFilter("hist.performed_at", filter.TimeRangeFilter)))

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
		},
		core.NullAuditTrailFilterData)
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
			ResourceType: []string{resourceType},
			ResourceId:   []string{resourceId},
		},
		core.AuditTrailSortParams{
			Limit: core.CreateNullInt32(1),
		},
		core.NullAuditTrailFilterData)
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

func GetModificationDiffFromEventId(eventId int64, role *core.Role) (map[string]interface{}, error) {
	rows, err := dbConn.Queryx(`
		SELECT diff
		FROM audit_resource_modification_diffs
		WHERE event_id = $1
	`, eventId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	rows.Next()

	diffData := types.JSONText{}
	err = rows.Scan(&diffData)
	if err != nil {
		return nil, err
	}

	ret := map[string]interface{}{}
	err = diffData.Unmarshal(&ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func retrieveResourceExtraData(resourceType string, resourceId string) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	var rows *sqlx.Rows
	var err error

	switch resourceType {
	case core.ResourceIdFileStorage:
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
	case core.ResourceIdDocMetadata:
		fileId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT r.category_id AS "cat_id"
			FROM file_metadata AS r
			WHERE r.id = $1
		`, fileId)
	case core.ResourceIdVendorProduct:
		productId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT p.vendor_id AS "vendor_id"
			FROM vendor_products AS p
			WHERE p.id = $1
		`, productId)
	case core.ResourceIdGLCat:
		catId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT c.parent_category_id AS "gl_parent_cat_id"
			FROM general_ledger_categories AS c
			WHERE c.id = $1
		`, catId)
	case core.ResourceIdGLAcc:
		accId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT a.parent_category_id AS "gl_parent_cat_id"
			FROM general_ledger_accounts AS a
			WHERE a.id = $1
		`, accId)
	case core.ResourceIdFlowNode:
		nodeId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT n.process_flow_id AS "process_flow_id"
			FROM process_flow_nodes AS n
			WHERE n.id = $1
		`, nodeId)
	case core.ResourceIdFlowNodeInput:
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
	case core.ResourceIdFlowNodeOutput:
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
	case core.ResourceIdSqlQueryMetadata:
		queryId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT r.db_id AS "db_id"
			FROM database_sql_metadata AS r
			WHERE r.id = $1
		`, queryId)
	case core.ResourceIdSqlQuery:
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
	case core.ResourceIdDatabaseConn:
		connId, err := strconv.ParseInt(resourceId, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = dbConn.Queryx(`
			SELECT r.db_id AS "db_id"
			FROM database_connection_info AS r
			WHERE r.id = $1
		`, connId)
	case core.ResourceIdSqlQueryRequest:
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

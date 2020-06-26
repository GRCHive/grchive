package database

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
)

func GetOverallPbcAnalytics(orgId int32, filter core.DocRequestFilterData) (map[core.DocRequestStatus]int32, error) {
	rows, err := dbConn.Queryx(fmt.Sprintf(`
		SELECT data.status, COUNT(data.id)
		FROM (
		SELECT get_pbc_request_status(req.*) AS status, req.id
		FROM document_requests AS req
		WHERE req.org_id = $1
			AND %s
		) AS data
		GROUP BY data.status
	`, buildDocRequestFilter("req", filter)), orgId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := map[core.DocRequestStatus]int32{}
	for rows.Next() {

		var status core.DocRequestStatus
		var count int32

		err = rows.Scan(&status, &count)
		if err != nil {
			return nil, err
		}

		data[status] = count
	}

	return data, nil
}

type PbcCategoryAnalyticsResult struct {
	Name string
	Data map[core.DocRequestStatus]int32
}

func GetGenericCategoryPbcAnalytics(query string, orgId int32, filter core.DocRequestFilterData) ([]PbcCategoryAnalyticsResult, error) {
	rows, err := dbConn.Queryx(fmt.Sprintf(query, buildDocRequestFilter("req", filter)), orgId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order := make([]string, 0)
	mapResults := map[string]*PbcCategoryAnalyticsResult{}
	for rows.Next() {
		var id string
		var status core.DocRequestStatus
		var count int32

		err = rows.Scan(&id, &status, &count)
		if err != nil {
			return nil, err
		}

		result, ok := mapResults[id]
		if !ok {
			mapResults[id] = &PbcCategoryAnalyticsResult{}
			result = mapResults[id]
			result.Name = id
			result.Data = map[core.DocRequestStatus]int32{}
			order = append(order, id)
		}

		result.Data[status] = count
	}

	data := make([]PbcCategoryAnalyticsResult, 0)
	for _, v := range order {
		data = append(data, *mapResults[v])
	}

	return data, nil

}

func GetAssigneeCategoryPbcAnalytics(orgId int32, filter core.DocRequestFilterData) ([]PbcCategoryAnalyticsResult, error) {
	return GetGenericCategoryPbcAnalytics(`
		SELECT 
			COALESCE(u.first_name || ' ' || u.last_name, 'No User') AS "cid",
			data.status,
			COUNT(data.id)
		FROM (
		SELECT get_pbc_request_status(req.*) AS status, req.id, req.assignee
		FROM document_requests AS req
		WHERE req.org_id = $1
			AND %s
		) AS data
		LEFT JOIN users AS u
			ON u.id = data.assignee
		GROUP BY data.status, cid
		ORDER BY cid
	`, orgId, filter)
}

func GetRequesterCategoryPbcAnalytics(orgId int32, filter core.DocRequestFilterData) ([]PbcCategoryAnalyticsResult, error) {
	return GetGenericCategoryPbcAnalytics(`
		SELECT 
			COALESCE(u.first_name || ' ' || u.last_name, 'No User') AS "cid",
			data.status,
			COUNT(data.id)
		FROM (
		SELECT get_pbc_request_status(req.*) AS status, req.id, req.requested_user_id
		FROM document_requests AS req
		WHERE req.org_id = $1
			AND %s
		) AS data
		LEFT JOIN users AS u
			ON u.id = data.requested_user_id
		GROUP BY data.status, cid
		ORDER BY cid
	`, orgId, filter)
}

func GetDocCatCategoryPbcAnalytics(orgId int32, filter core.DocRequestFilterData) ([]PbcCategoryAnalyticsResult, error) {
	return GetGenericCategoryPbcAnalytics(`
		SELECT 
			COALESCE(cat.name, 'N/A') AS "cid",
			data.status,
			COUNT(data.id)
		FROM (
		SELECT get_pbc_request_status(req.*) AS status, req.id, lnk.cat_id
		FROM document_requests AS req
		LEFT JOIN request_doc_cat_link AS lnk
			ON lnk.request_id = req.id
		WHERE req.org_id = $1
			AND %s
		) AS data
		LEFT JOIN process_flow_control_documentation_categories AS cat
			ON cat.id = data.cat_id
		GROUP BY data.status, cid
		ORDER BY cid
	`, orgId, filter)
}

func GetProcessFlowCategoryPbcAnalytics(orgId int32, filter core.DocRequestFilterData) ([]PbcCategoryAnalyticsResult, error) {
	return GetGenericCategoryPbcAnalytics(`
		SELECT 
			COALESCE(data.name, 'N/A') AS "cid",
			data.status,
			COUNT(data.id)
		FROM (
		SELECT DISTINCT get_pbc_request_status(req.*) AS status, req.id, flow.name
		FROM document_requests AS req
		LEFT JOIN request_control_link AS lnk
			ON lnk.request_id = req.id
		LEFT JOIN process_flow_control_node AS cn
			ON cn.control_id = lnk.control_id
		LEFT JOIN process_flow_nodes AS fn
			ON fn.id = cn.node_id
		LEFT JOIN process_flows AS flow
			ON flow.id = fn.process_flow_id
		WHERE req.org_id = $1
			AND %s
		) AS data
		GROUP BY data.status, cid
		ORDER BY cid
	`, orgId, filter)
}

func GetControlCategoryPbcAnalytics(orgId int32, filter core.DocRequestFilterData) ([]PbcCategoryAnalyticsResult, error) {
	return GetGenericCategoryPbcAnalytics(`
		SELECT 
			COALESCE(ctrl.identifier, 'N/A') AS "cid",
			data.status,
			COUNT(data.id)
		FROM (
		SELECT get_pbc_request_status(req.*) AS status, req.id, lnk.control_id
		FROM document_requests AS req
		LEFT JOIN request_control_link AS lnk
			ON lnk.request_id = req.id
		WHERE req.org_id = $1
			AND %s
		) AS data
		LEFT JOIN process_flow_controls AS ctrl
			ON ctrl.id = data.control_id
		GROUP BY data.status, cid
		ORDER BY cid
	`, orgId, filter)
}

func GetRiskCategoryPbcAnalytics(orgId int32, filter core.DocRequestFilterData) ([]PbcCategoryAnalyticsResult, error) {
	return GetGenericCategoryPbcAnalytics(`
		SELECT 
			COALESCE(rsk.identifier, 'N/A') AS "cid",
			data.status,
			COUNT(data.id)
		FROM (
		SELECT get_pbc_request_status(req.*) AS status, req.id, lnk.control_id
		FROM document_requests AS req
		LEFT JOIN request_control_link AS lnk
			ON lnk.request_id = req.id
		WHERE req.org_id = $1
			AND %s
		) AS data
		LEFT JOIN process_flow_risk_control AS rc
			ON rc.control_id = data.control_id
		LEFT JOIN process_flow_risks AS rsk
			ON rsk.id = rc.risk_id
		GROUP BY data.status, cid
		ORDER BY cid
	`, orgId, filter)
}

func GetGLCategoryPbcAnalytics(orgId int32, filter core.DocRequestFilterData) ([]PbcCategoryAnalyticsResult, error) {
	return GetGenericCategoryPbcAnalytics(`
		SELECT
			COALESCE(data.account_name, 'N/A') AS "cid",
			data.status,
			COUNT(data.id)
		FROM (
		SELECT DISTINCT get_pbc_request_status(req.*) AS status, req.id, acc.account_name
		FROM document_requests AS req
		LEFT JOIN request_control_link AS lnk
			ON lnk.request_id = req.id
		LEFT JOIN process_flow_control_node AS cn
			ON cn.control_id = lnk.control_id
		LEFT JOIN node_gl_link AS ngl
			ON ngl.node_id = cn.node_id
		LEFT JOIN general_ledger_accounts AS acc
			ON acc.id = ngl.gl_account_id
		WHERE req.org_id = $1
			AND ngl.gl_account_id IS NOT NULL
			AND %s
		) AS data
		GROUP BY data.status, cid
		ORDER BY cid
	`, orgId, filter)
}

func GetSystemCategoryPbcAnalytics(orgId int32, filter core.DocRequestFilterData) ([]PbcCategoryAnalyticsResult, error) {
	return GetGenericCategoryPbcAnalytics(`
		SELECT
			COALESCE(data.name, 'N/A') AS "cid",
			data.status,
			COUNT(data.id)
		FROM (
		SELECT DISTINCT get_pbc_request_status(req.*) AS status, req.id, sys.name
		FROM document_requests AS req
		LEFT JOIN request_control_link AS lnk
			ON lnk.request_id = req.id
		LEFT JOIN process_flow_control_node AS cn
			ON cn.control_id = lnk.control_id
		LEFT JOIN node_system_link AS nsl
			ON nsl.node_id = cn.node_id
		LEFT JOIN systems AS sys
			ON sys.id = nsl.system_id
		WHERE req.org_id = $1
			AND nsl.system_id IS NOT NULL
			AND %s
		) AS data
		GROUP BY data.status, cid
		ORDER BY cid
	`, orgId, filter)

}

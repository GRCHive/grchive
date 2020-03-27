package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"strconv"
	"strings"
)

func EditRisk(risk *core.Risk, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
		UPDATE process_flow_risks
		SET name = :name, description = :description, identifier = :identifier
		WHERE id = :id
	`, risk)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func DeleteRisks(nodeId int64, riskIds []int64, global bool, orgId int32, role *core.Role) error {
	if len(riskIds) == 0 {
		return nil
	}

	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	if nodeId != -1 {
		riskIdQuery := make([]string, 0)
		riskIdParams := append(make([]interface{}, 0), nodeId)
		for idx, id := range riskIds {
			riskIdParams = append(riskIdParams, id)
			riskIdQuery = append(riskIdQuery, fmt.Sprintf("$%d", idx+2))
		}

		_, err = tx.Exec(fmt.Sprintf(`
			DELETE FROM process_flow_risk_node
			WHERE node_id = $1
				AND risk_id IN (%s)
		`, strings.Join(riskIdQuery, ",")), riskIdParams...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if global {
		for _, id := range riskIds {
			_, err = tx.Exec(`
				DELETE FROM process_flow_risks
				WHERE id = $1
					AND org_id = $2
			`, id, orgId)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func AddRisksToNodeWithTx(riskIds []int64, nodeId int64, tx *sqlx.Tx, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	for _, id := range riskIds {
		_, err := tx.Exec(`
			INSERT INTO process_flow_risk_node (risk_id, node_id)
			VALUES ($1, $2)
		`, id, nodeId)

		if err != nil {
			return err
		}
	}
	return nil
}

func AddRisksToNode(riskIds []int64, nodeId int64, role *core.Role) error {
	tx := dbConn.MustBegin()
	err := AddRisksToNodeWithTx(riskIds, nodeId, tx, role)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func InsertNewRisk(risk *core.Risk, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	var err error

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO process_flow_risks (name, identifier, description, org_id)
		VALUES (:name, :identifier, :description, :org_id)
		RETURNING id
	`, risk)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&risk.Id)
	if err != nil {
		rows.Close()
		tx.Rollback()
		return err
	}
	rows.Close()

	err = tx.Commit()
	return err
}

func findAllRisksFromDbHelper(role *core.Role, stmt *sqlx.Stmt, args ...interface{}) ([]*core.Risk, error) {
	risks := []*core.Risk{}
	rows, err := stmt.Queryx(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := core.Risk{}
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}
		risks = append(risks, &r)
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, r := range risks {
		err = LogAuditSelectWithTx(r.OrgId, core.ResourceIdRisk, strconv.FormatInt(r.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return risks, tx.Commit()
}

func FindAllRisksForProcessFlow(flowId int64, role *core.Role) ([]*core.Risk, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	stmt, err := dbConn.Preparex(`
		SELECT 
			risk.*
		FROM process_flow_risks as risk
		INNER JOIN process_flow_risk_node AS risknode
			ON risknode.risk_id = risk.id
		INNER JOIN process_flow_nodes AS node
			ON risknode.node_id = node.id
		WHERE node.process_flow_id = $1
		GROUP BY risk.id
		ORDER BY risk.name ASC
	`)
	if err != nil {
		return nil, err
	}

	return findAllRisksFromDbHelper(role, stmt, flowId)
}

func FindAllRiskForOrganization(org *core.Organization, filter core.RiskFilterData, role *core.Role) ([]*core.Risk, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	sql := fmt.Sprintf(`
		SELECT 
			risk.*
		FROM process_flow_risks as risk
		LEFT JOIN process_flow_risk_control AS riskcontrol
			ON riskcontrol.risk_id = risk.id
		WHERE risk.org_id = $1
		GROUP BY risk.id
		HAVING
			%s
		ORDER BY risk.name ASC
	`,
		buildNumericFilter("COUNT(riskcontrol.control_id)", filter.NumControls),
	)
	stmt, err := dbConn.Preparex(sql)
	if err != nil {
		return nil, err
	}

	return findAllRisksFromDbHelper(role, stmt, org.Id)
}

func FindRisk(riskId int64, role *core.Role) (*core.Risk, error) {
	if !role.Permissions.HasAccess(core.ResourceProcessFlows, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}
	risk := core.Risk{}
	err := dbConn.Get(&risk, `
		SELECT *
		FROM process_flow_risks AS risk
		WHERE id = $1
	`, riskId)

	if err != nil {
		return nil, err
	}

	return &risk, LogAuditSelect(risk.OrgId, core.ResourceIdRisk, strconv.FormatInt(risk.Id, 10), role)
}

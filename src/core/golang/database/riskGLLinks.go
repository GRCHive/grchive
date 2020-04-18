package database

import (
	"gitlab.com/grchive/grchive/core"
)

func FindGeneralLedgerAccountsLinkedToRisk(riskId int64, orgId int32, role *core.Role) ([]*core.GeneralLedgerAccount, error) {
	if !role.Permissions.HasAccess(core.ResourceRisks, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	accounts := make([]*core.GeneralLedgerAccount, 0)
	err := dbConn.Select(&accounts, `
		SELECT DISTINCT acc.*
		FROM general_ledger_accounts AS acc
		INNER JOIN node_gl_link AS ngl
			ON ngl.gl_account_id = acc.id
		INNER JOIN process_flow_risk_node AS rn
			ON rn.node_id = ngl.node_id
		INNER JOIN process_flow_risks AS risk
			ON risk.id = rn.risk_id
		WHERE risk.id = $1 AND risk.org_id = $2
	`, riskId, orgId)

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func FindRisksLinkedToGeneralLedgerAccount(accountId int64, orgId int32, role *core.Role) ([]*core.Risk, error) {
	if !role.Permissions.HasAccess(core.ResourceRisks, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	risks := make([]*core.Risk, 0)
	err := dbConn.Select(&risks, `
		SELECT DISTINCT risk.*
		FROM process_flow_risks AS risk
		INNER JOIN process_flow_risk_node AS rn
			ON rn.risk_id = risk.id
		INNER JOIN node_gl_link AS ngl
			ON ngl.node_id = rn.node_id
		INNER JOIN general_ledger_accounts AS acc
			ON acc.id = ngl.gl_account_id
		WHERE acc.id = $1
			AND acc.org_id = $2
			AND risk.org_id = $2
	`, accountId, orgId)

	if err != nil {
		return nil, err
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	return risks, tx.Commit()
}

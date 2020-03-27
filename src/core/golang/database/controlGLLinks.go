package database

import (
	"gitlab.com/grchive/grchive/core"
	"strconv"
)

func FindGeneralLedgerAccountsLinkedToControl(controlId int64, orgId int32, role *core.Role) ([]*core.GeneralLedgerAccount, error) {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	accounts := make([]*core.GeneralLedgerAccount, 0)
	err := dbConn.Select(&accounts, `
		SELECT DISTINCT acc.*
		FROM general_ledger_accounts AS acc
		INNER JOIN node_gl_link AS ngl
			ON ngl.gl_account_id = acc.id
		INNER JOIN process_flow_control_node AS cn
			ON cn.node_id = ngl.node_id
		INNER JOIN process_flow_controls AS control
			ON control.id = cn.control_id
		WHERE control.id = $1 AND control.org_id = $2
	`, controlId, orgId)

	if err != nil {
		return nil, err
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, a := range accounts {
		err = LogAuditSelectWithTx(orgId, core.ResourceIdGLAcc, strconv.FormatInt(a.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return accounts, tx.Commit()
}

func FindControlsLinkedToGeneralLedgerAccount(accountId int64, orgId int32, role *core.Role) ([]*core.Control, error) {
	if !role.Permissions.HasAccess(core.ResourceControls, core.AccessView) ||
		!role.Permissions.HasAccess(core.ResourceGeneralLedger, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	controls := make([]*core.Control, 0)
	err := dbConn.Select(&controls, `
		SELECT DISTINCT control.*
		FROM process_flow_controls AS control
		INNER JOIN process_flow_control_node AS cn
			ON cn.control_id = control.id
		INNER JOIN node_gl_link AS ngl
			ON ngl.node_id = cn.node_id
		INNER JOIN general_ledger_accounts AS acc
			ON acc.id = ngl.gl_account_id
		WHERE acc.id = $1
			AND acc.org_id = $2
			AND control.org_id = $2
	`, accountId, orgId)

	if err != nil {
		return nil, err
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for _, c := range controls {
		err = LogAuditSelectWithTx(orgId, core.ResourceIdControl, strconv.FormatInt(c.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return controls, tx.Commit()
}

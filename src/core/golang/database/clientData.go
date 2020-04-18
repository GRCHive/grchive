package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"gitlab.com/grchive/grchive/core"
	"strconv"
)

func NewClientDataWithTx(data *core.ClientData, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceClientData, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO client_data (org_id, name, description)
		VALUES (:org_id, :name, :description)
		RETURNING id
	`, data)

	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()
	err = rows.Scan(&data.Id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateClientDataWithTx(data *core.ClientData, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceClientData, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
		UPDATE client_data
		SET name = :name,
			description = :description
		WHERE id = :id
			AND org_id = :org_id
	`, data)

	return err
}

func LinkClientDataToSourceWithTx(
	dataId int64,
	sourceId core.SourceId,
	sourceTarget map[string]interface{},
	orgId int32,
	role *core.Role,
	tx *sqlx.Tx,
) error {
	if !role.Permissions.HasAccess(core.ResourceClientData, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		return err
	}

	rawTarget, err := json.Marshal(sourceTarget)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO client_data_source_link (org_id, data_id, source_id, source_target)
		VALUES ($4, $1, $2, $3)
		ON CONFLICT (data_id) DO UPDATE SET
			source_id = EXCLUDED.source_id,
			source_target = EXCLUDED.source_target
	`, dataId, sourceId, string(rawTarget), orgId)
	return err
}

func getClientDataHelper(role *core.Role, conditions string, args ...interface{}) ([]*core.FullClientDataWithLink, error) {
	if !role.Permissions.HasAccess(core.ResourceClientData, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	data := make([]*core.FullClientDataWithLink, 0)
	rows, err := dbConn.Queryx(fmt.Sprintf(`
		SELECT
			data.id AS "data.id",
			data.org_id AS "data.org_id",
			data.name AS "data.name",
			data.description AS "data.description",
			link.org_id AS "link.org_id",
			link.data_id AS "link.data_id",
			link.source_id AS "link.source_id",
			link.source_target AS "link.source_target"
		FROM client_data AS data
		INNER JOIN client_data_source_link AS link
			ON link.data_id = data.id
		%s
	`, conditions), args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		newData := core.FullClientDataWithLink{}
		jsonData := types.JSONText{}
		err := rows.Scan(
			&newData.Data.Id,
			&newData.Data.OrgId,
			&newData.Data.Name,
			&newData.Data.Description,
			&newData.Link.OrgId,
			&newData.Link.DataId,
			&newData.Link.SourceId,
			&jsonData,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		err = jsonData.Unmarshal(&newData.Link.SourceTarget)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		data = append(data, &newData)

		err = LogAuditSelectWithTx(newData.Data.OrgId, core.ResourceIdClientData, strconv.FormatInt(newData.Data.Id, 10), role, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return data, tx.Commit()
}

func AllClientDataForOrganization(orgId int32, role *core.Role) ([]*core.FullClientDataWithLink, error) {
	return getClientDataHelper(role, `
		WHERE data.org_id = $1
	`, orgId)
}

func GetClientDataFromId(dataId int64, orgId int32, role *core.Role) (*core.FullClientDataWithLink, error) {
	data, err := getClientDataHelper(role, `
		WHERE data.id = $1 AND data.org_id = $2
	`, dataId, orgId)

	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("No client data found.")
	}

	return data[0], nil
}

func DeleteClientData(dataId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceClientData, core.AccessManage) {
		return core.ErrorUnauthorized
	}
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		DELETE FROM client_data
		WHERE id = $1 AND org_id = $2
	`, dataId, orgId)
	return tx.Commit()
}

func AllDataSourceOptions() ([]*core.DataSourceOption, error) {
	data := make([]*core.DataSourceOption, 0)
	err := dbConn.Select(&data, `
		SELECT *
		FROM data_source_options
	`)
	return data, err
}

func GetClientDataForCode(codeId int64, orgId int32, role *core.Role) (*core.ClientData, error) {
	if !role.Permissions.HasAccess(core.ResourceClientData, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	rows, err := dbConn.Queryx(`
		SELECT data.*
		FROM client_data AS data
		INNER JOIN code_to_client_data_link AS link
			ON link.data_id = data.id
		WHERE link.code_id = $1 AND link.org_id = $2
	`, codeId, orgId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	data := core.ClientData{}
	err = rows.StructScan(&data)
	return &data, err
}

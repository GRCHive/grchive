package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/sap_api"
)

func NewSapErpIntegrationWithTx(tx *sqlx.Tx, integrationId int64, setup sap_api.SapRfcConnectionOptions) error {
	_, err := tx.NamedExec(`
		INSERT INTO sap_erp_integration_info (integration_id, client, sysnr, host, real_hostname, username, password)
		VALUES (
			:integration_id,
			:setup.client,
			:setup.sysnr,
			:setup.host,
			:setup.real_hostname,
			:setup.username,
			:setup.password
		)
	`, struct {
		IntegrationId int64                           `db:"integration_id"`
		Setup         sap_api.SapRfcConnectionOptions `db:"setup"`
	}{
		IntegrationId: integrationId,
		Setup:         setup,
	})
	return err
}

func GetSapErpIntegration(integrationId int64) (*sap_api.SapRfcConnectionOptions, error) {
	connection := sap_api.SapRfcConnectionOptions{}
	err := dbConn.Get(&connection, `
		SELECT
			client,
			sysnr,
			host,
			real_hostname,
			username,
			password
		FROM sap_erp_integration_info
		WHERE integration_id = $1
	`, integrationId)
	return &connection, err
}

func EditSapErpIntegrationWithTx(tx *sqlx.Tx, integrationId int64, setup sap_api.SapRfcConnectionOptions) error {
	_, err := tx.NamedExec(`
		UPDATE sap_erp_integration_info
		SET client = :setup.client,
			sysnr = :setup.sysnr,
			host = :setup.host,
			real_hostname = :setup.real_hostname,
			username = :setup.username,
			password = :setup.password
		WHERE integration_id = :integration_id
	`, struct {
		IntegrationId int64                           `db:"integration_id"`
		Setup         sap_api.SapRfcConnectionOptions `db:"setup"`
	}{
		IntegrationId: integrationId,
		Setup:         setup,
	})
	return err
}

func AllSapErpRfc(integrationId int64) ([]*core.SapErpRfc, error) {
	rfc := make([]*core.SapErpRfc, 0)
	err := dbConn.Select(&rfc, `
		SELECT *
		FROM sap_erp_rfc
		WHERE integration_id = $1
		ORDER BY function_name ASC
	`, integrationId)
	return rfc, err
}

func NewSapErpRfcWithTx(tx *sqlx.Tx, rfc *core.SapErpRfc) error {
	rows, err := tx.NamedQuery(`
		INSERT INTO sap_erp_rfc (integration_id, function_name)
		VALUES (:integration_id, :function_name)
		RETURNING id
	`, rfc)
	if err != nil {
		return err
	}
	defer rows.Close()
	rows.Next()
	return rows.Scan(&rfc.Id)
}

func GetSapErpRfc(id int64) (*core.SapErpRfc, error) {
	rfc := core.SapErpRfc{}
	err := dbConn.Get(&rfc, `
		SELECT *
		FROM sap_erp_rfc
		WHERE id = $1
	`, id)
	return &rfc, err
}

func DeleteSapErpRfcWithTx(tx *sqlx.Tx, id int64) error {
	_, err := tx.Exec(`
		DELETE FROM sap_erp_rfc
		WHERE id = $1
	`, id)
	return err
}

func AllSapErpRfcVersions(rfcId int64) ([]*core.SapErpRfcVersion, error) {
	versions := make([]*core.SapErpRfcVersion, 0)
	err := dbConn.Select(&versions, `
		SELECT *
		FROM sap_erp_rfc_versions
		WHERE rfc_id = $1
	`, rfcId)
	return versions, err
}

func NewSapErpRfcVersionWithTx(tx *sqlx.Tx, rfcId int64) (*core.SapErpRfcVersion, error) {
	rows, err := tx.Queryx(`
		INSERT INTO sap_erp_rfc_versions (rfc_id, created_time)
		VALUES ($1, NOW())
		RETURNING *
	`, rfcId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	version := core.SapErpRfcVersion{}
	err = rows.StructScan(&version)
	return &version, err
}

func GetSapErpRfcVersion(versionId int64, rfcId int64) (*core.SapErpRfcVersion, error) {
	vers := core.SapErpRfcVersion{}
	err := dbConn.Get(&vers, `
		SELECT *
		FROM sap_erp_rfc_versions
		WHERE id = $1 AND rfc_id = $2
	`, versionId, rfcId)
	return &vers, err
}

func CompleteSapErpRfcVersionWithTx(tx *sqlx.Tx, versionId int64, success bool, data core.NullString) error {
	_, err := tx.Exec(`
		UPDATE sap_erp_rfc_versions
		SET success = $2,
			data = $3,
			finished_time = NOW()
		WHERE id = $1
	`, versionId, success, data)
	return err
}

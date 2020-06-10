package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/sap_api"
)

func NewIntegrationWithTx(tx *sqlx.Tx, integration *core.GenericIntegration) error {
	rows, err := tx.NamedQuery(`
		INSERT INTO integrations (org_id, type, name, description)
		VALUES (:org_id, :type, :name, :description)
		RETURNING id
	`, integration)
	if err != nil {
		return err
	}
	defer rows.Close()
	rows.Next()
	return rows.Scan(&integration.Id)
}

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

func LinkIntegrationWithSystemWithTx(tx *sqlx.Tx, integrationId int64, systemId int64, orgId int32) error {
	_, err := tx.Exec(`
		INSERT INTO integration_system_link (integration_id, system_id, org_id)
		VALUES ($1, $2, $3)
	`, integrationId, systemId, orgId)
	return err
}

func AllGenericIntegrationsForSystem(systemId int64) ([]*core.GenericIntegration, error) {
	integrations := make([]*core.GenericIntegration, 0)
	err := dbConn.Select(&integrations, `
		SELECT int.*
		FROM integrations AS int
		INNER JOIN integration_system_link AS isl
			ON isl.integration_id = int.id
		WHERE isl.system_id = $1
	`, systemId)
	return integrations, err
}

func GetGenericIntegration(integrationId int64) (*core.GenericIntegration, error) {
	integration := core.GenericIntegration{}
	err := dbConn.Get(&integration, `
		SELECT *
		FROM integrations
		WHERE id = $1
	`, integrationId)
	return &integration, err
}

func EditGenericIntegrationWithTx(tx *sqlx.Tx, integration *core.GenericIntegration) error {
	_, err := tx.NamedExec(`
		UPDATE integrations
		SET name = :name,
			description = :description
		WHERE id = :id
	`, integration)
	return err
}

func DeleteGenericIntegrationWithTx(tx *sqlx.Tx, integrationId int64) error {
	_, err := tx.Exec(`
		DELETE FROM integrations
		WHERE id = $1
	`, integrationId)
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

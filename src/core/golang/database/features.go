package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

// Returns an error if the feature was not enabled successfully either due to
// a database error or if the feature was already enabled.
func EnableFeatureForOrganizationWithTx(featureId core.FeatureId, orgId int32, role *core.Role, tx *sqlx.Tx) error {
	if !role.RoleMetadata.IsAdmin {
		return core.ErrorUnauthorized
	}

	_, err := tx.Exec(`
		INSERT INTO organization_enabled_features (org_id, feature_id, requested)
		VALUES ($2, $1, NOW())
	`, featureId, orgId)
	return err
}

func EnableFeatureForOrganization(featureId core.FeatureId, orgId int32, role *core.Role) error {
	tx := CreateTx()
	err := EnableFeatureForOrganizationWithTx(featureId, orgId, role, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func MarkFeatureAsFulfilled(featureId core.FeatureId, orgId int32, role *core.Role) error {
	if !role.RoleMetadata.IsAdmin {
		return core.ErrorUnauthorized
	}

	tx := CreateTx()
	_, err := tx.Exec(`
		UPDATE organization_enabled_features
		SET fulfilled = NOW()
		WHERE feature_id = $1 AND org_id = $2
	`, featureId, orgId)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func IsFeatureEnabledForOrganization(featureId core.FeatureId, orgId int32) (bool, bool, error) {
	rows, err := dbConn.Queryx(`
		SELECT 
			CASE
				WHEN fulfilled IS NULL THEN FALSE
				ELSE TRUE
			END AS enabled,
			CASE
				WHEN requested IS NULL THEN FALSE
				ELSE TRUE
			END AS pending
		FROM organization_enabled_features
		WHERE feature_id = $1 AND org_id = $2
	`, int64(featureId), orgId)

	if err != nil {
		return false, false, err
	}

	defer rows.Close()

	if !rows.Next() {
		return false, false, nil
	}

	type Result struct {
		Enabled bool `db:"enabled"`
		Pending bool `db:"pending"`
	}
	result := Result{}
	err = rows.StructScan(&result)
	if err != nil {
		return false, false, nil
	}

	return result.Enabled, result.Pending, err
}

func GetFeatureName(featureId core.FeatureId) (string, error) {
	name := ""
	err := dbConn.Get(&name, `
		SELECT name
		FROM available_features
		WHERE id = $1
	`, int64(featureId))
	return name, err
}

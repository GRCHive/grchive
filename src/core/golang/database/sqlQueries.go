package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/grchive/grchive/core"
)

func GetAllSqlQueryMetadataForDb(dbId int64, orgId int32, role *core.Role) ([]*core.DbSqlQueryMetadata, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSqlQuery, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	data := make([]*core.DbSqlQueryMetadata, 0)
	err := dbConn.Select(&data, `
		SELECT *
		FROM database_sql_metadata
		WHERE db_id = $1 AND org_id = $2
		ORDER BY name ASC
	`, dbId, orgId)
	return data, err
}

func GetAllSqlQueryVersionsForMetadata(metadataId int64, orgId int32, role *core.Role) ([]*core.DbSqlQuery, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSqlQuery, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	data := make([]*core.DbSqlQuery, 0)
	err := dbConn.Select(&data, `
		SELECT *
		FROM database_sql_queries
		WHERE metadata_id = $1 AND org_id = $2
		ORDER BY version_number DESC
	`, metadataId, orgId)
	return data, err
}

func GetSqlQueryMetadata(queryId int64, orgId int32, role *core.Role) (*core.DbSqlQueryMetadata, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSqlQuery, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	return nil, nil
}

func GetSqlQueryFromMetadataId(metadataId int64, orgId int32, version int32, role *core.Role) (*core.DbSqlQuery, error) {
	if !role.Permissions.HasAccess(core.ResourceDbSqlQuery, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	return nil, nil
}

func CreateSqlQueryMetadataWithTx(metadata *core.DbSqlQueryMetadata, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDbSqlQuery, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO database_sql_metadata (db_id, org_id, name, description)
		VALUES (:db_id, :org_id, :name, :description)
		RETURNING id
	`, metadata)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&metadata.Id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSqlQueryMetadataWithTx(metadata *core.DbSqlQueryMetadata, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDbSqlQuery, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		UPDATE database_sql_metadata 
		SET name = :name,
		 	description = :description
		WHERE id = :id
			AND org_id = :org_id
		RETURNING *
	`, metadata)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	err = rows.StructScan(metadata)
	if err != nil {
		return err
	}
	return nil
}

func CreateSqlQueryWithTx(query *core.DbSqlQuery, role *core.Role, tx *sqlx.Tx) error {
	if !role.Permissions.HasAccess(core.ResourceDbSqlQuery, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO database_sql_queries (metadata_id, version_number, upload_time, upload_user_id, org_id, query)
		SELECT :metadata_id, COALESCE(MAX(version_number), 0) + 1, :upload_time, :upload_user_id, :org_id, :query
		FROM database_sql_queries
		WHERE metadata_id = :metadata_id
			AND org_id = :org_id
		RETURNING id, version_number
	`, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&query.Id, &query.Version)
	if err != nil {
		return err
	}
	return nil
}

func DeleteSqlQuery(metadataId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDbSqlQuery, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	_, err := tx.Exec(`
		DELETE FROM database_sql_metadata 
		WHERE id = $1
			AND org_id = $2
	`, metadataId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

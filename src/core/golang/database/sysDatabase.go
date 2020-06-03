package database

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"gitlab.com/grchive/grchive/core"
	"strings"
)

func GetAllSupportedDatabaseTypes(role *core.Role) ([]*core.DatabaseType, error) {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	types := make([]*core.DatabaseType, 0)
	err := dbConn.Select(&types, `
		SELECT *
		FROM supported_databases
	`)
	return types, err
}

func InsertNewDatabase(db *core.Database, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	rows, err := tx.NamedQuery(`
		INSERT INTO database_resources(name, org_id, type_id, other_type, version)
		VALUES (:name, :org_id, :type_id, :other_type, :version)
		RETURNING id
	`, db)

	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&db.Id)
	if err != nil {
		rows.Close()
		tx.Rollback()
		return err
	}
	rows.Close()

	return tx.Commit()
}

func GetAllDatabasesForOrg(orgId int32, filter core.DatabaseFilterData, role *core.Role) ([]*core.Database, error) {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	dbs := make([]*core.Database, 0)
	err := dbConn.Select(&dbs, fmt.Sprintf(`
		SELECT *
		FROM database_resources
		WHERE org_id = $1 AND
			%s
	`,
		buildNumericFilter("type_id", filter.Type),
	), orgId)

	if err != nil {
		return nil, err
	}

	return dbs, nil
}

func GetAllDatabasesForOrgWithDeployment(orgId int32, deploymentType int32, filter core.DatabaseFilterData, role *core.Role) ([]*core.Database, error) {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	dbs := make([]*core.Database, 0)
	err := dbConn.Select(&dbs, fmt.Sprintf(`
		SELECT db.*
		FROM database_resources AS db
		INNER JOIN deployment_db_link AS link
			ON link.db_id = db.id
		INNER JOIN deployments AS dp
			ON dp.id = link.deployment_id
		WHERE db.org_id = $1
			AND dp.deployment_type = $2
			AND %s
	`,
		buildNumericFilter("db.type_id", filter.Type),
	), orgId, deploymentType)
	if err != nil {
		return nil, err
	}

	return dbs, nil
}

func GetDb(dbId int64, orgId int32, role *core.Role) (*core.Database, error) {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	db := core.Database{}
	err := dbConn.Get(&db, `
		SELECT *
		FROM database_resources
		WHERE id = $1
			AND org_id = $2
	`, dbId, orgId)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

func GetDbType(dbId int64, orgId int32, role *core.Role) (*core.DatabaseType, error) {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	dbType := core.DatabaseType{}
	err := dbConn.Get(&dbType, `
		SELECT sd.*
		FROM supported_databases AS sd
		INNER JOIN database_resources AS dr
			ON sd.id = dr.type_id
		WHERE dr.id = $1
			AND dr.org_id = $2
	`, dbId, orgId)
	return &dbType, err
}

func EditDb(db *core.Database, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessEdit) {
		return core.ErrorUnauthorized
	}
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
		UPDATE database_resources
		SET name = :name,
			type_id = :type_id,
			other_type = :other_type,
			version = :version
		WHERE id = :id
			AND org_id = :org_id
	`, db)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func DeleteDb(dbId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessManage) {
		return core.ErrorUnauthorized
	}
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM database_resources
		WHERE id = $1
			AND org_id = $2
	`, dbId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func InsertNewDatabaseConnection(conn *core.DatabaseConnection, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDbConnections, core.AccessManage) {
		return core.ErrorUnauthorized
	}

	connectionParameters, err := json.Marshal(conn.Parameters)
	if err != nil {
		return err
	}

	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	rows, err := tx.Queryx(`
		INSERT INTO database_connection_info (
			db_id,
			org_id,
			host, 
			port,
			dbname,
			parameters,
			username,
			password,
			salt)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9)
		RETURNING id`,
		conn.DbId,
		conn.OrgId,
		conn.Host,
		conn.Port,
		conn.DbName,
		string(connectionParameters),
		conn.Username,
		conn.Password,
		conn.Salt,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	rows.Next()
	err = rows.Scan(&conn.Id)
	if err != nil {
		rows.Close()
		tx.Rollback()
		return err
	}
	rows.Close()
	return tx.Commit()
}

// Returns nil, nil if no connection is found.
func FindDatabaseConnectionForDatabase(dbId int64, orgId int32, role *core.Role) (*core.DatabaseConnection, error) {
	if !role.Permissions.HasAccess(core.ResourceDbConnections, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	rows, err := dbConn.Queryx(`
		SELECT *
		FROM database_connection_info
		WHERE db_id = $1
			AND org_id = $2
	`, dbId, orgId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	conn := core.DatabaseConnection{}
	rawParameters := types.JSONText{}

	err = rows.Scan(&conn.Id,
		&conn.DbId,
		&conn.OrgId,
		&conn.Username,
		&conn.Password,
		&conn.Salt,
		&conn.Host,
		&conn.Port,
		&rawParameters,
		&conn.DbName,
	)
	if err != nil {
		return nil, err
	}

	err = rawParameters.Unmarshal(&conn.Parameters)
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

func DeleteDatabaseConnection(connId int64, dbId int64, orgId int32, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceDbConnections, core.AccessManage) {
		return core.ErrorUnauthorized
	}
	tx, err := CreateAuditTrailTx(role)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM database_connection_info
		WHERE id = $1
			AND db_id = $2
			AND org_id = $3
	`, connId, dbId, orgId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func FindDbIdsForSystem(sysId int64, orgId int32, role *core.Role) ([]int64, error) {
	if !role.Permissions.HasAccess(core.ResourceDatabases, core.AccessView) {
		return nil, core.ErrorUnauthorized
	}

	ids := make([]int64, 0)
	err := dbConn.Select(&ids, `
		SELECT db_id
		FROM database_system_link
		WHERE system_id = $1
			AND org_id = $2
	`, sysId, orgId)
	return ids, err
}

func LinkSystemsToDatabase(dbId int64, orgId int32, sysIds []int64, role *core.Role) error {
	if !role.Permissions.HasAccess(core.ResourceSystems, core.AccessEdit) ||
		!role.Permissions.HasAccess(core.ResourceDatabases, core.AccessEdit) {
		return core.ErrorUnauthorized
	}

	tx := dbConn.MustBegin()
	for _, sysId := range sysIds {
		_, err := tx.Exec(`
			INSERT INTO database_system_link (db_id, org_id, system_id)
			VALUES ($1, $2, $3)
		`, dbId, orgId, sysId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func GetDatabaseSettings(dbId int64) (*core.DatabaseSettings, error) {
	rows, err := dbConn.Queryx(`
		SELECT
			ds.db_id,
			ds.org_id,
			ds.auto_refresh_task,
			st.id IS NOT NULL AS "auto_refresh_enabled",
			rt.rrule AS "auto_refresh_rrule",
			ARRAY_TO_JSON(ARRAY_REMOVE(ARRAY_AGG(u.*), null)) AS "on_schema_change_notify_users"
		FROM database_settings AS ds
		LEFT JOIN scheduled_tasks AS st
			ON st.id = ds.auto_refresh_task
		LEFT JOIN recurring_tasks AS rt
			ON rt.event_id = st.id
		LEFT JOIN db_refresh_diff_message_recipients AS dmr
			ON dmr.db_id = ds.db_id
		LEFT JOIN users AS u
			ON u.id = dmr.user_id
		WHERE ds.db_id = $1
		GROUP BY ds.db_id, ds.org_id, ds.auto_refresh_task, st.id, rt.rrule
	`, dbId)

	if err != nil {
		return nil, err
	}

	rows.Next()

	settings := core.DatabaseSettings{}

	rawOnSchemaChangeNotifyUsers := []uint8{}

	err = rows.Scan(
		&settings.DbId,
		&settings.OrgId,
		&settings.AutoRefreshTaskId,
		&settings.AutoRefreshEnabled,
		&settings.AutoRefreshRRule,
		&rawOnSchemaChangeNotifyUsers,
	)

	if err != nil {
		return nil, err
	}

	settings.OnSchemaChangeNotifyUsers, err = readUserArray(rawOnSchemaChangeNotifyUsers)
	if err != nil {
		return nil, err
	}

	return &settings, err
}

func LinkDatabaseSettingToScheduledTaskWithTx(tx *sqlx.Tx, dbId int64, taskId int64) error {
	_, err := tx.Exec(`
		UPDATE database_settings
		SET auto_refresh_task = $2
		WHERE db_id = $1
	`, dbId, taskId)
	return err
}

func SyncSchemaChangeNotifyUsersWithTx(tx *sqlx.Tx, dbId int64, orgId int32, userIds []int64) error {
	if len(userIds) > 0 {
		// Do the sync in two steps:
		// 1) Delete user entries that don't match.
		// 2) Add users that don't exist.
		query, args, err := sqlx.In(`
			DELETE FROM db_refresh_diff_message_recipients
			WHERE db_id = ? 
				AND user_id NOT IN (?)
		`, dbId, userIds)

		if err != nil {
			return err
		}

		query = tx.Rebind(query)
		_, err = tx.Exec(query, args...)
		if err != nil {
			return err
		}

		queryBuilder := strings.Builder{}
		queryBuilder.WriteString("INSERT INTO db_refresh_diff_message_recipients (db_id, org_id, user_id) VALUES ")

		for idx, uid := range userIds {
			queryBuilder.WriteString(fmt.Sprintf("(%d, %d, %d)", dbId, orgId, uid))
			if idx != len(userIds)-1 {
				queryBuilder.WriteString(",")
			}
		}

		queryBuilder.WriteString("ON CONFLICT (db_id, user_id) DO NOTHING")
		_, err = tx.Exec(queryBuilder.String())
		return err
	} else {
		_, err := tx.Exec(`
			DELETE FROM db_refresh_diff_message_recipients
			WHERE db_id = $1
		`, dbId)
		return err
	}
}

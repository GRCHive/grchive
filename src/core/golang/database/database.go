package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gitlab.com/grchive/grchive/core"
	"time"
)

var dbConn *sqlx.DB

func Init() {
	dbConn = sqlx.MustConnect("postgres", core.EnvConfig.DatabaseConnString)
	dbConn.SetMaxOpenConns(10)
	dbConn.SetMaxIdleConns(5)
	dbConn.SetConnMaxLifetime(5 * time.Minute)
}

func CreateTx() *sqlx.Tx {
	return dbConn.MustBegin()
}

func UpgradeTxToAudit(tx *sqlx.Tx, role *core.Role) error {
	if role.UserId != -1 {
		_, err := tx.Exec(`
			SELECT set_current_role_for_user_id($1)
		`, role.UserId)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateAuditTrailTx(role *core.Role) (*sqlx.Tx, error) {
	tx := CreateTx()
	err := UpgradeTxToAudit(tx, role)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return tx, nil
}

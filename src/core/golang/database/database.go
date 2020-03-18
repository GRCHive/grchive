package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

type ListenHandler func(data string) error

const (
	NotifyChannelControlOwner       string = "controlowner"
	NotifyChannelDocRequestAssignee        = "docrequestassignee"
	NotifyChannelDocRequestStatus          = "docrequeststatus"
	NotifyChannelSqlRequestAssignee        = "sqlrequestassignee"
	NotifyChannelSqlRequestStatus          = "sqlrequeststatus"
)

func InitListeners(config map[string]ListenHandler) {
	minDuration, err := time.ParseDuration("5s")
	if err != nil {
		core.Error("Failed to parse min duration: " + err.Error())
	}

	maxDuration, err := time.ParseDuration("5m")
	if err != nil {
		core.Error("Failed to parse max duration: " + err.Error())
	}

	listener := pq.NewListener(
		core.EnvConfig.DatabaseConnString,
		minDuration,
		maxDuration,
		nil)

	for k := range config {
		err = listener.Listen(k)
		if err != nil {
			core.Error("Failed to listen to channel: " + err.Error())
		}
	}

	go func() {
		defer listener.Close()
		for {
			n := <-listener.Notify
			handler := config[n.Channel]
			err = handler(n.Extra)
			if err != nil {
				core.Warning("Failed to handle DB notify: " + err.Error())
			}
		}
	}()
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

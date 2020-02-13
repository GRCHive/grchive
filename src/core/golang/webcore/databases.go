package webcore

import (
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
)

func CreateNewDatabaseRefresh(dbId int64, orgId int32, role *core.Role) (*core.DbRefresh, error) {
	refresh, err := database.CreateNewDatabaseRefresh(dbId, orgId, role)
	if err != nil {
		return nil, err
	}

	DefaultRabbitMQ.SendMessage(PublishMessage{
		Exchange: DEFAULT_EXCHANGE,
		Queue:    DATABASE_REFRESH_QUEUE,
		Body: DatabaseRefreshMessage{
			Refresh: *refresh,
		},
	})

	return refresh, nil
}

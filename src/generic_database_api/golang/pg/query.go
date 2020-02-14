package pg_api

import (
	"gitlab.com/grchive/grchive/db_api/utility"
)

func (pg *PgDriver) RunQuery(query string) (*utility.SqlQueryResult, error) {
	rows, err := pg.connection.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return utility.CreateSqlQueryResultFromRows(rows)
}
